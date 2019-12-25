package tree

import (
	"fmt"
	"sync"
)

var errRankOfNullNode = fmt.Errorf("cannot determine rank of null-node")

type rankMap struct {
	lock  *sync.RWMutex
	count map[*Node]uint32
}

func newRankMap() *rankMap {
	return &rankMap{
		&sync.RWMutex{},
		make(map[*Node]uint32),
	}
}

func (rmap *rankMap) Put(n *Node, r uint32) (uint32, error) {
	if n == nil {
		return 0, errRankOfNullNode
	}
	if rmap == nil {
		rmap = newRankMap()
	}
	rmap.lock.RLock()
	rank := rmap.count[n]
	rmap.lock.RUnlock()
	rmap.lock.Lock()
	defer rmap.lock.Unlock()
	rmap.count[n] = r
	return rank, nil
}

func (rmap *rankMap) Get(n *Node) uint32 {
	if n == nil || rmap == nil {
		return 0
	}
	rmap.lock.RLock()
	defer rmap.lock.RUnlock()
	rank := rmap.count[n]
	return rank
}

func (rmap *rankMap) Inc(n *Node) (uint32, error) {
	if n == nil {
		return 0, errRankOfNullNode
	}
	if rmap == nil {
		rmap = newRankMap()
	}
	rmap.lock.Lock()
	defer rmap.lock.Unlock()
	rank := rmap.count[n]
	rmap.count[n] = rank + 1
	return rank, nil
}

func (rmap *rankMap) Clear(n *Node) uint32 {
	if n == nil || rmap == nil {
		return 0
	}
	rmap.lock.Lock()
	defer rmap.lock.Unlock()
	rank := rmap.count[n]
	delete(rmap.count, n)
	return rank
}

// --------------------------------------------------------------------------------

// TODO
// func (npkgs []nodePackage) Len() int {
// 	return
// }
