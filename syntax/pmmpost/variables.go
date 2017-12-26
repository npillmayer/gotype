/*
BSD License
Copyright (c) 2017, Norbert Pillmayer

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:
1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.
3. Neither the name of Tom Everett nor the names of its contributors
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

package pmmpost

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/npillmayer/gotype/gtcore/arithmetic"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax"
	pmmp "github.com/npillmayer/gotype/syntax/pmmpost/statements"
	"github.com/npillmayer/gotype/syntax/pmmpost/variables"
	dec "github.com/shopspring/decimal"
)

var T tracing.Trace = tracing.InterpreterTracer

/* Variables are complex things in MetaFont/MetaPost. These are legal:

   metafont> showvariable x;
   x=1
   x[]=numeric
   x[][]=numeric
   x[][][]=numeric
   x[][][][]=numeric
   x[]r=numeric
   x[]r[]=numeric
   ...

* Identifier-strings are called "tags". In the example above, 'x' is a tag
* and 'r' is a suffix.
*
* Array variables may be referenced without brackets, if the subscript is just
* a numeric literal, i.e. x[2]r and x2r refer to the same variable. We do
* not rely on the parser to decipher these kinds of variable names for us,
* but rather break up x2r16a => x[2]r[16]a by hand. However, the parser will
* split up array indices in brackets, for the subscript may be a complex expression
* ("x[ypart ((8,5) rotated 20)]" is a valid expression in MetaFont).
* Things are further complicated by the fact that subscripts are allowed to
* be decimals: x[1.2] is valid, and may be typed "x1.2".

   metafont> x[ypart ((8,5) rotated 20)] = 1;
   ## x7.4347=1

* I don't know if this makes sense in practice, but let's try to implement it --
* it might be fun!
*
* I did reject some of MetaFont's conventions, however, for the sake of simlicity:
* Types are inherited from the tag, i.e. if x is of type numeric, then x[2]r is
* of type numeric, too. This is different from MetaFont, where x2r may be of a
* different type than x2. Nevertheless, I'll stick to my interpretation,
* which I find less confusing.
*/

// === Variable Type Declarations ============================================

/* We will handle 5 types of variables and 2 types of suffixes.
 */

/*
const ( // variable types
	Undefined = iota
	pmmp.NumericType
	pmmp.PairType
	PathType
	ColorType
	PenType
	pmmp.ComplexArray
	pmmp.ComplexSuffix
)
*/

// Helper: get a type as string
func typeString(vt int) string {
	switch vt {
	case pmmp.Undefined:
		return "<undefined>"
	case pmmp.NumericType:
		return "numeric"
	case pmmp.PairType:
		return "pair"
	case pmmp.PathType:
		return "path"
	case pmmp.ColorType:
		return "color"
	case pmmp.PenType:
		return "pen"
	case pmmp.ComplexArray:
		return "[]"
	case pmmp.ComplexSuffix:
		return "<suffix>"
	}
	return fmt.Sprintf("<illegal type: %d>", vt)
}

// Helper: get a type from a string
func typeFromString(str string) int {
	switch str {
	case "numeric":
		return pmmp.NumericType
	case "pair":
		return pmmp.PairType
	case "path":
		return pmmp.PathType
	case "color":
		return pmmp.ColorType
	case "pen":
		return pmmp.PenType
	}
	return pmmp.Undefined
}

/* MetaFont declares variables explicitly ("numeric x;") or dynamically
 * ("x=1" => x is of type numeric). Dynamic variable use is permitted for
 * numeric variables only. All other types must be declared. Declaration
 * is for tags only, i.e. the "x" in "x2r". This differs from MetaFont,
 * where x2r can have a separate type from x.
 *
 * We build up doubly-linked tree of variable declarations to describe a
 * variable with a single defining tag. The tag is the entity that goes to
 * the symbol table (of a scope). Suffixes and subscripts are attached to
 * the tag, but invisible as symbols.
 *
 * Example:  numeric x; x2r := 7; x.b := 77;
 * Result:
 *   tag = "x"  of type pmmp.NumericType => into symbol table of a scope
 *    +-- suffix ".b" of type pmmp.ComplexSuffix         "x.b"
 *    +-- subscript "[]" of type pmmp.ComplexArray:      "x[]"
 *         +--- suffix ".r" of type pmmp.ComplexSuffix:  "x[].r"
 */
