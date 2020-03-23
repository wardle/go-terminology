package expression

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/wardle/go-terminology/expression/ecl"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
)

type cardinality struct {
	minimumValue int64
	maximumValue int64
	toMany       bool
}

type constraintOperator int

const (
	noConstraint constraintOperator = iota
	ancestorOf
	ancestorOrSelf
	childOf
	descendantOf
	descendantOrSelf
	parentOf
)

type comparisonOperator int

const (
	noComparison comparisonOperator = iota
	equals
	notEquals
	lessThanOrEquals
	lessThan
	greaterThanOrEquals
	greaterThan
)

type focusConcept struct {
	concept  *snomed.ConceptReference
	wildcard bool
}

// ApplyConstraint applies an expression in the "Expression Constraint Language" (ECL) to the specified (CG) expression
func ApplyConstraint(svc *terminology.Svc, exp *snomed.Expression, s string) (bool, error) {
	is := antlr.NewInputStream(s)
	lex := ecl.NewECLLexer(is)
	tokens := antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
	p := ecl.NewECLParser(tokens)
	p.RemoveErrorListeners() // remove default listeners, which includes console listener - so as to avoid printing all parse errors to console
	el := new(errorListener)
	p.AddErrorListener(el)
	visitor := &applyingECLVisitor{svc: svc, exp: exp}
	tree := p.Expressionconstraint()
	result, ok := visitor.Visit(tree).(bool)
	if !ok {
		return false, fmt.Errorf("parse error. got: %v", result)
	}
	if el.err != nil {
		return false, el.err
	}
	if len(visitor.errors) > 0 {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d error(s) parsing expression: ", len(visitor.errors)))
		for _, e := range visitor.errors {
			sb.WriteString(fmt.Sprintf("%s ", e))
		}
		return false, errors.New(strings.TrimSpace(sb.String()))
	}
	return result, nil
}

// applyingECLVisitor is used to apply the Expression Constraint Language (ecl).
type applyingECLVisitor struct {
	ecl.BaseECLVisitor
	svc    *terminology.Svc
	exp    *snomed.Expression
	errors []error
}

func (ev *applyingECLVisitor) addError(err error) {
	ev.errors = append(ev.errors, err)
}

func (ev *applyingECLVisitor) Visit(tree antlr.ParseTree) interface{} {
	if tree == nil {
		return nil
	}
	return tree.Accept(ev)
}

