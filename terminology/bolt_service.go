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
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"

	"github.com/boltdb/bolt"
	"github.com/wardle/go-terminology/snomed"
)

// boltService is a concrete file-based database service for SNOMED-CT
type boltService struct {
	db *bolt.DB
}

var defaultOptions = bolt.Options{
	Timeout:    10 * time.Second,
	NoGrowSync: false,
	ReadOnly:   false,
}

type bucket int

const (
	bkConcepts                   bucket = iota // concepts, keyed by SCTID (uint64)
	bkDescriptions                             // descriptions, keyed by SCTID (uint64)
	bkRelationships                            // relationships, keyed by SCTID (uint64)
	bkRefsetItems                              // refset items, keyed by their uuid (string)
	ixConceptDescriptions                      // key: concept_id-description_id
	ixConceptParentRelationships               // key: concept_id-relationship_id
	ixConceptChildRelationships                // key: concept_id-relationship_id
	ixComponentReferenceSets                   // key: component_id-refset_id
	ixReferenceSetComponentItems               // key: refset_id-component_id-reference_set_item_id
	ixReferenceSetItems                        // key: refset_id-reference_set_item_id
	ixRefsetTargetItems                        // key: refset_id-target_code-SPACE-reference_set_item_id
	ixConceptRecursiveParents                  // currently unused
	ixConceptDirectParents                     // currently unused
)

var bucketNames = [...][]byte{
	[]byte("con"), // key: sct_id value: concept
	[]byte("des"), // key: sct_id value: description
	[]byte("rel"), // key: sct_id value: relationship
	[]byte("ref"), // key: uuid value: component
	[]byte("cds"),
	[]byte("cpr"),
	[]byte("ccr"),
	[]byte("crs"),
	[]byte("rci"),
	[]byte("rsi"),
	[]byte("rti"),
	[]byte("crp"), // key: concept_id-parent_id
	[]byte("cdp"), // key: concept_id-parent_id
}

// ErrDatabaseNotInitialised is the error when database not properly initialised
var ErrDatabaseNotInitialised = errors.New("database not initialised")

// ErrNotFound is the error when something isn't found
var ErrNotFound = errors.New("Not found")

func (idx bucket) bucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket(bucketNames[idx])
}
func (idx bucket) createOrOpenBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	return tx.CreateBucketIfNotExists(bucketNames[idx])
}

func put(b *bolt.Bucket, key []byte, value proto.Message) error {
	d, err := proto.Marshal(value)
	if err != nil {
		return err
	}
	return b.Put(key, d)
}
func get(b *bolt.Bucket, key []byte, pb proto.Message) error {
	d := b.Get(key)
	if d == nil {
		return ErrNotFound
	}
	return proto.Unmarshal(d, pb)
}

func addIndexEntry(b *bolt.Bucket, key []byte, value []byte) error {
	k := bytes.Join([][]byte{key, value}, nil)
	return b.Put(k, nil)
}

func getIndexEntries(b *bolt.Bucket, key []byte) ([][]byte, error) {
	lp := len(key)
	c := b.Cursor()
	result := make([][]byte, 0)
	for k, _ := c.Seek(key); k != nil && bytes.HasPrefix(k, key); k, _ = c.Next() {
		result = append(result, k[lp:])
	}
	return result, nil
}

func compoundKey(keys ...[]byte) []byte {
	return bytes.Join(keys, nil)
}

// NewBoltService creates a new service at the specified location
func newBoltService(filename string, readOnly bool) (*boltService, error) {
	options := defaultOptions
	if readOnly {
		options.ReadOnly = true
	}
	db, err := bolt.Open(filename, 0600, &options)
	if err != nil {
		return nil, err
	}
	return &boltService{db: db}, nil
}

