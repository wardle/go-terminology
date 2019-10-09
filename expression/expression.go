// Package expression provides functionality to process SNOMED CT expressions
// Expressions are usually multiple SNOMED CT concepts combined together, much
// like a sentence is made up of words.
//
// SNOMED CT contains single concepts that actually represent expressions; these are usually
// historic or exist for ease of use. In order to appropriately determine equivalence, a range
// of functions are required to normalise any arbitrary concept or expression into a normalised
// form.
//
// The ANTLR parser was generated from the original ABNF source file
// using http://www.robertpinchbeck.com/abnf_to_antlr/Default.aspx and running
// java -jar ~/Downloads/antlr-4.7.2-complete.jar -Dlanguage=Go -package ecl -o ecl ECL.g4
// java -jar ~/Downloads/antlr-4.7.2-complete.jar -Dlanguage=Go -package cg -o cg CG.g4
// The compositional grammar (CG) is from https://confluence.ihtsdotools.org/display/DOCSCG/5.1+Normative+Specification
// The expression constraint grammar (ECL) is from https://confluence.ihtsdotools.org/pages/viewpage.action?pageId=28739405
package expression

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/wardle/go-terminology/expression/cg"
	"github.com/wardle/go-terminology/expression/ecl"
	"github.com/wardle/go-terminology/snomed"
)

// CreateSimpleExpression creates an expression from a single concept
func CreateSimpleExpression(concept *snomed.Concept) *snomed.Expression {
	exp := snomed.Expression{}
	clause := &snomed.Expression_Clause{}
	concepts := make([]*snomed.ConceptReference, 0)
	concepts = append(concepts, &snomed.ConceptReference{ConceptId: concept.Id})
	clause.FocusConcepts = concepts
	exp.Clause = clause
	return &exp
}

// Equal determines whether two expressions are exactly equal
func Equal(e1 *snomed.Expression, e2 *snomed.Expression) bool {
	renderer := NewCanonicalRenderer()
	s1, _ := renderer.Render(e1)
	s2, _ := renderer.Render(e2)
	return strings.EqualFold(s1, s2)
}

// Parse parses a SNOMED expression
func Parse(s string) (*snomed.Expression, error) {
	l := new(cgListener)
	is := antlr.NewInputStream(s)
	lexer := cg.NewCGLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := cg.NewCGParser(stream)
	p.RemoveErrorListeners() // remove default listeners, which includes console listener - so as to avoid printing all parse errors to console
	el := new(errorListener)
	p.AddErrorListener(el)
	antlr.ParseTreeWalkerDefault.Walk(l, p.Expression())
	return l.expression, el.err
}

// ParseError returns information about a parsing error
type ParseError struct {
	Line, Column   int
	OffendingToken string
	Msg            string
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("syntax error: line %d:%d %s", pe.Line, pe.Column, pe.Msg)
}

// todo: handle multiple errors, if they exist, and pass back a special error struct that permits return of structured error information
type errorListener struct {
	*antlr.DefaultErrorListener
	err error
}

func (el *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	token := ""
	if e != nil {
		if e.GetOffendingToken() != nil {
			token = e.GetOffendingToken().GetText()
		}
	}
	el.err = &ParseError{
		Line:           line,
		Column:         column,
		OffendingToken: token,
		Msg:            msg,
	}
}

// cgListener is an internal ANTLR listener.
// The generated code from ANTLR is quirky and not idiomatic and I am
// probably using the listener incorrectly. It would be better to use the visitor
// pattern, but that appears incomplete as of Antlr 4.7.1.
type cgListener struct {
	cg.BaseCGListener
	expression *snomed.Expression
	err        error
}

func (cgl *cgListener) EnterEquivalentto(c *cg.EquivalenttoContext) {
	cgl.expression.DefinitionStatus = snomed.Expression_EQUIVALENT_TO
}

func (cgl *cgListener) EnterSubtypeof(c *cg.SubtypeofContext) {
	cgl.expression.DefinitionStatus = snomed.Expression_SUBTYPE_OF
}

func (cgl *cgListener) EnterExpression(ctx *cg.ExpressionContext) {
	if cgl.expression == nil {
		cgl.expression = new(snomed.Expression)
		se, ok := ctx.Subexpression().(*cg.SubexpressionContext)
		if !ok {
			cgl.err = errors.New("No valid subexpression identified")
			return
		}
		clause, err := parseSubexpression(se)
		if err != nil {
			cgl.err = err
			return
		}
		cgl.expression.Clause = clause
		return
	}
}

