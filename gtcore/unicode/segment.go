package unicode

import (
	"bytes"
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
	Penalties() []int
}

type atom struct {
	r       rune
	length  int
	penalty int
}

var eotAtom atom = atom{rune(0), 0, 0}

func (a *atom) String() string {
	return fmt.Sprintf("[%+q p=%d]", a.r, a.penalty)
}

type Segmenter struct {
	deque                 *deque.Deque
	publisher             RunePublisher
	breakers              []UnicodeBreaker
	reader                io.RuneReader
	nullPenalty           int
	longestActiveMatch    int
	distanceToNextPenalty int
	atEOF                 bool
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
	s.distanceToNextPenalty = 10000
}

func (s *Segmenter) SetNullPenalty(null int) {
	s.nullPenalty = null
}

func (s *Segmenter) Next() ([]byte, int, error) {
	if s.reader == nil {
		return nil, 0, errors.New("segmenter not initialized: no input; must call Init()")
	}
	var match []byte
	l := 0
	err := s.readEnoughInput()
	if err != nil && err != io.EOF {
		return nil, 0, err
	}
	if s.distanceToNextPenalty < 1000 {
		fmt.Printf("distance to next penalty (%d) = %d\n", s.deque.At(s.distanceToNextPenalty-1), s.distanceToNextPenalty)
	}
	match, l = s.getTrailSegment()
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
	if s.atEOF {
		err = io.EOF
	} else {
		var r rune
		r, _, err = s.reader.ReadRune()
		fmt.Printf("rune = %+q\n", r)
		if err == nil {
			a := &atom{} // TODO get from pool
			a.r = r
			a.length = 1
			s.deque.PushFront(a)
		} else if err == io.EOF {
			s.deque.PushFront(&eotAtom)
			s.atEOF = true
			err = nil
		} else {
			fmt.Printf("ReadRune() returned error = %s\n", err)
			s.atEOF = true
		}
	}
	return err
}

func (s *Segmenter) readEnoughInput() (err error) {
	//for i := 0; s.deque.Len()-s.longestActiveMatch <= 0; {
	for i := 0; s.distanceToNextPenalty+s.longestActiveMatch >= s.deque.Len(); {
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
				s.insertPenalties(breaker.Penalties())
			}
			s.printQ()
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

func (s *Segmenter) insertPenalties(penalties []int) {
	l := s.deque.Len()
	for i, p := range penalties {
		atom := s.deque.At(i + 1).(*atom)
		atom.penalty += p
		fmt.Printf("l = %d, i = %d\n", l, i)
		if atom.penalty != s.nullPenalty && l-(i+1) < s.distanceToNextPenalty {
			s.distanceToNextPenalty = l - (i + 1)
			fmt.Printf("=> distance to penalty = %d\n", s.distanceToNextPenalty)
		}
	}
}

func (s *Segmenter) getTrailSegment() ([]byte, int) {
	fmt.Println("Collecting trail match:")
	buf := new(bytes.Buffer)
	l := s.deque.Len() - 1 // index of last
	var last *atom
	seglen := 0
	for i := 0; i < s.distanceToNextPenalty; i++ {
		fmt.Printf(".at[%d] = %s\n", l-i, s.deque.At(l-i))
		last = s.deque.PopBack().(*atom)
		written, _ := buf.WriteRune(last.r)
		seglen += written
	}
	s.distanceToNextPenalty = 10000
	for i := 0; i < s.deque.Len(); i++ {
		atom := s.deque.At(i).(*atom)
		if atom.penalty != s.nullPenalty {
			s.distanceToNextPenalty = s.deque.Len() - i // TODO may be off by 1
		}
	}
	fmt.Printf("new distance to penalty = %d\n", s.distanceToNextPenalty)
	return buf.Bytes(), seglen
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
	Penalties    []int
	nextStep     NfaStateFn
}

func NewRecognizer(codePointClass int, distance int, p []int, next NfaStateFn) *Recognizer {
	rec := &Recognizer{}
	rec.Expect = codePointClass
	rec.DistanceToGo = distance
	rec.Penalties = p
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

func (rec *Recognizer) RuneEvent(r rune, codePointClass int) []int {
	fmt.Printf("received rune event: %+q / %d\n", r, codePointClass)
	if rec.nextStep == nil {
		rec.DistanceToGo = 0
	} else {
		rec.nextStep = rec.nextStep(rec, r, codePointClass)
		if rec.MatchLen > 0 && rec.DistanceToGo == 0 { // accepted a match
			return rec.Penalties
		}
	}
	return nil
}

// ----------------------------------------------------------------------

type RuneSubscriber interface {
	RuneEvent(r rune, codePointClass int) []int
	Distance() int
	MatchLength() int
}

type RunePublisher interface {
	SubscribeMe(RuneSubscriber) RunePublisher
	PublishRuneEvent(r rune, codePointClass int) (longestDistance int, penalties []int)
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

func (rpub *runePublisherHeap) PublishRuneEvent(r rune, codePointClass int) (int, []int) {
	longest := 0
	it := rpub.pqueue.Iterator()
	penaltiesTotal := make([]int, 0, 2)
	for it.Next() {
		subscr := it.Value().(RuneSubscriber)
		penalties := subscr.RuneEvent(r, codePointClass)
		penaltiesTotal = AddPenalties(penaltiesTotal, penalties)
		d := subscr.MatchLength()
		if d > longest {
			longest = d
		}
	}
	fmt.Printf("lowest distance = %d\n", rpub.GetLowestDistance())
	for rpub.GetLowestDistance() == 0 {
		rpub.Pop() // drop all subscribers with distance == 0
	}
	return longest, penaltiesTotal
}

// The default function to aggregate break-penalties.
// Simply adds up all penalties.
func AddPenalties(total []int, penalties []int) []int {
	for i, p := range penalties {
		if i >= len(total) {
			total = append(total, p)
		} else {
			total[i] += p
		}
	}
	return total
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
