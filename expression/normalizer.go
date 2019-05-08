package expression

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
)

// Normalizer handles normalization of SNOMED CT expressions
// Normalization expands an expression into a normal-form, which makes it
// more readily computable. This essentially simplifies all terms as much as possible
// taking any complex compound single-form SNOMED codes and building the equivalent expression.
// Such an expression can then be used to determine equivalence or analytics.
// See https://confluence.ihtsdotools.org/display/DOCTSG/12.3.3+Building+Long+and+Short+Normal+Forms
// and https://confluence.ihtsdotools.org/display/DOCTSG/12.4+Transforming+Expressions+to+Normal+Forms
type Normalizer struct {
	svc *terminology.Svc
	e   *snomed.Expression
}

// NewNormalizer creates a new normalizer for the given expression
func NewNormalizer(svc *terminology.Svc, e *snomed.Expression) *Normalizer {
	e2 := proto.Clone(e)
	return &Normalizer{
		svc: svc,
		e:   e2.(*snomed.Expression),
	}
}

// getFocusConcepts() extracts the focus concepts from expression
func (n *Normalizer) getFocusConcepts() []*snomed.ConceptReference {
	return n.e.GetClause().GetFocusConcepts()
}

// normalizedFocusConcepts() returns normalized expressions for each focus concept
func (n *Normalizer) normalizedFocusConcepts() ([]*snomed.Expression, error) {
	fcs := n.getFocusConcepts()
	exps := make([]*snomed.Expression, len(fcs))
	for i, fc := range n.getFocusConcepts() {
		c, err := n.svc.Concept(fc.ConceptId)
		if err != nil {
			return nil, err
		}
		exps[i], err = NormalizeConcept(n.svc, c)
		if err != nil {
			return nil, err
		}
	}
	return exps, nil
}

// getRefinements() returns the refinements for the expression
// TODO: handle refinement groups as well
func (n *Normalizer) getRefinements() []*snomed.Expression_Refinement {
	return n.e.GetClause().GetRefinements()
}

// Normalize expands an expression into a normal-form, which makes it
// more readily computable. This essentially simplifies all terms as much as possible
// taking any complex compound single-form SNOMED codes and building the equivalent expression.
// Such an expression can then be used to determine equivalence or analytics.
// See https://confluence.ihtsdotools.org/display/DOCTSG/12.3.3+Building+Long+and+Short+Normal+Forms
// and https://confluence.ihtsdotools.org/display/DOCTSG/12.4+Transforming+Expressions+to+Normal+Forms
func Normalize(svc *terminology.Svc, e *snomed.Expression) (*snomed.Expression, error) {
	clause, err := normalizeClause(svc, e.GetClause())
	if err != nil {
		return nil, err
	}
	if clause == e.GetClause() { // if clause couldn't be normalised more, just return original expression
		return e, nil
	}
	exp := new(snomed.Expression)
	exp.Clause = clause
	exp.DefinitionStatus = e.DefinitionStatus
	return exp, nil
}

func normalizeClause(svc *terminology.Svc, clause *snomed.Expression_Clause) (*snomed.Expression_Clause, error) {
	conceptIds := make([]int64, 0)
	for _, c := range clause.GetFocusConcepts() {
		conceptIds = append(conceptIds, c.ConceptId)
	}
	concepts, err := svc.Concepts(conceptIds...)
	if err != nil {
		return nil, err
	}
	exps := make([]*snomed.Expression, 0)
	for _, concept := range concepts {
		e, err := NormalizeConcept(svc, concept)
		if err != nil {
			return nil, err
		}
		exps = append(exps, e)
	}
	e, err := mergeExpressions(svc, exps)
	if err != nil {
		return nil, err
	}

	return e.GetClause(), nil
}

// mergeExpressions merges the expressions into a single expression
func mergeExpressions(svc *terminology.Svc, exps []*snomed.Expression) (*snomed.Expression, error) {
	panic("not implemented")
}

// NormalizeConcept turns a single concept into its primitive components
func NormalizeConcept(svc *terminology.Svc, c *snomed.Concept) (*snomed.Expression, error) {
	primitive, err := svc.Primitive(c)
	if err != nil {
		return nil, err
	}
	exp := new(snomed.Expression)
	exp.Clause = new(snomed.Expression_Clause)
	focus := []*snomed.ConceptReference{{ConceptId: primitive.Id}}
	exp.Clause.FocusConcepts = focus
	// and now let's add the primitive versions of our defining relationships
	rels, err := svc.ParentRelationships(c.Id)
	if err != nil {
		return nil, err
	}
	attrs := make([]*snomed.Expression_Refinement, 0)
	unique := make(map[string]struct{}) // ensure only unique attributes recorded.
	for _, rel := range rels {
		if rel.GetTypeId() != snomed.IsA && rel.Active && rel.IsDefiningRelationship() {
			fmt.Printf("rel: %v\n", rel)
			typeID := rel.GetTypeId()
			childID := rel.GetDestinationId()
			relType, err := svc.Concept(typeID)
			if err != nil {
				return nil, err
			}
			child, err := svc.Concept(childID)
			if err != nil {
				return nil, err
			}
			primitiveChild, err := svc.Primitive(child) // get primitive of relationship, so potential for duplicates
			if err != nil {
				return nil, err
			}
			key := fmt.Sprintf("%d-%d", typeID, primitiveChild.Id)
			if _, done := unique[key]; !done {
				attr := new(snomed.Expression_Refinement)
				attr.RefinementConcept = &snomed.ConceptReference{ConceptId: relType.Id}
				attr.Value = &snomed.Expression_Refinement_ConceptValue{ConceptValue: &snomed.ConceptReference{ConceptId: primitiveChild.Id}}
				attrs = append(attrs, attr)
				unique[key] = struct{}{}
			}
		}
	}
	exp.Clause.Refinements = attrs
	return exp, nil
}