func parseSubexpression(ctx *cg.SubexpressionContext) (*snomed.Expression_Clause, error) {
	clause := new(snomed.Expression_Clause)
	fc, ok := ctx.Focusconcept().(*cg.FocusconceptContext)
	if !ok {
		return nil, errors.New("invalid expression: missing focus concept(s)")
	}
	focus, err := parseFocusConcepts(fc)
	if err != nil {
		return nil, err
	}
	clause.FocusConcepts = focus
	// do we have any refinements at all?
	refinements, ok := ctx.Refinement().(*cg.RefinementContext)
	if !ok {
		return clause, nil
	}
	// process any simple ungrouped refinements, if we have them
	set, ok := refinements.Attributeset().(*cg.AttributesetContext)
	if ok {
		refs, err := parseRefinements(set)
		if err != nil {
			return nil, err
		}
		clause.Refinements = refs
	}
	// process any grouped refinements
	grps, err := parseRefinementGroups(refinements.AllAttributegroup())
	if err != nil {
		return nil, err
	}
	clause.RefinementGroups = grps
	return clause, nil
}

func parseFocusConcepts(ctx *cg.FocusconceptContext) ([]*snomed.ConceptReference, error) {
	all := ctx.AllConceptreference()
	concepts := make([]*snomed.ConceptReference, len(all))
	for i, cr := range all {
		crr, ok := cr.(*cg.ConceptreferenceContext)
		if !ok {
			return nil, errors.New("invalid expression: missing concept reference")
		}
		concept, err := parseConceptReference(crr)
		if err != nil {
			return nil, err
		}
		concepts[i] = concept
	}
	return concepts, nil
}

func parseRefinementGroups(ctx []cg.IAttributegroupContext) ([]*snomed.Expression_RefinementGroup, error) {
	groups := make([]*snomed.Expression_RefinementGroup, len(ctx))
	for i, g := range ctx {
		ag := g.(*cg.AttributegroupContext)
		agg, ok := ag.Attributeset().(*cg.AttributesetContext)
		if !ok {
			return nil, errors.New("invalid expression: missing attribute set")
		}
		refinements, err := parseRefinements(agg)
		if err != nil {
			return nil, err
		}
		group := new(snomed.Expression_RefinementGroup)
		group.Refinements = refinements
		groups[i] = group
	}
	return groups, nil
}

func parseRefinements(ctx *cg.AttributesetContext) ([]*snomed.Expression_Refinement, error) {
	all := ctx.AllAttribute()
	attributes := make([]*snomed.Expression_Refinement, len(all))
	for i, attr := range all {
		attribute, err := parseAttribute(attr.(*cg.AttributeContext))
		if err != nil {
			return nil, err
		}
		attributes[i] = attribute
	}
	return attributes, nil
}

func parseAttribute(ctx *cg.AttributeContext) (*snomed.Expression_Refinement, error) {
	refinement := new(snomed.Expression_Refinement)
	name := ctx.Attributename().(*cg.AttributenameContext)
	nameRef, err := parseConceptReference(name.Conceptreference().(*cg.ConceptreferenceContext))
	if err != nil {
		return nil, err
	}
	refinement.RefinementConcept = nameRef
	value := ctx.Attributevalue().(*cg.AttributevalueContext)
	if v := value.Stringvalue(); v != nil {
		refinement.Value = &snomed.Expression_Refinement_StringValue{StringValue: v.GetText()}
	}
	if v := value.Expressionvalue(); v != nil {
		e := v.(*cg.ExpressionvalueContext)
		if e.Conceptreference() != nil {
			cr, err := parseConceptReference(e.Conceptreference().(*cg.ConceptreferenceContext))
			if err != nil {
				return nil, err
			}
			refinement.Value = &snomed.Expression_Refinement_ConceptValue{ConceptValue: cr}
		}
		if e.Subexpression() != nil {
			se, err := parseSubexpression(e.Subexpression().(*cg.SubexpressionContext))
			if err != nil {
				return nil, err
			}
			refinement.Value = &snomed.Expression_Refinement_ClauseValue{ClauseValue: se}
		}
	}
	if v := value.Numericvalue(); v != nil {
		n := v.(*cg.NumericvalueContext)
		if n.Integervalue() != nil {
			iv, err := strconv.ParseInt(n.Integervalue().(*cg.IntegervalueContext).GetText(), 10, 64)
			if err != nil {
				return nil, err
			}
			refinement.Value = &snomed.Expression_Refinement_IntValue{IntValue: iv}
		}
		if n.Decimalvalue() != nil {
			dv, err := strconv.ParseFloat(n.Decimalvalue().(*cg.DecimalvalueContext).GetText(), 64)
			if err != nil {
				return nil, err
			}
			refinement.Value = &snomed.Expression_Refinement_DoubleValue{DoubleValue: dv}
		}
	}
	return refinement, nil
}

