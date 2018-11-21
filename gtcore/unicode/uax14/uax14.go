package unicode

//

import (
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/emirpasic/gods/lists/arraylist"
	u "github.com/npillmayer/gotype/gtcore/unicode"
	"golang.org/x/text/unicode/rangetable"
)

/*
---------------------------------------------------------------------------

BSD License
Copyright (c) 2017-18, Norbert Pillmayer

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.

3. Neither the name of Norbert Pillmayer nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----------------------------------------------------------------------

/* UAX#14 is the Unicode Annex for Line Breaking (Line Wrap).
 * It defines a bunch of code-point classes and a set of rules
 * for how to place break points / break inhibitors.
*/

// === Constants and Code-Point Ranges ==================================

// Type for UAX#14 code-point classes
// Must be convertable to int
type UAX14Class int

//go:generate stringer -type=UAX14Class
const (
	// sequence matters!
	CMClass UAX14Class = iota
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
	sot       // pseudo class
	eot       // pseudo class
	optSpaces // pseudo class
)

// Range tables for UAX#14 code-point classes.
// Will be initialized with SetupUAX14Classes().
// Clients can check with unicode.Is(..., rune)
var CM, BA, LF, BK, CR, NL, SP, EX, QU, AL, PR,
	PO, OP, CP, IS, HY, SY, NU, AI, BB, GL, SA,
	JT, JV, JL, NS, ZW, ZWJ, WJ, CL, ID, CJ, H2,
	H3, EB, EM, HL, RI, XX *unicode.RangeTable

var uax14ClassFromString = map[string]UAX14Class{
	"CM": CMClass, "BA": BAClass, "LF": LFClass, "BK": BKClass, "CR": CRClass,
	"NL": NLClass, "SP": SPClass, "EX": EXClass, "QU": QUClass, "AL": ALClass,
	"PR": PRClass, "PO": POClass, "OP": OPClass, "CP": CPClass, "IS": ISClass,
	"HY": HYClass, "SY": SYClass, "NU": NUClass, "AI": AIClass, "BB": BBClass,
	"GL": GLClass, "SA": SAClass, "JT": JTClass, "JV": JVClass, "JL": JLClass,
	"NS": NSClass, "ZW": ZWClass, "ZWJ": ZWJClass, "WJ": WJClass, "CL": CLClass,
	"ID": IDClass, "CJ": CJClass, "H2": H2Class, "H3": H3Class, "EB": EBClass,
	"EM": EMClass, "HL": HLClass, "RI": RIClass, "XX": XXClass,
}

// Will be initialized in SetupUAX14Classes()
var rangeFromUAX14Class []*unicode.RangeTable

// Top-level client function:
// Get the line breaking/wrap class for a Unicode code-point
func UAX14ClassForRune(r rune) UAX14Class {
	if r == rune(0) {
		return eot
	}
	for lbc := CMClass; lbc < XXClass; lbc++ {
		urange := rangeFromUAX14Class[lbc]
		if urange == nil {
			fmt.Printf("-- no range for class %s\n", lbc)
		} else if unicode.Is(urange, r) {
			return lbc
		}
	}
	return XXClass
}

// Top-level preparation function:
// Create 'constants' for UAX#14 line breaking/wrap.
func SetupUAX14Classes() error {
	defer timeTrack(time.Now(), "setup of line breaking classes")
	if rangeFromUAX14Class != nil {
		return nil
	}
	lbcs, err := loadUnicodeLineBreakFile()
	if err != nil {
		return err
	}
	createRangesForClassesGlobal(lbcs)
	rangeFromUAX14Class = []*unicode.RangeTable{
		// sequence matters!
		CM, BA, LF, BK, CR, NL, SP, EX, QU, AL, PR,
		PO, OP, CP, IS, HY, SY, NU, AI, BB, GL, SA,
		JT, JV, JL, NS, ZW, ZWJ, WJ, CL, ID, CJ, H2,
		H3, EB, EM, HL, RI, XX,
	}
	return nil
}

