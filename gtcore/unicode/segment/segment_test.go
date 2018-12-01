package segment

import (
	"testing"
)

// --- ad hoc type for testing purposes----------------------------------

type item struct { // will implement RuneSubscriber
	done bool
}

func (it *item) Done() bool {
	return it.done
}

func (it *item) Unsubscribed()                              {}
func (it *item) RuneEvent(r rune, codePointClass int) []int { return nil }
func (it *item) MatchLength() int                           { return 1 }

// ----------------------------------------------------------------------

func TestQueue1(t *testing.T) {
	pq := &DefaultRunePublisher{}
	if pq.PopDone() != nil {
		t.Error("Should not be able to PopDone() on empty Q")
	}
	pq.Push(&item{done: true})
	if pq.Len() != 1 {
		t.Error("Len() should be 1 for Q with 1 item")
	}
	if !pq.Top().Done() {
		t.Error("single item in Q is done; access by Top() does not reflect this")
	}
	if pq.gap != 0 {
		t.Errorf("gap calculation after Push() to Q is wrong, should be 0, is %d", pq.gap)
	}
	pq.Fix(0)
	if pq.gap != 0 {
		t.Error("gap calculation in Q with length 1 is wrong")
	}
}

func TestQueue2(t *testing.T) {
	pq := &DefaultRunePublisher{}
	it := &item{done: true}
	pq.Push(it)
	pq.Push(&item{done: false})
	if pq.Len() != 2 {
		t.Error("Len() should be 2 for Q with 2 item")
	}
	if !pq.Top().Done() {
		t.Error("top item in Q is not done; should be")
	}
	if pq.gap != 1 {
		t.Errorf("gap calculation after Push()+Push() is wrong, should be 1, is %d", pq.gap)
	}
	it.done = false
	pq.Fix(1)
	if pq.Top().Done() {
		t.Error("top item in Q is done; should not be any more")
	}
	if pq.gap != 2 {
		t.Errorf("gap calculation after Fix() in Q with length 2 is wrong: %d", pq.gap)
	}
}

func TestQueue3(t *testing.T) {
	pq := &DefaultRunePublisher{}
	it := &item{done: false}
	pq.Push(it)
	pq.Push(&item{done: true})
	pq.Push(&item{done: false})
	if pq.Len() != 3 {
		t.Error("Len() should be 3 for Q with 3 item")
	}
	if !pq.Top().Done() {
		t.Error("top item in Q is not done; should be")
	}
	if pq.gap != 2 {
		t.Errorf("gap calculation after 3 x Push() is wrong: %d", pq.gap)
	}
	it.done = true
	pq.Fix(0)
	if pq.gap != 1 {
		t.Errorf("gap calculation after Fix() in Q with length 3 is wrong: %d", pq.gap)
	}
}

func TestQueue4(t *testing.T) {
	pq := &DefaultRunePublisher{}
	it1 := &item{done: false}
	pq.Push(it1)
	it2 := &item{done: false}
	pq.Push(it2)
	it3 := &item{done: false}
	pq.Push(it3)
	it4 := &item{done: false}
	pq.Push(it4)
	if pq.gap != 4 {
		t.Errorf("gap calculation after 4 x Push() is wrong: %d", pq.gap)
	}
	it3.done = true
	pq.Fix(2)
	it2.done = true
	pq.Fix(1)
	it1.done = true
	pq.Fix(0)
	if pq.gap != 1 {
		t.Errorf("gap calculation after Fix() is wrong: %d", pq.gap)
	}
	for j := 0; j < 3; j++ {
		if s := pq.PopDone(); s == nil {
			t.Error("top 3 items should have been done")
		}
	}
	if pq.Top().Done() {
		t.Error("top/only item in Q is done; should not be")
	}
}

/*

func TestRules1(t *testing.T) {
	setupRules()
	fmt.Printf("starts for %s = %d\n", QUClass, ruleStarts[QUClass])
}

func TestRules2(t *testing.T) {
	setupRules()
	fmt.Println("====================")
	fmt.Printf("rules set up\n")
	reader := strings.NewReader("(   (")
	lw := NewLineWrap()
	lw.Init(reader)
	_, n, err := lw.Next()
	fmt.Printf("%d bytes with err = %s\n", n, err)
	lw.printQ()
}

func TestPublishing(t *testing.T) {
	fmt.Println("~~~~~~~~~~~~~~~~~~~~")
	p := NewRunePublisher()
	step1 := newNfaStep(ALClass, 2, nil)
	step2 := newNfaStep(CLClass, 1, nil)
	step3 := newNfaStep(QUClass, 3, nil)
	p.SubscribeMe(step1).SubscribeMe(step2).SubscribeMe(step3)
	ldist := p.PublishRuneEvent('A', int(ALClass))
	p.Print()
	if p.Size() != 2 {
		t.Fail()
	}
	fmt.Printf("longest remaining distance = %d\n", ldist)
	if ldist != 2 {
		t.Fail()
	}
}
*/