func parseConceptReference(ctx *cg.ConceptreferenceContext) (ref *snomed.ConceptReference, err error) {
	ref = new(snomed.ConceptReference)
	ref.ConceptId, err = strconv.ParseInt(ctx.Conceptid().GetText(), 10, 64)
	if err != nil {
		return nil, err
	}
	if term := ctx.Term(); term != nil {
		ref.Term = ctx.Term().GetText()
	}
	return
}

// ParseConstraint parses an expression in the "Expression Constraint Language" (ECL)
func ParseConstraint(s string) (result *snomed.ExpressionConstraint, err error) {
	is := antlr.NewInputStream(s)
	lex := ecl.NewECLLexer(is)
	tokens := antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
	p := ecl.NewECLParser(tokens)
	visitor := new(eclVisitor)
	tree := p.Expressionconstraint()
	var ok bool
	if result, ok = visitor.Visit(tree).(*snomed.ExpressionConstraint); !ok {
		return nil, fmt.Errorf("error: did not parse valid expression constraint. got: %v", result)
	}
	if len(visitor.errors) > 0 {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d error(s) parsing expression: ", len(visitor.errors)))
		for _, e := range visitor.errors {
			sb.WriteString(fmt.Sprintf("%s ", e))
		}
		err = errors.New(strings.TrimSpace(sb.String()))
	}
	return
}

type eclVisitor struct {
	ecl.BaseECLVisitor
	errors []error
}

func (ev *eclVisitor) addError(err error) {
	ev.errors = append(ev.errors, err)
}

func (ev *eclVisitor) Visit(tree antlr.ParseTree) interface{} {
	if tree == nil {
		return nil
	}
	return tree.Accept(ev)
}

func (ev *eclVisitor) VisitExpressionconstraint(ctx *ecl.ExpressionconstraintContext) interface{} {
	result := new(snomed.ExpressionConstraint)
	if ctx.Compoundexpressionconstraint() != nil {
		if compound, ok := ev.Visit(ctx.Compoundexpressionconstraint()).(*snomed.ExpressionConstraint_Compound); ok {
			result.Type = &snomed.ExpressionConstraint_Compound_{
				Compound: compound,
			}
		} else {
			ev.addError(fmt.Errorf("invalid compound expression: %s", ctx.Compoundexpressionconstraint().GetText()))
		}
	}
	if dotted := ctx.Dottedexpressionconstraint(); dotted != nil {
		fmt.Printf("dotted")
	}
	if ctx.Subexpressionconstraint() != nil {
		if subexp, ok := ev.Visit(ctx.Subexpressionconstraint()).(*snomed.ExpressionConstraint_Subexpression); ok {
			result.Type = &snomed.ExpressionConstraint_Subexpression_{
				Subexpression: subexp,
			}
		} else {
			ev.addError(fmt.Errorf("invalid subexpression: %s", ctx.Subexpressionconstraint().GetText()))
		}
	}
	if ctx.Refinedexpressionconstraint() != nil {
		if refined, ok := ev.Visit(ctx.Refinedexpressionconstraint()).(*snomed.ExpressionConstraint_Refined); ok {
			result.Type = &snomed.ExpressionConstraint_Refined_{
				Refined: refined,
			}
		} else {
			ev.addError(fmt.Errorf("invalid refined expression: %s", ctx.Refinedexpressionconstraint().GetText()))
		}
	}
	fmt.Printf("Parsed: %s\nResult: %v\n\n", ctx.GetText(), result)
	return result
}
func (ev *eclVisitor) VisitRefinedexpressionconstraint(ctx *ecl.RefinedexpressionconstraintContext) interface{} {
	result := new(snomed.ExpressionConstraint_Refined)
	if ctx.Subexpressionconstraint() != nil {
		if subexp, ok := ev.Visit(ctx.Subexpressionconstraint()).(*snomed.ExpressionConstraint_Subexpression); ok {
			result.Subexpression = subexp
		} else {
			ev.addError(fmt.Errorf("invalid subexpression: %s", ctx.Subexpressionconstraint().GetText()))
		}
	}
	if ctx.Eclrefinement() != nil {
		if refinement, ok := ev.Visit(ctx.Eclrefinement()).(*snomed.ExpressionConstraint_Refinement); ok {
			result.Refinement = refinement
		} else {
			ev.addError(fmt.Errorf("invalid refinement: %s", ctx.Eclrefinement().GetText()))
		}
	}
	return result
}

