package unicode

import (
	"fmt"
	"testing"

	"github.com/npillmayer/gotype/gtcore/unicode"
)

func TestUCDReadLineWrap(t *testing.T) {
	loadUnicodeLineBreakFile()
}

func TestClassForRune1(t *testing.T) {
	SetupUAX14Classes()
	var r rune
	r = 'A'
	c := UAX14ClassForRune(r)
	fmt.Printf("%+q = %s\n", r, c)
}

func TestClassForRune2(t *testing.T) {
	SetupUAX14Classes()
	var r rune
	//r = 'A'
	r = '世'
	c := UAX14ClassForRune(r)
	fmt.Printf("%+q = %s\n", r, c)
}

func TestLineWrapNL(t *testing.T) {
	SetupUAX14Classes()
	publisher := unicode.NewRunePublisher()
	lw := NewUAX14LineWrap()
	lw.InitFor(publisher)
	lw.StartRulesFor('\n', int(NLClass))
	lw.ProceedWithRune('\n', int(NLClass))
	lw.ProceedWithRune('A', int(ALClass))
}

func TestLineWrapQU(t *testing.T) {
	SetupUAX14Classes()
	publisher := unicode.NewRunePublisher()
	lw := NewUAX14LineWrap()
	lw.InitFor(publisher)
	lw.StartRulesFor('"', int(QUClass))
	lw.ProceedWithRune('"', int(QUClass))
	lw.ProceedWithRune(' ', int(SPClass))
	lw.ProceedWithRune('(', int(OPClass))
	lw.ProceedWithRune(' ', int(SPClass))
}
