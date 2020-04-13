package expression

import (
	"fmt"

	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"google.golang.org/protobuf/proto"
)

// Normalizer handles normalization of SNOMED CT expressions in which
// we expand an expression into a normal-form, making it more readily
// computable. This essentially simplifies all terms as much as possible
// taking any complex compound single-form SNOMED codes and building the equivalent expression.
// Such an expression can then be used to determine equivalence or analytics.
// See https://confluence.ihtsdotools.org/display/DOCTSG/12.3.3+Building+Long+and+Short+Normal+Forms
// and https://confluence.ihtsdotools.org/display/DOCTSG/12.4+Transforming+Expressions+to+Normal+Forms
//
// Any SNOMED CT expression can be transformed to its normal form by replacing each reference to a fully
// defined concept with a nested expression representing the definition of that concept. Transformation
// rules then resolve redundancies, which may arise from expanding fully defined concepts, by removing
// less specific attribute values.
//
// The steps are:
// 1. Separate Information Model Context  - https://confluence.ihtsdotools.org/display/DOCTSG/12.4.1+Separate+Information+Model+Context
// 2. Normalising the expression - https://confluence.ihtsdotools.org/display/DOCTSG/12.4.2+Normalize+Expression
// 3.
//
type Normalizer struct {
	svc *terminology.Svc
}

// NewNormalizer creates a new normalizer for the given expression
func NewNormalizer(svc *terminology.Svc) *Normalizer {
	return &Normalizer{
		svc: svc,
	}
}

// Normalize normalizes the specified expression
func (n *Normalizer) Normalize(e *snomed.Expression) (*snomed.Expression, error) {
	panic("not implemented")
}

// normalizedFocusConcepts() returns normalized expressions for each focus concept
func (n *Normalizer) normalizedFocusConcepts(e *snomed.Expression) ([]*snomed.Expression, error) {
	fcs := e.GetClause().GetFocusConcepts()
	exps := make([]*snomed.Expression, len(fcs))
	for i, fc := range fcs {
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

// mergeRefinements attempts to merge the specified groups, returning success or failure
// together with the newly merged group if this has been possible
// This follows the rules from https://confluence.ihtsdotools.org/display/DOCTSG/12.4.10+Merging+Groups
// Firstly, at least one attribute in one of the groups is named matched by an attribute in other group
// Secondly, for each name-matched pair, the value should be identical or subsume the other
// If so, the two groups are merged.
// Note: this does *not* remove duplicates after merge.
//
// For some reason, we don't merge two groups with unrelated attibutes; this doesn't seem intuitive to me...
func (n *Normalizer) mergeRefinements(
	r1 []*snomed.Expression_Refinement,
	r2 []*snomed.Expression_Refinement) (bool, []*snomed.Expression_Refinement, error) {
	nameMatched := 0  // number of attributes in which names match
	valueMatched := 0 // equals or subsumes
	for _, r1r := range r1 {
		for _, r2r := range r2 {
			if r1r.RefinementConcept.ConceptId == r2r.RefinementConcept.ConceptId {
				nameMatched++
				if proto.Equal(r1r, r2r) {
					valueMatched++
					continue
				}
				if v1, ok := r1r.GetValue().(*snomed.Expression_Refinement_ConceptValue); ok {
					if v2, ok := r2r.GetValue().(*snomed.Expression_Refinement_ConceptValue); ok {
						if v1.ConceptValue.ConceptId == v2.ConceptValue.ConceptId {
							valueMatched++
							continue
						}
						c1, err := n.svc.Concept(v1.ConceptValue.ConceptId) // TODO: fetching the concepts is redundant
						if err != nil {                                     // TODO: change API so that this is unnecessary
							return false, nil, err
						}
						c2, err := n.svc.Concept(v2.ConceptValue.ConceptId)
						if err != nil {
							return false, nil, err
						}
						if n.svc.IsA(c1, c2.Id) || n.svc.IsA(c2, c1.Id) {
							valueMatched++
							continue
						}
					}
				}
			}
		}
	}
	if nameMatched == 0 || nameMatched > valueMatched {
		return false, nil, nil
	}
	// we've name matched at least one, and all name matched are also value matched, so we can merge
	result := make([]*snomed.Expression_Refinement, 0)
	for _, r := range r1 {
		result = append(result, r)
	}
	for _, r := range r2 {
		result = append(result, r)
	}
	return true, result, nil
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
// The merging of attributes is potentially difficult, particularly if there are attributes that subsume
// other attributes. Some of the consideration are mentioned here:
// https://confluence.ihtsdotools.org/display/DOCTSG/12.4.9+Attribute+Names+and+Attribute+Hierarchies
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
