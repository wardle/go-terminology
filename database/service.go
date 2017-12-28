package database

import (
	"fmt"

	"bitbucket.org/wardle/go-snomed/snomed"
	"golang.org/x/text/language"
)

// Service represents an opaque abstract SNOMED-CT persistence service.
// TODO: MW: tidy naming
type Service interface {
	GetConcept(conceptID int) (*snomed.Concept, error)
	GetConcepts(conceptIDs ...int) ([]*snomed.Concept, error)
	GetDescriptions(concept *snomed.Concept) ([]*snomed.Description, error)
	GetParentRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error)
	GetChildRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error)
	GetRecursiveChildrenIds(concept *snomed.Concept) ([]int, error)
	Close() error
}

// Snomed encapsulates a concrete persistent service and extends it by providing
// semantic inference and a useful, practical SNOMED-CT API which uses and encapsulates the
// underlying persistence service.
type Snomed struct {
	Service
	Language language.Tag
}

// GetPreferredDescription returns the preferred description for this concept in the default language for this service.
func (ds *Snomed) GetPreferredDescription(concept *snomed.Concept) (*snomed.Description, error) {
	return ds.GetPreferredDescriptionForLanguages(concept, []language.Tag{ds.Language})
}

// GetPreferredDescriptionForLanguages returns the preferred description for this concept in the languages specified
func (ds *Snomed) GetPreferredDescriptionForLanguages(concept *snomed.Concept, languages []language.Tag) (*snomed.Description, error) {
	preferred, err := ds.GetPreferredDescriptions(concept)
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
func (ds *Snomed) GetPreferredDescriptions(concept *snomed.Concept) ([]*snomed.Description, error) {
	descriptions, err := ds.GetDescriptions(concept)
	if err != nil {
		return nil, err
	}
	preferred := make([]*snomed.Description, 0, len(descriptions))
	for _, description := range descriptions {
		if description.Type.IsPreferred() {
			preferred = append(preferred, description)
		}
	}
	return preferred, nil
}

// GetSiblings returns the siblings of this concept, ie: those who share the same parents
func (ds *Snomed) GetSiblings(concept *snomed.Concept) ([]*snomed.Concept, error) {
	parents, err := ds.GetParents(concept)
	if err != nil {
		return nil, err
	}
	siblings := make([]*snomed.Concept, 0, 10)
	for _, parent := range parents {
		children, err := ds.GetChildren(parent)
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
func (ds *Snomed) GetAllParents(concept *snomed.Concept) ([]*snomed.Concept, error) {
	return ds.GetConcepts(concept.Parents...)
}

// GetParents returns the direct IS-A relations of the specified concept.
func (ds *Snomed) GetParents(concept *snomed.Concept) ([]*snomed.Concept, error) {
	return ds.GetParentsOfKind(concept, snomed.IsA)
}

// GetParentsOfKind returns the relations of the specified kinds (types) for the specified concept
func (ds *Snomed) GetParentsOfKind(concept *snomed.Concept, kinds ...snomed.Identifier) ([]*snomed.Concept, error) {
	relations, err := ds.GetParentRelationships(concept)
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
	return ds.GetConcepts(conceptIDs...)
}

// GetChildren returns the direct IS-A relations of the specified concept.
func (ds *Snomed) GetChildren(concept *snomed.Concept) ([]*snomed.Concept, error) {
	return ds.GetChildrenOfKind(concept, snomed.IsA)
}

// GetChildrenOfKind returns the relations of the specified kind (type) of the specified concept.
func (ds *Snomed) GetChildrenOfKind(concept *snomed.Concept, kind snomed.Identifier) ([]*snomed.Concept, error) {
	relations, err := ds.GetChildRelationships(concept)
	if err != nil {
		return nil, err
	}
	conceptIDs := make([]int, 0, len(relations))
	for _, relation := range relations {
		if relation.Type == kind {
			conceptIDs = append(conceptIDs, int(relation.Source))
		}
	}
	return ds.GetConcepts(conceptIDs...)
}

// FetchRecursiveChildren fetches all children of the given concept recursively.
// Use with caution with concepts at high levels of the hierarchy.
func (ds *Snomed) FetchRecursiveChildren(concept *snomed.Concept) ([]*snomed.Concept, error) {
	children, err := ds.GetRecursiveChildrenIds(concept)
	if err != nil {
		return nil, err
	}
	return ds.GetConcepts(children...)
}

// ConceptsForRelationship returns the concepts represented within a relationship
func (ds *Snomed) ConceptsForRelationship(rel *snomed.Relationship) (source *snomed.Concept, kind *snomed.Concept, target *snomed.Concept, err error) {
	concepts, err := ds.GetConcepts(int(rel.Source), int(rel.Type), int(rel.Target))
	if err != nil {
		return nil, nil, nil, err
	}
	return concepts[0], concepts[1], concepts[2], nil
}

// PathsToRoot returns the different possible paths to the root SNOMED-CT concept from this one.
func (ds *Snomed) PathsToRoot(concept *snomed.Concept) ([][]*snomed.Concept, error) {
	parents, err := ds.GetParents(concept)
	if err != nil {
		return nil, err
	}
	results := make([][]*snomed.Concept, 0, len(parents))
	if len(parents) == 0 {
		results = append(results, []*snomed.Concept{concept})
	}
	for _, parent := range parents {
		parentResults, err := ds.PathsToRoot(parent)
		if err != nil {
			return nil, err
		}
		for _, parentResult := range parentResults {
			r := append([]*snomed.Concept{concept}, parentResult...) // prepend current concept
			results = append(results, r)
		}
	}
	return results, nil
}

func debugPaths(paths [][]*snomed.Concept) {
	for i, path := range paths {
		fmt.Printf("Path %d: ", i)
		debugPath(path)
	}
}

func debugPath(path []*snomed.Concept) {
	for _, concept := range path {
		fmt.Printf("%d-", concept.ConceptID)
	}
	fmt.Print("\n")
}

// Genericise finds the best generic match for the given concept
// The "best" is chosen as the closest match to the specified concept and so
// if there are generic concepts which relate to one another, it will be the
// most specific (closest) match to the concept.
func (ds *Snomed) Genericise(concept *snomed.Concept, generics map[snomed.Identifier]*snomed.Concept) (*snomed.Concept, bool) {
	paths, err := ds.PathsToRoot(concept)
	if err != nil {
		return nil, false
	}
	var bestPath []*snomed.Concept
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
func (ds *Snomed) GenericiseToRoot(concept *snomed.Concept, root snomed.Identifier) (*snomed.Concept, error) {
	paths, err := ds.PathsToRoot(concept)
	if err != nil {
		return nil, err
	}
	var bestPath []*snomed.Concept
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