type PMMPVarDecl struct { // this is a tag, an array-subtype, or a suffix
	syntax.StdSymbol              // to use this we will have to override getName()
	Parent           *PMMPVarDecl // e.g., x <- [] <- suffix(a)
	BaseTag          *PMMPVarDecl // e.g., x // this pointer should never be nil
}

/* Expressive Stringer implementation.
 */
func (d *PMMPVarDecl) String() string {
	return fmt.Sprintf("<decl %s/%s>", d.GetFullName(), typeString(d.GetBaseType()))
}

/* Get isolated name of declaration partial (tag, array or suffix).
 */
func (d *PMMPVarDecl) GetName() string {
	if d.GetType() == pmmp.ComplexArray {
		return "[]"
	} else if d.GetType() == pmmp.ComplexSuffix {
		return "." + d.StdSymbol.GetName()
	} else {
		return d.StdSymbol.GetName()
	}
}

/* Get the full name of a type declaration, starting with the base tag.
 * x <- array <- suffix(a)  gives "x[].a", which is a bit more verbose
 * than MetaFont's response. I prefer this one.
 */
func (d *PMMPVarDecl) GetFullName() string {
	if d.Parent == nil {
		return d.StdSymbol.GetName()
	} else { // we are in a declaration for a complex type
		var s bytes.Buffer
		s.WriteString(d.Parent.GetFullName()) // recursive
		t := d.StdSymbol.GetType()
		if t == pmmp.ComplexArray {
			s.WriteString("[]")
		} else if t == pmmp.ComplexSuffix {
			s.WriteString(".") // MetaFont suppresses this if following an array partial
			s.WriteString(d.StdSymbol.GetName())
		} else {
			panic(fmt.Sprintf("illegal sub-type: %d", t))
		}
		return s.String()
	}
}

/* Returns the type of the base tag.
 */
func (d *PMMPVarDecl) GetBaseType() int {
	return d.BaseTag.StdSymbol.GetType()
}

/* Create and initialize a new variable type declaration.
 * This will be passed as a symbol-creator to the symbol table.
 */
func NewPMMPVarDecl(nm string) syntax.Symbol {
	sym := &PMMPVarDecl{}
	sym.Name = nm
	sym.Symtype = pmmp.Undefined
	sym.BaseTag = sym // this pointer should never be nil
	T.P("decl", sym.GetFullName()).Debug("atomic variable type declaration created")
	return sym
}

/* Convenience function to create and initialize a type declaration.
 * Callers provide a (usually complex) type and an optional parent.
 * If the parent is given and already has a child / suffix-partial with
 * the same signature as the one to create, this function will not create
 * a new partial, but provide the existing one.
 */
func CreatePMMPVarDecl(nm string, tp int, parent *PMMPVarDecl) *PMMPVarDecl {
	if parent != nil { // check if already exists as child of parent
		if parent.GetFirstChild() != nil {
			ch := parent.GetFirstChild().(*PMMPVarDecl)
			for ch != nil { // as long as there are children, i.e. partials
				if (tp == pmmp.ComplexSuffix && ch.GetName() == nm) ||
					(ch.GetType() == pmmp.ComplexArray && tp == pmmp.ComplexArray) {
					T.P("decl", ch.GetFullName()).Debug("variable type already declared")
					return ch // we're done
				}
				if c := ch.GetSibling(); c != nil { // move ch = ch->sibling
					ch = c.(*PMMPVarDecl)
				} else {
					ch = nil
				}
			}
		}
	}
	sym := NewPMMPVarDecl(nm).(*PMMPVarDecl) // not found, create a new one
	sym.Symtype = tp
	T.P("decl", sym.GetFullName()).Debug("variable type declaration created")
	if parent != nil {
		sym.AppendToVarDecl(parent)
	}
	return sym
}

