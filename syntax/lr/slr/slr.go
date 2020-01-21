/*
Package slr provides a SLR(1)-parser.
*/
package slr

import (
	"fmt"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/npillmayer/gotype/syntax/lr"
	"github.com/npillmayer/gotype/syntax/lr/sparse"
)

// T traces to the global SyntaxTracer.
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}

// A Token type, if you want to use it. Tokens of this type are returned
// by StdScanner.
//
// Clients may provide their own token data type.
type Token struct {
	Value  int
	Lexeme []byte
}

// Scanner is an interface the parser relies on.
type Scanner interface {
	MoveTo(position uint64)
	NextToken(expected []int) (tokval int, token interface{})
}

// A Parser type. Create and initialize one with slr.Parser(...)
type Parser struct {
	G       *lr.Grammar
	stack   *stack            // parser stack
	gotoT   *sparse.IntMatrix // GOTO table
	actionT *sparse.IntMatrix // ACTION table
	//accepting []int             // slice of accepting states
}

// NewParser creates a SLR(1) parser.
func NewParser(g *lr.Grammar, gotoTable *sparse.IntMatrix, actionTable *sparse.IntMatrix) *Parser {
	//func NewParser(g *lr.Grammar, gotoTable *sparse.IntMatrix, actionTable *sparse.IntMatrix,
	//	acceptingStates []int) *Parser {
	parser := &Parser{
		G:       g,
		stack:   newstack(),
		gotoT:   gotoTable,
		actionT: actionTable,
		//accepting: acceptingStates,
	}
	return parser
}

// Parse startes a new parse, given a start state and a scanner tokenizing the input.
// The parser must have been initialized.
//
// The parser returns true if the input string has been accepted.
func (p *Parser) Parse(S *lr.CFSMState, scan Scanner) (bool, error) {
	T().Debugf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	//T().Debugf("accepting states=%v", p.accepting)
	if p.G == nil {
		T().Errorf("parser not initialized")
		return false, fmt.Errorf("Parser not initialized")
	}
	var accepting bool
	p.stack.Push(S) // push the start state onto the stack
	// http://www.cse.unt.edu/~sweany/CSCE3650/HANDOUTS/LRParseAlg.pdf
	tokval, token := scan.NextToken(nil)
	done := false
	for !done {
		if token == nil {
			tokval = scanner.EOF
		}
		T().Debugf("got token %s/%d from scanner", token, tokval)
		state := p.stack.Peek()
		// if tokval == scanner.EOF {
		// 	done = true
		// }
		T().Debugf("looking for ACTION(%d, %d)", state.ID, tokval)
		action := p.actionT.Value(state.ID, tokval)
		T().Debugf("action = %s", valstring(action, p.actionT))
		if action == p.actionT.NullValue() {
			return false, fmt.Errorf("Syntax error at %d/%v", tokval, token)
		}
		if action == -1 { // shift action
			nextstate := int(p.gotoT.Value(state.ID, tokval))
			T().Debugf("shifting, next state = %d", nextstate)
			//p.stack.Push(&lr.CFSMState{ID: nextstate, Accept: contains(nextstate, p.accepting)})
			p.stack.Push(&lr.CFSMState{ID: nextstate})
			tokval, token = scan.NextToken(nil)
		} else if action > 0 { // reduce action
			nextstate := p.reduce(state.ID, p.G.Rule(int(action)))
			T().Debugf("next state = %d", nextstate)
			//p.stack.Push(&lr.CFSMState{ID: nextstate, Accept: contains(nextstate, p.accepting)})
			p.stack.Push(&lr.CFSMState{ID: nextstate})
		} else if action == -2 {
			//if contains(p.stack.Peek().ID, p.accepting) {
			T().Infof("ACCEPT")
			accepting = true
			done = true
		} else {
			done = true
		}
		T().Debugf("~~~ token %v processed ~~~~~~~~~~~~~~~~~~~~~~", token)
	}
	return accepting, nil
}

func (p *Parser) reduce(stateID int, rule *lr.Rule) int {
	T().Infof("reduce %v", rule)
	handle := reverse(rule.GetRHS())
	for range handle {
		p.stack.Pop()
		// if tos := p.stack.Pop(); sym.GetID() != tos.ID {
		// 	T().Errorf("Expected %v on top of stack, got %d", sym, tos)
		// }
	}
	lhs := rule.GetLHSSymbol()
	state := p.stack.Peek()
	nextstate := p.gotoT.Value(state.ID, lhs.GetID())
	return int(nextstate)
}

// ----------------------------------------------------------------------

type stack struct {
	arrstack *arraystack.Stack
}

func newstack() *stack {
	return &stack{
		arrstack: arraystack.New(),
	}
}

func (s *stack) Peek() *lr.CFSMState {
	state, ok := s.arrstack.Peek()
	if !ok {
		panic("Peek() on empty parser stack")
	}
	return state.(*lr.CFSMState)
}

func (s *stack) Push(state *lr.CFSMState) {
	s.arrstack.Push(state)
	for i, v := range s.arrstack.Values() {
		T().Debugf("[%2d] %d", i, v.(*lr.CFSMState).ID)
	}
}

func (s *stack) Pop() *lr.CFSMState {
	state, ok := s.arrstack.Pop()
	if !ok {
		panic("Pop() on empty parser stack")
	}
	return state.(*lr.CFSMState)
}

func reverse(syms []lr.Symbol) []lr.Symbol {
	r := append([]lr.Symbol(nil), syms...) // make copy first
	for i := len(syms)/2 - 1; i >= 0; i-- {
		opp := len(syms) - 1 - i
		syms[i], syms[opp] = syms[opp], syms[i]
	}
	return r
}

// ----------------------------------------------------------------------

func contains(el int, a []int) bool {
	for _, e := range a {
		if el == e {
			return true
		}
	}
	return false
}

func valstring(v int32, m *sparse.IntMatrix) string {
	if v == m.NullValue() {
		return "<none>"
	}
	return fmt.Sprintf("%d", v)
}
