package segment

// This is an adaption of an elegant queue data structure by
// Andrew J. Gillis (gammazero). The original is found under
// https://github.com/gammazero/deque :
// Extremely fast ring-buffer deque (double-ended queue) implementation.
//
// I had to substitute interface{} with a concrete struct to avoid
// frequent allocation/dealloation of small objects. I could have used
// a pool, but it's more efficient to rely on the smart memory
// allocation strategy of Deque. Many thanks to gammazero for the great work!
//
// Deque represents a single instance of the deque data structure.
type deque struct {
	buf   []atom
	head  int
	tail  int
	count int
}

// minCapacity is the smallest capacity that deque may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 16

// Len returns the number of elements currently stored in the queue.
func (q *deque) Len() int {
	return q.count
}

// PushBack appends an element to the back of the queue.  Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack().
func (q *deque) PushBack(r rune, penalty int) {
	q.growIfFull()
	q.buf[q.tail].r = r
	q.buf[q.tail].penalty = penalty
	// Calculate new tail position.
	q.tail = q.next(q.tail)
	q.count++
}

// PushFront prepends an element to the front of the queue.
func (q *deque) PushFront(r rune, penalty int) {
	q.growIfFull()

	// Calculate new head position.
	q.head = q.prev(q.head)
	q.buf[q.head].r = r
	q.buf[q.head].penalty = penalty
	q.count++
}

// PopFront removes and returns the element from the front of the queue.
// Implements FIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (q *deque) PopFront() (rune, int) {
	if q.count <= 0 {
		panic("deque: PopFront() called on empty queue")
	}
	r := q.buf[q.head].r
	penalty := q.buf[q.head].penalty
	//q.buf[q.head] = nil
	q.buf[q.head].r = 0
	q.buf[q.head].penalty = 0
	// Calculate new head position.
	q.head = q.next(q.head)
	q.count--

	q.shrinkIfExcess()
	return r, penalty
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (q *deque) PopBack() (rune, int) {
	if q.count <= 0 {
		panic("deque: PopBack() called on empty queue")
	}

	// Calculate new tail position
	q.tail = q.prev(q.tail)

	// Remove value at tail.
	r := q.buf[q.tail].r
	penalty := q.buf[q.tail].penalty
	//q.buf[q.tail] = nil
	q.buf[q.tail].r = 0
	q.buf[q.tail].penalty = 0
	q.count--

	q.shrinkIfExcess()
	return r, penalty
}

// Front returns the element at the front of the queue.  This is the element
// that would be returned by PopFront().  This call panics if the queue is
// empty.
func (q *deque) Front() (rune, int) {
	if q.count <= 0 {
		panic("deque: Front() called when empty")
	}
	r := q.buf[q.head].r
	penalty := q.buf[q.head].penalty
	return r, penalty
}

// Back returns the element at the back of the queue.  This is the element
// that would be returned by PopBack().  This call panics if the queue is
// empty.
func (q *deque) Back() (rune, int) {
	if q.count <= 0 {
		panic("deque: Back() called when empty")
	}
	r := q.buf[q.prev(q.tail)].r
	penalty := q.buf[q.prev(q.tail)].penalty
	return r, penalty
}

// At returns the element at index i in the queue without removing the element
// from the queue.  This method accepts only non-negative index values.  At(0)
// refers to the first element and is the same as Front().  At(Len()-1) refers
// to the last element and is the same as Back().  If the index is invalid, the
// call panics.
//
// The purpose of At is to allow Deque to serve as a more general purpose
// circular buffer, where items are only added to and removed from the the ends
// of the deque, but may be read from any place within the deque.  Consider the
// case of a fixed-size circular log buffer: A new entry is pushed onto one end
// and when full the oldest is popped from the other end.  All the log entries
// in the buffer must be readable without altering the buffer contents.
func (q *deque) At(i int) (rune, int) {
	if i < 0 || i >= q.count {
		panic("deque: At() called with index out of range")
	}
	at := (q.head + i) & (len(q.buf) - 1) // bitwise modulus
	r := q.buf[at].r
	penalty := q.buf[at].penalty
	return r, penalty
}

// SetAt modifies the element at index i in the queue. This method accepts
// only non-negative index values.
func (q *deque) SetAt(i int, r rune, penalty int) {
	if i < 0 || i >= q.count {
		panic("deque: At() called with index out of range")
	}
	at := (q.head + i) & (len(q.buf) - 1) // bitwise modulus
	q.buf[at].r = r
	q.buf[at].penalty = penalty
}

// Clear removes all elements from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at high frequency to avoid
// GC during reuse.  The queue will not be resized smaller as long as items are
// only added.  Only when items are removed is the queue subject to getting
// resized smaller.
func (q *deque) Clear() {
	// bitwise modulus
	modBits := len(q.buf) - 1
	for h := q.head; h != q.tail; h = (h + 1) & modBits {
		q.buf[h].r = 0
		q.buf[h].penalty = 0
	}
	q.head = 0
	q.tail = 0
	q.count = 0
}

// Rotate rotates the deque n steps front-to-back.  If n is negative, rotates
// back-to-front.  Having Deque provide Rotate() avoids resizing that could
// happen if implementing rotation using only Pop and Push methods.
func (q *deque) Rotate(n int) {
	if q.count <= 1 {
		return
	}
	// Rotating a multiple of q.count is same as no rotation.
	n %= q.count
	if n == 0 {
		return
	}

	modBits := len(q.buf) - 1
	// If no empty space in buffer, only move head and tail indexes.
	if q.head == q.tail {
		// Calculate new head and tail using bitwise modulus.
		q.head = (q.head + n) & modBits
		q.tail = (q.tail + n) & modBits
		return
	}

	if n < 0 {
		// Rotate back to front.
		for ; n < 0; n++ {
			// Calculate new head and tail using bitwise modulus.
			q.head = (q.head - 1) & modBits
			q.tail = (q.tail - 1) & modBits
			// Put tail value at head and remove value at tail.
			q.buf[q.head] = q.buf[q.tail]
			q.buf[q.tail].r = 0
			q.buf[q.tail].penalty = 0
		}
		return
	}

	// Rotate front to back.
	for ; n > 0; n-- {
		// Put head value at tail and remove value at head.
		q.buf[q.tail] = q.buf[q.head]
		q.buf[q.head].r = 0
		q.buf[q.head].penalty = 0
		// Calculate new head and tail using bitwise modulus.
		q.head = (q.head + 1) & modBits
		q.tail = (q.tail + 1) & modBits
	}
}

// prev returns the previous buffer position wrapping around buffer.
func (q *deque) prev(i int) int {
	return (i - 1) & (len(q.buf) - 1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (q *deque) next(i int) int {
	return (i + 1) & (len(q.buf) - 1) // bitwise modulus
}

// growIfFull resizes up if the buffer is full.
func (q *deque) growIfFull() {
	if len(q.buf) == 0 {
		q.buf = make([]atom, minCapacity)
		return
	}
	if q.count == len(q.buf) {
		q.resize()
	}
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (q *deque) shrinkIfExcess() {
	if len(q.buf) > minCapacity && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

// resize resizes the deque to fit exactly twice its current contents.  This is
// used to grow the queue when it is full, and also to shrink it when it is
// only a quarter full.
func (q *deque) resize() {
	newBuf := make([]atom, q.count<<1)
	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}
