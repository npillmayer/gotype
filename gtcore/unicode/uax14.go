package unicode

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/trees/binaryheap"
	godsutils "github.com/emirpasic/gods/utils"
	"github.com/gammazero/deque"
	"golang.org/x/text/unicode/rangetable"
)

type LineBreakingClass int

//go:generate stringer -type=LineBreakingClass
const (
	CMClass LineBreakingClass = iota
	BAClass
	LFClass
	BKClass
	CRClass
	NLClass
	SPClass
	EXClass
	QUClass
	ALClass
	PRClass
	POClass
	OPClass
	CPClass
	ISClass
	HYClass
	SYClass
	NUClass
	AIClass
	BBClass
	GLClass
	SAClass
	JTClass
	JVClass
	JLClass
	NSClass
	ZWClass
	ZWJClass
	WJClass
	CLClass
	IDClass
	CJClass
	H2Class
	H3Class
	EBClass
	EMClass
	HLClass
	RIClass
	XXClass
)

var CM, BA, LF, BK, CR, NL, SP, EX, QU, AL, PR, PO,
	OP, CP, IS, HY, SY, NU, AI, BB, GL, SA, JT,
	JV, JL, NS, ZW, ZWJ, WJ, CL, ID, CJ, H2, H3,
	EB, EM, HL, RI, XX *unicode.RangeTable

var lbcFromString = map[string]LineBreakingClass{
	"CM": CMClass, "BA": BAClass, "LF": LFClass, "BK": BKClass, "CR": CRClass, "NL": NLClass,
	"SP": SPClass, "EX": EXClass, "QU": QUClass, "AL": ALClass, "PR": PRClass,
	"PO": POClass, "OP": OPClass, "CP": CPClass, "IS": ISClass, "HY": HYClass,
	"SY": SYClass, "NU": NUClass, "AI": AIClass, "BB": BBClass, "GL": GLClass,
	"SA": SAClass, "JT": JTClass, "JV": JVClass, "JL": JLClass, "NS": NSClass,
	"ZW": ZWClass, "ZWJ": ZWJClass, "WJ": WJClass, "CL": CLClass, "ID": IDClass,
	"CJ": CJClass, "H2": H2Class, "H3": H3Class, "EB": EBClass, "EM": EMClass,
	"HL": HLClass, "RI": RIClass, "XX": XXClass,
}

var rangeFromLineBreakingClass []*unicode.RangeTable

func LineBreakingClassForRune(r rune) LineBreakingClass {
	for lbc := CMClass; lbc < XXClass; lbc++ {
		urange := rangeFromLineBreakingClass[lbc]
		if urange == nil {
			fmt.Printf("-- no range for class %s\n", lbc)
		}
		if unicode.Is(urange, r) {
			return lbc
		}
	}
	return XXClass
}

func setupLineBreakingClasses() error {
	defer timeTrack(time.Now(), "setup of line breaking classes")
	if rangeFromLineBreakingClass != nil {
		return nil
	}
	lbcs, err := loadUnicodeLineBreakFile()
	if err != nil {
		return err
	}
	createRangesForClassesGlobal(lbcs)
	rangeFromLineBreakingClass = []*unicode.RangeTable{
		CM, BA, LF, BK, CR, NL, SP, EX, QU, AL, PR, PO,
		OP, CP, IS, HY, SY, NU, AI, BB, GL, SA, JT,
		JV, JL, NS, ZW, ZWJ, WJ, CL, ID, CJ, H2, H3,
		EB, EM, HL, RI, XX,
	}
	return nil
}

func loadUnicodeLineBreakFile() ([]*arraylist.List, error) {
	f, err := os.Open("/Users/npi/prg/go/gotype/etc/LineBreak.txt")
	if err != nil {
		fmt.Printf("ERROR\n")
		return nil, err
	}
	defer f.Close()
	p := NewUCDParser(f)
	lbcs := make([]*arraylist.List, XXClass+1)
	i := 0
	for p.Next() {
		from, to := p.Range(0)
		brclzstr := p.String(1)
		brclz := lbcFromString[brclzstr]
		list := lbcs[brclz]
		if list == nil {
			list = arraylist.New()
		}
		for r := from; r <= to; r++ {
			list.Add(r)
		}
		lbcs[brclz] = list
		i++
	}
	err = p.Err()
	if err != nil {
		log.Fatal(err)
	}
	/*
		for cl := CMClass; cl <= XXClass; cl++ {
			fmt.Printf("class %s = %v\n", cl, lbcs[cl])
		}
	*/
	return lbcs, err
}

