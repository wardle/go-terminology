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

package terminology

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/wardle/go-terminology/snomed"
	"golang.org/x/text/language"
)

const (
	descriptorName = "sctdb.json"
	currentVersion = 0.1
)

// Svc encapsulates concrete persistent and search services and extends it by providing
// semantic inference and a useful, practical SNOMED-CT API.
type Svc struct {
	store
	search
	Descriptor
	Language language.Tag
}

// Descriptor provides a simple structure for file-backed database versioning
// and configuration.
type Descriptor struct {
	Version float32
}

// Statistics on the persistence store
type Statistics struct {
	concepts      int
	descriptions  int
	relationships int
	refsetItems   int
	refsets       []string
}

// Store represents the backend opaque abstract SNOMED-CT persistence service.
type store interface {
	GetConcept(conceptID int64) (*snomed.Concept, error)
	GetConcepts(conceptIsvc ...int64) ([]*snomed.Concept, error)
	GetDescriptions(concept *snomed.Concept) ([]*snomed.Description, error)
	GetParentRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error)
	GetChildRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error)
	GetAllChildrenIDs(concept *snomed.Concept) ([]int64, error)
	GetReferenceSets(componentID int64) ([]int64, error)
	GetReferenceSet(refset int64) (map[int64]bool, error)
	GetFromReferenceSet(refset int64, component int64) (*snomed.ReferenceSetItem, error)
	GetAllReferenceSets() ([]int64, error) // list of installed reference sets
	Put(components interface{}) error
	Iterate(fn func(*snomed.Concept) error) error
	GetStatistics() (Statistics, error)
	Close() error
}

// Search represents an opaque abstract SNOMED-CT search service.
type search interface {
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
func NewService(path string, readOnly bool) (*Svc, error) {
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
	bolt, err := newBoltService(filepath.Join(path, "bolt.db"), readOnly)
	if err != nil {
		return nil, err
	}
	bleve, err := newBleveService(filepath.Join(path, "index.bleve"), readOnly)
	if err != nil {
		return nil, err
	}
	return &Svc{store: bolt, search: bleve, Descriptor: *descriptor, Language: language.BritishEnglish}, nil
}

// Close closes any open resources in the backend implementations
func (svc *Svc) Close() error {
	if err := svc.store.Close(); err != nil {
		return err
	}
	return svc.store.Close()
}

func createOrOpenDescriptor(path string) (*Descriptor, error) {
	descriptorFilename := filepath.Join(path, descriptorName)
	if _, err := os.Stat(descriptorFilename); os.IsNotExist(err) {
		desc := &Descriptor{Version: currentVersion}
		return desc, saveDescriptor(path, desc)
	}
	data, err := ioutil.ReadFile(descriptorFilename)
	if err != nil {
		return nil, err
	}
	var desc Descriptor
	return &desc, json.Unmarshal(data, &desc)
}

func saveDescriptor(path string, descriptor *Descriptor) error {
	descriptorFilename := filepath.Join(path, descriptorName)
	data, err := json.Marshal(descriptor)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(descriptorFilename, data, 0644)
}

// IsA tests whether the given concept is a type of the specified
// This is a crude implementation which, probably, should be optimised or cached
// much like the old t_cached_parent_concepts table in the SQL version
func (svc *Svc) IsA(concept *snomed.Concept, parent int64) bool {
	if concept.Id == parent {
		return true
	}
	parents, err := svc.GetAllParents(concept)
	if err != nil {
		return false
	}
	for _, p := range parents {
		if p.Id == parent {
			return true
		}
	}
	return false
}

// GetFullySpecifiedName returns the FSN (fully specified name) for the given concept
func (svc *Svc) GetFullySpecifiedName(concept *snomed.Concept, refsetID int64) (*snomed.Description, error) {
	descs, err := svc.GetDescriptions(concept)
	if err != nil {
		return nil, err
	}
	return svc.getFullySpecifiedName(descs, refsetID)
}

// MustGetFullySpecifiedName returns the FSN for the given concept, or panics if there is an error or it is missing
func (svc *Svc) MustGetFullySpecifiedName(concept *snomed.Concept, refsetID int64) *snomed.Description {
	fsn, err := svc.GetFullySpecifiedName(concept, refsetID)
	if err != nil {
		panic(fmt.Errorf("Could not determine FSN for concept %d : %s", concept.Id, err))
	}
	return fsn
}

// GetPreferredSynonym returns the preferred synonym the specified concept based on the language reference set specified
func (svc *Svc) GetPreferredSynonym(c *snomed.Concept, refsetID int64) (*snomed.Description, error) {
	descs, err := svc.GetDescriptions(c)
	if err != nil {
		return nil, err
	}
	return svc.getPreferredSynonym(descs, refsetID)
}

// MustGetPreferredSynonym returns the preferred synonym for the specified concept
func (svc *Svc) MustGetPreferredSynonym(c *snomed.Concept, refsetID int64) *snomed.Description {
	d, err := svc.GetPreferredSynonym(c, refsetID)
	if err != nil {
		panic(fmt.Errorf("could not determine preferred synonym for concept %d : %s", c.Id, err))
	}
	return d
}

func (svc *Svc) getFullySpecifiedName(descs []*snomed.Description, refsetID int64) (*snomed.Description, error) {
	l := len(descs)
	refsets := make([]*snomed.ReferenceSetItem, l)
	for i, desc := range descs {
		if desc.IsFullySpecifiedName() {
			refset, err := svc.GetFromReferenceSet(refsetID, desc.Id)
			if err != nil {
				return nil, err
			}
			if refset != nil {
				refsets[i] = refset
			}
		}
	}
	for i, refset := range refsets {
		if refset != nil {
			if refset.GetLanguage().IsPreferred() {
				return descs[i], nil
			}
		}
	}
	return nil, fmt.Errorf("No fully specified name found in refset %d", refsetID)
}

func (svc *Svc) getPreferredSynonym(descs []*snomed.Description, refsetID int64) (*snomed.Description, error) {
	l := len(descs)
	refsets := make([]*snomed.ReferenceSetItem, l)
	for i, desc := range descs {
		if desc.IsSynonym() {
			refset, err := svc.GetFromReferenceSet(refsetID, desc.Id)
			if err != nil {
				return nil, err
			}
			if refset != nil {
				refsets[i] = refset
			}
		}
	}
	for i, refset := range refsets {
		if refset != nil {
			if refset.GetLanguage().IsPreferred() {
				return descs[i], nil
			}
		}
	}
	return nil, fmt.Errorf("No preferred description found in refset %d", refsetID)
}

// GetSiblings returns the siblings of this concept, ie: those who share the same parents
func (svc *Svc) GetSiblings(concept *snomed.Concept) ([]*snomed.Concept, error) {
	parents, err := svc.GetParents(concept)
	if err != nil {
		return nil, err
	}
	siblings := make([]*snomed.Concept, 0, 10)
	for _, parent := range parents {
		children, err := svc.GetChildren(parent)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			if child.Id != concept.Id {
				siblings = append(siblings, child)
			}
		}
	}
	return siblings, nil
}

