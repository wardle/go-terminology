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
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/wardle/go-terminology/snomed"
)

// boltService is a concrete file-based database service for SNOMED-CT
type boltService struct {
	db *bolt.DB
}

// Bucket structure
var (
	// Root buckets
	rbkConcepts      = []byte("Concepts")      // root bucket, containing concepts, keyed by id
	rbkDescriptions  = []byte("Descriptions")  // root bucket, containing descriptions, keyed by id
	rbkProperties    = []byte("Properties")    // root bucket, holding subbuckets named <conceptID> containing subbuckets (e.g. descriptions) containing all descriptions for that concept
	rbkReferenceSets = []byte("ReferenceSets") // root bucket, containing nested buckets named <refsetID> containing the items within that refset

	// Nested buckets "Properties"->"[conceptID]"->Bucket
	nbkParentRelationships = []byte("ParentRelationships") // nested bucket, containing parent relationships for this concept
	nbkChildRelationships  = []byte("ChildRelationships")  // nested bucket, containing child relationships for this concept
	nbkDescriptions        = []byte("Descriptions")        // nested bucket, containing descriptions for this concept
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
	case []*snomed.ReferenceSetItem:
		err = bs.putReferenceSets(components.([]*snomed.ReferenceSetItem))
	default:
		err = fmt.Errorf("unknown component type: %T", components)
	}
	return err
}

// GetConcept fetches a concept with the given identifier
func (bs *boltService) GetConcept(conceptID int64) (*snomed.Concept, error) {
	var c snomed.Concept
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(rbkConcepts)
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %s", rbkConcepts)
		}
		return mustReadFromBucket(bucket, conceptID, &c)
	})
	return &c, err
}

// GetConcepts returns a list of concepts with the given identifiers
func (bs *boltService) GetConcepts(conceptIDs ...int64) ([]*snomed.Concept, error) {
	result := make([]*snomed.Concept, len(conceptIDs))
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(rbkConcepts)
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %s", rbkConcepts)
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
		bucket, err := tx.CreateBucketIfNotExists(rbkConcepts)
		if err != nil {
			return err
		}
		for _, c := range concepts {
			if err = writeToBuckets(c.Id, c, bucket); err != nil {
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
		rootBucket, err := tx.CreateBucketIfNotExists(rbkDescriptions)
		if err != nil {
			return err
		}
		propsBucket, err := tx.CreateBucketIfNotExists(rbkProperties)
		if err != nil {
			return err
		}
		for _, d := range descriptions {
			conceptBucket, err := propsBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(int(d.ConceptId))))
			if err != nil {
				return err
			}
			descriptionsBucket, err := conceptBucket.CreateBucketIfNotExists(nbkDescriptions)
			if err != nil {
				return nil
			}
			if err := writeToBuckets(d.Id, d, descriptionsBucket, rootBucket); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetDescription returns the description with the given identifier
func (bs *boltService) GetDescription(descriptionID int64) (*snomed.Description, error) {
	var c snomed.Description
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(rbkDescriptions)
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %s", rbkDescriptions)
		}
		return mustReadFromBucket(bucket, descriptionID, &c)
	})
	return &c, err
}

// GetDescriptions returns the descriptions for this concept.
func (bs *boltService) GetDescriptions(concept *snomed.Concept) ([]*snomed.Description, error) {
	result := make([]*snomed.Description, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket, err := getPropertiesBucket(tx, concept.Id, nbkDescriptions)
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			var o snomed.Description
			err := proto.Unmarshal(v, &o)
			if err != nil {
				return err
			}
			result = append(result, &o)
			return nil
		})
		return nil
	})
	return result, err
}

// GetReferenceSets returns the refset identifiers to which this component is a member
func (bs *boltService) GetReferenceSets(referencedComponentID int64) ([]int64, error) {
	componentID := []byte(strconv.FormatInt(referencedComponentID, 10))
	result := make([]int64, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		referenceBucket := tx.Bucket([]byte(rbkReferenceSets))
		referenceBucket.ForEach(func(k, v []byte) error {
			refsetBucket := referenceBucket.Bucket(k)
			data := refsetBucket.Get(componentID)
			if data != nil {
				id, err := strconv.ParseInt(string(k), 10, 64)
				if err != nil {
					return err
				}
				result = append(result, id)
			}
			return nil
		})
		return nil
	})
	return result, err
}

