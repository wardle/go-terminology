package snomed

// Relationship provides a relationship between two concepts of a type defined by a concept.
type Relationship struct {
	RelationshipID Identifier
	Source         *Concept
	Target         *Concept
	Type           *Concept
}
