// Code generated from PMMPVar.g4 by ANTLR 4.7.2. DO NOT EDIT.

package grammar // PMMPVar
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 10, 36, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 3, 2, 3, 2, 3, 2, 7, 2, 14,
	10, 2, 12, 2, 14, 2, 17, 11, 2, 3, 2, 3, 2, 3, 3, 3, 3, 5, 3, 23, 10, 3,
	3, 4, 3, 4, 3, 4, 5, 4, 28, 10, 4, 3, 5, 3, 5, 3, 5, 3, 5, 5, 5, 34, 10,
	5, 3, 5, 2, 2, 6, 2, 4, 6, 8, 2, 2, 2, 36, 2, 10, 3, 2, 2, 2, 4, 22, 3,
	2, 2, 2, 6, 27, 3, 2, 2, 2, 8, 33, 3, 2, 2, 2, 10, 15, 5, 4, 3, 2, 11,
	14, 5, 6, 4, 2, 12, 14, 5, 8, 5, 2, 13, 11, 3, 2, 2, 2, 13, 12, 3, 2, 2,
	2, 14, 17, 3, 2, 2, 2, 15, 13, 3, 2, 2, 2, 15, 16, 3, 2, 2, 2, 16, 18,
	3, 2, 2, 2, 17, 15, 3, 2, 2, 2, 18, 19, 7, 5, 2, 2, 19, 3, 3, 2, 2, 2,
	20, 23, 7, 6, 2, 2, 21, 23, 7, 7, 2, 2, 22, 20, 3, 2, 2, 2, 22, 21, 3,
	2, 2, 2, 23, 5, 3, 2, 2, 2, 24, 25, 7, 9, 2, 2, 25, 28, 7, 7, 2, 2, 26,
	28, 7, 7, 2, 2, 27, 24, 3, 2, 2, 2, 27, 26, 3, 2, 2, 2, 28, 7, 3, 2, 2,
	2, 29, 34, 7, 8, 2, 2, 30, 31, 7, 3, 2, 2, 31, 32, 7, 8, 2, 2, 32, 34,
	7, 4, 2, 2, 33, 29, 3, 2, 2, 2, 33, 30, 3, 2, 2, 2, 34, 9, 3, 2, 2, 2,
	7, 13, 15, 22, 27, 33,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'['", "']'", "'@'", "", "", "", "'.'",
}
var symbolicNames = []string{
	"", "", "", "MARKER", "PATHTAG", "TAG", "DECIMAL", "DOT", "WS",
}

var ruleNames = []string{
	"variable", "tag", "suffix", "subscript",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type PMMPVarParser struct {
	*antlr.BaseParser
}

func NewPMMPVarParser(input antlr.TokenStream) *PMMPVarParser {
	this := new(PMMPVarParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "PMMPVar.g4"

	return this
}

// PMMPVarParser tokens.
const (
	PMMPVarParserEOF     = antlr.TokenEOF
	PMMPVarParserT__0    = 1
	PMMPVarParserT__1    = 2
	PMMPVarParserMARKER  = 3
	PMMPVarParserPATHTAG = 4
	PMMPVarParserTAG     = 5
	PMMPVarParserDECIMAL = 6
	PMMPVarParserDOT     = 7
	PMMPVarParserWS      = 8
)

// PMMPVarParser rules.
const (
	PMMPVarParserRULE_variable  = 0
	PMMPVarParserRULE_tag       = 1
	PMMPVarParserRULE_suffix    = 2
	PMMPVarParserRULE_subscript = 3
)

// IVariableContext is an interface to support dynamic dispatch.
type IVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVariableContext differentiates from other interfaces.
	IsVariableContext()
}

type VariableContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariableContext() *VariableContext {
	var p = new(VariableContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPVarParserRULE_variable
	return p
}

func (*VariableContext) IsVariableContext() {}

func NewVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableContext {
	var p = new(VariableContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPVarParserRULE_variable

	return p
}

func (s *VariableContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableContext) Tag() ITagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagContext)
}

func (s *VariableContext) MARKER() antlr.TerminalNode {
	return s.GetToken(PMMPVarParserMARKER, 0)
}

func (s *VariableContext) AllSuffix() []ISuffixContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISuffixContext)(nil)).Elem())
	var tst = make([]ISuffixContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISuffixContext)
		}
	}

	return tst
}

func (s *VariableContext) Suffix(i int) ISuffixContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISuffixContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISuffixContext)
}

func (s *VariableContext) AllSubscript() []ISubscriptContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISubscriptContext)(nil)).Elem())
	var tst = make([]ISubscriptContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISubscriptContext)
		}
	}

	return tst
}

func (s *VariableContext) Subscript(i int) ISubscriptContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISubscriptContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISubscriptContext)
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.ExitVariable(s)
	}
}

func (p *PMMPVarParser) Variable() (localctx IVariableContext) {
	localctx = NewVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, PMMPVarParserRULE_variable)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(8)
		p.Tag()
	}
	p.SetState(13)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<PMMPVarParserT__0)|(1<<PMMPVarParserTAG)|(1<<PMMPVarParserDECIMAL)|(1<<PMMPVarParserDOT))) != 0 {
		p.SetState(11)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case PMMPVarParserTAG, PMMPVarParserDOT:
			{
				p.SetState(9)
				p.Suffix()
			}

		case PMMPVarParserT__0, PMMPVarParserDECIMAL:
			{
				p.SetState(10)
				p.Subscript()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(15)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(16)
		p.Match(PMMPVarParserMARKER)
	}

	return localctx
}

