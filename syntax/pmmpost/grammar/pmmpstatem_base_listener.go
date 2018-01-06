// Generated from PMMPStatem.g4 by ANTLR 4.7.

package grammar // PMMPStatem
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BasePMMPStatemListener is a complete listener for a parse tree produced by PMMPStatemParser.
type BasePMMPStatemListener struct{}

var _ PMMPStatemListener = &BasePMMPStatemListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasePMMPStatemListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasePMMPStatemListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasePMMPStatemListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasePMMPStatemListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterFigures is called when production figures is entered.
func (s *BasePMMPStatemListener) EnterFigures(ctx *FiguresContext) {}

// ExitFigures is called when production figures is exited.
func (s *BasePMMPStatemListener) ExitFigures(ctx *FiguresContext) {}

// EnterFigure is called when production figure is entered.
func (s *BasePMMPStatemListener) EnterFigure(ctx *FigureContext) {}

// ExitFigure is called when production figure is exited.
func (s *BasePMMPStatemListener) ExitFigure(ctx *FigureContext) {}

// EnterBeginfig is called when production beginfig is entered.
func (s *BasePMMPStatemListener) EnterBeginfig(ctx *BeginfigContext) {}

// ExitBeginfig is called when production beginfig is exited.
func (s *BasePMMPStatemListener) ExitBeginfig(ctx *BeginfigContext) {}

// EnterEndfig is called when production endfig is entered.
func (s *BasePMMPStatemListener) EnterEndfig(ctx *EndfigContext) {}

// ExitEndfig is called when production endfig is exited.
func (s *BasePMMPStatemListener) ExitEndfig(ctx *EndfigContext) {}

// EnterStatementlist is called when production statementlist is entered.
func (s *BasePMMPStatemListener) EnterStatementlist(ctx *StatementlistContext) {}

// ExitStatementlist is called when production statementlist is exited.
func (s *BasePMMPStatemListener) ExitStatementlist(ctx *StatementlistContext) {}

// EnterStatement is called when production statement is entered.
func (s *BasePMMPStatemListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BasePMMPStatemListener) ExitStatement(ctx *StatementContext) {}

// EnterCompound is called when production compound is entered.
func (s *BasePMMPStatemListener) EnterCompound(ctx *CompoundContext) {}

// ExitCompound is called when production compound is exited.
func (s *BasePMMPStatemListener) ExitCompound(ctx *CompoundContext) {}

// EnterMultiequation is called when production multiequation is entered.
func (s *BasePMMPStatemListener) EnterMultiequation(ctx *MultiequationContext) {}

// ExitMultiequation is called when production multiequation is exited.
func (s *BasePMMPStatemListener) ExitMultiequation(ctx *MultiequationContext) {}

// EnterPathequation is called when production pathequation is entered.
func (s *BasePMMPStatemListener) EnterPathequation(ctx *PathequationContext) {}

// ExitPathequation is called when production pathequation is exited.
func (s *BasePMMPStatemListener) ExitPathequation(ctx *PathequationContext) {}

// EnterDeclaration is called when production declaration is entered.
func (s *BasePMMPStatemListener) EnterDeclaration(ctx *DeclarationContext) {}

// ExitDeclaration is called when production declaration is exited.
func (s *BasePMMPStatemListener) ExitDeclaration(ctx *DeclarationContext) {}

// EnterMptype is called when production mptype is entered.
func (s *BasePMMPStatemListener) EnterMptype(ctx *MptypeContext) {}

