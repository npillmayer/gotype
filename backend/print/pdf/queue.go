package pdf

import (
	"container/heap"
	"math"
	"sync"
)

// A heap (priority queue) of pages. The front element is the
// page with the lowest page number.
// pageheap is concurrency safe.
type pageheap struct {
	sync.RWMutex
	pages []*Page
}

// LowestPageNo returns the lowest page number of the pages in the
// pageheap. If the heap is empty, LowestPageNo will return math.MinInt16.
func (h pageheap) LowestPageNo() PageNo {
	h.RLock()
	defer h.RUnlock()
	if len(h.pages) == 0 {
		return math.MinInt16
	}
	return h.pages[0].pageNo
}

// AppendPage pushes a page onto the heap.
func (h *pageheap) AppendPage(page *Page) {
	heap.Push(h, page)
}

// NextPage pops the page with the lowest page number from the heap.
func (h *pageheap) NextPage() *Page {
	if h.Len() == 0 {
		return nil
	}
	return heap.Pop(h).(*Page)
}

// interfacace heap.Interface ----------------------------------------

// Len is the number of elements in the collection.
func (h pageheap) Len() int {
	h.RLock()
	defer h.RUnlock()
	return len(h.pages)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (h pageheap) Less(i, j int) bool {
	h.RLock()
	defer h.RUnlock()
	return h.pages[i].pageNo < h.pages[j].pageNo
}

// Swap swaps the elements with indexes i and j.
func (h pageheap) Swap(i, j int) {
	h.Lock()
	defer h.Unlock()
	h.pages[i], h.pages[j] = h.pages[j], h.pages[i]
}

// Push adds page as element Len()
func (h *pageheap) Push(page interface{}) {
	h.Lock()
	defer h.Unlock()
	h.pages = append(h.pages, page.(*Page))
}

// Pop removes and returs page Len() - 1.
func (h *pageheap) Pop() interface{} { // Len is the number of elements in the collection.
	h.Lock()
	defer h.Unlock()
	old := h.pages
	n := len(old)
	page := old[n-1]
	h.pages = old[0 : n-1]
	return page
}

// ----------------------------------------------------------------------
