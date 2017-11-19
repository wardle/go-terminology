package snomed

import (
	"sync"
)

// NaiveCache is an fairly naive in-memory cache used for development.
// It is designed to store arbitrary SNOMED-CT entities such as Concepts, Descriptions and Relationships.
// To make it easier to use, there are convenience methods to do the type-casting to put and get objects of
// these types which will simply wrap a better future implementation when implemented.
// This now is backed by a concurrent map instead ready for experimenting with concurrency.
// TODO: use an even better cache behind-the-scenes perhaps with persistence to filesystem?
type NaiveCache struct {
	cache sync.Map
}

// NewCache creates a new cache
func NewCache() *NaiveCache {
	return &NaiveCache{}
}

// Clear clears the cache
func (nc *NaiveCache) Clear() {
	nc.cache = *new(sync.Map)
}

// Get fetches a generic object from the cache using the specified identifier
// returning the object and a boolean indicating success (or not)
func (nc *NaiveCache) Get(id int) (interface{}, bool) {
	return nc.cache.Load(id)
}

// GetOrElse fethes a generic object from the cache or performs the callback specified, caching the result
func (nc *NaiveCache) GetOrElse(id int, f func(int) (interface{}, error)) (interface{}, error) {
	value, success := nc.Get(id)
	if success {
		return value, nil
	}
	value, err := f(id)
	if err != nil {
		return nil, err
	}
	nc.Put(id, value)
	return value, nil
}

// Put stores a generic object into the cache using the specified identifier
func (nc *NaiveCache) Put(id int, value interface{}) {
	nc.cache.Store(id, value)
}

// PutConcept stores a concept in the cache
func (nc *NaiveCache) PutConcept(conceptID int, concept *Concept) {
	nc.Put(conceptID, concept)
}

// GetConcept fetches a concept from the cache
func (nc *NaiveCache) GetConcept(conceptID int) (*Concept, bool) {
	value, success := nc.Get(conceptID)
	if !success {
		return nil, false
	}
	concept, success := value.(*Concept)
	return concept, success
}

// GetConceptOrElse fetches a concept from the cache or performs the callback specified, caching the result
func (nc *NaiveCache) GetConceptOrElse(conceptID int, f func(conceptID int) (interface{}, error)) (*Concept, error) {
	v, err := nc.GetOrElse(conceptID, f)
	if err != nil {
		return nil, err
	}
	return v.(*Concept), nil
}

// PutDescription stores a description in the cache
func (nc *NaiveCache) PutDescription(descriptionID int, description *Description) {
	nc.Put(descriptionID, description)
}

// GetDescription fetches a description from the cache
func (nc *NaiveCache) GetDescription(descriptionID int) (*Description, bool) {
	value, success := nc.Get(descriptionID)
	if !success {
		return nil, false
	}
	description, success := value.(*Description)
	return description, success
}

// GetDescriptionOrElse fetches a description from the cache or performs the callback specified, caching the result
func (nc *NaiveCache) GetDescriptionOrElse(descriptionID int, f func(descriptionID int) (interface{}, error)) (*Description, error) {
	v, err := nc.GetOrElse(descriptionID, f)
	if err != nil {
		return nil, err
	}
	return v.(*Description), nil
}

// PutRelationship stores a relationship in the cache
func (nc *NaiveCache) PutRelationship(descriptionID int, description *Description) {
	nc.Put(descriptionID, description)
}

// GetRelationship fetches a relationship from the cache
func (nc *NaiveCache) GetRelationship(relationshipID int) (*Relationship, bool) {
	value, success := nc.Get(relationshipID)
	if !success {
		return nil, false
	}
	relationship, success := value.(*Relationship)
	return relationship, true
}

// GetRelationshipOrElse fetches a relationship from the cache or performs the callback specified, caching the result
func (nc *NaiveCache) GetRelationshipOrElse(relationshipID int, f func(relationshipID int) (interface{}, error)) (*Relationship, error) {
	v, err := nc.GetOrElse(relationshipID, f)
	if err != nil {
		return nil, err
	}
	return v.(*Relationship), nil
}
