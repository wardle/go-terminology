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
	"sort"
	"strings"
	"time"

	"github.com/wardle/go-terminology/snomed"
	"golang.org/x/text/language"
)

const (
	descriptorName = "sctdb.json"
	currentVersion = 0.8
	storeKind      = "level"
	searchKind     = "bleve"
)

// Svc encapsulates concrete persistent and search services and extends it by providing
// semantic inference and a useful, practical SNOMED-CT API.
//
// The current priority of development is correct behaviour rather than optimisation,
// although most operations are extremely fast already.
//
// TODO(mw): once 'correct', profile and optimise the slowest paths, likely
// adding multiple caches for the most commonly derived data structures and
// putting more functionality within the backend transaction, when appropriate.
// It won't be until API complete that we'll understand the pinch points.
//
// It is likely that the transitive closure lists will need more caching, but it is
// unclear whether that is a simple flat list or, more likely now with more complex logic
// for expressions, the individual paths to root.
//
type Svc struct {
	path string
	store
	search
	Descriptor
	language.Matcher
}

// Descriptor provides a simple structure for file-backed database versioning
// and configuration.
type Descriptor struct {
	Version    float32
	StoreKind  string
	SearchKind string
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

	// Concept returns the specified concept
	Concept(conceptID int64) (*snomed.Concept, error)

	// Concepts returns the specified concepts
	Concepts(conceptIDs ...int64) ([]*snomed.Concept, error)

	// Description returns a specified description
	Description(descriptionID int64) (*snomed.Description, error)

	// Descriptions returns the descriptions for a given concept
	Descriptions(conceptID int64) ([]*snomed.Description, error)

	// ParentRelationships returns the parent relationships for a given concept
	ParentRelationships(conceptID int64) ([]*snomed.Relationship, error)

	// ChildRelationships returns the child relationships for a given concept
	ChildRelationships(conceptID int64) ([]*snomed.Relationship, error)

	// AllChildrenIDs returns all children for the specified concept
	AllChildrenIDs(conceptID int64, maximum int) ([]int64, error)

	// ComponentReferenceSets returns the reference set membership for a given SNOMED component
	ComponentReferenceSets(componentID int64) ([]int64, error)

	// GetReferenceSetItems returns all items in a specified reference set
	ReferenceSetComponents(refset int64) (map[int64]struct{}, error)

	// ComponentFromReferenceSet returns the specified components member(s) in the given reference set, if a member
	ComponentFromReferenceSet(refset int64, component int64) ([]*snomed.ReferenceSetItem, error)

	// MapTarget returns the items from the specified reference set with the given target
	MapTarget(refset int64, target string) ([]*snomed.ReferenceSetItem, error)

	// InstalledReferenceSets returns all installed reference sets
	InstalledReferenceSets() (map[int64]struct{}, error) // list of installed reference sets

	// Put is the standard polymorphic way of storing a component within the backing store
	Put(components interface{}) error

	// Iterate permits iteration across all concepts
	Iterate(fn func(*snomed.Concept) error) error

	// Statistics returns overall statistics for the backing store
	Statistics() (Statistics, error)

	// ClearPrecomputations clear precomputed indices
	ClearPrecomputations() error

	// PerformPrecomputations builds indices that can only be performed after a complete import.
	PerformPrecomputations() error

	// Close closes the opened backend store
	Close() error
}

// search represents the backend opaque abstract SNOMED-CT search service.
type search interface {
	Index(eds []*snomed.ExtendedDescription) error
	Search(sr *snomed.SearchRequest) ([]int64, error) //TODO: rename autocomplete
	Close() error
}

