// Package snomed is a simple proof-of-concept for a new SNOMED-CT microservice
package snomed

import (
	"errors"
)

// Concept is a SNOMED-CT concept
type Concept struct {
	ConceptID          int
	FullySpecifiedName string
	Status             *Status
	Parents            []int
}

// IsA determines if this concept a type of that concept?
func (c Concept) IsA(conceptID int) bool {
	for _, a := range c.Parents {
		if a == conceptID {
			return true
		}
	}
	return false
}

// ConceptStatus essentially records whether this concept is active or not
type ConceptStatus int

// Valid status codes
const (
	Current        ConceptStatus = iota //=0
	Retired                             //=1
	Duplicate                           //=2
	Outdated                            //=3
	Ambiguous                           //=4
	Erroneous                           //=5
	Limited                             //=6
	MovedElsewhere                      //=7
)

// IsActive returns whether this ConceptStatus should be regarded as "active"
func (s ConceptStatus) IsActive() bool {
	return s == Current
}

// Status of a SNOMED CT concept
type Status struct {
	Code     int
	Title    string
	IsActive bool
}

// map concept status code to its meaning
var statuses = make(map[int]*Status)

// initialiser
func init() {
	statuses[0] = &Status{0, "Current", true}            //: current (considered active)
	statuses[1] = &Status{1, "Retired", false}           //: Retired (considered inactive)
	statuses[2] = &Status{2, "Duplicate", false}         //: Duplicate (considered inactive)
	statuses[3] = &Status{3, "Outdated", false}          //: Outdated (considered inactive)
	statuses[4] = &Status{4, "Ambiguous", false}         //: Ambiguous (considered inactive)
	statuses[5] = &Status{5, "Erroneous", false}         //: Erroneous (considered inactive)
	statuses[6] = &Status{6, "Limited", true}            //: Limited (considered active)
	statuses[10] = &Status{10, "Moved elsewhere", false} //: Moved elsewhere (considered inactive)
	statuses[11] = &Status{11, "Pending move", true}     //: Pending move (considered active)
}

// return the status for the specified code
func lookupStatus(code int) *Status {
	return statuses[code]
}

// NewConcept creates a concept
func NewConcept(conceptID int, fullySpecifiedName string, statusID int, parents []int) (*Concept, error) {
	status := lookupStatus(statusID)
	if status != nil {
		return &Concept{conceptID, fullySpecifiedName, status, parents}, nil
	}
	return nil, errors.New("Invalid status code")
}

/*
// future proper modelling....

// Concept is an opaque SNOMED-CT concept
type Concept interface {
	conceptID() int             // return the concept identifier
	fullySpecifiedName() string // return the fully specified name
	status() *Status            // return the status
	//	parents() []Concept                                // return the IS-A parents for this concept
	//	children() []Concept                               // return the IS-A children for this concept
	//	preferredDescription(tag language.Tag) Description // return the preferred description for the specified locale
	//	isA(concept Concept) bool                          // is this concept a type of that concept
	//	childRelationships() []Relationship                // return the relationships in which this concept is the target
	//	parentRelationships() []Relationship               // return the relationships in which this concept is the source
}

// Description is a human-readable synonym for a given concept
type Description interface {
	descriptionId() int        // return the description identifier
	concept() *Concept         // return the concept that this represents
	languageTag() language.Tag // return the language of this synonym
}

// Relationship defines a relationship between one concept and another.
// The commonest is a IS-A relationship
type Relationship interface {
	source() *Concept       // source concept
	relationship() *Concept // type of the relationship
	target() *Concept       // target concept
}

// ConceptService represents an opaque way to fetch and navigate the SNOMED-CT model
type ConceptService interface {
	FetchById(conceptId int) (*Concept, error)
}
type DescriptionService interface {
	FetchById(descriptionId int) (*Description, error)
}

*/
