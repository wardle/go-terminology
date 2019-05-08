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
// java -jar ~/Downloads/antlr-4.7.1-complete.jar -Dlanguage=Go -package ecl -o ecl ECL.g4
// java -jar ~/Downloads/antlr-4.7.1-complete.jar -Dlanguage=Go -package cg -o cg CG.g4
// The compositional grammar (CG) is from https://confluence.ihtsdotools.org/display/DOCSCG/5.1+Normative+Specification
// The expression constraint grammar (ECL) is from https://confluence.ihtsdotools.org/pages/viewpage.action?pageId=28739405
package expression

import (
	"errors"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/wardle/go-terminology/expression/cg"
	"github.com/wardle/go-terminology/expression/ecl"
	"github.com/wardle/go-terminology/snomed"
	"strconv"
	"strings"
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

// ParseExpression parses a SNOMED expression
func ParseExpression(s string) (*snomed.Expression, error) {
	l := new(cgListener)
	is := antlr.NewInputStream(s)
	lexer := cg.NewCGLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := cg.NewCGParser(stream)
	antlr.ParseTreeWalkerDefault.Walk(l, p.Expression())
	return l.expression, l.err
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

func ParseExpressionConstraint(s string) error {
	l := new(eclListener)
	is := antlr.NewInputStream(s)
	lexer := cg.NewCGLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := ecl.NewECLParser(stream)
	antlr.ParseTreeWalkerDefault.Walk(l, p.Expressionconstraint())
	return l.err
}

type eclListener struct {
	ecl.BaseECLListener
	err error
}

func (el *eclListener) EnterExpressionconstraint(ctx *ecl.ExpressionconstraintContext) {
	if ctx.Subexpressionconstraint() != nil {
		fmt.Printf("subexpression: %v\n", ctx.Subexpressionconstraint())
	}
	if ctx.Compoundexpressionconstraint() != nil {
		fmt.Printf("compound expression: %v\n", ctx.Compoundexpressionconstraint())
	}
	if ctx.Dottedexpressionconstraint() != nil {
		fmt.Printf("dotted : %v\n", ctx.Dottedexpressionconstraint())
	}
	if ctx.Refinedexpressionconstraint() != nil {
		fmt.Printf("refined: %v\n", ctx.Refinedexpressionconstraint())
		subexpression := ctx.Refinedexpressionconstraint().(*ecl.RefinedexpressionconstraintContext).Subexpressionconstraint()
		refinement := ctx.Refinedexpressionconstraint().(*ecl.RefinedexpressionconstraintContext).Eclrefinement()
		fmt.Printf("|-- subexp: %v  refinement: %v", subexpression, refinement)
		fmt.Printf("refinement: %s", refinement.GetText())
	}
	fmt.Printf("%v\n", ctx)
}