// GetChildRelationships returns the child relationships for this concept.
// Child relationships are relationships in which this concept is the destination.
func (bs *boltService) GetChildRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error) {
	return bs.getRelationships(concept.Id, nbkChildRelationships)
}

// GetParentRelationships returns the parent relationships for this concept.
// Parent relationships are relationships in which this concept is the source.
func (bs *boltService) GetParentRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error) {
	return bs.getRelationships(concept.Id, nbkParentRelationships)
}

// getRelationships returns relationships using the specified property key.
func (bs *boltService) getRelationships(conceptID int64, key []byte) ([]*snomed.Relationship, error) {
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
			if err := proto.Unmarshal(v, &o); err != nil {
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
		propsBucket, err := tx.CreateBucketIfNotExists(rbkProperties)
		if err != nil {
			return err
		}
		for _, r := range relationships {
			sourceBucket, err := propsBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(int(r.SourceId))))
			if err != nil {
				return err
			}
			targetBucket, err := propsBucket.CreateBucketIfNotExists([]byte(strconv.Itoa(int(r.DestinationId))))
			if err != nil {
				return err
			}
			sParents, err := sourceBucket.CreateBucketIfNotExists(nbkParentRelationships)
			if err != nil {
				return err
			}
			sChildren, err := targetBucket.CreateBucketIfNotExists(nbkChildRelationships)
			if err != nil {
				return err
			}
			if err := writeToBuckets(r.Id, r, sParents, sChildren); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bs *boltService) putReferenceSets(refset []*snomed.ReferenceSetItem) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		referenceBucket, err := tx.CreateBucketIfNotExists(rbkReferenceSets)
		if err != nil {
			return err
		}
		for _, item := range refset {
			refsetID := []byte(strconv.FormatInt(item.GetRefsetId(), 10))
			referencedComponentID := []byte(strconv.FormatInt(item.GetReferencedComponentId(), 10))
			data, err := proto.Marshal(item)
			if err != nil {
				return err
			}
			refSetBucket, err := referenceBucket.CreateBucketIfNotExists(refsetID) // bucket for individual reference set
			if err != nil {
				return err
			}
			if err := refSetBucket.Put(referencedComponentID, data); err != nil {
				return err
			}
		}
		return nil
	})
}

// getPropertiesBucket returns the bucket holding properties for the concept specified, may be nil without an error!
func getPropertiesBucket(tx *bolt.Tx, conceptID int64, key []byte) (*bolt.Bucket, error) {
	propsBucket := tx.Bucket(rbkProperties)
	if propsBucket == nil {
		return nil, fmt.Errorf("missing bucket %s", rbkProperties)
	}
	conceptBucket := propsBucket.Bucket([]byte(strconv.FormatInt(conceptID, 10)))
	if conceptBucket == nil {
		return nil, nil
	}
	return conceptBucket.Bucket([]byte(key)), nil
}

// read an object from a bucket, returning nil and not initialising the structure if not found.
func readFromBucket(bucket *bolt.Bucket, id int, o proto.Message) error {
	key := []byte(strconv.Itoa(id))
	data := bucket.Get(key)
	if data == nil {
		return nil
	}
	return proto.Unmarshal(data, o)
}

// read an object from a bucket, throwing an error if not found
func mustReadFromBucket(bucket *bolt.Bucket, id int64, o proto.Message) error {
	key := []byte(strconv.FormatInt(id, 10))
	data := bucket.Get(key)
	if data == nil {
		return fmt.Errorf("no object found with identifier %d", id)
	}
	return proto.Unmarshal(data, o)
}

// helper method to write an object into multiple buckets
func writeToBuckets(id int64, o proto.Message, buckets ...*bolt.Bucket) error {
	data, err := proto.Marshal(o)
	if err != nil {
		return err
	}
	key := []byte(strconv.FormatInt(id, 10))
	for _, b := range buckets {
		if err := b.Put(key, data); err != nil {
			return err
		}
	}
	return nil
}

// Close releases all database resources.
func (bs *boltService) Close() error {
	return bs.db.Close()
}

