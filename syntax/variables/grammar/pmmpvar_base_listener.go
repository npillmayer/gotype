// Code generated from PMMPVar.g4 by ANTLR 4.7.2. DO NOT EDIT.

package grammar // PMMPVar
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BasePMMPVarListener is a complete listener for a parse tree produced by PMMPVarParser.
type BasePMMPVarListener struct{}

var _ PMMPVarListener = &BasePMMPVarListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasePMMPVarListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasePMMPVarListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasePMMPVarListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasePMMPVarListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterVariable is called when production variable is entered.
func (s *BasePMMPVarListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production variable is exited.
func (s *BasePMMPVarListener) ExitVariable(ctx *VariableContext) {}

// EnterPathtag is called when production pathtag is entered.
func (s *BasePMMPVarListener) EnterPathtag(ctx *PathtagContext) {}

// ExitPathtag is called when production pathtag is exited.
func (s *BasePMMPVarListener) ExitPathtag(ctx *PathtagContext) {}

// EnterSimpletag is called when production simpletag is entered.
func (s *BasePMMPVarListener) EnterSimpletag(ctx *SimpletagContext) {}

// ExitSimpletag is called when production simpletag is exited.
func (s *BasePMMPVarListener) ExitSimpletag(ctx *SimpletagContext) {}

// EnterSuffix is called when production suffix is entered.
func (s *BasePMMPVarListener) EnterSuffix(ctx *SuffixContext) {}

// ExitSuffix is called when production suffix is exited.
func (s *BasePMMPVarListener) ExitSuffix(ctx *SuffixContext) {}

// EnterSubscript is called when production subscript is entered.
func (s *BasePMMPVarListener) EnterSubscript(ctx *SubscriptContext) {}

// ExitSubscript is called when production subscript is exited.
func (s *BasePMMPVarListener) ExitSubscript(ctx *SubscriptContext) {}
