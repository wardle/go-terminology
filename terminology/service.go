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
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/wardle/go-terminology/snomed"
	"golang.org/x/text/language"
)

const (
	descriptorName = "sctdb.json"
	currentVersion = 3
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
	availableLanguages []language.Tag
}

// Descriptor provides a simple structure for file-backed database versioning
// and configuration.
type Descriptor struct {
	Version    int32
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
		return nil, fmt.Errorf("Incompatible database format v%d, needed v%d", descriptor.Version, currentVersion)
	}
	if descriptor.StoreKind != storeKind {
		return nil, fmt.Errorf("Incompatible database format '%s', needed %s", descriptor.StoreKind, storeKind)
	}
	if descriptor.SearchKind != searchKind {
		return nil, fmt.Errorf("Incompatible database format '%s', needed %s", descriptor.SearchKind, searchKind)
	}
	//store, err := newBoltService(filepath.Join(path, "bolt.db"), readOnly)
	store, err := newLevelService(filepath.Join(path, "level.db"), readOnly)
	if err != nil {
		return nil, err
	}
	bleve, err := newBleveIndex(filepath.Join(path, "bleve.db"), readOnly)
	if err != nil {
		return nil, err
	}
	svc := &Svc{path: path, store: store, search: bleve, Descriptor: *descriptor}
	// cache list of available languages from the current distribution
	if svc.availableLanguages, err = svc.AvailableLanguages(); err != nil {
		return nil, err
	}
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
func (svc *Svc) Put(context context.Context, components interface{}) error {
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
	var existing snomed.Concept
	return svc.store.Update(func(batch Batch) error {
		for _, c := range concepts {
			binary.BigEndian.PutUint64(key, uint64(c.Id))
			err := batch.Get(bkConcepts, key, &existing)
			if err == ErrNotFound {
				batch.Put(bkConcepts, key, c)
			} else {
				nt, err := ptypes.Timestamp(c.EffectiveTime)
				if err != nil {
					return err
				}
				ot, err := ptypes.Timestamp(existing.EffectiveTime)
				if err != nil {
					return err
				}
				if nt.After(ot) {
					batch.Put(bkConcepts, key, c)
				}
			}
		}
		return nil
	})
}

// PutDescriptions persists the specified descriptions
func (svc *Svc) putDescriptions(descriptions []*snomed.Description) error {
	dID := make([]byte, 8)
	var existing snomed.Description
	return svc.store.Update(func(batch Batch) error {
		for _, d := range descriptions {
			binary.BigEndian.PutUint64(dID, uint64(d.Id))
			err := batch.Get(bkDescriptions, dID, &existing)
			if err == ErrNotFound {
				batch.Put(bkDescriptions, dID, d)
			} else {
				nt, err := ptypes.Timestamp(d.EffectiveTime)
				if err != nil {
					return err
				}
				ot, err := ptypes.Timestamp(existing.EffectiveTime)
				if err != nil {
					return err
				}
				if nt.After(ot) {
					batch.Put(bkDescriptions, dID, d)
				}
			}
		}
		return nil
	})
}