// Load the Unicode UAX#14 definition file: LineBreak.txt
func loadUnicodeLineBreakFile() ([]*arraylist.List, error) {
	f, err := os.Open("/Users/npi/prg/go/gotype/etc/LineBreak.txt")
	if err != nil {
		//fmt.Printf("ERROR loading LineBreak.txt\n")
		return nil, err
	}
	defer f.Close()
	p := u.NewUCDParser(f)
	lbcs := make([]*arraylist.List, XXClass+1)
	i := 0
	for p.Next() {
		from, to := p.Range(0)
		brclzstr := p.String(1)
		brclz := uax14ClassFromString[brclzstr]
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
	ZWJ = createRangeTableFor(ZWJClass, lbcs)
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

func createRangeTableFor(lbc UAX14Class, lbcs []*arraylist.List) *unicode.RangeTable {
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

// === UAX#14 Rules =====================================================

type UAX14LineWrap struct {
	publisher    u.RunePublisher
	longestMatch int
	penalties    []int
	rules        map[UAX14Class][]*u.Recognizer // TODO map of StateFn s
}

func NewUAX14LineWrap() *UAX14LineWrap {
	uax14 := &UAX14LineWrap{}
	uax14.rules = map[UAX14Class][]*u.Recognizer{
		NLClass:   {u.NewRecognizer(int(NLClass), 2, []int{-1000}, UAX14NewLine)}, // TODO: create new one for every StartRulesFor(...), even better from pool
		QUClass:   {u.NewRecognizer(int(QUClass), 2, []int{1000}, UAX14Quote)},    // just keep StateFn here
		optSpaces: {u.NewRecognizer(int(SPClass), 1, nil, UAX14OptSpaces)},
	}
	return uax14
}

func (uax14 *UAX14LineWrap) InitFor(rpub u.RunePublisher) {
	uax14.publisher = rpub
}

func (uax14 *UAX14LineWrap) CodePointClassFor(r rune) int {
	return int(UAX14ClassForRune(r))
}

func (uax14 *UAX14LineWrap) StartRulesFor(r rune, cpClass int) {
	uax14c := UAX14Class(cpClass)
	rules := uax14.rules[uax14c]
	if len(rules) > 0 {
		fmt.Printf("starting rules for class = %s\n", uax14c)
	}
	for _, rule := range rules {
		uax14.publisher.SubscribeMe(rule)
	}
}

func (uax14 *UAX14LineWrap) ProceedWithRune(r rune, cpClass int) {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("proceeding with rune = %+q / %s\n", r, uax14c)
	uax14.longestMatch, uax14.penalties = uax14.publisher.PublishRuneEvent(r, int(uax14c))
	fmt.Printf("longest match = %d\n", uax14.LongestMatch())
}

func (uax14 *UAX14LineWrap) LongestMatch() int {
	return uax14.longestMatch
}

func (uax14 *UAX14LineWrap) Penalties() []int {
	return uax14.penalties
}

// --- Recognizer Rules -------------------------------------------------

func doAbort(rec *u.Recognizer) u.NfaStateFn {
	rec.DistanceToGo = 0
	rec.MatchLen = 0
	return abort
}

func abort(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	rec.DistanceToGo = 0
	rec.MatchLen = 0
	fmt.Println("-> abort")
	return abort
}

func accept(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("-> accept %s with lookahead = %s\n", UAX14Class(rec.Expect), uax14c)
	rec.DistanceToGo = 0
	return abort
}

func UAX14NewLine(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("fire rule NewLine for lookahead = %s\n", uax14c)
	if uax14c == NLClass || uax14c == LFClass {
		rec.DistanceToGo = 1
		rec.MatchLen++
		return accept
	}
	return doAbort(rec)
}

func UAX14Quote(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("fire rule Quote for lookahead = %s\n", uax14c)
	if uax14c == QUClass {
		rec.DistanceToGo = 2
		rec.MatchLen++
		rec.Expect = int(OPClass)
		return UAX14OptSpaces
	}
	return doAbort(rec)
}

func UAX14ThisClass(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("fire generic rule for %s with lookahead = %s\n", UAX14Class(rec.Expect), uax14c)
	if uax14c == UAX14Class(rec.Expect) {
		rec.DistanceToGo--
		rec.MatchLen++
		return accept
	}
	return doAbort(rec)
}

func UAX14OptSpaces(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	if uax14c == SPClass {
		fmt.Println("ignoring optional space")
		rec.MatchLen++
		return UAX14OptSpaces // repeat
	} else if uax14c == UAX14Class(rec.Expect) {
		rec.MatchLen++
		return accept // accept rec.Expect in next rec
	}
	return doAbort(rec)
}

// --- Util -------------------------------------------------------------

// Little helper for testing
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf(">>> timing: %s took %s\n", name, elapsed)
}
