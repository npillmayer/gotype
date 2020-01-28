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
	return fmt.Sprintf("%v ➞ %v ● %v", i.rule.LHS, i.rule.rhs[0:i.dot], i.rule.rhs[i.dot:])
}

// startItem returns an Earley item from a rule with the dot at position 0.
func startItem(r *Rule) (Item, *Symbol) {
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
func (i Item) Advance() (Item, *Symbol) {
	if i.dot >= len(i.rule.rhs) {
		return NullItem, nil
	}
	ii := Item{i.rule, i.dot + 1, i.Origin}
	return ii, ii.PeekSymbol()
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
		s = s + fmt.Sprintf("%v\\l", i)
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
