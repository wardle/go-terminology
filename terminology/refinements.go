package terminology

import (
	"github.com/wardle/go-terminology/snomed"
	"golang.org/x/text/language"
)

// Refinements determines the appropriate refinements for an arbitrary concept
// It is quite easy to do, we find the relationships and additionally determine
// whether the concept's attributes exist in the lateralisable reference set.
// TODO: this would be better deprecated in favour of using only expressions
// that would mean normalising any concept into an expression and *then* deriving
// possible refinements for that expression, instead.
func (svc *Svc) Refinements(conceptID int64, limit int, tags []language.Tag) (*snomed.RefinementResponse, error) {
	c, err := svc.Concept(conceptID)
	if err != nil {
		return nil, err
	}
	rels, err := svc.ParentRelationships(c.Id)
	if err != nil {
		return nil, err
	}
	attrs := make([]*snomed.RefinementResponse_Refinement, 0)
	properties := make(map[int64]struct{})
	for _, rel := range rels {
		if rel.Active && rel.TypeId != snomed.IsA {
			if _, done := properties[rel.DestinationId]; done {
				continue
			}
			properties[rel.DestinationId] = struct{}{}
			cc, err := svc.Concepts(rel.TypeId, rel.DestinationId)
			if err != nil {
				return nil, err
			}
			attr := new(snomed.RefinementResponse_Refinement)
			attr.Attribute, err = svc.ConceptReference(cc[0].Id, tags)
			if err != nil {
				return nil, err
			}
			attr.RootValue, err = svc.ConceptReference(cc[1].Id, tags)
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, attr)
			if rel.TypeId == snomed.BodyStructure || rel.TypeId == snomed.ProcedureSiteDirect || rel.TypeId == snomed.FindingSite {
				if _, done := properties[snomed.Side]; !done {
					islat, err := svc.IsLateralisable(rel.DestinationId)
					if err != nil {
						return nil, err
					}
					if islat {
						lat := new(snomed.RefinementResponse_Refinement)
						ll, err := svc.Concepts(snomed.Laterality, snomed.Side)
						if err != nil {
							return nil, err
						}
						lat.Attribute, err = svc.ConceptReference(ll[0].Id, tags)
						if err != nil {
							return nil, err
						}
						lat.RootValue, err = svc.ConceptReference(ll[1].Id, tags)
						if err != nil {
							return nil, err
						}
						attrs = append(attrs, lat)
					}
				}
			}

		}
	}
	response := new(snomed.RefinementResponse)
	response.Concept = c
	response.Refinements = attrs
	return response, nil
}

// IsLateralisable finds out whether the specific concept is lateralisable
func (svc *Svc) IsLateralisable(id int64) (bool, error) {
	rsis, err := svc.ComponentFromReferenceSet(snomed.LateralisableReferenceSet, id)
	if err != nil {
		return false, err
	}
	for _, rsi := range rsis {
		if rsi.Active {
			return true, nil
		}
	}
	return false, nil
}
