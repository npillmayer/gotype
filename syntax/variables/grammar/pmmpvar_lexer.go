// Code generated from PMMPVar.g4 by ANTLR 4.7.2. DO NOT EDIT.

package grammar

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 10, 63, 8,
	1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9,
	7, 4, 8, 9, 8, 4, 9, 9, 9, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3,
	5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 6, 5, 33, 10, 5, 13, 5, 14, 5, 34, 3,
	6, 6, 6, 38, 10, 6, 13, 6, 14, 6, 39, 3, 7, 6, 7, 43, 10, 7, 13, 7, 14,
	7, 44, 3, 7, 3, 7, 6, 7, 49, 10, 7, 13, 7, 14, 7, 50, 5, 7, 53, 10, 7,
	3, 8, 3, 8, 3, 9, 6, 9, 58, 10, 9, 13, 9, 14, 9, 59, 3, 9, 3, 9, 2, 2,
	10, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 10, 3, 2, 4, 4, 2,
	67, 92, 99, 124, 5, 2, 11, 12, 15, 15, 34, 34, 2, 68, 2, 3, 3, 2, 2, 2,
	2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2,
	2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2, 17, 3, 2, 2, 2, 3, 19, 3, 2, 2,
	2, 5, 21, 3, 2, 2, 2, 7, 23, 3, 2, 2, 2, 9, 25, 3, 2, 2, 2, 11, 37, 3,
	2, 2, 2, 13, 42, 3, 2, 2, 2, 15, 54, 3, 2, 2, 2, 17, 57, 3, 2, 2, 2, 19,
	20, 7, 93, 2, 2, 20, 4, 3, 2, 2, 2, 21, 22, 7, 95, 2, 2, 22, 6, 3, 2, 2,
	2, 23, 24, 7, 66, 2, 2, 24, 8, 3, 2, 2, 2, 25, 26, 7, 114, 2, 2, 26, 27,
	7, 99, 2, 2, 27, 28, 7, 118, 2, 2, 28, 29, 7, 106, 2, 2, 29, 30, 7, 48,
	2, 2, 30, 32, 3, 2, 2, 2, 31, 33, 9, 2, 2, 2, 32, 31, 3, 2, 2, 2, 33, 34,
	3, 2, 2, 2, 34, 32, 3, 2, 2, 2, 34, 35, 3, 2, 2, 2, 35, 10, 3, 2, 2, 2,
	36, 38, 9, 2, 2, 2, 37, 36, 3, 2, 2, 2, 38, 39, 3, 2, 2, 2, 39, 37, 3,
	2, 2, 2, 39, 40, 3, 2, 2, 2, 40, 12, 3, 2, 2, 2, 41, 43, 4, 50, 59, 2,
	42, 41, 3, 2, 2, 2, 43, 44, 3, 2, 2, 2, 44, 42, 3, 2, 2, 2, 44, 45, 3,
	2, 2, 2, 45, 52, 3, 2, 2, 2, 46, 48, 5, 15, 8, 2, 47, 49, 4, 50, 59, 2,
	48, 47, 3, 2, 2, 2, 49, 50, 3, 2, 2, 2, 50, 48, 3, 2, 2, 2, 50, 51, 3,
	2, 2, 2, 51, 53, 3, 2, 2, 2, 52, 46, 3, 2, 2, 2, 52, 53, 3, 2, 2, 2, 53,
	14, 3, 2, 2, 2, 54, 55, 7, 48, 2, 2, 55, 16, 3, 2, 2, 2, 56, 58, 9, 3,
	2, 2, 57, 56, 3, 2, 2, 2, 58, 59, 3, 2, 2, 2, 59, 57, 3, 2, 2, 2, 59, 60,
	3, 2, 2, 2, 60, 61, 3, 2, 2, 2, 61, 62, 8, 9, 2, 2, 62, 18, 3, 2, 2, 2,
	9, 2, 34, 39, 44, 50, 52, 59, 3, 8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'['", "']'", "'@'", "", "", "", "'.'",
}

var lexerSymbolicNames = []string{
	"", "", "", "MARKER", "PATHTAG", "TAG", "DECIMAL", "DOT", "WS",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "MARKER", "PATHTAG", "TAG", "DECIMAL", "DOT", "WS",
}

type PMMPVarLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewPMMPVarLexer(input antlr.CharStream) *PMMPVarLexer {

	l := new(PMMPVarLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "PMMPVar.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// PMMPVarLexer tokens.
const (
	PMMPVarLexerT__0    = 1
	PMMPVarLexerT__1    = 2
	PMMPVarLexerMARKER  = 3
	PMMPVarLexerPATHTAG = 4
	PMMPVarLexerTAG     = 5
	PMMPVarLexerDECIMAL = 6
	PMMPVarLexerDOT     = 7
	PMMPVarLexerWS      = 8
)
