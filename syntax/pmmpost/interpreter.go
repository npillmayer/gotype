package pmmpost

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

 * This is the implementation of an interpreter for "Poor Man's MetaPost",
 * my variant of the MetaPost graphical language. There is an accompanying
 * ANTLR grammar file, which describes the features and limitations of PMMPost.
 * I will sometimes refer to MetaFont, the original language underlying
 * MetaPost, as the grammar definitions are taken from Don Knuth's grammar
 * description in "The METAFONTBook".
 *
 * The implementation is tightly coupled to the ANTLR V4 parser generator.
 * ANTLR is a great tool and I see no use in being independent from it.

*/

import (
	"bytes"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/emirpasic/gods/stacks/linkedliststack"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/npillmayer/gotype/gtbackend/gfx"
	arithm "github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/syntax/pmmpost/corelang"
	pmmp "github.com/npillmayer/gotype/syntax/pmmpost/grammar"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	dec "github.com/shopspring/decimal"
)

// === Interpreter ===========================================================

/* We use AST-driven interpretation to execute the PMMPost input program. Input
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
 * Metafont is a dynamically scoped language. This means, functions can
 * access local variables from calling functions or groups. Nevertheless we
 * will find the definition of all variables (which are explicitly defined)
 * in a scope definition. This is mainly for type checking reasons and due
 * to the complex structure of MetaFont variable identifiers.
 */

/* Type PMMPostInterpreter.
 * This is an umbrella object to hold together the various tools needed to
 * execute steps 1 to 3 from above. It will orchestrate and instrument the
 * tools and execute them in the correct order. Also, this object will hold
 * and respect parameters we pass to the interpreter, so we can alter the
 * behaviour in certain aspects.
 */
type PMMPostInterpreter struct {
	ASTListener *ParseListener   // parse / AST listener
	runtime     *runtime.Runtime // runtime environment
	/*
		scopeTree     *runtime.ScopeTree        // collect scopes
		memFrameStack *runtime.MemoryFrameStack // runtime stack
		exprStack     *runtime.ExprStack        // eval expressions
		pathBuilder   *runtime.PathStack        // construct paths
	*/
}

/* Create a new Interpreter for "Poor Man's MetaPost". This is the top-level
 * object for this package.
 *
 * The interpreter manages an AST listener, a scope tree, a memory frame stack,
 * and stacks for numeric/pair and path expressions.
 *
 * Loads builtin symbols (variables and types).
 */
func NewPMMPostInterpreter() *PMMPostInterpreter {
	intp := &PMMPostInterpreter{}
	intp.runtime = runtime.NewRuntimeEnvironment(variables.NewPMMPVarDecl)
	/*
		intp.scopeTree = new(runtime.ScopeTree)                          // scopes for groups and functions
		intp.scopeTree.PushNewScope("globals", variables.NewPMMPVarDecl) // push global scope first
		intp.memFrameStack = new(runtime.MemoryFrameStack)               // initialize memory frame stack
		mf := intp.memFrameStack.PushNewMemoryFrame("global", nil)       // global memory
		mf.Scope = intp.scopeTree.Globals()                              // connect the global frame with the global scope
		intp.memFrameStack.Globals().SymbolTable = runtime.NewSymbolTable(variables.NewPMMPVarRef)
		intp.exprStack = runtime.NewExprStack()
		intp.pathBuilder = runtime.NewPathStack()
	*/
	pmmp.ScopeStack = intp.runtime.ScopeTree
	intp.loadBuiltinSymbols(intp.runtime.ScopeTree.Globals())             // load syms into new global scope
	intp.ASTListener = NewParseListener(intp.runtime.ScopeTree.Globals()) // listener for ANTLR
	intp.ASTListener.rt = intp.runtime
	intp.ASTListener.backend = &backend{}
	return intp
}

/* Load builtin symbols into a scope (usually the global scope).
 * TODO: nullpath, z@#, unitcircle, unitsquare
 */
