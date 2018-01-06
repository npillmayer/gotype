// Generated from Gallery.g4 by ANTLR 4.7.

package grammar // Gallery
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseGalleryListener is a complete listener for a parse tree produced by GalleryParser.
type BaseGalleryListener struct{}

var _ GalleryListener = &BaseGalleryListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseGalleryListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseGalleryListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseGalleryListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseGalleryListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStatementlist is called when production statementlist is entered.
func (s *BaseGalleryListener) EnterStatementlist(ctx *StatementlistContext) {}

// ExitStatementlist is called when production statementlist is exited.
func (s *BaseGalleryListener) ExitStatementlist(ctx *StatementlistContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseGalleryListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseGalleryListener) ExitStatement(ctx *StatementContext) {}

// EnterDeclaration is called when production declaration is entered.
func (s *BaseGalleryListener) EnterDeclaration(ctx *DeclarationContext) {}

// ExitDeclaration is called when production declaration is exited.
func (s *BaseGalleryListener) ExitDeclaration(ctx *DeclarationContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BaseGalleryListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BaseGalleryListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterCompound is called when production compound is entered.
func (s *BaseGalleryListener) EnterCompound(ctx *CompoundContext) {}

// ExitCompound is called when production compound is exited.
func (s *BaseGalleryListener) ExitCompound(ctx *CompoundContext) {}

// EnterConstraint is called when production constraint is entered.
func (s *BaseGalleryListener) EnterConstraint(ctx *ConstraintContext) {}

// ExitConstraint is called when production constraint is exited.
func (s *BaseGalleryListener) ExitConstraint(ctx *ConstraintContext) {}

// EnterEquation is called when production equation is entered.
func (s *BaseGalleryListener) EnterEquation(ctx *EquationContext) {}

// ExitEquation is called when production equation is exited.
func (s *BaseGalleryListener) ExitEquation(ctx *EquationContext) {}

// EnterOrientation is called when production orientation is entered.
func (s *BaseGalleryListener) EnterOrientation(ctx *OrientationContext) {}

// ExitOrientation is called when production orientation is exited.
func (s *BaseGalleryListener) ExitOrientation(ctx *OrientationContext) {}

// EnterSavecmd is called when production savecmd is entered.
func (s *BaseGalleryListener) EnterSavecmd(ctx *SavecmdContext) {}

// ExitSavecmd is called when production savecmd is exited.
func (s *BaseGalleryListener) ExitSavecmd(ctx *SavecmdContext) {}

// EnterShowcmd is called when production showcmd is entered.
func (s *BaseGalleryListener) EnterShowcmd(ctx *ShowcmdContext) {}

// ExitShowcmd is called when production showcmd is exited.
func (s *BaseGalleryListener) ExitShowcmd(ctx *ShowcmdContext) {}

// EnterProofcmd is called when production proofcmd is entered.
func (s *BaseGalleryListener) EnterProofcmd(ctx *ProofcmdContext) {}

// ExitProofcmd is called when production proofcmd is exited.
func (s *BaseGalleryListener) ExitProofcmd(ctx *ProofcmdContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseGalleryListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseGalleryListener) ExitExpression(ctx *ExpressionContext) {}

// EnterPathtertiary is called when production pathtertiary is entered.
func (s *BaseGalleryListener) EnterPathtertiary(ctx *PathtertiaryContext) {}

// ExitPathtertiary is called when production pathtertiary is exited.
func (s *BaseGalleryListener) ExitPathtertiary(ctx *PathtertiaryContext) {}

// EnterLonesecondary is called when production lonesecondary is entered.
func (s *BaseGalleryListener) EnterLonesecondary(ctx *LonesecondaryContext) {}

// ExitLonesecondary is called when production lonesecondary is exited.
func (s *BaseGalleryListener) ExitLonesecondary(ctx *LonesecondaryContext) {}

// EnterTerm is called when production term is entered.
func (s *BaseGalleryListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BaseGalleryListener) ExitTerm(ctx *TermContext) {}

// EnterPath is called when production path is entered.
func (s *BaseGalleryListener) EnterPath(ctx *PathContext) {}

// ExitPath is called when production path is exited.
func (s *BaseGalleryListener) ExitPath(ctx *PathContext) {}

// EnterCycle is called when production cycle is entered.
func (s *BaseGalleryListener) EnterCycle(ctx *CycleContext) {}

// ExitCycle is called when production cycle is exited.
func (s *BaseGalleryListener) ExitCycle(ctx *CycleContext) {}

// EnterTransform is called when production transform is entered.
func (s *BaseGalleryListener) EnterTransform(ctx *TransformContext) {}

// ExitTransform is called when production transform is exited.
func (s *BaseGalleryListener) ExitTransform(ctx *TransformContext) {}

// EnterLoneprimary is called when production loneprimary is entered.
func (s *BaseGalleryListener) EnterLoneprimary(ctx *LoneprimaryContext) {}

// ExitLoneprimary is called when production loneprimary is exited.
func (s *BaseGalleryListener) ExitLoneprimary(ctx *LoneprimaryContext) {}

// EnterFactor is called when production factor is entered.
func (s *BaseGalleryListener) EnterFactor(ctx *FactorContext) {}

// ExitFactor is called when production factor is exited.
func (s *BaseGalleryListener) ExitFactor(ctx *FactorContext) {}

// EnterTransformer is called when production transformer is entered.
func (s *BaseGalleryListener) EnterTransformer(ctx *TransformerContext) {}

// ExitTransformer is called when production transformer is exited.
func (s *BaseGalleryListener) ExitTransformer(ctx *TransformerContext) {}

// EnterFuncnumatom is called when production funcnumatom is entered.
func (s *BaseGalleryListener) EnterFuncnumatom(ctx *FuncnumatomContext) {}

// ExitFuncnumatom is called when production funcnumatom is exited.
func (s *BaseGalleryListener) ExitFuncnumatom(ctx *FuncnumatomContext) {}

// EnterScalarnumatom is called when production scalarnumatom is entered.
func (s *BaseGalleryListener) EnterScalarnumatom(ctx *ScalarnumatomContext) {}

// ExitScalarnumatom is called when production scalarnumatom is exited.
func (s *BaseGalleryListener) ExitScalarnumatom(ctx *ScalarnumatomContext) {}

// EnterInterpolation is called when production interpolation is entered.
func (s *BaseGalleryListener) EnterInterpolation(ctx *InterpolationContext) {}

// ExitInterpolation is called when production interpolation is exited.
func (s *BaseGalleryListener) ExitInterpolation(ctx *InterpolationContext) {}

// EnterSimplenumatom is called when production simplenumatom is entered.
func (s *BaseGalleryListener) EnterSimplenumatom(ctx *SimplenumatomContext) {}

// ExitSimplenumatom is called when production simplenumatom is exited.
func (s *BaseGalleryListener) ExitSimplenumatom(ctx *SimplenumatomContext) {}

// EnterPairpart is called when production pairpart is entered.
func (s *BaseGalleryListener) EnterPairpart(ctx *PairpartContext) {}

// ExitPairpart is called when production pairpart is exited.
func (s *BaseGalleryListener) ExitPairpart(ctx *PairpartContext) {}

// EnterPathpoint is called when production pathpoint is entered.
func (s *BaseGalleryListener) EnterPathpoint(ctx *PathpointContext) {}

// ExitPathpoint is called when production pathpoint is exited.
func (s *BaseGalleryListener) ExitPathpoint(ctx *PathpointContext) {}

// EnterReversepath is called when production reversepath is entered.
func (s *BaseGalleryListener) EnterReversepath(ctx *ReversepathContext) {}

// ExitReversepath is called when production reversepath is exited.
func (s *BaseGalleryListener) ExitReversepath(ctx *ReversepathContext) {}

// EnterSubpath is called when production subpath is entered.
func (s *BaseGalleryListener) EnterSubpath(ctx *SubpathContext) {}

// ExitSubpath is called when production subpath is exited.
func (s *BaseGalleryListener) ExitSubpath(ctx *SubpathContext) {}

// EnterEdgeconstraint is called when production edgeconstraint is entered.
func (s *BaseGalleryListener) EnterEdgeconstraint(ctx *EdgeconstraintContext) {}

// ExitEdgeconstraint is called when production edgeconstraint is exited.
func (s *BaseGalleryListener) ExitEdgeconstraint(ctx *EdgeconstraintContext) {}

// EnterBox is called when production box is entered.
func (s *BaseGalleryListener) EnterBox(ctx *BoxContext) {}

// ExitBox is called when production box is exited.
func (s *BaseGalleryListener) ExitBox(ctx *BoxContext) {}

// EnterScalarmulop is called when production scalarmulop is entered.
func (s *BaseGalleryListener) EnterScalarmulop(ctx *ScalarmulopContext) {}

// ExitScalarmulop is called when production scalarmulop is exited.
func (s *BaseGalleryListener) ExitScalarmulop(ctx *ScalarmulopContext) {}

// EnterNumtokenatom is called when production numtokenatom is entered.
func (s *BaseGalleryListener) EnterNumtokenatom(ctx *NumtokenatomContext) {}

// ExitNumtokenatom is called when production numtokenatom is exited.
func (s *BaseGalleryListener) ExitNumtokenatom(ctx *NumtokenatomContext) {}

// EnterDecimal is called when production decimal is entered.
func (s *BaseGalleryListener) EnterDecimal(ctx *DecimalContext) {}

// ExitDecimal is called when production decimal is exited.
func (s *BaseGalleryListener) ExitDecimal(ctx *DecimalContext) {}

// EnterVaratom is called when production varatom is entered.
func (s *BaseGalleryListener) EnterVaratom(ctx *VaratomContext) {}

// ExitVaratom is called when production varatom is exited.
func (s *BaseGalleryListener) ExitVaratom(ctx *VaratomContext) {}

// EnterLiteralpair is called when production literalpair is entered.
func (s *BaseGalleryListener) EnterLiteralpair(ctx *LiteralpairContext) {}

// ExitLiteralpair is called when production literalpair is exited.
func (s *BaseGalleryListener) ExitLiteralpair(ctx *LiteralpairContext) {}

// EnterSubexpression is called when production subexpression is entered.
func (s *BaseGalleryListener) EnterSubexpression(ctx *SubexpressionContext) {}

// ExitSubexpression is called when production subexpression is exited.
func (s *BaseGalleryListener) ExitSubexpression(ctx *SubexpressionContext) {}

// EnterExprgroup is called when production exprgroup is entered.
func (s *BaseGalleryListener) EnterExprgroup(ctx *ExprgroupContext) {}

// ExitExprgroup is called when production exprgroup is exited.
func (s *BaseGalleryListener) ExitExprgroup(ctx *ExprgroupContext) {}

// EnterVariable is called when production variable is entered.
func (s *BaseGalleryListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production variable is exited.
func (s *BaseGalleryListener) ExitVariable(ctx *VariableContext) {}

// EnterSubscript is called when production subscript is entered.
func (s *BaseGalleryListener) EnterSubscript(ctx *SubscriptContext) {}

// ExitSubscript is called when production subscript is exited.
func (s *BaseGalleryListener) ExitSubscript(ctx *SubscriptContext) {}

// EnterAnytag is called when production anytag is entered.
func (s *BaseGalleryListener) EnterAnytag(ctx *AnytagContext) {}

// ExitAnytag is called when production anytag is exited.
func (s *BaseGalleryListener) ExitAnytag(ctx *AnytagContext) {}