// Put a slice of SNOMED-CT components into persistent storage.
// This is polymorphic but expects a slice of SNOMED CT components
func (bs *boltService) Put(components interface{}) error {
	var err error
	switch components.(type) {
	case []*snomed.Concept:
		err = bs.putConcepts(components.([]*snomed.Concept))
	case []*snomed.Description:
		err = bs.putDescriptions(components.([]*snomed.Description))
	case []*snomed.Relationship:
		err = bs.putRelationships(components.([]*snomed.Relationship))
	case []*snomed.ReferenceSetItem:
		err = bs.putReferenceSets(components.([]*snomed.ReferenceSetItem))
	default:
		err = fmt.Errorf("unknown component type: %T", components)
	}
	return err
}

// Concept returns the concept with the given identifier
func (bs *boltService) Concept(conceptID int64) (*snomed.Concept, error) {
	key := make([]byte, 8)
	var c snomed.Concept
	err := bs.db.View(func(tx *bolt.Tx) error {
		concepts := bkConcepts.bucket(tx)
		if concepts == nil {
			return ErrDatabaseNotInitialised
		}
		binary.BigEndian.PutUint64(key, uint64(conceptID))
		return get(concepts, key, &c)
	})
	return &c, err
}

// GetConcepts returns a list of concepts with the given identifiers
func (bs *boltService) Concepts(conceptIDs ...int64) ([]*snomed.Concept, error) {
	key := make([]byte, 8)
	result := make([]*snomed.Concept, len(conceptIDs))
	err := bs.db.View(func(tx *bolt.Tx) error {
		concepts := bkConcepts.bucket(tx)
		if concepts == nil {
			return ErrDatabaseNotInitialised
		}
		for i, id := range conceptIDs {
			var c snomed.Concept
			binary.BigEndian.PutUint64(key, uint64(id))
			if err := get(concepts, key, &c); err != nil {
				return err
			}
			result[i] = &c
		}
		return nil
	})
	return result, err
}

// putConcepts persists the specified concepts
func (bs *boltService) putConcepts(concepts []*snomed.Concept) error {
	key := make([]byte, 8)
	return bs.db.Update(func(tx *bolt.Tx) error {
		bkConcepts, err := bkConcepts.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		for _, c := range concepts {
			binary.BigEndian.PutUint64(key, uint64(c.Id))
			if err := put(bkConcepts, key, c); err != nil {
				return err
			}
		}
		return nil
	})
}

// PutDescriptions persists the specified descriptions
func (bs *boltService) putDescriptions(descriptions []*snomed.Description) error {
	dID := make([]byte, 8)
	cID := make([]byte, 8)
	return bs.db.Update(func(tx *bolt.Tx) error {
		bd, err := bkDescriptions.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		conceptDescriptions, err := ixConceptDescriptions.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		for _, d := range descriptions {
			binary.BigEndian.PutUint64(dID, uint64(d.Id))
			if err := put(bd, dID, d); err != nil {
				return err
			}
			binary.BigEndian.PutUint64(cID, uint64(d.ConceptId))
			if err := addIndexEntry(conceptDescriptions, cID, dID); err != nil {
				return err
			}
		}
		return nil
	})
}

// Description returns the description with the given identifier
func (bs *boltService) Description(descriptionID int64) (*snomed.Description, error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(descriptionID))
	var c snomed.Description
	err := bs.db.View(func(tx *bolt.Tx) error {
		bd := bkDescriptions.bucket(tx)
		if bd == nil {
			return ErrDatabaseNotInitialised
		}
		return get(bd, key, &c)
	})
	return &c, err
}

// Descriptions returns the descriptions for a concept
func (bs *boltService) Descriptions(conceptID int64) (result []*snomed.Description, err error) {
	err = bs.db.View(func(tx *bolt.Tx) error {
		bd := bkDescriptions.bucket(tx)
		conceptDescriptions := ixConceptDescriptions.bucket(tx)
		if bd == nil || conceptDescriptions == nil {
			return ErrDatabaseNotInitialised
		}
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, uint64(conceptID))
		values, err := getIndexEntries(conceptDescriptions, key)
		if err != nil {
			return err
		}
		l := len(values)
		descs := make([]snomed.Description, l)
		result = make([]*snomed.Description, l)
		for i, v := range values {
			if err := get(bd, v, &descs[i]); err != nil {
				return err
			}
			result[i] = &descs[i]
		}
		return nil
	})
	return
}