func (intp *PMMPostInterpreter) loadBuiltinSymbols(scope *runtime.Scope) {
	originDef := corelang.Declare(intp.runtime, "origin", pmmp.PairType)
	origin := arithm.MakePair(arithm.ConstZero, arithm.ConstZero)
	_ = corelang.Variable(intp.runtime, originDef, origin, nil, true)
	upDef := corelang.Declare(intp.runtime, "up", pmmp.PairType)
	up := arithm.MakePair(arithm.ConstZero, arithm.ConstOne)
	_ = corelang.Variable(intp.runtime, upDef, up, nil, true)
	downDef := corelang.Declare(intp.runtime, "down", pmmp.PairType)
	down := arithm.MakePair(arithm.ConstZero, arithm.MinusOne)
	_ = corelang.Variable(intp.runtime, downDef, down, nil, true)
	rightDef := corelang.Declare(intp.runtime, "right", pmmp.PairType)
	right := arithm.MakePair(arithm.ConstOne, arithm.ConstZero)
	_ = corelang.Variable(intp.runtime, rightDef, right, nil, true)
	leftDef := corelang.Declare(intp.runtime, "left", pmmp.PairType)
	left := arithm.MakePair(arithm.MinusOne, arithm.ConstZero)
	_ = corelang.Variable(intp.runtime, leftDef, left, nil, true)
	_ = corelang.Declare(intp.runtime, "p", pmmp.PairType)
	_ = corelang.Declare(intp.runtime, "q", pmmp.PairType)
	_ = corelang.Declare(intp.runtime, "z", pmmp.PairType)
}

/* Set an output routine. Default is nil.
 */
func (intp *PMMPostInterpreter) SetOutputRoutine(o gfx.OutputRoutine) {
	if intp.ASTListener != nil {
		intp.ASTListener.backend.outputRoutine = o
	}
}

/* Parse and interpret a statement list.
 */
func (intp *PMMPostInterpreter) ParseStatements(input antlr.CharStream) {
	intp.ASTListener.ParseStatements(input)
}

// === AST driven parsing ====================================================

/* ANTLR will create an AST for us. We use a listener to attach to ANTLR's
 * AST walker. The listener manages scopes (with declarations) and memory frames
 * (with variable references and values).
 */