func createRangesForClassesGlobal(lbcs []*arraylist.List) {
	CM = createRangeTableFor(CMClass, lbcs)
	BA = createRangeTableFor(BAClass, lbcs)
	LF = createRangeTableFor(LFClass, lbcs)
	BK = createRangeTableFor(BKClass, lbcs)
	CR = createRangeTableFor(CRClass, lbcs)
	NL = createRangeTableFor(NLClass, lbcs)
	SP = createRangeTableFor(SPClass, lbcs)
	EX = createRangeTableFor(EXClass, lbcs)
	QU = createRangeTableFor(QUClass, lbcs)
	AL = createRangeTableFor(ALClass, lbcs)
	PR = createRangeTableFor(PRClass, lbcs)
	PO = createRangeTableFor(POClass, lbcs)
	OP = createRangeTableFor(OPClass, lbcs)
	CP = createRangeTableFor(CPClass, lbcs)
	IS = createRangeTableFor(ISClass, lbcs)
	HY = createRangeTableFor(HYClass, lbcs)
	SY = createRangeTableFor(SYClass, lbcs)
	NU = createRangeTableFor(NUClass, lbcs)
	AI = createRangeTableFor(AIClass, lbcs)
	BB = createRangeTableFor(BBClass, lbcs)
	GL = createRangeTableFor(GLClass, lbcs)
	SA = createRangeTableFor(SAClass, lbcs)
	JT = createRangeTableFor(JTClass, lbcs)
	JV = createRangeTableFor(JVClass, lbcs)
	JL = createRangeTableFor(JLClass, lbcs)
	NS = createRangeTableFor(NSClass, lbcs)
	ZW = createRangeTableFor(ZWClass, lbcs)
	ZWJ = createRangeTableFor(WJClass, lbcs)
	WJ = createRangeTableFor(WJClass, lbcs)
	CL = createRangeTableFor(CLClass, lbcs)
	ID = createRangeTableFor(IDClass, lbcs)
	CJ = createRangeTableFor(CJClass, lbcs)
	H2 = createRangeTableFor(H2Class, lbcs)
	H3 = createRangeTableFor(H3Class, lbcs)
	EB = createRangeTableFor(EBClass, lbcs)
	EM = createRangeTableFor(EMClass, lbcs)
	HL = createRangeTableFor(HLClass, lbcs)
	RI = createRangeTableFor(RIClass, lbcs)
	XX = createRangeTableFor(XXClass, lbcs)
}

func createRangeTableFor(lbc LineBreakingClass, lbcs []*arraylist.List) *unicode.RangeTable {
	listOfRunes := lbcs[lbc]
	var rtable *unicode.RangeTable
	if listOfRunes != nil {
		runes := make([]rune, listOfRunes.Size())
		bag := listOfRunes.Values()
		for i := 0; i < len(bag); i++ {
			runes[i] = bag[i].(rune)
		}
		rtable = rangetable.New([]rune(runes)...)
	} else {
		rtable = rangetable.New()
	}
	return rtable
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf(">>> timing: %s took %s\n", name, elapsed)
}

type rule struct {
	breakAt       int
	penaltyBefore int
	penaltyAfter  int
	lbcs          []LineBreakingClass
}

func newRule(at int, penB int, penA int, lbcs ...LineBreakingClass) *rule {
	r := &rule{}
	r.breakAt = at
	r.penaltyBefore = penB
	r.penaltyAfter = penA
	r.lbcs = make([]LineBreakingClass, len(lbcs))
	for i, lbc := range lbcs {
		r.lbcs[i] = lbc
	}
	return r
}

func (rule *rule) String() string {
	var b bytes.Buffer
	var i int
	for i = 0; i < rule.breakAt; i++ {
		b.WriteString(" ")
		b.WriteString(rule.lbcs[i].String())
	}
	b.WriteString(" -->")
	for ; i < len(rule.lbcs); i++ {
		b.WriteString(" ")
		b.WriteString(rule.lbcs[i].String())
	}
	return b.String()
}

type ruleset []*rule

var rules ruleset
var ruleStarts map[LineBreakingClass][]int

const (
	optSpaces LineBreakingClass = iota + XXClass
	anyClass
	sot
	eot
)

func setupRules() {
	rulecnt := 3
	rules = make(ruleset, rulecnt)
	rules[0] = newRule(1, 0, 1000, QUClass, optSpaces, OPClass)
	rules[1] = newRule(0, 0, -1000, NLClass)
	rules[2] = newRule(0, 1000, 1000, OPClass)
	ruleStarts = make(map[LineBreakingClass][]int)
	for i := 0; i < rulecnt; i++ {
		rule := rules[i]
		starter := rule.lbcs[0]
		rs := ruleStarts[starter]
		if rs == nil {
			rs = make([]int, 1, 1)
			rs[0] = i
		} else {
			rs = append(rs, i)
		}
		ruleStarts[starter] = rs
	}
}

type rulestate struct {
	rule *rule
	at   int
}

type segment struct {
	lbc       LineBreakingClass
	penalties [10]int
	r         rune
	length    int
}

func (seg *segment) String() string {
	return fmt.Sprintf("[%+q #%d (%s)]", seg.r, seg.length, seg.lbc)
}

type LineWrap struct {
	deque     *deque.Deque
	publisher RunePublisher
	//trailPosDist int
	reader io.RuneReader
	atEOF  bool
}

func NewLineWrap() *LineWrap {
	lw := &LineWrap{}
	lw.deque = &deque.Deque{}
	lw.publisher = NewRunePublisher()
	return lw
}

