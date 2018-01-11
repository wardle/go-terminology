// Copyright 2018 Mark Wardle / Eldrix Ltd
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wardle/go-terminology/snomed"
	"golang.org/x/text/language"
)

const (
	currentVersion = 0.1
)

// Snomed encapsulates concrete persistent and search services and extends it by providing
// semantic inference and a useful, practical SNOMED-CT API.
type Snomed struct {
	Store
	Search
	ServiceDescriptor
	Language language.Tag
}

// ServiceDescriptor provides a simple structure for file-backed database versioning
// and configuration.
type ServiceDescriptor struct {
	Version float32
}

// Store represents the backend opaque abstract SNOMED-CT persistence service.
type Store interface {
	GetConcept(conceptID int) (*snomed.Concept, error)
	GetConcepts(conceptIDs ...int) ([]*snomed.Concept, error)
	GetDescriptions(concept *snomed.Concept) ([]*snomed.Description, error)
	GetParentRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error)
	GetChildRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error)
	GetAllChildrenIDs(concept *snomed.Concept) ([]int, error)
	PutConcepts(concepts []*snomed.Concept) error
	PutDescriptions(descriptions []*snomed.Description) error
	PutRelationships(relationships []*snomed.Relationship) error
	Close() error
}

// Search represents an opaque abstract SNOMED-CT search service.
type Search interface {
	// Search executes a search request and returns description identifiers
	Search(search *SearchRequest) ([]int, error)
	Close() error
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

// NewService opens or creates a service at the specified location.
func NewService(path string, readOnly bool) (*Snomed, error) {
	err := os.MkdirAll(path, 0771)
	if err != nil {
		return nil, err
	}
	descriptor, err := createOrOpenDescriptor(path)
	if err != nil {
		return nil, err
	}
	if descriptor.Version != currentVersion {
		return nil, fmt.Errorf("Incompatible database format v%f, needed %f", descriptor.Version, currentVersion)
	}
	bolt, err := NewBoltService(filepath.Join(path, "bolt.db"), readOnly)
	if err != nil {
		return nil, err
	}
	bleve, err := NewBleveService(filepath.Join(path, "index.bleve"), readOnly)
	if err != nil {
		return nil, err
	}
	return &Snomed{Store: bolt, Search: bleve, ServiceDescriptor: *descriptor, Language: language.BritishEnglish}, nil
}

// Close closes any open resources in the backend implementations
func (ds *Snomed) Close() error {
	if err := ds.Store.Close(); err != nil {
		return err
	}
	return ds.Store.Close()
}

func createOrOpenDescriptor(path string) (*ServiceDescriptor, error) {
	descriptorFilename := filepath.Join(path, "sctdb.json")
	if _, err := os.Stat(descriptorFilename); os.IsNotExist(err) {
		desc := &ServiceDescriptor{Version: currentVersion}
		data, err := json.Marshal(desc)
		if err != nil {
			return nil, err
		}
		ioutil.WriteFile(descriptorFilename, data, 0644)
		return desc, nil
	}
	data, err := ioutil.ReadFile(descriptorFilename)
	if err != nil {
		return nil, err
	}
	var desc ServiceDescriptor
	return &desc, json.Unmarshal(data, &desc)
}

// IsA tests whether the given concept is a type of the specified
// This is a crude implementation which, probably, should be optimised or cached
// much like the old t_cached_parent_concepts table in the SQL version
func (ds *Snomed) IsA(concept *snomed.Concept, parent snomed.Identifier) bool {
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
func (ds *Snomed) GetFullySpecifiedName(concept *snomed.Concept) (*snomed.Description, error) {
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
func (ds *Snomed) MustGetFullySpecifiedName(concept *snomed.Concept) *snomed.Description {
	fsn, err := ds.GetFullySpecifiedName(concept)
	if err != nil {
		panic(fmt.Errorf("Could not determine FSN for concept %d : %s", concept.ID, err))
	}
	return fsn
}

// GetPreferredDescription returns the preferred description for this concept in the default language for this service.
func (ds *Snomed) GetPreferredDescription(concept *snomed.Concept) (*snomed.Description, error) {
	return ds.GetPreferredDescriptionForLanguages(concept, []language.Tag{ds.Language})
}

// GetPreferredDescriptionForLanguages returns the preferred description for this concept in the languages specified
// TODO(mw): this is now wrong as SNOMED-CT RF2 uses subsets to handle language preferences
// TODO(mw): implement new
func (ds *Snomed) GetPreferredDescriptionForLanguages(concept *snomed.Concept, languages []language.Tag) (*snomed.Description, error) {
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
func (ds *Snomed) GetPreferredDescriptions(concept *snomed.Concept) ([]*snomed.Description, error) {
	descriptions, err := ds.GetDescriptions(concept)
	if err != nil {
		return nil, err
	}
	preferred := make([]*snomed.Description, 0, len(descriptions))
	for _, description := range descriptions {
		if description.IsSynonym() {
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
			if child.ID != concept.ID {
				siblings = append(siblings, child)
			}
		}
	}
	return siblings, nil
}

// GetAllParents returns all of the parents (recursively) for a given concept
func (ds *Snomed) GetAllParents(concept *snomed.Concept) ([]*snomed.Concept, error) {
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

func (ds *Snomed) getAllParents(concept *snomed.Concept, parents map[snomed.Identifier]bool) error {
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
func (ds *Snomed) GetParents(concept *snomed.Concept) ([]*snomed.Concept, error) {
	return ds.GetParentsOfKind(concept, snomed.IsAConceptID)
}

// GetParentsOfKind returns the active relations of the specified kinds (types) for the specified concept
func (ds *Snomed) GetParentsOfKind(concept *snomed.Concept, kinds ...snomed.Identifier) ([]*snomed.Concept, error) {
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
func (ds *Snomed) GetChildren(concept *snomed.Concept) ([]*snomed.Concept, error) {
	return ds.GetChildrenOfKind(concept, snomed.IsAConceptID)
}

// GetChildrenOfKind returns the relations of the specified kind (type) of the specified concept.
func (ds *Snomed) GetChildrenOfKind(concept *snomed.Concept, kind snomed.Identifier) ([]*snomed.Concept, error) {
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

// GetAllChildren fetches all children of the given concept recursively.
// Use with caution with concepts at high levels of the hierarchy.
func (ds *Snomed) GetAllChildren(concept *snomed.Concept) ([]*snomed.Concept, error) {
	children, err := ds.GetAllChildrenIDs(concept)
	if err != nil {
		return nil, err
	}
	return ds.GetConcepts(children...)
}

// ConceptsForRelationship returns the concepts represented within a relationship
func (ds *Snomed) ConceptsForRelationship(rel *snomed.Relationship) (source *snomed.Concept, kind *snomed.Concept, target *snomed.Concept, err error) {
	concepts, err := ds.GetConcepts(int(rel.SourceID), int(rel.TypeID), int(rel.DestinationID))
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
		fmt.Printf("%d-", concept.ID)
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
func (ds *Snomed) GenericiseToRoot(concept *snomed.Concept, root snomed.Identifier) (*snomed.Concept, error) {
	paths, err := ds.PathsToRoot(concept)
	if err != nil {
		return nil, err
	}
	var bestPath []*snomed.Concept
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
