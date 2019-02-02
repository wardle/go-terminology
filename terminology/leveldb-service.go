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
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/wardle/go-terminology/snomed"
)

// levelService is a concrete file-based database service for SNOMED-CT using goleveldb
type levelService struct {
	db *leveldb.DB
}

func (ls *levelService) put(batch *leveldb.Batch, b bucket, key []byte, value proto.Message) error {
	d, err := proto.Marshal(value)
	if err != nil {
		return err
	}
	k := bytes.Join([][]byte{b.name(), key}, nil)
	batch.Put(k, d)
	return nil
}

func (ls *levelService) get(b bucket, key []byte, pb proto.Message) error {
	d, err := ls.db.Get(bytes.Join([][]byte{b.name(), key}, nil), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return ErrNotFound
		}
		return err
	}
	return proto.Unmarshal(d, pb)
}

func (ls *levelService) addIndexEntry(batch *leveldb.Batch, b bucket, key []byte, value []byte) {
	k := bytes.Join([][]byte{b.name(), key, value}, nil)
	batch.Put(k, []byte{'.'})
}

func (ls *levelService) getIndexEntries(b bucket, key []byte) ([][]byte, error) {
	prefix := bytes.Join([][]byte{b.name(), key}, nil)
	lp := len(prefix)
	iter := ls.db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	result := make([][]byte, 0)
	for iter.Next() {
		k := iter.Key()
		entry := k[lp:]
		entry2 := make([]byte, len(entry))
		copy(entry2, entry)
		result = append(result, entry2) // we have to store a copy
	}
	return result, iter.Error()
}

// NewBoltService creates a new service at the specified location
func newLevelService(filename string, readOnly bool) (*levelService, error) {
	opts := opt.Options{ReadOnly: readOnly}
	db, err := leveldb.OpenFile(filename, &opts)
	if err != nil {
		return nil, err
	}
	return &levelService{db: db}, nil
}

// Put a slice of SNOMED-CT components into persistent storage.
// This is polymorphic but expects a slice of SNOMED CT components
func (ls *levelService) Put(components interface{}) error {
	var err error
	switch components.(type) {
	case []*snomed.Concept:
		err = ls.putConcepts(components.([]*snomed.Concept))
	case []*snomed.Description:
		err = ls.putDescriptions(components.([]*snomed.Description))
	case []*snomed.Relationship:
		err = ls.putRelationships(components.([]*snomed.Relationship))
	case []*snomed.ReferenceSetItem:
		err = ls.putReferenceSets(components.([]*snomed.ReferenceSetItem))
	default:
		err = fmt.Errorf("unknown component type: %T", components)
	}
	return err
}

// Concept returns the concept with the given identifier
func (ls *levelService) Concept(conceptID int64) (*snomed.Concept, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(conceptID))
	var c snomed.Concept
	err := ls.get(bkConcepts, key, &c)
	return &c, err
}

// GetConcepts returns a list of concepts with the given identifiers
func (ls *levelService) Concepts(conceptIDs ...int64) ([]*snomed.Concept, error) {
	key := make([]byte, 8)
	l := len(conceptIDs)
	r1 := make([]snomed.Concept, l)
	r2 := make([]*snomed.Concept, l)
	for i, id := range conceptIDs {
		binary.BigEndian.PutUint64(key, uint64(id))
		if err := ls.get(bkConcepts, key, &r1[i]); err != nil {
			return nil, err
		}
		r2[i] = &r1[i]
	}
	return r2, nil
}

// putConcepts persists the specified concepts
func (ls *levelService) putConcepts(concepts []*snomed.Concept) error {
	key := make([]byte, 8)
	batch := new(leveldb.Batch)
	for _, c := range concepts {
		binary.BigEndian.PutUint64(key, uint64(c.Id))
		if err := ls.put(batch, bkConcepts, key, c); err != nil {
			return err
		}
	}
	return ls.db.Write(batch, nil)
}

// PutDescriptions persists the specified descriptions
func (ls *levelService) putDescriptions(descriptions []*snomed.Description) error {
	dID := make([]byte, 8)
	cID := make([]byte, 8)
	batch := new(leveldb.Batch)
	for _, d := range descriptions {
		binary.BigEndian.PutUint64(dID, uint64(d.Id))
		if err := ls.put(batch, bkDescriptions, dID, d); err != nil {
			return err
		}
		binary.BigEndian.PutUint64(cID, uint64(d.ConceptId))
		ls.addIndexEntry(batch, ixConceptDescriptions, cID, dID)
	}
	return ls.db.Write(batch, nil)
}

// Description returns the description with the given identifier
func (ls *levelService) Description(descriptionID int64) (*snomed.Description, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(descriptionID))
	var c snomed.Description
	err := ls.get(bkDescriptions, key, &c)
	return &c, err
}