func (lw *LineWrap) Init(reader io.RuneReader) {
	lw.reader = reader
}

func (lw *LineWrap) Next() ([]byte, int, error) {
	err := lw.readEnoughInput()
	/*
		if bytes, penalty := lw.getSegment(); bytes != nil {
			return bytes, penalty, err
		} else {
			return lw.flush(), 0, err
		}
	*/
	return nil, 0, err
}

func (lw *LineWrap) frontSegment() *segment {
	return lw.deque.Front().(*segment)
}

func (lw *LineWrap) trailSegment() *segment {
	return lw.deque.Back().(*segment)
}

func (lw *LineWrap) readRune() (err error) {
	fmt.Println("--------------------")
	fmt.Printf("reading next rune\n")
	if !lw.atEOF {
		r, _, err := lw.reader.ReadRune()
		fmt.Printf("r = %+q\n", r)
		if err == nil {
			lbc := LineBreakingClassForRune(r)
			if !lw.collectSpace(r, lbc) {
				seg := &segment{}
				seg.r = r
				seg.lbc = lbc
				seg.length = 1
				lw.deque.PushFront(seg)
			}
		} else {
			fmt.Printf("ReadRune() returned error = %s\n", err)
			lw.atEOF = true
		}
	} else {
		err = io.EOF
	}
	return err
}

func (lw *LineWrap) collectSpace(r rune, lbc LineBreakingClass) bool {
	collected := false
	fmt.Printf("checking for space: %+q = %s\n", r, lbc)
	if lw.deque.Len() > 0 {
		if lbc == SPClass {
			front := lw.frontSegment()
			fmt.Printf("front = %s\n", front)
			if lbc == SPClass && front.lbc == SPClass {
				front.length++
				fmt.Printf("|SP| = %d\n", front.length)
				collected = true // TODO: conserve apppended space rune
			}
		}
	}
	return collected
}

func (lw *LineWrap) maxActiveRuleTrailLength() int {
	return 5
}

func (lw *LineWrap) readEnoughInput() (err error) {
	for i := 0; lw.deque.Len()-lw.maxActiveRuleTrailLength() <= 0; {
		err = lw.readRune()
		if err != nil {
			break
		}
		if lw.deque.Len() > 0 {
			seg := lw.frontSegment()
			brclass := LineBreakingClassForRune(seg.r)
			// start up rules with brclass as first step
			inxs := ruleStarts[brclass]
			fmt.Printf("starting rules: %v\n", inxs)
			for inx, r := range inxs {
				fmt.Printf("  rule #%d: %s\n", inx, rules[r])
				rule := rules[r]
				lw.startRule(rule)
			}
		} else {
			fmt.Println("Q empty")
		}
		i++
		if i > 7 {
			break
		}
	}
	return err
}

func (lw *LineWrap) printQ() {
	fmt.Printf("Q #%d: ", lw.deque.Len())
	for i := 0; i < lw.deque.Len(); i++ {
		fmt.Printf(" - %s", lw.deque.At(i))
	}
	fmt.Println(" .")
}

// ----------------------------------------------------------------------

//const nStates int = 5
//const nClasses int = int(XXClass) + 1

//type stateFn func(*lexer) stateFn

/*
func run() {
	for state := startState; state != nil; {
		state = state(lexer)
	}
}
*/

type stateFn func(*nfaStep) *nfaStep

type nfaStep struct {
	expect       LineBreakingClass
	distanceToGo int
	nextStep     stateFn
}

func newNfaStep(lbc LineBreakingClass, distance int, next stateFn) *nfaStep {
	step := &nfaStep{}
	step.expect = lbc
	step.distanceToGo = distance
	step.nextStep = next
	return step
}

func (step *nfaStep) String() string {
	return fmt.Sprintf("[%s -> %d]", step.expect, step.distanceToGo)
}

func (step *nfaStep) Distance() int {
	return step.distanceToGo
}

func (step *nfaStep) RuneEvent(r rune, breakingClass int) {
	d := step.distanceToGo
	if step.distanceToGo > 0 {
		step.distanceToGo--
	}
	fmt.Printf("received rune event; distance = %d -> %d\n", d, step.distanceToGo)
}

// ----------------------------------------------------------------------

type RuneSubscriber interface {
	RuneEvent(r rune, breakingClass int)
	Distance() int
}

type RunePublisher interface {
	SubscribeMe(RuneSubscriber) RunePublisher
	PublishRuneEvent(r rune, breakingClass int) (longestDistance int)
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

func (rpub *runePublisherHeap) PublishRuneEvent(r rune, breakingClass int) int {
	longest := 0
	it := rpub.pqueue.Iterator()
	for it.Next() {
		subscr := it.Value().(RuneSubscriber)
		subscr.RuneEvent(r, breakingClass)
		d := subscr.Distance()
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
	step1 := s1.(*nfaStep)
	step2 := s2.(*nfaStep)
	return godsutils.IntComparator(step1.distanceToGo, step2.distanceToGo) // '<'
}
