package uax14

import (
	"fmt"

	u "github.com/npillmayer/gotype/gtcore/unicode"
)

/*
---------------------------------------------------------------------------

BSD License

Copyright (c) 2017-18, Norbert Pillmayer (norbert@pillmayer.com)

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
*/

// These are the UAX#14 line break rules, slightly adapted.
// See http://www.unicode.org/Public/UCD/latest/ucd/auxiliary/LineBreakTest.html
//
// Format: <rule-no>\t<LHS>\t<op>\t<RHS>
//
var rulesLineBreakTests string = `
0.3		÷	eot
4.0	BK	!
5.01	CR	×	LF
5.02	CR	!
5.03	LF	!
5.04	NL	!
6.0		×	( BK | CR | LF | NL )
7.01		×	SP
7.02		×	ZW
8.0	ZW SP*	÷
8.1	ZWJ_O	×
9.0	[^ SP BK CR LF NL ZW]	×	CM
11.01		×	WJ
11.02	WJ	×
12.0	GL	×
12.1	[^ SP BA HY CM]	×	GL
12.2	[^ BA HY CM] CM+	×	GL
12.3	^ CM+	×	GL
13.01		×	EX
13.02	[^ NU CM]	×	(CL | CP | IS | SY)
13.03	[^ NU CM] CM+	×	(CL | CP | IS | SY)
13.04	^ CM+	×	(CL | CP | IS | SY)
14.0	OP SP*	×
15.0	QU SP*	×	OP
16.0	(CL | CP) SP*	×	NS
17.0	B2 SP*	×	B2
18.0	SP	÷
19.01		×	QU
19.02	QU	×
20.01		÷	CB
20.02	CB	÷
21.01		×	BA
21.02		×	HY
21.03		×	NS
21.04	BB	×
21.1	HL (HY | BA)	×
21.2	SY	×	HL
22.01	(AL | HL)	×	IN
22.02	EX	×	IN
22.03	(ID | EB | EM)	×	IN
22.04	IN	×	IN
22.05	NU	×	IN
23.02	(AL | HL)	×	NU
23.03	NU	×	(AL | HL)
23.12	PR	×	(ID | EB | EM)
23.13	(ID | EB | EM)	×	PO
24.02	(PR | PO)	×	(AL | HL)
24.03	(AL | HL)	×	(PR | PO)
25.01	(PR | PO)	×	( OP | HY )? NU
25.02	( OP | HY )	×	NU
25.03	NU	×	(NU | SY | IS)
25.04	NU (NU | SY | IS)*	×	(NU | SY | IS | CL | CP)
25.05	NU (NU | SY | IS)* (CL | CP)?	×	(PO | PR)
26.01	JL	×	JL | JV | H2 | H3
26.02	JV | H2	×	JV | JT
26.03	JT | H3	×	JT
27.01	JL | JV | JT | H2 | H3	×	IN
27.02	JL | JV | JT | H2 | H3	×	PO
27.03	PR	×	JL | JV | JT | H2 | H3
28.0	(AL | HL)	×	(AL | HL)
29.0	IS	×	(AL | HL)
30.01	(AL | HL | NU)	×	OP
30.02	CP	×	(AL | HL | NU)
30.11	^ (RI RI)* RI	×	RI
30.12	[^RI] (RI RI)* RI	×	RI
30.13	RI	÷	RI
30.2	EB	×	EM
`

const (
	GlueBANG  int = -10000
	GlueBREAK int = -500
	GlueJOIN  int = 10000
)

var penaltyForOp = map[string]int{
	"!": GlueBANG,
	"÷": GlueBREAK,
	"×": GlueJOIN,
}

// ----------------------------------------------------------------------

func Rule05_NewLine(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("fire rule 05_NewLine for lookahead = %s\n", uax14c)
	if uax14c == BKClass || uax14c == NLClass || uax14c == LFClass {
		rec.DistanceToGo = 1
		rec.MatchLen++
		return doAccept(rec, GlueBANG, GlueJOIN)
	} else if uax14c == CRClass {
		rec.DistanceToGo = 1
		rec.MatchLen++
		return Rule05_CRLF
	}
	return doAbort(rec)
}

func Rule05_CRLF(rec *u.Recognizer, r rune, cpClass int) u.NfaStateFn {
	uax14c := UAX14Class(cpClass)
	fmt.Printf("fire rule 05_CRLF for lookahead = %s\n", uax14c)
	if uax14c == LFClass {
		return doAccept(rec, GlueBANG, GlueJOIN)
	}
	return doAccept(rec, 0, GlueBANG, GlueJOIN)
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

// ----------------------------------------------------------------------
/*
type rule struct {
	ID  string
	lhs string
	rhs string
	op  int
}

func parseRuleLine(line string) *rule {
	var r *rule
	r = &rule{}
	s := line
	i := strings.IndexByte(s, '\t')
	r.ID = strings.TrimSpace(s[:i])
	s = s[i+1:]
	i = strings.IndexByte(s, '\t')
	r.lhs = strings.TrimSpace(s[:i])
	s = s[i+1:]
	i = strings.IndexByte(s, '\t')
	op, _ := penaltyForOp[strings.TrimSpace(s[:i])]
	r.op = op
	s = s[i+1:]
	r.rhs = strings.TrimSpace(s)
	fmt.Println(r)
	return r
}

func parseRules(rulesStr string) []*rule {
	var err error
	r := strings.NewReader(rulesStr)
	scanner := bufio.NewScanner(r)
	var rules []*rule = make([]*rule, 0, 20)
	for scanner.Scan() && err == nil {
		line := scanner.Text()
		if line == "" {
			continue
		}
		fmt.Printf("new rule line = %s\n", line)
		parseRuleLine(line)
	}
	return rules
}
*/
