package snomed

import (
	"fmt"
)

// NaiveCache is an extraordinarily naive in-memory cache used for development.
// It is designed to store arbitrary SNOMED-CT entities such as Concepts, Descriptions and Relationships.
// To make it easier to use, there are convenience methods to do the type-casting to put and get objects of
// these types which will simply wrap a better future implementation when implemented.
// TODO: use a better cache behind-the-scenes.
type NaiveCache struct {
	cache map[int]interface{}
}

// NewCache creates a new cache
func NewCache() *NaiveCache {
	return &NaiveCache{make(map[int]interface{})}
}

// Get fetches a generic object from the cache using the specified identifier
// returning the object and a boolean indicating success (or not)
func (nc NaiveCache) Get(id int) (interface{}, bool) {
	value := nc.cache[id]
	if value != nil {
		return value, true
	}
	return nil, false
}

// Put stores a generic object into the cache using the specified identifier
func (nc NaiveCache) Put(id int, value interface{}) {
	nc.cache[id] = value
}

// PutConcept stores a concept in the cache
func (nc NaiveCache) PutConcept(conceptID int, concept *Concept) {
	nc.Put(conceptID, concept)
}

// GetConcept fetches a concept from the cache
func (nc NaiveCache) GetConcept(conceptID int) (*Concept, bool) {
	value, success := nc.Get(conceptID)
	if success {
		concept, ok := value.(*Concept)
		if ok {
			return concept, true
		}
	}
	return nil, false
}

// GetConceptOrElse fetches a concept from the cache or performs the callback specified, caching the result
func (nc NaiveCache) GetConceptOrElse(conceptID int, f func(conceptID int) (*Concept, error)) (*Concept, error) {
	concept, success := nc.GetConcept(conceptID)
	if success {
		return concept, nil
	}
	return f(conceptID)
}

// PutDescription stores a description in the cache
func (nc NaiveCache) PutDescription(descriptionID int, description *Description) {
	nc.Put(descriptionID, description)
}

// GetDescription fetches a description from the cache
func (nc NaiveCache) GetDescription(descriptionID int) (*Description, bool) {
	value, success := nc.Get(descriptionID)
	if success {
		description, ok := value.(*Description)
		if ok {
			return description, true
		}
	}
	return nil, false
}

// GetDescriptionOrElse fetches a description from the cache or performs the callback specified, caching the result
func (nc NaiveCache) GetDescriptionOrElse(descriptionID int, f func(descriptionID int) (*Description, error)) (*Description, error) {
	description, success := nc.GetDescription(descriptionID)
	if success {
		return description, nil
	}
	return f(descriptionID)
}

// PutRelationship stores a relationship in the cache
func (nc NaiveCache) PutRelationship(descriptionID int, description *Description) {
	nc.Put(descriptionID, description)
}

// GetRelationship fetches a relationship from the cache
func (nc NaiveCache) GetRelationship(relationshipID int) (*Relationship, bool) {
	value, success := nc.Get(relationshipID)
	if success {
		relationship, ok := value.(*Relationship)
		if ok {
			return relationship, true
		}
	}
	return nil, false
}

// GetRelationshipOrElse fetches a relationship from the cache or performs the callback specified, caching the result
func (nc NaiveCache) GetRelationshipOrElse(relationshipID int, f func(relationshipID int) (*Relationship, error)) (*Relationship, error) {
	relationship, success := nc.GetRelationship(relationshipID)
	if success {
		return relationship, nil
	}
	return f(relationshipID)
}

func (nc NaiveCache) String() string {
	return fmt.Sprintf("Cache with %d entries", len(nc.cache))
}