// ITagContext is an interface to support dynamic dispatch.
type ITagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagContext differentiates from other interfaces.
	IsTagContext()
}

type TagContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagContext() *TagContext {
	var p = new(TagContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPVarParserRULE_tag
	return p
}

func (*TagContext) IsTagContext() {}

func NewTagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagContext {
	var p = new(TagContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPVarParserRULE_tag

	return p
}

func (s *TagContext) GetParser() antlr.Parser { return s.parser }

func (s *TagContext) CopyFrom(ctx *TagContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *TagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type PathtagContext struct {
	*TagContext
}

func NewPathtagContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PathtagContext {
	var p = new(PathtagContext)

	p.TagContext = NewEmptyTagContext()
	p.parser = parser
	p.CopyFrom(ctx.(*TagContext))

	return p
}

func (s *PathtagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathtagContext) PATHTAG() antlr.TerminalNode {
	return s.GetToken(PMMPVarParserPATHTAG, 0)
}

func (s *PathtagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.EnterPathtag(s)
	}
}

func (s *PathtagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.ExitPathtag(s)
	}
}

type SimpletagContext struct {
	*TagContext
}

func NewSimpletagContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimpletagContext {
	var p = new(SimpletagContext)

	p.TagContext = NewEmptyTagContext()
	p.parser = parser
	p.CopyFrom(ctx.(*TagContext))

	return p
}

func (s *SimpletagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpletagContext) TAG() antlr.TerminalNode {
	return s.GetToken(PMMPVarParserTAG, 0)
}

func (s *SimpletagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.EnterSimpletag(s)
	}
}

func (s *SimpletagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.ExitSimpletag(s)
	}
}

func (p *PMMPVarParser) Tag() (localctx ITagContext) {
	localctx = NewTagContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, PMMPVarParserRULE_tag)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(20)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPVarParserPATHTAG:
		localctx = NewPathtagContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(18)
			p.Match(PMMPVarParserPATHTAG)
		}

	case PMMPVarParserTAG:
		localctx = NewSimpletagContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(19)
			p.Match(PMMPVarParserTAG)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ISuffixContext is an interface to support dynamic dispatch.
type ISuffixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSuffixContext differentiates from other interfaces.
	IsSuffixContext()
}

type SuffixContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySuffixContext() *SuffixContext {
	var p = new(SuffixContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPVarParserRULE_suffix
	return p
}

func (*SuffixContext) IsSuffixContext() {}

func NewSuffixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SuffixContext {
	var p = new(SuffixContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPVarParserRULE_suffix

	return p
}

func (s *SuffixContext) GetParser() antlr.Parser { return s.parser }

func (s *SuffixContext) DOT() antlr.TerminalNode {
	return s.GetToken(PMMPVarParserDOT, 0)
}

func (s *SuffixContext) TAG() antlr.TerminalNode {
	return s.GetToken(PMMPVarParserTAG, 0)
}

func (s *SuffixContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SuffixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SuffixContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.EnterSuffix(s)
	}
}

func (s *SuffixContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.ExitSuffix(s)
	}
}

func (p *PMMPVarParser) Suffix() (localctx ISuffixContext) {
	localctx = NewSuffixContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, PMMPVarParserRULE_suffix)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(25)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPVarParserDOT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(22)
			p.Match(PMMPVarParserDOT)
		}
		{
			p.SetState(23)
			p.Match(PMMPVarParserTAG)
		}

	case PMMPVarParserTAG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(24)
			p.Match(PMMPVarParserTAG)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ISubscriptContext is an interface to support dynamic dispatch.
type ISubscriptContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSubscriptContext differentiates from other interfaces.
	IsSubscriptContext()
}

type SubscriptContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubscriptContext() *SubscriptContext {
	var p = new(SubscriptContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPVarParserRULE_subscript
	return p
}

func (*SubscriptContext) IsSubscriptContext() {}

func NewSubscriptContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubscriptContext {
	var p = new(SubscriptContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPVarParserRULE_subscript

	return p
}

func (s *SubscriptContext) GetParser() antlr.Parser { return s.parser }

func (s *SubscriptContext) DECIMAL() antlr.TerminalNode {
	return s.GetToken(PMMPVarParserDECIMAL, 0)
}

func (s *SubscriptContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubscriptContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubscriptContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.EnterSubscript(s)
	}
}

func (s *SubscriptContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPVarListener); ok {
		listenerT.ExitSubscript(s)
	}
}

func (p *PMMPVarParser) Subscript() (localctx ISubscriptContext) {
	localctx = NewSubscriptContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, PMMPVarParserRULE_subscript)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(31)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPVarParserDECIMAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(27)
			p.Match(PMMPVarParserDECIMAL)
		}

	case PMMPVarParserT__0:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(28)
			p.Match(PMMPVarParserT__0)
		}
		{
			p.SetState(29)
			p.Match(PMMPVarParserDECIMAL)
		}
		{
			p.SetState(30)
			p.Match(PMMPVarParserT__1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}