// compoundExpressionConstraint = conjunctionExpressionConstraint / disjunctionExpressionConstraint / exclusionExpressionConstraint
func (ev *eclVisitor) VisitCompoundexpressionconstraint(ctx *ecl.CompoundexpressionconstraintContext) interface{} {
	if ctx.Conjunctionexpressionconstraint() != nil {
		if and, ok := ev.Visit(ctx.Conjunctionexpressionconstraint()).(*snomed.ExpressionConstraint_Compound); ok {
			return and
		}
	}
	if ctx.Disjunctionexpressionconstraint() != nil {
		if or, ok := ev.Visit(ctx.Disjunctionexpressionconstraint()).(*snomed.ExpressionConstraint_Compound); ok {
			return or
		}
	}
	if ctx.Exclusionexpressionconstraint() != nil {
		if not, ok := ev.Visit(ctx.Exclusionexpressionconstraint()).(*snomed.ExpressionConstraint_Compound); ok {
			return not
		}
	}
	return nil
}

// conjunctionExpressionConstraint = subExpressionConstraint 1*(ws conjunction ws subExpressionConstraint)
func (ev *eclVisitor) VisitConjunctionexpressionconstraint(ctx *ecl.ConjunctionexpressionconstraintContext) interface{} {
	result := new(snomed.ExpressionConstraint_Compound)
	result.Logical = snomed.ExpressionConstraint_AND
	terms := make([]*snomed.ExpressionConstraint_Subexpression, 0)
	for i, subexp := range ctx.AllSubexpressionconstraint() {
		if se, ok := ev.Visit(subexp).(*snomed.ExpressionConstraint_Subexpression); ok {
			if i == 0 {
				result.Subexpression = se
			} else {
				terms = append(terms, se)
			}
		} else {
			ev.addError(fmt.Errorf("invalid subexpression: %s", subexp.GetText()))
			return nil
		}
	}
	result.Terms = terms
	return result
}

func (ev *eclVisitor) VisitDisjunctionexpressionconstraint(ctx *ecl.DisjunctionexpressionconstraintContext) interface{} {
	result := new(snomed.ExpressionConstraint_Compound)
	result.Logical = snomed.ExpressionConstraint_OR
	terms := make([]*snomed.ExpressionConstraint_Subexpression, 0)
	for i, subexp := range ctx.AllSubexpressionconstraint() {
		if se, ok := ev.Visit(subexp).(*snomed.ExpressionConstraint_Subexpression); ok {
			if i == 0 {
				result.Subexpression = se
			} else {
				terms = append(terms, se)
			}
		} else {
			ev.addError(fmt.Errorf("invalid subexpression: %s", subexp.GetText()))
			return nil
		}
	}
	result.Terms = terms
	return result

}

func (ev *eclVisitor) VisitExclusionexpressionconstraint(ctx *ecl.ExclusionexpressionconstraintContext) interface{} {
	result := new(snomed.ExpressionConstraint_Compound)
	result.Logical = snomed.ExpressionConstraint_NOT
	terms := make([]*snomed.ExpressionConstraint_Subexpression, 0)
	for i, subexp := range ctx.AllSubexpressionconstraint() {
		if se, ok := ev.Visit(subexp).(*snomed.ExpressionConstraint_Subexpression); ok {
			if i == 0 {
				result.Subexpression = se
			} else {
				terms = append(terms, se)
			}
		} else {
			ev.addError(fmt.Errorf("invalid subexpression: %s", subexp.GetText()))
			return nil
		}
	}
	result.Terms = terms
	return result
}

