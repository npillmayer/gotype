package layout

import (
	"sync"

	"github.com/npillmayer/gotype/engine/dom"
)

type domToBoxAssoc sync.Map

func newAssoc() *domToBoxAssoc {
	return &domToBoxAssoc{}
}

func (m *domToBoxAssoc) Put(domnode *dom.W3CNode, c *Container) {
	(*sync.Map)(m).Store(domnode, c)
}

func (m *domToBoxAssoc) Get(domnode *dom.W3CNode) (*Container, bool) {
	c, ok := (*sync.Map)(m).Load(domnode)
	if !ok {
		return nil, false
	}
	return c.(*Container), ok
}
