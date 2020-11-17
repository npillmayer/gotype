package bidi

import (
	"github.com/npillmayer/gotype/syntax/lr/sppf"
	"github.com/npillmayer/gotype/syntax/terex"
)

type BidiTreeOp struct {
	name string
}

func (op *BidiTreeOp) Rewrite(list *terex.GCons, env *terex.Environment) terex.Element {
	T().Debugf(env.Dump())
	return terex.Elem(list)
}

func (op *BidiTreeOp) Descend(sppf.RuleCtxt) bool {
	return true
}

func (op *BidiTreeOp) Name() string {
	return op.name
}

func (op *BidiTreeOp) Operator() terex.Operator {
	return op
}

// String is part of interface terex.Operator
func (op *BidiTreeOp) String() string {
	return op.name
}

// Call is part of interface terex.Operator
func (op *BidiTreeOp) Call(el terex.Element, env *terex.Environment) terex.Element {
	return terex.Elem(nil)
}

// Quote is part of interface terex.Operator
func (op *BidiTreeOp) Quote(el terex.Element, env *terex.Environment) terex.Element {
	return el
}

func newBidiTreeOp(name string) *BidiTreeOp {
	return &BidiTreeOp{
		name: name,
	}
}

var _ terex.Operator = newBidiTreeOp("Hello")