func (ev *eclVisitor) VisitEclrefinement(ctx *ecl.EclrefinementContext) interface{} {
	result := new(snomed.ExpressionConstraint_Refinement)
	if ctx.Subrefinement() != nil {
		if subrefinement, ok := ev.Visit(ctx.Subrefinement()).(*snomed.ExpressionConstraint_Subrefinement); ok {
			result.Subrefinement = subrefinement
		} else {
			ev.addError(fmt.Errorf("invalid refinement: %s", ctx.Subrefinement().GetText()))
		}
	}
	if ctx.Conjunctionrefinementset() != nil {
		if conjunction, ok := ev.Visit(ctx.Conjunctionrefinementset()).(*snomed.ExpressionConstraint_ConjunctionRefinementSet); ok {
			result.Logical = &snomed.ExpressionConstraint_Refinement_ConjunctionRefinementSet{
				ConjunctionRefinementSet: conjunction,
			}
		} else {
			ev.addError(fmt.Errorf("invalid conjunction refinement set: %s", ctx.Conjunctionrefinementset().GetText()))
		}
	}
	if ctx.Disjunctionrefinementset() != nil {
		if disjunction, ok := ev.Visit(ctx.Disjunctionrefinementset()).(*snomed.ExpressionConstraint_DisjunctionRefinementSet); ok {
			result.Logical = &snomed.ExpressionConstraint_Refinement_DisjunctionRefinementSet{
				DisjunctionRefinementSet: disjunction,
			}
		} else {
			ev.addError(fmt.Errorf("invalid disjunction refinement set: %s", ctx.Disjunctionrefinementset().GetText()))
		}
	}
	return result
}

// conjunctionRefinementSet = 1*(ws conjunction ws subRefinement)
func (ev *eclVisitor) VisitConjunctionrefinementset(ctx *ecl.ConjunctionrefinementsetContext) interface{} {
	result := new(snomed.ExpressionConstraint_ConjunctionRefinementSet)
	subrefinements := make([]*snomed.ExpressionConstraint_Subrefinement, 0)
	for _, sr := range ctx.AllSubrefinement() {
		if subrefinement, ok := ev.Visit(sr).(*snomed.ExpressionConstraint_Subrefinement); ok {
			subrefinements = append(subrefinements, subrefinement)
		} else {
			ev.addError(fmt.Errorf("invalid subrefinement: %s", ctx.GetText()))
		}
	}
	result.And = subrefinements
	return result
}

func (ev *eclVisitor) VisitDisjunctionrefinementset(ctx *ecl.DisjunctionrefinementsetContext) interface{} {
	result := new(snomed.ExpressionConstraint_DisjunctionRefinementSet)
	subrefinements := make([]*snomed.ExpressionConstraint_Subrefinement, 0)
	for _, sr := range ctx.AllSubrefinement() {
		if subrefinement, ok := ev.Visit(sr).(*snomed.ExpressionConstraint_Subrefinement); ok {
			subrefinements = append(subrefinements, subrefinement)
		} else {
			ev.addError(fmt.Errorf("invalid subrefinement: %s", ctx.GetText()))
		}
	}
	result.Or = subrefinements
	return result
}

// subRefinement = eclAttributeSet / eclAttributeGroup / "(" ws eclRefinement ws ")"
func (ev *eclVisitor) VisitSubrefinement(ctx *ecl.SubrefinementContext) interface{} {
	result := new(snomed.ExpressionConstraint_Subrefinement)
	if ctx.Eclattributeset() != nil {
		if attributeSet, ok := ev.Visit(ctx.Eclattributeset()).(*snomed.ExpressionConstraint_AttributeSet); ok {
			result.Value = &snomed.ExpressionConstraint_Subrefinement_AttributeSet{
				AttributeSet: attributeSet,
			}
		} else {
			ev.addError(fmt.Errorf("invalid attribute set: %s", ctx.Eclattributeset().GetText()))
		}
	}
	if ctx.Eclattributegroup() != nil {
		if attributeGroup, ok := ev.Visit(ctx.Eclattributegroup()).(*snomed.ExpressionConstraint_AttributeGroup); ok {
			result.Value = &snomed.ExpressionConstraint_Subrefinement_AttributeGroup{
				AttributeGroup: attributeGroup,
			}
		} else {
			ev.addError(fmt.Errorf("invalid attribute group: %s", ctx.Eclattributegroup().GetText()))
		}
	}
	if ctx.Eclrefinement() != nil {
		if refinement, ok := ev.Visit(ctx.Eclrefinement()).(*snomed.ExpressionConstraint_Refinement); ok {
			result.Value = &snomed.ExpressionConstraint_Subrefinement_Refinement{
				Refinement: refinement,
			}
		} else {
			ev.addError(fmt.Errorf("invalid refinement: %s", ctx.Eclrefinement().GetText()))
		}
	}
	return result
}

// eclAttributeGroup = ["[" cardinality "]" ws] "{" ws eclAttributeSet ws "}"
func (ev *eclVisitor) VisitEclattributegroup(ctx *ecl.EclattributegroupContext) interface{} {
	result := new(snomed.ExpressionConstraint_AttributeGroup)
	if ctx.Cardinality() != nil {
		if cardinality, ok := ev.Visit(ctx.Cardinality()).(*snomed.ExpressionConstraint_Cardinality); ok {
			result.Cardinality = cardinality
		} else {
			ev.addError(fmt.Errorf("invalid cardinality: %s", ctx.Cardinality().GetText()))
		}
	}
	if attributeset, ok := ev.Visit(ctx.Eclattributeset()).(*snomed.ExpressionConstraint_AttributeSet); ok {
		result.AttributeSet = attributeset
	} else {
		ev.addError(fmt.Errorf("invalid attribute set: %s", ctx.Eclattributeset().GetText()))
	}
	return result
}