func (svc *Svc) indexDescriptions(batch Batch, ds []*snomed.Description) {
	cID := make([]byte, 8)
	dID := make([]byte, 8)
	for _, d := range ds {
		binary.BigEndian.PutUint64(cID, uint64(d.ConceptId))
		binary.BigEndian.PutUint64(dID, uint64(d.Id))
		batch.AddIndexEntry(ixConceptDescriptions, cID, dID)
	}
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

// DescriptionIds returns the description identifiers for a concept
func (svc *Svc) descriptionIDs(conceptID int64) (result []int64, err error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(conceptID))
	svc.store.View(func(batch Batch) error {
		values, err := batch.GetIndexEntries(ixConceptDescriptions, key)
		if err != nil {
			return err
		}
		l := len(values)
		result = make([]int64, l)
		for i, v := range values {
			result[i] = int64(binary.BigEndian.Uint64(v))
		}
		return nil
	})
	return
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

// PutRelationship persists the specified relationships
func (svc *Svc) putRelationships(relationships []*snomed.Relationship) error {
	rID := make([]byte, 8)
	var existing snomed.Relationship
	return svc.store.Update(func(batch Batch) error {
		for _, r := range relationships {
			binary.BigEndian.PutUint64(rID, uint64(r.Id))
			err := batch.Get(bkRelationships, rID, &existing)
			if err == ErrNotFound {
				batch.Put(bkRelationships, rID, r)
			} else {
				nt, err := ptypes.Timestamp(r.EffectiveTime)
				if err != nil {
					return err
				}
				ot, err := ptypes.Timestamp(existing.EffectiveTime)
				if err != nil {
					return err
				}
				if nt.After(ot) {
					batch.Put(bkRelationships, rID, r)
				}
			}
		}
		return nil
	})
}

func (svc *Svc) indexRelationships(batch Batch, rs []*snomed.Relationship) {
	rID := make([]byte, 8)
	sourceID := make([]byte, 8)
	destinationID := make([]byte, 8)
	for _, r := range rs {
		binary.BigEndian.PutUint64(rID, uint64(r.Id))
		binary.BigEndian.PutUint64(sourceID, uint64(r.SourceId))
		binary.BigEndian.PutUint64(destinationID, uint64(r.DestinationId))
		batch.AddIndexEntry(ixConceptParentRelationships, sourceID, rID)
		batch.AddIndexEntry(ixConceptChildRelationships, destinationID, rID)
		if r.TypeId == snomed.IsA && r.Active {
			batch.AddIndexEntry(ixConceptParents, sourceID, destinationID)
			batch.AddIndexEntry(ixConceptChildren, destinationID, sourceID)
		}
	}
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

// ReferenceSetItem returns the specified reference set item
func (svc *Svc) ReferenceSetItem(itemID string) (*snomed.ReferenceSetItem, error) {
	var item snomed.ReferenceSetItem
	return &item, svc.store.View(func(batch Batch) error {
		return batch.Get(bkRefsetItems, []byte(itemID), &item)
	})
}

func (svc *Svc) putReferenceSets(refset []*snomed.ReferenceSetItem) error {
	var existing snomed.ReferenceSetItem
	return svc.store.Update(func(batch Batch) error {
		for _, item := range refset {
			itemID := []byte(item.Id)
			err := batch.Get(bkRefsetItems, itemID, &existing)
			if err == ErrNotFound {
				batch.Put(bkRefsetItems, itemID, item)
			} else {
				nt, err := ptypes.Timestamp(item.EffectiveTime)
				if err != nil {
					return err
				}
				ot, err := ptypes.Timestamp(existing.EffectiveTime)
				if err != nil {
					return err
				}
				if nt.After(ot) {
					batch.Put(bkRefsetItems, itemID, item)
				}
			}
		}
		return nil
	})
}

func (svc *Svc) indexRefsetItems(batch Batch, rs []*snomed.ReferenceSetItem) {
	refsetID := make([]byte, 8)
	referencedComponentID := make([]byte, 8)
	for _, r := range rs {
		itemID := []byte(r.Id)
		binary.BigEndian.PutUint64(referencedComponentID, uint64(r.ReferencedComponentId))
		binary.BigEndian.PutUint64(refsetID, uint64(r.RefsetId))
		batch.AddIndexEntry(ixComponentReferenceSets, referencedComponentID, refsetID)
		batch.AddIndexEntry(ixReferenceSetComponentItems, compoundKey(refsetID, referencedComponentID), itemID)
		// support cross maps such as simple maps and complex maps
		var target string
		if simpleMap := r.GetSimpleMap(); simpleMap != nil {
			target = simpleMap.GetMapTarget()
		} else if complexMap := r.GetComplexMap(); complexMap != nil {
			target = complexMap.GetMapTarget()
		}
		if target != "" {
			batch.AddIndexEntry(ixRefsetTargetItems, compoundKey(refsetID, []byte(target+" ")), itemID)
		}
		// keep track of installed reference sets
		batch.AddIndexEntry(ixReferenceSets, refsetID, nil)
	}
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

// IsInReferenceSet returns whether the specified component is in the specified reference set
func (svc *Svc) IsInReferenceSet(referencedComponentID int64, refsetID int64) (result bool, err error) {
	k := make([]byte, 8)
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(k, uint64(referencedComponentID))
	binary.BigEndian.PutUint64(v, uint64(refsetID))
	err = svc.store.View(func(batch Batch) error {
		result, err = batch.CheckIndexEntry(ixComponentReferenceSets, k, v)
		return nil
	})
	return
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
				fmt.Fprintf(os.Stderr, fmt.Sprintf("error fetching refset item with identifier %s: %s", v, err))
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

// AllChildrenIDs returns all children for the specified concept.
// As this is potentially a very large number, you must specify an indicative maximum number.
func (svc *Svc) AllChildrenIDs(ctx context.Context, conceptID int64, maximum int) (children []int64, err error) {
	ch := svc.StreamAllChildrenIDs(ctx, conceptID, maximum)
	for c := range ch {
		if c.Err != nil {
			return nil, c.Err
		}
		children = append(children, c.ID)
	}
	return
}

// ConceptIDStream wraps a concept identifier with an error, for use in streaming
type ConceptIDStream struct {
	ID  int64
	Err error
}

// ConceptStream wraps a concept identifier with an error, for use in streaming
type ConceptStream struct {
	*snomed.Concept
	Err error
}

// ConceptReferenceStream wraps a concept reference with an error, helpful when streaming via channels
type ConceptReferenceStream struct {
	*snomed.ConceptReference
	Err error
}

// DescriptionStream is for use in streaming descriptions
type DescriptionStream struct {
	*snomed.Description
	Err error
}

// ExtendedDescriptionStream is for use in streaming extended descriptions
type ExtendedDescriptionStream struct {
	*snomed.ExtendedDescription
	Err error
}

// StreamConceptReferences is a helper function to turn a stream of identifiers into a stream of (more useful)
// ConceptReferences.
func (svc *Svc) StreamConceptReferences(ctx context.Context, concepts <-chan ConceptIDStream, nchannels int, tags []language.Tag) <-chan ConceptReferenceStream {
	out := make(chan ConceptReferenceStream)
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < nchannels; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case c := <-concepts:
					if c.ID == 0 && c.Err == nil {
						return
					}
					if c.Err != nil {
						out <- ConceptReferenceStream{Err: c.Err}
						close(done)
						return
					}
					cr, err := svc.ConceptReference(c.ID, tags)
					if err != nil {
						out <- ConceptReferenceStream{Err: err}
						close(done)
						return
					}
					select {
					case out <- ConceptReferenceStream{ConceptReference: cr}:
					case <-done:
						return
					}
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// StreamAllChildrenIDs streams all of the children of the specified concept.
// The maximum is used as an indicative measure, rather than an absolutely exact target.
func (svc *Svc) StreamAllChildrenIDs(ctx context.Context, conceptID int64, maximum int) <-chan ConceptIDStream {
	var allChildren sync.Map // already processed concepts
	done := make(chan struct{})
	results := make(chan ConceptIDStream)
	work := make(chan int64, maximum) // concepts to be processed
	var wg sync.WaitGroup             // count of work
	wg.Add(1)
	work <- conceptID // send first item of work
	process := func(id int64) error {
		defer wg.Done()
		if _, exists := allChildren.LoadOrStore(id, struct{}{}); exists {
			return nil
		}
		if id != conceptID { // send out a result, only if not original concept
			select {
			case results <- ConceptIDStream{ID: id}:
			case <-ctx.Done():
				return ctx.Err()
			case <-done:
				return nil
			}
		}
		children, err := svc.Children(id) // find more work
		if err != nil {
			return err
		}
		for _, child := range children {
			wg.Add(1)
			select {
			case work <- child:
			case <-ctx.Done():
				return ctx.Err()
			case <-done:
				return nil
			default:
				// we've exhausted our buffer, so give up gracefully
				return fmt.Errorf("too many children")
			}
		}
		return nil
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case c := <-work:
					if c == 0 {
						return
					}
					if err := process(c); err != nil {
						results <- ConceptIDStream{Err: err}
						close(done) // abort
					}
				case <-done:
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(work)
		close(results)
	}()
	return results
}

// IterateConcepts is an iterator for all concepts, useful for pre-processing and pre-computations
func (svc *Svc) IterateConcepts(ctx context.Context) <-chan ConceptStream {
	out := make(chan ConceptStream)
	go func() {
		defer close(out)
		err := svc.store.View(func(batch Batch) error {
			return batch.Iterate(bkConcepts, nil, func(key, value []byte) error {
				concept := new(snomed.Concept)
				if err := proto.Unmarshal(value, concept); err != nil {
					return err
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- ConceptStream{Concept: concept}:
				}
				return nil
			})
		})
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return
			}
			select {
			case <-ctx.Done():
				return
			case out <- ConceptStream{Err: err}:
				return
			}
		}
	}()
	return out
}

func (svc *Svc) iterateDescriptions(ctx context.Context, batchSize int) <-chan []*snomed.Description {
	ch := make(chan []*snomed.Description)
	go func() {
		defer close(ch)
		err := svc.store.View(func(batch Batch) error {
			job := make([]*snomed.Description, 0, batchSize)
			err := batch.Iterate(bkDescriptions, nil, func(key, value []byte) error {
				d := new(snomed.Description)
				if err := proto.Unmarshal(value, d); err != nil {
					return err
				}
				job = append(job, d)
				if len(job) == batchSize {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ch <- job:
					}
					job = make([]*snomed.Description, 0, batchSize)
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
			if len(job) > 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ch <- job:
				}
			}
			return nil
		})
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return
			}
			panic(err)
		}
	}()
	return ch
}
func (svc *Svc) iterateRelationships(ctx context.Context, batchSize int) <-chan []*snomed.Relationship {
	ch := make(chan []*snomed.Relationship)
	go func() {
		defer close(ch)
		err := svc.store.View(func(batch Batch) error {
			job := make([]*snomed.Relationship, 0, batchSize)
			err := batch.Iterate(bkRelationships, nil, func(key, value []byte) error {
				d := new(snomed.Relationship)
				if err := proto.Unmarshal(value, d); err != nil {
					return err
				}
				job = append(job, d)
				if len(job) == batchSize {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ch <- job:
					}
					job = make([]*snomed.Relationship, 0, batchSize)
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
			if len(job) > 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ch <- job:
				}
			}
			return nil
		})
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return
			}
			panic(err)
		}
	}()
	return ch
}
func (svc *Svc) iterateRefsetItems(ctx context.Context, batchSize int) <-chan []*snomed.ReferenceSetItem {
	ch := make(chan []*snomed.ReferenceSetItem)
	go func() {
		defer close(ch)
		err := svc.store.View(func(batch Batch) error {
			job := make([]*snomed.ReferenceSetItem, 0, batchSize)
			err := batch.Iterate(bkRefsetItems, nil, func(key, value []byte) error {
				d := new(snomed.ReferenceSetItem)
				if err := proto.Unmarshal(value, d); err != nil {
					return err
				}
				job = append(job, d)
				if len(job) == batchSize {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ch <- job:
					}
					job = make([]*snomed.ReferenceSetItem, 0, batchSize)
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
			if len(job) > 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ch <- job:
				}
			}
			return nil
		})
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return
			}
			panic(err)
		}
	}()
	return ch
}