func (bs *boltService) GetReferenceSetItems(refset int64) (map[int64]bool, error) {
	refsetID := []byte(strconv.FormatInt(refset, 10))
	result := make(map[int64]bool)
	err := bs.db.View(func(tx *bolt.Tx) error {
		referenceBucket := tx.Bucket([]byte(rbkReferenceSets))
		if referenceBucket == nil {
			return fmt.Errorf("no bucket found for refsets")
		}
		bucket := referenceBucket.Bucket(refsetID)
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %d", refset)
		}
		err := bucket.ForEach(func(k, v []byte) error {
			id, err := strconv.ParseInt(string(k), 10, 64)
			if err != nil {
				return err
			}
			result[id] = true
			return nil
		})
		return err
	})
	return result, err
}

// GetFromReferenceSet gets the specified components from the specified refset, or error
func (bs *boltService) GetFromReferenceSet(refset int64, component int64) (*snomed.ReferenceSetItem, error) {
	var result snomed.ReferenceSetItem
	found := false
	err := bs.db.View(func(tx *bolt.Tx) error {
		referenceBucket := tx.Bucket([]byte(rbkReferenceSets))
		if referenceBucket == nil {
			return fmt.Errorf("no bucket found to store refsets")
		}
		bucket := referenceBucket.Bucket([]byte(strconv.Itoa(int(refset))))
		if bucket == nil {
			return fmt.Errorf("refset %d not installed", refset)
		}
		if err := mustReadFromBucket(bucket, component, &result); err == nil {
			found = true
		}
		return nil
	})
	if found {
		return &result, err
	}
	return nil, err
}

// GetAllReferenceSets returns a list of installed reference sets
func (bs *boltService) GetAllReferenceSets() ([]int64, error) {
	result := make([]int64, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		referenceBucket := tx.Bucket([]byte(rbkReferenceSets))
		if referenceBucket != nil {
			referenceBucket.ForEach(func(k, v []byte) error {
				id, err := strconv.ParseInt(string(k), 10, 64)
				if err != nil {
					return err
				}
				result = append(result, id)
				return nil
			})
		}
		return nil
	})
	return result, err
}

// GetAllChildrenIDs returns the recursive children for this concept.
// This is a potentially large number, depending on where in the hierarchy the concept sits.
// TODO(mw): change to use transitive closure table
func (bs *boltService) GetAllChildrenIDs(concept *snomed.Concept) ([]int64, error) {
	allChildren := make(map[int64]bool)
	err := bs.recursiveChildren(concept.Id, allChildren)
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
	children, err := bs.getRelationships(conceptID, nbkChildRelationships)
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

// Iterate is a crude iterator for all concepts, useful for pre-processing and pre-computations
func (bs *boltService) Iterate(fn func(*snomed.Concept) error) error {
	return bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(rbkConcepts))
		var concept snomed.Concept
		return bucket.ForEach(func(k, v []byte) error {
			if err := proto.Unmarshal(v, &concept); err != nil {
				return err
			}
			return fn(&concept)
		})
	})
}

// GetStatistics returns statistics for the backend store
// This is crude and inefficient at the moment
// TODO(wardle): improve efficiency and speed
func (bs *boltService) GetStatistics() (Statistics, error) {
	stats := Statistics{}
	refsetNames := make([]string, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		// concepts
		cBucket := tx.Bucket([]byte(rbkConcepts))
		stats.concepts = cBucket.Stats().KeyN
		// descriptions
		dBucket := tx.Bucket([]byte(rbkDescriptions))
		stats.descriptions = dBucket.Stats().KeyN

		// reference sets
		rs := tx.Bucket([]byte(rbkReferenceSets))
		stats.refsetItems = rs.Stats().KeyN
		refsets := make([]int64, 0)
		c := rs.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v == nil { // if value is nil, then we have a subbucket
				id, err := strconv.ParseInt(string(k), 10, 64)
				if err != nil {
					return err
				}
				refsets = append(refsets, id)
			}
		}
		concepts, err := bs.GetConcepts(refsets...)
		if err != nil {
			return err
		}
		for _, c := range concepts {
			descs, err := bs.GetDescriptions(c)
			if err != nil {
				return err
			}
			if len(descs) > 0 {
				refsetName := fmt.Sprintf("%s (%d)", descs[0].Term, c.Id)
				refsetNames = append(refsetNames, refsetName)
			}
		}
		return err
	})
	stats.refsets = refsetNames
	return stats, err
}
