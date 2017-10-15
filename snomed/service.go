package snomed

// Service is an abstract way of handling the SNOMED-CT model
type Service interface {
	GetConcept(conceptID int) (*Concept, error)
	GetParents(conceptID int) []*Concept
	GetChildren(conceptID int) []*Concept
	GetRecursiveParents(concept *Concept) []*Concept
	GetRecursiveChildren(concept *Concept) []*Concept
	GetDescriptions(concept *Concept) []*Description
	GetPreferredDescription(concept *Concept, language string) *Description //TODO: change to language code
	//GetSourceRelationships(concept *Concept) []*snomed.Relationship                // return relationships for which this concept is a target
	//GetTargetRelationships(concept *Concept) []*snomed.Relationship
}