func (svc *Svc) countBucket(b bucket, count *uint64) {
	err := svc.store.View(func(batch Batch) error {
		return batch.Iterate(b, nil, func(key, value []byte) error {
			atomic.AddUint64(count, 1)
			return nil
		})
	})
	if err != nil {
		panic(err)
	}
}

// Statistics returns statistics for the backend store
func (svc *Svc) Statistics(lang string, verbose bool) (Statistics, error) {
	stats := Statistics{}
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		return stats, err
	}
	stats.searchIndex, err = svc.search.Statistics()
	if err != nil {
		return stats, err
	}
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	if verbose {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					for _, r := range `-\|/` {
						fmt.Printf("\r%c %4s: %d concepts, %d descriptions, %d relationships and %d refset items...",
							r,
							time.Since(start),
							atomic.LoadUint64(&stats.concepts),
							atomic.LoadUint64(&stats.descriptions),
							atomic.LoadUint64(&stats.relationships),
							atomic.LoadUint64(&stats.refsetItems))
						time.Sleep(50 * time.Millisecond)
					}
				}

			}
		}()
	}
	var cWg, dWg, relWg, riWg, rfWg sync.WaitGroup
	cWg.Add(1)
	go func() {
		defer cWg.Done()
		svc.countBucket(bkConcepts, &stats.concepts)
	}()
	dWg.Add(1)
	go func() {
		defer dWg.Done()
		svc.countBucket(bkDescriptions, &stats.descriptions)
	}()
	relWg.Add(1)
	go func() {
		defer relWg.Done()
		svc.countBucket(bkRelationships, &stats.relationships)
	}()
	riWg.Add(1)
	go func() {
		defer riWg.Done()
		svc.countBucket(bkRefsetItems, &stats.refsetItems)
	}()
	rfWg.Add(1)
	go func() {
		defer rfWg.Done()
		refsets, err := svc.InstalledReferenceSets()
		if err != nil {
			panic(err)
		}
		stats.refsets = make([]string, 0)
		for refset := range refsets {
			rsd, err := svc.PreferredSynonym(refset, tags)
			if err != nil {
				panic(err)
			}
			stats.refsets = append(stats.refsets, rsd.Term)
		}
	}()
	cWg.Wait()
	dWg.Wait()
	relWg.Wait()
	riWg.Wait()
	rfWg.Wait()
	cancel()
	if verbose {
		fmt.Println()
	}
	return stats, nil
}

