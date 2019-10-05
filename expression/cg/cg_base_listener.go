// Code generated from CG.g4 by ANTLR 4.7.2. DO NOT EDIT.

package cg // CG
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseCGListener is a complete listener for a parse tree produced by CGParser.
type BaseCGListener struct{}

var _ CGListener = &BaseCGListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseCGListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseCGListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseCGListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseCGListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseCGListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseCGListener) ExitExpression(ctx *ExpressionContext) {}

// EnterSubexpression is called when production subexpression is entered.
func (s *BaseCGListener) EnterSubexpression(ctx *SubexpressionContext) {}

// ExitSubexpression is called when production subexpression is exited.
func (s *BaseCGListener) ExitSubexpression(ctx *SubexpressionContext) {}

// EnterDefinitionstatus is called when production definitionstatus is entered.
func (s *BaseCGListener) EnterDefinitionstatus(ctx *DefinitionstatusContext) {}

// ExitDefinitionstatus is called when production definitionstatus is exited.
func (s *BaseCGListener) ExitDefinitionstatus(ctx *DefinitionstatusContext) {}

// EnterEquivalentto is called when production equivalentto is entered.
func (s *BaseCGListener) EnterEquivalentto(ctx *EquivalenttoContext) {}

// ExitEquivalentto is called when production equivalentto is exited.
func (s *BaseCGListener) ExitEquivalentto(ctx *EquivalenttoContext) {}

// EnterSubtypeof is called when production subtypeof is entered.
func (s *BaseCGListener) EnterSubtypeof(ctx *SubtypeofContext) {}

// ExitSubtypeof is called when production subtypeof is exited.
func (s *BaseCGListener) ExitSubtypeof(ctx *SubtypeofContext) {}

// EnterFocusconcept is called when production focusconcept is entered.
func (s *BaseCGListener) EnterFocusconcept(ctx *FocusconceptContext) {}

// ExitFocusconcept is called when production focusconcept is exited.
func (s *BaseCGListener) ExitFocusconcept(ctx *FocusconceptContext) {}

// EnterConceptreference is called when production conceptreference is entered.
func (s *BaseCGListener) EnterConceptreference(ctx *ConceptreferenceContext) {}

// ExitConceptreference is called when production conceptreference is exited.
func (s *BaseCGListener) ExitConceptreference(ctx *ConceptreferenceContext) {}

// EnterConceptid is called when production conceptid is entered.
func (s *BaseCGListener) EnterConceptid(ctx *ConceptidContext) {}

// ExitConceptid is called when production conceptid is exited.
func (s *BaseCGListener) ExitConceptid(ctx *ConceptidContext) {}

// EnterTerm is called when production term is entered.
func (s *BaseCGListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BaseCGListener) ExitTerm(ctx *TermContext) {}

// EnterRefinement is called when production refinement is entered.
func (s *BaseCGListener) EnterRefinement(ctx *RefinementContext) {}

// ExitRefinement is called when production refinement is exited.
func (s *BaseCGListener) ExitRefinement(ctx *RefinementContext) {}

// EnterAttributegroup is called when production attributegroup is entered.
func (s *BaseCGListener) EnterAttributegroup(ctx *AttributegroupContext) {}

// ExitAttributegroup is called when production attributegroup is exited.
func (s *BaseCGListener) ExitAttributegroup(ctx *AttributegroupContext) {}

// EnterAttributeset is called when production attributeset is entered.
func (s *BaseCGListener) EnterAttributeset(ctx *AttributesetContext) {}

// ExitAttributeset is called when production attributeset is exited.
func (s *BaseCGListener) ExitAttributeset(ctx *AttributesetContext) {}

// EnterAttribute is called when production attribute is entered.
func (s *BaseCGListener) EnterAttribute(ctx *AttributeContext) {}

// ExitAttribute is called when production attribute is exited.
func (s *BaseCGListener) ExitAttribute(ctx *AttributeContext) {}

// EnterAttributename is called when production attributename is entered.
func (s *BaseCGListener) EnterAttributename(ctx *AttributenameContext) {}

// ExitAttributename is called when production attributename is exited.
func (s *BaseCGListener) ExitAttributename(ctx *AttributenameContext) {}

// EnterAttributevalue is called when production attributevalue is entered.
func (s *BaseCGListener) EnterAttributevalue(ctx *AttributevalueContext) {}

// ExitAttributevalue is called when production attributevalue is exited.
func (s *BaseCGListener) ExitAttributevalue(ctx *AttributevalueContext) {}

// EnterExpressionvalue is called when production expressionvalue is entered.
func (s *BaseCGListener) EnterExpressionvalue(ctx *ExpressionvalueContext) {}

// ExitExpressionvalue is called when production expressionvalue is exited.
func (s *BaseCGListener) ExitExpressionvalue(ctx *ExpressionvalueContext) {}

// EnterStringvalue is called when production stringvalue is entered.
func (s *BaseCGListener) EnterStringvalue(ctx *StringvalueContext) {}

// ExitStringvalue is called when production stringvalue is exited.
func (s *BaseCGListener) ExitStringvalue(ctx *StringvalueContext) {}

// EnterNumericvalue is called when production numericvalue is entered.
func (s *BaseCGListener) EnterNumericvalue(ctx *NumericvalueContext) {}

// ExitNumericvalue is called when production numericvalue is exited.
func (s *BaseCGListener) ExitNumericvalue(ctx *NumericvalueContext) {}

// EnterIntegervalue is called when production integervalue is entered.
func (s *BaseCGListener) EnterIntegervalue(ctx *IntegervalueContext) {}

