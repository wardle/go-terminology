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
	"encoding/gob"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/wardle/go-terminology/snomed"
)

// boltService is a concrete file-based database service for SNOMED-CT
// TODO(wardle):switch to protobuf
type boltService struct {
	db *bolt.DB
}

// Bucket structure
const (
	// Root buckets
	bkConcepts      = "Concepts"      // root bucket, the main bucket for different types of SNOMED-CT component, holds each component keyed by ID
	bkProperties    = "Properties"    // root bucket, holding subbuckets named <conceptID> containing properties and values as slices of identifiers
	bkReferenceSets = "ReferenceSets" // root bucket, containing reference sets named <refsetID>

	// Nested buckets "Properties"->"[conceptID]"->Bucket
	bkParentRelationships = "ParentRelationships" // nested bucket, containing parent relationships for this concept
	bkChildRelationships  = "ChildRelationships"  // nested bucket, containing child relationships for this concept
	bkDescriptions        = "Descriptions"        // nested bucket, containing descriptions for this concept
)

// assert that, at compile-time, this database service is a valid implementation of a persistence store
var _ store = (*boltService)(nil)

var defaultOptions = &bolt.Options{
	Timeout:    0,
	NoGrowSync: false,
	ReadOnly:   false,
}
var readOnlyOptions = &bolt.Options{
	Timeout:    0,
	NoGrowSync: false,
	ReadOnly:   true,
}

// NewBoltService creates a new service at the specified location
func newBoltService(filename string, readOnly bool) (*boltService, error) {
	options := defaultOptions
	if readOnly {
		options = readOnlyOptions
	}
	db, err := bolt.Open(filename, 0644, options)
	if err != nil {
		return nil, err
	}
	return &boltService{db: db}, nil
}

// Put a slice of SNOMED-CT components into persistent storage.
// This is polymorphic but expects a slice of a core SNOMED CT component
func (bs *boltService) Put(components interface{}) error {
	var err error
	switch components.(type) {
	case []*snomed.Concept:
		err = bs.putConcepts(components.([]*snomed.Concept))
	case []*snomed.Description:
		err = bs.putDescriptions(components.([]*snomed.Description))
	case []*snomed.Relationship:
		err = bs.putRelationships(components.([]*snomed.Relationship))
	case []*snomed.LanguageReferenceSet:
		err = bs.putLanguageReferenceSets(components.([]*snomed.LanguageReferenceSet))
	default:
		err = fmt.Errorf("unknown component type: %T", components)
	}
	return err
}

// GetConcept fetches a concept with the given identifier
func (bs *boltService) GetConcept(conceptID int) (*snomed.Concept, error) {
	var c snomed.Concept
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bkConcepts))
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %s", bkConcepts)
		}
		return mustReadFromBucket(bucket, conceptID, &c)
	})
	return &c, err
}

// GetConcepts returns a list of concepts with the given identifiers
func (bs *boltService) GetConcepts(conceptIDs ...int) ([]*snomed.Concept, error) {
	result := make([]*snomed.Concept, len(conceptIDs))
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bkConcepts))
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %s", bkConcepts)
		}
		for i, id := range conceptIDs {
			var c snomed.Concept
			if err := mustReadFromBucket(bucket, id, &c); err != nil {
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
	return bs.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bkConcepts))
		if err != nil {
			return err
		}
		for _, c := range concepts {
			id := int(c.ID)
			if err = writeToBuckets(id, c, bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

// PutDescriptions persists the specified descriptions
// This 1) writes the description into generic components bucket and 2) adds the description id to the concept
func (bs *boltService) putDescriptions(descriptions []*snomed.Description) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		propsBucket, err := tx.CreateBucketIfNotExists([]byte(bkProperties))
		if err != nil {
			return err
		}
		for _, d := range descriptions {
			conceptBucket, err := propsBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(int(d.ConceptID))))
			if err != nil {
				return err
			}
			descriptionsBucket, err := conceptBucket.CreateBucketIfNotExists([]byte(bkDescriptions))
			if err != nil {
				return nil
			}
			if err := writeToBuckets(int(d.ID), d, descriptionsBucket); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetDescriptions returns the descriptions for this concept.
func (bs *boltService) GetDescriptions(concept *snomed.Concept) ([]*snomed.Description, error) {
	result := make([]*snomed.Description, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket, err := getPropertiesBucket(tx, int(concept.ID), bkDescriptions)
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			var o snomed.Description
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)
			if err := dec.Decode(&o); err != nil {
				return err
			}
			result = append(result, &o)
			return nil
		})
		return nil
	})
	return result, err
}