// eclAttribute = ["[" cardinality "]" ws] [reverseFlag ws] eclAttributeName ws (expressionComparisonOperator ws subExpressionConstraint / numericComparisonOperator ws "#" numericValue / stringComparisonOperator ws QM stringValue QM)
func (ev *eclVisitor) VisitEclattribute(ctx *ecl.EclattributeContext) interface{} {
	result := new(snomed.ExpressionConstraint_Attribute)
	if ctx.Cardinality() != nil {
		if cardinality, ok := ev.Visit(ctx.Cardinality()).(*snomed.ExpressionConstraint_Cardinality); ok {
			result.Cardinality = cardinality
		} else {
			ev.addError(fmt.Errorf("invalid cardinality: %s", ctx.Cardinality().GetText()))
		}
	}
	if ctx.Reverseflag() != nil {
		result.Reverse = true
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
func (ev *eclVisitor) VisitEclattributename(ctx *ecl.EclattributenameContext) interface{} {
	return ev.Visit(ctx.Subexpressionconstraint())
}

// expressionComparisonOperator = "=" / "!="
func (ev *eclVisitor) VisitExpressioncomparisonoperator(ctx *ecl.ExpressioncomparisonoperatorContext) interface{} {
	return ev.parseComparisonOperator(ctx.GetText())
}

// numericComparisonOperator = "=" / "!=" / "<=" / "<" / ">=" / ">"
func (ev *eclVisitor) VisitNumericcomparisonoperator(ctx *ecl.NumericcomparisonoperatorContext) interface{} {
	return ev.parseComparisonOperator(ctx.GetText())
}

func (ev *eclVisitor) VisitStringcomparisonoperator(ctx *ecl.StringcomparisonoperatorContext) interface{} {
	return ev.parseComparisonOperator(ctx.GetText())
}

func (ev *eclVisitor) parseComparisonOperator(operator string) snomed.ExpressionConstraint_ComparisonOperator {
	switch operator {
	case "=":
		return snomed.ExpressionConstraint_EQUALS
	case "!=":
		return snomed.ExpressionConstraint_NOT_EQUALS
	case "<=":
		return snomed.ExpressionConstraint_LESS_THAN_OR_EQUALS
	case "<":
		return snomed.ExpressionConstraint_LESS_THAN
	case ">=":
		return snomed.ExpressionConstraint_GREATER_THAN_OR_EQUALS
	case ">":
		return snomed.ExpressionConstraint_GREATER_THAN
	default:
		ev.addError(fmt.Errorf("invalid comparison operator: %s", operator))
		return 0
	}
}

// cardinality = minValue to maxValue
// minValue = nonNegativeIntegerValue
func (ev *eclVisitor) VisitCardinality(ctx *ecl.CardinalityContext) interface{} {
	result := new(snomed.ExpressionConstraint_Cardinality)
	if ctx.Minvalue() != nil {
		if v, ok := ev.Visit(ctx.Minvalue()).(int64); ok {
			result.MinimumValue = v
		} else {
			ev.addError(fmt.Errorf("invalid integer: %s", ctx.Minvalue().GetText()))
		}
	}
	if ctx.Maxvalue() != nil {
		if v, ok := ev.Visit(ctx.Maxvalue()).(int64); ok {
			result.MaximumValue = v
		} else {
			ev.addError(fmt.Errorf("invalid integer: %s", ctx.Maxvalue().GetText()))
		}
	}
	return result
}

func (ev *eclVisitor) VisitMinvalue(ctx *ecl.MinvalueContext) interface{} {
	v, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		ev.addError(fmt.Errorf("invalid integer value: %s", ctx.GetText()))
	}
	return v
}

// Visit a parse tree produced by ECLParser#maxvalue.
func (ev *eclVisitor) VisitMaxvalue(ctx *ecl.MaxvalueContext) interface{} {
	v, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		ev.addError(fmt.Errorf("invalid integer value: %s", ctx.GetText()))
	}
	return v
}

