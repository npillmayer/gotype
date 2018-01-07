package gallery

/*
----------------------------------------------------------------------

BSD License
Copyright (c) 2017, Norbert Pillmayer <norbert@pillmayer.com>

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:
1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.
3. Neither the name of Norbert Pillmayer nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----------------------------------------------------------------------

 * This is the implementation of an interpreter for 'Gallery', a DSL for
 * placing frames on pages. See the Wiki for further details.
 *
 * The implementation is tightly coupled to the ANTLR V4 parser generator.
 * ANTLR is a great tool and I see no use in being independent from it.

*/

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	arithm "github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/syntax/pmmpost/corelang"
	"github.com/npillmayer/gotype/syntax/gallery/grammar"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	dec "github.com/shopspring/decimal"
)

// === Interpreter ===========================================================

/* We use AST-driven interpretation to execute the Gallery input program. Input
 * is more or less a list of statements and function definitions.
 *
 * We will annotate the AST with scope-information, holding the symbols of
 * dynamic scopes. Scopes stem from either:
 * - function definitions (macros in MetaFont: def and vardef)
 * - compound statements, i.e. groups (begingroup ... endgroup)
 *
 * The interpreter relies on the scopes and definitions constructed earlier.
 * It manages a memory frame stack to track the calling sequence of functions
 * and groups.
 *
 * So the overall picture looks like this:
 * 1. ANTLR V4 constructs an AST for us.
 * 2. We use a listener to walk the AST and execute the statements.
 *
 * Metafont, and therefore Gallery, is a dynamically scoped language. This means,
 * functions can access local variables from calling functions or groups.
 * Nevertheless we will find the definition of all variables (which are explicitly
 * defined) in a scope definition. This is mainly for type checking reasons and
 * due to the complex structure of MetaFont variable identifiers.
 */

/* Type GalleryInterpreter.
 * This is an umbrella object to hold together the various tools needed to
 * execute steps 1 to 3 from above. It will orchestrate and instrument the
 * tools and execute them in the correct order. Also, this object will hold
 * and respect parameters we pass to the interpreter, so we can alter the
 * behaviour in certain aspects.
 */
type GalleryInterpreter struct {
	ASTListener *ParseListener   // parse / AST listener
	runtime     *runtime.Runtime // runtime environment
}

/* Create a new Interpreter for "Poor Man's MetaPost". This is the top-level
 * object for this package.
 *
 * The interpreter manages an AST listener, a scope tree, a memory frame stack,
 * and stacks for numeric/pair and path expressions.
 *
 * Loads builtin symbols (variables and types).
 */
func NewGalleryInterpreter() *GalleryInterpreter {
	intp := &GalleryInterpreter{}
	intp.runtime = runtime.NewRuntimeEnvironment(variables.NewPMMPVarDecl)
	corelang.LoadBuiltinSymbols(intp.runtime.ScopeTree.Globals())         // load syms into new global scope
	intp.ASTListener = NewParseListener(intp.runtime.ScopeTree.Globals()) // listener for ANTLR
	intp.ASTListener.rt = intp.runtime
	return intp
}

/* Parse and interpret a statement list.
 */
func (intp *GalleryInterpreter) ParseStatements(input antlr.CharStream) {
	intp.ASTListener.ParseStatements(input)
}

// === AST driven parsing ====================================================

/* ANTLR will create an AST for us. We use a listener to attach to ANTLR's
 * AST walker. The listener manages scopes (with declarations) and memory frames
 * (with variable references and values).
 */
type ParseListener struct {
	*grammar.BaseGalleryListener // build on top of ANTLR's base 'class'
	statemParser                 *grammar.GalleryParser
	annotations                  map[interface{}]Annotation // node annotations
	expectingLvalue              bool                       // do not evaluate variable
	rt                           *runtime.Runtime           // runtime environment
}

/* We will annotate the AST. Functions and groups will get a scope, filled with
 * statically local variable definitions. Scope information is tracked using a
 * stack. Walking the AST results in scopes being pushed and popped
 * on and off the scope stack (dynamically forming a scope tree). Scope tree
 * scopes are always linking backward to their parent, with the global scope
 * being the root scope. Every scope holds a symbol table.
 *
 * Whenever we identify a scope, we fill it with the local symbols and then
 * attach it to the corresponding AST node. AST nodes will be either
 * - function definitions (macros in MetaFont: def and vardef)
 * - compound statements, i.e. groups (begingroup ... endgroup)
 */
