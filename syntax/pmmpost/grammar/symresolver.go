/* Grammar for a poor man's implementation of John Hobby's MetaPost.
 * The lexer and parser code is produced by ANTLR V4.
 *
 * MetaFont's expression syntax relies heavily on semantic checks during
 * parsing. The simplest example is an equation
 *
 *      a = b
 *
 * which may be either of numeric type or pair type. Numeric expressions
 * differ vastly in syntax from pair expressions, thus there are two options:
 *
 * (1) Just treat all expressions the same and apply semantic checks as we
 *     go. This is the way MetaFont itself handles ambiguities. The
 *     advantage of this approach is the ability to throw appropriate
 *     errors whenever we encounter type errors.
 * (2) Introduce a semantic predicate whenever we read a variable reference.
 *     This may lead to strange parser error messages sometimes, but the
 *     resulting grammar is much clearer.
 *
 * This implementation chooses (2). The necessary predicate is defined in this
 * package.
 */
package grammar

import (
	"regexp"

	"github.com/npillmayer/gotype/gtcore/config/tracing"
	"github.com/npillmayer/gotype/syntax/runtime"
)

const ( // variable types
	Undefined = iota
	NumericType
	PairType
	PathType
	ColorType
	PenType
	ComplexArray
	ComplexSuffix
)

var T tracing.Trace = tracing.InterpreterTracer

var ScopeStack *runtime.ScopeTree // to be set from outside (e.g., the parser)

var tagStart *regexp.Regexp = regexp.MustCompile("^[a-z]+") // pattern for tags

/* Semantic predicate for the lexer: Checks for TAGs and MIXEDTAGs if
 * they are declared as being of type pair. Needs the package variable
 * ScopeStack to be set.
 *
 * The method will check a prefix of the terminal using the pattern
 * "^[a-z]+" and feed this string into the scope-stack's symbol resolver.
 */
func ispair(tagname string) bool {
	tagprefix := tagStart.FindString(tagname)
	if tagprefix == "" {
		return false
	} else {
		tag, _ := ScopeStack.Current().ResolveSymbol(tagname)
		//fmt.Printf("tag = %v, %v\n", tag, scope)
		if tag == nil {
			return false
		} else {
			var t runtime.Typable
			var ok bool
			if t, ok = tag.(runtime.Typable); ok {
				ispair := t.GetType() == runtime.PairType
				T.P("tag", tagprefix).Debugf("ispair? = %v", ispair)
				return t.GetType() == runtime.PairType
			}
		}
		return true
	}
}
