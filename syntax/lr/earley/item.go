package earley

import "fmt"

// Item is an Earley item.
type Item struct {
	rule   *Rule
	dot    int    // position of dot within rule
	Origin uint64 // start position in input
}

// NullItem is the invalid item.
var NullItem = Item{nil, 0, 0}

// Return an Earley-item for a rule wth the dot at the start of RHS
// func (r *Rule) startItem() (item, Symbol) {
// 	if r.IsEps() {
// 		return item{r, 0}, nil
// 	}
// 	return item{r, 0}, r.rhs[0]
// }

func (i Item) String() string {
	return fmt.Sprintf("%v ➞ %v ● %v", i.rule.LHS, i.rule.rhs[0:i.dot], i.rule.rhs[i.dot:])
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
