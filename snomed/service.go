package snomed

import (
	"fmt"
	"golang.org/x/text/language"
)

// Service represents an opaque abstract SNOMED-CT persistence service.
// TODO: MW: tidy naming
type Service interface {
	GetConcept(conceptID int) (*Concept, error)
	GetConcepts(conceptIDs ...int) ([]*Concept, error)
	GetDescriptions(concept *Concept) ([]*Description, error)
	GetParentRelationships(concept *Concept) ([]*Relationship, error)
	GetChildRelationships(concept *Concept) ([]*Relationship, error)
	GetRecursiveChildrenIds(concept *Concept) ([]int, error)
	Close() error
}

// Snomed encapsulates a concrete persistent service and extensct it by providing
// semantic inference and a useful, practical SNOMED-CT API which uses and encapsulates the
// underlying persistence service.
type Snomed struct {
	Service
	Language language.Tag
}

// GetPreferredDescription returns the preferred description for this concept in the default language for this service.
func (sct *Snomed) GetPreferredDescription(concept *Concept) (*Description, error) {
	return sct.GetPreferredDescriptionForLanguages(concept, []language.Tag{sct.Language})
}

// GetPreferredDescriptionForLanguages returns the preferred description for this concept in the languages specified
func (sct *Snomed) GetPreferredDescriptionForLanguages(concept *Concept, languages []language.Tag) (*Description, error) {
	preferred, err := sct.GetPreferredDescriptions(concept)
	if err != nil {
		return nil, err
	}
	matcher := language.NewMatcher(languages)
	tags := make([]language.Tag, 0, len(preferred))
	for _, d := range preferred {
		tags = append(tags, d.LanguageCode)
	}
	_, index, _ := matcher.Match(tags...)
	return preferred[index], nil
}

// GetPreferredDescriptions returns the preferred descriptions for the given concept
func (sct *Snomed) GetPreferredDescriptions(concept *Concept) ([]*Description, error) {
	descriptions, err := sct.GetDescriptions(concept)
	if err != nil {
		return nil, err
	}
	preferred := make([]*Description, 0, len(descriptions))
	for _, description := range descriptions {
		if description.Type.IsPreferred() {
			preferred = append(preferred, description)
		}
	}
	return preferred, nil
}

// GetSiblings returns the siblings of this concept, ie: those who share the same parents
func (sct *Snomed) GetSiblings(concept *Concept) ([]*Concept, error) {
	parents, err := sct.GetParents(concept)
	if err != nil {
		return nil, err
	}
	siblings := make([]*Concept, 0, 10)
	for _, parent := range parents {
		children, err := sct.GetChildren(parent)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			if child.ConceptID != concept.ConceptID {
				siblings = append(siblings, child)
			}
		}
	}
	return siblings, nil
}

// GetAllParents returns all of the parents (recursively) for a given concept
func (sct *Snomed) GetAllParents(concept *Concept) ([]*Concept, error) {
	return sct.GetConcepts(concept.Parents...)
}

// GetParents returns the direct IS-A relations of the specified concept.
func (sct *Snomed) GetParents(concept *Concept) ([]*Concept, error) {
	return sct.GetParentsOfKind(concept, IsA)
}

// GetParentsOfKind returns the relations of the specified kinsct (types) for the specified concept
func (sct *Snomed) GetParentsOfKind(concept *Concept, kinds ...Identifier) ([]*Concept, error) {
	relations, err := sct.GetParentRelationships(concept)
	if err != nil {
		return nil, err
	}
	conceptIDs := make([]int, 0, len(relations))
	for _, relation := range relations {
		for _, kind := range kinds {
			if relation.Type == kind {
				conceptIDs = append(conceptIDs, int(relation.Target))
			}
		}
	}
	return sct.GetConcepts(conceptIDs...)
}

// GetChildren returns the direct IS-A relations of the specified concept.
func (sct *Snomed) GetChildren(concept *Concept) ([]*Concept, error) {
	return sct.GetChildrenOfKind(concept, IsA)
}