// Descriptions returns the descriptions for a concept
func (ls *levelService) Descriptions(conceptID int64) (result []*snomed.Description, err error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(conceptID))
	values, err := ls.getIndexEntries(ixConceptDescriptions, key)
	if err != nil {
		return nil, err
	}
	l := len(values)
	descs := make([]snomed.Description, l)
	result = make([]*snomed.Description, l)
	for i, v := range values {
		if err := ls.get(bkDescriptions, v, &descs[i]); err != nil {
			return nil, err
		}
		result[i] = &descs[i]
	}
	return
}

// PutRelationship persists the specified relationship
func (ls *levelService) putRelationships(relationships []*snomed.Relationship) error {
	rID := make([]byte, 8)
	sourceID := make([]byte, 8)
	destinationID := make([]byte, 8)
	batch := new(leveldb.Batch)
	for _, r := range relationships {
		binary.BigEndian.PutUint64(rID, uint64(r.Id))
		binary.BigEndian.PutUint64(sourceID, uint64(r.SourceId))
		binary.BigEndian.PutUint64(destinationID, uint64(r.DestinationId))
		if err := ls.put(batch, bkRelationships, rID, r); err != nil {
			return err
		}
		ls.addIndexEntry(batch, ixConceptParentRelationships, sourceID, rID)
		ls.addIndexEntry(batch, ixConceptChildRelationships, destinationID, rID)
	}
	return ls.db.Write(batch, nil)
}

// ChildRelationships returns the child relationships for this concept.
// Child relationships are relationships in which this concept is the destination.
func (ls *levelService) ChildRelationships(conceptID int64) ([]*snomed.Relationship, error) {
	return ls.getRelationships(conceptID, ixConceptChildRelationships)
}

// GetParentRelationships returns the parent relationships for this concept.
// Parent relationships are relationships in which this concept is the source.
func (ls *levelService) ParentRelationships(conceptID int64) ([]*snomed.Relationship, error) {
	return ls.getRelationships(conceptID, ixConceptParentRelationships)
}

// getRelationships returns relationships using the specified bucket
func (ls *levelService) getRelationships(conceptID int64, idx bucket) ([]*snomed.Relationship, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(conceptID))
	entries, err := ls.getIndexEntries(idx, key)
	if err != nil {
		return nil, err
	}
	l := len(entries)
	relationships := make([]snomed.Relationship, l)
	result := make([]*snomed.Relationship, l)
	for i, id := range entries {
		if err := ls.get(bkRelationships, id, &relationships[i]); err != nil {
			return nil, err
		}
		result[i] = &relationships[i]
	}
	return result, err
}

func (ls *levelService) putReferenceSets(refset []*snomed.ReferenceSetItem) error {
	referencedComponentID := make([]byte, 8)
	refsetID := make([]byte, 8)
	batch := new(leveldb.Batch)
	for _, item := range refset {
		itemID := []byte(item.Id)
		if err := ls.put(batch, bkRefsetItems, itemID, item); err != nil {
			return nil
		}
		binary.BigEndian.PutUint64(referencedComponentID, uint64(item.ReferencedComponentId))
		binary.BigEndian.PutUint64(refsetID, uint64(item.RefsetId))
		ls.addIndexEntry(batch, ixComponentReferenceSets, referencedComponentID, refsetID)
		ls.addIndexEntry(batch, ixReferenceSetComponentItems, compoundKey(refsetID, referencedComponentID), itemID)
		// support cross maps such as simple maps and complex maps
		var target string
		if simpleMap := item.GetSimpleMap(); simpleMap != nil {
			target = simpleMap.GetMapTarget()
		} else if complexMap := item.GetComplexMap(); complexMap != nil {
			target = complexMap.GetMapTarget()
		}
		if target != "" {
			ls.addIndexEntry(batch, ixRefsetTargetItems, compoundKey(refsetID, []byte(target+" ")), itemID)
		}
		// keep track of installed reference sets
		ls.addIndexEntry(batch, ixReferenceSets, refsetID, nil)
	}
	return ls.db.Write(batch, nil)
}

// ComponentReferenceSets returns the refset identifiers to which this component is a member
func (ls *levelService) ComponentReferenceSets(referencedComponentID int64) (result []int64, err error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(referencedComponentID))
	values, err := ls.getIndexEntries(ixComponentReferenceSets, key)
	if err != nil {
		return nil, err
	}
	result = make([]int64, len(values))
	for i, v := range values {
		result[i] = int64(binary.BigEndian.Uint64(v))
	}
	return
}

