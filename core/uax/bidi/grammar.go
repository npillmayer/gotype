package bidi

import (
	"errors"
	"sync"

	"github.com/npillmayer/gotype/syntax/lr/earley"
	"github.com/npillmayer/gotype/syntax/lr/sppf"

	"github.com/npillmayer/gotype/syntax/lr"
	"golang.org/x/text/unicode/bidi"
)

var globalBidiGrammar *lr.LRAnalysis

var initParser sync.Once

func getParser() *earley.Parser {
	initParser.Do(func() {
		globalBidiGrammar = NewBidiGrammar()
		globalBidiGrammar.Grammar().Dump()
	})
	parser := earley.NewParser(globalBidiGrammar, earley.GenerateTree(true), earley.StoreTokens(true))
	if parser == nil {
		panic("Could not created Bidi grammar parser")
	}
	return parser
}

// NewBidiGrammar creates a new grammar for parsing Bidi runs in a paragraph. It is usually
// not called by clients directly, but rather used transparently with a call to Parse.
// It is included in the API for advanced usage, like extending or modifying the grammar.
func NewBidiGrammar() *lr.LRAnalysis {
	b := lr.NewGrammarBuilder("UAX#9")
	b.LHS("Run").T(bc(bidi.LRI)).N("OptNI").N("L").N("OptNI").T(bc(bidi.PDI)).End()
	//b.LHS("Run").T(bc(bidi.LRI)).N("L").T(bc(bidi.PDI)).End()
	//b.LHS("Run").N("LRI").N("L").N("PDI").End()
	//b.LHS("Run").T(bc(bidi.RLI)).N("R").T(bc(bidi.PDI)).End()
	b.LHS("L").N("LParenRun").End()
	b.LHS("LParenRun").T(bc(LPAREN)).N("OptNI").N("L").N("OptNI").T(bc(RPAREN)).End()
	// b.LHS("R").N("RParenRun").End()
	// b.LHS("RParenRun").T(bc(LPAREN)).N("R").T(bc(RPAREN)).End()
	b.LHS("NI").T(bc(bidi.B)).End()
	b.LHS("NI").T(bc(bidi.S)).End()
	b.LHS("NI").T(bc(bidi.WS)).End()
	b.LHS("NI").T(bc(bidi.ON)).End()
	b.LHS("NI").N("NI").N("NI").End()
	b.LHS("L").T(bc(bidi.L)).End()
	b.LHS("R").T(bc(bidi.R)).End()
	// b.LHS("AL").T(bc(bidi.AL)).End()
	// b.LHS("EN").T(bc(bidi.EN)).End()
	// b.LHS("LEN").T(bc(LEN)).End()
	// b.LHS("AN").T(bc(bidi.AN)).End()
	// b.LHS("R").N("AL").End()                            // W3
	// b.LHS("EN").N("EN").T(bc(bidi.ES)).N("EN").End()    // W4
	// b.LHS("EN").N("EN").T(bc(bidi.CS)).N("EN").End()    // W4
	// b.LHS("AN").N("AN").T(bc(bidi.CS)).N("AN").End()    // W4
	// b.LHS("LEN").N("LEN").T(bc(bidi.ES)).N("LEN").End() // W4
	// b.LHS("LEN").N("LEN").T(bc(bidi.CS)).N("LEN").End() // W4
	// b.LHS("EN").N("EN").T(bc(bidi.ET)).End()            // W5
	// b.LHS("EN").T(bc(bidi.ET)).N("EN").End()            // W5
	// b.LHS("LEN").N("LEN").T(bc(bidi.ET)).End()          // W5
	// b.LHS("LEN").T(bc(bidi.ET)).N("LEN").End()          // W5
	b.LHS("NI").T(bc(bidi.CS)).End() // W6
	b.LHS("NI").T(bc(bidi.ET)).End() // W6
	b.LHS("NI").T(bc(bidi.ES)).End() // W6
	// b.LHS("ON").T(bc(bidi.CS)).End()                    // W6
	// b.LHS("ON").T(bc(bidi.ET)).End()                    // W6
	// b.LHS("ON").T(bc(bidi.ES)).End()                    // W6
	b.LHS("L").N("LEN").End()                // W7
	b.LHS("L").N("L").N("NI").N("L").End()   // N1
	b.LHS("R").N("R").N("NI").N("R").End()   // N1
	b.LHS("L").N("L").N("NI").N("LEN").End() //
	b.LHS("L").N("L").N("L").End()           //
	b.LHS("OptNI").N("NI").End()             //
	b.LHS("OptNI").Epsilon()                 //
	g, err := b.Grammar()
	if err != nil {
		panic(err)
	}
	return lr.Analysis(g)
}

// Parse parses a paragraph of text, given by a scanner, and parses it according to the
// Unicode UAX#9 Bidi algorithm.
func Parse(scanner *Scanner) (bool, *sppf.Forest, error) {
	if scanner == nil {
		return false, nil, errors.New("Expected parameter scanner to be non-nil")
	}
	parser := getParser()
	var parsetree *sppf.Forest
	accept, err := parser.Parse(scanner, nil)
	if accept {
		parsetree = parser.ParseForest()
	}
	return accept, parsetree, err
}

// --- Helpers ---------------------------------------------------------------

func bc(tokval bidi.Class) (string, int) {
	return ":" + ClassString(tokval), int(tokval)
}