// Search searches the SNOMED CT hierarchy
func (svc *Svc) Search(req *snomed.SearchRequest, tags []language.Tag) (*snomed.SearchResponse, error) {
	requestedMax := req.MaximumHits
	if req.MaximumHits < 100 { // if request is for less than 100 hits, choose 100 and we throw away others
		req.MaximumHits = 100 // because we need to sort by our own score
	}
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
		pd, err := svc.PreferredSynonym(d.ConceptId, tags)
		if err != nil {
			return nil, err
		}
		items[i].PreferredTerm = pd.Term
	}
	result := make([]*snomed.SearchResponse_Item, len(descriptionIDs))
	for i := range items {
		result[i] = &items[i]
	}
	sort.Slice(result, func(i, j int) bool {
		return len(result[i].Term) < len(result[j].Term)
	})
	if requestedMax > 0 && len(result) > int(requestedMax) {
		result = result[0:requestedMax]
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
	parents, err := svc.AllParentIDs(concept.Id)
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
func (svc *Svc) FullySpecifiedName(concept *snomed.Concept, tags []language.Tag) (*snomed.Description, error) {
	descs, err := svc.Descriptions(concept.Id)
	if err != nil {
		return nil, err
	}
	return svc.languageMatch(descs, snomed.FullySpecifiedName, tags)
}

// MustGetFullySpecifiedName returns the FSN for the given concept, or panics if there is an error or it is missing
// from the language reference sets specified, in order of preference
func (svc *Svc) MustGetFullySpecifiedName(concept *snomed.Concept, tags []language.Tag) *snomed.Description {
	fsn, err := svc.FullySpecifiedName(concept, tags)
	if err != nil {
		panic(fmt.Errorf("Could not determine FSN for concept %d : %s", concept.Id, err))
	}
	return fsn
}

// PreferredSynonymByReferenceSet determines the preferred synonym by virtue of
// member of the description in the specified (language) reference set.
// This is a more appropriate way of determining preferred synonym for concepts within, for example,
// the UK dm+d.
// See https://www.nhsbsa.nhs.uk/sites/default/files/2018-10/doc_SnomedCTUKDrugExtensionModel%20-%20v1.0.pdf
// and see references to the "dm+d realm description refset".
// This falls back to standard language based preferred term.
func (svc *Svc) PreferredSynonymByReferenceSet(conceptID int64, refsetID int64, tags []language.Tag) (*snomed.Description, error) {
	descs, err := svc.descriptionIDs(conceptID)
	if err != nil {
		return nil, err
	}
	for _, dID := range descs {
		refsetItems, err := svc.ComponentFromReferenceSet(refsetID, dID)
		if err != nil {
			return nil, err
		}
		if len(refsetItems) == 0 {
			continue
		}
		d, err := svc.Description(dID)
		if err != nil {
			return nil, err
		}
		if d.IsSynonym() == false {
			continue
		}
		for _, refsetItem := range refsetItems {
			if refsetItem.Active && refsetItem.GetLanguage().IsPreferred() {
				return d, nil
			}
		}
	}
	return svc.PreferredSynonym(conceptID, tags)
}

// PreferredSynonym returns the preferred synonym the specified concept based
// on the language preferences specified, in order of preference
func (svc *Svc) PreferredSynonym(conceptID int64, tags []language.Tag) (*snomed.Description, error) {
	descs, err := svc.Descriptions(conceptID)
	if err != nil {
		return nil, err
	}
	return svc.languageMatch(descs, snomed.Synonym, tags)
}

// MustGetPreferredSynonym returns the preferred synonym for the specified concept, using the
// language preferences specified, in order of preference
func (svc *Svc) MustGetPreferredSynonym(conceptID int64, tags []language.Tag) *snomed.Description {
	d, err := svc.PreferredSynonym(conceptID, tags)
	if err != nil {
		panic(fmt.Errorf("could not determine preferred synonym for concept %d : %s", conceptID, err))
	}
	return d
}

// languageMatch finds the best match for the type of description using the language preferences supplied.
func (svc *Svc) languageMatch(descs []*snomed.Description, typeID snomed.DescriptionTypeID, tags []language.Tag) (*snomed.Description, error) {
	d, found, err := svc.refsetLanguageMatch(descs, typeID, tags)
	if !found && err == nil {
		return svc.simpleLanguageMatch(descs, typeID, tags)
	}
	return d, err
}

// simpleLanguageMatch attempts to match a requested language using only the
// language codes in each of the descriptions, without recourse to a language refset.
// this is useful as a fallback in case a concept isn't included in the known language refset
// (e.g. the UK DM+D) or if a specific language reference set isn't installed.
func (svc *Svc) simpleLanguageMatch(descs []*snomed.Description, typeID snomed.DescriptionTypeID, tags []language.Tag) (*snomed.Description, error) {
	dTags := make([]language.Tag, 0)
	ds := make([]*snomed.Description, 0)
	// make matching deterministic, by ensuring the list of descriptions is ordered in a deterministic way
	sort.Slice(descs, func(i, j int) bool {
		return descs[i].LanguageCode < descs[j].LanguageCode
	})
	for _, desc := range descs {
		if desc.TypeId == int64(typeID) {
			dTags = append(dTags, desc.LanguageTag())
			ds = append(ds, desc)
		}
	}
	if len(ds) == 0 { // we matched no description
		return nil, fmt.Errorf("No descriptions matched type %d in list %v", typeID, descs)
	}
	matcher := language.NewMatcher(dTags)
	_, i, _ := matcher.Match(tags...)
	return ds[i], nil
}

// refsetLanguageMatch attempts to match the required language by using known language reference sets
func (svc *Svc) refsetLanguageMatch(descs []*snomed.Description, typeID snomed.DescriptionTypeID, tags []language.Tag) (*snomed.Description, bool, error) {
	if len(svc.availableLanguages) == 0 {
		return nil, false, nil // apparently no language reference sets installed. give up now
	}
	matcher := language.NewMatcher(svc.availableLanguages)
	_, i, _ := matcher.Match(tags...)
	preferred := LanguageForTag(svc.availableLanguages[i])
	for _, desc := range descs {
		if desc.TypeId == int64(typeID) {
			refsetItems, err := svc.ComponentFromReferenceSet(preferred.LanguageReferenceSetIdentifier(), desc.Id)
			if err != nil {
				return nil, false, err
			}
			for _, refset := range refsetItems {
				if refset.Active && refset.GetLanguage().IsPreferred() {
					return desc, true, nil
				}
			}
		}
	}
	return nil, false, nil
}

// Parents returns the parents of the specified concept
func (svc *Svc) Parents(conceptID int64) ([]int64, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(conceptID))
	var result []int64
	return result, svc.store.View(func(batch Batch) error {
		entries, err := batch.GetIndexEntries(ixConceptParents, key)
		if err != nil {
			return err
		}
		result = make([]int64, len(entries))
		for i, v := range entries {
			result[i] = int64(binary.BigEndian.Uint64(v))
		}
		return nil
	})
}

