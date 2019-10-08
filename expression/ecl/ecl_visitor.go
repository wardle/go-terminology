// Code generated from ECL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package ecl // ECL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by ECLParser.
type ECLVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ECLParser#expressionconstraint.
	VisitExpressionconstraint(ctx *ExpressionconstraintContext) interface{}

	// Visit a parse tree produced by ECLParser#refinedexpressionconstraint.
	VisitRefinedexpressionconstraint(ctx *RefinedexpressionconstraintContext) interface{}

	// Visit a parse tree produced by ECLParser#compoundexpressionconstraint.
	VisitCompoundexpressionconstraint(ctx *CompoundexpressionconstraintContext) interface{}

	// Visit a parse tree produced by ECLParser#conjunctionexpressionconstraint.
	VisitConjunctionexpressionconstraint(ctx *ConjunctionexpressionconstraintContext) interface{}

	// Visit a parse tree produced by ECLParser#disjunctionexpressionconstraint.
	VisitDisjunctionexpressionconstraint(ctx *DisjunctionexpressionconstraintContext) interface{}

	// Visit a parse tree produced by ECLParser#exclusionexpressionconstraint.
	VisitExclusionexpressionconstraint(ctx *ExclusionexpressionconstraintContext) interface{}

	// Visit a parse tree produced by ECLParser#dottedexpressionconstraint.
	VisitDottedexpressionconstraint(ctx *DottedexpressionconstraintContext) interface{}

	// Visit a parse tree produced by ECLParser#dottedexpressionattribute.
	VisitDottedexpressionattribute(ctx *DottedexpressionattributeContext) interface{}

	// Visit a parse tree produced by ECLParser#subexpressionconstraint.
	VisitSubexpressionconstraint(ctx *SubexpressionconstraintContext) interface{}

	// Visit a parse tree produced by ECLParser#eclfocusconcept.
	VisitEclfocusconcept(ctx *EclfocusconceptContext) interface{}

	// Visit a parse tree produced by ECLParser#dot.
	VisitDot(ctx *DotContext) interface{}

	// Visit a parse tree produced by ECLParser#memberof.
	VisitMemberof(ctx *MemberofContext) interface{}

	// Visit a parse tree produced by ECLParser#eclconceptreference.
	VisitEclconceptreference(ctx *EclconceptreferenceContext) interface{}

	// Visit a parse tree produced by ECLParser#conceptid.
	VisitConceptid(ctx *ConceptidContext) interface{}

	// Visit a parse tree produced by ECLParser#term.
	VisitTerm(ctx *TermContext) interface{}

	// Visit a parse tree produced by ECLParser#wildcard.
	VisitWildcard(ctx *WildcardContext) interface{}

	// Visit a parse tree produced by ECLParser#constraintoperator.
	VisitConstraintoperator(ctx *ConstraintoperatorContext) interface{}

	// Visit a parse tree produced by ECLParser#descendantof.
	VisitDescendantof(ctx *DescendantofContext) interface{}

	// Visit a parse tree produced by ECLParser#descendantorselfof.
	VisitDescendantorselfof(ctx *DescendantorselfofContext) interface{}

	// Visit a parse tree produced by ECLParser#childof.
	VisitChildof(ctx *ChildofContext) interface{}

	// Visit a parse tree produced by ECLParser#ancestorof.
	VisitAncestorof(ctx *AncestorofContext) interface{}

	// Visit a parse tree produced by ECLParser#ancestororselfof.
	VisitAncestororselfof(ctx *AncestororselfofContext) interface{}

	// Visit a parse tree produced by ECLParser#parentof.
	VisitParentof(ctx *ParentofContext) interface{}

	// Visit a parse tree produced by ECLParser#conjunction.
	VisitConjunction(ctx *ConjunctionContext) interface{}

	// Visit a parse tree produced by ECLParser#disjunction.
	VisitDisjunction(ctx *DisjunctionContext) interface{}

	// Visit a parse tree produced by ECLParser#exclusion.
	VisitExclusion(ctx *ExclusionContext) interface{}

	// Visit a parse tree produced by ECLParser#eclrefinement.
	VisitEclrefinement(ctx *EclrefinementContext) interface{}

	// Visit a parse tree produced by ECLParser#conjunctionrefinementset.
	VisitConjunctionrefinementset(ctx *ConjunctionrefinementsetContext) interface{}

	// Visit a parse tree produced by ECLParser#disjunctionrefinementset.
	VisitDisjunctionrefinementset(ctx *DisjunctionrefinementsetContext) interface{}

	// Visit a parse tree produced by ECLParser#subrefinement.
	VisitSubrefinement(ctx *SubrefinementContext) interface{}

	// Visit a parse tree produced by ECLParser#eclattributeset.
	VisitEclattributeset(ctx *EclattributesetContext) interface{}

	// Visit a parse tree produced by ECLParser#conjunctionattributeset.
	VisitConjunctionattributeset(ctx *ConjunctionattributesetContext) interface{}

	// Visit a parse tree produced by ECLParser#disjunctionattributeset.
	VisitDisjunctionattributeset(ctx *DisjunctionattributesetContext) interface{}

	// Visit a parse tree produced by ECLParser#subattributeset.
	VisitSubattributeset(ctx *SubattributesetContext) interface{}

	// Visit a parse tree produced by ECLParser#eclattributegroup.
	VisitEclattributegroup(ctx *EclattributegroupContext) interface{}

	// Visit a parse tree produced by ECLParser#eclattribute.
	VisitEclattribute(ctx *EclattributeContext) interface{}

	// Visit a parse tree produced by ECLParser#cardinality.
	VisitCardinality(ctx *CardinalityContext) interface{}

	// Visit a parse tree produced by ECLParser#minvalue.
	VisitMinvalue(ctx *MinvalueContext) interface{}

	// Visit a parse tree produced by ECLParser#to.
	VisitTo(ctx *ToContext) interface{}

	// Visit a parse tree produced by ECLParser#maxvalue.
	VisitMaxvalue(ctx *MaxvalueContext) interface{}

	// Visit a parse tree produced by ECLParser#many.
	VisitMany(ctx *ManyContext) interface{}

	// Visit a parse tree produced by ECLParser#reverseflag.
	VisitReverseflag(ctx *ReverseflagContext) interface{}

	// Visit a parse tree produced by ECLParser#eclattributename.
	VisitEclattributename(ctx *EclattributenameContext) interface{}

	// Visit a parse tree produced by ECLParser#expressioncomparisonoperator.
	VisitExpressioncomparisonoperator(ctx *ExpressioncomparisonoperatorContext) interface{}

	// Visit a parse tree produced by ECLParser#numericcomparisonoperator.
	VisitNumericcomparisonoperator(ctx *NumericcomparisonoperatorContext) interface{}

	// Visit a parse tree produced by ECLParser#stringcomparisonoperator.
	VisitStringcomparisonoperator(ctx *StringcomparisonoperatorContext) interface{}

	// Visit a parse tree produced by ECLParser#numericvalue.
	VisitNumericvalue(ctx *NumericvalueContext) interface{}

	// Visit a parse tree produced by ECLParser#stringvalue.
	VisitStringvalue(ctx *StringvalueContext) interface{}

	// Visit a parse tree produced by ECLParser#integervalue.
	VisitIntegervalue(ctx *IntegervalueContext) interface{}

	// Visit a parse tree produced by ECLParser#decimalvalue.
	VisitDecimalvalue(ctx *DecimalvalueContext) interface{}

	// Visit a parse tree produced by ECLParser#nonnegativeintegervalue.
	VisitNonnegativeintegervalue(ctx *NonnegativeintegervalueContext) interface{}

	// Visit a parse tree produced by ECLParser#sctid.
	VisitSctid(ctx *SctidContext) interface{}

	// Visit a parse tree produced by ECLParser#ws.
	VisitWs(ctx *WsContext) interface{}

	// Visit a parse tree produced by ECLParser#mws.
	VisitMws(ctx *MwsContext) interface{}

	// Visit a parse tree produced by ECLParser#comment.
	VisitComment(ctx *CommentContext) interface{}

	// Visit a parse tree produced by ECLParser#nonstarchar.
	VisitNonstarchar(ctx *NonstarcharContext) interface{}

	// Visit a parse tree produced by ECLParser#starwithnonfslash.
	VisitStarwithnonfslash(ctx *StarwithnonfslashContext) interface{}

	// Visit a parse tree produced by ECLParser#nonfslash.
	VisitNonfslash(ctx *NonfslashContext) interface{}

	// Visit a parse tree produced by ECLParser#sp.
	VisitSp(ctx *SpContext) interface{}

	// Visit a parse tree produced by ECLParser#htab.
	VisitHtab(ctx *HtabContext) interface{}

	// Visit a parse tree produced by ECLParser#cr.
	VisitCr(ctx *CrContext) interface{}

	// Visit a parse tree produced by ECLParser#lf.
	VisitLf(ctx *LfContext) interface{}

	// Visit a parse tree produced by ECLParser#qm.
	VisitQm(ctx *QmContext) interface{}

	// Visit a parse tree produced by ECLParser#bs.
	VisitBs(ctx *BsContext) interface{}

	// Visit a parse tree produced by ECLParser#digit.
	VisitDigit(ctx *DigitContext) interface{}

	// Visit a parse tree produced by ECLParser#zero.
	VisitZero(ctx *ZeroContext) interface{}

	// Visit a parse tree produced by ECLParser#digitnonzero.
	VisitDigitnonzero(ctx *DigitnonzeroContext) interface{}

	// Visit a parse tree produced by ECLParser#nonwsnonpipe.
	VisitNonwsnonpipe(ctx *NonwsnonpipeContext) interface{}

	// Visit a parse tree produced by ECLParser#anynonescapedchar.
	VisitAnynonescapedchar(ctx *AnynonescapedcharContext) interface{}

	// Visit a parse tree produced by ECLParser#escapedchar.
	VisitEscapedchar(ctx *EscapedcharContext) interface{}

	// Visit a parse tree produced by ECLParser#utf8_2.
	VisitUtf8_2(ctx *Utf8_2Context) interface{}

	// Visit a parse tree produced by ECLParser#utf8_3.
	VisitUtf8_3(ctx *Utf8_3Context) interface{}

	// Visit a parse tree produced by ECLParser#utf8_4.
	VisitUtf8_4(ctx *Utf8_4Context) interface{}

	// Visit a parse tree produced by ECLParser#utf8_tail.
	VisitUtf8_tail(ctx *Utf8_tailContext) interface{}
}
