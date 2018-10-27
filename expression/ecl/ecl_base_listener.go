// Generated from ECL.g4 by ANTLR 4.7.

package ecl // ECL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseECLListener is a complete listener for a parse tree produced by ECLParser.
type BaseECLListener struct{}

var _ ECLListener = &BaseECLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseECLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseECLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseECLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseECLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterExpressionconstraint is called when production expressionconstraint is entered.
func (s *BaseECLListener) EnterExpressionconstraint(ctx *ExpressionconstraintContext) {}

// ExitExpressionconstraint is called when production expressionconstraint is exited.
func (s *BaseECLListener) ExitExpressionconstraint(ctx *ExpressionconstraintContext) {}

// EnterRefinedexpressionconstraint is called when production refinedexpressionconstraint is entered.
func (s *BaseECLListener) EnterRefinedexpressionconstraint(ctx *RefinedexpressionconstraintContext) {}

// ExitRefinedexpressionconstraint is called when production refinedexpressionconstraint is exited.
func (s *BaseECLListener) ExitRefinedexpressionconstraint(ctx *RefinedexpressionconstraintContext) {}

// EnterCompoundexpressionconstraint is called when production compoundexpressionconstraint is entered.
func (s *BaseECLListener) EnterCompoundexpressionconstraint(ctx *CompoundexpressionconstraintContext) {
}

// ExitCompoundexpressionconstraint is called when production compoundexpressionconstraint is exited.
func (s *BaseECLListener) ExitCompoundexpressionconstraint(ctx *CompoundexpressionconstraintContext) {}

// EnterConjunctionexpressionconstraint is called when production conjunctionexpressionconstraint is entered.
func (s *BaseECLListener) EnterConjunctionexpressionconstraint(ctx *ConjunctionexpressionconstraintContext) {
}

// ExitConjunctionexpressionconstraint is called when production conjunctionexpressionconstraint is exited.
func (s *BaseECLListener) ExitConjunctionexpressionconstraint(ctx *ConjunctionexpressionconstraintContext) {
}

// EnterDisjunctionexpressionconstraint is called when production disjunctionexpressionconstraint is entered.
func (s *BaseECLListener) EnterDisjunctionexpressionconstraint(ctx *DisjunctionexpressionconstraintContext) {
}

// ExitDisjunctionexpressionconstraint is called when production disjunctionexpressionconstraint is exited.
func (s *BaseECLListener) ExitDisjunctionexpressionconstraint(ctx *DisjunctionexpressionconstraintContext) {
}

// EnterExclusionexpressionconstraint is called when production exclusionexpressionconstraint is entered.
func (s *BaseECLListener) EnterExclusionexpressionconstraint(ctx *ExclusionexpressionconstraintContext) {
}

// ExitExclusionexpressionconstraint is called when production exclusionexpressionconstraint is exited.
func (s *BaseECLListener) ExitExclusionexpressionconstraint(ctx *ExclusionexpressionconstraintContext) {
}

// EnterDottedexpressionconstraint is called when production dottedexpressionconstraint is entered.
func (s *BaseECLListener) EnterDottedexpressionconstraint(ctx *DottedexpressionconstraintContext) {}

// ExitDottedexpressionconstraint is called when production dottedexpressionconstraint is exited.
func (s *BaseECLListener) ExitDottedexpressionconstraint(ctx *DottedexpressionconstraintContext) {}

// EnterDottedexpressionattribute is called when production dottedexpressionattribute is entered.
func (s *BaseECLListener) EnterDottedexpressionattribute(ctx *DottedexpressionattributeContext) {}

// ExitDottedexpressionattribute is called when production dottedexpressionattribute is exited.
func (s *BaseECLListener) ExitDottedexpressionattribute(ctx *DottedexpressionattributeContext) {}

// EnterSubexpressionconstraint is called when production subexpressionconstraint is entered.
func (s *BaseECLListener) EnterSubexpressionconstraint(ctx *SubexpressionconstraintContext) {}

