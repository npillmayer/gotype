// Generated from Gallery.g4 by ANTLR 4.7.

package grammar // Gallery
import "github.com/antlr/antlr4/runtime/Go/antlr"

// GalleryListener is a complete listener for a parse tree produced by GalleryParser.
type GalleryListener interface {
	antlr.ParseTreeListener

	// EnterStatementlist is called when entering the statementlist production.
	EnterStatementlist(c *StatementlistContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterDeclaration is called when entering the declaration production.
	EnterDeclaration(c *DeclarationContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterCompound is called when entering the compound production.
	EnterCompound(c *CompoundContext)

	// EnterConstraint is called when entering the constraint production.
	EnterConstraint(c *ConstraintContext)

	// EnterEquation is called when entering the equation production.
	EnterEquation(c *EquationContext)

	// EnterOrientation is called when entering the orientation production.
	EnterOrientation(c *OrientationContext)

	// EnterSavecmd is called when entering the savecmd production.
	EnterSavecmd(c *SavecmdContext)

	// EnterShowcmd is called when entering the showcmd production.
	EnterShowcmd(c *ShowcmdContext)

	// EnterProofcmd is called when entering the proofcmd production.
	EnterProofcmd(c *ProofcmdContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterPathtertiary is called when entering the pathtertiary production.
	EnterPathtertiary(c *PathtertiaryContext)

	// EnterLonesecondary is called when entering the lonesecondary production.
	EnterLonesecondary(c *LonesecondaryContext)

	// EnterTerm is called when entering the term production.
	EnterTerm(c *TermContext)

	// EnterPath is called when entering the path production.
	EnterPath(c *PathContext)

	// EnterCycle is called when entering the cycle production.
	EnterCycle(c *CycleContext)

	// EnterTransform is called when entering the transform production.
	EnterTransform(c *TransformContext)

	// EnterLoneprimary is called when entering the loneprimary production.
	EnterLoneprimary(c *LoneprimaryContext)

	// EnterFactor is called when entering the factor production.
	EnterFactor(c *FactorContext)

	// EnterTransformer is called when entering the transformer production.
	EnterTransformer(c *TransformerContext)

	// EnterFuncnumatom is called when entering the funcnumatom production.
	EnterFuncnumatom(c *FuncnumatomContext)

	// EnterScalarnumatom is called when entering the scalarnumatom production.
	EnterScalarnumatom(c *ScalarnumatomContext)

	// EnterInterpolation is called when entering the interpolation production.
	EnterInterpolation(c *InterpolationContext)

	// EnterSimplenumatom is called when entering the simplenumatom production.
	EnterSimplenumatom(c *SimplenumatomContext)

	// EnterPairpart is called when entering the pairpart production.
	EnterPairpart(c *PairpartContext)

	// EnterPathpoint is called when entering the pathpoint production.
	EnterPathpoint(c *PathpointContext)

	// EnterReversepath is called when entering the reversepath production.
	EnterReversepath(c *ReversepathContext)

	// EnterSubpath is called when entering the subpath production.
	EnterSubpath(c *SubpathContext)

	// EnterEdgeconstraint is called when entering the edgeconstraint production.
	EnterEdgeconstraint(c *EdgeconstraintContext)

	// EnterBox is called when entering the box production.
	EnterBox(c *BoxContext)

	// EnterScalarmulop is called when entering the scalarmulop production.
	EnterScalarmulop(c *ScalarmulopContext)

	// EnterNumtokenatom is called when entering the numtokenatom production.
	EnterNumtokenatom(c *NumtokenatomContext)

	// EnterDecimal is called when entering the decimal production.
	EnterDecimal(c *DecimalContext)

	// EnterVaratom is called when entering the varatom production.
	EnterVaratom(c *VaratomContext)

	// EnterLiteralpair is called when entering the literalpair production.
	EnterLiteralpair(c *LiteralpairContext)

	// EnterSubexpression is called when entering the subexpression production.
	EnterSubexpression(c *SubexpressionContext)

	// EnterExprgroup is called when entering the exprgroup production.
	EnterExprgroup(c *ExprgroupContext)

	// EnterVariable is called when entering the variable production.
	EnterVariable(c *VariableContext)

	// EnterSubscript is called when entering the subscript production.
	EnterSubscript(c *SubscriptContext)

	// EnterAnytag is called when entering the anytag production.
	EnterAnytag(c *AnytagContext)

	// ExitStatementlist is called when exiting the statementlist production.
	ExitStatementlist(c *StatementlistContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitDeclaration is called when exiting the declaration production.
	ExitDeclaration(c *DeclarationContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitCompound is called when exiting the compound production.
	ExitCompound(c *CompoundContext)

	// ExitConstraint is called when exiting the constraint production.
	ExitConstraint(c *ConstraintContext)

	// ExitEquation is called when exiting the equation production.
	ExitEquation(c *EquationContext)

	// ExitOrientation is called when exiting the orientation production.
	ExitOrientation(c *OrientationContext)

	// ExitSavecmd is called when exiting the savecmd production.
	ExitSavecmd(c *SavecmdContext)

	// ExitShowcmd is called when exiting the showcmd production.
	ExitShowcmd(c *ShowcmdContext)

	// ExitProofcmd is called when exiting the proofcmd production.
	ExitProofcmd(c *ProofcmdContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitPathtertiary is called when exiting the pathtertiary production.
	ExitPathtertiary(c *PathtertiaryContext)

	// ExitLonesecondary is called when exiting the lonesecondary production.
	ExitLonesecondary(c *LonesecondaryContext)

	// ExitTerm is called when exiting the term production.
	ExitTerm(c *TermContext)

	// ExitPath is called when exiting the path production.
	ExitPath(c *PathContext)

	// ExitCycle is called when exiting the cycle production.
	ExitCycle(c *CycleContext)

	// ExitTransform is called when exiting the transform production.
	ExitTransform(c *TransformContext)

	// ExitLoneprimary is called when exiting the loneprimary production.
	ExitLoneprimary(c *LoneprimaryContext)

	// ExitFactor is called when exiting the factor production.
	ExitFactor(c *FactorContext)

	// ExitTransformer is called when exiting the transformer production.
	ExitTransformer(c *TransformerContext)

	// ExitFuncnumatom is called when exiting the funcnumatom production.
	ExitFuncnumatom(c *FuncnumatomContext)

	// ExitScalarnumatom is called when exiting the scalarnumatom production.
	ExitScalarnumatom(c *ScalarnumatomContext)

	// ExitInterpolation is called when exiting the interpolation production.
	ExitInterpolation(c *InterpolationContext)

	// ExitSimplenumatom is called when exiting the simplenumatom production.
	ExitSimplenumatom(c *SimplenumatomContext)

	// ExitPairpart is called when exiting the pairpart production.
	ExitPairpart(c *PairpartContext)

	// ExitPathpoint is called when exiting the pathpoint production.
	ExitPathpoint(c *PathpointContext)

	// ExitReversepath is called when exiting the reversepath production.
	ExitReversepath(c *ReversepathContext)

	// ExitSubpath is called when exiting the subpath production.
	ExitSubpath(c *SubpathContext)

	// ExitEdgeconstraint is called when exiting the edgeconstraint production.
	ExitEdgeconstraint(c *EdgeconstraintContext)

	// ExitBox is called when exiting the box production.
	ExitBox(c *BoxContext)

	// ExitScalarmulop is called when exiting the scalarmulop production.
	ExitScalarmulop(c *ScalarmulopContext)

	// ExitNumtokenatom is called when exiting the numtokenatom production.
	ExitNumtokenatom(c *NumtokenatomContext)

	// ExitDecimal is called when exiting the decimal production.
	ExitDecimal(c *DecimalContext)

	// ExitVaratom is called when exiting the varatom production.
	ExitVaratom(c *VaratomContext)

	// ExitLiteralpair is called when exiting the literalpair production.
	ExitLiteralpair(c *LiteralpairContext)

	// ExitSubexpression is called when exiting the subexpression production.
	ExitSubexpression(c *SubexpressionContext)

	// ExitExprgroup is called when exiting the exprgroup production.
	ExitExprgroup(c *ExprgroupContext)

	// ExitVariable is called when exiting the variable production.
	ExitVariable(c *VariableContext)

	// ExitSubscript is called when exiting the subscript production.
	ExitSubscript(c *SubscriptContext)

	// ExitAnytag is called when exiting the anytag production.
	ExitAnytag(c *AnytagContext)
}