/* Append a complex type partial (suffix or array) to a parent identifier.
 * Will not append the partial, if a partial with this name already
 * exists (as a child).
 */
func (d *PMMPVarDecl) AppendToVarDecl(v *PMMPVarDecl) *PMMPVarDecl {
	if v == nil {
		panic("attempt to append type declaration to nil-tag")
	}
	t := d.StdSymbol.GetType()
	if t != pmmp.ComplexSuffix && t != pmmp.ComplexArray {
		panic(fmt.Sprintf("attempt to append simple type (%d) to tag", t))
	}
	d.BaseTag = v.BaseTag
	d.Parent = v
	v.AppendChild(d)
	T.P("decl", d.GetFullName()).Debugf("new variable type declaration")
	return d
}

/* Show the variable declarations for a tag.
 */
func (d *PMMPVarDecl) showDeclarations(s *bytes.Buffer) *bytes.Buffer {
	if s == nil {
		s = new(bytes.Buffer)
	}
	s.WriteString(fmt.Sprintf("%s : %s\n", d.GetFullName(), typeString(d.BaseTag.GetBaseType())))
	ch := d.GetFirstChild()
	for ; ch != nil; ch = ch.GetSibling() {
		s = ch.(*PMMPVarDecl).showDeclarations(s)
	}
	return s
}

// check interface assignability
var _ syntax.Symbol = &PMMPVarDecl{}
var _ syntax.Typable = &PMMPVarDecl{}
var _ syntax.TreeNode = &PMMPVarDecl{}

// === Variable References / Usage ===========================================

/* Variable reference look like "x", "x2", "hello.world" or "a[4.32].b".
 * Variable references always refer to variable declarations (see code
 * segments above.), which define the type and structure of the variable.
 *
 * The declaration may have partials of type subscript. For every such
 * partial the reference needs a decimal subscript, which we will store
 * in an array of subscripts.
 *
 * Example: x[2.8].b[1] => subscripts = [2.8, 1]
 *
 * Variable references can have a value (of type interface{}).
 */
type PMMPVarRef struct {
	syntax.StdSymbol               // store by normalized name
	cachedName       string        // store full name
	decl             *PMMPVarDecl  // type declaration for this variable
	subscripts       []dec.Decimal // list of subscripts, first to last
	value            interface{}   // if known: has a value (numeric or pair)
}

/* Variables of type pair will use two sub-symbols for the x-part and
 * y-part of the pair respectively. We will connect them using the
 * sibling-link (x-part) and child-link (y-part) of the PMMPVarRef.
 * Both parts link back to the pair variable.
 *
 * We need a different serial ID for the y-part, as it will be used as a
 * variable index in a system of linear equations LEQ.
 */
type PairPartRef struct {
	Id      int         // serial ID
	pairvar *PMMPVarRef // pair parent
	value   interface{} // if known: has a value (numeric)
}

/* Create a variable reference. Low level method.
 */
func CreatePMMPVarRef(decl *PMMPVarDecl, value interface{}, indices []dec.Decimal) *PMMPVarRef {
	if decl.GetType() == pmmp.PairType {
		return CreatePMMPPairTypeVarRef(decl, value, indices)
	} else {
		v := &PMMPVarRef{
			decl:       decl,
			subscripts: indices,
			value:      value,
		}
		v.SetType(decl.GetType())
		v.Id = newVarSerial() // TODO: check, when this is needed (now: id leak)
		//T.Debugf("created var ref: subscripts = %v", indices)
		return v
	}
}

/* Create a pair variable reference. Low level method.
 */