// ExitSubexpressionconstraint is called when production subexpressionconstraint is exited.
func (s *BaseECLListener) ExitSubexpressionconstraint(ctx *SubexpressionconstraintContext) {}

// EnterEclfocusconcept is called when production eclfocusconcept is entered.
func (s *BaseECLListener) EnterEclfocusconcept(ctx *EclfocusconceptContext) {}

// ExitEclfocusconcept is called when production eclfocusconcept is exited.
func (s *BaseECLListener) ExitEclfocusconcept(ctx *EclfocusconceptContext) {}

// EnterDot is called when production dot is entered.
func (s *BaseECLListener) EnterDot(ctx *DotContext) {}

// ExitDot is called when production dot is exited.
func (s *BaseECLListener) ExitDot(ctx *DotContext) {}

// EnterMemberof is called when production memberof is entered.
func (s *BaseECLListener) EnterMemberof(ctx *MemberofContext) {}

// ExitMemberof is called when production memberof is exited.
func (s *BaseECLListener) ExitMemberof(ctx *MemberofContext) {}

// EnterEclconceptreference is called when production eclconceptreference is entered.
func (s *BaseECLListener) EnterEclconceptreference(ctx *EclconceptreferenceContext) {}

// ExitEclconceptreference is called when production eclconceptreference is exited.
func (s *BaseECLListener) ExitEclconceptreference(ctx *EclconceptreferenceContext) {}

// EnterConceptid is called when production conceptid is entered.
func (s *BaseECLListener) EnterConceptid(ctx *ConceptidContext) {}

// ExitConceptid is called when production conceptid is exited.
func (s *BaseECLListener) ExitConceptid(ctx *ConceptidContext) {}

// EnterTerm is called when production term is entered.
func (s *BaseECLListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BaseECLListener) ExitTerm(ctx *TermContext) {}

// EnterWildcard is called when production wildcard is entered.
func (s *BaseECLListener) EnterWildcard(ctx *WildcardContext) {}

// ExitWildcard is called when production wildcard is exited.
func (s *BaseECLListener) ExitWildcard(ctx *WildcardContext) {}

// EnterConstraintoperator is called when production constraintoperator is entered.
func (s *BaseECLListener) EnterConstraintoperator(ctx *ConstraintoperatorContext) {}

// ExitConstraintoperator is called when production constraintoperator is exited.
func (s *BaseECLListener) ExitConstraintoperator(ctx *ConstraintoperatorContext) {}

// EnterDescendantof is called when production descendantof is entered.
func (s *BaseECLListener) EnterDescendantof(ctx *DescendantofContext) {}

// ExitDescendantof is called when production descendantof is exited.
func (s *BaseECLListener) ExitDescendantof(ctx *DescendantofContext) {}

// EnterDescendantorselfof is called when production descendantorselfof is entered.
func (s *BaseECLListener) EnterDescendantorselfof(ctx *DescendantorselfofContext) {}

// ExitDescendantorselfof is called when production descendantorselfof is exited.
func (s *BaseECLListener) ExitDescendantorselfof(ctx *DescendantorselfofContext) {}

// EnterChildof is called when production childof is entered.
func (s *BaseECLListener) EnterChildof(ctx *ChildofContext) {}

// ExitChildof is called when production childof is exited.
func (s *BaseECLListener) ExitChildof(ctx *ChildofContext) {}

// EnterAncestorof is called when production ancestorof is entered.
func (s *BaseECLListener) EnterAncestorof(ctx *AncestorofContext) {}

// ExitAncestorof is called when production ancestorof is exited.
func (s *BaseECLListener) ExitAncestorof(ctx *AncestorofContext) {}

// EnterAncestororselfof is called when production ancestororselfof is entered.
func (s *BaseECLListener) EnterAncestororselfof(ctx *AncestororselfofContext) {}

// ExitAncestororselfof is called when production ancestororselfof is exited.
func (s *BaseECLListener) ExitAncestororselfof(ctx *AncestororselfofContext) {}

// EnterParentof is called when production parentof is entered.
func (s *BaseECLListener) EnterParentof(ctx *ParentofContext) {}