// PutRelationship persists the specified relationship
func (bs *boltService) putRelationships(relationships []*snomed.Relationship) error {
	rID := make([]byte, 8)
	sourceID := make([]byte, 8)
	destinationID := make([]byte, 8)
	return bs.db.Update(func(tx *bolt.Tx) error {
		relBucket, err := bkRelationships.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		parentRelationships, err := ixConceptParentRelationships.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		childRelationships, err := ixConceptChildRelationships.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		for _, r := range relationships {
			binary.BigEndian.PutUint64(rID, uint64(r.Id))
			binary.BigEndian.PutUint64(sourceID, uint64(r.SourceId))
			binary.BigEndian.PutUint64(destinationID, uint64(r.DestinationId))
			if err := put(relBucket, rID, r); err != nil {
				return err
			}
			if err := addIndexEntry(parentRelationships, sourceID, rID); err != nil {
				return err
			}
			if err := addIndexEntry(childRelationships, destinationID, rID); err != nil {
				return err
			}
		}
		return nil
	})
}

// ChildRelationships returns the child relationships for this concept.
// Child relationships are relationships in which this concept is the destination.
func (bs *boltService) ChildRelationships(conceptID int64) ([]*snomed.Relationship, error) {
	return bs.getRelationships(conceptID, ixConceptChildRelationships)
}

// GetParentRelationships returns the parent relationships for this concept.
// Parent relationships are relationships in which this concept is the source.
func (bs *boltService) ParentRelationships(conceptID int64) ([]*snomed.Relationship, error) {
	return bs.getRelationships(conceptID, ixConceptParentRelationships)
}

// getRelationships returns relationships using the specified bucket
func (bs *boltService) getRelationships(conceptID int64, idx bucket) (result []*snomed.Relationship, err error) {
	err = bs.db.View(func(tx *bolt.Tx) error {
		relBucket := bkRelationships.bucket(tx)
		b := idx.bucket(tx)
		if relBucket == nil || b == nil {
			return ErrDatabaseNotInitialised
		}
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, uint64(conceptID))
		entries, err := getIndexEntries(b, key)
		if err != nil {
			return err
		}
		l := len(entries)
		relationships := make([]snomed.Relationship, l)
		result = make([]*snomed.Relationship, l)
		for i, id := range entries {
			if err := get(relBucket, id, &relationships[i]); err != nil {
				return err
			}
			result[i] = &relationships[i]
		}
		return nil
	})
	return
}

func (bs *boltService) putReferenceSets(refset []*snomed.ReferenceSetItem) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		rb, err := bkRefsetItems.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		componentReferenceSets, err := ixComponentReferenceSets.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		referenceSetComponentItems, err := ixReferenceSetComponentItems.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		referenceSetItems, err := ixReferenceSetItems.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		refsetTargetItems, err := ixRefsetTargetItems.createOrOpenBucket(tx)
		if err != nil {
			return err
		}
		referencedComponentID := make([]byte, 8)
		refsetID := make([]byte, 8)
		for _, item := range refset {
			itemID := []byte(item.Id)
			if err := put(rb, itemID, item); err != nil {
				return nil
			}
			binary.BigEndian.PutUint64(referencedComponentID, uint64(item.ReferencedComponentId))
			binary.BigEndian.PutUint64(refsetID, uint64(item.RefsetId))
			if err := addIndexEntry(componentReferenceSets, referencedComponentID, refsetID); err != nil {
				return err
			}
			if err := addIndexEntry(referenceSetComponentItems, compoundKey(refsetID, referencedComponentID), itemID); err != nil {
				return err
			}
			if err := addIndexEntry(referenceSetItems, refsetID, itemID); err != nil {
				return err
			}
			var target string
			if simpleMap := item.GetSimpleMap(); simpleMap != nil {
				target = simpleMap.GetMapTarget()
			} else if complexMap := item.GetComplexMap(); complexMap != nil {
				target = complexMap.GetMapTarget()
			}
			if target != "" {
				if err := addIndexEntry(refsetTargetItems, compoundKey(refsetID, []byte(target+" ")), itemID); err != nil {
					return err
				}
			}

		}
		return nil
	})
}

