package db

import (
	"bitbucket.org/wardle/go-snomed/rf2"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"strconv"
)

// BoltService is a concrete file-based database service for SNOMED-CT
type BoltService struct {
	db *bolt.DB
}

// Bucket names
const (
	bkConcepts            = "Concepts"              // concept by conceptID
	bkDescriptions        = "Descriptions"          // description by descriptionID
	bkRelationships       = "Relationships"         // relationship by relationshipID
	bkParentRelationships = "ParentRelationships"   // parent relationships by conceptID
	bkChildRelationships  = "ChildRelationships"    // child relationships by conceptID
	bkConceptDescriptions = "DescriptionsByConcept" // descriptions by conceptID
)

// this is to ensure that, at compile-time, our database service is a valid implementation of Service
//var _ Service = (*BoltService)(nil)

// NewBoltService creates a new service at the specified location
func NewBoltService(filename string) (*BoltService, error) {
	db, err := bolt.Open(filename, 0644, nil)
	if err != nil {
		return nil, err
	}
	return &BoltService{db: db}, nil
}

// GetConcepts returns a list of concepts with the given identifiers
func (bs *BoltService) GetConcepts(conceptIDs ...int) ([]*rf2.Concept, error) {
	l := len(conceptIDs)
	result := make([]*rf2.Concept, l)
	err := bs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bkConcepts))
		if bucket == nil {
			return fmt.Errorf("no bucket found with name: %s", bkConcepts)
		}
		for _, conceptID := range conceptIDs {
			var concept rf2.Concept
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
func (bs *BoltService) GetConcept(conceptID int) (*rf2.Concept, error) {
	var concept rf2.Concept
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
func (bs *BoltService) PutConcepts(concepts ...*rf2.Concept) error {
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

func writeToBucket(bucket *bolt.Bucket, id int, o interface{}) error {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(o)
	if err != nil {
		return err
	}
	key := []byte(strconv.Itoa(id))
	return bucket.Put(key, buf.Bytes())
}

// Close releases all database resources.
func (bs *BoltService) Close() error {
	return bs.db.Close()
}