// subExpressionConstraint = [constraintOperator ws] [memberOf ws] (eclFocusConcept / "(" ws expressionConstraint ws ")")
func (ev *eclVisitor) VisitSubexpressionconstraint(ctx *ecl.SubexpressionconstraintContext) interface{} {
	result := new(snomed.ExpressionConstraint_Subexpression)
	if ctx.Constraintoperator() != nil {
		if constraint, ok := ev.Visit(ctx.Constraintoperator()).(snomed.ExpressionConstraint_Constraint); ok {
			fmt.Printf("setting constraint to %v\n", constraint)
			result.Constraint = constraint
		} else {
			ev.addError(fmt.Errorf("invalid constraint: %s", ctx.Constraintoperator().GetText()))
		}
	}
	if ctx.Memberof() != nil {
		result.MemberOf = true
	}
	if ctx.Eclfocusconcept() != nil {
		if fc, ok := ev.Visit(ctx.Eclfocusconcept()).(*snomed.ExpressionConstraint_FocusConcept); ok {
			result.Value = &snomed.ExpressionConstraint_Subexpression_FocusConcept{
				FocusConcept: fc,
			}
		} else {
			ev.addError(fmt.Errorf("invalid focus concept: %s", ctx.Eclfocusconcept().GetText()))
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
func (ev *eclVisitor) VisitEclattributeset(ctx *ecl.EclattributesetContext) interface{} {
	result := new(snomed.ExpressionConstraint_AttributeSet)
	if ctx.Subattributeset() != nil {
		if subAttributeSet, ok := ev.Visit(ctx.Subattributeset()).(*snomed.ExpressionConstraint_SubattributeSet); ok {
			result.SubattributeSet = subAttributeSet
		} else {
			ev.addError(fmt.Errorf("invalid subattribute set: %s", ctx.Subattributeset().GetText()))
		}
	}
	if ctx.Conjunctionattributeset() != nil {
		if and, ok := ev.Visit(ctx.Conjunctionattributeset()).(*snomed.ExpressionConstraint_ConjunctionAttributeSet); ok {
			result.Logical = &snomed.ExpressionConstraint_AttributeSet_Conjunction{
				Conjunction: and,
			}
		} else {
			ev.addError(fmt.Errorf("invalid conjunction attribute set: %s", ctx.Conjunctionattributeset().GetText()))
		}
	}
	if ctx.Disjunctionattributeset() != nil {
		if or, ok := ev.Visit(ctx.Disjunctionattributeset()).(*snomed.ExpressionConstraint_DisjunctionAttributeSet); ok {
			result.Logical = &snomed.ExpressionConstraint_AttributeSet_Disjunction{
				Disjunction: or,
			}
		} else {
			ev.addError(fmt.Errorf("invalid disjunction attribute set: %s", ctx.Disjunctionattributeset().GetText()))
		}
	}
	return result
}

// subAttributeSet = eclAttribute / "(" ws eclAttributeSet ws ")"
func (ev *eclVisitor) VisitSubattributeset(ctx *ecl.SubattributesetContext) interface{} {
	result := new(snomed.ExpressionConstraint_SubattributeSet)
	if ctx.Eclattribute() != nil {
		if attribute, ok := ev.Visit(ctx.Eclattribute()).(*snomed.ExpressionConstraint_Attribute); ok {
			result.Value = &snomed.ExpressionConstraint_SubattributeSet_Attribute{
				Attribute: attribute,
			}
		} else {
			ev.addError(fmt.Errorf("invalid attribute: %s", ctx.Eclattribute().GetText()))
		}
	}
	if ctx.Eclattributeset() != nil {
		if attributeSet, ok := ev.Visit(ctx.Eclattributeset()).(*snomed.ExpressionConstraint_AttributeSet); ok {
			result.Value = &snomed.ExpressionConstraint_SubattributeSet_AttributeSet{
				AttributeSet: attributeSet,
			}
		} else {
			ev.addError(fmt.Errorf("invalid attribute set: %s", ctx.Eclattributeset().GetText()))
		}
	}
	return result
}

// conjunctionAttributeSet = 1*(ws conjunction ws subAttributeSet)
func (ev *eclVisitor) VisitConjunctionattributeset(ctx *ecl.ConjunctionattributesetContext) interface{} {
	result := new(snomed.ExpressionConstraint_ConjunctionAttributeSet)
	terms := make([]*snomed.ExpressionConstraint_SubattributeSet, 0)
	for _, subattributeset := range ctx.AllSubattributeset() {
		if sas, ok := ev.Visit(subattributeset).(*snomed.ExpressionConstraint_SubattributeSet); ok {
			terms = append(terms, sas)
		} else {
			ev.addError(fmt.Errorf("invalid subattribute set : %s", ctx.GetText()))
		}
	}
	result.And = terms
	return result
}

// disjunctionAttributeSet = 1*(ws disjunction ws subAttributeSet)
func (ev *eclVisitor) VisitDisjunctionattributeset(ctx *ecl.DisjunctionattributesetContext) interface{} {
	result := new(snomed.ExpressionConstraint_DisjunctionAttributeSet)
	terms := make([]*snomed.ExpressionConstraint_SubattributeSet, 0)
	for _, subattributeset := range ctx.AllSubattributeset() {
		if sas, ok := ev.Visit(subattributeset).(*snomed.ExpressionConstraint_SubattributeSet); ok {
			terms = append(terms, sas)
		} else {
			ev.addError(fmt.Errorf("invalid subattribute set : %s", ctx.GetText()))
		}
	}
	result.Or = terms
	return result

}

func (ev *eclVisitor) VisitConstraintoperator(ctx *ecl.ConstraintoperatorContext) interface{} {
	fmt.Printf("parsing constraint: '%s'\nctx:%v\n", ctx.GetText(), ctx)
	switch {
	case ctx.Ancestorof() != nil:
		fmt.Printf("ancestor of\n")
		return snomed.ExpressionConstraint_ANCESTOR_OF
	case ctx.Ancestororselfof() != nil:
		fmt.Printf("ancestor or self of\n")
		return snomed.ExpressionConstraint_ANCESTOR_OR_SELF_OF
	case ctx.Childof() != nil:
		fmt.Printf("child of: %v\n", ctx.Childof())
		return snomed.ExpressionConstraint_CHILD_OF
	case ctx.Descendantof() != nil:
		return snomed.ExpressionConstraint_DESCENDANT_OF
	case ctx.Descendantorselfof() != nil:
		return snomed.ExpressionConstraint_DESCENDANT_OR_SELF_OF
	case ctx.Parentof() != nil:
		return snomed.ExpressionConstraint_PARENT_OF
	default:
		ev.addError(fmt.Errorf("unknown constraint operator: %s", ctx.GetText()))
		return nil
	}
}

// eclFocusConcept = eclConceptReference / wildCard
func (ev *eclVisitor) VisitEclfocusconcept(ctx *ecl.EclfocusconceptContext) interface{} {
	fc := new(snomed.ExpressionConstraint_FocusConcept)
	if ctx.Eclconceptreference() != nil {
		if cr, ok := ev.Visit(ctx.Eclconceptreference()).(*snomed.ConceptReference); ok {
			fc.ConceptReference = cr
		} else {
			ev.addError(fmt.Errorf("invalid concept reference: %s", ctx.GetText()))
		}
	}
	if ctx.Wildcard() != nil {
		fc.Wildcard = true
	}
	return fc
}

// eclConceptReference = conceptId [ws "|" ws term ws "|"]
func (ev *eclVisitor) VisitEclconceptreference(ctx *ecl.EclconceptreferenceContext) interface{} {
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
func (ev *eclVisitor) VisitConceptid(ctx *ecl.ConceptidContext) interface{} {
	return ev.Visit(ctx.Sctid())
}

// returns string
func (ev *eclVisitor) VisitTerm(ctx *ecl.TermContext) interface{} {
	return ctx.GetText()
}

// TODO: check escaping rules - do we need to parse - probably?
// ABNF: stringValue = 1*(anyNonEscapedChar / escapedChar)
func (ev *eclVisitor) VisitStringvalue(ctx *ecl.StringvalueContext) interface{} {
	return ctx.GetText()
}

// numericValue = ["-"/"+"] (decimalValue / integerValue)
// returns either an int64 or float64
func (ev *eclVisitor) VisitNumericvalue(ctx *ecl.NumericvalueContext) interface{} {
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
func (ev *eclVisitor) VisitIntegervalue(ctx *ecl.IntegervalueContext) interface{} {
	v, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		ev.addError(fmt.Errorf("invalid integer value: %s", ctx.GetText()))
	}
	return v
}

// returns float64
func (ev *eclVisitor) VisitDecimalvalue(ctx *ecl.DecimalvalueContext) interface{} {
	v, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		ev.addError(fmt.Errorf("invalid double value: %s", ctx.GetText()))
	}
	return v
}

// returns a SNOMED CT identifier (int64)
func (ev *eclVisitor) VisitSctid(ctx *ecl.SctidContext) interface{} {
	sct, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		ev.addError(err)
	}
	return sct
}