// GetChildrenOfKind returns the relations of the specified kind (type) of the specified concept.
func (sct *Snomed) GetChildrenOfKind(concept *Concept, kind Identifier) ([]*Concept, error) {
	relations, err := sct.GetChildRelationships(concept)
	if err != nil {
		return nil, err
	}
	conceptIDs := make([]int, 0, len(relations))
	for _, relation := range relations {
		if relation.Type == kind {
			conceptIDs = append(conceptIDs, int(relation.Source))
		}
	}
	return sct.GetConcepts(conceptIDs...)
}

// GetRecursiveChildren fetches all children of the given concept recursively.
// Use with caution with concepts at high levels of the hierarchy.
func (sct *Snomed) GetRecursiveChildren(concept *Concept) ([]*Concept, error) {
	children, err := sct.GetRecursiveChildrenIds(concept)
	if err != nil {
		return nil, err
	}
	return sct.GetConcepts(children...)
}

// ConceptsForRelationship returns the concepts represented within a relationship
func (sct *Snomed) ConceptsForRelationship(rel *Relationship) (source *Concept, kind *Concept, target *Concept, err error) {
	concepts, err := sct.GetConcepts(int(rel.Source), int(rel.Type), int(rel.Target))
	if err != nil {
		return nil, nil, nil, err
	}
	return concepts[0], concepts[1], concepts[2], nil
}

// PathsToRoot returns the different possible paths to the root SNOMED-CT concept from this one.
func (sct *Snomed) PathsToRoot(concept *Concept) ([][]*Concept, error) {
	parents, err := sct.GetParents(concept)
	if err != nil {
		return nil, err
	}
	results := make([][]*Concept, 0, len(parents))
	if len(parents) == 0 {
		results = append(results, []*Concept{concept})
	}
	for _, parent := range parents {
		parentResults, err := sct.PathsToRoot(parent)
		if err != nil {
			return nil, err
		}
		for _, parentResult := range parentResults {
			r := append([]*Concept{concept}, parentResult...) // prepend current concept
			results = append(results, r)
		}
	}
	return results, nil
}

func debugPaths(paths [][]*Concept) {
	for i, path := range paths {
		fmt.Printf("Path %d: ", i)
		debugPath(path)
	}
}

func debugPath(path []*Concept) {
	for _, concept := range path {
		fmt.Printf("%d-", concept.ConceptID)
	}
	fmt.Print("\n")
}

// Genericise finsct the best generic match for the given concept
// The "best" is chosen as the closest match to the specified concept and so
// if there are generic concepts which relate to one another, it will be the
// most specific (closest) match to the concept.
func (sct *Snomed) Genericise(concept *Concept, generics map[Identifier]*Concept) (*Concept, bool) {
	paths, err := sct.PathsToRoot(concept)
	if err != nil {
		return nil, false
	}
	var bestPath []*Concept
	bestPos := -1
	for _, path := range paths {
		for i, concept := range path {
			if generics[concept.ConceptID] != nil {
				if i > 0 && (bestPos == -1 || bestPos > i) {
					bestPos = i
					bestPath = path
				}
			}
		}
	}
	if bestPos == -1 {
		return nil, false
	}
	return bestPath[bestPos], true
}

// GenericiseToRoot walks the SNOMED-CT IS-A hierarchy to find the most general concept
// beneath the specified root.
// This finds the shortest path from the concept to the specified root and then
// returns one concept *down* from that root.
func (sct *Snomed) GenericiseToRoot(concept *Concept, root Identifier) (*Concept, error) {
	paths, err := sct.PathsToRoot(concept)
	if err != nil {
		return nil, err
	}
	var bestPath []*Concept
	bestPos := -1
	for _, path := range paths {
		for i, concept := range path {
			if concept.ConceptID == root {
				if i > 0 && (bestPos == -1 || bestPos > i) {
					bestPos = i
					bestPath = path
				}
			}
		}
	}
	if bestPos == -1 {
		return nil, fmt.Errorf("Root concept of %d not found for concept %d", root, concept.ConceptID)
	}
	return bestPath[bestPos-1], nil
}