// ExitParentof is called when production parentof is exited.
func (s *BaseECLListener) ExitParentof(ctx *ParentofContext) {}

// EnterConjunction is called when production conjunction is entered.
func (s *BaseECLListener) EnterConjunction(ctx *ConjunctionContext) {}

// ExitConjunction is called when production conjunction is exited.
func (s *BaseECLListener) ExitConjunction(ctx *ConjunctionContext) {}

// EnterDisjunction is called when production disjunction is entered.
func (s *BaseECLListener) EnterDisjunction(ctx *DisjunctionContext) {}

// ExitDisjunction is called when production disjunction is exited.
func (s *BaseECLListener) ExitDisjunction(ctx *DisjunctionContext) {}

// EnterExclusion is called when production exclusion is entered.
func (s *BaseECLListener) EnterExclusion(ctx *ExclusionContext) {}

// ExitExclusion is called when production exclusion is exited.
func (s *BaseECLListener) ExitExclusion(ctx *ExclusionContext) {}

// EnterEclrefinement is called when production eclrefinement is entered.
func (s *BaseECLListener) EnterEclrefinement(ctx *EclrefinementContext) {}

// ExitEclrefinement is called when production eclrefinement is exited.
func (s *BaseECLListener) ExitEclrefinement(ctx *EclrefinementContext) {}

// EnterConjunctionrefinementset is called when production conjunctionrefinementset is entered.
func (s *BaseECLListener) EnterConjunctionrefinementset(ctx *ConjunctionrefinementsetContext) {}

// ExitConjunctionrefinementset is called when production conjunctionrefinementset is exited.
func (s *BaseECLListener) ExitConjunctionrefinementset(ctx *ConjunctionrefinementsetContext) {}

// EnterDisjunctionrefinementset is called when production disjunctionrefinementset is entered.
func (s *BaseECLListener) EnterDisjunctionrefinementset(ctx *DisjunctionrefinementsetContext) {}

// ExitDisjunctionrefinementset is called when production disjunctionrefinementset is exited.
func (s *BaseECLListener) ExitDisjunctionrefinementset(ctx *DisjunctionrefinementsetContext) {}

// EnterSubrefinement is called when production subrefinement is entered.
func (s *BaseECLListener) EnterSubrefinement(ctx *SubrefinementContext) {}

// ExitSubrefinement is called when production subrefinement is exited.
func (s *BaseECLListener) ExitSubrefinement(ctx *SubrefinementContext) {}

// EnterEclattributeset is called when production eclattributeset is entered.
func (s *BaseECLListener) EnterEclattributeset(ctx *EclattributesetContext) {}

// ExitEclattributeset is called when production eclattributeset is exited.
func (s *BaseECLListener) ExitEclattributeset(ctx *EclattributesetContext) {}

// EnterConjunctionattributeset is called when production conjunctionattributeset is entered.
func (s *BaseECLListener) EnterConjunctionattributeset(ctx *ConjunctionattributesetContext) {}

// ExitConjunctionattributeset is called when production conjunctionattributeset is exited.
func (s *BaseECLListener) ExitConjunctionattributeset(ctx *ConjunctionattributesetContext) {}

// EnterDisjunctionattributeset is called when production disjunctionattributeset is entered.
func (s *BaseECLListener) EnterDisjunctionattributeset(ctx *DisjunctionattributesetContext) {}

// ExitDisjunctionattributeset is called when production disjunctionattributeset is exited.
func (s *BaseECLListener) ExitDisjunctionattributeset(ctx *DisjunctionattributesetContext) {}

// EnterSubattributeset is called when production subattributeset is entered.
func (s *BaseECLListener) EnterSubattributeset(ctx *SubattributesetContext) {}

// ExitSubattributeset is called when production subattributeset is exited.
func (s *BaseECLListener) ExitSubattributeset(ctx *SubattributesetContext) {}

// EnterEclattributegroup is called when production eclattributegroup is entered.
func (s *BaseECLListener) EnterEclattributegroup(ctx *EclattributegroupContext) {}