// ComponentReferenceSets returns the refset identifiers to which this component is a member
func (bs *boltService) ComponentReferenceSets(referencedComponentID int64) (result []int64, err error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(referencedComponentID))
	err = bs.db.View(func(tx *bolt.Tx) error {
		componentReferenceSets := ixComponentReferenceSets.bucket(tx)
		if componentReferenceSets == nil {
			return ErrDatabaseNotInitialised
		}
		values, err := getIndexEntries(componentReferenceSets, key)
		if err != nil {
			return err
		}
		result = make([]int64, len(values))
		for i, v := range values {
			result[i] = int64(binary.BigEndian.Uint64(v))
		}
		return nil
	})
	return
}

func (bs *boltService) MapTarget(refset int64, target string) (result []*snomed.ReferenceSetItem, err error) {
	refsetID := make([]byte, 8)
	binary.BigEndian.PutUint64(refsetID, uint64(refset))
	err = bs.db.View(func(tx *bolt.Tx) error {
		rib := bkRefsetItems.bucket(tx)
		refsetTargetItems := ixRefsetTargetItems.bucket(tx)
		if rib == nil || refsetTargetItems == nil {
			return ErrDatabaseNotInitialised
		}
		key := compoundKey(refsetID, []byte(target+" ")) // ensure delimiter between refset-target and value
		values, err := getIndexEntries(refsetTargetItems, key)
		if err != nil {
			return err
		}
		l := len(values)
		items := make([]snomed.ReferenceSetItem, l)
		result = make([]*snomed.ReferenceSetItem, l)
		for i, v := range values {
			if err := get(rib, v, &items[i]); err != nil {
				return err
			}
			result[i] = &items[i]
		}
		return nil
	})
	return
}

// Close releases all database resources.
func (bs *boltService) Close() error {
	return bs.db.Close()
}

// ReferenceSetComponents returns the components within a given reference set
func (bs *boltService) ReferenceSetComponents(refset int64) (map[int64]struct{}, error) {
	result := make(map[int64]struct{})
	err := bs.db.View(func(tx *bolt.Tx) error {
		components := ixReferenceSetComponentItems.bucket(tx)
		if components == nil {
			return ErrDatabaseNotInitialised
		}
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, uint64(refset))
		values, err := getIndexEntries(components, key)
		if err != nil {
			return err
		}
		for _, v := range values {
			result[int64(binary.BigEndian.Uint64(v[:8]))] = struct{}{}
		}
		return nil
	})
	return result, err
}

// ComponentFromReferenceSet gets the specified components from the specified refset, or error
func (bs *boltService) ComponentFromReferenceSet(refset int64, component int64) (result []*snomed.ReferenceSetItem, err error) {
	refsetID := make([]byte, 8)
	componentID := make([]byte, 8)
	binary.BigEndian.PutUint64(refsetID, uint64(refset))
	binary.BigEndian.PutUint64(componentID, uint64(component))
	err = bs.db.View(func(tx *bolt.Tx) error {
		refsetItemBucket := bkRefsetItems.bucket(tx)
		referenceSetComponentItems := ixReferenceSetComponentItems.bucket(tx)
		if refsetItemBucket == nil || referenceSetComponentItems == nil {
			return ErrDatabaseNotInitialised
		}
		key := compoundKey(refsetID, componentID)
		values, err := getIndexEntries(referenceSetComponentItems, key)
		if err != nil {
			return err
		}
		l := len(values)
		items := make([]snomed.ReferenceSetItem, l)
		result = make([]*snomed.ReferenceSetItem, l)
		for i, v := range values {
			if err := get(refsetItemBucket, v, &items[i]); err != nil {
				return err
			}
			result[i] = &items[i]
		}
		return nil
	})
	return
}