type ParseListener struct {
	*pmmp.BasePMMPStatemListener // build on top of ANTLR's base 'class'
	statemParser                 *pmmp.PMMPStatemParser
	annotations                  map[interface{}]Annotation // node annotations
	expectingLvalue              bool                       // do not evaluate variable
	rt                           *runtime.Runtime           // runtime environment
	backend                      *backend                   // where to draw to
	//varParser                    *variables.PMMPVarParser   // sub-parser for variables
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
	//tree := pl.statemParser.Figure() // start at top level rule: figure  TODO
	tree := pl.statemParser.Statement()
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
	lexer := pmmp.NewPMMPStatemLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&TracingErrorListener{})
	stream := antlr.NewCommonTokenStream(lexer, 0)
	if pl.statemParser == nil {
		pl.statemParser = pmmp.NewPMMPStatemParser(stream)
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

/* A picture has been completed.
 */
func (pl *ParseListener) ExitFigure(ctx *pmmp.FigureContext) {
	if pl.backend.picture != nil {
		T.Debug("figure complete")
		image := pl.backend.picture.AsImage()
		if pl.backend.outputRoutine != nil {
			pl.backend.outputRoutine.Shipout(pl.backend.picture.Name, image)
		} else {
			T.Error("no output routine set")
		}
		pl.backend.picture = nil
	}
}

/* Start a figure. Name, width and height are given.
 *
 * 'beginfig' '(' LABEL ',' DECIMALTOKEN UNIT ',' DECIMALTOKEN UNIT ')' SEMIC
 *
 */
func (pl *ParseListener) ExitBeginfig(ctx *pmmp.BeginfigContext) {
	d := ctx.DECIMALTOKEN(0)
	if d != nil {
		w, _ := dec.NewFromString(d.GetText())
		w = scaleDimension(w, ctx.UNIT(0).GetText())
		fw, _ := w.Float64()
		h, _ := dec.NewFromString(ctx.DECIMALTOKEN(1).GetText())
		h = scaleDimension(h, ctx.UNIT(1).GetText())
		fh, _ := h.Float64()
		name := strings.Trim(ctx.LABEL().GetText(), "\"")
		T.P("figure", name).Debugf("dimension %s x %s", w, h)
		pl.backend.picture = gfx.NewPicture(name, fw, fh)
		corelang.Begingroup(pl.rt, "figure")
		wdecl := corelang.Declare(pl.rt, "w", pmmp.NumericType)
		_ = corelang.Variable(pl.rt, wdecl, w, nil, true)
		hdecl := corelang.Declare(pl.rt, "h", pmmp.NumericType)
		_ = corelang.Variable(pl.rt, hdecl, h, nil, true)
	} else {
		T.Error("parse error, no figure completed")
	}
}

/* End a figure.
 */
func (pl *ParseListener) ExitEndfig(ctx *pmmp.EndfigContext) {
	if pl.backend.picture != nil {
		corelang.Endgroup(pl.rt)
	}
}

/* Draw command: draw a path. Draws a path using current pen and current color.
 *
 * drawCmd : 'draw' pathexpression
 *
 * pathexpression is TOS of path builder stack.
 */
func (pl *ParseListener) ExitDrawCmd(ctx *pmmp.DrawCmdContext) {
	pn, _ := pl.rt.PathBuilder.Pop()
	path := pn.Path
	pl.backend.picture.Draw(path)
}

/* Fill command: fill a cloxed path. Fills a path using current color.
 *
 * fillCmd : 'fill' pathexpression
 *
 * pathexpression is TOS of path builder stack.
 */
func (pl *ParseListener) ExitFillCmd(ctx *pmmp.FillCmdContext) {
	pn, _ := pl.rt.PathBuilder.Pop()
	path := pn.Path
	if pl.backend.picture != nil {
		pl.backend.picture.Fill(path)
	}
}

/* Pickup a pen. Example: "pickup pencircle scaled 3 withcolor #f080cc".
 *
 * 'pickup' PEN ( 'scaled' DECIMALTOKEN )? ( 'withcolor' COLOR )?
 *
 * The pen is used for subsequent drawing and filling commands.
 */
func (pl *ParseListener) ExitPickupCmd(ctx *pmmp.PickupCmdContext) {
	diam := arithm.ConstOne
	if ctx.DECIMALTOKEN() != nil {
		diam, _ = dec.NewFromString(ctx.DECIMALTOKEN().GetText())
	}
	var color color.Color = color.Black
	if ctx.COLOR() != nil {
		color = colorFromHex(ctx.COLOR().GetText())
	}
	pentype := ctx.PEN().GetText()
	pl.backend.pickupPen(pentype, diam, color)
}

/* Read an equation and put it into the LEQ solver.
 *
 * equation : numtertiary ( EQUALS numtertiary )+
 *            pairtertiary ( EQUALS pairtertiary )+
 *
 * Equations may be chained, i.e. a=b=c. All operands are on the expression
 * stack. Operates right-associative.
 */
func (pl *ParseListener) ExitMultiequation(ctx *pmmp.MultiequationContext) {
	prev, _ := pl.rt.ExprStack.Pop() // walk the chain from right to left
	i := ctx.GetChildCount() / 2     // number of expr nodes -1
	for ; i > 0; i-- {               // invariant: prev is LHS of equation to the right
		tos := pl.rt.ExprStack.Top() // LHS of current equation
		pl.rt.ExprStack.Push(prev)   // LHS of equation to the right
		v := corelang.GetVariableFromExpression(pl.rt, prev)
		T.Debugf("var for equation = TOS: %v", v)
		prev = tos // advance prev
		pl.rt.ExprStack.EquateTOS2OS()
	}
}

/* Complete a path equation. Path equations are essentially assignments,
 * as we do not allow paths to contain unknown points.
 *
 * pathatom EQUALS pathexpression          # pathequation
 *
 */
func (pl *ParseListener) ExitPathequation(ctx *pmmp.PathequationContext) {
	pe, ok := pl.rt.PathBuilder.Pop()
	pv := pl.rt.PathBuilder.Top()
	if ok {
		vref := pv.Symbol.(*variables.PMMPVarRef)
		//T.P("var", vref.GetName()).Debugf("set value = %s", pe.Path.String())
		vref.SetValue(pe.Path)
	} else {
		T.Error("cannot perform path equation: not enough parameters on the stack")
	}
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
func (pl *ParseListener) ExitAssignment(ctx *pmmp.AssignmentContext) {
	e, _ := pl.rt.ExprStack.Pop()       // the expression value
	lvalue, ok := pl.rt.ExprStack.Pop() // the lvalue
	if !ok || !lvalue.IsValid() {
		T.Debug("operands broken for assignment")
	} else {
		v := corelang.GetVariableFromExpression(pl.rt, lvalue)
		if v != nil {
			//varname := lvalue.GetXPolyn().TraceString(pl.exprStack)
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

/* Start a new scope: "begingroup".
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
func (pl *ParseListener) EnterCompound(ctx *pmmp.CompoundContext) {
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
func (pl *ParseListener) ExitCompound(ctx *pmmp.CompoundContext) {
	corelang.Endgroup(pl.rt)
	pl.Summary()
}

/* Start a new scope within expressions.
 *
 * numatom : BEGINGROUP statementlist numtertiary ENDGROUP  # exprgroup
 *
 * MetaFont allows begingroup ... ; expr endgroup  within expressions,
 * providing brackets around some pmmp and returning a (sub-)
 * expression.
 */
func (pl *ParseListener) EnterExprgroup(ctx *pmmp.ExprgroupContext) {
	groupscope, _ := corelang.Begingroup(pl.rt, "expr-group")
	pl.annotate(ctx, groupscope, "") // Annotate the AST node with this scope
}

/* See rule ExitCompound.
 *
 * numatom : BEGINGROUP statementlist numtertiary ENDGROUP  # exprgroup
 *
 * Additionally leave the return expression on the stack.
 */
func (pl *ParseListener) ExitExprgroup(ctx *pmmp.ExprgroupContext) {
	corelang.Endgroup(pl.rt)
	// the return expression is already on the stack
	pl.Summary()
}

/* End a scope: "endgroup". Restores all 'save'd variables and declarations.
 * Restore of variable declarations is automatic due to popping the top scope,
 * where the saved variable definitions live.
 * In MetaFont, variables declared inside expression-groups live on as "capsules",
 * if they are still unknown and if they are handed back as the result-expression
 * of an expression-group. However, this is not relevant for compound pmmp
 * (a different kind of group).
 */
func (pl *ParseListener) EnterPairexprgroup(ctx *pmmp.PairexprgroupContext) {
	groupscope, _ := corelang.Begingroup(pl.rt, "expr-group")
	pl.annotate(ctx, groupscope, "") // Annotate the AST node with this scope
}

/* See rule ExitCompound.
 *
 * pairatom : BEGINGROUP statementlist pairtertiary ENDGROUP
 *
 * Additionally leave the return expression (type pair) on the stack.
 */
func (pl *ParseListener) ExitPairexprgroup(ctx *pmmp.PairexprgroupContext) {
	corelang.Endgroup(pl.rt)
	// the return expression is already on the stack
	pl.Summary()
}

/* Save a tag within a group. The tag will be restored at the end of the
 * group. Save-commands within global scope will be ignored.
 *
 * saveStmt : 'save' tag (',' tag)*
 */
func (pl *ParseListener) ExitSaveStmt(ctx *pmmp.SaveStmtContext) {
	T.Debugf("saving %d tags", len(ctx.AllTag()))
	for _, tag := range ctx.AllTag() { // list of tags
		corelang.Save(pl.rt, tag.GetText())
	}
}

/* Finish a declaration statement. Example: "numeric a, b, c;"
 * If var decl is new, insert a new symbol in global scope. If var decl
 * already exists, erase all variables and re-enter var decl (MetaFont
 * semantics). If var decl has been "saved" in the current or in an outer scope,
 * make this tag a new undefined symbol.
 *
 * declaration :  mptype tag ( ',' tag )*
 */
func (pl *ParseListener) ExitDeclaration(ctx *pmmp.DeclarationContext) {
	T.P("type", ctx.Mptype().GetText()).Infof("declaration of %d tags", len(ctx.AllTag()))
	mftype := variables.TypeFromString(ctx.Mptype().GetText())
	if mftype == pmmp.Undefined {
		T.Error("unknown type: %s", ctx.Mptype().GetText())
		T.Error("assuming type numeric")
		mftype = pmmp.NumericType
	}
	for _, tag := range ctx.AllTag() { // iterator over tag list
		sym := corelang.Declare(pl.rt, tag.GetText(), mftype)
		var s *bytes.Buffer
		s = sym.ShowDeclarations(s)
		T.Infof("%s", s.String())
	}
}

/* Variable reference as a numeric expression primary. Example: "x2.r".
 *
 *      TAG ( subscript | TAG | MIXEDTAG )*
 * MIXEDTAG ( subscript | TAG | MIXEDTAG )*
 */
func (pl *ParseListener) ExitVariable(ctx *pmmp.VariableContext) {
	t := ctx.GetText()
	T.P("var", t).Debug("num-expr variable, verbose")
	s := pl.collectVarRefParts(t, ctx.GetChildren())
	vref := pl.makeCanonicalAndResolve(s, false)
	corelang.PushVariable(pl.rt, vref, false)
}

/* Variable reference as a pair expression primary. Example: "z3r".
 *
 *      PTAG ( subscript | anytag )*                  # pairvariable
 * MIXEDPTAG ( subscript | anytag )*                  # pairvariable
 */
func (pl *ParseListener) ExitPairvariable(ctx *pmmp.PairvariableContext) {
	t := ctx.GetText()
	T.P("var", t).Debug("pair-expr variable, verbose")
	s := pl.collectVarRefParts(t, ctx.GetChildren())
	vref := pl.makeCanonicalAndResolve(s, false)
	corelang.PushVariable(pl.rt, vref, false)
}

/* Variable as a path expression primary. Example: "P3". The variable and
 * its value are pushed to the PathBuilder.
 *
 * PATHTAG ( subscript | anytag )*       # pathvariable
 */
func (pl *ParseListener) ExitPathvariable(ctx *pmmp.PathvariableContext) {
	t := ctx.GetText()
	T.P("var", t).Debug("path-expr variable, verbose")
	s := pl.collectVarRefParts(t, ctx.GetChildren())
	vref := pl.makeCanonicalAndResolve(s, false)
	if vref.GetType() != pmmp.PathType { // automacically created as numeric
		T.P("var", t).Debug("setting type to path")
		vref.Decl.SetType(pmmp.PathType) // change type of declaration
		vref.SetType(pmmp.PathType)      // change type live variable
	}
	T.P("var", vref.GetName()).Debugf("pushing path = %.28v", vref.GetValue())
	pl.rt.PathBuilder.PushPath(vref, pathValue(vref.GetValue()))
}

// Helper
func pathValue(pv interface{}) arithm.Path {
	if pv == nil {
		return nil
	} else {
		return pv.(arithm.Path)
	}
}

/* Lvalue for assignments. This is in fact a variable reference.
 * We set a flag 'expectingLvalue' to suppress value substitution
 * when the variable is put onto the expression stack.
 */
func (pl *ParseListener) EnterLvalue(ctx *pmmp.LvalueContext) {
	pl.expectingLvalue = true
}

/* Lvalue for assignments. This is in fact a variable reference.
 */
func (pl *ParseListener) ExitLvalue(ctx *pmmp.LvalueContext) {
	t := ctx.GetText()
	T.P("var", t).Debug("lvalue variable, verbose")
	s := pl.collectVarRefParts(t, ctx.GetChildren())
	vref := pl.makeCanonicalAndResolve(s, pl.expectingLvalue)
	corelang.PushVariable(pl.rt, vref, pl.expectingLvalue)
	T.P("var", vref.GetName()).Debugf("lvalue type is %s", variables.TypeString(vref.GetType()))
	pl.expectingLvalue = false
}

/* Add/subtract 2 numeric terms.
 *
 * numtertiary (PLUS|MINUS) numsecondary
 */
func (pl *ParseListener) ExitNumtertiary(ctx *pmmp.NumtertiaryContext) {
	if ctx.PLUS() != nil {
		pl.rt.ExprStack.AddTOS2OS()
	} else if ctx.MINUS() != nil {
		pl.rt.ExprStack.SubtractTOS2OS()
	} // fallthrough for sole numsecondary
}

/* Add/subtract 2 pair terms.
 *
 * pairtertiary (PLUS|MINUS) pairsecondary
 */
func (pl *ParseListener) ExitPairtertiary(ctx *pmmp.PairtertiaryContext) {
	if ctx.PLUS() != nil {
		pl.rt.ExprStack.AddTOS2OS()
	} else if ctx.MINUS() != nil {
		pl.rt.ExprStack.SubtractTOS2OS()
	} // fallthrough for sole pairsecondary
}

/* Multiply 2 numeric factors.
 *
 * numprimary
 * numsecondary (TIMES|OVER) numprimary
 */
func (pl *ParseListener) ExitNumsecondary(ctx *pmmp.NumsecondaryContext) {
	if ctx.TIMES() != nil {
		pl.rt.ExprStack.MultiplyTOS2OS()
	} else if ctx.OVER() != nil {
		pl.rt.ExprStack.DivideTOS2OS()
	} // fallthrough for sole numprimary
}

/* Construct a path from path fragments.
 *
 * pathfragm ( PATHJOIN pathfragm )* cycle?
 *
 * Each fragment is either a known pair or a sub-path. Fragment AST nodes
 * are already labeled with a string marker. Pairs and sub-paths lie
 * in reverse order on either the expression stack or the path stack.
 * We therefore push each value on a small helper stack and collect them
 * in reverse order, appending them on the resulting path. The completed
 * path is pushed onto the path stack.
 */
func (pl *ParseListener) ExitPathtertiary(ctx *pmmp.PathtertiaryContext) {
	children := ctx.GetChildren()
	var stack *linkedliststack.Stack = linkedliststack.New()
	var isCycle, joinMsg bool
	for i := len(children) - 1; i >= 0; i-- { // process fragments in reverse order
		ch := ctx.GetChild(i)
		if isTerminal(ch) {
			if !joinMsg {
				T.Debugf("currently no differentiation between path joins")
				joinMsg = true
			}
		} else {
			_, t := pl.getAnnotation(ch) // read annotated label for fragment
			//T.Debugf("%d: %s = %s", i, t, getCtxText(ch))
			if t == "pair" {
				pr, ok := pl.rt.ExprStack.PopAsPair()
				if ok { // is a known pair
					//T.Debugf("adding pair to path: %s", pr.String())
					stack.Push(pr)
				} else {
					T.Error("cannot add unknown pair to path")
				}
			} else if t == "subpath" {
				pnode, _ := pl.rt.PathBuilder.Pop()
				//T.Debugf("fragment subpath %v / %v", pnode.Symbol, pnode.Path)
				//v := pnode.Symbol
				stack.Push(pnode.Path)
			} else { // cycle
				isCycle = true
			}
		}
	}
	//T.Debugf("stack is of size %d", stack.Size())
	// TODO this should be refactored to live in pathbuilder.go
	path := arithm.NullPath()
	tos, ok := stack.Pop()
	for ok { // as long as there are items on our helper stack
		if pr, ispair := tos.(arithm.Pair); ispair { // add pair
			path.AddPoint(pr)
		} else { // add subpath
			//T.Debugf("TOS = %v", tos)
			spath, _ := tos.(arithm.Path)
			path.AddSubpath(spath)
		}
		tos, ok = stack.Pop()
	}
	if path.Length() > 1 && isCycle {
		path.Cycle()
	}
	T.Debugf("new path = %s", path.String())
	pl.rt.PathBuilder.PushPath(nil, path) // push anonymous path
}

/* Fragment of a path: either a pair or a sub-path.
 *
 * pathfragm : pathsecondary
 *             pairtertiary
 *
 * We annotate the AST node with a string label: either "subpath" or "pair".
 */
func (pl *ParseListener) ExitPathfragm(ctx *pmmp.PathfragmContext) {
	if ctx.Pathsecondary() != nil {
		pl.annotate(ctx, nil, "subpath")
	} else if ctx.Pairtertiary() != nil {
		pl.annotate(ctx, nil, "pair")
	}
}

/* (1) Multiply or divide a pair by a number, or
 * (2) multiply a number by a pair.
 *
 *  pairprimary                                  // nothing to do
 *  pairsecondary (TIMES|OVER) numprimary
 *  numsecondary TIMES pairprimary
 */
func (pl *ParseListener) ExitPairsecond(ctx *pmmp.PairsecondContext) {
	if ctx.TIMES() != nil {
		pl.rt.ExprStack.MultiplyTOS2OS()
	} else if ctx.OVER() != nil {
		pl.rt.ExprStack.DivideTOS2OS()
	} // fallthrough for sole pairprimary
}

/* Apply a transform to a pair secondary.
 *
 * pairsecondary transformer                     # transform
 */
func (pl *ParseListener) ExitTransform(ctx *pmmp.TransformContext) {
	_, t := pl.getAnnotation(ctx.Transformer())
	T.Debugf("transform: %s", t)
	if t == "scaled" {
		scale, ok := pl.rt.ExprStack.PopAsNumeric()
		if !ok {
			T.P("transf", t).Error("need known numeric scale")
		} else {
			pl.rt.ExprStack.PushConstant(scale)
			pl.rt.ExprStack.MultiplyTOS2OS()
		}
	} else if t == "rotated" {
		angle, ok := pl.rt.ExprStack.PopAsNumeric()
		if !ok {
			T.P("transf", t).Error("need known numeric angle")
		} else {
			pl.rt.ExprStack.PushConstant(angle)
			pl.rt.ExprStack.Rotate2OSbyTOS()
		}
	} else if t == "shifted" {
		shift, ok := pl.rt.ExprStack.PopAsPair()
		if !ok {
			T.P("transf", t).Error("need known numeric pair")
		} else {
			pl.rt.ExprStack.PushPairConstant(shift)
			pl.rt.ExprStack.AddTOS2OS()
		}
	} else {
		T.P("transf", t).Error("unknown transform")
	}
}

/* Prepare a transform to apply to a pairsecondary.
 *
 * transformer : SCALED numprimary | ROTATED numprimary | SHIFTED pairprimary
 */
func (pl *ParseListener) ExitTransformer(ctx *pmmp.TransformerContext) {
	t := ctx.GetChild(0).(antlr.ParseTree).GetText()
	T.P("transf", t).Debug("transformer")
	pl.annotate(ctx, nil, t)
}

/* Scale a numeric (known or unknown) atom with a numeric coefficient.
 *
 * numprimary : scalarmulop numatom        # scalarnumatom
 *
 * TOS is a polynomial (known or unknown), 2OS is numeric constant. We just
 * multiply them.
 */
func (pl *ParseListener) ExitScalarnumatom(ctx *pmmp.ScalarnumatomContext) {
	pl.rt.ExprStack.MultiplyTOS2OS()
}

/* Apply a numeric function to a known numeric argument.
 *
 * numprimary : MATHFUNC numatom           # funcnumatom
 *
 * MATHFUNC : 'floor' | 'ceil' | 'sqrt' ;
 */
func (pl *ParseListener) ExitFuncnumatom(ctx *pmmp.FuncnumatomContext) {
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
			c = corelang.Mathfunc(c, fname)
			pl.rt.ExprStack.PushConstant(c)
		}
	}
}

/* Numeric interpolation, i.e. n[a,b].
 *
 *      numatom '[' numtertiary ',' numtertiary ']'
 * numtokenatom '[' numtertiary ',' numtertiary ']'
 *
 * All three expressions will be on the expression stack. We will convert
 * n[a,b] => a - na + nb.
 */
func (pl *ParseListener) ExitInterpolation(ctx *pmmp.InterpolationContext) {
	pl.rt.ExprStack.Interpolate()
}

/* Pair interpolation, i.e. n[z1,z2].
 *
 *      numatom '[' pairtertiary ',' pairtertiary ']'       # pairinterpolation
 * numtokenatom '[' pairtertiary ',' pairtertiary ']'       # pairinterpolation
 *
 * All three expressions will be on the expression stack. We will convert
 * n[a,b] => a - na + nb.
 */
func (pl *ParseListener) ExitPairinterpolation(ctx *pmmp.PairinterpolationContext) {
	pl.rt.ExprStack.Interpolate()
}

/* Length of a pair (i.e., distance from origin). Argument must be a known
 * pair.
 *
 * numprimary : LENGTH pairprimary                        # pointdistance
 */
func (pl *ParseListener) ExitPointdistance(ctx *pmmp.PointdistanceContext) {
	pl.rt.ExprStack.LengthTOS()
}

/* X-part or y-part of a pair variable.
 *
 * numprimary : PAIRPART pairprimary                      # pairpart
 */
func (pl *ParseListener) ExitPairpart(ctx *pmmp.PairpartContext) {
	e := pl.rt.ExprStack.Top()
	//T.Infof("pairpart of: %v", e)
	if e.IsPair() { // otherwise just leave numeric expression on stack
		e, _ = pl.rt.ExprStack.Pop()
		if c, isconst := e.GetXPolyn().IsConstant(); isconst {
			pl.rt.ExprStack.PushConstant(c) // just push the value
		} else {
			if v := corelang.GetVariableFromExpression(pl.rt, e); v != nil {
				//T.Infof("pair on the stack: %v", v)
				part := ctx.PAIRPART().GetText() // xpart or ypart
				if v.IsPair() {
					if part == "xpart" {
						pl.rt.ExprStack.PushVariable(v.XPart(), nil)
					} else {
						pl.rt.ExprStack.PushVariable(v.YPart(), nil)
					}
				} else {
					T.P("var", v.GetName()).Errorf("cannot take %s from numeric", part)
				}
			}
		}
	}
}

/* Rule for coefficients to variables, e.g.  +x, -y2.
 *
 * scalarmulop : ( PLUS | MINUS )          // leave +1 | -1 on stack
 *             | numtokenatom              // nothing to to
 */
func (pl *ParseListener) ExitScalarmulop(ctx *pmmp.ScalarmulopContext) {
	if ctx.PLUS() != nil {
		pl.rt.ExprStack.PushConstant(arithm.ConstOne) // put 1 on stack
	} else if ctx.MINUS() != nil {
		pl.rt.ExprStack.PushConstant(arithm.MinusOne) // -1 on stack
	}
}

/* Numeric prefix for a variable, e.g., 3x, 1/2y.r, -0.25z.
 *
 * numtokenatom : DECIMALTOKEN '/' DECIMALTOKEN
 *              | DECIMALTOKEN
 */
func (pl *ParseListener) ExitNumtokenatom(ctx *pmmp.NumtokenatomContext) {
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

/* Numeric expression primary: decimal constant, possibly including a unit.
 * Example: "3.14162mm".
 */
func (pl *ParseListener) ExitDecimal(ctx *pmmp.DecimalContext) {
	d := ctx.DECIMALTOKEN()
	num, _ := dec.NewFromString(d.GetText()) // TODO error possible ?!?
	if u := ctx.UNIT(); u != nil {
		num = num.Mul(unit2numeric(u.GetText())) // multiply with unit value
	}
	T.P("token", num.String()).Debug("numeric token")
	pl.rt.ExprStack.PushConstant(num) // put decimal number on expression stack
}

/* Literal pair, i.e. a point with 2 corrdinates (x-part, y-part).
 *
 * pairatom : '(' numtertiary ',' numtertiary ')'        # literalpair
 *
 */
func (pl *ParseListener) ExitLiteralpair(ctx *pmmp.LiteralpairContext) {
	ey, _ := pl.rt.ExprStack.Pop()
	ex, _ := pl.rt.ExprStack.Pop()
	e := runtime.NewPairExpression(ex.GetXPolyn(), ey.GetXPolyn())
	T.Debugf("pair atom %s", e.String())
	pl.rt.ExprStack.Push(e)
}

/* Get a point from a path. The path must be known.
 *
 * pairprimary: POINT numtertiary OF pathprimary          # pathpoint
 *
 * Expects the path to be on the path stack. Will leave a pair on the
 * expression stack.
 */
func (pl *ParseListener) ExitPathpoint(ctx *pmmp.PathpointContext) {
	pnode, ok := pl.rt.PathBuilder.Pop()
	if ok {
		num, _ := pl.rt.ExprStack.PopAsNumeric()
		i := int(num.IntPart())
		pr := pnode.Path.GetPoint(i)
		T.P("op", "point-of").Debugf("point #%d is %s", i, pr)
		pl.rt.ExprStack.PushPairConstant(pr)
	} else {
		T.P("op", "point-of").Error("expected path on the stack")
	}
}

/* Reverse a path. Produces a reversed copy and puts it onto the path
 * stack.
 *
 * pathprimary:  REVERSE pathprimary            # reversepath
 *
 * The pathprimary is expected on the path stack. It will be popped and
 * replaced by the reversed path.
 */
func (pl *ParseListener) ExitReversepath(ctx *pmmp.ReversepathContext) {
	pnode, ok := pl.rt.PathBuilder.Pop()
	if ok {
		path := pnode.Path.Copy()
		path.Reverse()
		T.P("op", "reverse").Debugf("reversed = %.30s", path.String())
		pl.rt.PathBuilder.PushPath(nil, path)
	} else {
		T.P("op", "reverse").Error("expected path on the stack")
	}
}

/* Re-use a part of a path.
 *
 * pathprimary:  SUBPATH pairtertiary OF pathprimary   # subpath
 *
 * The original path is not destoyed. It is expected to be on the path stack
 * and will be replaced by the subpath.
 */
func (pl *ParseListener) ExitSubpath(ctx *pmmp.SubpathContext) {
	pnode, ok := pl.rt.PathBuilder.Pop()
	if ok {
		fromto, _ := pl.rt.ExprStack.PopAsPair()
		path := pnode.Path.Copy()
		path.Subpath(int(fromto.XPart().IntPart()), int(fromto.YPart().IntPart()))
		pl.rt.PathBuilder.PushPath(nil, path)
	} else {
		T.P("op", "reverse").Error("expected path on the stack")
	}
}

/* Tracing command: showvariable <tag>.
 */
func (pl *ParseListener) ExitShowvariableCmd(ctx *pmmp.ShowvariableCmdContext) {
	tag := ctx.Tag().GetText()
	T.P("tag", tag).Infof("## showvariable %s;", tag)
	output := corelang.Showvariable(pl.rt, tag)
	writer := T.Writer()
	writer.Write([]byte(output))
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
func (pl *ParseListener) ExitSubscript(ctx *pmmp.SubscriptContext) {
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

/* Create a color from a hex string. Returns black if the string cannot
 * be interpreted as a color hex code.
 */
func colorFromHex(hex string) color.Color {
	c, err := colorful.Hex(hex)
	if err == nil {
		return c
	} else {
		return color.Black
	}
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
		_, ok = ctx.(pmmp.IAnytagContext)
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
