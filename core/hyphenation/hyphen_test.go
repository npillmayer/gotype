package hyphenation

import (
	"fmt"
	"testing"

	"github.com/npillmayer/gotype/core/config/configtestadapter"
	"github.com/npillmayer/gotype/core/config/gconf"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
)

func init() {
	gconf.Initialize(configtestadapter.New())
}

var germanDict, usDict *Dictionary

func init() {
	germanDict = LoadPatterns(gconf.GetString("etc-dir") + "/pattern/hyph-de-1996.tex")
	fmt.Printf("%s\n", germanDict.Identifier)
	usDict = LoadPatterns(gconf.GetString("etc-dir") + "/pattern/hyph-en-us.tex")
	//usDict = LoadPatterns("/Users/npi/prg/go/gotype/etc/hyph-en-us.tex")
	fmt.Printf("%s\n", usDict.Identifier)
}

func TestDEPatterns(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	//fmt.Printf("Ausnahme = %s\n", dict.HyphenationString("Ausnahme"))
	h := germanDict.HyphenationString("Ausnahme")
	if h != "Aus-nah-me" {
		t.Fail()
	}
}

func TestDEPatterns2(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	s := germanDict.Hyphenate("Ausnahme")
	t.Logf("Ausnahme = %v (%d)\n", s, len(s))
	if len(s) != 3 || s[0] != "Aus" {
		t.Fail()
	}
}

func TestUSPatterns(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	h := usDict.HyphenationString("hello")
	if h != "hel-lo" {
		t.Logf("hello should be hel-lo, is %s", h)
		t.Fail()
	}
	h = usDict.HyphenationString("table") // exception dictionary
	if h != "ta-ble" {
		t.Logf("table should be ta-ble, is %s", h)
		t.Fail()
	}
	h = usDict.HyphenationString("computer")
	if h != "com-put-er" {
		t.Logf("computer should be com-put-er, is %s", h)
		t.Fail()
	}
	h = usDict.HyphenationString("algorithm")
	if h != "al-go-rithm" {
		t.Logf("algorithm should be al-go-rithm, is %s", h)
		t.Fail()
	}
	h = usDict.HyphenationString("concatenation")
	if h != "con-cate-na-tion" {
		t.Logf("concatenation should be con-cate-na-tion, is %s", h)
		t.Fail()
	}
	h = usDict.HyphenationString("quick")
	if h != "quick" {
		t.Logf("quick should be quick, is %s", h)
		t.Fail()
	}
	h = usDict.HyphenationString("king")
	if h != "king" {
		t.Logf("king should be king, is %s", h)
		t.Fail()
	}
}