func CreatePMMPPairTypeVarRef(decl *PMMPVarDecl, value interface{}, indices []dec.Decimal) *PMMPVarRef {
	v := &PMMPVarRef{
		decl:       decl,
		subscripts: indices,
		value:      value,
	}
	v.SetType(pmmp.PairType)
	v.Id = newVarSerial() // TODO: check, when this is needed (now: id leak)
	var pair arithmetic.Pair
	var ok bool
	ypart := &PairPartRef{
		Id:      newVarSerial(),
		pairvar: v,
	}
	xpart := &PairPartRef{
		Id:      v.Id,
		pairvar: v,
	}
	v.SetSibling(xpart)
	v.SetFirstChild(ypart)
	if pair, ok = value.(arithmetic.Pair); ok {
		xpart.value = pair.XPart()
		ypart.value = pair.YPart()
	}
	return v
}

/* Symbol-creator for symbol table: creates tag symbol.
 * Do not use this for pair variables !!
 */
func NewPMMPVarRef(tagName string) syntax.Symbol {
	T.P("tag", tagName).Debug("tag for variable reference created")
	v := &PMMPVarRef{}
	v.Id = newVarSerial()
	return v
}

/* Expressive Stringer implementation.
 */
func (v *PMMPVarRef) String() string {
	return fmt.Sprintf("<var %s=%v w/ %s>", v.GetFullName(), v.value, v.decl)
}

/* This method returns the full nomalized name, i.e. "x[2].r".
 * This enables us to store the variable in a symbol table.
 * Overrides GetName() of interface Symbol.
 */
func (v *PMMPVarRef) GetName() string {
	if len(v.cachedName) == 0 {
		v.cachedName = v.GetFullName()
	}
	return v.cachedName
}

// --- Variable Type Pair ----------------------------------------------------

/* Pair parts (x-part or y-part) return the name of their parent pair symbol,
 * prepending "xpart" or "ypart" respectively. This name is constant and
 * may be used to store the pair part in a symbol table.
 */
func (ppart *PairPartRef) GetName() string {
	if ppart.pairvar.GetID() == ppart.GetID() {
		return "xpart " + ppart.pairvar.GetName()
	} else {
		return "ypart " + ppart.pairvar.GetName()
	}
}

/* Returns the serial ID for a pair variable's part.
 */
func (ppart *PairPartRef) GetID() int {
	return ppart.Id
}

/* Predicate: is this variable of type pair?
 */
func (v *PMMPVarRef) IsPair() bool {
	return v.GetType() == pmmp.PairType
}

/* Get the x-part of a pair variable
 */
func (v *PMMPVarRef) XPart() *PairPartRef {
	if !v.IsPair() {
		T.P("var", v.GetName()).Error("cannot access x-part of non-pair")
		return nil
	}
	return v.GetSibling().(*PairPartRef)
}

/* Get the y-part of a pair variable
 */
func (v *PMMPVarRef) YPart() *PairPartRef {
	if !v.IsPair() {
		T.P("var", v.GetName()).Error("cannot access y-part of non-pair")
		return nil
	}
	return v.GetFirstChild().(*PairPartRef)
}

/* Get the x-part value of a pair.
 *
 * Interface syntax.Assignable
 */
func (ppart *PairPartRef) GetValue() interface{} {
	return ppart.value
}

/* Interface syntax.Assignable
 */
func (ppart *PairPartRef) SetValue(val interface{}) {
	T.P("var", ppart.GetName()).Debugf("new value: %v", val)
	ppart.value = val
	ppart.pairvar.PullValue()
}

/* Interface syntax.Assignable
 */
func (ppart *PairPartRef) IsKnown() bool {
	return (ppart.value != nil)
}

/* Filler for interface TreeNode. Never called.
 */
func (ppart *PairPartRef) GetSibling() syntax.TreeNode {
	return nil
}

/* Filler for interface TreeNode. Never called.
 */
func (ppart *PairPartRef) SetSibling(syntax.TreeNode) {
}

/* Filler for interface TreeNode. Never called.
 */
func (ppart *PairPartRef) GetFirstChild() syntax.TreeNode {
	return nil
}

/* Filler for interface TreeNode. Never called.
 */