func (ls *levelService) MapTarget(refset int64, target string) (result []*snomed.ReferenceSetItem, err error) {
	refsetID := make([]byte, 8)
	binary.BigEndian.PutUint64(refsetID, uint64(refset))
	key := compoundKey(refsetID, []byte(target+" ")) // ensure delimiter between refset-target and value
	values, err := ls.getIndexEntries(ixRefsetTargetItems, key)
	if err != nil {
		return nil, err
	}
	l := len(values)
	items := make([]snomed.ReferenceSetItem, l)
	result = make([]*snomed.ReferenceSetItem, l)
	for i, v := range values {
		if err := ls.get(bkRefsetItems, v, &items[i]); err != nil {
			return nil, err
		}
		result[i] = &items[i]
	}
	return
}

// Close releases all database resources.
func (ls *levelService) Close() error {
	return ls.db.Close()
}

// ReferenceSetComponents returns the components within a given reference set
func (ls *levelService) ReferenceSetComponents(refset int64) (map[int64]struct{}, error) {
	result := make(map[int64]struct{})
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(refset))
	values, err := ls.getIndexEntries(ixReferenceSetComponentItems, key)
	if err != nil {
		return nil, err
	}
	for _, v := range values {
		result[int64(binary.BigEndian.Uint64(v[:8]))] = struct{}{}
	}
	return result, err
}

// ComponentFromReferenceSet gets the specified components from the specified refset, or error
func (ls *levelService) ComponentFromReferenceSet(refset int64, component int64) (result []*snomed.ReferenceSetItem, err error) {
	refsetID := make([]byte, 8)
	componentID := make([]byte, 8)
	binary.BigEndian.PutUint64(refsetID, uint64(refset))
	binary.BigEndian.PutUint64(componentID, uint64(component))
	key := compoundKey(refsetID, componentID)
	values, err := ls.getIndexEntries(ixReferenceSetComponentItems, key)
	if err != nil {
		return nil, err
	}
	l := len(values)
	items := make([]snomed.ReferenceSetItem, l)
	result = make([]*snomed.ReferenceSetItem, l)
	for i, v := range values {
		if err := ls.get(bkRefsetItems, v, &items[i]); err != nil {
			return nil, err
		}
		result[i] = &items[i]
	}
	return
}

// InstalledReferenceSets returns a list of installed reference sets
func (ls *levelService) InstalledReferenceSets() (map[int64]struct{}, error) {
	result := make(map[int64]struct{})
	iter := ls.db.NewIterator(util.BytesPrefix(ixReferenceSets.name()), nil)
	defer iter.Release()

	for iter.Next() {
		refsetID := int64(binary.BigEndian.Uint64(iter.Key()))
		result[refsetID] = struct{}{}
	}
	err := iter.Error()
	return result, err
}

// AllChildrenIDs returns the recursive children for this concept.
// This is a potentially large number, depending on where in the hierarchy the concept sits.
// TODO(mw): change to use transitive closure table
func (ls *levelService) AllChildrenIDs(conceptID int64, maximum int) ([]int64, error) {
	allChildren := make(map[int64]struct{})
	err := ls.recursiveChildren(conceptID, allChildren, maximum)
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
func (ls *levelService) recursiveChildren(conceptID int64, allChildren map[int64]struct{}, maximum int) error {
	if len(allChildren) > maximum {
		return fmt.Errorf("Too many children; aborted")
	}
	children, err := ls.getRelationships(conceptID, ixConceptChildRelationships)
	if err != nil {
		return err
	}
	for _, child := range children {
		if child.TypeId == snomed.IsA {
			childID := child.SourceId
			if _, exists := allChildren[childID]; !exists {
				if err := ls.recursiveChildren(childID, allChildren, maximum); err != nil {
					return err
				}
				allChildren[childID] = struct{}{}
			}
		}
	}
	return nil
}

// ClearPrecomputations clear precomputed indices
func (ls *levelService) ClearPrecomputations() error {
	return nil
}

// PerformPrecomputations builds indices that can only be performed after a complete import.
func (ls *levelService) PerformPrecomputations() error {
	return nil
}

// Iterate is a crude iterator for all concepts, useful for pre-processing and pre-computations
func (ls *levelService) Iterate(fn func(*snomed.Concept) error) error {
	iter := ls.db.NewIterator(util.BytesPrefix(bkConcepts.name()), nil)
	defer iter.Release()
	for iter.Next() {
		concept := snomed.Concept{}
		if err := proto.Unmarshal(iter.Value(), &concept); err != nil {
			return err
		}
		if err := fn(&concept); err != nil {
			return err
		}
	}
	return iter.Error()
}

// Statistics returns statistics for the backend store
func (ls *levelService) Statistics() (Statistics, error) {
	stats := Statistics{}
	return stats, nil
}
