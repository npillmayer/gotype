// Generated from Gallery.g4 by ANTLR 4.7.

package grammar // Gallery
import "github.com/antlr/antlr4/runtime/Go/antlr"

// GalleryListener is a complete listener for a parse tree produced by GalleryParser.
type GalleryListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterTypedecl is called when entering the typedecl production.
	EnterTypedecl(c *TypedeclContext)

	// EnterLocaldecl is called when entering the localdecl production.
	EnterLocaldecl(c *LocaldeclContext)

	// EnterParameterdecl is called when entering the parameterdecl production.
	EnterParameterdecl(c *ParameterdeclContext)

	// EnterSavecmd is called when entering the savecmd production.
	EnterSavecmd(c *SavecmdContext)

	// EnterShowcmd is called when entering the showcmd production.
	EnterShowcmd(c *ShowcmdContext)

	// EnterProofcmd is called when entering the proofcmd production.
	EnterProofcmd(c *ProofcmdContext)

	// EnterLetcmd is called when entering the letcmd production.
	EnterLetcmd(c *LetcmdContext)

	// EnterPathjoin is called when entering the pathjoin production.
	EnterPathjoin(c *PathjoinContext)

	// EnterStatementlist is called when entering the statementlist production.
	EnterStatementlist(c *StatementlistContext)

	// EnterCompound is called when entering the compound production.
	EnterCompound(c *CompoundContext)

	// EnterEmpty is called when entering the empty production.
	EnterEmpty(c *EmptyContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterConstraint is called when entering the constraint production.
	EnterConstraint(c *ConstraintContext)

	// EnterEquation is called when entering the equation production.
	EnterEquation(c *EquationContext)

	// EnterOrientation is called when entering the orientation production.
	EnterOrientation(c *OrientationContext)

	// EnterToken is called when entering the token production.
	EnterToken(c *TokenContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterPathtertiary is called when entering the pathtertiary production.
	EnterPathtertiary(c *PathtertiaryContext)

	// EnterTerm is called when entering the term production.
	EnterTerm(c *TermContext)

	// EnterPath is called when entering the path production.
	EnterPath(c *PathContext)

	// EnterCycle is called when entering the cycle production.
	EnterCycle(c *CycleContext)

	// EnterTransform is called when entering the transform production.
	EnterTransform(c *TransformContext)

	// EnterFactor is called when entering the factor production.
	EnterFactor(c *FactorContext)

	// EnterFuncatom is called when entering the funcatom production.
	EnterFuncatom(c *FuncatomContext)

	// EnterScalaratom is called when entering the scalaratom production.
	EnterScalaratom(c *ScalaratomContext)

	// EnterInterpolation is called when entering the interpolation production.
	EnterInterpolation(c *InterpolationContext)

	// EnterSimpleatom is called when entering the simpleatom production.
	EnterSimpleatom(c *SimpleatomContext)

	// EnterPairpart is called when entering the pairpart production.
	EnterPairpart(c *PairpartContext)

	// EnterPointof is called when entering the pointof production.
	EnterPointof(c *PointofContext)

	// EnterReversepath is called when entering the reversepath production.
	EnterReversepath(c *ReversepathContext)

	// EnterSubpath is called when entering the subpath production.
	EnterSubpath(c *SubpathContext)

	// EnterEdgeconstraint is called when entering the edgeconstraint production.
	EnterEdgeconstraint(c *EdgeconstraintContext)

	// EnterBox is called when entering the box production.
	EnterBox(c *BoxContext)

	// EnterEdgepath is called when entering the edgepath production.
	EnterEdgepath(c *EdgepathContext)

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

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitTypedecl is called when exiting the typedecl production.
	ExitTypedecl(c *TypedeclContext)

	// ExitLocaldecl is called when exiting the localdecl production.
	ExitLocaldecl(c *LocaldeclContext)

	// ExitParameterdecl is called when exiting the parameterdecl production.
	ExitParameterdecl(c *ParameterdeclContext)

	// ExitSavecmd is called when exiting the savecmd production.
	ExitSavecmd(c *SavecmdContext)

	// ExitShowcmd is called when exiting the showcmd production.
	ExitShowcmd(c *ShowcmdContext)

	// ExitProofcmd is called when exiting the proofcmd production.
	ExitProofcmd(c *ProofcmdContext)

	// ExitLetcmd is called when exiting the letcmd production.
	ExitLetcmd(c *LetcmdContext)

	// ExitPathjoin is called when exiting the pathjoin production.
	ExitPathjoin(c *PathjoinContext)

	// ExitStatementlist is called when exiting the statementlist production.
	ExitStatementlist(c *StatementlistContext)

	// ExitCompound is called when exiting the compound production.
	ExitCompound(c *CompoundContext)

	// ExitEmpty is called when exiting the empty production.
	ExitEmpty(c *EmptyContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitConstraint is called when exiting the constraint production.
	ExitConstraint(c *ConstraintContext)

	// ExitEquation is called when exiting the equation production.
	ExitEquation(c *EquationContext)

	// ExitOrientation is called when exiting the orientation production.
	ExitOrientation(c *OrientationContext)

	// ExitToken is called when exiting the token production.
	ExitToken(c *TokenContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitPathtertiary is called when exiting the pathtertiary production.
	ExitPathtertiary(c *PathtertiaryContext)

	// ExitTerm is called when exiting the term production.
	ExitTerm(c *TermContext)

	// ExitPath is called when exiting the path production.
	ExitPath(c *PathContext)

	// ExitCycle is called when exiting the cycle production.
	ExitCycle(c *CycleContext)

	// ExitTransform is called when exiting the transform production.
	ExitTransform(c *TransformContext)

	// ExitFactor is called when exiting the factor production.
	ExitFactor(c *FactorContext)

	// ExitFuncatom is called when exiting the funcatom production.
	ExitFuncatom(c *FuncatomContext)

	// ExitScalaratom is called when exiting the scalaratom production.
	ExitScalaratom(c *ScalaratomContext)

	// ExitInterpolation is called when exiting the interpolation production.
	ExitInterpolation(c *InterpolationContext)

	// ExitSimpleatom is called when exiting the simpleatom production.
	ExitSimpleatom(c *SimpleatomContext)

	// ExitPairpart is called when exiting the pairpart production.
	ExitPairpart(c *PairpartContext)

	// ExitPointof is called when exiting the pointof production.
	ExitPointof(c *PointofContext)

	// ExitReversepath is called when exiting the reversepath production.
	ExitReversepath(c *ReversepathContext)

	// ExitSubpath is called when exiting the subpath production.
	ExitSubpath(c *SubpathContext)

	// ExitEdgeconstraint is called when exiting the edgeconstraint production.
	ExitEdgeconstraint(c *EdgeconstraintContext)

	// ExitBox is called when exiting the box production.
	ExitBox(c *BoxContext)

	// ExitEdgepath is called when exiting the edgepath production.
	ExitEdgepath(c *EdgepathContext)

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
