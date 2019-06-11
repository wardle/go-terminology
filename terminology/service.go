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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	path   string
	store  Store
	search Search
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
	store, err := newLevelService(filepath.Join(path, "level.db"), readOnly)
	if err != nil {
		return nil, err
	}
	bleve, err := newBleveIndex(filepath.Join(path, "bleve.db"), readOnly)
	if err != nil {
		return nil, err
	}
	svc := &Svc{path: path, store: store, search: bleve, Descriptor: *descriptor}
	svc.Matcher = newMatcher(svc) // TODO: this is weird
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

// Put a slice of SNOMED-CT components into persistent storage.
// This is polymorphic but expects a slice of SNOMED CT components
func (svc *Svc) Put(components interface{}) error {
	var err error
	switch components.(type) {
	case []*snomed.Concept:
		err = svc.putConcepts(components.([]*snomed.Concept))
	case []*snomed.Description:
		err = svc.putDescriptions(components.([]*snomed.Description))
	case []*snomed.Relationship:
		err = svc.putRelationships(components.([]*snomed.Relationship))
	case []*snomed.ReferenceSetItem:
		err = svc.putReferenceSets(components.([]*snomed.ReferenceSetItem))
	default:
		err = fmt.Errorf("unknown component type: %T", components)
	}
	return err
}

// Concept returns the concept with the given identifier
func (svc *Svc) Concept(conceptID int64) (*snomed.Concept, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(conceptID))
	var c snomed.Concept
	return &c, svc.store.View(func(batch Batch) error {
		return batch.Get(bkConcepts, key, &c)
	})
}

// Concepts returns a list of concepts with the given identifiers
func (svc *Svc) Concepts(conceptIDs ...int64) ([]*snomed.Concept, error) {
	key := make([]byte, 8)
	l := len(conceptIDs)
	r1 := make([]snomed.Concept, l)
	r2 := make([]*snomed.Concept, l)
	err := svc.store.View(func(batch Batch) error {
		for i, id := range conceptIDs {
			binary.BigEndian.PutUint64(key, uint64(id))
			if err := batch.Get(bkConcepts, key, &r1[i]); err != nil {
				return err
			}
			r2[i] = &r1[i]
		}
		return nil
	})
	return r2, err
}

// putConcepts persists the specified concepts
func (svc *Svc) putConcepts(concepts []*snomed.Concept) error {
	key := make([]byte, 8)
	return svc.store.Update(func(batch Batch) error {
		for _, c := range concepts {
			binary.BigEndian.PutUint64(key, uint64(c.Id))
			batch.Put(bkConcepts, key, c)
		}
		return nil
	})
}

// PutDescriptions persists the specified descriptions
func (svc *Svc) putDescriptions(descriptions []*snomed.Description) error {
	dID := make([]byte, 8)
	cID := make([]byte, 8)
	return svc.store.Update(func(batch Batch) error {
		for _, d := range descriptions {
			binary.BigEndian.PutUint64(dID, uint64(d.Id))
			batch.Put(bkDescriptions, dID, d)
			binary.BigEndian.PutUint64(cID, uint64(d.ConceptId))
			batch.AddIndexEntry(ixConceptDescriptions, cID, dID)
		}
		return nil
	})
}

// Description returns the description with the given identifier
func (svc *Svc) Description(descriptionID int64) (*snomed.Description, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(descriptionID))
	var c snomed.Description
	return &c, svc.store.View(func(batch Batch) error {
		return batch.Get(bkDescriptions, key, &c)
	})
}

// Descriptions returns the descriptions for a concept
func (svc *Svc) Descriptions(conceptID int64) (result []*snomed.Description, err error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(conceptID))
	svc.store.View(func(batch Batch) error {
		values, err := batch.GetIndexEntries(ixConceptDescriptions, key)
		if err != nil {
			return err
		}
		l := len(values)
		descs := make([]snomed.Description, l)
		result = make([]*snomed.Description, l)
		for i, v := range values {
			if err := batch.Get(bkDescriptions, v, &descs[i]); err != nil {
				return err
			}
			result[i] = &descs[i]
		}
		return nil
	})
	return
}

// PutRelationship persists the specified relationship
func (svc *Svc) putRelationships(relationships []*snomed.Relationship) error {
	rID := make([]byte, 8)
	sourceID := make([]byte, 8)
	destinationID := make([]byte, 8)
	return svc.store.Update(func(batch Batch) error {
		for _, r := range relationships {
			binary.BigEndian.PutUint64(rID, uint64(r.Id))
			binary.BigEndian.PutUint64(sourceID, uint64(r.SourceId))
			binary.BigEndian.PutUint64(destinationID, uint64(r.DestinationId))
			batch.Put(bkRelationships, rID, r)
			batch.AddIndexEntry(ixConceptParentRelationships, sourceID, rID)
			batch.AddIndexEntry(ixConceptChildRelationships, destinationID, rID)
		}
		return nil
	})
}

// ChildRelationships returns the child relationships for this concept.
// Child relationships are relationships in which this concept is the destination.
func (svc *Svc) ChildRelationships(conceptID int64) ([]*snomed.Relationship, error) {
	return svc.getRelationships(conceptID, ixConceptChildRelationships)
}

