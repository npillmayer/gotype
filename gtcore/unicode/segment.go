package unicode

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/emirpasic/gods/trees/binaryheap"
	godsutils "github.com/emirpasic/gods/utils"
	"github.com/gammazero/deque"
)

type UnicodeBreaker interface {
	InitFor(RunePublisher)
	CodePointClassFor(rune) int
	StartRulesFor(rune, int)
	ProceedWithRune(rune, int)
	LongestMatch() int
}

type atom struct {
	penalties [10]int
	r         rune
	length    int
}

func (a *atom) String() string {
	return fmt.Sprintf("[%+q #%d]", a.r, a.length)
}

type Segmenter struct {
	deque              *deque.Deque
	publisher          RunePublisher
	breakers           []UnicodeBreaker
	reader             io.RuneReader
	longestActiveMatch int
	atEOF              bool
}

func NewSegmenter(breakers ...UnicodeBreaker) *Segmenter {
	s := &Segmenter{}
	s.breakers = breakers
	return s
}

func (s *Segmenter) Init(reader io.RuneReader) {
	if reader == nil {
		reader = strings.NewReader("")
	}
	s.reader = reader
	if s.deque == nil {
		s.deque = &deque.Deque{}         // Q of atoms
		s.publisher = NewRunePublisher() // for publishing rune events to breakers
		for _, breaker := range s.breakers {
			breaker.InitFor(s.publisher)
		}
	} else {
		s.deque.Clear()
		s.longestActiveMatch = 0
		s.atEOF = false
	}
}

func (s *Segmenter) Next() ([]byte, int, error) {
	if s.reader == nil {
		return nil, 0, errors.New("segmenter not initialized: no input; must call Init()")
	}
	var match []byte
	l := 0
	err := s.readEnoughInput()
	if err != io.EOF {
		return nil, 0, err
	}
	match, l, err = s.getTrailMatch()
	return match, l, err
}

func (s *Segmenter) frontAtom() *atom {
	return s.deque.Front().(*atom)
}

func (s *Segmenter) trailAtom() *atom {
	return s.deque.Back().(*atom)
}

func (s *Segmenter) readRune() (err error) {
	fmt.Println("-------- reading next rune -----------")
	if !s.atEOF {
		var r rune
		r, _, err = s.reader.ReadRune()
		fmt.Printf("rune = %+q\n", r)
		if err == nil {
			a := &atom{}
			a.r = r
			a.length = 1
			s.deque.PushFront(a)
		} else {
			fmt.Printf("ReadRune() returned error = %s\n", err)
			s.atEOF = true
		}
	} else {
		err = io.EOF
	}
	return err
}

func (s *Segmenter) readEnoughInput() (err error) {
	for i := 0; s.deque.Len()-s.longestActiveMatch <= 0; {
		err = s.readRune()
		if err != nil {
			break
		}
		if s.deque.Len() > 0 {
			s.longestActiveMatch = 0
			a := s.frontAtom()
			for _, breaker := range s.breakers {
				cpClass := breaker.CodePointClassFor(a.r)
				breaker.StartRulesFor(a.r, cpClass)
				breaker.ProceedWithRune(a.r, cpClass)
				if breaker.LongestMatch() > s.longestActiveMatch {
					s.longestActiveMatch = breaker.LongestMatch()
				}
			}
		} else {
			fmt.Println("Q empty")
		}
		i++
		if i > 20 { // TODO eliminate (used for debugging purposes only)
			break
		}
	}
	return err
}

func (s *Segmenter) getTrailMatch() ([]byte, int, error) {
	return []byte("abc"), 3, io.EOF
}

func (s *Segmenter) printQ() {
	fmt.Printf("Q #%d: ", s.deque.Len())
	for i := 0; i < s.deque.Len(); i++ {
		fmt.Printf(" - %s", s.deque.At(i))
	}
	fmt.Println(" .")
}

// ----------------------------------------------------------------------

type NfaStateFn func(*Recognizer, rune, int) NfaStateFn