// NewService opens or creates a service at the specified location.
func NewService(path string, readOnly bool) (*Svc, error) {
	err := os.MkdirAll(path, 0771)
	if err != nil {
		return nil, err
	}
	descriptor, err := createOrOpenDescriptor(path, storeKind, searchKind)
	if err != nil {
		return nil, err
	}
	if descriptor.Version != currentVersion {
		return nil, fmt.Errorf("Incompatible database format v%1f, needed %1f", descriptor.Version, currentVersion)
	}
	if descriptor.StoreKind != storeKind {
		return nil, fmt.Errorf("Incompatible database format '%s', needed %s", descriptor.StoreKind, storeKind)
	}
	if descriptor.SearchKind != searchKind {
		return nil, fmt.Errorf("Incompatible database format '%s', needed %s", descriptor.SearchKind, searchKind)
	}
	bolt, err := newLevelService(filepath.Join(path, "level.db"), readOnly)
	if err != nil {
		return nil, err
	}
	bleve, err := newBleveIndex(filepath.Join(path, "bleve.db"), readOnly)
	if err != nil {
		return nil, err
	}
	svc := &Svc{path: path, store: bolt, search: bleve, Descriptor: *descriptor, Matcher: newMatcher(bolt)}
	return svc, nil
}

// Close closes any open resources in the backend implementations
func (svc *Svc) Close() error {
	if svc.search != nil {
		if err := svc.search.Close(); err != nil {
			return err
		}
	}
	return svc.store.Close()
}

