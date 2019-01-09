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

*/

import (
	"fmt"
	"image/color"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/npillmayer/gotype/backend/gfx"
	arithm "github.com/npillmayer/gotype/core/arithmetic"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/path"
	"github.com/npillmayer/gotype/core/polygon"
	"github.com/npillmayer/gotype/syntax/corelang"
	"github.com/npillmayer/gotype/syntax/pmmpost/grammar"
	"github.com/npillmayer/gotype/syntax/runtime"
	"github.com/npillmayer/gotype/syntax/variables"
	dec "github.com/shopspring/decimal"
)

// === Interpreter ===========================================================

/*
We use AST-driven interpretation to execute the input program. Input
is more or less a list of statements and function definitions.
We will annotate the AST with scope-information, holding the symbols of
dynamic scopes. Scopes stem from either:

- function definitions (macros in MetaFont: def and vardef)

- compound statements, i.e. groups (begingroup ... endgroup)

The interpreter relies on the scopes and definitions constructed earlier.
It manages a memory frame stack to track the calling sequence of functions
and groups.

So the overall picture looks like this:

1. ANTLR V4 constructs an AST for us.

2. We use a listener to walk the AST and execute the statements.

Metafont, and therefore PMMPost, is a dynamically scoped language. This means,
functions can access local variables from calling functions or groups.
Nevertheless we will find the definition of all variables (which are explicitly
defined) in a scope definition. This is mainly for type checking reasons and
due to the complex structure of MetaFont variable identifiers.

PMMPostInterpreter

This is an umbrella object to hold together the various tools needed to
execute steps 1 to 3 from above. It will orchestrate and instrument the
tools and execute them in the correct order. Also, this object will hold
and respect parameters we pass to the interpreter, so we can alter the
behaviour in certain aspects.
*/
type PMMPostInterpreter struct {
	ASTListener *PMMPostParseListener // parse / AST listener
	runtime     *runtime.Runtime      // runtime environment
	scripting   *corelang.Scripting   // Lua scripting subsystem
}

// Create a new Interpreter for the "Poor Man's MetaPost". This is the top-level
// object for this package.
//
// The interpreter manages an AST listener, a scope tree, a memory frame stack,
// and stacks for numeric/pair and path expressions.
//
// Loads builtin symbols (variables and types) if argument is true.
//
func NewPMMPostInterpreter(loadBuiltins bool) *PMMPostInterpreter {
	T = tracing.InterpreterTracer
	intp := &PMMPostInterpreter{}
	intp.runtime = runtime.NewRuntimeEnvironment(variables.NewPMMPVarDecl)
	intp.scripting = corelang.NewScripting(intp.runtime) // scripting subsystem
	if loadBuiltins {
		corelang.LoadBuiltinSymbols(intp.runtime, intp.scripting) // load syms into global scope
		intp.loadAdditionalBuiltinSymbols(intp.runtime, intp.scripting)
	}
	intp.ASTListener = NewParseListener(intp.runtime, intp.scripting) // listener for ANTLR
	return intp
}

// Set an output routine for drawing. Default is nil.
func (intp *PMMPostInterpreter) SetOutputRoutine(o gfx.OutputRoutine) {
	if intp.ASTListener != nil {
		intp.ASTListener.backend.outputRoutine = o
	}
}

// Parse and interpret a statement list.
//
func (intp *PMMPostInterpreter) ParseStatements(input []byte) {
	// func (intp *PMMPostInterpreter) ParseStatements(input antlr.CharStream) {
	inputStream := antlr.NewInputStream(string(input))
	intp.ASTListener.ParseStatements(inputStream)
}

// Load additionl builtins for this language (added to the core symbols)
func (intp *PMMPostInterpreter) loadAdditionalBuiltinSymbols(rt *runtime.Runtime,
	scripting *corelang.Scripting) {
	//
	//scripting.RegisterHook("TODO: z", ping)
}

// Value returns the current value of a variable.
// Return values are: canonical name of the variable & value of the
// variable as string.
func (intp *PMMPostInterpreter) Value(variable string) (string, string) {
	variable = strings.TrimSpace(variable)
	r, _ := utf8.DecodeRuneInString(variable)
	if !unicode.IsLetter(r) {
		return variable, "<illegal variable name>"
	}
	vtree := variables.ParseVariableFromString(variable, nil)
	if vtree == nil {
		return variable, "<illegal variable name>"
	}
	vref := variables.GetVarRefFromVarSyntax(vtree, intp.runtime.ScopeTree)
	if vref == nil {
		return variable, "<undefined>"
	}
	v := corelang.MakeCanonicalAndResolve(intp.runtime, variable, false)
	if v == nil {
		t := variables.TypeString(vref.GetType())
		return vref.GetFullName(), fmt.Sprintf("<unset %s>", t)
	}
	return v.GetFullName(), v.ValueString()
}

// === AST driven parsing ====================================================

