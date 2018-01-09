package db

import (
	"bitbucket.org/wardle/go-snomed/rf2"
	"bitbucket.org/wardle/go-snomed/snomed"
	"fmt"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
)

// Store represents the backend opaque abstract SNOMED-CT persistence service.
type Store interface {
	GetConcept(conceptID int) (*rf2.Concept, error)
	GetConcepts(conceptIDs ...int) ([]*rf2.Concept, error)
	GetDescriptions(concept *rf2.Concept) ([]*rf2.Description, error)
	GetParentRelationships(concept *rf2.Concept) ([]*rf2.Relationship, error)
	GetChildRelationships(concept *rf2.Concept) ([]*rf2.Relationship, error)
	//GetRecursiveChildrenIds(concept *rf2.Concept) ([]int, error)
	Close() error
}

// Search represents an opaque abstract SNOMED-CT search service.
type Search interface {
	// Search executes a search request and returns description identifiers
	Search(search *SearchRequest) ([]int, error)
}

// SearchRequest is used to set the parameters on which to search
type SearchRequest struct {
	Terms             string // search terms
	Limit             int    // max number of results
	Modules           []int  // limit search to specific modules
	RecursiveParents  []int  // limit search to specific recursive parents
	DirectParents     []int  // limit search to specific direct parents
	OnlyActiveConcept bool   // limit search to only active concepts
}

// Status provides status information about the SNOMED-CT service.
type Status struct {
	concepts       int
	descriptions   int
	relationships  int
	refSets        int
	hasPrecomputed bool
	hasIndex       bool
}

// Snomed encapsulates concrete persistent and search services and extends it by providing
// semantic inference and a useful, practical SNOMED-CT API.
type Snomed struct {
	Store
	Search
	Language language.Tag
}

// NewService opens or creates a service at the specified location.
func NewService(path string, readOnly bool) (*Snomed, error) {
	err := os.MkdirAll(path, 0771)
	if err != nil {
		return nil, err
	}
	bolt, err := NewBoltService(filepath.Join(path, "bolt.db"), readOnly)
	if err != nil {
		return nil, err
	}
	bleve, err := NewBleveService(filepath.Join(path, "index.bleve"), readOnly)
	if err != nil {
		return nil, err
	}
	return &Snomed{Store: bolt, Search: bleve, Language: language.BritishEnglish}, nil
}

// Status returns useful status information
func (ds *Snomed) Status() Status {
	return Status{}
}

// IsA tests whether the given concept is a type of the specified
// This is a crude implementation which, probably, should be optimised or cached
// much like the old t_cached_parent_concepts table in the SQL version
func (ds *Snomed) IsA(concept *rf2.Concept, parent snomed.Identifier) bool {
	if concept.ID == parent {
		return true
	}
	parents, err := ds.GetAllParents(concept)
	if err != nil {
		return false
	}
	for _, p := range parents {
		if p.ID == parent {
			return true
		}
	}
	return false
}

// GetFullySpecifiedName returns the FSN (fully specified name) for the given concept
func (ds *Snomed) GetFullySpecifiedName(concept *rf2.Concept) (*rf2.Description, error) {
	descriptions, err := ds.GetDescriptions(concept)
	if err != nil {
		return nil, err
	}
	for _, d := range descriptions {
		if d.IsFullySpecifiedName() {
			return d, nil
		}
	}
	return nil, fmt.Errorf("no fsn found for concept %d", concept.ID)
}

// MustGetFullySpecifiedName returns the FSN for the given concept, or panics if there is an error or it is missing
func (ds *Snomed) MustGetFullySpecifiedName(concept *rf2.Concept) *rf2.Description {
	fsn, err := ds.GetFullySpecifiedName(concept)
	if err != nil {
		panic(err)
	}
	return fsn
}

// GetPreferredDescription returns the preferred description for this concept in the default language for this service.
func (ds *Snomed) GetPreferredDescription(concept *rf2.Concept) (*rf2.Description, error) {
	return ds.GetPreferredDescriptionForLanguages(concept, []language.Tag{ds.Language})
}

