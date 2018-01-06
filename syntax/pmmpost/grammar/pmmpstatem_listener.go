// Generated from PMMPStatem.g4 by ANTLR 4.7.

package grammar // PMMPStatem
import "github.com/antlr/antlr4/runtime/Go/antlr"

// PMMPStatemListener is a complete listener for a parse tree produced by PMMPStatemParser.
type PMMPStatemListener interface {
	antlr.ParseTreeListener

	// EnterFigures is called when entering the figures production.
	EnterFigures(c *FiguresContext)

	// EnterFigure is called when entering the figure production.
	EnterFigure(c *FigureContext)

	// EnterBeginfig is called when entering the beginfig production.
	EnterBeginfig(c *BeginfigContext)

	// EnterEndfig is called when entering the endfig production.
	EnterEndfig(c *EndfigContext)

	// EnterStatementlist is called when entering the statementlist production.
	EnterStatementlist(c *StatementlistContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterCompound is called when entering the compound production.
	EnterCompound(c *CompoundContext)

	// EnterMultiequation is called when entering the multiequation production.
	EnterMultiequation(c *MultiequationContext)

	// EnterPathequation is called when entering the pathequation production.
	EnterPathequation(c *PathequationContext)

	// EnterDeclaration is called when entering the declaration production.
	EnterDeclaration(c *DeclarationContext)

	// EnterMptype is called when entering the mptype production.
	EnterMptype(c *MptypeContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterLvalue is called when entering the lvalue production.
	EnterLvalue(c *LvalueContext)

	// EnterCommand is called when entering the command production.
	EnterCommand(c *CommandContext)

	// EnterSaveStmt is called when entering the saveStmt production.
	EnterSaveStmt(c *SaveStmtContext)

	// EnterShowvariableCmd is called when entering the showvariableCmd production.
	EnterShowvariableCmd(c *ShowvariableCmdContext)

	// EnterDrawCmd is called when entering the drawCmd production.
	EnterDrawCmd(c *DrawCmdContext)

	// EnterFillCmd is called when entering the fillCmd production.
	EnterFillCmd(c *FillCmdContext)

	// EnterPickupCmd is called when entering the pickupCmd production.
	EnterPickupCmd(c *PickupCmdContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterNumtertiary is called when entering the numtertiary production.
	EnterNumtertiary(c *NumtertiaryContext)

	// EnterNumsecondary is called when entering the numsecondary production.
	EnterNumsecondary(c *NumsecondaryContext)

	// EnterFuncnumatom is called when entering the funcnumatom production.
	EnterFuncnumatom(c *FuncnumatomContext)

	// EnterScalarnumatom is called when entering the scalarnumatom production.
	EnterScalarnumatom(c *ScalarnumatomContext)

	// EnterInterpolation is called when entering the interpolation production.
	EnterInterpolation(c *InterpolationContext)

	// EnterSimplenumatom is called when entering the simplenumatom production.
	EnterSimplenumatom(c *SimplenumatomContext)

	// EnterPointdistance is called when entering the pointdistance production.
	EnterPointdistance(c *PointdistanceContext)

	// EnterPairpart is called when entering the pairpart production.
	EnterPairpart(c *PairpartContext)

	// EnterScalarmulop is called when entering the scalarmulop production.
	EnterScalarmulop(c *ScalarmulopContext)

	// EnterNumtokenatom is called when entering the numtokenatom production.
	EnterNumtokenatom(c *NumtokenatomContext)

	// EnterInternal is called when entering the internal production.
	EnterInternal(c *InternalContext)

	// EnterWhatever is called when entering the whatever production.
	EnterWhatever(c *WhateverContext)

	// EnterVariable is called when entering the variable production.
	EnterVariable(c *VariableContext)

	// EnterDecimal is called when entering the decimal production.
	EnterDecimal(c *DecimalContext)

	// EnterSubexpression is called when entering the subexpression production.
	EnterSubexpression(c *SubexpressionContext)

	// EnterExprgroup is called when entering the exprgroup production.
	EnterExprgroup(c *ExprgroupContext)

	// EnterSubscript is called when entering the subscript production.
	EnterSubscript(c *SubscriptContext)

	// EnterPairtertiary is called when entering the pairtertiary production.
	EnterPairtertiary(c *PairtertiaryContext)

	// EnterPairsecond is called when entering the pairsecond production.
	EnterPairsecond(c *PairsecondContext)

	// EnterTransform is called when entering the transform production.
	EnterTransform(c *TransformContext)

	// EnterTransformer is called when entering the transformer production.
	EnterTransformer(c *TransformerContext)

	// EnterSimplepairatom is called when entering the simplepairatom production.
	EnterSimplepairatom(c *SimplepairatomContext)

	// EnterScalarmuloppair is called when entering the scalarmuloppair production.
	EnterScalarmuloppair(c *ScalarmuloppairContext)

	// EnterPathpoint is called when entering the pathpoint production.
	EnterPathpoint(c *PathpointContext)

	// EnterPairinterpolation is called when entering the pairinterpolation production.
	EnterPairinterpolation(c *PairinterpolationContext)

	// EnterLiteralpair is called when entering the literalpair production.
	EnterLiteralpair(c *LiteralpairContext)

	// EnterPairvariable is called when entering the pairvariable production.
	EnterPairvariable(c *PairvariableContext)

	// EnterSubpairexpression is called when entering the subpairexpression production.
	EnterSubpairexpression(c *SubpairexpressionContext)

	// EnterPairexprgroup is called when entering the pairexprgroup production.
	EnterPairexprgroup(c *PairexprgroupContext)

	// EnterPathexpression is called when entering the pathexpression production.
	EnterPathexpression(c *PathexpressionContext)

	// EnterPathtertiary is called when entering the pathtertiary production.
	EnterPathtertiary(c *PathtertiaryContext)

	// EnterPathfragm is called when entering the pathfragm production.
	EnterPathfragm(c *PathfragmContext)

	// EnterCycle is called when entering the cycle production.
	EnterCycle(c *CycleContext)

	// EnterPathsecondary is called when entering the pathsecondary production.
	EnterPathsecondary(c *PathsecondaryContext)

	// EnterAtomicpath is called when entering the atomicpath production.
	EnterAtomicpath(c *AtomicpathContext)

	// EnterReversepath is called when entering the reversepath production.
	EnterReversepath(c *ReversepathContext)

	// EnterSubpath is called when entering the subpath production.
	EnterSubpath(c *SubpathContext)

	// EnterPathvariable is called when entering the pathvariable production.
	EnterPathvariable(c *PathvariableContext)

	// EnterTag is called when entering the tag production.
	EnterTag(c *TagContext)

	// EnterAnytag is called when entering the anytag production.
	EnterAnytag(c *AnytagContext)

	// ExitFigures is called when exiting the figures production.
	ExitFigures(c *FiguresContext)

	// ExitFigure is called when exiting the figure production.
	ExitFigure(c *FigureContext)

	// ExitBeginfig is called when exiting the beginfig production.
	ExitBeginfig(c *BeginfigContext)

	// ExitEndfig is called when exiting the endfig production.
	ExitEndfig(c *EndfigContext)

	// ExitStatementlist is called when exiting the statementlist production.
	ExitStatementlist(c *StatementlistContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitCompound is called when exiting the compound production.
	ExitCompound(c *CompoundContext)

	// ExitMultiequation is called when exiting the multiequation production.
	ExitMultiequation(c *MultiequationContext)

	// ExitPathequation is called when exiting the pathequation production.
	ExitPathequation(c *PathequationContext)

	// ExitDeclaration is called when exiting the declaration production.
	ExitDeclaration(c *DeclarationContext)

	// ExitMptype is called when exiting the mptype production.
	ExitMptype(c *MptypeContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitLvalue is called when exiting the lvalue production.
	ExitLvalue(c *LvalueContext)

	// ExitCommand is called when exiting the command production.
	ExitCommand(c *CommandContext)

	// ExitSaveStmt is called when exiting the saveStmt production.
	ExitSaveStmt(c *SaveStmtContext)

	// ExitShowvariableCmd is called when exiting the showvariableCmd production.
	ExitShowvariableCmd(c *ShowvariableCmdContext)

	// ExitDrawCmd is called when exiting the drawCmd production.
	ExitDrawCmd(c *DrawCmdContext)

	// ExitFillCmd is called when exiting the fillCmd production.
	ExitFillCmd(c *FillCmdContext)

	// ExitPickupCmd is called when exiting the pickupCmd production.
	ExitPickupCmd(c *PickupCmdContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitNumtertiary is called when exiting the numtertiary production.
	ExitNumtertiary(c *NumtertiaryContext)

	// ExitNumsecondary is called when exiting the numsecondary production.
	ExitNumsecondary(c *NumsecondaryContext)

	// ExitFuncnumatom is called when exiting the funcnumatom production.
	ExitFuncnumatom(c *FuncnumatomContext)

	// ExitScalarnumatom is called when exiting the scalarnumatom production.
	ExitScalarnumatom(c *ScalarnumatomContext)

	// ExitInterpolation is called when exiting the interpolation production.
	ExitInterpolation(c *InterpolationContext)

	// ExitSimplenumatom is called when exiting the simplenumatom production.
	ExitSimplenumatom(c *SimplenumatomContext)

	// ExitPointdistance is called when exiting the pointdistance production.
	ExitPointdistance(c *PointdistanceContext)

	// ExitPairpart is called when exiting the pairpart production.
	ExitPairpart(c *PairpartContext)

	// ExitScalarmulop is called when exiting the scalarmulop production.
	ExitScalarmulop(c *ScalarmulopContext)

	// ExitNumtokenatom is called when exiting the numtokenatom production.
	ExitNumtokenatom(c *NumtokenatomContext)

	// ExitInternal is called when exiting the internal production.
	ExitInternal(c *InternalContext)

	// ExitWhatever is called when exiting the whatever production.
	ExitWhatever(c *WhateverContext)

	// ExitVariable is called when exiting the variable production.
	ExitVariable(c *VariableContext)

	// ExitDecimal is called when exiting the decimal production.
	ExitDecimal(c *DecimalContext)

	// ExitSubexpression is called when exiting the subexpression production.
	ExitSubexpression(c *SubexpressionContext)

	// ExitExprgroup is called when exiting the exprgroup production.
	ExitExprgroup(c *ExprgroupContext)

	// ExitSubscript is called when exiting the subscript production.
	ExitSubscript(c *SubscriptContext)

	// ExitPairtertiary is called when exiting the pairtertiary production.
	ExitPairtertiary(c *PairtertiaryContext)

	// ExitPairsecond is called when exiting the pairsecond production.
	ExitPairsecond(c *PairsecondContext)

	// ExitTransform is called when exiting the transform production.
	ExitTransform(c *TransformContext)

	// ExitTransformer is called when exiting the transformer production.
	ExitTransformer(c *TransformerContext)

	// ExitSimplepairatom is called when exiting the simplepairatom production.
	ExitSimplepairatom(c *SimplepairatomContext)

	// ExitScalarmuloppair is called when exiting the scalarmuloppair production.
	ExitScalarmuloppair(c *ScalarmuloppairContext)

	// ExitPathpoint is called when exiting the pathpoint production.
	ExitPathpoint(c *PathpointContext)

	// ExitPairinterpolation is called when exiting the pairinterpolation production.
	ExitPairinterpolation(c *PairinterpolationContext)

	// ExitLiteralpair is called when exiting the literalpair production.
	ExitLiteralpair(c *LiteralpairContext)

	// ExitPairvariable is called when exiting the pairvariable production.
	ExitPairvariable(c *PairvariableContext)

	// ExitSubpairexpression is called when exiting the subpairexpression production.
	ExitSubpairexpression(c *SubpairexpressionContext)

	// ExitPairexprgroup is called when exiting the pairexprgroup production.
	ExitPairexprgroup(c *PairexprgroupContext)

	// ExitPathexpression is called when exiting the pathexpression production.
	ExitPathexpression(c *PathexpressionContext)

	// ExitPathtertiary is called when exiting the pathtertiary production.
	ExitPathtertiary(c *PathtertiaryContext)

	// ExitPathfragm is called when exiting the pathfragm production.
	ExitPathfragm(c *PathfragmContext)

	// ExitCycle is called when exiting the cycle production.
	ExitCycle(c *CycleContext)

	// ExitPathsecondary is called when exiting the pathsecondary production.
	ExitPathsecondary(c *PathsecondaryContext)

	// ExitAtomicpath is called when exiting the atomicpath production.
	ExitAtomicpath(c *AtomicpathContext)

	// ExitReversepath is called when exiting the reversepath production.
	ExitReversepath(c *ReversepathContext)

	// ExitSubpath is called when exiting the subpath production.
	ExitSubpath(c *SubpathContext)

	// ExitPathvariable is called when exiting the pathvariable production.
	ExitPathvariable(c *PathvariableContext)

	// ExitTag is called when exiting the tag production.
	ExitTag(c *TagContext)

	// ExitAnytag is called when exiting the anytag production.
	ExitAnytag(c *AnytagContext)
}
