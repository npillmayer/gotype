/*
Package sparse implements a simple type for sparse integer matrices.
It is mainly used for parser tables (GOTO-table and ACTION-table).
Every entry in the table is either a single int32 or a pair (int32,int32).

This implementation uses the COO algorithm (a.k.a. triplet-encoding).

   https://medium.com/@jmaxg3/101-ways-to-store-a-sparse-matrix-c7f2bf15a229
   https://www.coin-or.org/Ipopt/documentation/node38.html

----------------------------------------------------------------------

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

3. Neither the name of Norbert Pillmayer or the names of its contributors
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
package sparse

/*
Type for a spare matrix of integer values. Construct with

    M := NewIntMatrix(10, 10, -1)  // last parameter is M's null-value

Now

    M.Set(2, 3, 4711)              // set a value
    v := M.Value(2, 3              // returns 4711
    cnt := M.ValueCount()          // returns 1
    v = M.Value(10, 10)            // returns -1, i.e. the null-value

Values cannot be deleted, but may be overwritten with the null-value. Space for
null-values is not re-claimed.
*/
type IntMatrix struct {
	values  []triplet
	rowcnt  int
	colcnt  int
	nullval int32
}

// Triplet values to store
type triplet struct {
	row, col int
	value    intPair
}

// Create a new matrix for int, size m x n.
func NewIntMatrix(m, n int, nullValue int32) *IntMatrix {
	return &IntMatrix{
		values:  []triplet{},
		rowcnt:  m,
		colcnt:  n,
		nullval: nullValue,
	}
}

// Return row count.
func (m *IntMatrix) M() int {
	return m.rowcnt
}

// Return column count.
func (m *IntMatrix) N() int {
	return m.colcnt
}

// Return this matrix' null value
func (m *IntMatrix) NullValue() int32 {
	return m.nullval
}

// Return the number of values in the matrix.
func (m *IntMatrix) ValueCount() int {
	return len(m.values)
}

// Return the primary value at position (i,j), or NullValue
func (m *IntMatrix) Value(i, j int) int32 {
	for _, t := range m.values {
		if !t.storedLeftOf(i, j) { // have skipped all lesser indices
			if t.storedAt(i, j) {
				return t.value.first()
			}
			break
		}
	}
	return m.nullval
}

// Return the pair of values at position (i,j), or (NullValue, NullValue)
func (m *IntMatrix) Values(i, j int) (int32, int32) {
	for _, t := range m.values {
		if !t.storedLeftOf(i, j) { // have skipped all lesser indices
			if t.storedAt(i, j) {
				return t.value.first(), t.value.second()
			}
			break
		}
	}
	return m.nullval, m.nullval
}

// Set a value in the matrix at position (i,j).
func (m *IntMatrix) Set(i, j int, value int32) *IntMatrix {
	return m.setOrAdd(i, j, value, false)
}

func (m *IntMatrix) Add(i, j int, value int32) *IntMatrix {
	return m.setOrAdd(i, j, value, true)
}

func (m *IntMatrix) setOrAdd(i, j int, value int32, doAdd bool) *IntMatrix {
	at := 0 // will be position of new value
	for k, t := range m.values {
		if !t.storedLeftOf(i, j) { // have skipped all lesser indices
			if t.storedAt(i, j) { // value already present
				if doAdd {
					v := m.values[k].value
					m.values[k].value = addIntValue(v, value, m.nullval) // add new value
				} else {
					m.values[k].value = newIntPair(value, 0) // set new value
				}
				return m // and done
			}
			break // no old value present
		}
		at++
	}
	tnew := triplet{row: i, col: j, value: newIntPair(value, 0)}
	// the following 3 lines have to work for k being the right edge of v or not
	m.values = append(m.values, tnew)    // make room
	copy(m.values[at+1:], m.values[at:]) // copy remainder values one index to right
	m.values[at] = tnew                  // if not append-case: insert new triplet
	return m
}

func addIntValue(v intPair, n int32, nullval int32) intPair {
	if v.first() == nullval {
		return newIntPair(n, 0)
	} else if v.second() == nullval {
		return v.setSecond(n)
	} else {
		// entry is full. what to do?
		return v.setSecond(n) // overwrite second
	}
}

func (t *triplet) storedLeftOf(i, j int) bool {
	return t.row < i || t.row == i && t.col < j
}

func (t *triplet) storedAt(i, j int) bool {
	return (t.row == i && t.col == j)
}

// we will store 2 int32 in one position
type intPair int64

func newIntPair(a, b int32) intPair {
	pr := int64(b)
	pr = pr<<32 + int64(a)
	return intPair(pr)
}

func (pr intPair) first() int32 {
	return int32(pr)
}

func (pr intPair) second() int32 {
	return int32(pr >> 32)
}

func (pr intPair) setFirst(a int32) intPair {
	pr = pr + intPair(int64(a))
	return pr
}

func (pr intPair) setSecond(b int32) intPair {
	pr = pr + intPair(int64(b)<<32)
	return pr
}
