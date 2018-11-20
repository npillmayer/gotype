package unicode

import (
	"fmt"
	"strings"
	"testing"
)

func TestUCDReadLineWrap(t *testing.T) {
	loadUnicodeLineBreakFile()
}

func TestLBC1(t *testing.T) {
	setupLineBreakingClasses()
	var r rune
	r = 'A'
	c := LineBreakingClassForRune(r)
	fmt.Printf("%+q = %s\n", r, c)
}

func TestLBC2(t *testing.T) {
	setupLineBreakingClasses()
	var r rune
	//r = 'A'
	r = '世'
	c := LineBreakingClassForRune(r)
	fmt.Printf("%+q = %s\n", r, c)
}

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
