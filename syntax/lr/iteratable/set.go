package iteratable

import "sort"

/*
BSD License

Copyright (c) 2017–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */

// Set is a speical purpose set type, suitable mainly for implementing algorithms
// around scanners, parsers, etc. These kinds of algorihms are often more straightforward
// to describe as set constructions and operations.
type Set struct {
	items           []interface{}
	inx             int         // iteration index
	exhaust         bool        // iteration strives to exhaust the set
	stagnationItem  interface{} // detected a Take/Add cycle for this item
	stagnationCount int         // count sequence of Take/Add cycles
	order           Order       // if an order is provided, elements may be sorted
}

// NewSet creates a new iteratable set. Clients may provide initial elements and an
// estimation for the capacity needed.
func NewSet(size int) *Set {
	if size <= 0 {
		size = 10
	}
	s := &Set{
		items: make([]interface{}, 0, size),
	}
	return s
}

// Add adds a new item to a set, if it is not already present.
// If an once-iteration has already read item, and the item is added again to
// the set, the once-iteration will not process the item again.
// For an exhaust-iteration, the added item is put at the end of the list.
func (s *Set) Add(item interface{}) {
	if item == nil {
		return
	}
	if s == nil {
		s = NewSet(-1)
		s.Add(item)
		return
	}
	if s.Contains(item) {
		return
	}
	s.items = append(s.items, item) // grow items-slice
	//fmt.Printf("s.items=%v\n", s.items)
	if s.exhaust { // end of list for exhaustion is position 0
		copy(s.items[1:], s.items[:len(s.items)-1]) // overwrites new item
		s.items[0] = item
		s.inx++
		if item == s.stagnationItem {
			s.stagnationCount++
		} else {
			s.stagnationCount = 0
		}
	}
	//fmt.Printf("added %v, len(S)=%d\n", item, len(s.items))
	return
}

// Remove removes an item from a set, if present. It will return the item if it
// was found, nil otherwise.
func (s *Set) Remove(item interface{}) interface{} {
	if s == nil || item == nil {
		return nil
	}
	var i int
	var m interface{}
	for i, m = range s.items {
		if m == item {
			break
		}
	}
	if m == nil { // item not found
		return nil
	}
	copy(s.items[i:], s.items[i+1:])
	s.items[len(s.items)-1] = nil
	s.items = s.items[:len(s.items)-1]
	if i < s.inx {
		s.inx--
	}
	return m
}

// Contains returns true, if item is contained in the set, false otherwise.
func (s *Set) Contains(item interface{}) bool {
	if s == nil || item == nil {
		return false
	}
	for _, m := range s.items {
		if m == item {
			return true
		}
	}
	return false
}

// Size returns the number of items in the set.
func (s *Set) Size() int {
	if s == nil {
		return 0
	}
	return len(s.items)
}

// Empty returns true, if the set contains items, false otherwise.
func (s *Set) Empty() bool {
	return s.Size() == 0
}

// Equals returns true if both sets contain the same elements.
func (s *Set) Equals(other *Set) bool {
	for _, m := range other.items {
		if !s.Contains(m) {
			return false
		}
	}
	return s.Size() == other.Size()
}

// Union merges the elements of two sets.
func (s *Set) Union(other *Set) *Set {
	if other == nil {
		return s
	}
	if s == nil {
		s = other
		return s
	}
	for _, m := range other.items {
		s.Add(m)
	}
	return s
}

// Intersection returns the common elements of two given sets.
func (s *Set) Intersection(other *Set) *Set {
	if other == nil {
		return s
	}
	if s == nil {
		s = other
		return s
	}
	for _, item := range s.items {
		if !other.Contains(item) {
			s.Remove(item)
		}
	}
	return nil
}

// Difference returns s without the elements contained in other, too.
func (s *Set) Difference(other *Set) *Set {
	if s == nil || other == nil {
		return s
	}
	for _, item := range s.items {
		if other.Contains(item) {
			s.Remove(item)
		}
	}
	return s
}

// Subset removes all elements of a set not fulfilling a condition.
// The condition is given as a boolean function.
// If it is nil, s is returned unchanged.
func (s *Set) Subset(predicate func(interface{}) bool) *Set {
	if s == nil || len(s.items) == 0 {
		return NewSet(-1)
	}
	if predicate != nil {
		for _, item := range s.Values() {
			if !predicate(item) {
				s.Remove(item)
			}
		}
	}
	return s
}

// Each applies a mapper function to each element in the set.
func (s *Set) Each(mapper func(interface{})) {
	if s == nil || len(s.items) == 0 || mapper == nil {
		return
	}
	for _, m := range s.items {
		mapper(m)
	}
}

// Values returns all items of a set.
func (s *Set) Values() []interface{} {
	if s == nil || len(s.items) == 0 {
		return []interface{}{}
	}
	r := make([]interface{}, len(s.items))
	copy(r, s.items)
	return r
}

