// Generated from CG.g4 by ANTLR 4.7.

package cg // CG
import "github.com/antlr/antlr4/runtime/Go/antlr"

// CGListener is a complete listener for a parse tree produced by CGParser.
type CGListener interface {
	antlr.ParseTreeListener

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterSubexpression is called when entering the subexpression production.
	EnterSubexpression(c *SubexpressionContext)

	// EnterDefinitionstatus is called when entering the definitionstatus production.
	EnterDefinitionstatus(c *DefinitionstatusContext)

	// EnterEquivalentto is called when entering the equivalentto production.
	EnterEquivalentto(c *EquivalenttoContext)

	// EnterSubtypeof is called when entering the subtypeof production.
	EnterSubtypeof(c *SubtypeofContext)

	// EnterFocusconcept is called when entering the focusconcept production.
	EnterFocusconcept(c *FocusconceptContext)

	// EnterConceptreference is called when entering the conceptreference production.
	EnterConceptreference(c *ConceptreferenceContext)

	// EnterConceptid is called when entering the conceptid production.
	EnterConceptid(c *ConceptidContext)

	// EnterTerm is called when entering the term production.
	EnterTerm(c *TermContext)

	// EnterRefinement is called when entering the refinement production.
	EnterRefinement(c *RefinementContext)

	// EnterAttributegroup is called when entering the attributegroup production.
	EnterAttributegroup(c *AttributegroupContext)

	// EnterAttributeset is called when entering the attributeset production.
	EnterAttributeset(c *AttributesetContext)

	// EnterAttribute is called when entering the attribute production.
	EnterAttribute(c *AttributeContext)

	// EnterAttributename is called when entering the attributename production.
	EnterAttributename(c *AttributenameContext)

	// EnterAttributevalue is called when entering the attributevalue production.
	EnterAttributevalue(c *AttributevalueContext)

	// EnterExpressionvalue is called when entering the expressionvalue production.
	EnterExpressionvalue(c *ExpressionvalueContext)

	// EnterStringvalue is called when entering the stringvalue production.
	EnterStringvalue(c *StringvalueContext)

	// EnterNumericvalue is called when entering the numericvalue production.
	EnterNumericvalue(c *NumericvalueContext)

	// EnterIntegervalue is called when entering the integervalue production.
	EnterIntegervalue(c *IntegervalueContext)

	// EnterDecimalvalue is called when entering the decimalvalue production.
	EnterDecimalvalue(c *DecimalvalueContext)

	// EnterSctid is called when entering the sctid production.
	EnterSctid(c *SctidContext)

	// EnterWs is called when entering the ws production.
	EnterWs(c *WsContext)

	// EnterSp is called when entering the sp production.
	EnterSp(c *SpContext)

	// EnterHtab is called when entering the htab production.
	EnterHtab(c *HtabContext)

	// EnterCr is called when entering the cr production.
	EnterCr(c *CrContext)

	// EnterLf is called when entering the lf production.
	EnterLf(c *LfContext)

	// EnterQm is called when entering the qm production.
	EnterQm(c *QmContext)

	// EnterBs is called when entering the bs production.
	EnterBs(c *BsContext)

	// EnterDigit is called when entering the digit production.
	EnterDigit(c *DigitContext)

	// EnterZero is called when entering the zero production.
	EnterZero(c *ZeroContext)

	// EnterDigitnonzero is called when entering the digitnonzero production.
	EnterDigitnonzero(c *DigitnonzeroContext)

	// EnterNonwsnonpipe is called when entering the nonwsnonpipe production.
	EnterNonwsnonpipe(c *NonwsnonpipeContext)

	// EnterAnynonescapedchar is called when entering the anynonescapedchar production.
	EnterAnynonescapedchar(c *AnynonescapedcharContext)

	// EnterEscapedchar is called when entering the escapedchar production.
	EnterEscapedchar(c *EscapedcharContext)

	// EnterUtf8_2 is called when entering the utf8_2 production.
	EnterUtf8_2(c *Utf8_2Context)

	// EnterUtf8_3 is called when entering the utf8_3 production.
	EnterUtf8_3(c *Utf8_3Context)

	// EnterUtf8_4 is called when entering the utf8_4 production.
	EnterUtf8_4(c *Utf8_4Context)

	// EnterUtf8_tail is called when entering the utf8_tail production.
	EnterUtf8_tail(c *Utf8_tailContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitSubexpression is called when exiting the subexpression production.
	ExitSubexpression(c *SubexpressionContext)

	// ExitDefinitionstatus is called when exiting the definitionstatus production.
	ExitDefinitionstatus(c *DefinitionstatusContext)

	// ExitEquivalentto is called when exiting the equivalentto production.
	ExitEquivalentto(c *EquivalenttoContext)

	// ExitSubtypeof is called when exiting the subtypeof production.
	ExitSubtypeof(c *SubtypeofContext)

	// ExitFocusconcept is called when exiting the focusconcept production.
	ExitFocusconcept(c *FocusconceptContext)

	// ExitConceptreference is called when exiting the conceptreference production.
	ExitConceptreference(c *ConceptreferenceContext)

	// ExitConceptid is called when exiting the conceptid production.
	ExitConceptid(c *ConceptidContext)

	// ExitTerm is called when exiting the term production.
	ExitTerm(c *TermContext)

	// ExitRefinement is called when exiting the refinement production.
	ExitRefinement(c *RefinementContext)

	// ExitAttributegroup is called when exiting the attributegroup production.
	ExitAttributegroup(c *AttributegroupContext)

	// ExitAttributeset is called when exiting the attributeset production.
	ExitAttributeset(c *AttributesetContext)

	// ExitAttribute is called when exiting the attribute production.
	ExitAttribute(c *AttributeContext)

	// ExitAttributename is called when exiting the attributename production.
	ExitAttributename(c *AttributenameContext)

	// ExitAttributevalue is called when exiting the attributevalue production.
	ExitAttributevalue(c *AttributevalueContext)

	// ExitExpressionvalue is called when exiting the expressionvalue production.
	ExitExpressionvalue(c *ExpressionvalueContext)

	// ExitStringvalue is called when exiting the stringvalue production.
	ExitStringvalue(c *StringvalueContext)

	// ExitNumericvalue is called when exiting the numericvalue production.
	ExitNumericvalue(c *NumericvalueContext)

	// ExitIntegervalue is called when exiting the integervalue production.
	ExitIntegervalue(c *IntegervalueContext)

	// ExitDecimalvalue is called when exiting the decimalvalue production.
	ExitDecimalvalue(c *DecimalvalueContext)

	// ExitSctid is called when exiting the sctid production.
	ExitSctid(c *SctidContext)

	// ExitWs is called when exiting the ws production.
	ExitWs(c *WsContext)

	// ExitSp is called when exiting the sp production.
	ExitSp(c *SpContext)

	// ExitHtab is called when exiting the htab production.
	ExitHtab(c *HtabContext)

	// ExitCr is called when exiting the cr production.
	ExitCr(c *CrContext)

	// ExitLf is called when exiting the lf production.
	ExitLf(c *LfContext)

	// ExitQm is called when exiting the qm production.
	ExitQm(c *QmContext)

	// ExitBs is called when exiting the bs production.
	ExitBs(c *BsContext)

	// ExitDigit is called when exiting the digit production.
	ExitDigit(c *DigitContext)

	// ExitZero is called when exiting the zero production.
	ExitZero(c *ZeroContext)

	// ExitDigitnonzero is called when exiting the digitnonzero production.
	ExitDigitnonzero(c *DigitnonzeroContext)

	// ExitNonwsnonpipe is called when exiting the nonwsnonpipe production.
	ExitNonwsnonpipe(c *NonwsnonpipeContext)

	// ExitAnynonescapedchar is called when exiting the anynonescapedchar production.
	ExitAnynonescapedchar(c *AnynonescapedcharContext)

	// ExitEscapedchar is called when exiting the escapedchar production.
	ExitEscapedchar(c *EscapedcharContext)

	// ExitUtf8_2 is called when exiting the utf8_2 production.
	ExitUtf8_2(c *Utf8_2Context)

	// ExitUtf8_3 is called when exiting the utf8_3 production.
	ExitUtf8_3(c *Utf8_3Context)

	// ExitUtf8_4 is called when exiting the utf8_4 production.
	ExitUtf8_4(c *Utf8_4Context)

	// ExitUtf8_tail is called when exiting the utf8_tail production.
	ExitUtf8_tail(c *Utf8_tailContext)
}
