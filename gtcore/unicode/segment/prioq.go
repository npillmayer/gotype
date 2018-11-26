package segment

import "fmt"

// DefaultRunePblisher is a type to organize RuneSubscribers.
//
// Rune publishers have to maintain a list of subscribers. Subscribers are
// then notified on the arrival of new new runes (code-points) by sending
// them rune-events. When a subscriber is done with consuming runes (subscribers
// are often short-lived), it signals Done()=true.
//
// The DefaultRunePublisher data structure "prioritizes" subscribers with
// Done()=true within a queue.
// It maintains a "gap" position between done and not-done. The queue grows as
// needed.
//
// A DefaultRunePublisher implements RunePublisher
type DefaultRunePublisher struct {
	q   []RuneSubscriber // queue is slice of subscribers
	gap int              // index of first subscriber which is Done(), may be out of range
}

// Number of subscribers held.
func (pq DefaultRunePublisher) Len() int { return len(pq.q) }

func (pq DefaultRunePublisher) empty() bool { return len(pq.q) == 0 }

func (pq DefaultRunePublisher) at(i int) RuneSubscriber {
	return pq.q[i]
}

// Top subscriber in queue. If there is at last one Done() subscriber, top()
// will return one.
func (pq DefaultRunePublisher) Top() RuneSubscriber {
	if pq.Len() == 0 {
		return nil
	}
	return pq.q[pq.Len()-1]
}

// The Done()-flag of a subscriber has changed: inform the queue to let
// it re-organize.
func (pq *DefaultRunePublisher) Fix(at int) {
	if at < pq.Len() {
		if pq.q[at].Done() {
			pq.bubbleUp(at)
		} else {
			pq.bubbleDown(at)
		}
	}
}

// Put a new subscriber into the queue.
func (pq *DefaultRunePublisher) Push(subscr RuneSubscriber) {
	l := pq.Len() // index of new item
	pq.q = append(pq.q, subscr)
	if !pq.Top().Done() {
		pq.bubbleDown(l)
	}
}

// Pop the topmost subscriber.
func (pq *DefaultRunePublisher) Pop() RuneSubscriber {
	old := pq.q
	n := len(old)
	subscr := old[n-1]
	pq.q = old[0 : n-1]
	if pq.gap > pq.Len() {
		pq.gap--
	}
	return subscr
}

// Pop the topmost subscriber if it is Done(), otherwise return nil.
func (pq *DefaultRunePublisher) PopDone() RuneSubscriber {
	if pq.Top().Done() {
		return pq.Pop()
	}
	return nil
}

// Pre-requisite: subscriber at positition is Done().
func (pq *DefaultRunePublisher) bubbleUp(i int) {
	if i < pq.gap {
		if pq.gap < pq.Len() {
			pq.q[i], pq.q[pq.gap] = pq.q[pq.gap], pq.q[i] // swap
		} else if i < pq.Len()-1 { // gap is out of range
			last := pq.Len() - 1
			pq.q[i], pq.q[last] = pq.q[last], pq.q[i] // swap with topmost
		}
		pq.gap--
	}
}

// Pre-requisite: subscriber at positition is not Done().
func (pq *DefaultRunePublisher) bubbleDown(i int) {
	if pq.gap <= i {
		if pq.gap < i {
			pq.q[i], pq.q[pq.gap] = pq.q[pq.gap], pq.q[i] // swap
		}
		pq.gap++
	}
}

func (rpub *DefaultRunePublisher) print() {
	fmt.Printf("Publisher of length %d:\n", rpub.Len())
	for i := rpub.Len() - 1; i >= 0; i++ {
		subscr := rpub.at(i)
		fmt.Printf(" - [%d] %s\n", i, subscr)
	}
}
