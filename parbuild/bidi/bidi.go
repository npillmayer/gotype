package bidi

/*

===========================================================================

Paragraphs:  try to find out global direction / embedding level.
----------

P1: split into paragraphs

     +--P2: iterate over runes, excluding matching isolations, and search for strong R,AL runes
X  --+
	 +--P3: set embedding level to 1, if R,AL found, leave it 0 otherwise

X: may instead be set from outside.

=> global direction D

===========================================================================

- Stack of Triple(embedding level, directional override status, isolate status)
  -----
  Initialized with (D, N, false), with  N=neutral, stack max level 125

- Isolation level count (and for error handling: overflow count)


Rules:
======

X1: Init and iterate over runes, applying X2..X8

X2, X3: LRE rules, skipped for now

X4, X5: explicit overrides RLO and LRO

X5a, b, c: Isolates
-X5a,b: push(L+1, N, true)                 // error handling ignored for now
-X5c: like X, push(l+1, <undecided>, true) // find out natural direction with strong rune

X6: directional overrides affect neutral runes

X6a: PDI close directional isolation level, pop, decrement isolate count

X7: PDE, skipped for now

------------------------------------------------------------


*/

type NfaStateFn func(*Recognizer, rune, int) NfaStateFn

type Recognizer struct {
	MatchLen int        // length of active match
	nextStep NfaStateFn // next step of a DFA
}
