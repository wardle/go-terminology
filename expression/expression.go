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
// java -jar ~/Downloads/antlr-4.7.1-complete.jar -Dlanguage=Go -package ecl -o ecl ECL.g4 -visitor
// java -jar ~/Downloads/antlr-4.7.1-complete.jar -Dlanguage=Go -package cg -o cg CG.g4 -visitor
// The compositional grammar (CG) is from https://confluence.ihtsdotools.org/display/DOCSCG/5.1+Normative+Specification
// The expression contraint grammar (ECL) is from https://confluence.ihtsdotools.org/pages/viewpage.action?pageId=28739405
package expression

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/wardle/go-terminology/expression/cg"
	"github.com/wardle/go-terminology/snomed"
	"strconv"
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

// Normalize expands an expression into a normal-form, which makes it
// more readily computable. This essentially simplifies all terms as much as possible
// taking any complex compound single-form SNOMED codes and building the equivalent expression.
// Such an expression can then be used to determine equivalence or analytics.
func Normalize(e *snomed.Expression) *snomed.Expression {
	panic("Not implemented")
}

// ParseExpression parses a SNOMED expression
func ParseExpression(s string) (*snomed.Expression, error) {
	l := new(parserListener)
	is := antlr.NewInputStream(s)
	lexer := cg.NewCGLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := cg.NewCGParser(stream)
	antlr.ParseTreeWalkerDefault.Walk(l, p.Expression())
	return l.expression, l.err
}

// parserListener is an internal ANTLR listener.
// The generated code from ANTLR is quirky and not idiomatic and I am
// probably using the listener incorrectly. It would be better to use the visitor
// pattern, but that appears incomplete at of Antlr 4.7.1.
type parserListener struct {
	cg.BaseCGListener
	expression *snomed.Expression
	err        error
}

func (pl *parserListener) EnterEquivalentto(c *cg.EquivalenttoContext) {
	pl.expression.DefinitionStatus = snomed.Expression_EQUIVALENT_TO
}

func (pl *parserListener) EnterSubtypeof(c *cg.SubtypeofContext) {
	pl.expression.DefinitionStatus = snomed.Expression_SUBTYPE_OF
}

func (pl *parserListener) EnterExpression(ctx *cg.ExpressionContext) {
	if pl.expression == nil {
		pl.expression = new(snomed.Expression)
		se, ok := ctx.Subexpression().(*cg.SubexpressionContext)
		if !ok {
			pl.err = errors.New("No valid subexpression identified")
			return
		}
		clause, err := parseSubexpression(se)
		if err != nil {
			pl.err = err
			return
		}
		pl.expression.Clause = clause
		return
	}
}

func parseSubexpression(ctx *cg.SubexpressionContext) (*snomed.Expression_Clause, error) {
	clause := new(snomed.Expression_Clause)
	focus, err := parseFocusConcepts(ctx.Focusconcept().(*cg.FocusconceptContext))
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
		concept, err := parseConceptReference(cr.(*cg.ConceptreferenceContext))
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
		refinements, err := parseRefinements(ag.Attributeset().(*cg.AttributesetContext))
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
	ref.Term = ctx.Term().GetText()
	return
}