// ParentRelationships returns the parent relationships for this concept.
// Parent relationships are relationships in which this concept is the source.
func (svc *Svc) ParentRelationships(conceptID int64) ([]*snomed.Relationship, error) {
	return svc.getRelationships(conceptID, ixConceptParentRelationships)
}

// getRelationships returns relationships using the specified bucket
func (svc *Svc) getRelationships(conceptID int64, idx bucket) ([]*snomed.Relationship, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(conceptID))
	var result []*snomed.Relationship
	return result, svc.store.View(func(batch Batch) error {
		entries, err := batch.GetIndexEntries(idx, key)
		if err != nil {
			return err
		}
		l := len(entries)
		relationships := make([]snomed.Relationship, l)
		result = make([]*snomed.Relationship, l)
		for i, id := range entries {
			if err := batch.Get(bkRelationships, id, &relationships[i]); err != nil {
				return err
			}
			result[i] = &relationships[i]
		}
		return nil
	})
}

func (svc *Svc) putReferenceSets(refset []*snomed.ReferenceSetItem) error {
	referencedComponentID := make([]byte, 8)
	refsetID := make([]byte, 8)
	return svc.store.Update(func(batch Batch) error {
		for _, item := range refset {
			itemID := []byte(item.Id)
			batch.Put(bkRefsetItems, itemID, item)
			binary.BigEndian.PutUint64(referencedComponentID, uint64(item.ReferencedComponentId))
			binary.BigEndian.PutUint64(refsetID, uint64(item.RefsetId))
			batch.AddIndexEntry(ixComponentReferenceSets, referencedComponentID, refsetID)
			batch.AddIndexEntry(ixReferenceSetComponentItems, compoundKey(refsetID, referencedComponentID), itemID)
			// support cross maps such as simple maps and complex maps
			var target string
			if simpleMap := item.GetSimpleMap(); simpleMap != nil {
				target = simpleMap.GetMapTarget()
			} else if complexMap := item.GetComplexMap(); complexMap != nil {
				target = complexMap.GetMapTarget()
			}
			if target != "" {
				batch.AddIndexEntry(ixRefsetTargetItems, compoundKey(refsetID, []byte(target+" ")), itemID)
			}
			// keep track of installed reference sets
			batch.AddIndexEntry(ixReferenceSets, refsetID, nil)
		}
		return nil
	})
}

// ComponentReferenceSets returns the refset identifiers to which this component is a member
func (svc *Svc) ComponentReferenceSets(referencedComponentID int64) ([]int64, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(referencedComponentID))
	var result []int64
	return result, svc.store.View(func(batch Batch) error {
		values, err := batch.GetIndexEntries(ixComponentReferenceSets, key)
		if err != nil {
			return err
		}
		result = make([]int64, len(values))
		for i, v := range values {
			result[i] = int64(binary.BigEndian.Uint64(v))
		}
		return nil
	})
}

