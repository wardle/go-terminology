// Generated from ECL.g4 by ANTLR 4.7.

package ecl // ECL
import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseECLVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseECLVisitor) VisitExpressionconstraint(ctx *ExpressionconstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitRefinedexpressionconstraint(ctx *RefinedexpressionconstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitCompoundexpressionconstraint(ctx *CompoundexpressionconstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitConjunctionexpressionconstraint(ctx *ConjunctionexpressionconstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDisjunctionexpressionconstraint(ctx *DisjunctionexpressionconstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitExclusionexpressionconstraint(ctx *ExclusionexpressionconstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDottedexpressionconstraint(ctx *DottedexpressionconstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDottedexpressionattribute(ctx *DottedexpressionattributeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitSubexpressionconstraint(ctx *SubexpressionconstraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitEclfocusconcept(ctx *EclfocusconceptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDot(ctx *DotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitMemberof(ctx *MemberofContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitEclconceptreference(ctx *EclconceptreferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitConceptid(ctx *ConceptidContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitTerm(ctx *TermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitWildcard(ctx *WildcardContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitConstraintoperator(ctx *ConstraintoperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDescendantof(ctx *DescendantofContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDescendantorselfof(ctx *DescendantorselfofContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitChildof(ctx *ChildofContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitAncestorof(ctx *AncestorofContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitAncestororselfof(ctx *AncestororselfofContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitParentof(ctx *ParentofContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitConjunction(ctx *ConjunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDisjunction(ctx *DisjunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitExclusion(ctx *ExclusionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitEclrefinement(ctx *EclrefinementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitConjunctionrefinementset(ctx *ConjunctionrefinementsetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDisjunctionrefinementset(ctx *DisjunctionrefinementsetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitSubrefinement(ctx *SubrefinementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitEclattributeset(ctx *EclattributesetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitConjunctionattributeset(ctx *ConjunctionattributesetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDisjunctionattributeset(ctx *DisjunctionattributesetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitSubattributeset(ctx *SubattributesetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitEclattributegroup(ctx *EclattributegroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitEclattribute(ctx *EclattributeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitCardinality(ctx *CardinalityContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitMinvalue(ctx *MinvalueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitTo(ctx *ToContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitMaxvalue(ctx *MaxvalueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitMany(ctx *ManyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitReverseflag(ctx *ReverseflagContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitEclattributename(ctx *EclattributenameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitExpressioncomparisonoperator(ctx *ExpressioncomparisonoperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitNumericcomparisonoperator(ctx *NumericcomparisonoperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitStringcomparisonoperator(ctx *StringcomparisonoperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitNumericvalue(ctx *NumericvalueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitStringvalue(ctx *StringvalueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitIntegervalue(ctx *IntegervalueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDecimalvalue(ctx *DecimalvalueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitNonnegativeintegervalue(ctx *NonnegativeintegervalueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitSctid(ctx *SctidContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitWs(ctx *WsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitMws(ctx *MwsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitComment(ctx *CommentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitNonstarchar(ctx *NonstarcharContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitStarwithnonfslash(ctx *StarwithnonfslashContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitNonfslash(ctx *NonfslashContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitSp(ctx *SpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitHtab(ctx *HtabContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitCr(ctx *CrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitLf(ctx *LfContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitQm(ctx *QmContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitBs(ctx *BsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDigit(ctx *DigitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitZero(ctx *ZeroContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitDigitnonzero(ctx *DigitnonzeroContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitNonwsnonpipe(ctx *NonwsnonpipeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitAnynonescapedchar(ctx *AnynonescapedcharContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitEscapedchar(ctx *EscapedcharContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitUtf8_2(ctx *Utf8_2Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitUtf8_3(ctx *Utf8_3Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitUtf8_4(ctx *Utf8_4Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseECLVisitor) VisitUtf8_tail(ctx *Utf8_tailContext) interface{} {
	return v.VisitChildren(ctx)
}