// ExitMptype is called when production mptype is exited.
func (s *BasePMMPStatemListener) ExitMptype(ctx *MptypeContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BasePMMPStatemListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BasePMMPStatemListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterLvalue is called when production lvalue is entered.
func (s *BasePMMPStatemListener) EnterLvalue(ctx *LvalueContext) {}

// ExitLvalue is called when production lvalue is exited.
func (s *BasePMMPStatemListener) ExitLvalue(ctx *LvalueContext) {}

// EnterCommand is called when production command is entered.
func (s *BasePMMPStatemListener) EnterCommand(ctx *CommandContext) {}

// ExitCommand is called when production command is exited.
func (s *BasePMMPStatemListener) ExitCommand(ctx *CommandContext) {}

// EnterSaveStmt is called when production saveStmt is entered.
func (s *BasePMMPStatemListener) EnterSaveStmt(ctx *SaveStmtContext) {}

// ExitSaveStmt is called when production saveStmt is exited.
func (s *BasePMMPStatemListener) ExitSaveStmt(ctx *SaveStmtContext) {}

// EnterShowvariableCmd is called when production showvariableCmd is entered.
func (s *BasePMMPStatemListener) EnterShowvariableCmd(ctx *ShowvariableCmdContext) {}

// ExitShowvariableCmd is called when production showvariableCmd is exited.
func (s *BasePMMPStatemListener) ExitShowvariableCmd(ctx *ShowvariableCmdContext) {}

// EnterDrawCmd is called when production drawCmd is entered.
func (s *BasePMMPStatemListener) EnterDrawCmd(ctx *DrawCmdContext) {}

// ExitDrawCmd is called when production drawCmd is exited.
func (s *BasePMMPStatemListener) ExitDrawCmd(ctx *DrawCmdContext) {}

// EnterFillCmd is called when production fillCmd is entered.
func (s *BasePMMPStatemListener) EnterFillCmd(ctx *FillCmdContext) {}

// ExitFillCmd is called when production fillCmd is exited.
func (s *BasePMMPStatemListener) ExitFillCmd(ctx *FillCmdContext) {}

// EnterPickupCmd is called when production pickupCmd is entered.
func (s *BasePMMPStatemListener) EnterPickupCmd(ctx *PickupCmdContext) {}

// ExitPickupCmd is called when production pickupCmd is exited.
func (s *BasePMMPStatemListener) ExitPickupCmd(ctx *PickupCmdContext) {}

// EnterExpression is called when production expression is entered.
func (s *BasePMMPStatemListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BasePMMPStatemListener) ExitExpression(ctx *ExpressionContext) {}

// EnterNumtertiary is called when production numtertiary is entered.
func (s *BasePMMPStatemListener) EnterNumtertiary(ctx *NumtertiaryContext) {}

// ExitNumtertiary is called when production numtertiary is exited.
func (s *BasePMMPStatemListener) ExitNumtertiary(ctx *NumtertiaryContext) {}

// EnterNumsecondary is called when production numsecondary is entered.
func (s *BasePMMPStatemListener) EnterNumsecondary(ctx *NumsecondaryContext) {}

// ExitNumsecondary is called when production numsecondary is exited.
func (s *BasePMMPStatemListener) ExitNumsecondary(ctx *NumsecondaryContext) {}

// EnterFuncnumatom is called when production funcnumatom is entered.
func (s *BasePMMPStatemListener) EnterFuncnumatom(ctx *FuncnumatomContext) {}

// ExitFuncnumatom is called when production funcnumatom is exited.
func (s *BasePMMPStatemListener) ExitFuncnumatom(ctx *FuncnumatomContext) {}

// EnterScalarnumatom is called when production scalarnumatom is entered.
func (s *BasePMMPStatemListener) EnterScalarnumatom(ctx *ScalarnumatomContext) {}

// ExitScalarnumatom is called when production scalarnumatom is exited.
func (s *BasePMMPStatemListener) ExitScalarnumatom(ctx *ScalarnumatomContext) {}

// EnterInterpolation is called when production interpolation is entered.
func (s *BasePMMPStatemListener) EnterInterpolation(ctx *InterpolationContext) {}

// ExitInterpolation is called when production interpolation is exited.
func (s *BasePMMPStatemListener) ExitInterpolation(ctx *InterpolationContext) {}

// EnterSimplenumatom is called when production simplenumatom is entered.
func (s *BasePMMPStatemListener) EnterSimplenumatom(ctx *SimplenumatomContext) {}

// ExitSimplenumatom is called when production simplenumatom is exited.
func (s *BasePMMPStatemListener) ExitSimplenumatom(ctx *SimplenumatomContext) {}

// EnterPointdistance is called when production pointdistance is entered.
func (s *BasePMMPStatemListener) EnterPointdistance(ctx *PointdistanceContext) {}

// ExitPointdistance is called when production pointdistance is exited.
func (s *BasePMMPStatemListener) ExitPointdistance(ctx *PointdistanceContext) {}

// EnterPairpart is called when production pairpart is entered.
func (s *BasePMMPStatemListener) EnterPairpart(ctx *PairpartContext) {}

// ExitPairpart is called when production pairpart is exited.
func (s *BasePMMPStatemListener) ExitPairpart(ctx *PairpartContext) {}

// EnterScalarmulop is called when production scalarmulop is entered.
func (s *BasePMMPStatemListener) EnterScalarmulop(ctx *ScalarmulopContext) {}

// ExitScalarmulop is called when production scalarmulop is exited.
func (s *BasePMMPStatemListener) ExitScalarmulop(ctx *ScalarmulopContext) {}

// EnterNumtokenatom is called when production numtokenatom is entered.
func (s *BasePMMPStatemListener) EnterNumtokenatom(ctx *NumtokenatomContext) {}

// ExitNumtokenatom is called when production numtokenatom is exited.
func (s *BasePMMPStatemListener) ExitNumtokenatom(ctx *NumtokenatomContext) {}

// EnterInternal is called when production internal is entered.
func (s *BasePMMPStatemListener) EnterInternal(ctx *InternalContext) {}

// ExitInternal is called when production internal is exited.
func (s *BasePMMPStatemListener) ExitInternal(ctx *InternalContext) {}

// EnterWhatever is called when production whatever is entered.
func (s *BasePMMPStatemListener) EnterWhatever(ctx *WhateverContext) {}

// ExitWhatever is called when production whatever is exited.
func (s *BasePMMPStatemListener) ExitWhatever(ctx *WhateverContext) {}

// EnterVariable is called when production variable is entered.
func (s *BasePMMPStatemListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production variable is exited.
func (s *BasePMMPStatemListener) ExitVariable(ctx *VariableContext) {}

// EnterDecimal is called when production decimal is entered.
func (s *BasePMMPStatemListener) EnterDecimal(ctx *DecimalContext) {}

// ExitDecimal is called when production decimal is exited.
func (s *BasePMMPStatemListener) ExitDecimal(ctx *DecimalContext) {}

// EnterSubexpression is called when production subexpression is entered.
func (s *BasePMMPStatemListener) EnterSubexpression(ctx *SubexpressionContext) {}

// ExitSubexpression is called when production subexpression is exited.
func (s *BasePMMPStatemListener) ExitSubexpression(ctx *SubexpressionContext) {}

// EnterExprgroup is called when production exprgroup is entered.
func (s *BasePMMPStatemListener) EnterExprgroup(ctx *ExprgroupContext) {}

// ExitExprgroup is called when production exprgroup is exited.
func (s *BasePMMPStatemListener) ExitExprgroup(ctx *ExprgroupContext) {}

// EnterSubscript is called when production subscript is entered.
func (s *BasePMMPStatemListener) EnterSubscript(ctx *SubscriptContext) {}

// ExitSubscript is called when production subscript is exited.
func (s *BasePMMPStatemListener) ExitSubscript(ctx *SubscriptContext) {}

// EnterPairtertiary is called when production pairtertiary is entered.
func (s *BasePMMPStatemListener) EnterPairtertiary(ctx *PairtertiaryContext) {}

// ExitPairtertiary is called when production pairtertiary is exited.
func (s *BasePMMPStatemListener) ExitPairtertiary(ctx *PairtertiaryContext) {}

// EnterPairsecond is called when production pairsecond is entered.
func (s *BasePMMPStatemListener) EnterPairsecond(ctx *PairsecondContext) {}

// ExitPairsecond is called when production pairsecond is exited.
func (s *BasePMMPStatemListener) ExitPairsecond(ctx *PairsecondContext) {}

// EnterTransform is called when production transform is entered.
func (s *BasePMMPStatemListener) EnterTransform(ctx *TransformContext) {}

// ExitTransform is called when production transform is exited.
func (s *BasePMMPStatemListener) ExitTransform(ctx *TransformContext) {}

// EnterTransformer is called when production transformer is entered.
func (s *BasePMMPStatemListener) EnterTransformer(ctx *TransformerContext) {}

// ExitTransformer is called when production transformer is exited.
func (s *BasePMMPStatemListener) ExitTransformer(ctx *TransformerContext) {}

// EnterSimplepairatom is called when production simplepairatom is entered.
func (s *BasePMMPStatemListener) EnterSimplepairatom(ctx *SimplepairatomContext) {}

// ExitSimplepairatom is called when production simplepairatom is exited.
func (s *BasePMMPStatemListener) ExitSimplepairatom(ctx *SimplepairatomContext) {}

// EnterScalarmuloppair is called when production scalarmuloppair is entered.
func (s *BasePMMPStatemListener) EnterScalarmuloppair(ctx *ScalarmuloppairContext) {}

// ExitScalarmuloppair is called when production scalarmuloppair is exited.
func (s *BasePMMPStatemListener) ExitScalarmuloppair(ctx *ScalarmuloppairContext) {}

// EnterPathpoint is called when production pathpoint is entered.
func (s *BasePMMPStatemListener) EnterPathpoint(ctx *PathpointContext) {}

// ExitPathpoint is called when production pathpoint is exited.
func (s *BasePMMPStatemListener) ExitPathpoint(ctx *PathpointContext) {}

// EnterPairinterpolation is called when production pairinterpolation is entered.
func (s *BasePMMPStatemListener) EnterPairinterpolation(ctx *PairinterpolationContext) {}

// ExitPairinterpolation is called when production pairinterpolation is exited.
func (s *BasePMMPStatemListener) ExitPairinterpolation(ctx *PairinterpolationContext) {}

// EnterLiteralpair is called when production literalpair is entered.
func (s *BasePMMPStatemListener) EnterLiteralpair(ctx *LiteralpairContext) {}

// ExitLiteralpair is called when production literalpair is exited.
func (s *BasePMMPStatemListener) ExitLiteralpair(ctx *LiteralpairContext) {}

// EnterPairvariable is called when production pairvariable is entered.
func (s *BasePMMPStatemListener) EnterPairvariable(ctx *PairvariableContext) {}

// ExitPairvariable is called when production pairvariable is exited.
func (s *BasePMMPStatemListener) ExitPairvariable(ctx *PairvariableContext) {}

// EnterSubpairexpression is called when production subpairexpression is entered.
func (s *BasePMMPStatemListener) EnterSubpairexpression(ctx *SubpairexpressionContext) {}

// ExitSubpairexpression is called when production subpairexpression is exited.
func (s *BasePMMPStatemListener) ExitSubpairexpression(ctx *SubpairexpressionContext) {}

// EnterPairexprgroup is called when production pairexprgroup is entered.
func (s *BasePMMPStatemListener) EnterPairexprgroup(ctx *PairexprgroupContext) {}

// ExitPairexprgroup is called when production pairexprgroup is exited.
func (s *BasePMMPStatemListener) ExitPairexprgroup(ctx *PairexprgroupContext) {}

// EnterPathexpression is called when production pathexpression is entered.
func (s *BasePMMPStatemListener) EnterPathexpression(ctx *PathexpressionContext) {}

// ExitPathexpression is called when production pathexpression is exited.
func (s *BasePMMPStatemListener) ExitPathexpression(ctx *PathexpressionContext) {}

// EnterPathtertiary is called when production pathtertiary is entered.
func (s *BasePMMPStatemListener) EnterPathtertiary(ctx *PathtertiaryContext) {}

// ExitPathtertiary is called when production pathtertiary is exited.
func (s *BasePMMPStatemListener) ExitPathtertiary(ctx *PathtertiaryContext) {}

// EnterPathfragm is called when production pathfragm is entered.
func (s *BasePMMPStatemListener) EnterPathfragm(ctx *PathfragmContext) {}

// ExitPathfragm is called when production pathfragm is exited.
func (s *BasePMMPStatemListener) ExitPathfragm(ctx *PathfragmContext) {}

// EnterCycle is called when production cycle is entered.
func (s *BasePMMPStatemListener) EnterCycle(ctx *CycleContext) {}

// ExitCycle is called when production cycle is exited.
func (s *BasePMMPStatemListener) ExitCycle(ctx *CycleContext) {}

// EnterPathsecondary is called when production pathsecondary is entered.
func (s *BasePMMPStatemListener) EnterPathsecondary(ctx *PathsecondaryContext) {}

// ExitPathsecondary is called when production pathsecondary is exited.
func (s *BasePMMPStatemListener) ExitPathsecondary(ctx *PathsecondaryContext) {}

// EnterAtomicpath is called when production atomicpath is entered.
func (s *BasePMMPStatemListener) EnterAtomicpath(ctx *AtomicpathContext) {}

// ExitAtomicpath is called when production atomicpath is exited.
func (s *BasePMMPStatemListener) ExitAtomicpath(ctx *AtomicpathContext) {}

// EnterReversepath is called when production reversepath is entered.
func (s *BasePMMPStatemListener) EnterReversepath(ctx *ReversepathContext) {}

// ExitReversepath is called when production reversepath is exited.
func (s *BasePMMPStatemListener) ExitReversepath(ctx *ReversepathContext) {}

// EnterSubpath is called when production subpath is entered.
func (s *BasePMMPStatemListener) EnterSubpath(ctx *SubpathContext) {}

// ExitSubpath is called when production subpath is exited.
func (s *BasePMMPStatemListener) ExitSubpath(ctx *SubpathContext) {}

// EnterPathvariable is called when production pathvariable is entered.
func (s *BasePMMPStatemListener) EnterPathvariable(ctx *PathvariableContext) {}

// ExitPathvariable is called when production pathvariable is exited.
func (s *BasePMMPStatemListener) ExitPathvariable(ctx *PathvariableContext) {}

// EnterTag is called when production tag is entered.
func (s *BasePMMPStatemListener) EnterTag(ctx *TagContext) {}

// ExitTag is called when production tag is exited.
func (s *BasePMMPStatemListener) ExitTag(ctx *TagContext) {}

// EnterAnytag is called when production anytag is entered.
func (s *BasePMMPStatemListener) EnterAnytag(ctx *AnytagContext) {}

// ExitAnytag is called when production anytag is exited.
func (s *BasePMMPStatemListener) ExitAnytag(ctx *AnytagContext) {}