// GetPreferredDescriptionForLanguages returns the preferred description for this concept in the languages specified
// TODO(mw): this is now wrong as SNOMED-CT RF2 uses subsets to handle language preferences
// TODO(mw): implement new
func (ds *Snomed) GetPreferredDescriptionForLanguages(concept *rf2.Concept, languages []language.Tag) (*rf2.Description, error) {
	preferred, err := ds.GetPreferredDescriptions(concept)
	if err != nil {
		return nil, err
	}
	matcher := language.NewMatcher(languages)
	tags := make([]language.Tag, 0, len(preferred))
	for _, d := range preferred {
		tags = append(tags, d.LanguageTag())
	}
	_, index, _ := matcher.Match(tags...)
	return preferred[index], nil
}

// GetPreferredDescriptions returns the preferred descriptions for the given concept
func (ds *Snomed) GetPreferredDescriptions(concept *rf2.Concept) ([]*rf2.Description, error) {
	descriptions, err := ds.GetDescriptions(concept)
	if err != nil {
		return nil, err
	}
	preferred := make([]*rf2.Description, 0, len(descriptions))
	for _, description := range descriptions {
		if description.IsSynonym() {
			preferred = append(preferred, description)
		}
	}
	return preferred, nil
}

// GetSiblings returns the siblings of this concept, ie: those who share the same parents
func (ds *Snomed) GetSiblings(concept *rf2.Concept) ([]*rf2.Concept, error) {
	parents, err := ds.GetParents(concept)
	if err != nil {
		return nil, err
	}
	siblings := make([]*rf2.Concept, 0, 10)
	for _, parent := range parents {
		children, err := ds.GetChildren(parent)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			if child.ID != concept.ID {
				siblings = append(siblings, child)
			}
		}
	}
	return siblings, nil
}

// GetAllParents returns all of the parents (recursively) for a given concept
func (ds *Snomed) GetAllParents(concept *rf2.Concept) ([]*rf2.Concept, error) {
	parents := make(map[snomed.Identifier]bool)
	err := ds.getAllParents(concept, parents)
	if err != nil {
		return nil, err
	}
	keys := make([]int, len(parents))
	i := 0
	for k := range parents {
		keys[i] = int(k)
		i++
	}
	return ds.GetConcepts(keys...)
}

func (ds *Snomed) getAllParents(concept *rf2.Concept, parents map[snomed.Identifier]bool) error {
	ps, err := ds.GetParents(concept)
	if err != nil {
		return err
	}
	for _, p := range ps {
		parents[p.ID] = true
		ds.getAllParents(p, parents)
	}
	return nil
}

// GetParents returns the direct IS-A relations of the specified concept.
func (ds *Snomed) GetParents(concept *rf2.Concept) ([]*rf2.Concept, error) {
	return ds.GetParentsOfKind(concept, snomed.IsAConceptID)
}

// GetParentsOfKind returns the active relations of the specified kinds (types) for the specified concept
func (ds *Snomed) GetParentsOfKind(concept *rf2.Concept, kinds ...snomed.Identifier) ([]*rf2.Concept, error) {
	relations, err := ds.GetParentRelationships(concept)
	if err != nil {
		return nil, err
	}
	conceptIDs := make([]int, 0, len(relations))
	for _, relation := range relations {
		if relation.Active {
			for _, kind := range kinds {
				if relation.TypeID == kind {
					conceptIDs = append(conceptIDs, int(relation.DestinationID))
				}
			}
		}
	}
	return ds.GetConcepts(conceptIDs...)
}

// GetChildren returns the direct IS-A relations of the specified concept.
func (ds *Snomed) GetChildren(concept *rf2.Concept) ([]*rf2.Concept, error) {
	return ds.GetChildrenOfKind(concept, snomed.IsAConceptID)
}

