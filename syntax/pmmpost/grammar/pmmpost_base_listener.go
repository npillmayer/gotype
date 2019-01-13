// Code generated from PMMPost.g4 by ANTLR 4.7.2. DO NOT EDIT.

package grammar // PMMPost
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BasePMMPostListener is a complete listener for a parse tree produced by PMMPostParser.
type BasePMMPostListener struct{}

var _ PMMPostListener = &BasePMMPostListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasePMMPostListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasePMMPostListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasePMMPostListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasePMMPostListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterBeginfig is called when production beginfig is entered.
func (s *BasePMMPostListener) EnterBeginfig(ctx *BeginfigContext) {}

// ExitBeginfig is called when production beginfig is exited.
func (s *BasePMMPostListener) ExitBeginfig(ctx *BeginfigContext) {}

// EnterEndfig is called when production endfig is entered.
func (s *BasePMMPostListener) EnterEndfig(ctx *EndfigContext) {}

// ExitEndfig is called when production endfig is exited.
func (s *BasePMMPostListener) ExitEndfig(ctx *EndfigContext) {}

// EnterStatement is called when production statement is entered.
func (s *BasePMMPostListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BasePMMPostListener) ExitStatement(ctx *StatementContext) {}

// EnterTypedecl is called when production typedecl is entered.
func (s *BasePMMPostListener) EnterTypedecl(ctx *TypedeclContext) {}

// ExitTypedecl is called when production typedecl is exited.
func (s *BasePMMPostListener) ExitTypedecl(ctx *TypedeclContext) {}

// EnterLocaldecl is called when production localdecl is entered.
func (s *BasePMMPostListener) EnterLocaldecl(ctx *LocaldeclContext) {}

// ExitLocaldecl is called when production localdecl is exited.
func (s *BasePMMPostListener) ExitLocaldecl(ctx *LocaldeclContext) {}

// EnterSavecmd is called when production savecmd is entered.
func (s *BasePMMPostListener) EnterSavecmd(ctx *SavecmdContext) {}

// ExitSavecmd is called when production savecmd is exited.
func (s *BasePMMPostListener) ExitSavecmd(ctx *SavecmdContext) {}

// EnterShowcmd is called when production showcmd is entered.
func (s *BasePMMPostListener) EnterShowcmd(ctx *ShowcmdContext) {}

// ExitShowcmd is called when production showcmd is exited.
func (s *BasePMMPostListener) ExitShowcmd(ctx *ShowcmdContext) {}

// EnterProofcmd is called when production proofcmd is entered.
func (s *BasePMMPostListener) EnterProofcmd(ctx *ProofcmdContext) {}

// ExitProofcmd is called when production proofcmd is exited.
func (s *BasePMMPostListener) ExitProofcmd(ctx *ProofcmdContext) {}

// EnterLetcmd is called when production letcmd is entered.
func (s *BasePMMPostListener) EnterLetcmd(ctx *LetcmdContext) {}

// ExitLetcmd is called when production letcmd is exited.
func (s *BasePMMPostListener) ExitLetcmd(ctx *LetcmdContext) {}

// EnterCmdpickup is called when production cmdpickup is entered.
func (s *BasePMMPostListener) EnterCmdpickup(ctx *CmdpickupContext) {}

// ExitCmdpickup is called when production cmdpickup is exited.
func (s *BasePMMPostListener) ExitCmdpickup(ctx *CmdpickupContext) {}

// EnterCmddraw is called when production cmddraw is entered.
func (s *BasePMMPostListener) EnterCmddraw(ctx *CmddrawContext) {}

// ExitCmddraw is called when production cmddraw is exited.
func (s *BasePMMPostListener) ExitCmddraw(ctx *CmddrawContext) {}

// EnterCmdfill is called when production cmdfill is entered.
func (s *BasePMMPostListener) EnterCmdfill(ctx *CmdfillContext) {}