// GetChildRelationships returns the child relationships for this concept.
// Child relationships are relationships in which this concept is the destination.
func (bs *boltService) GetChildRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error) {
	return bs.getRelationships(int(concept.ID), bkChildRelationships)
}

// GetParentRelationships returns the parent relationships for this concept.
// Parent relationships are relationships in which this concept is the source.
func (bs *boltService) GetParentRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error) {
	return bs.getRelationships(int(concept.ID), bkParentRelationships)
}

// getRelationships returns relationships using the specified property key.
func (bs *boltService) getRelationships(conceptID int, key string) ([]*snomed.Relationship, error) {
	result := make([]*snomed.Relationship, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket, err := getPropertiesBucket(tx, conceptID, key)
		if err != nil {
			return err
		}
		if bucket == nil { // if we have no property bucket, then we have no relationships
			return nil
		}
		bucket.ForEach(func(k, v []byte) error {
			var o snomed.Relationship
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)
			if err := dec.Decode(&o); err != nil {
				return err
			}
			result = append(result, &o)
			return nil
		})
		return nil
	})
	return result, err
}

// PutRelationship persists the specified relationship
// TODO(mw): add more optimisations and precaching for each relationship
func (bs *boltService) putRelationships(relationships []*snomed.Relationship) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		propsBucket, err := tx.CreateBucketIfNotExists([]byte(bkProperties))
		if err != nil {
			return err
		}
		for _, r := range relationships {
			sourceBucket, err := propsBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(int(r.SourceID))))
			if err != nil {
				return err
			}
			targetBucket, err := propsBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(int(r.DestinationID))))
			if err != nil {
				return err
			}
			sParents, err := sourceBucket.CreateBucketIfNotExists([]byte(bkParentRelationships))
			if err != nil {
				return err
			}
			sChildren, err := targetBucket.CreateBucketIfNotExists([]byte(bkChildRelationships))
			if err != nil {
				return err
			}
			if err := writeToBuckets(int(r.ID), r, sParents, sChildren); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bs *boltService) putLanguageReferenceSets(refset []*snomed.LanguageReferenceSet) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		referenceBucket, err := tx.CreateBucketIfNotExists([]byte(bkReferenceSets))
		if err != nil {
			return err
		}
		for _, item := range refset {
			if err := bs.putReferenceSetItem(referenceBucket, item.RefsetID, item.ReferencedComponentID, item); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bs *boltService) putReferenceSetItem(bucket *bolt.Bucket, refsetID snomed.Identifier, referencedComponentID snomed.Identifier, item interface{}) error {
	refSetBucket, err := bucket.CreateBucketIfNotExists([]byte(strconv.Itoa(int(refsetID)))) // bucket for individual reference set
	if err != nil {
		return err
	}
	return writeToBuckets(int(referencedComponentID), item, refSetBucket) // add the referenced component keyed by referenced component ID
}

// getPropertiesBucket returns the bucket holding properties for the concept specified, may be nil without an error!
func getPropertiesBucket(tx *bolt.Tx, conceptID int, key string) (*bolt.Bucket, error) {
	propsBucket := tx.Bucket([]byte(bkProperties))
	if propsBucket == nil {
		return nil, fmt.Errorf("missing bucket %s", bkProperties)
	}
	conceptBucket := propsBucket.Bucket([]byte(strconv.Itoa(conceptID)))
	if conceptBucket == nil {
		return nil, nil
	}
	return conceptBucket.Bucket([]byte(key)), nil
}