type Annotation struct {
	scope *runtime.Scope
	text  string
}

/* Construct a new AST listener.
 */
func NewParseListener(globalScope *runtime.Scope) *ParseListener {
	pl := &ParseListener{} // no need to initialize base class
	pl.annotations = make(map[interface{}]Annotation)
	pl.annotate("global", globalScope, "")
	return pl
}

/* We use ANTLR V4 for parsing the statement grammar.
 */
func (pl *ParseListener) ParseStatements(input antlr.CharStream) {
	pl.LazyCreateParser(input)
	tree := pl.statemParser.Statementlist()
	sexpr := antlr.TreesStringTree(tree, nil, pl.statemParser)
	T.Debugf("### STATEMENT = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(pl, tree)
}

// We create our own type of error listener for the ANTLR parser
type TracingErrorListener struct {
	*antlr.DefaultErrorListener // use default as base class
}

/* Our error listener prints an error to the trace.
 */
func (c *TracingErrorListener) SyntaxError(r antlr.Recognizer, sym interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	//
	at := fmt.Sprintf("%s:%s", strconv.Itoa(line), strconv.Itoa(column))
	T.P("line", at).Errorf("%.44s", msg)
}

/* Helper function to create an ANTLR V4 parser. This function should cache
 * the parser for re-use, but currently I do not understand how to do this
 * in the Go version of ANTLR. According to a forum discussion, creating a
 * new parser every time seems to be the accepted mode of operation and should
 * not carry too much of a performance penalty.
 */
func (pl *ParseListener) LazyCreateParser(input antlr.CharStream) {
	// We let ANTLR to the heavy lifting.
	lexer := grammar.NewGalleryLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&TracingErrorListener{})
	stream := antlr.NewCommonTokenStream(lexer, 0)
	if pl.statemParser == nil {
		pl.statemParser = grammar.NewGalleryParser(stream)
		pl.statemParser.RemoveErrorListeners()
		pl.statemParser.AddErrorListener(&TracingErrorListener{})
		pl.statemParser.BuildParseTrees = true
	} else {
		pl.statemParser.SetInputStream(stream) // this should work
	}
}

/* Annotate an AST node, i.e., attach a scope information.
 */
func (pl *ParseListener) annotate(node interface{}, scope *runtime.Scope, text string) {
	pl.annotations[node] = Annotation{scope, text}
}

/* Get the annotation for an AST node.
 */
func (pl *ParseListener) getAnnotation(node interface{}) (*runtime.Scope, string) {
	if a, found := pl.annotations[node]; found {
		return a.scope, a.text
	}
	return nil, ""
}

/* Print out a summary of all the scopes and symbols collected up to now.
 */
func (pl *ParseListener) Summary() {
	pl.rt.ExprStack.Summary()
	T.Info("Summary of symbols:")
	for _, annot := range pl.annotations {
		scope := annot.scope
		if scope != nil {
			for _, sym := range scope.Symbols().Table {
				T.P("scope", scope.GetName()).Infof("%v", sym)
			}
		}
	}
}

// === Listener Callbacks for Statement Parsing ==============================

/* TODO: Improve error handling.
 * Currently just traces to T.
 */
func (pl *ParseListener) VisitErrorNode(node antlr.ErrorNode) {
	T.Errorf("parser error: %s", node.GetText())
}

/* Helper to trace terminal symbols. Just traces to T.
 */
func (pl *ParseListener) VisitTerminal(node antlr.TerminalNode) {
	//T.Debug("@@@ terminal: %s  = %v", node.GetText(), node.GetSymbol())
}



// --- Variable Handling -----------------------------------------------------

/* Internal method: Construct a valid variable reference string from parts on
 * the stack.
 *
 * Collect fragments of a variable reference, e.g. "x[k+1]r".
 * Subscripts should be found on the expression stack and inserted as
 * numeric constants, i.e. resulting in "x[5]r" (if k=4).
 *
 * Parameter t is the text of the variable ref literal, e.g. "x[k+1]r".
 * It is split by the parser into:
 *
 * . "x" -> TAG x
 * . subscript { k+1 }
 * . "r" -> TAG r
 */