// GetChildrenOfKind returns the relations of the specified kind (type) of the specified concept.
func (ds *Snomed) GetChildrenOfKind(concept *rf2.Concept, kind snomed.Identifier) ([]*rf2.Concept, error) {
	relations, err := ds.GetChildRelationships(concept)
	if err != nil {
		return nil, err
	}
	conceptIDs := make([]int, 0, len(relations))
	for _, relation := range relations {
		if relation.Active {
			if relation.TypeID == kind {
				conceptIDs = append(conceptIDs, int(relation.SourceID))
			}
		}
	}
	return ds.GetConcepts(conceptIDs...)
}

/*
// FetchRecursiveChildren fetches all children of the given concept recursively.
// Use with caution with concepts at high levels of the hierarchy.
func (ds *Snomed) FetchRecursiveChildren(concept *rf2.Concept) ([]*rf2.Concept, error) {
	children, err := ds.GetRecursiveChildrenIds(concept)
	if err != nil {
		return nil, err
	}
	return ds.GetConcepts(children...)
}
*/

// ConceptsForRelationship returns the concepts represented within a relationship
func (ds *Snomed) ConceptsForRelationship(rel *rf2.Relationship) (source *rf2.Concept, kind *rf2.Concept, target *rf2.Concept, err error) {
	concepts, err := ds.GetConcepts(int(rel.SourceID), int(rel.TypeID), int(rel.DestinationID))
	if err != nil {
		return nil, nil, nil, err
	}
	return concepts[0], concepts[1], concepts[2], nil
}

// PathsToRoot returns the different possible paths to the root SNOMED-CT concept from this one.
func (ds *Snomed) PathsToRoot(concept *rf2.Concept) ([][]*rf2.Concept, error) {
	parents, err := ds.GetParents(concept)
	if err != nil {
		return nil, err
	}
	results := make([][]*rf2.Concept, 0, len(parents))
	if len(parents) == 0 {
		results = append(results, []*rf2.Concept{concept})
	}
	for _, parent := range parents {
		parentResults, err := ds.PathsToRoot(parent)
		if err != nil {
			return nil, err
		}
		for _, parentResult := range parentResults {
			r := append([]*rf2.Concept{concept}, parentResult...) // prepend current concept
			results = append(results, r)
		}
	}
	return results, nil
}

func debugPaths(paths [][]*rf2.Concept) {
	for i, path := range paths {
		fmt.Printf("Path %d: ", i)
		debugPath(path)
	}
}

func debugPath(path []*rf2.Concept) {
	for _, concept := range path {
		fmt.Printf("%d-", concept.ID)
	}
	fmt.Print("\n")
}

// Genericise finds the best generic match for the given concept
// The "best" is chosen as the closest match to the specified concept and so
// if there are generic concepts which relate to one another, it will be the
// most specific (closest) match to the concept.
func (ds *Snomed) Genericise(concept *rf2.Concept, generics map[snomed.Identifier]*rf2.Concept) (*rf2.Concept, bool) {
	paths, err := ds.PathsToRoot(concept)
	if err != nil {
		return nil, false
	}
	var bestPath []*rf2.Concept
	bestPos := -1
	for _, path := range paths {
		for i, concept := range path {
			if generics[concept.ID] != nil {
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
func (ds *Snomed) GenericiseToRoot(concept *rf2.Concept, root snomed.Identifier) (*rf2.Concept, error) {
	paths, err := ds.PathsToRoot(concept)
	if err != nil {
		return nil, err
	}
	var bestPath []*rf2.Concept
	bestPos := -1
	for _, path := range paths {
		for i, concept := range path {
			if concept.ID == root {
				if i > 0 && (bestPos == -1 || bestPos > i) {
					bestPos = i
					bestPath = path
				}
			}
		}
	}
	if bestPos == -1 {
		return nil, fmt.Errorf("Root concept of %d not found for concept %d", root, concept.ID)
	}
	return bestPath[bestPos-1], nil
}
