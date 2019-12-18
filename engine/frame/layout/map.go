package layout

import (
	"sync"

	"github.com/npillmayer/gotype/engine/tree"
)

type styleToBoxAssoc sync.Map

func newAssoc() *styleToBoxAssoc {
	return &styleToBoxAssoc{}
}

func (m *styleToBoxAssoc) Put(sn *tree.Node, c *Container) {
	(*sync.Map)(m).Store(sn, c)
}

func (m *styleToBoxAssoc) Get(sn *tree.Node) (*Container, bool) {
	c, ok := (*sync.Map)(m).Load(sn)
	if !ok {
		return nil, false
	}
	return c.(*Container), ok
}