// ExitCmdfill is called when production cmdfill is exited.
func (s *BasePMMPostListener) ExitCmdfill(ctx *CmdfillContext) {}

// EnterDrawCmd is called when production drawCmd is entered.
func (s *BasePMMPostListener) EnterDrawCmd(ctx *DrawCmdContext) {}

// ExitDrawCmd is called when production drawCmd is exited.
func (s *BasePMMPostListener) ExitDrawCmd(ctx *DrawCmdContext) {}

// EnterFillCmd is called when production fillCmd is entered.
func (s *BasePMMPostListener) EnterFillCmd(ctx *FillCmdContext) {}

// ExitFillCmd is called when production fillCmd is exited.
func (s *BasePMMPostListener) ExitFillCmd(ctx *FillCmdContext) {}

// EnterPickupCmd is called when production pickupCmd is entered.
func (s *BasePMMPostListener) EnterPickupCmd(ctx *PickupCmdContext) {}

// ExitPickupCmd is called when production pickupCmd is exited.
func (s *BasePMMPostListener) ExitPickupCmd(ctx *PickupCmdContext) {}

// EnterPathjoin is called when production pathjoin is entered.
func (s *BasePMMPostListener) EnterPathjoin(ctx *PathjoinContext) {}

// ExitPathjoin is called when production pathjoin is exited.
func (s *BasePMMPostListener) ExitPathjoin(ctx *PathjoinContext) {}

// EnterCurspec is called when production curspec is entered.
func (s *BasePMMPostListener) EnterCurspec(ctx *CurspecContext) {}

// ExitCurspec is called when production curspec is exited.
func (s *BasePMMPostListener) ExitCurspec(ctx *CurspecContext) {}

// EnterDirspec is called when production dirspec is entered.
func (s *BasePMMPostListener) EnterDirspec(ctx *DirspecContext) {}

// ExitDirspec is called when production dirspec is exited.
func (s *BasePMMPostListener) ExitDirspec(ctx *DirspecContext) {}

// EnterBasicpathjoin is called when production basicpathjoin is entered.
func (s *BasePMMPostListener) EnterBasicpathjoin(ctx *BasicpathjoinContext) {}

// ExitBasicpathjoin is called when production basicpathjoin is exited.
func (s *BasePMMPostListener) ExitBasicpathjoin(ctx *BasicpathjoinContext) {}

// EnterControls is called when production controls is entered.
func (s *BasePMMPostListener) EnterControls(ctx *ControlsContext) {}

// ExitControls is called when production controls is exited.
func (s *BasePMMPostListener) ExitControls(ctx *ControlsContext) {}

// EnterStatementlist is called when production statementlist is entered.
func (s *BasePMMPostListener) EnterStatementlist(ctx *StatementlistContext) {}

// ExitStatementlist is called when production statementlist is exited.
func (s *BasePMMPostListener) ExitStatementlist(ctx *StatementlistContext) {}

// EnterVardef is called when production vardef is entered.
func (s *BasePMMPostListener) EnterVardef(ctx *VardefContext) {}

// ExitVardef is called when production vardef is exited.
func (s *BasePMMPostListener) ExitVardef(ctx *VardefContext) {}

// EnterCompound is called when production compound is entered.
func (s *BasePMMPostListener) EnterCompound(ctx *CompoundContext) {}

// ExitCompound is called when production compound is exited.
func (s *BasePMMPostListener) ExitCompound(ctx *CompoundContext) {}

// EnterEmpty is called when production empty is entered.
func (s *BasePMMPostListener) EnterEmpty(ctx *EmptyContext) {}

