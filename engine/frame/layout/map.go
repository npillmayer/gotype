package layout

import (
	"sync"

	"github.com/npillmayer/gotype/engine/dom"
)

type domToBoxAssoc struct {
	sync.RWMutex
	m map[*dom.W3CNode]*Container
}

func newAssoc() *domToBoxAssoc {
	return &domToBoxAssoc{
		m: make(map[*dom.W3CNode]*Container),
	}
}

func (d2c *domToBoxAssoc) Put(domnode *dom.W3CNode, c *Container) {
	d2c.Lock()
	defer d2c.Unlock()
	d2c.m[domnode] = c
}

func (d2c *domToBoxAssoc) Get(domnode *dom.W3CNode) (*Container, bool) {
	d2c.RLock()
	defer d2c.RUnlock()
	c, ok := d2c.m[domnode]
	return c, ok
}

func (d2c *domToBoxAssoc) Length() int {
	d2c.RLock()
	defer d2c.RUnlock()
	return len(d2c.m)
}