func (ppart *PairPartRef) SetFirstChild(tn syntax.TreeNode) {
}

/* Get the full normalized (canonical) name of a variable,  i.e. "x[2].r".
 */
func (v *PMMPVarRef) GetFullName() string {
	var suffixes []string
	d := v.decl
	if d == nil {
		return fmt.Sprintf("<undeclared variable: %s>", v.GetName())
	}
	subscriptcount := len(v.subscripts) - 1
	for sfx := v.decl; sfx != nil; sfx = sfx.Parent { // iterate backwards
		//T.Printf("sfx = %v", sfx)
		if sfx.GetType() == pmmp.ComplexArray {
			s := "[" + v.subscripts[subscriptcount].String() + "]"
			suffixes = append(suffixes, s)
			subscriptcount -= 1
		} else {
			suffixes = append(suffixes, sfx.GetName())
		}
	}
	// now reverse the slice of suffixes
	for i := 0; i < len(suffixes)/2; i++ { // swap suffixes in place
		j := len(suffixes) - i - 1
		suffixes[i], suffixes[j] = suffixes[j], suffixes[i]
	}
	return strings.Join(suffixes, "")
}

/* Interface syntax.Assignable
 */
func (v *PMMPVarRef) GetValue() interface{} {
	return v.value
}

/* Interface syntax.Assignable
 */
func (v *PMMPVarRef) SetValue(val interface{}) {
	T.P("var", v.GetName()).Debugf("new value: %v", val)
	v.value = val
	if v.IsPair() {
		var xpart *PairPartRef = v.GetSibling().(*PairPartRef)
		var ypart *PairPartRef = v.GetFirstChild().(*PairPartRef)
		if val == nil {
			xpart.value = nil
			ypart.value = nil
		} else {
			var pairval arithmetic.Pair = val.(arithmetic.Pair)
			xpart.value = pairval.XPart()
			ypart.value = pairval.YPart()
		}
	}
}

/* Whenever a pair part (x-part or y-part) is set, it sends a message to
 * the parent pair variable to pull the value. If both parts are known and
 * a numeric value is set, the parent pair creates a combined pair value.
 */
func (v *PMMPVarRef) PullValue() {
	if v.IsPair() {
		var ppart1, ppart2 *PairPartRef
		ppart1 = v.GetSibling().(*PairPartRef)
		ppart2 = v.GetFirstChild().(*PairPartRef)
		if ppart1 != nil && ppart2 != nil {
			if ppart1.GetValue() != nil && ppart2.GetValue() != nil {
				v.value = &arithmetic.SimplePair{
					X: ppart1.GetValue().(dec.Decimal),
					Y: ppart2.GetValue().(dec.Decimal),
				}
				T.P("var", v.GetName()).Debugf("pair value = %s",
					v.value.(arithmetic.Pair).String())
			}
		}
	}
}

/* Get the value of a variable as a string. For unknown values this method
 * returns "<unknown>".
 */
func (v *PMMPVarRef) ValueString() string {
	if v.IsPair() {
		var xvalue, yvalue string
		xpart := v.XPart().value
		if xpart == nil {
			xvalue = fmt.Sprintf("xpart %s", v.GetName())
		} else {
			xvalue = xpart.(dec.Decimal).String()
		}
		ypart := v.YPart().value
		if ypart == nil {
			yvalue = fmt.Sprintf("ypart %s", v.GetName())
		} else {
			yvalue = ypart.(dec.Decimal).String()
		}
		return fmt.Sprintf("(%s,%s)", xvalue, yvalue)
	} else {
		if v.IsKnown() {
			return v.value.(dec.Decimal).String()
		} else {
			return "<numeric>"
		}
	}
	return "?"
}

/* Interface syntax.Assignable
 */
func (v *PMMPVarRef) IsKnown() bool {
	return (v.value != nil)
}