// InstalledReferenceSets returns a list of installed reference sets
func (bs *boltService) InstalledReferenceSets() ([]int64, error) {
	return bs.AllChildrenIDs(snomed.ReferenceSetConcept)
}

// AllChildrenIDs returns the recursive children for this concept.
// This is a potentially large number, depending on where in the hierarchy the concept sits.
// TODO(mw): change to use transitive closure table
func (bs *boltService) AllChildrenIDs(conceptID int64) ([]int64, error) {
	allChildren := make(map[int64]bool)
	err := bs.recursiveChildren(conceptID, allChildren)
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
func (bs *boltService) recursiveChildren(conceptID int64, allChildren map[int64]bool) error {
	children, err := bs.getRelationships(conceptID, ixConceptChildRelationships)
	if err != nil {
		return err
	}
	for _, child := range children {
		if child.TypeId == snomed.IsA {
			childID := child.SourceId
			if allChildren[childID] == false {
				allChildren[childID] = true
				if err != nil {
					return err
				}
				err = bs.recursiveChildren(childID, allChildren)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ClearPrecomputations clear precomputed indices
func (bs *boltService) ClearPrecomputations() error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(bucketNames[ixConceptDirectParents]); err != nil {
			return err
		}
		if err := tx.DeleteBucket(bucketNames[ixConceptRecursiveParents]); err != nil {
			return err
		}
		return nil
	})
}

// PerformPrecomputations builds indices that can only be performed after a complete import.
func (bs *boltService) PerformPrecomputations() error {
	return nil
}

// Iterate is a crude iterator for all concepts, useful for pre-processing and pre-computations
func (bs *boltService) Iterate(fn func(*snomed.Concept) error) error {
	return bs.db.View(func(tx *bolt.Tx) error {
		concepts := bkConcepts.bucket(tx)
		if concepts == nil {
			return ErrDatabaseNotInitialised
		}
		concepts.ForEach(func(k, v []byte) error {
			concept := snomed.Concept{}
			if err := proto.Unmarshal(v, &concept); err != nil {
				return err
			}
			fn(&concept)
			return nil
		})
		return nil
	})
}

// Statistics returns statistics for the backend store
func (bs *boltService) Statistics() (Statistics, error) {
	stats := Statistics{}
	refsetNames := make([]string, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		if concepts := bkConcepts.bucket(tx); concepts != nil {
			stats.concepts = concepts.Stats().KeyN
		}
		if descriptions := bkDescriptions.bucket(tx); descriptions != nil {
			stats.descriptions = descriptions.Stats().KeyN
		}
		if relationships := bkRelationships.bucket(tx); relationships != nil {
			stats.relationships = relationships.Stats().KeyN
		}
		if refsetItems := bkRefsetItems.bucket(tx); refsetItems != nil {
			stats.refsetItems = refsetItems.Stats().KeyN
		}
		return nil
	})

	refsets, err := bs.InstalledReferenceSets()
	if err != nil {
		return stats, err
	}
	for _, refset := range refsets {
		descs, err := bs.Descriptions(refset)
		if err != nil {
			return stats, err
		}
		if len(descs) > 0 {
			refsetName := fmt.Sprintf("%s (%d)", descs[0].Term, refset)
			refsetNames = append(refsetNames, refsetName)
		}
	}
	stats.refsets = refsetNames

	return stats, err
}
