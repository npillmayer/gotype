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
	"github.com/antlr/antlr4/runtime/Go/antlr"
	arithm "github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/syntax/corelang"
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
	ASTListener *GalleryParseListener // parse / AST listener
	runtime     *runtime.Runtime      // runtime environment
}

/* Create a new Interpreter for "Poor Man's MetaPost". This is the top-level
 * object for this package.
 *
 * The interpreter manages an AST listener, a scope tree, a memory frame stack,
 * and stacks for numeric/pair and path expressions.
 *
 * Loads builtin symbols (variables and types).
 */
func NewGalleryInterpreter(loadBuiltins bool) *GalleryInterpreter {
	intp := &GalleryInterpreter{}
	intp.runtime = runtime.NewRuntimeEnvironment(variables.NewPMMPVarDecl)
	if loadBuiltins {
		corelang.LoadBuiltinSymbols(intp.runtime) // load syms into global scope
	}
	intp.ASTListener = NewParseListener(intp.runtime) // listener for ANTLR
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
type GalleryParseListener struct {
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
func NewParseListener(rt *runtime.Runtime) *GalleryParseListener {
	pl := &GalleryParseListener{} // no need to initialize base class
	pl.rt = rt
	pl.annotations = make(map[interface{}]Annotation)
	pl.annotate("global", rt.ScopeTree.Globals(), "")
	return pl
}

/* We use ANTLR V4 for parsing the statement grammar.
 */
func (pl *GalleryParseListener) ParseStatements(input antlr.CharStream) {
	pl.LazyCreateParser(input)
	tree := pl.statemParser.Program()
	sexpr := antlr.TreesStringTree(tree, nil, pl.statemParser)
	T.Debugf("### program = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(pl, tree)
}

/* Helper function to create an ANTLR V4 parser. This function should cache
 * the parser for re-use, but currently I do not understand how to do this
 * in the Go version of ANTLR. According to a forum discussion, creating a
 * new parser every time seems to be the accepted mode of operation and should
 * not carry too much of a performance penalty.
 */
func (pl *GalleryParseListener) LazyCreateParser(input antlr.CharStream) {
	// We let ANTLR to the heavy lifting.
	lexer := grammar.NewGalleryLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&corelang.TracingErrorListener{})
	stream := antlr.NewCommonTokenStream(lexer, 0)
	if pl.statemParser == nil {
		pl.statemParser = grammar.NewGalleryParser(stream)
		pl.statemParser.RemoveErrorListeners()
		pl.statemParser.AddErrorListener(&corelang.TracingErrorListener{})
		pl.statemParser.BuildParseTrees = true
	} else {
		pl.statemParser.SetInputStream(stream) // this should work
	}
}

/* Annotate an AST node, i.e., attach a scope information.
 */
func (pl *GalleryParseListener) annotate(node interface{}, scope *runtime.Scope, text string) {
	pl.annotations[node] = Annotation{scope, text}
}

/* Get the annotation for an AST node.
 */
func (pl *GalleryParseListener) getAnnotation(node interface{}) (*runtime.Scope, string) {
	if a, found := pl.annotations[node]; found {
		return a.scope, a.text
	}
	return nil, ""
}

/* Print out a summary of all the scopes and symbols collected up to now.
 */
func (pl *GalleryParseListener) Summary() {
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

/* Print an error to the trace.
 */
func (pl *GalleryParseListener) VisitErrorNode(node antlr.ErrorNode) {
	T.Errorf("parser error: %s", node.GetText())
}

/* Helper to trace terminal symbols. Just traces to T.
func (pl *GalleryParseListener) VisitTerminal(node antlr.TerminalNode) {
	//T.Debug("@@@ terminal: %s  = %v", node.GetText(), node.GetSymbol())
}
*/

// --- Statements ------------------------------------------------------------

/* Start a new scope / compound statement: "begingroup".
 *
 * MetaFont uses dynamic scopes with the "begingroup ... endgroup" notation,
 * but in a unique way: users have to "save" variables explicitly
 * (The METAFONTBook, chapter 17).
 *
 * We produce a new scope for the scope tree and a new memory frame onto
 * the memory frame stack. The scope holds declarations for local variables.
 * The memory frame holds variable references (pointing to the the decls
 * in the scope).
 *
 * We reserve a symbol in the group's scope on "save" statements only.
 * Saved tags are of type tag/undefined initially.
 *
 * Notice: It is a bit of an overkill to use dynamic scopes and memory frames
 * in parallel, but I'll stick to this traditional approach for clarity.
 * Furthermore, this approach comes handy for the runtime-handling of
 * definitions (a.k.a. macros) and/or functions.
 */
func (pl *GalleryParseListener) EnterCompound(ctx *grammar.CompoundContext) {
	groupscope, _ := corelang.Begingroup(pl.rt, "compound-group")
	pl.annotate(ctx, groupscope, "") // Annotate the AST node with this scope
}

/* End a scope: "endgroup". Restores all 'save'd variables and declarations.
 * Restore of variable declarations is automatic due to popping the top scope,
 * where the saved variable definitions live.
 * In MetaFont, variables declared inside expression-groups live on as "capsules",
 * if they are still unknown and if they are handed back as the result-expression
 * of an expression-group. However, this is not relevant for compound statements
 * (a different kind of group).
 */
func (pl *GalleryParseListener) ExitCompound(ctx *grammar.CompoundContext) {
	corelang.Endgroup(pl.rt)
	pl.Summary()
}

/* Annotate AST to expect an lvalue for assignments (which is a variable
 * reference). We set a flag 'expectingLvalue' to suppress value substitution
 * when the variable is put onto the expression stack.
 *
 * assignment :  variable ASSIGN expression
 */
func (pl *GalleryParseListener) EnterAssignment(ctx *grammar.AssignmentContext) {
	pl.expectingLvalue = true
}

/* Variable assignment.
 *
 * assignment : lvalue ASSIGN numtertiary
 *
 * Both operands are on the expression stack: variable and expression
 *
 * (1) Retract lvalue from the resolver's table (make a capsule)
 * (3) Unset the value of lvalue
 * (3) Re-incarnate lvalue (get a new ID for it)
 * (4) Create equation on expression stack
 */
func (pl *GalleryParseListener) ExitAssignment(ctx *grammar.AssignmentContext) {
	e, _ := pl.rt.ExprStack.Pop()       // the expression value
	lvalue, ok := pl.rt.ExprStack.Pop() // the lvalue
	if !ok || !lvalue.IsValid() {
		T.Debug("operands broken for assignment")
	} else {
		if v := corelang.GetVariableFromExpression(pl.rt, lvalue); v != nil {
			varname := v.GetName()
			t := v.GetType()
			T.P("var", varname).Debugf("lvalue type is %s", variables.TypeString(t))
			if !pl.rt.ExprStack.CheckTypeMatch(lvalue, e) {
				T.P("var", varname).Errorf("type mismatch")
				panic("type mismatch in assignment")
			} else {
				corelang.Assign(pl.rt, v, e)
			}
		}
	}
}

/* Read an equation and put it into the LEQ solver.
 *
 * equation : expression ( EQUALS expression )+
 *
 * Equations may be chained, i.e. a=b=c. All operands are on the expression
 * stack. Operates right-associative.
 */
func (pl *GalleryParseListener) ExitEquation(ctx *grammar.EquationContext) {
	prev, _ := pl.rt.ExprStack.Pop() // walk the chain from right to left
	i := ctx.GetChildCount() / 2     // number of expr nodes -1
	for ; i > 0; i-- {               // invariant: prev is LHS of equation to the right
		tos := pl.rt.ExprStack.Top() // LHS of current equation
		pl.rt.ExprStack.Push(prev)   // LHS of equation to the right
		prev = tos                   // advance prev
		pl.rt.ExprStack.EquateTOS2OS()
	}
}

// --- Commands --------------------------------------------------------------

/* Save a tag within a group. The tag will be restored at the end of the
 * group. Save-commands within global scope will be ignored.
 *
 * SAVE TAG (COMMA TAG)*          # savecmd
 */
func (pl *GalleryParseListener) ExitSaveStmt(ctx *grammar.SavecmdContext) {
	T.Debugf("saving %d tags", len(ctx.AllTAG()))
	for _, tag := range ctx.AllTAG() { // list of tags
		corelang.Save(pl.rt, tag.GetText())
	}
}

/* Tracing command: show <tag>.
 *
 * SHOW TAG (COMMA TAG)*          # showcmd
 */
func (pl *GalleryParseListener) ExitShowcmd(ctx *grammar.ShowcmdContext) {
	for _, tag := range ctx.AllTAG() { // list of tags
		t := tag.GetText()
		T.P("tag", t).Infof("## show %s;", t)
		output := corelang.Showvariable(pl.rt, t)
		writer := T.Writer()
		writer.Write([]byte(output))
	}
}

// --- Expressions -----------------------------------------------------------

/* Add/subtract 2 factors. Factors may be of any type.
 *
 * tertiary : secondary                                  # term
 *          | tertiary (PLUS|MINUS) secondary            # term
 */
func (pl *GalleryParseListener) ExitTerm(ctx *grammar.TermContext) {
	if ctx.PLUS() != nil {
		pl.rt.ExprStack.AddTOS2OS()
	} else if ctx.MINUS() != nil {
		pl.rt.ExprStack.SubtractTOS2OS()
	} // fallthrough for sole secondary
}

/* Multiply/divide 2 factors. Factors may be of any kind.
 * For pairs:
 * (1) Multiply or divide a pair by a number, or
 * (2) multiply a number by a pair.
 *
 * secondary : primary                                    # factor
 *           | secondary (TIMES|OVER) primary             # factor
 */
func (pl *GalleryParseListener) ExitFactor(ctx *grammar.FactorContext) {
	if ctx.TIMES() != nil {
		pl.rt.ExprStack.MultiplyTOS2OS()
	} else if ctx.OVER() != nil {
		pl.rt.ExprStack.DivideTOS2OS()
	} // fallthrough for sole primary
}

/*
 */
/* Apply a function to a known argument.
 *
 * primary :   MATHFUNC atom                  # funcatom
 *
 * MATHFUNC : 'floor' | 'ceil' | 'sqrt' | @func ;
 */
func (pl *GalleryParseListener) ExitFuncatom(ctx *grammar.FuncatomContext) {
	fname := ctx.MATHFUNC().GetText()
	T.P("mathf", fname).Debug("applying function")
	e, ok := pl.rt.ExprStack.Pop()
	if !ok || !e.IsValid() {
		T.P("mathf", fname).Error("no arg present for function")
	} else {
		c, isconst := e.GetXPolyn().IsConstant()
		if !isconst {
			T.P("mathf", fname).Error("not implemented: f(<unknown>)")
		} else {
			// TODO: call Lua function
			c = corelang.Mathfunc(c, fname)
			pl.rt.ExprStack.PushConstant(c)
		}
	}
}

/* Put x-part or y-part of a pair variable on the expression stack.
 * The variable may be known or unknown.
 *
 * primary :  PAIRPART primary                            # pairpart
 */
func (pl *GalleryParseListener) ExitPairpart(ctx *grammar.PairpartContext) {
	part := ctx.PAIRPART().GetText() // xpart or ypart
	e := pl.rt.ExprStack.Top()
	T.Debugf("pairpart of: %v", e)
	if e.IsPair() {
		e, _ = pl.rt.ExprStack.Pop()
		if c, isconst := e.GetXPolyn().IsConstant(); isconst && part == "xpart" {
			pl.rt.ExprStack.PushConstant(c) // just push the value
		} else if c, isconst := e.GetYPolyn().IsConstant(); isconst && part == "ypart" {
			pl.rt.ExprStack.PushConstant(c) // just push the value
		} else {
			if v := corelang.GetVariableFromExpression(pl.rt, e); v != nil {
				T.Debugf("pair on the stack: %v", v)
				if part == "xpart" {
					pl.rt.ExprStack.PushVariable(v.XPart(), nil)
				} else {
					pl.rt.ExprStack.PushVariable(v.YPart(), nil)
				}
			} else { // this cannot happen...
				T.Errorf("not implemented: %s of anonymous unknown expression")
			}
		}
	} else {
		T.Errorf("cannot take %s of non-pair", part)
	}
}

/* Numeric expression primary: decimal constant, possibly including a unit.
 * Example: "3.14162mm".
 *
 * atom : DECIMALTOKEN UNIT?                    # decimal
 */
func (pl *GalleryParseListener) ExitDecimal(ctx *grammar.DecimalContext) {
	d := ctx.DECIMALTOKEN()
	num, _ := dec.NewFromString(d.GetText())
	if u := ctx.UNIT(); u != nil {
		//num = num.Mul(corelang.Unit2numeric(u.GetText())) // multiply with unit value
		num = corelang.ScaleDimension(num, u.GetText())
	}
	T.P("token", num.String()).Debug("numeric token")
	pl.rt.ExprStack.PushConstant(num) // put decimal number on expression stack
}

/* Variable reference as an expression primary. Example: "x2.r".
 * May be an lvalue (which must not be replaced by its value).
 *
 * variable : MIXEDTAG ( subscript | anytag )*
 *          | TAG ( subscript | anytag )*
 *          | LAMBDAARG
 * anytag   : TAG
 *          | MIXEDTAG
 *
 * TODO check for @#
 */
func (pl *GalleryParseListener) ExitVariable(ctx *grammar.VariableContext) {
	t := ctx.GetText()
	T.P("var", t).Debug("variable, verbose")
	s := corelang.CollectVarRefParts(pl.rt, t, ctx.GetChildren())
	vref := corelang.MakeCanonicalAndResolve(pl.rt, s, true) // create if not defined
	corelang.PushVariable(pl.rt, vref, pl.expectingLvalue)
	if pl.expectingLvalue {
		T.P("var", vref.GetName()).Debugf("lvalue type is %s",
			variables.TypeString(vref.GetType()))
		pl.expectingLvalue = false
	}
}

/* A variable subscript.
 */
func (pl *GalleryParseListener) ExitSubscript(ctx *grammar.SubscriptContext) {
	if ctx.DECIMALTOKEN() != nil {
		c, _ := dec.NewFromString(ctx.DECIMALTOKEN().GetText())
		pl.rt.ExprStack.PushConstant(c)
	}
}

/* Literal pair, i.e. a point with 2 corrdinates (x-part, y-part).
 *
 * atom : LPAREN tertiary COMMA tertiary RPAREN          # literalpair
 */
func (pl *GalleryParseListener) ExitLiteralpair(ctx *grammar.LiteralpairContext) {
	ey, _ := pl.rt.ExprStack.Pop()
	ex, _ := pl.rt.ExprStack.Pop()
	e := runtime.NewPairExpression(ex.GetXPolyn(), ey.GetXPolyn())
	T.Debugf("pair atom %s", e.String())
	pl.rt.ExprStack.Push(e)
}

/* Scale a (known or unknown) atom with a numeric coefficient.
 * The numeric coefficient is on the stack and may be applied to atoms
 * of different type.
 *
 * primary : scalarmulop atom        # scalarnumatom
 *
 * TOS is a polynomial (known or unknown), 2OS is numeric constant. We just
 * multiply them.
 *
 * Note that this allows a decimal constant to be prefixed by a decimal
 * constant:  1/3 12 yields 4. This is not allowed in MetaFont, but you
 * can always accomplish the same effect with a group:
 *
 * metafont> a=1/3 begingroup 12 endgroup;
 * ## a=4
 */
func (pl *GalleryParseListener) ExitScalaratom(ctx *grammar.ScalaratomContext) {
	pl.rt.ExprStack.MultiplyTOS2OS()
}

/* Attach a (signed) coefficient to a variable, e.g.  +3x, -1/3y.  We just have
 * to leave a numeric constant on the stack.
 *
 * scalarmulop : (PLUS|MINUS)? numtokenatom
 *
 * We have to handle the MINUS case only, as the rule for numtokenatom
 * already left a numeric constant on the expression stack.
 */
func (pl *GalleryParseListener) ExitScalarmulop(ctx *grammar.ScalarmulopContext) {
	if ctx.MINUS() != nil {
		pl.rt.ExprStack.PushConstant(arithm.MinusOne) // -1 on stack
		pl.rt.ExprStack.MultiplyTOS2OS()              // multiply with numtokenatom
	}
}

/* Numeric prefix for a variable, e.g., 3x, 1/2y.r, 0.25z.
 *
 * numtokenatom : DECIMALTOKEN '/' DECIMALTOKEN
 *              | DECIMALTOKEN
 */
func (pl *GalleryParseListener) ExitNumtokenatom(ctx *grammar.NumtokenatomContext) {
	numbers := ctx.AllDECIMALTOKEN()
	num1, _ := dec.NewFromString(numbers[0].GetText())
	T.P("token", num1.String()).Debug("numeric token")
	pl.rt.ExprStack.PushConstant(num1) // put decimal number on expression stack
	if len(numbers) > 1 {
		num2, _ := dec.NewFromString(numbers[1].GetText())
		T.P("token", num2.String()).Debug("numeric token")
		pl.rt.ExprStack.PushConstant(num2)
		pl.rt.ExprStack.DivideTOS2OS()
	}
}

/* Start a new scope within expressions.
 *
 * atom :  BEGINGROUP statementlist tertiary ENDGROUP     # exprgroup
 *
 * MetaFont allows begingroup ... ; expr endgroup  within expressions,
 * providing brackets around some statements and returning a (sub-)
 * expression.
 */
func (pl *GalleryParseListener) EnterExprgroup(ctx *grammar.ExprgroupContext) {
	groupscope, _ := corelang.Begingroup(pl.rt, "expr-group")
	pl.annotate(ctx, groupscope, "") // Annotate the AST node with this scope
}

/* See rule ExitCompound.
 *
 * atom :  BEGINGROUP statementlist tertiary ENDGROUP     # exprgroup
 *
 * Additionally leave the return expression on the stack.
 */
func (pl *GalleryParseListener) ExitExprgroup(ctx *grammar.ExprgroupContext) {
	corelang.Endgroup(pl.rt)
	// the return expression is already on the stack
	//pl.Summary()
}

// --- Variable Handling -----------------------------------------------------

/* Construct a valid variable reference string from parts on the stack.
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
 *
func (pl *GalleryParseListener) collectVarRefParts(t string, children []antlr.Tree) string {
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
	varname := vname.String()
	T.P("var", varname).Debug("collected parts")
	return varname
}

/* Get or create a variable reference. To get the canonical representation of
 * the variable reference, we parse it and construct a small AST. This AST
 * is fed into GetVarRefFromVarSyntax(). The resulting variable reference
 * struct is used to find the memory location of the variable reference.
 *
 * The reference lives in a memory frame, so we first locate it, then put
 * it on the expression stack. If the variable has a known value, we will
 * put the value onto the stack (otherwise the variable reference).
 *
func (pl *GalleryParseListener) makeCanonicalAndResolve(v string, create bool) *variables.PMMPVarRef {
	vtree := variables.ParseVariableFromString(v, &TracingErrorListener{})
	vref := variables.GetVarRefFromVarSyntax(vtree, pl.rt.ScopeTree)
	vref, _ = corelang.FindVariableReferenceInMemory(pl.rt, vref, create)
	return vref
}
*/
