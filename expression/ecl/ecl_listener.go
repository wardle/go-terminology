// Generated from ECL.g4 by ANTLR 4.7.

package ecl // ECL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// ECLListener is a complete listener for a parse tree produced by ECLParser.
type ECLListener interface {
	antlr.ParseTreeListener

	// EnterExpressionconstraint is called when entering the expressionconstraint production.
	EnterExpressionconstraint(c *ExpressionconstraintContext)

	// EnterRefinedexpressionconstraint is called when entering the refinedexpressionconstraint production.
	EnterRefinedexpressionconstraint(c *RefinedexpressionconstraintContext)

	// EnterCompoundexpressionconstraint is called when entering the compoundexpressionconstraint production.
	EnterCompoundexpressionconstraint(c *CompoundexpressionconstraintContext)

	// EnterConjunctionexpressionconstraint is called when entering the conjunctionexpressionconstraint production.
	EnterConjunctionexpressionconstraint(c *ConjunctionexpressionconstraintContext)

	// EnterDisjunctionexpressionconstraint is called when entering the disjunctionexpressionconstraint production.
	EnterDisjunctionexpressionconstraint(c *DisjunctionexpressionconstraintContext)

	// EnterExclusionexpressionconstraint is called when entering the exclusionexpressionconstraint production.
	EnterExclusionexpressionconstraint(c *ExclusionexpressionconstraintContext)

	// EnterDottedexpressionconstraint is called when entering the dottedexpressionconstraint production.
	EnterDottedexpressionconstraint(c *DottedexpressionconstraintContext)

	// EnterDottedexpressionattribute is called when entering the dottedexpressionattribute production.
	EnterDottedexpressionattribute(c *DottedexpressionattributeContext)

	// EnterSubexpressionconstraint is called when entering the subexpressionconstraint production.
	EnterSubexpressionconstraint(c *SubexpressionconstraintContext)

	// EnterEclfocusconcept is called when entering the eclfocusconcept production.
	EnterEclfocusconcept(c *EclfocusconceptContext)

	// EnterDot is called when entering the dot production.
	EnterDot(c *DotContext)

	// EnterMemberof is called when entering the memberof production.
	EnterMemberof(c *MemberofContext)

	// EnterEclconceptreference is called when entering the eclconceptreference production.
	EnterEclconceptreference(c *EclconceptreferenceContext)

	// EnterConceptid is called when entering the conceptid production.
	EnterConceptid(c *ConceptidContext)

	// EnterTerm is called when entering the term production.
	EnterTerm(c *TermContext)

	// EnterWildcard is called when entering the wildcard production.
	EnterWildcard(c *WildcardContext)

	// EnterConstraintoperator is called when entering the constraintoperator production.
	EnterConstraintoperator(c *ConstraintoperatorContext)

	// EnterDescendantof is called when entering the descendantof production.
	EnterDescendantof(c *DescendantofContext)

	// EnterDescendantorselfof is called when entering the descendantorselfof production.
	EnterDescendantorselfof(c *DescendantorselfofContext)

	// EnterChildof is called when entering the childof production.
	EnterChildof(c *ChildofContext)

	// EnterAncestorof is called when entering the ancestorof production.
	EnterAncestorof(c *AncestorofContext)

	// EnterAncestororselfof is called when entering the ancestororselfof production.
	EnterAncestororselfof(c *AncestororselfofContext)

	// EnterParentof is called when entering the parentof production.
	EnterParentof(c *ParentofContext)

	// EnterConjunction is called when entering the conjunction production.
	EnterConjunction(c *ConjunctionContext)

	// EnterDisjunction is called when entering the disjunction production.
	EnterDisjunction(c *DisjunctionContext)

	// EnterExclusion is called when entering the exclusion production.
	EnterExclusion(c *ExclusionContext)

	// EnterEclrefinement is called when entering the eclrefinement production.
	EnterEclrefinement(c *EclrefinementContext)

	// EnterConjunctionrefinementset is called when entering the conjunctionrefinementset production.
	EnterConjunctionrefinementset(c *ConjunctionrefinementsetContext)

	// EnterDisjunctionrefinementset is called when entering the disjunctionrefinementset production.
	EnterDisjunctionrefinementset(c *DisjunctionrefinementsetContext)

	// EnterSubrefinement is called when entering the subrefinement production.
	EnterSubrefinement(c *SubrefinementContext)

	// EnterEclattributeset is called when entering the eclattributeset production.
	EnterEclattributeset(c *EclattributesetContext)

	// EnterConjunctionattributeset is called when entering the conjunctionattributeset production.
	EnterConjunctionattributeset(c *ConjunctionattributesetContext)

	// EnterDisjunctionattributeset is called when entering the disjunctionattributeset production.
	EnterDisjunctionattributeset(c *DisjunctionattributesetContext)

	// EnterSubattributeset is called when entering the subattributeset production.
	EnterSubattributeset(c *SubattributesetContext)

	// EnterEclattributegroup is called when entering the eclattributegroup production.
	EnterEclattributegroup(c *EclattributegroupContext)

	// EnterEclattribute is called when entering the eclattribute production.
	EnterEclattribute(c *EclattributeContext)

	// EnterCardinality is called when entering the cardinality production.
	EnterCardinality(c *CardinalityContext)

	// EnterMinvalue is called when entering the minvalue production.
	EnterMinvalue(c *MinvalueContext)

	// EnterTo is called when entering the to production.
	EnterTo(c *ToContext)

	// EnterMaxvalue is called when entering the maxvalue production.
	EnterMaxvalue(c *MaxvalueContext)

	// EnterMany is called when entering the many production.
	EnterMany(c *ManyContext)

	// EnterReverseflag is called when entering the reverseflag production.
	EnterReverseflag(c *ReverseflagContext)

	// EnterEclattributename is called when entering the eclattributename production.
	EnterEclattributename(c *EclattributenameContext)

	// EnterExpressioncomparisonoperator is called when entering the expressioncomparisonoperator production.
	EnterExpressioncomparisonoperator(c *ExpressioncomparisonoperatorContext)

	// EnterNumericcomparisonoperator is called when entering the numericcomparisonoperator production.
	EnterNumericcomparisonoperator(c *NumericcomparisonoperatorContext)

	// EnterStringcomparisonoperator is called when entering the stringcomparisonoperator production.
	EnterStringcomparisonoperator(c *StringcomparisonoperatorContext)

	// EnterNumericvalue is called when entering the numericvalue production.
	EnterNumericvalue(c *NumericvalueContext)

	// EnterStringvalue is called when entering the stringvalue production.
	EnterStringvalue(c *StringvalueContext)

	// EnterIntegervalue is called when entering the integervalue production.
	EnterIntegervalue(c *IntegervalueContext)

	// EnterDecimalvalue is called when entering the decimalvalue production.
	EnterDecimalvalue(c *DecimalvalueContext)

	// EnterNonnegativeintegervalue is called when entering the nonnegativeintegervalue production.
	EnterNonnegativeintegervalue(c *NonnegativeintegervalueContext)

	// EnterSctid is called when entering the sctid production.
	EnterSctid(c *SctidContext)

	// EnterWs is called when entering the ws production.
	EnterWs(c *WsContext)

	// EnterMws is called when entering the mws production.
	EnterMws(c *MwsContext)

	// EnterComment is called when entering the comment production.
	EnterComment(c *CommentContext)

	// EnterNonstarchar is called when entering the nonstarchar production.
	EnterNonstarchar(c *NonstarcharContext)

	// EnterStarwithnonfslash is called when entering the starwithnonfslash production.
	EnterStarwithnonfslash(c *StarwithnonfslashContext)

	// EnterNonfslash is called when entering the nonfslash production.
	EnterNonfslash(c *NonfslashContext)

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

	// ExitExpressionconstraint is called when exiting the expressionconstraint production.
	ExitExpressionconstraint(c *ExpressionconstraintContext)

	// ExitRefinedexpressionconstraint is called when exiting the refinedexpressionconstraint production.
	ExitRefinedexpressionconstraint(c *RefinedexpressionconstraintContext)

	// ExitCompoundexpressionconstraint is called when exiting the compoundexpressionconstraint production.
	ExitCompoundexpressionconstraint(c *CompoundexpressionconstraintContext)

	// ExitConjunctionexpressionconstraint is called when exiting the conjunctionexpressionconstraint production.
	ExitConjunctionexpressionconstraint(c *ConjunctionexpressionconstraintContext)

	// ExitDisjunctionexpressionconstraint is called when exiting the disjunctionexpressionconstraint production.
	ExitDisjunctionexpressionconstraint(c *DisjunctionexpressionconstraintContext)

	// ExitExclusionexpressionconstraint is called when exiting the exclusionexpressionconstraint production.
	ExitExclusionexpressionconstraint(c *ExclusionexpressionconstraintContext)

	// ExitDottedexpressionconstraint is called when exiting the dottedexpressionconstraint production.
	ExitDottedexpressionconstraint(c *DottedexpressionconstraintContext)

	// ExitDottedexpressionattribute is called when exiting the dottedexpressionattribute production.
	ExitDottedexpressionattribute(c *DottedexpressionattributeContext)

	// ExitSubexpressionconstraint is called when exiting the subexpressionconstraint production.
	ExitSubexpressionconstraint(c *SubexpressionconstraintContext)

	// ExitEclfocusconcept is called when exiting the eclfocusconcept production.
	ExitEclfocusconcept(c *EclfocusconceptContext)

	// ExitDot is called when exiting the dot production.
	ExitDot(c *DotContext)

	// ExitMemberof is called when exiting the memberof production.
	ExitMemberof(c *MemberofContext)

	// ExitEclconceptreference is called when exiting the eclconceptreference production.
	ExitEclconceptreference(c *EclconceptreferenceContext)

	// ExitConceptid is called when exiting the conceptid production.
	ExitConceptid(c *ConceptidContext)

	// ExitTerm is called when exiting the term production.
	ExitTerm(c *TermContext)

	// ExitWildcard is called when exiting the wildcard production.
	ExitWildcard(c *WildcardContext)

	// ExitConstraintoperator is called when exiting the constraintoperator production.
	ExitConstraintoperator(c *ConstraintoperatorContext)

	// ExitDescendantof is called when exiting the descendantof production.
	ExitDescendantof(c *DescendantofContext)

	// ExitDescendantorselfof is called when exiting the descendantorselfof production.
	ExitDescendantorselfof(c *DescendantorselfofContext)

	// ExitChildof is called when exiting the childof production.
	ExitChildof(c *ChildofContext)

	// ExitAncestorof is called when exiting the ancestorof production.
	ExitAncestorof(c *AncestorofContext)

	// ExitAncestororselfof is called when exiting the ancestororselfof production.
	ExitAncestororselfof(c *AncestororselfofContext)

	// ExitParentof is called when exiting the parentof production.
	ExitParentof(c *ParentofContext)

	// ExitConjunction is called when exiting the conjunction production.
	ExitConjunction(c *ConjunctionContext)

	// ExitDisjunction is called when exiting the disjunction production.
	ExitDisjunction(c *DisjunctionContext)

	// ExitExclusion is called when exiting the exclusion production.
	ExitExclusion(c *ExclusionContext)

	// ExitEclrefinement is called when exiting the eclrefinement production.
	ExitEclrefinement(c *EclrefinementContext)

	// ExitConjunctionrefinementset is called when exiting the conjunctionrefinementset production.
	ExitConjunctionrefinementset(c *ConjunctionrefinementsetContext)

	// ExitDisjunctionrefinementset is called when exiting the disjunctionrefinementset production.
	ExitDisjunctionrefinementset(c *DisjunctionrefinementsetContext)

	// ExitSubrefinement is called when exiting the subrefinement production.
	ExitSubrefinement(c *SubrefinementContext)

	// ExitEclattributeset is called when exiting the eclattributeset production.
	ExitEclattributeset(c *EclattributesetContext)

	// ExitConjunctionattributeset is called when exiting the conjunctionattributeset production.
	ExitConjunctionattributeset(c *ConjunctionattributesetContext)

	// ExitDisjunctionattributeset is called when exiting the disjunctionattributeset production.
	ExitDisjunctionattributeset(c *DisjunctionattributesetContext)

	// ExitSubattributeset is called when exiting the subattributeset production.
	ExitSubattributeset(c *SubattributesetContext)

	// ExitEclattributegroup is called when exiting the eclattributegroup production.
	ExitEclattributegroup(c *EclattributegroupContext)

	// ExitEclattribute is called when exiting the eclattribute production.
	ExitEclattribute(c *EclattributeContext)

	// ExitCardinality is called when exiting the cardinality production.
	ExitCardinality(c *CardinalityContext)

	// ExitMinvalue is called when exiting the minvalue production.
	ExitMinvalue(c *MinvalueContext)

	// ExitTo is called when exiting the to production.
	ExitTo(c *ToContext)

	// ExitMaxvalue is called when exiting the maxvalue production.
	ExitMaxvalue(c *MaxvalueContext)

	// ExitMany is called when exiting the many production.
	ExitMany(c *ManyContext)

	// ExitReverseflag is called when exiting the reverseflag production.
	ExitReverseflag(c *ReverseflagContext)

	// ExitEclattributename is called when exiting the eclattributename production.
	ExitEclattributename(c *EclattributenameContext)

	// ExitExpressioncomparisonoperator is called when exiting the expressioncomparisonoperator production.
	ExitExpressioncomparisonoperator(c *ExpressioncomparisonoperatorContext)

	// ExitNumericcomparisonoperator is called when exiting the numericcomparisonoperator production.
	ExitNumericcomparisonoperator(c *NumericcomparisonoperatorContext)

	// ExitStringcomparisonoperator is called when exiting the stringcomparisonoperator production.
	ExitStringcomparisonoperator(c *StringcomparisonoperatorContext)

	// ExitNumericvalue is called when exiting the numericvalue production.
	ExitNumericvalue(c *NumericvalueContext)

	// ExitStringvalue is called when exiting the stringvalue production.
	ExitStringvalue(c *StringvalueContext)

	// ExitIntegervalue is called when exiting the integervalue production.
	ExitIntegervalue(c *IntegervalueContext)

	// ExitDecimalvalue is called when exiting the decimalvalue production.
	ExitDecimalvalue(c *DecimalvalueContext)

	// ExitNonnegativeintegervalue is called when exiting the nonnegativeintegervalue production.
	ExitNonnegativeintegervalue(c *NonnegativeintegervalueContext)

	// ExitSctid is called when exiting the sctid production.
	ExitSctid(c *SctidContext)

	// ExitWs is called when exiting the ws production.
	ExitWs(c *WsContext)

	// ExitMws is called when exiting the mws production.
	ExitMws(c *MwsContext)

	// ExitComment is called when exiting the comment production.
	ExitComment(c *CommentContext)

	// ExitNonstarchar is called when exiting the nonstarchar production.
	ExitNonstarchar(c *NonstarcharContext)

	// ExitStarwithnonfslash is called when exiting the starwithnonfslash production.
	ExitStarwithnonfslash(c *StarwithnonfslashContext)

	// ExitNonfslash is called when exiting the nonfslash production.
	ExitNonfslash(c *NonfslashContext)

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