// Children returns the children of the specified concept
func (svc *Svc) Children(conceptID int64) ([]int64, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(conceptID))
	var result []int64
	return result, svc.store.View(func(batch Batch) error {
		entries, err := batch.GetIndexEntries(ixConceptChildren, key)
		if err != nil {
			return err
		}
		result = make([]int64, len(entries))
		for i, v := range entries {
			result[i] = int64(binary.BigEndian.Uint64(v))
		}
		return nil
	})
}

// Siblings returns the siblings of this concept, ie: those who share the same parents
func (svc *Svc) Siblings(conceptID int64) ([]int64, error) {
	parents, err := svc.Parents(conceptID)
	if err != nil {
		return nil, err
	}
	siblings := make([]int64, 0)
	for _, parent := range parents {
		children, err := svc.Children(parent)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			if child != conceptID {
				siblings = append(siblings, child)
			}
		}
	}
	return siblings, nil
}

// AllParents returns all of the parents (recursively) for a given concept
func (svc *Svc) AllParents(conceptID int64) ([]*snomed.Concept, error) {
	parents, err := svc.AllParentIDs(conceptID)
	if err != nil {
		return nil, err
	}
	return svc.Concepts(parents...)
}

// AllParentIDs returns a list of the identifiers for all parents
// TODO(mw): switch to using transitive closure
func (svc *Svc) AllParentIDs(conceptID int64) ([]int64, error) {
	parents := make(map[int64]struct{})
	err := svc.allParents(conceptID, parents)
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

func (svc *Svc) allParents(conceptID int64, parents map[int64]struct{}) error {
	ps, err := svc.Parents(conceptID)
	if err != nil {
		return err
	}
	for _, p := range ps {
		if _, ok := parents[p]; ok { // have we already processed this?
			continue
		}
		parents[p] = struct{}{}
		svc.allParents(p, parents)
	}
	return nil
}

// ParentIDsOfKind returns the active relations of the specified kinds (types) for the specified concept
// Unfortunately, SNOMED-CT isn't perfect and there are some duplicate relationships so
// we filter these and return only unique results
func (svc *Svc) ParentIDsOfKind(conceptID int64, kinds ...int64) ([]int64, error) {
	relations, err := svc.ParentRelationships(conceptID)
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

// ChildrenOfKind returns the relations of the specified kind (type) of the specified concept.
func (svc *Svc) ChildrenOfKind(conceptID int64, kind int64) ([]int64, error) {
	relations, err := svc.ChildRelationships(conceptID)
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
	return result, nil
}

// AllChildren fetches all children of the given concept recursively.
// Use with caution with concepts at high levels of the hierarchy.
func (svc *Svc) AllChildren(ctx context.Context, concept *snomed.Concept, maximum int) ([]*snomed.Concept, error) {
	children, err := svc.AllChildrenIDs(ctx, concept.Id, maximum)
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
func (svc *Svc) PathsToRoot(conceptID int64) ([][]int64, error) {
	parents, err := svc.Parents(conceptID)
	if err != nil {
		return nil, err
	}
	results := make([][]int64, 0, len(parents))
	if len(parents) == 0 {
		results = append(results, []int64{conceptID})
	}
	for _, parent := range parents {
		parentResults, err := svc.PathsToRoot(parent)
		if err != nil {
			return nil, err
		}
		for _, parentResult := range parentResults {
			r := append([]int64{conceptID}, parentResult...) // prepend current concept
			results = append(results, r)
		}
	}
	return results, nil
}

func debugPaths(paths [][]int64) {
	for i, path := range paths {
		fmt.Printf("Path %d: ", i)
		debugPath(path)
	}
}

func debugPath(path []int64) {
	for _, id := range path {
		fmt.Printf("%d-", id)
	}
	fmt.Print("\n")
}

// GenericiseTo returns the matches for the given concept within the set specified, ordered
// from best to worst. It scores by looking at the paths from concept to root.
func (svc *Svc) GenericiseTo(conceptID int64, includeParents bool, generics map[int64]struct{}) ([]int64, error) {
	allGenerics := make(map[int64]struct{})
	for generic := range generics {
		if _, exists := allGenerics[generic]; !exists {
			allGenerics[generic] = struct{}{}
			if includeParents {
				if err := svc.allParents(generic, allGenerics); err != nil {
					return nil, err
				}
			}
		}
	}
	if _, ok := allGenerics[conceptID]; ok { // if original concept is in the refset, return it
		return []int64{conceptID}, nil
	}
	// find parents of our concept that intersect with the refset+/-refset's parents
	paths, err := svc.PathsToRoot(conceptID)
	if err != nil {
		return nil, err
	}
	results := make(map[int64]float64)
	for _, path := range paths {
		score, concept := scorePath(path, allGenerics)
		if score == 0.0 {
			continue
		}
		if existingScore, ok := results[concept]; ok {
			if score > existingScore {
				results[concept] = score
			}
		} else {
			results[concept] = score
		}
	}
	concepts := make([]int64, 0, len(results))
	for concept := range results {
		concepts = append(concepts, concept)
	}
	sort.Slice(concepts, func(i, j int) bool {
		return results[concepts[i]] > results[concepts[j]]
	})
	return concepts, nil
}

// GenericiseToBest returns the best match for the given concept in the set of concepts specified.
// The "best" is chosen as the closest match to the specified concept and so
// if there are generic concepts which relate to one another, it will be the
// most specific (closest) match to the concept. To determine this, we use
// the closest match of the longest path.
func (svc *Svc) GenericiseToBest(conceptID int64, generics map[int64]struct{}) (int64, error) {
	result, err := svc.GenericiseTo(conceptID, false, generics)
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		result, err = svc.GenericiseTo(conceptID, true, generics)
		if err != nil {
			return 0, err
		}
		if len(result) == 0 {
			return 0, nil
		}
	}
	return result[0], nil
}

// scorePath determines a score for the path, approximating to the most specific concept in the path
// found in the subset (generics).
func scorePath(path []int64, generics map[int64]struct{}) (float64, int64) {
	for i, concept := range path {
		if _, ok := generics[concept]; ok {
			return 1.0 - float64(i)/float64(len(path)), concept
		}
	}
	return 0.0, 0
}

// LongestPathToRoot returns the longest path to the root concept from the specified concept
func (svc *Svc) LongestPathToRoot(conceptID int64) (longest []int64, err error) {
	paths, err := svc.PathsToRoot(conceptID)
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
func (svc *Svc) ShortestPathToRoot(conceptID int64) (shortest []int64, err error) {
	paths, err := svc.PathsToRoot(conceptID)
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
func (svc *Svc) GenericiseToRoot(conceptID int64, root int64) (int64, error) {
	paths, err := svc.PathsToRoot(conceptID)
	if err != nil {
		return 0, err
	}
	var bestPath []int64
	bestPos := -1
	for _, path := range paths {
		for i, concept := range path {
			if concept == root {
				if i > 0 && (bestPos == -1 || bestPos > i) {
					bestPos = i
					bestPath = path
				}
			}
		}
	}
	if bestPos == -1 {
		return 0, fmt.Errorf("Root concept of %d not found for concept %d", root, conceptID)
	}
	return bestPath[bestPos-1], nil
}

// Primitive finds the closest primitive for the specified concept in the hierarchy
func (svc *Svc) Primitive(concept *snomed.Concept) (*snomed.Concept, error) {
	if concept.IsPrimitive() {
		return concept, nil
	}
	paths, err := svc.PathsToRoot(concept.Id)
	if err != nil {
		return nil, err
	}
	bestLength := -1
	var best *snomed.Concept
	for _, path := range paths {
		for i, c := range path {
			concept, err := svc.Concept(c)
			if err != nil {
				return nil, err
			}
			if concept.IsPrimitive() && (bestLength == -1 || bestLength > i) {
				bestLength = i
				best = concept
			}
		}
	}
	return best, nil
}

// ExtendedConcept returns a denormalised representation of a SNOMED CT concept
func (svc *Svc) ExtendedConcept(conceptID int64, tags []language.Tag) (result *snomed.ExtendedConcept, err error) {
	var preferred *snomed.Description
	var relationships []*snomed.Relationship
	var mux sync.Mutex
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		r := snomed.ExtendedConcept{}
		errs := make([]error, 4)
		r.Concept, errs[0] = svc.Concept(conceptID)
		r.ConceptRefsets, errs[1] = svc.ComponentReferenceSets(conceptID)
		r.AllParentIds, errs[2] = svc.AllParentIDs(conceptID)
		r.DirectParentIds, errs[3] = svc.ParentIDsOfKind(conceptID, snomed.IsA)
		mux.Lock()
		defer mux.Unlock()
		for _, e := range errs {
			if e != nil {
				err = e
				return
			}
		}
		result = &r
	}()
	go func() {
		defer wg.Done()
		d, e1 := svc.PreferredSynonym(conceptID, tags)
		mux.Lock()
		defer mux.Unlock()
		preferred = d
		if e1 != nil && err == nil {
			err = e1
		}
	}()
	go func() {
		defer wg.Done()
		rels, e1 := svc.ParentRelationships(conceptID)
		mux.Lock()
		defer mux.Unlock()
		relationships = rels
		if e1 != nil && err == nil {
			err = e1
		}
	}()
	wg.Wait()
	result.PreferredDescription = preferred
	result.Relationships = relationships
	return
}

// ConceptReference creates a reference for the specified concept.
// This is generally more useful than simply getting the Concept itself!
func (svc *Svc) ConceptReference(conceptID int64, tags []language.Tag) (*snomed.ConceptReference, error) {
	var d *snomed.Description
	r := new(snomed.ConceptReference)
	r.ConceptId = conceptID
	d, err := svc.PreferredSynonym(conceptID, tags)
	if err != nil {
		return nil, err
	}
	r.Term = d.Term
	return r, nil
}

// ClearPrecomputations clears all pre-computations and indices
func (svc *Svc) ClearPrecomputations() error {
	// delete all indices
	svc.store.Update(func(b Batch) error {
		var wg sync.WaitGroup
		for idx := ixConceptDescriptions; idx < lastIndex; idx++ {
			wg.Add(1)
			go func(i bucket) {
				defer wg.Done()
				b.ClearIndexEntries(i)
			}(idx)
		}
		wg.Wait()
		return nil
	})
	// close, delete and recreate (empty) search index
	svc.search.Close()
	path := filepath.Join(svc.path, "bleve.db")
	os.RemoveAll(path)
	search, err := newBleveIndex(path, false)
	svc.search = search
	return err
}

// PerformPrecomputations runs all pre-computations and generation of indices
func (svc *Svc) PerformPrecomputations(ctx context.Context, batchSize int, verbose bool) error {
	start := time.Now()
	ncpu := runtime.NumCPU()
	if batchSize == 0 {
		batchSize = 5000
	}
	descriptions := svc.iterateDescriptions(ctx, batchSize)
	relationships := svc.iterateRelationships(ctx, batchSize)
	refsetItems := svc.iterateRefsetItems(ctx, batchSize)
	var dWg, rWg, rsWg sync.WaitGroup
	var nd, nr, nri uint32 // counts of components processed
	done := make(chan bool)
	if verbose {
		fmt.Printf("Indexing....\n")
		go func() {
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				default:
					fmt.Printf("\rIndexing %s: processed %d descriptions, %d relationships and %d reference set items...",
						time.Since(start), atomic.LoadUint32(&nd), atomic.LoadUint32(&nr), atomic.LoadUint32(&nri))
					time.Sleep(500 * time.Microsecond)
				}
			}
		}()
	}
	for i := 0; i < ncpu; i++ {
		dWg.Add(1)
		go func() {
			defer dWg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case ds := <-descriptions:
					if ds == nil {
						return
					}
					err := svc.store.Update(func(batch Batch) error {
						svc.indexDescriptions(batch, ds)
						atomic.AddUint32(&nd, uint32(len(ds)))
						return nil
					})
					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}
	dWg.Wait()
	for i := 0; i < ncpu; i++ {
		rWg.Add(1)
		go func() {
			defer rWg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case rs := <-relationships:
					if rs == nil {
						return
					}
					err := svc.store.Update(func(batch Batch) error {
						svc.indexRelationships(batch, rs)
						atomic.AddUint32(&nr, uint32(len(rs)))
						return nil
					})
					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}
	rWg.Wait()
	for i := 0; i < ncpu; i++ {
		rsWg.Add(1)
		go func() {
			defer rsWg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case rs := <-refsetItems:
					if rs == nil {
						return
					}
					err := svc.store.Update(func(batch Batch) error {
						svc.indexRefsetItems(batch, rs)
						atomic.AddUint32(&nri, uint32(len(rs)))
						return nil
					})
					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}
	rsWg.Wait()
	close(done)
	// and now we have finished indexing, let's build search index
	if verbose {
		fmt.Printf("\nBuilding search index...\n")
	}

	if err := svc.buildSearchIndices(ctx, verbose); err != nil {
		return err
	}
	var err error
	if svc.availableLanguages, err = svc.AvailableLanguages(); err != nil { // refresh list of available languages
		return err
	}
	if verbose {
		fmt.Printf("\nPrecomputations complete. Total time: %s\n", time.Since(start))
	}
	return nil
}
func (svc *Svc) buildSearchIndices(ctx context.Context, verbose bool) error {
	tags, _, err := language.ParseAcceptLanguage("en-GB") // TODO: better language handling for search index
	if err != nil {
		return err
	}
	batchSize := 10000
	var total int64
	start := time.Now()
	eds := svc.iterateExtendedDescriptions(ctx, tags)
	var wg sync.WaitGroup
	for i := 0; i <= runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := make([]*snomed.ExtendedDescription, 0)
			for ed := range eds {
				batch = append(batch, ed.ExtendedDescription)
				if len(batch) == batchSize {
					atomic.AddInt64(&total, int64(batchSize))
					if verbose {
						elapsed := time.Since(start)
						count := atomic.LoadInt64(&total)
						fmt.Fprintf(os.Stderr, "\rSearch index: processed %d descriptions in %s. Mean time per description: %s...", count, elapsed, elapsed/time.Duration(count))
					}
					if err := svc.search.Index(batch); err != nil {
						panic(err)
					}
					batch = make([]*snomed.ExtendedDescription, 0)
				}
			}
			if err = svc.search.Index(batch); err != nil {
				panic(err)
			}
			atomic.AddInt64(&total, int64(len(batch)))
		}()
	}
	wg.Wait()
	fmt.Fprintf(os.Stderr, "\nProcessed total: %d descriptions in %s.\n", atomic.LoadInt64(&total), time.Since(start))
	return nil
}