// ExitEclattributegroup is called when production eclattributegroup is exited.
func (s *BaseECLListener) ExitEclattributegroup(ctx *EclattributegroupContext) {}

// EnterEclattribute is called when production eclattribute is entered.
func (s *BaseECLListener) EnterEclattribute(ctx *EclattributeContext) {}

// ExitEclattribute is called when production eclattribute is exited.
func (s *BaseECLListener) ExitEclattribute(ctx *EclattributeContext) {}

// EnterCardinality is called when production cardinality is entered.
func (s *BaseECLListener) EnterCardinality(ctx *CardinalityContext) {}

// ExitCardinality is called when production cardinality is exited.
func (s *BaseECLListener) ExitCardinality(ctx *CardinalityContext) {}

// EnterMinvalue is called when production minvalue is entered.
func (s *BaseECLListener) EnterMinvalue(ctx *MinvalueContext) {}

// ExitMinvalue is called when production minvalue is exited.
func (s *BaseECLListener) ExitMinvalue(ctx *MinvalueContext) {}

// EnterTo is called when production to is entered.
func (s *BaseECLListener) EnterTo(ctx *ToContext) {}

// ExitTo is called when production to is exited.
func (s *BaseECLListener) ExitTo(ctx *ToContext) {}

// EnterMaxvalue is called when production maxvalue is entered.
func (s *BaseECLListener) EnterMaxvalue(ctx *MaxvalueContext) {}

// ExitMaxvalue is called when production maxvalue is exited.
func (s *BaseECLListener) ExitMaxvalue(ctx *MaxvalueContext) {}

// EnterMany is called when production many is entered.
func (s *BaseECLListener) EnterMany(ctx *ManyContext) {}

// ExitMany is called when production many is exited.
func (s *BaseECLListener) ExitMany(ctx *ManyContext) {}

// EnterReverseflag is called when production reverseflag is entered.
func (s *BaseECLListener) EnterReverseflag(ctx *ReverseflagContext) {}

// ExitReverseflag is called when production reverseflag is exited.
func (s *BaseECLListener) ExitReverseflag(ctx *ReverseflagContext) {}

// EnterEclattributename is called when production eclattributename is entered.
func (s *BaseECLListener) EnterEclattributename(ctx *EclattributenameContext) {}

// ExitEclattributename is called when production eclattributename is exited.
func (s *BaseECLListener) ExitEclattributename(ctx *EclattributenameContext) {}

// EnterExpressioncomparisonoperator is called when production expressioncomparisonoperator is entered.
func (s *BaseECLListener) EnterExpressioncomparisonoperator(ctx *ExpressioncomparisonoperatorContext) {
}

// ExitExpressioncomparisonoperator is called when production expressioncomparisonoperator is exited.
func (s *BaseECLListener) ExitExpressioncomparisonoperator(ctx *ExpressioncomparisonoperatorContext) {}

// EnterNumericcomparisonoperator is called when production numericcomparisonoperator is entered.
func (s *BaseECLListener) EnterNumericcomparisonoperator(ctx *NumericcomparisonoperatorContext) {}

// ExitNumericcomparisonoperator is called when production numericcomparisonoperator is exited.
func (s *BaseECLListener) ExitNumericcomparisonoperator(ctx *NumericcomparisonoperatorContext) {}

// EnterStringcomparisonoperator is called when production stringcomparisonoperator is entered.
func (s *BaseECLListener) EnterStringcomparisonoperator(ctx *StringcomparisonoperatorContext) {}

// ExitStringcomparisonoperator is called when production stringcomparisonoperator is exited.
func (s *BaseECLListener) ExitStringcomparisonoperator(ctx *StringcomparisonoperatorContext) {}

// EnterNumericvalue is called when production numericvalue is entered.
func (s *BaseECLListener) EnterNumericvalue(ctx *NumericvalueContext) {}

// ExitNumericvalue is called when production numericvalue is exited.
func (s *BaseECLListener) ExitNumericvalue(ctx *NumericvalueContext) {}

