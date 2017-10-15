package snomed

// Service is an abstract way of handling the SNOMED-CT model
type Service interface {
	GetConcept(conceptID int) (*Concept, error)                             // Find a concept by identifier
	GetParents(conceptID int) []*Concept                                    // Get parents of this concept (defined by IS-A relationships)
	GetChildren(conceptID int) []*Concept                                   // Get children of this concept (defined by IS-A relationships)
	GetAllParents(concept *Concept) []*Concept                              // Get all parents (recursively) of this concept
	GetAllChildren(concept *Concept) []*Concept                             // Get all children (recursively) of this concept
	GetDescriptions(concept *Concept) []*Description                        // Get all descriptions (synonyms) for this concept
	GetPreferredDescription(concept *Concept, language string) *Description // TODO: change to language code
	GetSourceRelationships(concept *Concept) []*Relationship                // return relationships for which this concept is a target
	GetTargetRelationships(concept *Concept) []*Relationship                // return relationships for which this concept is a source
}
