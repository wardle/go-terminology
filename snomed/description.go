package snomed

// Description is a synonym for a concept.
type Description struct {
	DescriptionID int
	Term          string
	Language      string // TODO: use proper language codes
	Concept       *Concept
}