/* Set a new ID for a variable reference. Whenever variables become
 * re-incarnated, a new serial ID is needed. Re-incarnation happens,
 * whenever a variable goes out of scope, but is still relevant in the
 * LEQ-system. The variables' name continues to live on in a new incarnation,
 * while the out-of-scope variable lives on with the old serial.
 *
 * Returns the old serial ID.
 */
func (v *PMMPVarRef) reincarnate() int {
	oldserial := v.GetID()
	v.Id = newVarSerial()
	if v.IsPair() {
		ypartid := newVarSerial()
		v.XPart().Id = v.Id
		v.YPart().Id = ypartid
	}
	return oldserial
}

var varSerial = 1 // serial number counter for variables, always > 0 !

/* Get a new unique serial ID for variables.
 */
func newVarSerial() int {
	serial := varSerial
	varSerial++
	T.Debugf("creating new serial ID %d", serial)
	return serial
}

// check interface assignability
var _ syntax.Symbol = &PMMPVarRef{}
var _ syntax.Typable = &PMMPVarRef{}
var _ syntax.TreeNode = &PMMPVarRef{}
var _ syntax.Assignable = &PMMPVarRef{}

var _ syntax.Symbol = &PairPartRef{}
var _ syntax.TreeNode = &PairPartRef{}
var _ syntax.Assignable = &PairPartRef{}

// === Variable Parser =======================================================

/* We use a small ANTLR V4 sub-grammar for parsing variable references.
 * We'll attach a listener to ANTLR's AST walker.
 */
type VarParseListener struct {
	*variables.BasePMMPVarListener // build on top of ANTLR's base 'class'
	scopeTree                      *syntax.ScopeTree
	def                            *PMMPVarDecl
	ref                            *PMMPVarRef
	suffixes                       *sll.List
}

/* Construct a new variable parse listener.
 */
func NewVarParseListener() *VarParseListener {
	pl := &VarParseListener{} // no need to initialize base class
	return pl
}

/* Parse a variable from a string. This function will try to find an existing
 * declaration and extend it accordingly. If no declaration can be found, this
 * function will construct a new one from the variable reference.
 *
 * To find a variable in memory we first construct the "canonical" form.
 * Variables live in symbol tables and are resolved by name, so the exact
 * spelling of the variable's name is important. The canonical form uses
 * brackets for subscripts and dots for suffixes.
 *
 * Examples:  x => x, x3 => x[3], z4r => z[4].r, hello.world => hello.world
 */
func (pl *ParseListener) ParseVarFromString(vstr string) antlr.RuleContext {
	// We let ANTLR to the heavy lifting. This may change in the future,
	// as it would be fairly straightforward to implement this by hand.
	input := antlr.NewInputStream(vstr + "@")
	T.Debugf("parsing variable ref = %s", vstr)
	varlexer := variables.NewPMMPVarLexer(input)
	stream := antlr.NewCommonTokenStream(varlexer, 0)
	// TODO: make the parser re-usable....!
	// TODO: We can re-use the parser, but not the lexer (see ANTLR book, chapter 10).
	pl.varParser = variables.NewPMMPVarParser(stream)
	pl.varParser.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	pl.varParser.BuildParseTrees = true
	tree := pl.varParser.Variable()
	sexpr := antlr.TreesStringTree(tree, nil, pl.varParser)
	T.Debugf("### VAR = %s", sexpr)
	return tree
}

/* After having ANTLR create us a parse tree for a variable identifier, we
 * will construct a variable reference from the parse tree. This var ref
 * has an initial serial which is unique. This may not be what you want:
 * usually you will try to find an existing incarnation (with a lower serial)
 * in the memory (see method FindVariableReferenceInMemory).
 *
 * We walk the ANTLR parse tree using a listener (VarParseListener).
 */
func (pl *ParseListener) GetPMMPVarRefFromVarSyntax(vtree antlr.RuleContext) *PMMPVarRef {
	listener := NewVarParseListener()
	listener.scopeTree = pl.interpreter.scopeTree      // variables declarations have to be found
	listener.suffixes = sll.New()                      // list to collect suffixes
	antlr.ParseTreeWalkerDefault.Walk(listener, vtree) // fills listener.ref
	return listener.ref
}

