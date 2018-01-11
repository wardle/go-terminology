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
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/wardle/go-terminology/snomed"
)

// BoltService is a concrete file-based database service for SNOMED-CT
type BoltService struct {
	db *bolt.DB
}

// Bucket names
const (
	bkConcepts            = "Concepts"            // concept by conceptID
	bkDescriptions        = "Descriptions"        // description by descriptionID - within bucket of each concept
	bkParentRelationships = "ParentRelationships" // parent relationships by conceptID - within bucket of each concept
	bkChildRelationships  = "ChildRelationships"  // child relationships by conceptID - within bucket of each concept
	bkReferenceSets       = "ReferenceSets"       // reference sets by ID with subbuckets for the refsetID
)

// this is to ensure that, at compile-time, our database service is a valid implementation of a persistence store
var _ Store = (*BoltService)(nil)

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
func NewBoltService(filename string, readOnly bool) (*BoltService, error) {
	options := defaultOptions
	if readOnly {
		options = readOnlyOptions
	}
	db, err := bolt.Open(filename, 0644, options)
	if err != nil {
		return nil, err
	}
	return &BoltService{db: db}, nil
}

// GetBoltDB returns the underlying bolt database
func (bs *BoltService) GetBoltDB() *bolt.DB {
	return bs.db
}

// GetConcepts returns a list of concepts with the given identifiers
func (bs *BoltService) GetConcepts(conceptIDs ...int) ([]*snomed.Concept, error) {
	l := len(conceptIDs)
	result := make([]*snomed.Concept, 0, l)
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bkConcepts))
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %s", bkConcepts)
		}
		for _, conceptID := range conceptIDs {
			var concept snomed.Concept
			err := readFromBucket(bucket, conceptID, &concept)
			if err != nil {
				return err
			}
			result = append(result, &concept)
		}
		return nil
	})
	return result, err
}

// GetConcept fetches a concept with the given identifier
func (bs *BoltService) GetConcept(conceptID int) (*snomed.Concept, error) {
	var concept snomed.Concept
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bkConcepts))
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %s", bkConcepts)
		}
		return readFromBucket(bucket, conceptID, &concept)
	})
	if err != nil {
		return nil, err
	}
	return &concept, nil
}

// PutConcepts persists the specified concepts
func (bs *BoltService) PutConcepts(concepts []*snomed.Concept) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bkConcepts))
		if err != nil {
			return err
		}
		for _, c := range concepts {
			err = writeToBucket(bucket, int(c.ID), c)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// PutDescriptions persists the specified descriptions
func (bs *BoltService) PutDescriptions(descriptions []*snomed.Description) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		for _, d := range descriptions {
			cBucket, err := tx.CreateBucketIfNotExists([]byte(strconv.Itoa(int(d.ConceptID)))) // concept bucket
			if err != nil {
				return err
			}
			dBucket, err := cBucket.CreateBucketIfNotExists([]byte(bkDescriptions)) // create descriptions sub-bucket
			if err != nil {
				return err
			}
			err = writeToBucket(dBucket, int(d.ID), d)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// PutRelationships persists the specified relations
// TODO(mw): add more optimisations and precaching for each relationship
// note: this duplicates the relationship, possibly optimising walking the hierarchies
// at the expense of disk and memory usage
// TODO(mw): prove this premature optimisation actually works, rather than normalising
// and simply tracking the identifiers and then doing separate lookups...
func (bs *BoltService) PutRelationships(relationships []*snomed.Relationship) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		for _, r := range relationships {
			source := []byte(strconv.Itoa(int(r.SourceID)))
			target := []byte(strconv.Itoa(int(r.DestinationID)))
			sBucket, err := tx.CreateBucketIfNotExists(source) // bucket for source concept
			tBucket, err := tx.CreateBucketIfNotExists(target) // bucket for target concept
			sParents, err := sBucket.CreateBucketIfNotExists([]byte(bkParentRelationships))
			tChildren, err := tBucket.CreateBucketIfNotExists([]byte(bkChildRelationships))
			err = writeToBuckets(int(r.ID), r, sParents, tChildren)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// helper method to read an object from a bucket
func readFromBucket(bucket *bolt.Bucket, id int, o interface{}) error {
	key := []byte(strconv.Itoa(id))
	data := bucket.Get(key)
	if data == nil {
		return fmt.Errorf("No object found with identifier %d", id)
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(o)
}

// helper method to write an object into a single bucket
func writeToBucket(bucket *bolt.Bucket, id int, o interface{}) error {
	return writeToBuckets(id, o, bucket)
}

// helper method to write an object into one or more buckets.
func writeToBuckets(id int, o interface{}, buckets ...*bolt.Bucket) error {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(o)
	if err != nil {
		return err
	}
	key := []byte(strconv.Itoa(id))
	for _, b := range buckets {
		err = b.Put(key, buf.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

// Close releases all database resources.
func (bs *BoltService) Close() error {
	return bs.db.Close()
}

// GetDescriptions returns the descriptions for this concept.
func (bs *BoltService) GetDescriptions(concept *snomed.Concept) ([]*snomed.Description, error) {
	all := make([]*snomed.Description, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		cBkt := tx.Bucket([]byte(strconv.Itoa(int(concept.ID)))) // get individual concept bucket
		if cBkt == nil {
			return fmt.Errorf("No bucket found for concept %d", concept.ID)
		}
		if dBkt := cBkt.Bucket([]byte(bkDescriptions)); dBkt != nil {
			dBkt.ForEach(func(k, v []byte) error {
				var d snomed.Description
				buf := bytes.NewBuffer(v)
				dec := gob.NewDecoder(buf)
				dec.Decode(&d)
				all = append(all, &d)
				return nil
			})
		}
		return nil
	})
	return all, err
}

// GetChildRelationships returns the child relationships for this concept.
// Child relationships are relationships in which this concept is the destination.
func (bs *BoltService) GetChildRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error) {
	return bs.getRelationships(concept, []byte(bkChildRelationships))
}

// GetParentRelationships returns the parent relationships for this concept.
// Parent relationships are relationships in which this concept is the source.
func (bs *BoltService) GetParentRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error) {
	return bs.getRelationships(concept, []byte(bkParentRelationships))
}

// GetAllChildrenIDs returns the recursive children for this concept.
// This is a potentially large number, depending on where in the hierarchy the concept sits.
// TODO(mw): change to use transitive closure table
func (bs *BoltService) GetAllChildrenIDs(concept *snomed.Concept) ([]int, error) {
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
// TODO: use transitive closure precached table a la java version
func (bs *BoltService) recursiveChildren(conceptID int, allChildren map[int]bool) error {
	children, err := bs.getRelationshipsByID(conceptID, []byte(bkChildRelationships))
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

// helper method to get either parent or child relationships for a concept
func (bs *BoltService) getRelationships(concept *snomed.Concept, bucket []byte) ([]*snomed.Relationship, error) {
	return bs.getRelationshipsByID(int(concept.ID), bucket)
}

// helper method to get either parent or child relationships for a concept

func (bs *BoltService) getRelationshipsByID(conceptID int, bucket []byte) ([]*snomed.Relationship, error) {
	all := make([]*snomed.Relationship, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		cBkt := tx.Bucket([]byte(strconv.Itoa(conceptID))) // get individual concept bucket
		rBkt := cBkt.Bucket(bucket)
		if rBkt == nil { // if there is no bucket, then there are no relationships
			return nil
		}
		rBkt.ForEach(func(k, v []byte) error {
			var r snomed.Relationship
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)
			dec.Decode(&r)
			all = append(all, &r)
			return nil
		})
		return nil
	})
	return all, err
}