// ExitEmpty is called when production empty is exited.
func (s *BasePMMPostListener) ExitEmpty(ctx *EmptyContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BasePMMPostListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BasePMMPostListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterConstraint is called when production constraint is entered.
func (s *BasePMMPostListener) EnterConstraint(ctx *ConstraintContext) {}

// ExitConstraint is called when production constraint is exited.
func (s *BasePMMPostListener) ExitConstraint(ctx *ConstraintContext) {}

// EnterEquation is called when production equation is entered.
func (s *BasePMMPostListener) EnterEquation(ctx *EquationContext) {}

// ExitEquation is called when production equation is exited.
func (s *BasePMMPostListener) ExitEquation(ctx *EquationContext) {}

// EnterOrientation is called when production orientation is entered.
func (s *BasePMMPostListener) EnterOrientation(ctx *OrientationContext) {}

// ExitOrientation is called when production orientation is exited.
func (s *BasePMMPostListener) ExitOrientation(ctx *OrientationContext) {}

// EnterToken is called when production token is entered.
func (s *BasePMMPostListener) EnterToken(ctx *TokenContext) {}

// ExitToken is called when production token is exited.
func (s *BasePMMPostListener) ExitToken(ctx *TokenContext) {}

// EnterExpression is called when production expression is entered.
func (s *BasePMMPostListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BasePMMPostListener) ExitExpression(ctx *ExpressionContext) {}

// EnterPathtertiary is called when production pathtertiary is entered.
func (s *BasePMMPostListener) EnterPathtertiary(ctx *PathtertiaryContext) {}

// ExitPathtertiary is called when production pathtertiary is exited.
func (s *BasePMMPostListener) ExitPathtertiary(ctx *PathtertiaryContext) {}

// EnterTerm is called when production term is entered.
func (s *BasePMMPostListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BasePMMPostListener) ExitTerm(ctx *TermContext) {}

// EnterPath is called when production path is entered.
func (s *BasePMMPostListener) EnterPath(ctx *PathContext) {}

// ExitPath is called when production path is exited.
func (s *BasePMMPostListener) ExitPath(ctx *PathContext) {}

// EnterCycle is called when production cycle is entered.
func (s *BasePMMPostListener) EnterCycle(ctx *CycleContext) {}

// ExitCycle is called when production cycle is exited.
func (s *BasePMMPostListener) ExitCycle(ctx *CycleContext) {}

// EnterTransform is called when production transform is entered.
func (s *BasePMMPostListener) EnterTransform(ctx *TransformContext) {}

// ExitTransform is called when production transform is exited.
func (s *BasePMMPostListener) ExitTransform(ctx *TransformContext) {}

// EnterFactor is called when production factor is entered.
func (s *BasePMMPostListener) EnterFactor(ctx *FactorContext) {}

// ExitFactor is called when production factor is exited.
func (s *BasePMMPostListener) ExitFactor(ctx *FactorContext) {}

// EnterFuncatom is called when production funcatom is entered.
func (s *BasePMMPostListener) EnterFuncatom(ctx *FuncatomContext) {}

// ExitFuncatom is called when production funcatom is exited.
func (s *BasePMMPostListener) ExitFuncatom(ctx *FuncatomContext) {}

// EnterScalaratom is called when production scalaratom is entered.
func (s *BasePMMPostListener) EnterScalaratom(ctx *ScalaratomContext) {}

// ExitScalaratom is called when production scalaratom is exited.
func (s *BasePMMPostListener) ExitScalaratom(ctx *ScalaratomContext) {}

// EnterInterpolation is called when production interpolation is entered.
func (s *BasePMMPostListener) EnterInterpolation(ctx *InterpolationContext) {}

// ExitInterpolation is called when production interpolation is exited.
func (s *BasePMMPostListener) ExitInterpolation(ctx *InterpolationContext) {}

// EnterSimpleatom is called when production simpleatom is entered.
func (s *BasePMMPostListener) EnterSimpleatom(ctx *SimpleatomContext) {}

// ExitSimpleatom is called when production simpleatom is exited.
func (s *BasePMMPostListener) ExitSimpleatom(ctx *SimpleatomContext) {}

// EnterPairpart is called when production pairpart is entered.
func (s *BasePMMPostListener) EnterPairpart(ctx *PairpartContext) {}

// ExitPairpart is called when production pairpart is exited.
func (s *BasePMMPostListener) ExitPairpart(ctx *PairpartContext) {}

// EnterPointof is called when production pointof is entered.
func (s *BasePMMPostListener) EnterPointof(ctx *PointofContext) {}

// ExitPointof is called when production pointof is exited.
func (s *BasePMMPostListener) ExitPointof(ctx *PointofContext) {}

// EnterReversepath is called when production reversepath is entered.
func (s *BasePMMPostListener) EnterReversepath(ctx *ReversepathContext) {}

// ExitReversepath is called when production reversepath is exited.
func (s *BasePMMPostListener) ExitReversepath(ctx *ReversepathContext) {}

// EnterSubpath is called when production subpath is entered.
func (s *BasePMMPostListener) EnterSubpath(ctx *SubpathContext) {}

// ExitSubpath is called when production subpath is exited.
func (s *BasePMMPostListener) ExitSubpath(ctx *SubpathContext) {}

// EnterScalarmulop is called when production scalarmulop is entered.
func (s *BasePMMPostListener) EnterScalarmulop(ctx *ScalarmulopContext) {}

// ExitScalarmulop is called when production scalarmulop is exited.
func (s *BasePMMPostListener) ExitScalarmulop(ctx *ScalarmulopContext) {}

// EnterNumtokenatom is called when production numtokenatom is entered.
func (s *BasePMMPostListener) EnterNumtokenatom(ctx *NumtokenatomContext) {}

// ExitNumtokenatom is called when production numtokenatom is exited.
func (s *BasePMMPostListener) ExitNumtokenatom(ctx *NumtokenatomContext) {}

// EnterDecimal is called when production decimal is entered.
func (s *BasePMMPostListener) EnterDecimal(ctx *DecimalContext) {}

// ExitDecimal is called when production decimal is exited.
func (s *BasePMMPostListener) ExitDecimal(ctx *DecimalContext) {}

// EnterVaratom is called when production varatom is entered.
func (s *BasePMMPostListener) EnterVaratom(ctx *VaratomContext) {}

// ExitVaratom is called when production varatom is exited.
func (s *BasePMMPostListener) ExitVaratom(ctx *VaratomContext) {}

// EnterLiteralpair is called when production literalpair is entered.
func (s *BasePMMPostListener) EnterLiteralpair(ctx *LiteralpairContext) {}

// ExitLiteralpair is called when production literalpair is exited.
func (s *BasePMMPostListener) ExitLiteralpair(ctx *LiteralpairContext) {}

// EnterSubexpression is called when production subexpression is entered.
func (s *BasePMMPostListener) EnterSubexpression(ctx *SubexpressionContext) {}

// ExitSubexpression is called when production subexpression is exited.
func (s *BasePMMPostListener) ExitSubexpression(ctx *SubexpressionContext) {}

// EnterExprgroup is called when production exprgroup is entered.
func (s *BasePMMPostListener) EnterExprgroup(ctx *ExprgroupContext) {}

// ExitExprgroup is called when production exprgroup is exited.
func (s *BasePMMPostListener) ExitExprgroup(ctx *ExprgroupContext) {}

// EnterVariable is called when production variable is entered.
func (s *BasePMMPostListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production variable is exited.
func (s *BasePMMPostListener) ExitVariable(ctx *VariableContext) {}

// EnterSubscript is called when production subscript is entered.
func (s *BasePMMPostListener) EnterSubscript(ctx *SubscriptContext) {}

// ExitSubscript is called when production subscript is exited.
func (s *BasePMMPostListener) ExitSubscript(ctx *SubscriptContext) {}

// EnterAnytag is called when production anytag is entered.
func (s *BasePMMPostListener) EnterAnytag(ctx *AnytagContext) {}

// ExitAnytag is called when production anytag is exited.
func (s *BasePMMPostListener) ExitAnytag(ctx *AnytagContext) {}