// Helper struct to collect suffixes
type varsuffix struct {
	text   string
	number bool
}

/* Listener callback, receiving a complete variable reference.
 *
 * variable : tag (suffix | subscript)* MARKER
 *
 * A variable has been referenced. We will have to find the declaration of
 * this variable and push a variable reference onto the expression stack.
 * A variable reference looks something like this: "x2a.b" or "y[3.14]".
 * Some complications may arise:
 * - No declaration for the variable's tag can be found: we will create
 *   a declaration for the tag in global scope with type numeric.
 * - The declaration is incomplete, i.e. the tag is declared, but not the
 *   suffix(es). We will extend the declaration appropriately.
 * - We will have to create a subscript vector for the var reference.
 *   we'll collect them (together with suffixes) in a list.
 * Example:
 * Parser reads "x2a", thus tag="x" + subscript="2" + suffix="a"
 * We create (if not yet known) vl.def="x[].a" and vl.ref="x[2].a"
 *
 * The MARKER will be ignored.
 */
func (vl *VarParseListener) ExitVariable(ctx *variables.VariableContext) {
	tag := ctx.Tag().GetText()
	T.P("tag", tag).Debugf("looking for declaration for tag")
	sym, scope := vl.scopeTree.Current().ResolveSymbol(tag)
	if sym != nil {
		vl.def = sym.(*PMMPVarDecl) // scopes are assumed to create these
		T.P("decl", vl.def.GetFullName()).Debugf("found %v in scope %s", vl.def, scope.GetName())
	} else { // variable declaration for tag not found => create it
		sym, _ = vl.scopeTree.Globals().DefineSymbol(tag)
		vl.def = sym.(*PMMPVarDecl)      // scopes are assumed to create these
		vl.def.SetType(pmmp.NumericType) // un-declared variables default to type numeric
		T.P("decl", vl.def.GetName()).Debugf("created %v in global scope", vl.def)
	} // now def declaration of <tag> is in vl.def
	// produce declarations for suffixes, if necessary
	it := vl.suffixes.Iterator()
	subscrCount := 0
	for it.Next() {
		i, vs := it.Index(), it.Value().(varsuffix)
		T.P("decl", vl.def.GetFullName()).Debugf("appending suffix #%d: %s", i, vs)
		if vs.number { // subscript
			vl.def = CreatePMMPVarDecl("<array>", pmmp.ComplexArray, vl.def)
			subscrCount += 1
		} else { // tag suffix
			vl.def = CreatePMMPVarDecl(vs.text, pmmp.ComplexSuffix, vl.def)
		}
	}
	T.P("decl", vl.def.GetFullName()).Debugf("full declared type: %v", vl.def)
	// now create variable ref and push onto expression stack
	var subscripts []dec.Decimal = make([]dec.Decimal, subscrCount, subscrCount+1)
	it = vl.suffixes.Iterator()
	for it.Next() { // extract subscripts -> array
		_, vs2 := it.Index(), it.Value().(varsuffix)
		if vs2.number { // subscript
			d, _ := dec.NewFromString(vs2.text)
			subscripts = append(subscripts, d)
		}
	}
	vl.ref = CreatePMMPVarRef(vl.def, nil, subscripts)
	T.P("var", vl.ref.GetName()).Debugf("var ref %v", vl.ref)
}

/* Variable parsing: Collect a suffix.
 */
func (vl *VarParseListener) ExitSuffix(ctx *variables.SuffixContext) {
	tag := ctx.TAG().GetText()
	T.Debugf("suffix tag: %s", tag)
	vl.suffixes.Add(varsuffix{tag, false})
}

/* Variable parsing: Collect a numeric subscript.
 */
func (vl *VarParseListener) ExitSubscript(ctx *variables.SubscriptContext) {
	d := ctx.DECIMAL().GetText()
	T.Debugf("subscript: %s", d)
	vl.suffixes.Add(varsuffix{d, true})
}