// GetAllParents returns all of the parents (recursively) for a given concept
func (svc *Svc) GetAllParents(concept *snomed.Concept) ([]*snomed.Concept, error) {
	parents, err := svc.GetAllParentIDs(concept)
	if err != nil {
		return nil, err
	}
	return svc.GetConcepts(parents...)
}

// GetAllParentIDs returns a list of the identifiers for all parents
func (svc *Svc) GetAllParentIDs(concept *snomed.Concept) ([]int64, error) {
	parents := make(map[int64]bool)
	err := svc.getAllParents(concept, parents)
	if err != nil {
		return nil, err
	}
	keys := make([]int64, len(parents))
	i := 0
	for k := range parents {
		keys[i] = k
		i++
	}
	return keys, nil
}

func (svc *Svc) getAllParents(concept *snomed.Concept, parents map[int64]bool) error {
	ps, err := svc.GetParents(concept)
	if err != nil {
		return err
	}
	for _, p := range ps {
		parents[p.Id] = true
		svc.getAllParents(p, parents)
	}
	return nil
}

// GetParents returns the direct IS-A relations of the specified concept.
func (svc *Svc) GetParents(concept *snomed.Concept) ([]*snomed.Concept, error) {
	return svc.GetParentsOfKind(concept, snomed.IsAConceptID)
}

// GetParentsOfKind returns the active relations of the specified kinsvc (types) for the specified concept
func (svc *Svc) GetParentsOfKind(concept *snomed.Concept, kinds ...int64) ([]*snomed.Concept, error) {
	result, err := svc.GetParentIDsOfKind(concept, kinds...)
	if err != nil {
		return nil, err
	}
	return svc.GetConcepts(result...)
}

// GetParentIDsOfKind returns the active relations of the specified kinsvc (types) for the specified concept
// Unfortunately, SNOMED-CT isn't perfect and there are some duplicate relationships so
// we filter these and return only unique results
func (svc *Svc) GetParentIDsOfKind(concept *snomed.Concept, kinds ...int64) ([]int64, error) {
	relations, err := svc.GetParentRelationships(concept)
	if err != nil {
		return nil, err
	}
	conceptIDs := make(map[int64]struct{})
	for _, relation := range relations {
		if relation.Active {
			for _, kind := range kinds {
				if relation.TypeId == kind {
					conceptIDs[relation.DestinationId] = struct{}{}
				}
			}
		}
	}
	result := make([]int64, 0, len(conceptIDs))
	for id := range conceptIDs {
		result = append(result, id)
	}
	return result, nil
}

// GetChildren returns the direct IS-A relations of the specified concept.
func (svc *Svc) GetChildren(concept *snomed.Concept) ([]*snomed.Concept, error) {
	return svc.GetChildrenOfKind(concept, snomed.IsAConceptID)
}