func createOrOpenDescriptor(path string, storeKind string, searchKind string) (*Descriptor, error) {
	descriptorFilename := filepath.Join(path, descriptorName)
	if _, err := os.Stat(descriptorFilename); os.IsNotExist(err) {
		desc := &Descriptor{
			Version:    currentVersion,
			StoreKind:  storeKind,
			SearchKind: searchKind,
		}
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

// Search searches the SNOMED CT hierarchy
func (svc *Svc) Search(req *snomed.SearchRequest, tags []language.Tag) (*snomed.SearchResponse, error) {
	descriptionIDs, err := svc.search.Search(req)
	if err != nil {
		return nil, err
	}
	items := make([]snomed.SearchResponse_Item, len(descriptionIDs))

	for i, dID := range descriptionIDs {
		if dID == 0 {
			continue
		}
		d, err := svc.Description(dID)
		if err != nil {
			return nil, err
		}
		items[i].Term = d.Term
		items[i].ConceptId = d.ConceptId
		pd, ok, err := svc.PreferredSynonym(d.ConceptId, tags)
		if err != nil {
			return nil, err
		}
		if ok {
			items[i].PreferredTerm = pd.Term
		} else {
			items[i].PreferredTerm = d.Term // fallback to using term instead of preferred term
		}
	}
	result := make([]*snomed.SearchResponse_Item, len(descriptionIDs))
	for i := range items {
		result[i] = &items[i]
	}
	response := new(snomed.SearchResponse)
	response.Items = result
	return response, nil
}

// IsA tests whether the given concept is a type of the specified
// This is a crude implementation which, probably, should be optimised or cached
// much like the old t_cached_parent_concepts table in the SQL version
func (svc *Svc) IsA(concept *snomed.Concept, parent int64) bool {
	if concept.Id == parent {
		return true
	}
	parents, err := svc.AllParents(concept)
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

// FullySpecifiedName returns the FSN (fully specified name) for the given concept, from the
// language reference sets specified, in order of preference
func (svc *Svc) FullySpecifiedName(concept *snomed.Concept, tags []language.Tag) (*snomed.Description, bool, error) {
	descs, err := svc.Descriptions(concept.Id)
	if err != nil {
		return nil, false, err
	}
	return svc.languageMatch(descs, snomed.FullySpecifiedName, tags)
}

// MustGetFullySpecifiedName returns the FSN for the given concept, or panics if there is an error or it is missing
// from the language reference sets specified, in order of preference
func (svc *Svc) MustGetFullySpecifiedName(concept *snomed.Concept, tags []language.Tag) *snomed.Description {
	fsn, found, err := svc.FullySpecifiedName(concept, tags)
	if !found || err != nil {
		panic(fmt.Errorf("Could not determine FSN for concept %d : %s", concept.Id, err))
	}
	return fsn
}

// PreferredSynonym returns the preferred synonym the specified concept based
// on the language preferences specified, in order of preference
func (svc *Svc) PreferredSynonym(conceptID int64, tags []language.Tag) (*snomed.Description, bool, error) {
	descs, err := svc.Descriptions(conceptID)
	if err != nil {
		return nil, false, err
	}
	return svc.languageMatch(descs, snomed.Synonym, tags)
}

// MustGetPreferredSynonym returns the preferred synonym for the specified concept, using the
// language preferences specified, in order of preference
func (svc *Svc) MustGetPreferredSynonym(conceptID int64, tags []language.Tag) *snomed.Description {
	d, found, err := svc.PreferredSynonym(conceptID, tags)
	if err != nil || !found {
		panic(fmt.Errorf("could not determine preferred synonym for concept %d : %s", conceptID, err))
	}
	return d
}

// languageMatch finds the best match for the type of description using the language preferences supplied.
func (svc *Svc) languageMatch(descs []*snomed.Description, typeID snomed.DescriptionTypeID, tags []language.Tag) (*snomed.Description, bool, error) {
	d, found, err := svc.refsetLanguageMatch(descs, typeID, tags)
	if !found && err == nil {
		return svc.simpleLanguageMatch(descs, typeID, tags)
	}
	return d, found, err
}

// simpleLanguageMatch attempts to match a requested language using only the
// language codes in each of the descriptions, without recourse to a language refset.
// this is useful as a fallback in case a concept isn't included in the known language refset
// (e.g. the UK DM+D) or if a specific language reference set isn't installed.
func (svc *Svc) simpleLanguageMatch(descs []*snomed.Description, typeID snomed.DescriptionTypeID, tags []language.Tag) (*snomed.Description, bool, error) {
	dTags := make([]language.Tag, 0)
	ds := make([]*snomed.Description, 0)
	for _, desc := range descs {
		if desc.TypeId == int64(typeID) {
			dTags = append(dTags, desc.LanguageTag())
			ds = append(ds, desc)
		}
	}
	if len(ds) == 0 { // we matched no description
		return nil, false, fmt.Errorf("No descriptions matched type %d in list %v", typeID, descs)
	}
	matcher := language.NewMatcher(dTags)
	_, i, _ := matcher.Match(tags...)
	return ds[i], true, nil
}

// refsetLanguageMatch attempts to match the required language by using known language reference sets
func (svc *Svc) refsetLanguageMatch(descs []*snomed.Description, typeID snomed.DescriptionTypeID, tags []language.Tag) (*snomed.Description, bool, error) {
	preferred := svc.Match(tags)
	for _, desc := range descs {
		if desc.TypeId == int64(typeID) {
			refsetItems, err := svc.ComponentFromReferenceSet(preferred.LanguageReferenceSetIdentifier(), desc.Id)
			if err != nil {
				return nil, false, err
			}
			for _, refset := range refsetItems {
				if refset.GetLanguage().IsPreferred() {
					return desc, true, nil
				}
			}
		}
	}
	return nil, false, nil
}

// Siblings returns the siblings of this concept, ie: those who share the same parents
func (svc *Svc) Siblings(concept *snomed.Concept) ([]*snomed.Concept, error) {
	parents, err := svc.Parents(concept)
	if err != nil {
		return nil, err
	}
	siblings := make([]*snomed.Concept, 0, 10)
	for _, parent := range parents {
		children, err := svc.Children(parent)
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

// AllParents returns all of the parents (recursively) for a given concept
func (svc *Svc) AllParents(concept *snomed.Concept) ([]*snomed.Concept, error) {
	parents, err := svc.AllParentIDs(concept)
	if err != nil {
		return nil, err
	}
	return svc.Concepts(parents...)
}

// AllParentIDs returns a list of the identifiers for all parents
// TODO(mw): switch to using transitive closure
func (svc *Svc) AllParentIDs(concept *snomed.Concept) ([]int64, error) {
	parents := make(map[int64]bool)
	err := svc.allParents(concept, parents)
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

func (svc *Svc) allParents(concept *snomed.Concept, parents map[int64]bool) error {
	ps, err := svc.Parents(concept)
	if err != nil {
		return err
	}
	for _, p := range ps {
		parents[p.Id] = true
		svc.allParents(p, parents)
	}
	return nil
}

// Parents returns the direct IS-A relations of the specified concept.
func (svc *Svc) Parents(concept *snomed.Concept) ([]*snomed.Concept, error) {
	return svc.ParentsOfKind(concept, snomed.IsA)
}

// ParentsOfKind returns the active relations of the specified kinds (types) for the specified concept
func (svc *Svc) ParentsOfKind(concept *snomed.Concept, kinds ...int64) ([]*snomed.Concept, error) {
	result, err := svc.ParentIDsOfKind(concept, kinds...)
	if err != nil {
		return nil, err
	}
	return svc.Concepts(result...)
}

// ParentIDsOfKind returns the active relations of the specified kinds (types) for the specified concept
// Unfortunately, SNOMED-CT isn't perfect and there are some duplicate relationships so
// we filter these and return only unique results
func (svc *Svc) ParentIDsOfKind(concept *snomed.Concept, kinds ...int64) ([]int64, error) {
	relations, err := svc.ParentRelationships(concept.Id)
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

// Children returns the direct IS-A relations of the specified concept.
func (svc *Svc) Children(concept *snomed.Concept) ([]*snomed.Concept, error) {
	return svc.ChildrenOfKind(concept, snomed.IsA)
}

// ChildrenOfKind returns the relations of the specified kind (type) of the specified concept.
func (svc *Svc) ChildrenOfKind(concept *snomed.Concept, kind int64) ([]*snomed.Concept, error) {
	relations, err := svc.ChildRelationships(concept.Id)
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
	return svc.Concepts(result...)
}

// AllChildren fetches all children of the given concept recursively.
// Use with caution with concepts at high levels of the hierarchy.
func (svc *Svc) AllChildren(concept *snomed.Concept, maximum int) ([]*snomed.Concept, error) {
	children, err := svc.AllChildrenIDs(concept.Id, maximum)
	if err != nil {
		return nil, err
	}
	return svc.Concepts(children...)
}

// ConceptsForRelationship returns the concepts represented within a relationship
func (svc *Svc) ConceptsForRelationship(rel *snomed.Relationship) (source *snomed.Concept, kind *snomed.Concept, target *snomed.Concept, err error) {
	concepts, err := svc.Concepts(rel.SourceId, rel.TypeId, rel.DestinationId)
	if err != nil {
		return nil, nil, nil, err
	}
	return concepts[0], concepts[1], concepts[2], nil
}

// PathsToRoot returns the different possible paths to the root SNOMED-CT concept from this one.
// The passed in concept will be the first entry of each path, the SNOMED root will be the last.
func (svc *Svc) PathsToRoot(concept *snomed.Concept) ([][]*snomed.Concept, error) {
	parents, err := svc.Parents(concept)
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
func (svc *Svc) GenericiseTo(concept *snomed.Concept, generics map[int64]struct{}) (*snomed.Concept, bool) {
	if _, ok := generics[concept.Id]; ok {
		return concept, true
	}
	paths, err := svc.PathsToRoot(concept)
	if err != nil {
		return nil, false
	}
	sort.Slice(paths, func(i, j int) bool { // sort our paths in order of length
		return len(paths[i]) < len(paths[j])
	})

	var bestPath []*snomed.Concept
	bestPos, bestLength := -1, 0
	for _, path := range paths {
		for i, concept := range path {
			if _, ok := generics[concept.Id]; ok {
				if bestPos == -1 || bestPos > i || (bestPos == i && len(path) > bestLength) {
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

// ShortestPathToRoot returns the shortest path to the root concept from the specified concept
func (svc *Svc) ShortestPathToRoot(concept *snomed.Concept) (shortest []*snomed.Concept, err error) {
	paths, err := svc.PathsToRoot(concept)
	if err != nil {
		return nil, err
	}
	shortestLength := -1
	for _, path := range paths {
		length := len(path)
		if shortestLength == -1 || shortestLength > length {
			shortest = path
			shortestLength = length
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

// Primitive finds the closest primitive for the specified concept in the hierarchy
func (svc *Svc) Primitive(concept *snomed.Concept) (*snomed.Concept, error) {
	if concept.IsPrimitive() {
		return concept, nil
	}
	paths, err := svc.PathsToRoot(concept)
	if err != nil {
		return nil, err
	}
	bestLength := -1
	var best *snomed.Concept
	for _, path := range paths {
		for i, c := range path {
			if c.IsPrimitive() && (bestLength == -1 || bestLength > i) {
				bestLength = i
				best = c
			}
		}
	}
	return best, nil
}

// ExtendedConcept returns a denormalised representation of a SNOMED CT concept
func (svc *Svc) ExtendedConcept(conceptID int64, tags []language.Tag) (*snomed.ExtendedConcept, error) {
	c, err := svc.Concept(conceptID)
	if err != nil {
		return nil, err
	}
	result := snomed.ExtendedConcept{}
	result.Concept = c
	refsets, err := svc.ComponentReferenceSets(c.Id)
	if err != nil {
		return nil, err
	}
	result.ConceptRefsets = refsets
	relationships, err := svc.ParentRelationships(c.Id)
	if err != nil {
		return nil, err
	}
	result.Relationships = relationships
	recursiveParentIDs, err := svc.AllParentIDs(c)
	if err != nil {
		return nil, err
	}
	result.RecursiveParentIds = recursiveParentIDs
	directParents, err := svc.ParentIDsOfKind(c, snomed.IsA)
	if err != nil {
		return nil, err
	}
	result.DirectParentIds = directParents
	result.PreferredDescription = svc.MustGetPreferredSynonym(conceptID, tags)
	return &result, nil
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

// ClearPrecomputations clears all pre-computations and indices
func (svc *Svc) ClearPrecomputations() error {
	if err := svc.store.ClearPrecomputations(); err != nil {
		return err
	}
	svc.search.Close()
	path := filepath.Join(svc.path, "bleve.db")
	os.RemoveAll(path)
	search, err := newBleveIndex(path, false)
	svc.search = search
	return err
}

// PerformPrecomputations runs all pre-computations and generation of indices
func (svc *Svc) PerformPrecomputations(verbose bool) error {
	if err := svc.store.PerformPrecomputations(); err != nil {
		return err
	}
	tags, _, err := language.ParseAcceptLanguage("en-GB") // TODO: better language handling for search index
	if err != nil {
		return err
	}
	batchSize := 50000
	batch := make([]*snomed.ExtendedDescription, 0)
	total := 0
	start := time.Now()
	svc.iterateExtendedDescriptions(tags, func(ed *snomed.ExtendedDescription) error {
		batch = append(batch, ed)
		if len(batch) == batchSize {
			total += len(batch)
			if verbose {
				elapsed := time.Since(start)
				fmt.Fprintf(os.Stderr, "\rProcessed %d descriptions in %s. Mean time per description: %s...", total, elapsed, elapsed/time.Duration(total))
			}
			if err := svc.search.Index(batch); err != nil {
				return nil
			}
			batch = make([]*snomed.ExtendedDescription, 0)
		}
		return nil
	})
	if err = svc.search.Index(batch); err != nil {
		return err
	}
	total += len(batch)
	fmt.Fprintf(os.Stderr, "\nProcessed total: %d descriptions in %s.\n", total, time.Since(start))
	return nil
}