// ANTLR will create an AST for us. We use a listener to attach to ANTLR's
// AST walker. The listener manages scopes (with declarations) and memory frames
// (with variable references and values).
//
type PMMPostParseListener struct {
	*grammar.BasePMMPostListener // build on top of ANTLR's base 'class'
	statemParser                 *grammar.PMMPostParser
	annotations                  map[interface{}]Annotation // node annotations
	expectingLvalue              bool                       // do not evaluate variable
	rt                           *runtime.Runtime           // runtime environment
	scripting                    *corelang.Scripting        // (Lua) scripting environment
	backend                      *backend                   // where to draw to
}

// We will annotate the AST. Functions and groups will get a scope, filled with
// statically local variable definitions. Scope information is tracked using a
// stack. Walking the AST results in scopes being pushed and popped
// on and off the scope stack (dynamically forming a scope tree). Scope tree
// scopes are always linking backward to their parent, with the global scope
// being the root scope. Every scope holds a symbol table.
//
// Whenever we identify a scope, we fill it with the local symbols and then
// attach it to the corresponding AST node. AST nodes will be either
//
// - function definitions (macros in MetaFont: def and vardef)
//
// - compound statements, i.e. groups (begingroup ... endgroup)
//
type Annotation struct {
	scope *runtime.Scope
	text  string
}

// Construct a new AST listener.
func NewParseListener(rt *runtime.Runtime, scr *corelang.Scripting) *PMMPostParseListener {
	pl := &PMMPostParseListener{} // no need to initialize base class
	pl.rt = rt
	pl.scripting = scr
	pl.annotations = make(map[interface{}]Annotation)
	pl.annotate("global", rt.ScopeTree.Globals(), "")
	return pl
}

// We use ANTLR V4 for parsing the statement grammar.
func (pl *PMMPostParseListener) ParseStatements(input antlr.CharStream) {
	pl.LazyCreateParser(input)
	tree := pl.statemParser.Statement() // TODO: Figures()
	//tree := pl.statemParser.Figures()
	sexpr := antlr.TreesStringTree(tree, nil, pl.statemParser)
	T.Debugf("### figures = %s", sexpr)
	antlr.ParseTreeWalkerDefault.Walk(pl, tree)
}

// Create an ANTLR V4 parser. This function should cache
// the parser for re-use, but currently I do not understand how to do this
// in the Go version of ANTLR. According to a forum discussion, creating a
// new parser every time seems to be the accepted mode of operation and should
// not carry too much of a performance penalty.
//
func (pl *PMMPostParseListener) LazyCreateParser(input antlr.CharStream) {
	// We let ANTLR to the heavy lifting.
	lexer := grammar.NewPMMPostLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&corelang.TracingErrorListener{})
	stream := antlr.NewCommonTokenStream(lexer, 0)
	if pl.statemParser == nil {
		pl.statemParser = grammar.NewPMMPostParser(stream)
		pl.statemParser.RemoveErrorListeners()
		pl.statemParser.AddErrorListener(&corelang.TracingErrorListener{})
		pl.statemParser.BuildParseTrees = true
	} else {
		pl.statemParser.SetInputStream(stream) // this should work
	}
}

// Annotate an AST node, i.e., attach a scope information.
func (pl *PMMPostParseListener) annotate(node interface{}, scope *runtime.Scope, text string) {
	pl.annotations[node] = Annotation{scope, text}
}

// Get the annotation for an AST node.
func (pl *PMMPostParseListener) getAnnotation(node interface{}) (*runtime.Scope, string) {
	if a, found := pl.annotations[node]; found {
		return a.scope, a.text
	}
	return nil, ""
}