// EnterStringvalue is called when production stringvalue is entered.
func (s *BaseECLListener) EnterStringvalue(ctx *StringvalueContext) {}

// ExitStringvalue is called when production stringvalue is exited.
func (s *BaseECLListener) ExitStringvalue(ctx *StringvalueContext) {}

// EnterIntegervalue is called when production integervalue is entered.
func (s *BaseECLListener) EnterIntegervalue(ctx *IntegervalueContext) {}

// ExitIntegervalue is called when production integervalue is exited.
func (s *BaseECLListener) ExitIntegervalue(ctx *IntegervalueContext) {}

// EnterDecimalvalue is called when production decimalvalue is entered.
func (s *BaseECLListener) EnterDecimalvalue(ctx *DecimalvalueContext) {}

// ExitDecimalvalue is called when production decimalvalue is exited.
func (s *BaseECLListener) ExitDecimalvalue(ctx *DecimalvalueContext) {}

// EnterNonnegativeintegervalue is called when production nonnegativeintegervalue is entered.
func (s *BaseECLListener) EnterNonnegativeintegervalue(ctx *NonnegativeintegervalueContext) {}

// ExitNonnegativeintegervalue is called when production nonnegativeintegervalue is exited.
func (s *BaseECLListener) ExitNonnegativeintegervalue(ctx *NonnegativeintegervalueContext) {}

// EnterSctid is called when production sctid is entered.
func (s *BaseECLListener) EnterSctid(ctx *SctidContext) {}

// ExitSctid is called when production sctid is exited.
func (s *BaseECLListener) ExitSctid(ctx *SctidContext) {}

// EnterWs is called when production ws is entered.
func (s *BaseECLListener) EnterWs(ctx *WsContext) {}

// ExitWs is called when production ws is exited.
func (s *BaseECLListener) ExitWs(ctx *WsContext) {}

// EnterMws is called when production mws is entered.
func (s *BaseECLListener) EnterMws(ctx *MwsContext) {}

// ExitMws is called when production mws is exited.
func (s *BaseECLListener) ExitMws(ctx *MwsContext) {}

// EnterComment is called when production comment is entered.
func (s *BaseECLListener) EnterComment(ctx *CommentContext) {}

// ExitComment is called when production comment is exited.
func (s *BaseECLListener) ExitComment(ctx *CommentContext) {}

// EnterNonstarchar is called when production nonstarchar is entered.
func (s *BaseECLListener) EnterNonstarchar(ctx *NonstarcharContext) {}

// ExitNonstarchar is called when production nonstarchar is exited.
func (s *BaseECLListener) ExitNonstarchar(ctx *NonstarcharContext) {}

// EnterStarwithnonfslash is called when production starwithnonfslash is entered.
func (s *BaseECLListener) EnterStarwithnonfslash(ctx *StarwithnonfslashContext) {}

// ExitStarwithnonfslash is called when production starwithnonfslash is exited.
func (s *BaseECLListener) ExitStarwithnonfslash(ctx *StarwithnonfslashContext) {}

// EnterNonfslash is called when production nonfslash is entered.
func (s *BaseECLListener) EnterNonfslash(ctx *NonfslashContext) {}

// ExitNonfslash is called when production nonfslash is exited.
func (s *BaseECLListener) ExitNonfslash(ctx *NonfslashContext) {}

// EnterSp is called when production sp is entered.
func (s *BaseECLListener) EnterSp(ctx *SpContext) {}

// ExitSp is called when production sp is exited.
func (s *BaseECLListener) ExitSp(ctx *SpContext) {}

// EnterHtab is called when production htab is entered.
func (s *BaseECLListener) EnterHtab(ctx *HtabContext) {}

// ExitHtab is called when production htab is exited.
func (s *BaseECLListener) ExitHtab(ctx *HtabContext) {}

// EnterCr is called when production cr is entered.
func (s *BaseECLListener) EnterCr(ctx *CrContext) {}