// MapTarget returns the simple and complex maps for which the specified target, is the target
func (svc *Svc) MapTarget(refset int64, target string) ([]*snomed.ReferenceSetItem, error) {
	refsetID := make([]byte, 8)
	binary.BigEndian.PutUint64(refsetID, uint64(refset))
	key := compoundKey(refsetID, []byte(target+" ")) // ensure delimiter between refset-target and value
	var result []*snomed.ReferenceSetItem
	err := svc.store.View(func(batch Batch) error {
		values, err := batch.GetIndexEntries(ixRefsetTargetItems, key)
		if err != nil {
			return err
		}
		l := len(values)
		items := make([]snomed.ReferenceSetItem, l)
		result = make([]*snomed.ReferenceSetItem, l)
		for i, v := range values {
			if err := batch.Get(bkRefsetItems, v, &items[i]); err != nil {
				return err
			}
			result[i] = &items[i]
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(result, func(i, j int) bool {
		var gi, gj int64 = -1, -1 // groups, ensure deterministic sort order irrespective of map type
		var pi, pj int64 = -1, -1 // priorities
		if cmi := result[i].GetComplexMap(); cmi != nil {
			gi = cmi.GetMapGroup()
			pi = cmi.GetMapPriority()
		}
		if cmj := result[j].GetComplexMap(); cmj != nil {
			gj = cmj.GetMapGroup()
			pj = cmj.GetMapPriority()
		}
		if gi != gj {
			return gi < gj
		}
		return pi < pj
	})
	return result, nil
}

// ReferenceSetComponents returns the components within a given reference set
func (svc *Svc) ReferenceSetComponents(refset int64) (map[int64]struct{}, error) {
	result := make(map[int64]struct{})
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(refset))
	return result, svc.store.View(func(batch Batch) error {
		values, err := batch.GetIndexEntries(ixReferenceSetComponentItems, key)
		if err != nil {
			return err
		}
		for _, v := range values {
			result[int64(binary.BigEndian.Uint64(v[:8]))] = struct{}{}
		}
		return nil
	})

}

// ComponentFromReferenceSet gets the specified components from the specified refset, or error
func (svc *Svc) ComponentFromReferenceSet(refset int64, component int64) ([]*snomed.ReferenceSetItem, error) {
	refsetID := make([]byte, 8)
	componentID := make([]byte, 8)
	binary.BigEndian.PutUint64(refsetID, uint64(refset))
	binary.BigEndian.PutUint64(componentID, uint64(component))
	key := compoundKey(refsetID, componentID)
	var result []*snomed.ReferenceSetItem
	return result, svc.store.View(func(batch Batch) error {

		values, err := batch.GetIndexEntries(ixReferenceSetComponentItems, key)
		if err != nil {
			return err
		}
		l := len(values)
		items := make([]snomed.ReferenceSetItem, l)
		result = make([]*snomed.ReferenceSetItem, l)
		for i, v := range values {
			if err := batch.Get(bkRefsetItems, v, &items[i]); err != nil {
				return err
			}
			result[i] = &items[i]
		}
		return nil
	})

}

// GetAssociations returns the associations for the specified concept
// e.g. to get the SAME_AS associations for a concept,
// GetAssociations(conceptID, snomed.SameAsReferenceSet)
func (svc *Svc) GetAssociations(conceptID int64, refsetID int64) ([]int64, error) {
	items, err := svc.ComponentFromReferenceSet(refsetID, conceptID)
	if err != nil {
		return nil, err
	}
	result := make([]int64, len(items))
	for i, item := range items {
		result[i] = item.GetAssociation().GetTargetComponentId()
	}
	return result, nil
}

// InstalledReferenceSets returns a list of installed reference sets
func (svc *Svc) InstalledReferenceSets() (map[int64]struct{}, error) {
	result := make(map[int64]struct{})
	return result, svc.store.View(func(batch Batch) error {
		entries, err := batch.GetIndexEntries(ixReferenceSets, nil)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			refsetID := int64(binary.BigEndian.Uint64(entry))
			result[refsetID] = struct{}{}

		}
		return nil
	})
}

// AllChildrenIDs returns the recursive children for this concept.
// This is a potentially large number, depending on where in the hierarchy the concept sits.
// TODO(mw): change to use transitive closure table
func (svc *Svc) AllChildrenIDs(conceptID int64, maximum int) ([]int64, error) {
	allChildren := make(map[int64]struct{})
	err := svc.recursiveChildren(conceptID, allChildren, maximum)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(allChildren))
	for id := range allChildren {
		ids = append(ids, id)
	}
	return ids, nil
}

// this is a brute-force, non-cached temporary version which actually fetches the id
// TODO(mwardle): benchmark and possibly use transitive closure precached table a la java version
func (svc *Svc) recursiveChildren(conceptID int64, allChildren map[int64]struct{}, maximum int) error {
	if len(allChildren) > maximum {
		return fmt.Errorf("Too many children; aborted")
	}
	children, err := svc.getRelationships(conceptID, ixConceptChildRelationships)
	if err != nil {
		return err
	}
	for _, child := range children {
		if child.TypeId == snomed.IsA {
			childID := child.SourceId
			if _, exists := allChildren[childID]; !exists {
				if err := svc.recursiveChildren(childID, allChildren, maximum); err != nil {
					return err
				}
				allChildren[childID] = struct{}{}
			}
		}
	}
	return nil
}

// Iterate is a crude iterator for all concepts, useful for pre-processing and pre-computations
func (svc *Svc) Iterate(fn func(*snomed.Concept) error) error {
	concept := snomed.Concept{}
	return svc.store.View(func(batch Batch) error {
		return batch.Iterate(bkConcepts, nil, func(key, value []byte) error {
			if err := proto.Unmarshal(value, &concept); err != nil {
				return err
			}
			if err := fn(&concept); err != nil {
				return err
			}
			return nil
		})
	})
}

// Statistics returns statistics for the backend store
func (svc *Svc) Statistics(lang string) (Statistics, error) {
	stats := Statistics{}
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		return stats, err
	}
	refsets, err := svc.InstalledReferenceSets()
	if err != nil {
		return stats, err
	}
	stats.refsets = make([]string, 0)
	for refset := range refsets {
		rsd, ok, err := svc.PreferredSynonym(refset, tags)
		if err != nil {
			return stats, err
		}
		if ok {
			stats.refsets = append(stats.refsets, rsd.Term)
		} else {
			stats.refsets = append(stats.refsets, strconv.FormatInt(refset, 10))
		}
	}
	return stats, nil
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
	parents, err := svc.AllParentIDs(concept)
	if err != nil {
		return false
	}
	for _, p := range parents {
		if p == parent {
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
		active := true
		for _, p := range path {
			if p.Active == false {
				active = false
				break
			}
		}
		if active == false {
			continue
		}
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

// ClearPrecomputations clears all pre-computations and indices
func (svc *Svc) ClearPrecomputations() error {
	svc.search.Close()
	path := filepath.Join(svc.path, "bleve.db")
	os.RemoveAll(path)
	search, err := newBleveIndex(path, false)
	svc.search = search
	return err
}

// PerformPrecomputations runs all pre-computations and generation of indices
func (svc *Svc) PerformPrecomputations(verbose bool) error {
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
