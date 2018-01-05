package hyphenation

import (
	"fmt"
	"testing"
)

var germanDict, usDict *Dictionnary

func init() {
	germanDict = LoadPatterns("/Users/npi/prg/go/gotype/etc/hyph-de-1996.tex")
	fmt.Printf("%s\n", germanDict.Identifier)
	usDict = LoadPatterns("/Users/npi/prg/go/gotype/etc/hyph-en-us.tex")
	fmt.Printf("%s\n", usDict.Identifier)
}

func TestDEPatterns(t *testing.T) {
	//fmt.Printf("Ausnahme = %s\n", dict.HypenationString("Ausnahme"))
	h := germanDict.HypenationString("Ausnahme")
	if h != "Aus-nah-me" {
		t.Fail()
	}
}

func TestDEPatterns2(t *testing.T) {
	s := germanDict.Hyphenate("Ausnahme")
	t.Logf("Ausnahme = %v (%d)\n", s, len(s))
	if len(s) != 3 || s[0] != "Aus" {
		t.Fail()
	}
}

func TestUSPatterns(t *testing.T) {
	h := usDict.HypenationString("table") // exception dictionnary
	if h != "ta-ble" {
		t.Fail()
	}
	h = usDict.HypenationString("computer")
	if h != "com-put-er" {
		t.Fail()
	}
	h = usDict.HypenationString("algorithm")
	if h != "al-go-rithm" {
		t.Fail()
	}
	h = usDict.HypenationString("concatenation")
	if h != "con-cate-na-tion" {
		t.Fail()
	}
}
