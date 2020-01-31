package lr

import (
	"fmt"

	"github.com/npillmayer/gotype/syntax/lr/iteratable"
)

// Item is an Earley item.
type Item struct {
	rule   *Rule
	dot    int    // position of dot within rule
	Origin uint64 // start position in input
}

// NullItem is the invalid item.
var NullItem = Item{nil, 0, 0}

func (i Item) String() string {
	return fmt.Sprintf("%v ➞ %v ● %v  (%d)", i.rule.LHS, i.rule.rhs[0:i.dot],
		i.rule.rhs[i.dot:], i.Origin)
}

// Rule returns the grammar rule of this item.
func (i Item) Rule() *Rule {
	return i.rule
}

// StartItem returns an Earley item from a rule with the dot at position 0.
func StartItem(r *Rule) (Item, *Symbol) {
	if r.IsEps() {
		return Item{r, 0, 0}, nil
	}
	return Item{r, 0, 0}, r.rhs[0]
}

// PeekSymbol returns the symbol after the dot, if any.
func (i Item) PeekSymbol() *Symbol {
	if i.dot >= len(i.rule.rhs) {
		return nil
	}
	return i.rule.rhs[i.dot]
}

// Prefix returns a slice, so the result should probably be considered read-only.
// It returns the symbols of the RHS before the dot.
func (i Item) Prefix() []*Symbol {
	return i.rule.rhs[:i.dot]
}

// Advance advances the dot of an item over the next symbol.
// Returns NullItem if the dot is already past the last symbol.
func (i Item) Advance() Item {
	if i.dot >= len(i.rule.rhs) {
		return NullItem
	}
	ii := Item{i.rule, i.dot + 1, i.Origin}
	return ii
}

// --- Item sets --------------------------------------------------------

// A set of Earley items ( A -> B *C D ).

func newItemSet() *iteratable.Set {
	return iteratable.NewSet(0)
}

func asItem(i interface{}) Item {
	if item, ok := i.(Item); ok {
		return item
	}
	panic("not an item")
}

// Prepare an item set for export to Graphviz.
func forGraphviz(iset *iteratable.Set) string {
	items := iset.Values()
	if len(items) == 0 {
		return "err\\n"
	}
	s := ""
	for _, i := range items {
		item := i.(Item)
		s = s + fmt.Sprintf("%v\\n", item)
	}
	return s
}

// Dump an item set to the tracer.
func Dump(iset *iteratable.Set) {
	items := iset.Values()
	//T().Debug("--- item set ------------")
	for k, i := range items {
		T().Debugf("item %2d = %v", k, i)
	}
	//T().Debug("-------------------------")
}

// --- Spans ------------------------------------------------------------

// Span is a small type for capturing a length of input token run. For every
// terminal and non-terminal, a parse tree/forest will track which input positions
// this symbol covers. A span denotes a start position and the position just
// behind the end.
type Span [2]uint64 // (x…y)

// From returns the start value of a span.
func (s *Span) From() uint64 {
	return s[0]
}

// To returns the end value of a span.
func (s *Span) To() uint64 {
	return s[1]
}

// Len returns the length of (x…y)
func (s *Span) Len() uint64 {
	return s[1] - s[0]
}

func (s *Span) String() string {
	return fmt.Sprintf("(%d…%d)", s[0], s[1])
}
