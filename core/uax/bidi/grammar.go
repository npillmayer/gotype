package bidi

import (
	"errors"
	"io/ioutil"
	"sync"

	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/sppf"

	"github.com/npillmayer/gotype/syntax/lr"
	"golang.org/x/text/unicode/bidi"
)

var globalBidiGrammar *lr.LRAnalysis
var globalParser *earley.Parser

var initParser sync.Once

func getParser() *earley.Parser {
	initParser.Do(func() {
		globalBidiGrammar = NewBidiGrammar()
		globalBidiGrammar.Grammar().Dump()
		globalParser = earley.NewParser(globalBidiGrammar, earley.GenerateTree(true), earley.StoreTokens(true))
		if globalParser == nil {
			panic("Could not created Bidi grammar parser")
		}
	})
	return globalParser
}

// NewBidiGrammar creates a new grammar for parsing Bidi runs in a paragraph. It is usually
// not called by clients directly, but rather used transparently with a call to Parse.
// It is included in the API for advanced usage, like extending or modifying the grammar.
func NewBidiGrammar() *lr.LRAnalysis {
	b := lr.NewGrammarBuilder("UAX#9")
	b.LHS("Run").T(bc(bidi.LRI)).N("L").T(bc(bidi.PDI)).End()
	b.LHS("Run").T(bc(bidi.LRI)).N("EN").T(bc(bidi.PDI)).End()
	b.LHS("Run").T(bc(bidi.RLI)).N("R").T(bc(bidi.PDI)).End()
	b.LHS("L").T(bc(bidi.L)).End()
	b.LHS("R").T(bc(bidi.R)).End()
	b.LHS("AL").T(bc(bidi.AL)).End()
	b.LHS("EN").T(bc(bidi.EN)).End()
	b.LHS("AN").T(bc(bidi.AN)).End()
	b.LHS("R").N("AL").End()                         // W3
	b.LHS("EN").N("EN").T(bc(bidi.ES)).N("EN").End() // W4
	b.LHS("EN").N("EN").T(bc(bidi.CS)).N("EN").End() // W4
	b.LHS("AN").N("AN").T(bc(bidi.CS)).N("AN").End() // W4
	g, err := b.Grammar()
	if err != nil {
		panic(err)
	}
	return lr.Analysis(g)
}

// Parse parses a paragraph of text, given by a scanner, and parses it according to the
// Unicode UAX#9 Bidi algorithm.
func Parse(scanner *Scanner) (bool, error) {
	if scanner == nil {
		return false, errors.New("Expected parameter scanner to be non-nil")
	}
	parser := getParser()
	T().Debugf("===========================================================")
	accept, err := parser.Parse(scanner, nil)
	if accept {
		parsetree := parser.ParseForest()
		tmpfile, err := ioutil.TempFile(".", "bidi-*.dot")
		if err != nil {
			panic("cannot open tmp file")
		}
		sppf.ToGraphViz(parsetree, tmpfile)
		T().Infof("Exported parse tree to %s", tmpfile.Name())
	}
	return accept, err
}

// --- Helpers ---------------------------------------------------------------

func bc(tokval bidi.Class) (string, int) {
	return ":" + ClassString(tokval), int(tokval)
}