// ExitCr is called when production cr is exited.
func (s *BaseECLListener) ExitCr(ctx *CrContext) {}

// EnterLf is called when production lf is entered.
func (s *BaseECLListener) EnterLf(ctx *LfContext) {}

// ExitLf is called when production lf is exited.
func (s *BaseECLListener) ExitLf(ctx *LfContext) {}

// EnterQm is called when production qm is entered.
func (s *BaseECLListener) EnterQm(ctx *QmContext) {}

// ExitQm is called when production qm is exited.
func (s *BaseECLListener) ExitQm(ctx *QmContext) {}

// EnterBs is called when production bs is entered.
func (s *BaseECLListener) EnterBs(ctx *BsContext) {}

// ExitBs is called when production bs is exited.
func (s *BaseECLListener) ExitBs(ctx *BsContext) {}

// EnterDigit is called when production digit is entered.
func (s *BaseECLListener) EnterDigit(ctx *DigitContext) {}

// ExitDigit is called when production digit is exited.
func (s *BaseECLListener) ExitDigit(ctx *DigitContext) {}

// EnterZero is called when production zero is entered.
func (s *BaseECLListener) EnterZero(ctx *ZeroContext) {}

// ExitZero is called when production zero is exited.
func (s *BaseECLListener) ExitZero(ctx *ZeroContext) {}

// EnterDigitnonzero is called when production digitnonzero is entered.
func (s *BaseECLListener) EnterDigitnonzero(ctx *DigitnonzeroContext) {}

// ExitDigitnonzero is called when production digitnonzero is exited.
func (s *BaseECLListener) ExitDigitnonzero(ctx *DigitnonzeroContext) {}

// EnterNonwsnonpipe is called when production nonwsnonpipe is entered.
func (s *BaseECLListener) EnterNonwsnonpipe(ctx *NonwsnonpipeContext) {}

// ExitNonwsnonpipe is called when production nonwsnonpipe is exited.
func (s *BaseECLListener) ExitNonwsnonpipe(ctx *NonwsnonpipeContext) {}

// EnterAnynonescapedchar is called when production anynonescapedchar is entered.
func (s *BaseECLListener) EnterAnynonescapedchar(ctx *AnynonescapedcharContext) {}

// ExitAnynonescapedchar is called when production anynonescapedchar is exited.
func (s *BaseECLListener) ExitAnynonescapedchar(ctx *AnynonescapedcharContext) {}

// EnterEscapedchar is called when production escapedchar is entered.
func (s *BaseECLListener) EnterEscapedchar(ctx *EscapedcharContext) {}

// ExitEscapedchar is called when production escapedchar is exited.
func (s *BaseECLListener) ExitEscapedchar(ctx *EscapedcharContext) {}

// EnterUtf8_2 is called when production utf8_2 is entered.
func (s *BaseECLListener) EnterUtf8_2(ctx *Utf8_2Context) {}

// ExitUtf8_2 is called when production utf8_2 is exited.
func (s *BaseECLListener) ExitUtf8_2(ctx *Utf8_2Context) {}

// EnterUtf8_3 is called when production utf8_3 is entered.
func (s *BaseECLListener) EnterUtf8_3(ctx *Utf8_3Context) {}

// ExitUtf8_3 is called when production utf8_3 is exited.
func (s *BaseECLListener) ExitUtf8_3(ctx *Utf8_3Context) {}

// EnterUtf8_4 is called when production utf8_4 is entered.
func (s *BaseECLListener) EnterUtf8_4(ctx *Utf8_4Context) {}

// ExitUtf8_4 is called when production utf8_4 is exited.
func (s *BaseECLListener) ExitUtf8_4(ctx *Utf8_4Context) {}

// EnterUtf8_tail is called when production utf8_tail is entered.
func (s *BaseECLListener) EnterUtf8_tail(ctx *Utf8_tailContext) {}

// ExitUtf8_tail is called when production utf8_tail is exited.
func (s *BaseECLListener) ExitUtf8_tail(ctx *Utf8_tailContext) {}