// read an object from a bucket, returning nil and not initialising the structure if not found.
func readFromBucket(bucket *bolt.Bucket, id int, o interface{}) error {
	key := []byte(strconv.Itoa(id))
	data := bucket.Get(key)
	if data == nil {
		return nil
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(o)
}

// read an object from a bucket, throwing an error if not found
func mustReadFromBucket(bucket *bolt.Bucket, id int, o interface{}) error {
	key := []byte(strconv.Itoa(id))
	data := bucket.Get(key)
	if data == nil {
		return fmt.Errorf("no object found with identifier %d", id)
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(o)
}

// helper method to write an object into multiple buckets
func writeToBuckets(id int, o interface{}, buckets ...*bolt.Bucket) error {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(o); err != nil {
		return err
	}
	key := []byte(strconv.Itoa(id))
	for _, b := range buckets {
		if err := b.Put(key, buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// Close releases all database resources.
func (bs *boltService) Close() error {
	return bs.db.Close()
}

func (bs *boltService) GetReferenceSet(refset snomed.Identifier) (map[snomed.Identifier]bool, error) {
	refsetID := []byte(strconv.Itoa(int(refset)))
	result := make(map[snomed.Identifier]bool)
	err := bs.db.View(func(tx *bolt.Tx) error {
		referenceBucket := tx.Bucket([]byte(bkReferenceSets))
		if referenceBucket == nil {
			return fmt.Errorf("no bucket found for refsets")
		}
		bucket := referenceBucket.Bucket(refsetID)
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %d", refset)
		}
		err := bucket.ForEach(func(k, v []byte) error {
			id, err := strconv.Atoi(string(k))
			if err != nil {
				return err
			}
			result[snomed.Identifier(id)] = true
			return nil
		})
		return err
	})
	return result, err
}

// GetFromReferenceSet gets the specified components from the specified refset, or error
func (bs *boltService) GetFromReferenceSet(refset snomed.Identifier, component snomed.Identifier, result interface{}) (bool, error) {
	found := false
	err := bs.db.View(func(tx *bolt.Tx) error {
		referenceBucket := tx.Bucket([]byte(bkReferenceSets))
		if referenceBucket == nil {
			return fmt.Errorf("no bucket found to store refsets")
		}
		bucket := referenceBucket.Bucket([]byte(strconv.Itoa(int(refset))))
		if bucket == nil {
			return fmt.Errorf("refset %d not installed", refset)
		}
		if err := mustReadFromBucket(bucket, int(component), result); err == nil {
			found = true
		}
		return nil
	})
	return found, err
}

// GetReferenceSets returns a list of installed reference sets
func (bs *boltService) GetReferenceSets() ([]snomed.Identifier, error) {
	result := make([]snomed.Identifier, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		referenceBucket := tx.Bucket([]byte(bkReferenceSets))
		referenceBucket.ForEach(func(k, v []byte) error {
			id, err := strconv.Atoi(string(k))
			if err != nil {
				return err
			}
			result = append(result, snomed.Identifier(id))
			return nil
		})
		return nil
	})
	return result, err
}

// GetAllChildrenIDs returns the recursive children for this concept.
// This is a potentially large number, depending on where in the hierarchy the concept sits.
// TODO(mw): change to use transitive closure table
func (bs *boltService) GetAllChildrenIDs(concept *snomed.Concept) ([]int, error) {
	allChildren := make(map[int]bool)
	err := bs.recursiveChildren(int(concept.ID), allChildren)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(allChildren))
	for id := range allChildren {
		ids = append(ids, id)
	}
	return ids, nil
}

// this is a brute-force, non-cached temporary version which actually fetches the id
// TODO(mwardle): benchmark and possibly use transitive closure precached table a la java version
func (bs *boltService) recursiveChildren(conceptID int, allChildren map[int]bool) error {
	children, err := bs.getRelationships(conceptID, bkChildRelationships)
	if err != nil {
		return err
	}
	for _, child := range children {
		if child.TypeID == snomed.IsAConceptID {
			childID := int(child.SourceID)
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

// Iterate is a crude iterator for all concepts, useful for pre-processing and pre-computations
func (bs *boltService) Iterate(fn func(*snomed.Concept) error) error {
	return bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bkConcepts))
		var concept snomed.Concept
		return bucket.ForEach(func(k, v []byte) error {
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)
			if err := dec.Decode(&concept); err != nil {
				return err
			}
			return fn(&concept)
		})
	})
}