// Print out a summary of all the scopes and symbols collected up to now.
func (pl *PMMPostParseListener) Summary() {
	pl.rt.ExprStack.Summary()
	T.Infof("Summary of symbols:")
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

// Print an error to the trace.
func (pl *PMMPostParseListener) VisitErrorNode(node antlr.ErrorNode) {
	T.Errorf("parser error: %s", node.GetText())
}

/* Helper to trace terminal symbols. Just traces to T.
func (pl *PMMPostParseListener) VisitTerminal(node antlr.TerminalNode) {
	//T.Debugf("@@@ terminal: %s  = %v", node.GetText(), node.GetSymbol())
}
*/

// --- Statements ------------------------------------------------------------

// A picture has been completed.
func (pl *PMMPostParseListener) ExitFigure(ctx *grammar.FigureContext) {
	if pl.backend.picture != nil {
		T.Debugf("figure complete")
		pl.backend.picture.Shipout()
		/*
			//image := pl.backend.picture.AsImage()
			if pl.backend.outputRoutine != nil {
				//pl.backend.outputRoutine.Shipout(pl.backend.picture, "png")
			} else {
				T.Errorf("no output routine set")
			}
		*/
		pl.backend.picture = nil
	}
}

// Start a figure. Name, width and height are given.
//
//    'beginfig' '(' LABEL ',' DECIMALTOKEN UNIT ',' DECIMALTOKEN UNIT ')' SEMIC
//
func (pl *PMMPostParseListener) ExitBeginfig(ctx *grammar.BeginfigContext) {
	d := ctx.DECIMALTOKEN(0)
	if d != nil {
		w, _ := dec.NewFromString(d.GetText())
		w = corelang.ScaleDimension(w, ctx.UNIT(0).GetText())
		fw, _ := w.Float64()
		h, _ := dec.NewFromString(ctx.DECIMALTOKEN(1).GetText())
		h = corelang.ScaleDimension(h, ctx.UNIT(1).GetText())
		fh, _ := h.Float64()
		name := strings.Trim(ctx.LABEL().GetText(), "\"")
		T.P("figure", name).Debugf("dimension %s x %s", w, h)
		pl.backend.picture = gfx.NewPicture(name, fw, fh, "png") // TODO: format?
		corelang.Begingroup(pl.rt, "figure")
		wdecl := corelang.Declare(pl.rt, "w", variables.NumericType)
		_ = corelang.Variable(pl.rt, wdecl, w, nil, true)
		hdecl := corelang.Declare(pl.rt, "h", variables.NumericType)
		_ = corelang.Variable(pl.rt, hdecl, h, nil, true)
	} else {
		T.Errorf("parse error, no figure completed")
	}
}

// End a figure.
func (pl *PMMPostParseListener) ExitEndfig(ctx *grammar.EndfigContext) {
	if pl.backend.picture != nil {
		corelang.Endgroup(pl.rt)
	}
}

// Pickup a pen. Example: "pickup pencircle scaled 3 withcolor #f080cc".
//
//    'pickup' PEN ( 'scaled' DECIMALTOKEN )? ( 'withcolor' COLOR )?
//
// The pen is used for subsequent drawing and filling commands.
//
func (pl *PMMPostParseListener) ExitPickupCmd(ctx *grammar.PickupCmdContext) {
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

// Draw command: draw a path. Draws a path using current pen and current color.
//
//    drawCmd : 'draw' pathexpression
//
// pathexpression is TOS of path builder stack.
//
func (pl *PMMPostParseListener) ExitDrawCmd(ctx *grammar.DrawCmdContext) {
	e, _ := pl.rt.ExprStack.Pop()
	path := e.Other.(*path.Path)
	pl.backend.picture.Draw(gfx.NewDrawablePath(path, path.Controls))
}

// Fill command: fill a cloxed path. Fills a path using current color.
//
//    fillCmd : 'fill' pathexpression
//
// pathexpression is TOS of path builder stack.
//
func (pl *PMMPostParseListener) ExitFillCmd(ctx *grammar.FillCmdContext) {
	e, _ := pl.rt.ExprStack.Pop()
	path := e.Other.(*path.Path)
	if pl.backend.picture != nil {
		pl.backend.picture.Fill(gfx.NewDrawablePath(path, path.Controls))
	}
}

// Finish a declaration statement. Example: "numeric a, b, c;"
// If var decl is new, insert a new symbol in global scope. If var decl
// already exists, erase all variables and re-enter var decl (MetaFont
// semantics). If var decl has been "saved" in the current or in an outer scope,
// make this tag a new undefined symbol.
//
//    TYPE TAG ( COMMA TAG )*                        # typedecl
//
func (pl *PMMPostParseListener) ExitTypedecl(ctx *grammar.TypedeclContext) {
	T.P("type", ctx.TYPE().GetText()).Debugf("declaration of %d tags", len(ctx.AllTAG()))
	mftype := variables.TypeFromString(ctx.TYPE().GetText())
	if mftype == variables.Undefined {
		T.Errorf("unknown type: %s", ctx.TYPE().GetText())
		T.Errorf("assuming type numeric")
		mftype = variables.NumericType
	}
	for _, tag := range ctx.AllTAG() { // iterator over tag list
		_ = corelang.Declare(pl.rt, tag.GetText(), mftype)
	}
}

// Finish a vardef. Vardefs are primary tags, which call a Lua script.
//
// 	vardef : VARDEF TAG ( COMMA TAG )*
//
func (pl *PMMPostParseListener) ExitVardef(ctx *grammar.VardefContext) {
	T.P("def", "vardef").Debugf("declaration of %d tags", len(ctx.AllTAG()))
	for _, tag := range ctx.AllTAG() { // iterator over tag list
		_ = corelang.Declare(pl.rt, tag.GetText(), variables.VardefType)
	}
}

// Start a new scope / compound statement: "begingroup".
//
// MetaFont uses dynamic scopes with the "begingroup ... endgroup" notation,
// but in a unique way: users have to "save" variables explicitly
// (The METAFONTBook, chapter 17).
//
// We produce a new scope for the scope tree and a new memory frame onto
// the memory frame stack. The scope holds declarations for local variables.
// The memory frame holds variable references (pointing to the the decls
// in the scope).
//
// We reserve a symbol in the group's scope on "save" statements only.
// Saved tags are of type tag/undefined initially.
//
// Notice: It is a bit of an overkill to use dynamic scopes and memory frames
// in parallel, but I'll stick to this traditional approach for clarity.
// Furthermore, this approach comes handy for the runtime-handling of
// definitions (a.k.a. macros) and/or functions.
//
func (pl *PMMPostParseListener) EnterCompound(ctx *grammar.CompoundContext) {
	groupscope, _ := corelang.Begingroup(pl.rt, "compound-group")
	pl.annotate(ctx, groupscope, "") // Annotate the AST node with this scope
}

// End a scope: "endgroup". Restores all 'save'd variables and declarations.
// Restore of variable declarations is automatic due to popping the top scope,
// where the saved variable definitions live.
//
// In MetaFont, variables declared inside expression-groups live on as "capsules",
// if they are still unknown and used in unsolved equations with global variables.
//
func (pl *PMMPostParseListener) ExitCompound(ctx *grammar.CompoundContext) {
	corelang.Endgroup(pl.rt)
	pl.Summary()
}

// Annotate AST to expect an lvalue for assignments (which is a variable
// reference). We set a flag 'expectingLvalue' to suppress value substitution
// when the variable is put onto the expression stack.
//
//    assignment :  variable ASSIGN expression
//
func (pl *PMMPostParseListener) EnterAssignment(ctx *grammar.AssignmentContext) {
	pl.expectingLvalue = true
}

// Variable assignment.
//
//    assignment : lvalue ASSIGN numtertiary
//
// Both operands are on the expression stack: variable and expression
//
// (1) Retract lvalue from the resolver's table (make a capsule)
//
// (3) Unset the value of lvalue
//
// (3) Re-incarnate lvalue (get a new ID for it)
//
// (4) Create equation on expression stack
//
func (pl *PMMPostParseListener) ExitAssignment(ctx *grammar.AssignmentContext) {
	e, _ := pl.rt.ExprStack.Pop()       // the expression value
	lvalue, ok := pl.rt.ExprStack.Pop() // the lvalue
	if !ok || !lvalue.IsValid() {
		T.Debugf("operands broken for assignment")
	} else {
		if v := corelang.GetVariableFromExpression(pl.rt, lvalue); v != nil {
			varname := v.GetName()
			t := v.GetType()
			T.P("var", varname).Debugf("lvalue type is %s", variables.TypeString(t))
			if !pl.rt.ExprStack.CheckTypeMatch(lvalue, e) {
				T.P("var", varname).Errorf("type mismatch in assignment")
				panic("type mismatch in assignment")
			} else {
				corelang.Assign(pl.rt, v, e)
			}
		}
	}
}

// Read an equation and put it into the LEQ solver.
//
//    equation : expression ( EQUALS expression )+
//
// Equations may be chained, i.e.
//
//    a=b=c
//
// All operands are on the expression stack. Operates right-associative.
//
func (pl *PMMPostParseListener) ExitEquation(ctx *grammar.EquationContext) {
	prev, _ := pl.rt.ExprStack.Pop() // walk the chain from right to left
	i := ctx.GetChildCount() / 2     // number of expr nodes -1
	for ; i > 0; i-- {               // invariant: prev is LHS of equation to the right
		tos := pl.rt.ExprStack.Top() // LHS of current equation
		pl.rt.ExprStack.Push(prev)   // LHS of equation to the right
		prev = tos                   // advance prev
		pl.equate()                  // perform equation depending on operands' type
	}
}

// Equation for different expression types.
func (pl *PMMPostParseListener) equate() {
	rhs, _ := pl.rt.ExprStack.Pop()
	lhs := pl.rt.ExprStack.Top()
	etype := getExprType(lhs)
	T.Debugf("%s equation: %v = %v", variables.TypeString(etype), lhs, rhs)
	switch etype {
	case variables.NumericType:
		pl.rt.ExprStack.Push(rhs)
		pl.rt.ExprStack.EquateTOS2OS()
	case variables.PairType:
		pl.rt.ExprStack.Push(rhs)
		pl.rt.ExprStack.EquateTOS2OS()
	case variables.PathType:
		lhs.Other = rhs.Other // now the path is on the stack
		sym := pl.rt.ExprStack.GetVariable(lhs)
		vref, _ := sym.(*variables.PMMPVarRef)
		vref.SetValue(rhs.Other)
	default:
		T.P("op", "=").Errorf("not implemented: equation for type %s", variables.TypeString(etype))
	}
}

// --- Commands --------------------------------------------------------------

// Save a tag within a group. The tag will be restored at the end of the
// group. Save-commands within global scope will be ignored.
//
//    SAVE TAG (COMMA TAG)*          # savecmd
//
func (pl *PMMPostParseListener) ExitSaveStmt(ctx *grammar.SavecmdContext) {
	T.Debugf("saving %d tags", len(ctx.AllTAG()))
	for _, tag := range ctx.AllTAG() { // list of tags
		corelang.Save(pl.rt, tag.GetText())
	}
}

// Tracing command: show <tag>.
//
//    SHOW TAG (COMMA TAG)*          # showcmd
//
func (pl *PMMPostParseListener) ExitShowcmd(ctx *grammar.ShowcmdContext) {
	for _, tag := range ctx.AllTAG() { // list of tags
		t := tag.GetText()
		T.P("tag", t).Infof("## show %s;", t)
		output := corelang.Showvariable(pl.rt, t)
		//writer := T.Writer()
		//writer.Write([]byte(output))
		T.Infof(output)
	}
}

// --- Expressions -----------------------------------------------------------

// Top level expression.
//
//    expression : tertiary
//               | expression PATHCLIPOP tertiary
//
// Path clipping operations are
//
//    PATHCLIPOP : 'union' | 'intersection' | 'difference'
//
func (pl *PMMPostParseListener) ExitExpression(ctx *grammar.ExpressionContext) {
	if ctx.PATHCLIPOP() != nil {
		e, _ := pl.rt.ExprStack.Pop()
		p := e.Other.(path.Path)
		/*
			e, _ := pl.rt.ExprStack.Pop()
			p1 := e.Other.(path.Path)
			e, _ = pl.rt.ExprStack.Pop()
			p2 := e.Other.(path.Path)
			var p path.Path
				switch ctx.PATHCLIPOP().GetText() {
				case "union":
					p = polygon.Union(p1, p2)
				case "intersection":
					p = polygon.Union(p1, p2)
				case "difference":
					p = polygon.Union(p1, p2)
				default:
					T.P("op", ctx.PATHCLIPOP().GetText()).Errorf("unknown clipping operation")
					p = p1
				}
		*/
		T.P("op", ctx.PATHCLIPOP().GetText()).Errorf("unknown clipping operation")
		pl.rt.ExprStack.Push(runtime.NewOtherExpression(p))
	}
}

// Add/subtract 2 factors. Factors may be of type numeric or pair.
//
//   tertiary : secondary                                  # term
//            | tertiary (PLUS|MINUS) secondary            # term
//
func (pl *PMMPostParseListener) ExitTerm(ctx *grammar.TermContext) {
	if ctx.PLUS() != nil {
		pl.rt.ExprStack.AddTOS2OS()
	} else if ctx.MINUS() != nil {
		pl.rt.ExprStack.SubtractTOS2OS()
	} // fallthrough for lone secondary
}

// Multiply/divide 2 factors. Factors may be of any kind.
// For pairs:
//
// (1) Multiply or divide a pair by a number, or
//
// (2) multiply a number by a pair.
//
//   secondary : primary                                    # factor
//             | secondary (TIMES|OVER) primary             # factor
//
func (pl *PMMPostParseListener) ExitFactor(ctx *grammar.FactorContext) {
	if ctx.TIMES() != nil {
		pl.rt.ExprStack.MultiplyTOS2OS()
	} else if ctx.OVER() != nil {
		pl.rt.ExprStack.DivideTOS2OS()
	} // fallthrough for lone primary
}

// Create a path by joining path fragments.
//
//   path : secondary ( pathjoin secondary )+ cycle?
//
// Each fragment is either a known pair or a sub-path. Fragment AST nodes
// are already labeled with a string marker. Pairs and sub-paths lie
// in reverse order on either the expression stack or the path stack.
//
// The completed path is pushed onto the path stack.
//
func (pl *PMMPostParseListener) ExitPath(ctx *grammar.PathContext) {
	builder := NewPathBuilder()
	l := len(ctx.AllSecondary()) - 1
	for i, _ := range ctx.AllSecondary() {
		if e, _ := pl.rt.ExprStack.Pop(); e.IsPair {
			if pr, isconst := e.GetConstPair(); isconst {
				T.P("op", "path").Debugf(" fragment %d is pair: %s", l-i, e.String())
				builder.CollectKnot(pr)
			} else {
				T.Errorf("polygon elements must be known")
			}
		} else if e.Other != nil {
			if p, ispath := e.Other.(*path.Path); ispath {
				builder.CollectSubpath(p)
			} else {
				T.P("op", "path").Errorf(" fragment %d is unknown type", l-i)
			}
		} else {
			T.P("op", "path").Errorf(" fragment %d is empty", l-i)
		}
	}
	if ctx.Cycle() != nil {
		builder.Cycle()
	}
	path := builder.MakePath()
	e := runtime.NewOtherExpression(path) // return the new path on the stack
	pl.rt.ExprStack.Push(e)
}

// Apply a chain of affine transforms to a pair or a path.
//
//   secondary ( TRANSFORM primary )+           # transform
//
func (pl *PMMPostParseListener) ExitTransform(ctx *grammar.TransformContext) {
	l := len(ctx.AllTRANSFORM())           // count transforms
	transforms := make([]string, l)        // array for transform operators
	for i, t := range ctx.AllTRANSFORM() { // for every transform
		transforms[i] = t.GetText()
	}
	trArgs := make([]*runtime.ExprNode, l) // array for transform parameters
	for i := l - 1; i >= 0; i-- {          // store all transform parameters
		trArgs[i], _ = pl.rt.ExprStack.Pop() // t.-parameters are on the stack
		T.Debugf("transform %d/%d =  %s %v", i, l, ctx.TRANSFORM(i), trArgs[i])
	}
	// now check secondary (on top of stack): pair, path or transform variable
	e := pl.rt.ExprStack.Top() // leave it on the stack
	etype := getExprType(e)
	switch etype {
	case variables.PairType:
		T.Debugf("transform of pair")
		pl.transformPair(e, transforms, trArgs) // transform a pair
	case variables.PathType:
		T.Debugf("transform of polygon")
		pl.transformPath(e, transforms, trArgs) // transform a path
	default: // TODO: transform of transform
		T.Errorf("cannot transform type %s", variables.TypeString(etype))
	}
}

func (pl *PMMPostParseListener) transformPair(e *runtime.ExprNode, transforms []string,
	parameters []*runtime.ExprNode) {
	//
	_, isknown := e.GetConstPair()
	for i, t := range transforms {
		switch t {
		case "scaled":
			_, isconst := parameters[i].GetConstNumeric()
			if !isknown && !isconst {
				T.Errorf("not implemented: <unknown pair> scaled <unknown>")
			} else {
				pl.rt.ExprStack.Push(parameters[i])
				pl.rt.ExprStack.MultiplyTOS2OS()
				T.P("op", "scale").Debugf("result = %v", pl.rt.ExprStack.Top())
			}
		case "shifted":
			_, isconst := parameters[i].GetConstPair()
			if !isknown && !isconst {
				T.Errorf("not implemented: <unknown pair> shifted <unknown>")
			} else {
				pl.rt.ExprStack.Push(parameters[i])
				pl.rt.ExprStack.AddTOS2OS()
				T.P("op", "shift").Debugf("result = %v", pl.rt.ExprStack.Top())
			}
		case "rotated":
			_, isconst := parameters[i].GetConstNumeric()
			if !isconst {
				T.Errorf("not implemented: rotated <unknown>")
			} else {
				pl.rt.ExprStack.Push(parameters[i])
				pl.rt.ExprStack.Rotate2OSbyTOS()
				T.P("op", "rotate").Debugf("result = %v", pl.rt.ExprStack.Top())
			}
		default:
			T.Errorf("ignoring unknown transform: %s", t)
		}
	}
}

func (pl *PMMPostParseListener) transformPath(e *runtime.ExprNode, transforms []string,
	parameters []*runtime.ExprNode) {
	//
	p, _ := e.Other.(polygon.Polygon) // secondary must be known path
	affinetr := arithm.Identity()     // we'll construct an overall transform
	for i, t := range transforms {    // apply all transforms to affinetr
		switch t {
		case "scaled":
			_, isconst := parameters[i].GetConstNumeric()
			if !isconst {
				T.Errorf("not implemented: <path> scaled <unknown>")
			} else {
				T.Errorf("not implemented: scaling of <path>")
			}
		case "shifted":
			pr, isconst := parameters[i].GetConstPair()
			if !isconst {
				T.Errorf("not implemented: <path> shifted <unknown>")
			} else {
				a := arithm.Translation(pr)
				affinetr = affinetr.Combine(a)
			}
		case "rotated":
			angle, isconst := parameters[i].GetConstNumeric()
			if !isconst {
				T.Errorf("not implemented: <path> rotated <unknown>")
			} else {
				a := arithm.Rotation(angle.Mul(arithm.Deg2Rad))
				affinetr = affinetr.Combine(a)
			}
		default:
			T.Errorf("ignoring unknown transform: %s", t)
		}
	}
	// now apply the resulting transform to the path/polygon
	e.Other = polygon.Transform(p, affinetr)
}

// Apply a function to a known argument. Functions may be internal math
// functions (working on type numeric) or calls to the (Lua) scripting
// sub-system.
//
//   primary  :  MATHFUNC atom                  # funcatom
//   MATHFUNC : 'floor' | 'ceil' | 'sqrt' | @func ;
//
func (pl *PMMPostParseListener) ExitFuncatom(ctx *grammar.FuncatomContext) {
	fname := ctx.MATHFUNC().GetText()
	T.P("func", fname).Debugf("applying function")
	e, ok := pl.rt.ExprStack.Pop()
	if ok {
		var val interface{}
		etype := getExprType(e)
		switch etype {
		case variables.NumericType:
			if val, ok = e.GetConstNumeric(); !ok {
				T.P("func", fname).Errorf("not implemented: f(<unknown numeric>)")
				val = nil
			}
		case variables.PairType:
			if val, ok = e.GetConstPair(); !ok {
				T.P("func", fname).Errorf("not implemented: f(<unknown pair>)")
				val = nil
			}
		case variables.PathType:
			val = e.Other
		default:
			T.P("func", fname).Errorf("not implemented: f(%s)", variables.TypeString(etype))
		}
		if val != nil {
			e, vars := corelang.CallFunc(val, fname, pl.scripting)
			announceVariables(vars, pl.rt.ExprStack)
			pl.rt.ExprStack.Push(e)
		} else {
			pl.rt.ExprStack.PushConstant(arithm.ConstZero)
		}
	}
}

// Numeric interpolation, i.e. "n[a,b]".
//
//    primary
//       | numtokenatom [ tertiary , tertiary ]       # interpolation
//       |         atom [ tertiary , tertiary ]       # interpolation
//
// All three expressions will be on the expression stack. We will convert
// n[a,b] => a - na + nb.
//
func (pl *PMMPostParseListener) ExitInterpolation(ctx *grammar.InterpolationContext) {
	pl.rt.ExprStack.Interpolate()
}

// Put x-part or y-part of a pair variable on the expression stack.
// The variable may be known or unknown.
//
//    primary :  PAIRPART primary                            # pairpart
//
func (pl *PMMPostParseListener) ExitPairpart(ctx *grammar.PairpartContext) {
	part := ctx.PAIRPART().GetText() // xpart or ypart
	e := pl.rt.ExprStack.Top()
	T.Debugf("pairpart of: %v", e)
	if e.IsPair {
		e, _ = pl.rt.ExprStack.Pop()
		if c, isconst := e.XPolyn.IsConstant(); isconst && part == "xpart" {
			pl.rt.ExprStack.PushConstant(c) // just push the value
		} else if c, isconst := e.YPolyn.IsConstant(); isconst && part == "ypart" {
			pl.rt.ExprStack.PushConstant(c) // just push the value
		} else {
			if v := corelang.GetVariableFromExpression(pl.rt, e); v != nil {
				T.Debugf("part of pair on the stack: %v", v)
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

// Numeric expression primary: decimal constant, possibly including a unit.
//
// Example: "3.14162mm".
//
//    atom : DECIMALTOKEN UNIT?                    # decimal
//
func (pl *PMMPostParseListener) ExitDecimal(ctx *grammar.DecimalContext) {
	d := ctx.DECIMALTOKEN()
	num, _ := dec.NewFromString(d.GetText())
	if u := ctx.UNIT(); u != nil {
		//num = num.Mul(corelang.Unit2numeric(u.GetText())) // multiply with unit value
		num = corelang.ScaleDimension(num, u.GetText())
	}
	T.P("token", num.String()).Debugf("numeric token")
	pl.rt.ExprStack.PushConstant(num) // put decimal number on expression stack
}

// Variable reference as an expression primary. Example: "x2.r".
// May be an lvalue (which must not be replaced by its value).
//
//    variable : MIXEDTAG ( subscript | anytag )*
//             | TAG ( subscript | anytag )*
//             | LAMBDAARG
//    anytag   : TAG
//             | MIXEDTAG
//
// TODO check for @#
//
func (pl *PMMPostParseListener) ExitVariable(ctx *grammar.VariableContext) {
	t := ctx.GetText()
	T.P("var", t).Debugf("variable, verbose")
	var vref *variables.PMMPVarRef
	if t == "whatever" {
		vref = corelang.Whatever(pl.rt)
		corelang.PushVariable(pl.rt, vref, pl.expectingLvalue)
	} else {
		s := corelang.CollectVarRefParts(pl.rt, t, ctx.GetChildren())
		vref = corelang.MakeCanonicalAndResolve(pl.rt, s, true)
		if vref.GetType() == variables.VardefType {
			pl.callVardef(vref, pl.expectingLvalue)
		} else {
			corelang.PushVariable(pl.rt, vref, pl.expectingLvalue)
		}
	}
	if pl.expectingLvalue {
		T.P("var", vref.GetName()).Debugf("lvalue type is %s",
			variables.TypeString(vref.GetType()))
		pl.expectingLvalue = false
	}
}

// Call Lua for a varref.
func (pl *PMMPostParseListener) callVardef(vref *variables.PMMPVarRef, expectLvalue bool) {
	e, vars := corelang.CallVardef(vref, pl.scripting)
	T.P("vardef", vref.GetName()).Debugf("return value from Lua call: %v", e)
	announceVariables(vars, pl.rt.ExprStack)
	pl.rt.ExprStack.Push(e)
}

// Helper
func announceVariables(vars []*variables.PMMPVarRef, es *runtime.ExprStack) {
	for _, v := range vars {
		es.AnnounceVariable(v)
	}
}

// A variable subscript.
func (pl *PMMPostParseListener) ExitSubscript(ctx *grammar.SubscriptContext) {
	if ctx.DECIMALTOKEN() != nil {
		c, _ := dec.NewFromString(ctx.DECIMALTOKEN().GetText())
		pl.rt.ExprStack.PushConstant(c)
	}
}

// Literal pair, i.e. a point with 2 corrdinates (x-part, y-part).
//
//    atom : LPAREN tertiary COMMA tertiary RPAREN          # literalpair
//
func (pl *PMMPostParseListener) ExitLiteralpair(ctx *grammar.LiteralpairContext) {
	ey, _ := pl.rt.ExprStack.Pop()
	ex, _ := pl.rt.ExprStack.Pop()
	e := runtime.NewPairExpression(ex.XPolyn, ey.XPolyn)
	T.Debugf("pair atom %s", e.String())
	pl.rt.ExprStack.Push(e)
}

// Scale a (known or unknown) atom with a numeric coefficient.
// The numeric coefficient is on the stack and may be applied to atoms
// of different type.
//
//    primary : scalarmulop atom        # scalarnumatom
//
// TOS is a polynomial (known or unknown), 2OS is numeric constant. We just
// multiply them.
//
// Note that this allows a decimal constant to be prefixed by a decimal
// constant:  1/3 12 yields 4. This is not allowed in MetaFont, but you
// can always accomplish the same effect with a group:
//
//    metafont> a=1/3 begingroup 12 endgroup;
//    ## a=4
//
func (pl *PMMPostParseListener) ExitScalaratom(ctx *grammar.ScalaratomContext) {
	pl.rt.ExprStack.MultiplyTOS2OS()
}

// Attach a (signed) coefficient to a variable, e.g.  +3x, -1/3y.  We just have
// to leave a numeric constant on the stack.
//
//    scalarmulop : (PLUS|MINUS)? numtokenatom
//
// We have to handle the MINUS case only, as the rule for numtokenatom
// already left a numeric constant on the expression stack.
//
func (pl *PMMPostParseListener) ExitScalarmulop(ctx *grammar.ScalarmulopContext) {
	if ctx.MINUS() != nil {
		pl.rt.ExprStack.PushConstant(arithm.MinusOne) // -1 on stack
		pl.rt.ExprStack.MultiplyTOS2OS()              // multiply with numtokenatom
	}
}

// Numeric prefix for a variable, e.g., "3x", "1/2y.r", "0.25z".
//
//    numtokenatom : DECIMALTOKEN '/' DECIMALTOKEN
//                 | DECIMALTOKEN
//
//
func (pl *PMMPostParseListener) ExitNumtokenatom(ctx *grammar.NumtokenatomContext) {
	numbers := ctx.AllDECIMALTOKEN()
	num1, _ := dec.NewFromString(numbers[0].GetText())
	T.P("token", num1.String()).Debugf("numeric token")
	pl.rt.ExprStack.PushConstant(num1) // put decimal number on expression stack
	if len(numbers) > 1 {
		num2, _ := dec.NewFromString(numbers[1].GetText())
		T.P("token", num2.String()).Debugf("numeric token")
		pl.rt.ExprStack.PushConstant(num2)
		pl.rt.ExprStack.DivideTOS2OS()
	}
}

// Start a new scope within expressions.
//
//    atom :  BEGINGROUP statementlist tertiary ENDGROUP     # exprgroup
//
// MetaFont allows "begingroup ... ; <expr> endgroup"  within expressions,
// providing brackets around some statements and returning a (sub-)
// expression.
//
func (pl *PMMPostParseListener) EnterExprgroup(ctx *grammar.ExprgroupContext) {
	groupscope, _ := corelang.Begingroup(pl.rt, "expr-group")
	pl.annotate(ctx, groupscope, "") // Annotate the AST node with this scope
}

// See rule ExitCompound.
//
//    atom :  BEGINGROUP statementlist tertiary ENDGROUP     # exprgroup
//
// Additionally leave the return expression on the stack.
//
func (pl *PMMPostParseListener) ExitExprgroup(ctx *grammar.ExprgroupContext) {
	corelang.Endgroup(pl.rt)
	// the return expression is already on the stack
	//pl.Summary()
}

// ---------------------------------------------------------------------------

// Helper: get the (variable) type of an item of the expression stack
//func (pl *PMMPostParseListener) getExprType(e *runtime.ExprNode) int {
func getExprType(e *runtime.ExprNode) int {
	if e != nil {
		if e.IsValid() {
			if e.IsPair {
				return variables.PairType
			} else {
				/*
					sym := pl.rt.ExprStack.GetVariable(e)
					if sym != nil {
						if t, ok := sym.(runtime.Typed); ok {
							return t.GetType()
						}
					} else if _, ok := e.GetConstNumeric(); ok {
						return variables.NumericType
					}
				*/
				if e.XPolyn.IsValid() {
					return variables.NumericType
				}
			}
		}
		if e.Other != nil { // neither numeric nor pair, now test other types
			var ok bool
			if _, ok = e.Other.(path.Path); ok {
				return variables.PathType
			}
		}
	}
	return variables.Undefined
}

// Create a color from a hex string. Returns black if the string cannot
// be interpreted as a color hex code.
func colorFromHex(hex string) color.Color {
	c, err := colorful.Hex(hex)
	if err == nil {
		return c
	} else {
		return color.Black
	}
}