// ExitIntegervalue is called when production integervalue is exited.
func (s *BaseCGListener) ExitIntegervalue(ctx *IntegervalueContext) {}

// EnterDecimalvalue is called when production decimalvalue is entered.
func (s *BaseCGListener) EnterDecimalvalue(ctx *DecimalvalueContext) {}

// ExitDecimalvalue is called when production decimalvalue is exited.
func (s *BaseCGListener) ExitDecimalvalue(ctx *DecimalvalueContext) {}

// EnterSctid is called when production sctid is entered.
func (s *BaseCGListener) EnterSctid(ctx *SctidContext) {}

// ExitSctid is called when production sctid is exited.
func (s *BaseCGListener) ExitSctid(ctx *SctidContext) {}

// EnterWs is called when production ws is entered.
func (s *BaseCGListener) EnterWs(ctx *WsContext) {}

// ExitWs is called when production ws is exited.
func (s *BaseCGListener) ExitWs(ctx *WsContext) {}

// EnterSp is called when production sp is entered.
func (s *BaseCGListener) EnterSp(ctx *SpContext) {}

// ExitSp is called when production sp is exited.
func (s *BaseCGListener) ExitSp(ctx *SpContext) {}

// EnterHtab is called when production htab is entered.
func (s *BaseCGListener) EnterHtab(ctx *HtabContext) {}

// ExitHtab is called when production htab is exited.
func (s *BaseCGListener) ExitHtab(ctx *HtabContext) {}

// EnterCr is called when production cr is entered.
func (s *BaseCGListener) EnterCr(ctx *CrContext) {}

// ExitCr is called when production cr is exited.
func (s *BaseCGListener) ExitCr(ctx *CrContext) {}

// EnterLf is called when production lf is entered.
func (s *BaseCGListener) EnterLf(ctx *LfContext) {}

// ExitLf is called when production lf is exited.
func (s *BaseCGListener) ExitLf(ctx *LfContext) {}

// EnterQm is called when production qm is entered.
func (s *BaseCGListener) EnterQm(ctx *QmContext) {}

// ExitQm is called when production qm is exited.
func (s *BaseCGListener) ExitQm(ctx *QmContext) {}

// EnterBs is called when production bs is entered.
func (s *BaseCGListener) EnterBs(ctx *BsContext) {}

// ExitBs is called when production bs is exited.
func (s *BaseCGListener) ExitBs(ctx *BsContext) {}

// EnterDigit is called when production digit is entered.
func (s *BaseCGListener) EnterDigit(ctx *DigitContext) {}

// ExitDigit is called when production digit is exited.
func (s *BaseCGListener) ExitDigit(ctx *DigitContext) {}

// EnterZero is called when production zero is entered.
func (s *BaseCGListener) EnterZero(ctx *ZeroContext) {}

// ExitZero is called when production zero is exited.
func (s *BaseCGListener) ExitZero(ctx *ZeroContext) {}

// EnterDigitnonzero is called when production digitnonzero is entered.
func (s *BaseCGListener) EnterDigitnonzero(ctx *DigitnonzeroContext) {}

// ExitDigitnonzero is called when production digitnonzero is exited.
func (s *BaseCGListener) ExitDigitnonzero(ctx *DigitnonzeroContext) {}

// EnterNonwsnonpipe is called when production nonwsnonpipe is entered.
func (s *BaseCGListener) EnterNonwsnonpipe(ctx *NonwsnonpipeContext) {}

// ExitNonwsnonpipe is called when production nonwsnonpipe is exited.
func (s *BaseCGListener) ExitNonwsnonpipe(ctx *NonwsnonpipeContext) {}

// EnterAnynonescapedchar is called when production anynonescapedchar is entered.
func (s *BaseCGListener) EnterAnynonescapedchar(ctx *AnynonescapedcharContext) {}

// ExitAnynonescapedchar is called when production anynonescapedchar is exited.
func (s *BaseCGListener) ExitAnynonescapedchar(ctx *AnynonescapedcharContext) {}

// EnterEscapedchar is called when production escapedchar is entered.
func (s *BaseCGListener) EnterEscapedchar(ctx *EscapedcharContext) {}

// ExitEscapedchar is called when production escapedchar is exited.
func (s *BaseCGListener) ExitEscapedchar(ctx *EscapedcharContext) {}

// EnterUtf8_2 is called when production utf8_2 is entered.
func (s *BaseCGListener) EnterUtf8_2(ctx *Utf8_2Context) {}

// ExitUtf8_2 is called when production utf8_2 is exited.
func (s *BaseCGListener) ExitUtf8_2(ctx *Utf8_2Context) {}

// EnterUtf8_3 is called when production utf8_3 is entered.
func (s *BaseCGListener) EnterUtf8_3(ctx *Utf8_3Context) {}

// ExitUtf8_3 is called when production utf8_3 is exited.
func (s *BaseCGListener) ExitUtf8_3(ctx *Utf8_3Context) {}

// EnterUtf8_4 is called when production utf8_4 is entered.
func (s *BaseCGListener) EnterUtf8_4(ctx *Utf8_4Context) {}

// ExitUtf8_4 is called when production utf8_4 is exited.
func (s *BaseCGListener) ExitUtf8_4(ctx *Utf8_4Context) {}

// EnterUtf8_tail is called when production utf8_tail is entered.
func (s *BaseCGListener) EnterUtf8_tail(ctx *Utf8_tailContext) {}

// ExitUtf8_tail is called when production utf8_tail is exited.
func (s *BaseCGListener) ExitUtf8_tail(ctx *Utf8_tailContext) {}