func (ev *applyingECLVisitor) VisitExpressionconstraint(ctx *ecl.ExpressionconstraintContext) interface{} {
	if ctx.Compoundexpressionconstraint() != nil {
		if result, ok := ev.Visit(ctx.Compoundexpressionconstraint()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing compound expression: %s", ctx.Compoundexpressionconstraint().GetText()))
	}
	if dotted := ctx.Dottedexpressionconstraint(); dotted != nil {
		if result, ok := ev.Visit(ctx.Dottedexpressionconstraint()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing dotted expression: %s", ctx.Dottedexpressionconstraint().GetText()))
	}
	if ctx.Subexpressionconstraint() != nil {
		if result, ok := ev.Visit(ctx.Subexpressionconstraint()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing subexpression: %s", ctx.Subexpressionconstraint().GetText()))
	}
	if ctx.Refinedexpressionconstraint() != nil {
		if result, ok := ev.Visit(ctx.Refinedexpressionconstraint()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing refined expression: %s", ctx.Refinedexpressionconstraint().GetText()))
	}
	return false
}

func (ev *applyingECLVisitor) VisitRefinedexpressionconstraint(ctx *ecl.RefinedexpressionconstraintContext) interface{} {
	if ctx.Subexpressionconstraint() != nil {
		if result, ok := ev.Visit(ctx.Subexpressionconstraint()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing subexpression: %s", ctx.Subexpressionconstraint().GetText()))
	}
	if ctx.Eclrefinement() != nil {
		if result, ok := ev.Visit(ctx.Eclrefinement()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing refinement: %s", ctx.Eclrefinement().GetText()))
	}
	return false
}

// compoundExpressionConstraint = conjunctionExpressionConstraint / disjunctionExpressionConstraint / exclusionExpressionConstraint
func (ev *applyingECLVisitor) VisitCompoundexpressionconstraint(ctx *ecl.CompoundexpressionconstraintContext) interface{} {
	if ctx.Conjunctionexpressionconstraint() != nil {
		if result, ok := ev.Visit(ctx.Conjunctionexpressionconstraint()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing conjunction expression: %s", ctx.Conjunctionexpressionconstraint().GetText()))
	}
	if ctx.Disjunctionexpressionconstraint() != nil {
		if result, ok := ev.Visit(ctx.Disjunctionexpressionconstraint()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing disjunction expression: %s", ctx.Disjunctionexpressionconstraint().GetText()))
	}
	if ctx.Exclusionexpressionconstraint() != nil {
		if result, ok := ev.Visit(ctx.Exclusionexpressionconstraint()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing exclusion expression: %s", ctx.Exclusionexpressionconstraint().GetText()))
	}
	return false
}

// conjunctionExpressionConstraint = subExpressionConstraint 1*(ws conjunction ws subExpressionConstraint)
func (ev *applyingECLVisitor) VisitConjunctionexpressionconstraint(ctx *ecl.ConjunctionexpressionconstraintContext) interface{} {
	result := true
	for _, subexp := range ctx.AllSubexpressionconstraint() {
		if term, ok := ev.Visit(subexp).(bool); ok {
			result = result && term
		} else {
			ev.addError(fmt.Errorf("invalid subexpression in conjunction: %s", subexp.GetText()))
			return false
		}
	}
	return result
}

func (ev *applyingECLVisitor) VisitDisjunctionexpressionconstraint(ctx *ecl.DisjunctionexpressionconstraintContext) interface{} {
	result := false
	for _, subexp := range ctx.AllSubexpressionconstraint() {
		if term, ok := ev.Visit(subexp).(bool); ok {
			result = result || term
		} else {
			ev.addError(fmt.Errorf("error parsing disjunction subexpression: %s", subexp.GetText()))
			return false
		}
	}
	return result
}

func (ev *applyingECLVisitor) VisitExclusionexpressionconstraint(ctx *ecl.ExclusionexpressionconstraintContext) interface{} {
	result := false
	for _, subexp := range ctx.AllSubexpressionconstraint() {
		if term, ok := ev.Visit(subexp).(bool); ok {
			result = result || term
		} else {
			ev.addError(fmt.Errorf("errro parsing exclusion subexpression: %s", subexp.GetText()))
			return false
		}
	}
	return !result
}

func (ev *applyingECLVisitor) VisitEclrefinement(ctx *ecl.EclrefinementContext) interface{} {
	if ctx.Subrefinement() != nil {
		if result, ok := ev.Visit(ctx.Subrefinement()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing refinement: %s", ctx.Subrefinement().GetText()))
	}
	if ctx.Conjunctionrefinementset() != nil {
		if result, ok := ev.Visit(ctx.Conjunctionrefinementset()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing conjunction refinement set: %s", ctx.Conjunctionrefinementset().GetText()))
	}
	if ctx.Disjunctionrefinementset() != nil {
		if result, ok := ev.Visit(ctx.Disjunctionrefinementset()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing disjunction refinement set: %s", ctx.Disjunctionrefinementset().GetText()))
	}
	return false
}

// conjunctionRefinementSet = 1*(ws conjunction ws subRefinement)
func (ev *applyingECLVisitor) VisitConjunctionrefinementset(ctx *ecl.ConjunctionrefinementsetContext) interface{} {
	result := true
	for _, sr := range ctx.AllSubrefinement() {
		if term, ok := ev.Visit(sr).(bool); ok {
			result = result && term
		} else {
			ev.addError(fmt.Errorf("error parsing subrefinement: %s", ctx.GetText()))
			return false
		}
	}
	return result
}

func (ev *applyingECLVisitor) VisitDisjunctionrefinementset(ctx *ecl.DisjunctionrefinementsetContext) interface{} {
	result := false
	for _, sr := range ctx.AllSubrefinement() {
		if term, ok := ev.Visit(sr).(bool); ok {
			result = result || term
		} else {
			ev.addError(fmt.Errorf("error parsing subrefinement: %s", ctx.GetText()))
			return false
		}
	}
	return result
}

// subRefinement = eclAttributeSet / eclAttributeGroup / "(" ws eclRefinement ws ")"
func (ev *applyingECLVisitor) VisitSubrefinement(ctx *ecl.SubrefinementContext) interface{} {
	if ctx.Eclattributeset() != nil {
		if result, ok := ev.Visit(ctx.Eclattributeset()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing attribute set: %s", ctx.Eclattributeset().GetText()))
	}
	if ctx.Eclattributegroup() != nil {
		if result, ok := ev.Visit(ctx.Eclattributegroup()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing attribute group: %s", ctx.Eclattributegroup().GetText()))
	}
	if ctx.Eclrefinement() != nil {
		if result, ok := ev.Visit(ctx.Eclrefinement()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("error parsing refinement: %s", ctx.Eclrefinement().GetText()))
	}
	return false
}

// eclAttributeGroup = ["[" cardinality "]" ws] "{" ws eclAttributeSet ws "}"
//	cardinality : minvalue to maxvalue;
// 	minvalue : nonnegativeintegervalue;
//  to : (PERIOD PERIOD);
//	maxvalue : nonnegativeintegervalue | many;
func (ev *applyingECLVisitor) VisitEclattributegroup(ctx *ecl.EclattributegroupContext) interface{} {
	var card *cardinality
	var ok bool
	if ctx.Cardinality() != nil {
		if card, ok = ev.Visit(ctx.Cardinality()).(*cardinality); !ok {
			ev.addError(fmt.Errorf("error parsing cardinality: %s", ctx.Cardinality().GetText()))
		}
	}
	// check that cardinality is valid
	// TODO: move this check into cardinality visitor
	if card != nil && card.minimumValue == 0 && card.toMany {
		ev.addError(fmt.Errorf("invalid cardinality: redundant to have 0 to many cardinality: %s", ctx.GetText()))
		return false
	}
	// TODO: not implemented
	if card != nil {
		ev.addError(fmt.Errorf("error: cardinality for attribute groups not implemented: %s", ctx.GetText()))
		return false
	}

	if result, ok := ev.Visit(ctx.Eclattributeset()).(bool); ok {
		return result
	}
	ev.addError(fmt.Errorf("error parsing attribute set: %s", ctx.Eclattributeset().GetText()))
	return false
}

/*

// eclAttribute = ["[" cardinality "]" ws] [reverseFlag ws] eclAttributeName ws (expressionComparisonOperator ws subExpressionConstraint / numericComparisonOperator ws "#" numericValue / stringComparisonOperator ws QM stringValue QM)
func (ev *applyingECLVisitor) VisitEclattribute(ctx *ecl.EclattributeContext) interface{} {
	var card *cardinality
	ok := false
	reverse := false
	if ctx.Cardinality() != nil {
		if card, ok = ev.Visit(ctx.Cardinality()).(*cardinality); !ok {
			ev.addError(fmt.Errorf("invalid cardinality: %s", ctx.Cardinality().GetText()))
		}
	}
	if ctx.Reverseflag() != nil {
		reverse = true
	}
	if attributeName, ok := ev.Visit(ctx.Eclattributename()).(*snomed.ExpressionConstraint_Subexpression); ok {
		result.Name = attributeName
	} else {
		ev.addError(fmt.Errorf("invalid attribute name: %s", ctx.Eclattributename().GetText()))
	}
	if ctx.Expressioncomparisonoperator() != nil && ctx.Subexpressionconstraint() != nil {
		if operator, ok := ev.Visit(ctx.Expressioncomparisonoperator()).(snomed.ExpressionConstraint_ComparisonOperator); ok {
			result.Operator = operator
		} else {
			ev.addError(fmt.Errorf("invalid expression comparison operator: %s", ctx.Expressioncomparisonoperator().GetText()))
		}
		if exp, ok := ev.Visit(ctx.Subexpressionconstraint()).(*snomed.ExpressionConstraint_Subexpression); ok {
			result.Value = &snomed.ExpressionConstraint_Attribute_SubexpressionValue{
				SubexpressionValue: exp,
			}
		} else {
			ev.addError(fmt.Errorf("invalid subexpression in attribute: %s", ctx.GetText()))
		}
	}
	if ctx.Numericcomparisonoperator() != nil && ctx.Numericvalue() != nil {
		if operator, ok := ev.Visit(ctx.Numericcomparisonoperator()).(snomed.ExpressionConstraint_ComparisonOperator); ok {
			result.Operator = operator
		} else {
			ev.addError(fmt.Errorf("invalid numeric comparison operator: %s", ctx.Numericcomparisonoperator().GetText()))
		}
		if v, ok := ev.Visit(ctx.Numericvalue()).(float64); ok {
			result.Value = &snomed.ExpressionConstraint_Attribute_DoubleValue{
				DoubleValue: v,
			}
		}
		if v, ok := ev.Visit(ctx.Numericvalue()).(int64); ok {
			result.Value = &snomed.ExpressionConstraint_Attribute_IntegerValue{
				IntegerValue: v,
			}
		}
	}
	if ctx.Stringcomparisonoperator() != nil && ctx.Stringvalue() != nil {
		if operator, ok := ev.Visit(ctx.Stringcomparisonoperator()).(snomed.ExpressionConstraint_ComparisonOperator); ok {
			result.Operator = operator
		} else {
			ev.addError(fmt.Errorf("invalid string comparison operator: %s", ctx.Numericcomparisonoperator().GetText()))
		}
		result.Value = &snomed.ExpressionConstraint_Attribute_StringValue{
			StringValue: ctx.Stringvalue().GetText(),
		}
	}
	return result
}

// eclAttributeName = subExpressionConstraint
func (ev *applyingECLVisitor) VisitEclattributename(ctx *ecl.EclattributenameContext) interface{} {
	return ev.Visit(ctx.Subexpressionconstraint())
}

// expressionComparisonOperator = "=" / "!="
func (ev *applyingECLVisitor) VisitExpressioncomparisonoperator(ctx *ecl.ExpressioncomparisonoperatorContext) interface{} {
	return ev.parseComparisonOperator(ctx.GetText())
}

// numericComparisonOperator = "=" / "!=" / "<=" / "<" / ">=" / ">"
func (ev *applyingECLVisitor) VisitNumericcomparisonoperator(ctx *ecl.NumericcomparisonoperatorContext) interface{} {
	return ev.parseComparisonOperator(ctx.GetText())
}

func (ev *applyingECLVisitor) VisitStringcomparisonoperator(ctx *ecl.StringcomparisonoperatorContext) interface{} {
	return ev.parseComparisonOperator(ctx.GetText())
}

func (ev *applyingECLVisitor) parseComparisonOperator(operator string) snomed.ExpressionConstraint_ComparisonOperator {
	switch operator {
	case "=":
		return equals
	case "!=":
		return notEquals
	case "<=":
		return lessThanOrEquals
	case "<":
		return lessThan
	case ">=":
		return greaterThanOrEquals
	case ">":
		return greaterThan
	default:
		ev.addError(fmt.Errorf("invalid comparison operator: %s", operator))
		return 0
	}
}

// cardinality = minValue to maxValue
// minValue = nonNegativeIntegerValue
func (ev *applyingECLVisitor) VisitCardinality(ctx *ecl.CardinalityContext) interface{} {
	result := new(cardinality)
	if ctx.Minvalue() != nil {
		if v, ok := ev.Visit(ctx.Minvalue()).(int64); ok {
			result.minimumValue = v
		} else {
			ev.addError(fmt.Errorf("invalid integer: %s", ctx.Minvalue().GetText()))
			return nil
		}
	}
	if ctx.Maxvalue() != nil {
		if v, ok := ev.Visit(ctx.Maxvalue()).(int64); ok {
			result.maximumValue = v
		} else {
			ev.addError(fmt.Errorf("invalid integer: %s", ctx.Maxvalue().GetText()))
			return nil
		}
	}
	return result
}

func (ev *applyingECLVisitor) VisitMinvalue(ctx *ecl.MinvalueContext) interface{} {
	v, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		ev.addError(fmt.Errorf("invalid integer value: %s", ctx.GetText()))
	}
	return v
}

// Visit a parse tree produced by ECLParser#maxvalue.
func (ev *applyingECLVisitor) VisitMaxvalue(ctx *ecl.MaxvalueContext) interface{} {
	v, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		ev.addError(fmt.Errorf("invalid integer value: %s", ctx.GetText()))
	}
	return v
}

// subExpressionConstraint = [constraintOperator ws] [memberOf ws] (eclFocusConcept / "(" ws expressionConstraint ws ")")
func (ev *applyingECLVisitor) VisitSubexpressionconstraint(ctx *ecl.SubexpressionconstraintContext) interface{} {
	var cop constraintOperator
	var ok bool
	var memberOf bool
	var fc *focusConcept
	if ctx.Constraintoperator() != nil {
		if cop, ok = ev.Visit(ctx.Constraintoperator()).(constraintOperator); !ok {
			ev.addError(fmt.Errorf("error parsing constraint operator: %s", ctx.GetText()))
			return false
		}
	}
	if ctx.Memberof() != nil {
		memberOf = true
	}
	if ctx.Eclfocusconcept() != nil {
		if fc, ok = ev.Visit(ctx.Eclfocusconcept()).(*focusConcept); !ok {
			ev.addError(fmt.Errorf("error parsing focus concept: %s", ctx.Eclfocusconcept().GetText()))
			return false
		}
	}
	if ctx.Expressionconstraint() != nil {
		if exp, ok := ev.Visit(ctx.Expressionconstraint()).(*snomed.ExpressionConstraint); ok {
			result.Value = &snomed.ExpressionConstraint_Subexpression_Expression{
				Expression: exp,
			}
		} else {
			ev.addError(fmt.Errorf("invalid expression: %s", ctx.Expressionconstraint().GetText()))
		}
	}
	return result
}

// eclAttributeSet = subAttributeSet ws [conjunctionAttributeSet / disjunctionAttributeSet]
func (ev *applyingECLVisitor) VisitEclattributeset(ctx *ecl.EclattributesetContext) interface{} {
	result := false
	if ctx.Subattributeset() != nil {
		if term, ok := ev.Visit(ctx.Subattributeset()).(bool); ok {
			result = term
		} else {
			ev.addError(fmt.Errorf("error parsing subattribute set: %s", ctx.Subattributeset().GetText()))
			return false
		}
	}
	if ctx.Conjunctionattributeset() != nil {
		if term, ok := ev.Visit(ctx.Conjunctionattributeset()).(bool); ok {
			result = result && term
		} else {
			ev.addError(fmt.Errorf("error parsing conjunction attribute set: %s", ctx.Conjunctionattributeset().GetText()))
			return false
		}
	}
	if ctx.Disjunctionattributeset() != nil {
		if term, ok := ev.Visit(ctx.Disjunctionattributeset()).(bool); ok {
			result = result || term
		} else {
			ev.addError(fmt.Errorf("invalid disjunction attribute set: %s", ctx.Disjunctionattributeset().GetText()))
			return false
		}
	}
	return result
}

// subAttributeSet = eclAttribute / "(" ws eclAttributeSet ws ")"
func (ev *applyingECLVisitor) VisitSubattributeset(ctx *ecl.SubattributesetContext) interface{} {
	if ctx.Eclattribute() != nil {
		if result, ok := ev.Visit(ctx.Eclattribute()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("invalid attribute: %s", ctx.Eclattribute().GetText()))
	}
	if ctx.Eclattributeset() != nil {
		if result, ok := ev.Visit(ctx.Eclattributeset()).(bool); ok {
			return result
		}
		ev.addError(fmt.Errorf("invalid attribute set: %s", ctx.Eclattributeset().GetText()))

	}
	return false
}

// conjunctionAttributeSet = 1*(ws conjunction ws subAttributeSet)
func (ev *applyingECLVisitor) VisitConjunctionattributeset(ctx *ecl.ConjunctionattributesetContext) interface{} {
	result := true
	for _, subattributeset := range ctx.AllSubattributeset() {
		if term, ok := ev.Visit(subattributeset).(bool); ok {
			result = result && term
		} else {
			ev.addError(fmt.Errorf("invalid subattribute set : %s", ctx.GetText()))
			return false
		}
	}
	return result
}

// disjunctionAttributeSet = 1*(ws disjunction ws subAttributeSet)
func (ev *applyingECLVisitor) VisitDisjunctionattributeset(ctx *ecl.DisjunctionattributesetContext) interface{} {
	result := false
	for _, subattributeset := range ctx.AllSubattributeset() {
		if term, ok := ev.Visit(subattributeset).(bool); ok {
			result = result || term
		} else {
			ev.addError(fmt.Errorf("error parsing subattribute set : %s", ctx.GetText()))
			return false
		}
	}
	return result
}

func (ev *applyingECLVisitor) VisitConstraintoperator(ctx *ecl.ConstraintoperatorContext) interface{} {
	fmt.Printf("parsing constraint: '%s'\nctx:%v\n", ctx.GetText(), ctx)
	switch {
	case ctx.Ancestorof() != nil:
		fmt.Printf("ancestor of\n")
		return ancestorOf
	case ctx.Ancestororselfof() != nil:
		fmt.Printf("ancestor or self of\n")
		return ancestorOrSelf
	case ctx.Childof() != nil:
		fmt.Printf("child of: %v\n", ctx.Childof())
		return childOf
	case ctx.Descendantof() != nil:
		return descendantOf
	case ctx.Descendantorselfof() != nil:
		return descendantOrSelf
	case ctx.Parentof() != nil:
		return parentOf
	default:
		ev.addError(fmt.Errorf("unknown constraint operator: %s", ctx.GetText()))
		return nil
	}
}

// eclFocusConcept = eclConceptReference / wildCard
func (ev *applyingECLVisitor) VisitEclfocusconcept(ctx *ecl.EclfocusconceptContext) interface{} {
	fc := new(focusConcept)
	if ctx.Eclconceptreference() != nil {
		if cr, ok := ev.Visit(ctx.Eclconceptreference()).(*snomed.ConceptReference); ok {
			fc.concept = cr
		} else {
			ev.addError(fmt.Errorf("invalid concept reference: %s", ctx.GetText()))
		}
	}
	if ctx.Wildcard() != nil {
		fc.wildcard = true
	}
	return fc
}

// eclConceptReference = conceptId [ws "|" ws term ws "|"]
func (ev *applyingECLVisitor) VisitEclconceptreference(ctx *ecl.EclconceptreferenceContext) interface{} {
	cr := new(snomed.ConceptReference)
	if ctx.Conceptid() != nil {
		if conceptID, ok := ev.Visit(ctx.Conceptid()).(int64); ok {
			cr.ConceptId = conceptID
		} else {
			ev.addError(fmt.Errorf("invalid concept id: %s", ctx.Conceptid().GetText()))
		}
	}
	if ctx.Term() != nil {
		if term, ok := ev.Visit(ctx.Term()).(string); ok {
			cr.Term = term
		}
	}
	return cr
}

// returns int64
// conceptId = sctId
func (ev *applyingECLVisitor) VisitConceptid(ctx *ecl.ConceptidContext) interface{} {
	return ev.Visit(ctx.Sctid())
}

// returns string
func (ev *applyingECLVisitor) VisitTerm(ctx *ecl.TermContext) interface{} {
	return ctx.GetText()
}

// TODO: check escaping rules - do we need to parse - probably?
// ABNF: stringValue = 1*(anyNonEscapedChar / escapedChar)
func (ev *applyingECLVisitor) VisitStringvalue(ctx *ecl.StringvalueContext) interface{} {
	return ctx.GetText()
}

// numericValue = ["-"/"+"] (decimalValue / integerValue)
// returns either an int64 or float64
func (ev *applyingECLVisitor) VisitNumericvalue(ctx *ecl.NumericvalueContext) interface{} {
	sign := 1
	if ctx.DASH() != nil {
		sign = -1
	}
	if ctx.Integervalue() != nil {
		return int64(sign) * ev.Visit(ctx.Integervalue()).(int64)
	}
	if ctx.Decimalvalue() != nil {
		return float64(sign) * ev.Visit(ctx.Decimalvalue()).(float64)
	}
	ev.addError(fmt.Errorf("invalid numeric value: %s", ctx.GetText()))
	return 0
}

// returns int64
func (ev *applyingECLVisitor) VisitIntegervalue(ctx *ecl.IntegervalueContext) interface{} {
	v, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		ev.addError(fmt.Errorf("invalid integer value: %s", ctx.GetText()))
	}
	return v
}

// returns float64
func (ev *applyingECLVisitor) VisitDecimalvalue(ctx *ecl.DecimalvalueContext) interface{} {
	v, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		ev.addError(fmt.Errorf("invalid double value: %s", ctx.GetText()))
	}
	return v
}

// returns a SNOMED CT identifier (int64)
func (ev *applyingECLVisitor) VisitSctid(ctx *ecl.SctidContext) interface{} {
	sct, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		ev.addError(err)
	}
	return sct
}

type predicate interface {
	Test(svc *terminology.Svc, exp *snomed.Expression) bool
}

// fixedPredicate will always return the result as yet, irrespective of the expression on test
type fixedPredicate struct {
	result bool
}

func (fp *fixedPredicate) Test(svc *terminology.Svc, exp *snomed.Expression) bool {
	return fp.result
}

// memberPredicate will return true when a focus concept is a member of the set specified
type memberPredicate struct {
	members map[int64]struct{}
}

func (fp *memberPredicate) Test(svc *terminology.Svc, exp *snomed.Expression) bool {
	for _, fc := range exp.GetClause().GetFocusConcepts() {
		if _, ok := fp.members[fc.ConceptId]; ok {
			return true
		}
	}
	return false
}
*/