// GetChildrenOfKind returns the relations of the specified kind (type) of the specified concept.
func (svc *Svc) GetChildrenOfKind(concept *snomed.Concept, kind int64) ([]*snomed.Concept, error) {
	relations, err := svc.GetChildRelationships(concept)
	if err != nil {
		return nil, err
	}
	conceptIDs := make(map[int64]struct{})
	for _, relation := range relations {
		if relation.Active {
			if relation.TypeId == kind {
				conceptIDs[relation.SourceId] = struct{}{}
			}
		}
	}
	result := make([]int64, 0, len(conceptIDs))
	for id := range conceptIDs {
		result = append(result, id)
	}
	return svc.GetConcepts(result...)
}

// GetAllChildren fetches all children of the given concept recursively.
// Use with caution with concepts at high levels of the hierarchy.
func (svc *Svc) GetAllChildren(concept *snomed.Concept) ([]*snomed.Concept, error) {
	children, err := svc.GetAllChildrenIDs(concept)
	if err != nil {
		return nil, err
	}
	return svc.GetConcepts(children...)
}

// ConceptsForRelationship returns the concepts represented within a relationship
func (svc *Svc) ConceptsForRelationship(rel *snomed.Relationship) (source *snomed.Concept, kind *snomed.Concept, target *snomed.Concept, err error) {
	concepts, err := svc.GetConcepts(rel.SourceId, rel.TypeId, rel.DestinationId)
	if err != nil {
		return nil, nil, nil, err
	}
	return concepts[0], concepts[1], concepts[2], nil
}

// PathsToRoot returns the different possible paths to the root SNOMED-CT concept from this one.
func (svc *Svc) PathsToRoot(concept *snomed.Concept) ([][]*snomed.Concept, error) {
	parents, err := svc.GetParents(concept)
	if err != nil {
		return nil, err
	}
	results := make([][]*snomed.Concept, 0, len(parents))
	if len(parents) == 0 {
		results = append(results, []*snomed.Concept{concept})
	}
	for _, parent := range parents {
		parentResults, err := svc.PathsToRoot(parent)
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
		fmt.Printf("%d-", concept.Id)
	}
	fmt.Print("\n")
}

// GenericiseTo returns the best generic match for the given concept
// The "best" is chosen as the closest match to the specified concept and so
// if there are generic concepts which relate to one another, it will be the
// most specific (closest) match to the concept. To determine this, we use
// the closest match of the longest path.
func (svc *Svc) GenericiseTo(concept *snomed.Concept, generics map[int64]bool) (*snomed.Concept, bool) {
	if generics[concept.Id] {
		return concept, true
	}
	paths, err := svc.PathsToRoot(concept)
	if err != nil {
		return nil, false
	}
	var bestPath []*snomed.Concept
	bestPos, bestLength := -1, 0
	for _, path := range paths {
		for i, concept := range path {
			if generics[concept.Id] {
				if i >= 0 && (bestPos == -1 || bestPos > i || (bestPos == i && len(path) > bestLength)) {
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

// LongestPathToRoot returns the longest path to the root concept from the specified concept
func (svc *Svc) LongestPathToRoot(concept *snomed.Concept) (longest []*snomed.Concept, err error) {
	paths, err := svc.PathsToRoot(concept)
	if err != nil {
		return nil, err
	}
	longestLength := 0
	for _, path := range paths {
		length := len(path)
		if length >= longestLength {
			longest = path
			longestLength = length
		}
	}
	return
}

// GenericiseToRoot walks the SNOMED-CT IS-A hierarchy to find the most general concept
// beneath the specified root.
// This finds the shortest path from the concept to the specified root and then
// returns one concept *down* from that root.
func (svc *Svc) GenericiseToRoot(concept *snomed.Concept, root int64) (*snomed.Concept, error) {
	paths, err := svc.PathsToRoot(concept)
	if err != nil {
		return nil, err
	}
	var bestPath []*snomed.Concept
	bestPos := -1
	for _, path := range paths {
		for i, concept := range path {
			if concept.Id == root {
				if i > 0 && (bestPos == -1 || bestPos > i) {
					bestPos = i
					bestPath = path
				}
			}
		}
	}
	if bestPos == -1 {
		return nil, fmt.Errorf("Root concept of %d not found for concept %d", root, concept.Id)
	}
	return bestPath[bestPos-1], nil
}

func (st Statistics) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Number of concepts: %d\n", st.concepts))
	b.WriteString(fmt.Sprintf("Number of descriptions: %d\n", st.descriptions))
	b.WriteString(fmt.Sprintf("Number of relationships: %d\n", st.relationships))
	b.WriteString(fmt.Sprintf("Number of reference set items: %d\n", st.refsetItems))
	b.WriteString(fmt.Sprintf("Number of installed refsets: %d:\n", len(st.refsets)))
	for _, s := range st.refsets {
		b.WriteString(fmt.Sprintf("  Installed refset: %s\n", s))
	}
	return b.String()
}