// Copy makes a copy of a set.
func (s *Set) Copy() *Set {
	if s == nil {
		return nil
	}
	r := NewSet(len(s.items))
	if len(s.items) > 0 {
		r.items = append(r.items, s.items...)
	}
	return r
}

// First returns the first (random) element of a set, or nil.
func (s *Set) First() interface{} {
	if s == nil || len(s.items) == 0 {
		return nil
	}
	return s.items[0]
}

// FirstMatch is a more efficient shortcut for S.Copy().Subset(...).First(),
// i.e. it finds a random element in S with a given condition.
// The condition is given as a boolean function.
// If it is nil, it is treated as "any".
func (s *Set) FirstMatch(predicate func(interface{}) bool) interface{} {
	if s == nil || len(s.items) == 0 {
		return nil
	}
	for _, item := range s.items {
		if predicate == nil || predicate(item) {
			return item
		}
	}
	return nil
}

// --- Iteration --------------------------------------------------------

// Iterate sets up a set for iteration.
//
// We do not use a separate iterator type, thus effectively locking a set to
// a single client at a time. However, for succinct implementations of certain
// algorithms, this way it's much more straightfoward to read and understand.
// Moreover, we allow for modification of the set while it is being iterated.
//
// So, with
//
//     var S iteratable.Set{...}
//
// we could write iterator code something like this:
//
//     iter := S.Iterator()
//     for iter.Next() {
//          item := iter.Item()
//     }
//
// but instead, we rather write
//
//     S.InterateOnce()                 S.Exhaust()
//     for S.Next() {                   for S.Next() {
//          item := S.Item()                 item := S.Take()
//     }                                }
//
// This may seem like a neglectable difference, even violating established OO practice
// (and in a certain way it is), but some algorithms are easier to follow without
// having iterators all over the place, which are used within a single function anyway.
//
// If exhaust is set, the iteration will return items until no more items are
// in the set, otherwise it will iterate over each item once.
func (s *Set) Iterate(exhaust bool) {
	if exhaust {
		s.exhaust = true
		s.inx = len(s.items)
	} else {
		s.inx = -1
	}
}

// IterateOnce is a shortcut for Iterate(false).
// Will iterate over each item once.
func (s *Set) IterateOnce() {
	s.Iterate(false)
}

// Exhaust is a shortcut for Iterate(true).
// The iteration will return items until no more items are in the set.
func (s *Set) Exhaust() {
	s.Iterate(true)
}

// If a set is set up for exhaustion, the iteration will return items until no more items are
// in the set, otherwise it will iterate over each item once.
const (
	Once         bool = false
	Exhaustively bool = true
)

// Next moves the iteration cursor by one position.
// It returns true, if an item is available at the new position.
func (s *Set) Next() bool {
	if s.exhaust {
		s.inx--
		if s.inx < 0 {
			return false
		}
	} else {
		s.inx++
		if s.inx >= len(s.items) {
			return false
		}
	}
	return true
}

// Item returns the item at the currenct iteration position, if any.
func (s *Set) Item() interface{} {
	if s.inx >= 0 && s.inx < len(s.items) {
		return s.items[s.inx]
	}
	return nil
}

// Take returns the item at the current iteration position, if any.
// The item is removed from the set.
func (s *Set) Take() interface{} {
	item := s.Item()
	if item != nil {
		if s.inx < len(s.items)-1 {
			copy(s.items[s.inx:], s.items[s.inx+1:])
		}
		s.items[len(s.items)-1] = nil
		s.items = s.items[:len(s.items)-1]
	}
	s.stagnationItem = item
	return item
}

// Stagnates signals that every item in the set has been scheduled at least once
// for re-consideration.
func (s *Set) Stagnates() bool {
	return s.stagnationCount == len(s.items)
}

// --- Sorting ----------------------------------------------------------------

// Order is a comparator used for sorting the elements of a set.
// It should return true, if x < y for a given interpretation of '<'.
type Order func(x, y interface{}) bool

// Sort sorts the elements of a set according to a given order.
// It returns the input set with elements sorted (not a copy).
func (s *Set) Sort(order Order) *Set {
	if s == nil || len(s.items) == 0 || order == nil {
		return s
	}
	oldOrder := s.order
	s.order = order
	sort.Sort(s)
	s.order = oldOrder
	return s
}

// Len is for sort.interface; not intended for public use.
func (s *Set) Len() int {
	return s.Size()
}

// Swap is for sort.interface; not intended for public use.
func (s *Set) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Swap is for sort.interface; not intended for public use.
func (s *Set) Less(i, j int) bool {
	return s.order(s.items[i], s.items[j])
}

// --- Helpers ----------------------------------------------------------

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