func (pl *ParseListener) collectVarRefParts(t string, children []antlr.Tree) string {
	var vname bytes.Buffer
	for _, ch := range children {
		T.Debugf("collecting var ref part: %s", getCtxText(ch))
		if isTerminal(ch) { // just copy string parts to output
			T.Debugf("adding suffix verbatim: %s", getCtxText(ch))
			vname.WriteString(ch.(antlr.ParseTree).GetText())
		} else { // non-terminal is a subscript-expression
			subscript, ok := pl.rt.ExprStack.Pop() // take subscript from stack
			if !ok {
				T.P("var", t).Error("expected subscript on expression stack")
				T.P("var", t).Error("substituting 0 instead")
				vname.WriteString("[0]")
			} else {
				c, isconst := subscript.GetXPolyn().IsConstant()
				if !isconst { // we cannot handle unknown subscripts
					T.P("var", t).Error("subscript must be known numeric")
					T.P("var", t).Errorf("substituting 0 for %s",
						pl.rt.ExprStack.TraceString(subscript))
					vname.WriteString("[0]")
				} else {
					vname.WriteString("[")
					vname.WriteString(c.String())
					vname.WriteString("]")
				}
			}
		}
	}
	return vname.String()
}

/* A variable subscript.
 */
func (pl *ParseListener) ExitSubscript(ctx *grammar.SubscriptContext) {
	if ctx.DECIMALTOKEN() != nil {
		c, _ := dec.NewFromString(ctx.DECIMALTOKEN().GetText())
		pl.rt.ExprStack.PushConstant(c)
	}
}

/* Get or create a variable reference. To get the canonical representation of
 * the variable reference, we parse it and construct a small AST. This AST
 * is fed into getPMMPVarRefFromVarSyntax(). The resulting variable reference
 * struct is used to find the memory location of the variable reference.
 *
 * The reference lives in a memory frame, so we first locate it, then put
 * it on the expression stack. If the variable has a known value, we will
 * put the value onto the stack (otherwise the variable reference).
 */
func (pl *ParseListener) makeCanonicalAndResolve(v string, expectLvalue bool) *variables.PMMPVarRef {
	vtree := variables.ParseVariableFromString(v, &TracingErrorListener{})
	vref := variables.GetVarRefFromVarSyntax(vtree, pl.rt.ScopeTree)
	vref, _ = corelang.FindVariableReferenceInMemory(pl.rt, vref, true) // allocate if not found
	return vref
}

// --- Utilities -------------------------------------------------------------

/* TODO complete this.
 * Return scaled points for high level units (cm, mm, pt, in, ...)
 */
func unit2numeric(u string) dec.Decimal {
	switch u {
	case "in":
		return dec.NewFromFloat(0.01388888)
	}
	return arithm.ConstOne
}

/* Scale a numeric value by a unit.
 */
func scaleDimension(dimen dec.Decimal, unit string) dec.Decimal {
	u := unit2numeric(unit)
	return dimen.Mul(u)
}

/* I do not always understand ANTLR V4's Go runtime typing and tree semantics
 * (rather poorly documented), so I introduce some helpers. Some of these are
 * probably unnecessary for a better versed ANTLR Go user...
 */
func isTerminal(node antlr.Tree) bool {
	_, ok := node.GetPayload().(antlr.RuleNode)
	//fmt.Printf("node is terminal: %v\n", !ok)
	return !ok
}

/* I do not always understand ANTLR V4's Go runtime typing and tree semantics
 * (rather poorly documented), so I introduce some helpers. Some of these are
 * probably unnecessary for a better versed ANTLR Go user...
 */
func isTag(node antlr.Tree) bool {
	var ok bool
	var r antlr.RuleNode
	r, ok = node.GetPayload().(antlr.RuleNode)
	if ok {
		ctx := r.GetRuleContext()
		_, ok = ctx.(grammar.IAnytagContext)
	}
	//fmt.Printf("node is anytag: %v\n", ok)
	return ok
}

/* I do not always understand ANTLR V4's Go runtime typing and tree semantics
 * (rather poorly documented), so I introduce some helpers. Some of these are
 * probably unnecessary for a better versed ANTLR Go user...
 */
func getCtxText(ctx antlr.Tree) string {
	t := ctx.(antlr.ParseTree).GetText()
	return t
}
