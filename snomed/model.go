// Package snomed is a simple proof-of-concept for a new SNOMED-CT microservice
package snomed

import (
	"fmt"
	"golang.org/x/text/language"
	"strconv"
)

// Important fixed SNOMED-CT concepts
var (
	SctDisease                     Identifier = 64572001
	SctCentralNervousSystemDisease Identifier = 23853001
	SctDiseases                    Identifier = 64572001 // parent concept of all diseases and syndromes within SNOMED-CT

)

// Types of relationship
// TODO: this is not exhaustive (yet)
// See children of 410662002: https://termbrowser.nhs.uk/?perspective=full&conceptId1=410662002
var (
	IsA                            Identifier = 116680003
	FindingSite                    Identifier = 363698007
	HasActiveIngredient            Identifier = 127489000
	HasAMP                         Identifier = 10362701000001108
	HasARP                         Identifier = 12223201000001101
	HasBasisOfStrength             Identifier = 10363001000001101
	HasDispensedDoseForm           Identifier = 10362901000001105
	HasDoseForm                    Identifier = 411116001
	HasExcipient                   Identifier = 8653101000001104
	HasSpecificActiveIngredient    Identifier = 10362801000001104
	HasTradeFamilyGroup            Identifier = 9191701000001107
	HasVMP                         Identifier = 10362601000001103
	HasVmpNonAvailabilityIndicator Identifier = 8940601000001102
	HasVmpPrescribingStatus        Identifier = 8940001000001105
	VrpPrescribingStatus           Identifier = 12223501000001103
)

// Concept is a SNOMED-CT concept
type Concept struct {
	ConceptID          Identifier
	FullySpecifiedName string
	Status             *Status
	Parents            []int // cache all recursive parents for optimised IS-A testing
}

// Description is a synonym for a concept.
type Description struct {
	DescriptionID        Identifier
	Status               DescriptionStatus
	Type                 DescriptionType
	InitialCapitalStatus int // TODO: make useful / meaningful
	LanguageCode         language.Tag
	Term                 string
}

// DescriptionStatus maps description status codes to meaning
type DescriptionStatus int

// Description status codes
const (
	DescriptionCurrent           DescriptionStatus = 0
	DescriptionNonCurrent        DescriptionStatus = 1
	DescriptionDuplicate         DescriptionStatus = 2
	DescriptionOutdated          DescriptionStatus = 3
	DescriptionErroneous         DescriptionStatus = 4
	DescriptionLimited           DescriptionStatus = 5
	DescriptionInappropriate     DescriptionStatus = 6
	DescriptionConceptNonCurrent DescriptionStatus = 7
	DescriptionMovedElsewhere    DescriptionStatus = 8
	DescriptionPendingMove       DescriptionStatus = 9
)

// IsCurrent returns whether the given status code is current
func (ds DescriptionStatus) IsCurrent() bool {
	return ds == DescriptionCurrent
}

// DescriptionType maps description types codes to meaning
type DescriptionType int

// Description type codes
const (
	Unspecified        DescriptionType = 0
	Preferred          DescriptionType = 1
	Synonym            DescriptionType = 2
	FullySpecifiedName DescriptionType = 3
)

// IsPreferred returns whether the given type code is preferred
func (dt DescriptionType) IsPreferred() bool {
	return dt == Preferred
}

// Relationship provides a relationship between two concepts of a type defined by a concept.
type Relationship struct {
	RelationshipID Identifier
	Source         Identifier
	Type           Identifier
	Target         Identifier
}

// IsA determines if this concept a type of that concept?
func (c *Concept) IsA(conceptID Identifier) bool {
	if c.ConceptID == conceptID {
		return true
	}
	id := int(conceptID)
	for _, a := range c.Parents {
		if a == id {
			return true
		}
	}
	return false
}

func (c *Concept) String() string {
	return c.FullySpecifiedName + " (" + strconv.Itoa(int(c.ConceptID)) + ")"
}

// ConceptStatus essentially records whether this concept is active or not
type ConceptStatus int

// Valid status codes
const (
	Current        ConceptStatus = 0
	Retired                      = 1
	Duplicate                    = 2
	Outdated                     = 3
	Ambiguous                    = 4
	Erroneous                    = 5
	Limited                      = 6
	MovedElsewhere               = 7
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
func NewConcept(conceptID Identifier, fullySpecifiedName string, statusID int, parents []int) (*Concept, error) {
	status := lookupStatus(statusID)
	if status != nil {
		return &Concept{conceptID, fullySpecifiedName, status, parents}, nil
	}
	return nil, fmt.Errorf("invalid status code: %d", statusID)
}

// NewRelationship creates a relationship between concepts
func NewRelationship(relationshipID Identifier, sourceConceptID Identifier, typeConceptID Identifier, targetConceptID Identifier) *Relationship {
	return &Relationship{relationshipID, sourceConceptID, typeConceptID, targetConceptID}
}

// Predicate defines an arbitrary way of selecting a concept based on some logic not defined here.
// As there are three possible results, not applicable, true and false, the result is a pointer to
// a boolean and so can be nil, true or false.
type Predicate interface {
	Test(concept *Concept) *bool
}

// Qualifier defines a chainable set of filters to include or exclude concepts.
// TODO: modify to support arbitrary boolean logic in a nestable qualifier pattern
type Qualifier struct {
	predicate Predicate
	next      *Qualifier
}

// rule is a simple way of building an internal predicate based on IS-A relationships
type rule struct {
	identifier Identifier
	include    bool
}

func (r rule) Test(concept *Concept) *bool {
	if concept.IsA(r.identifier) {
		return &r.include
	}
	return nil
}

// NewQualifier builds a simple include/exclude chainable filter for concepts
func NewQualifier() *Qualifier {
	return &Qualifier{nil, nil}
}

func (q *Qualifier) addRule(identifier Identifier, include bool) *Qualifier {
	rule := &rule{identifier, include}
	if q.predicate == nil {
		q.predicate = rule
	} else {
		next := &Qualifier{rule, nil}
		q.next = next
	}
	return q
}

// Include adds an additional include filter to an existing chain
func (q *Qualifier) Include(identifier Identifier) *Qualifier {
	return q.addRule(identifier, true)
}

// Exclude adds an additional exclude filter to an existing chain
func (q *Qualifier) Exclude(identifier Identifier) *Qualifier {
	return q.addRule(identifier, true)
}

// Apply applies the chain of predicates to the given concept
func (q *Qualifier) test(concept *Concept, shouldInclude bool) bool {
	if q.predicate == nil {
		return shouldInclude
	}
	r := q.predicate.Test(concept)
	if r != nil {
		shouldInclude = *r
	}
	if n := q.next; n != nil {
		shouldInclude = n.test(concept, shouldInclude)
	}
	return shouldInclude
}

// Test is a simple predicate that determines whether the given concept should be included
func (q Qualifier) Test(concept *Concept) bool {
	return q.test(concept, true)
}
