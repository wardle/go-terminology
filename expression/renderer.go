package expression

import (
	"sort"
	"strconv"
	"strings"

	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
)

// Renderer renders a SNOMED CT expression as text such that it can be roundtripped
// back to an expression via parsing. This means that generated text meets the syntax
// in the compositional grammar.
type Renderer struct {
	svc         *terminology.Svc
	hideTerms   bool           // hide terms, default false. Should be true for canonical
	updateTerms bool           // update terms to preferred terms from live service, default false.
	sort        bool           // sort focus concepts, refinements, attributes and attribute groups as per canonical form
	tags        []language.Tag // preferred language for terms, used if updating
}

// NewDefaultRenderer returns a renderer with the default formatting options
func NewDefaultRenderer() *Renderer {
	return &Renderer{
		svc:         nil,
		hideTerms:   false,
		updateTerms: false,
		sort:        false,
	}
}

// NewCanonicalRenderer returns a renderer that formats expressions canonically.
// This adopts the rules outlined here: https://confluence.ihtsdotools.org/display/DOCTSG/12.4.29+Canonical+Representation
func NewCanonicalRenderer() *Renderer {
	return &Renderer{
		svc:         nil,
		hideTerms:   true,
		sort:        true,
		updateTerms: false,
	}
}

// NewUpdatingRenderer returns a renderer that updates terms according to the preferred synonyms
func NewUpdatingRenderer(svc *terminology.Svc, tags []language.Tag) *Renderer {
	return &Renderer{
		svc:         svc,
		hideTerms:   false,
		updateTerms: true,
		sort:        false,
		tags:        tags,
	}
}

// Render is a simple helper to render the specified expression using a default renderer.
func Render(exp *snomed.Expression) string {
	r, err := NewDefaultRenderer().Render(exp)
	if err != nil {
		panic(err)
	}
	return r
}

// Render renders a SNOMED CT expression according to the configured rendering rules
func (r *Renderer) Render(exp *snomed.Expression) (string, error) {
	var sb strings.Builder
	if err := r.renderExpression(&sb, exp); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (r *Renderer) renderConcept(cr *snomed.ConceptReference) (string, error) {
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(cr.ConceptId, 10))
	if r.hideTerms == false {
		var term string
		if r.svc != nil && r.updateTerms && r.tags != nil {
			d, err := r.svc.PreferredSynonym(cr.ConceptId, r.tags)
			if err != nil {
				return "", err
			}
			term = d.Term
		}
		if term == "" && cr.Term != "" {
			term = cr.Term
		}
		sb.WriteString("|")
		sb.WriteString(term)
		sb.WriteString("|")

	}
	return sb.String(), nil
}

func (r *Renderer) renderRefinement(refinement *snomed.Expression_Refinement) (string, error) {
	var sb strings.Builder
	concept, err := r.renderConcept(refinement.GetRefinementConcept())
	if err != nil {
		return "", err
	}
	sb.WriteString(concept)
	sb.WriteString("=")
	value := refinement.GetValue()
	if clauseValue, ok := value.(*snomed.Expression_Refinement_ClauseValue); ok {
		sb.WriteString("(")
		clause, err := r.renderClause(clauseValue.ClauseValue)
		if err != nil {
			return "", err
		}
		sb.WriteString(clause)
		sb.WriteString(")")
	}
	if conceptValue, ok := value.(*snomed.Expression_Refinement_ConceptValue); ok {
		c, err := r.renderConcept(conceptValue.ConceptValue)
		if err != nil {
			return "", err
		}
		sb.WriteString(c)
	}
	if doubleValue, ok := value.(*snomed.Expression_Refinement_DoubleValue); ok {
		sb.WriteString("#")
		sb.WriteString(strconv.FormatFloat(doubleValue.DoubleValue, 'e', -1, 64))
	}
	if intValue, ok := value.(*snomed.Expression_Refinement_IntValue); ok {
		sb.WriteString("#")
		sb.WriteString(strconv.FormatInt(intValue.IntValue, 10))
	}
	if stringValue, ok := value.(*snomed.Expression_Refinement_StringValue); ok {
		sb.WriteString(stringValue.StringValue)
	}
	return sb.String(), nil
}

func (r *Renderer) renderExpression(sb *strings.Builder, e *snomed.Expression) error {
	// deliberately omit equivalent-to, as this is default
	//if e.DefinitionStatus == Expression_EQUIVALENT_TO {
	//	sb.WriteString("===")
	//}
	if e.DefinitionStatus == snomed.Expression_SUBTYPE_OF {
		sb.WriteString("<<<")
	}
	clause, err := r.renderClause(e.GetClause())
	if err != nil {
		return err
	}
	sb.WriteString(clause)
	return nil
}

func (r *Renderer) renderRefinements(refinements []*snomed.Expression_Refinement) (string, error) {
	var err error
	var sb strings.Builder
	rr := make([]string, len(refinements))
	for i, refinement := range refinements {
		rr[i], err = r.renderRefinement(refinement)
		if err != nil {
			return "", err
		}
	}
	if r.sort {
		sort.Strings(rr)
	}
	sb.WriteString(strings.Join(rr, ","))
	return sb.String(), nil
}

func (r *Renderer) renderGroups(groups []*snomed.Expression_RefinementGroup) (string, error) {
	var err error
	var sb strings.Builder
	rg := make([]string, len(groups))
	for i, group := range groups {
		rg[i], err = r.renderGroup(group)
		if err != nil {
			return "", err
		}
	}
	if r.sort {
		sort.Strings(rg)
	}
	sb.WriteString(strings.Join(rg, ""))
	return sb.String(), nil
}

func (r *Renderer) renderGroup(group *snomed.Expression_RefinementGroup) (string, error) {
	var sb strings.Builder
	sb.WriteString("{")
	refinements, err := r.renderRefinements(group.GetRefinements())
	if err != nil {
		return "", err
	}
	sb.WriteString(refinements)
	sb.WriteString("}")
	return sb.String(), nil
}

func (r *Renderer) renderClause(clause *snomed.Expression_Clause) (string, error) {
	var sb strings.Builder
	// process focus concepts
	concepts := clause.GetFocusConcepts()
	rf := make([]string, len(concepts))
	var err error
	for i, concept := range concepts {
		rf[i], err = r.renderConcept(concept)
		if err != nil {
			return "", err
		}
	}
	if r.sort {
		sort.Strings(rf)
	}
	sb.WriteString(strings.Join(rf, "+"))
	refinements := clause.GetRefinements()
	groups := clause.GetRefinementGroups()
	if len(refinements) == 0 && len(groups) == 0 {
		return sb.String(), nil
	}
	sb.WriteString(":")
	rr, err := r.renderRefinements(refinements)
	if err != nil {
		return "", err
	}
	sb.WriteString(rr)
	if len(refinements) > 0 && len(groups) > 0 {
		sb.WriteString(",")
	}
	rg, err := r.renderGroups(groups)
	if err != nil {
		return "", err
	}
	sb.WriteString(rg)
	return sb.String(), nil
}