type Recognizer struct {
	Expect       int
	DistanceToGo int
	MatchLen     int
	nextStep     NfaStateFn
}

func NewRecognizer(codePointClass int, distance int, next NfaStateFn) *Recognizer {
	rec := &Recognizer{}
	rec.Expect = codePointClass
	rec.DistanceToGo = distance
	rec.nextStep = next
	return rec
}

func (rec *Recognizer) String() string {
	if rec == nil {
		return "[nil rule]"
	}
	return fmt.Sprintf("[%s -> %d]", rec.Expect, rec.DistanceToGo)
}

func (rec *Recognizer) Distance() int {
	return rec.DistanceToGo
}

func (rec *Recognizer) MatchLength() int {
	return rec.MatchLen
}

func (rec *Recognizer) RuneEvent(r rune, codePointClass int) {
	fmt.Printf("received rune event: %+q / %d\n", r, codePointClass)
	//d := rec.DistanceToGo
	if rec.nextStep == nil {
		rec.DistanceToGo = 0
	} else {
		rec.nextStep = rec.nextStep(rec, r, codePointClass)
	}
}

// ----------------------------------------------------------------------

type RuneSubscriber interface {
	RuneEvent(r rune, codePointClass int)
	Distance() int
	MatchLength() int
}

type RunePublisher interface {
	SubscribeMe(RuneSubscriber) RunePublisher
	PublishRuneEvent(r rune, codePointClass int) (longestDistance int)
	GetLowestDistance() int
}

type runePublisherHeap struct {
	pqueue *binaryheap.Heap
}

func NewRunePublisher() *runePublisherHeap {
	rpub := &runePublisherHeap{}
	rpub.pqueue = binaryheap.NewWith(nfaStepComparator)
	return rpub
}

func (rpub *runePublisherHeap) Peek() RuneSubscriber {
	subscr, _ := rpub.pqueue.Peek()
	return subscr.(RuneSubscriber)
}

func (rpub *runePublisherHeap) Push(subscr RuneSubscriber) {
	rpub.pqueue.Push(subscr)
}

func (rpub *runePublisherHeap) Pop() RuneSubscriber {
	subscr, _ := rpub.pqueue.Pop()
	return subscr.(RuneSubscriber)
}

func (rpub *runePublisherHeap) Empty() bool {
	return rpub.pqueue.Empty()
}

func (rpub *runePublisherHeap) Size() int {
	return rpub.pqueue.Size()
}

func (rpub *runePublisherHeap) PublishRuneEvent(r rune, codePointClass int) int {
	longest := 0
	it := rpub.pqueue.Iterator()
	for it.Next() {
		subscr := it.Value().(RuneSubscriber)
		subscr.RuneEvent(r, codePointClass)
		d := subscr.MatchLength()
		if d > longest {
			longest = d
		}
	}
	fmt.Printf("lowest distance = %d\n", rpub.GetLowestDistance())
	for rpub.GetLowestDistance() == 0 {
		rpub.Pop() // drop all subscribers with distance == 0
	}
	return longest
}

func (rpub *runePublisherHeap) SubscribeMe(rsub RuneSubscriber) RunePublisher {
	rpub.Push(rsub)
	return rpub
}

func (rpub *runePublisherHeap) GetLowestDistance() int {
	if rpub.Empty() {
		return -1
	}
	return rpub.Peek().Distance()
}

func (rpub *runePublisherHeap) Print() {
	fmt.Printf("Publisher of length %d:\n", rpub.Size())
	it := rpub.pqueue.Iterator()
	for it.Next() {
		subscr := it.Value().(RuneSubscriber)
		fmt.Printf(" - %s\n", subscr)
	}
}

func nfaStepComparator(s1, s2 interface{}) int {
	rec1 := s1.(*Recognizer)
	rec2 := s2.(*Recognizer)
	return godsutils.IntComparator(rec1.DistanceToGo, rec2.DistanceToGo) // '<'
}
