// Package rf2 defines the specification for SNOMED-CT releases in the RF2 format.
// See https://confluence.ihtsdotools.org/display/DOCRELFMT/SNOMED+CT+Release+File+Specifications
//
// These are, in large part, raw representations of the release files with some small
// additions, predominantly relating to valid enumerations, to aid computability.
//
// The structures are not intended for use in a working terminology server which needs
// optimised structures in order to be performant.
//
// These structures are designed to cope with importing any SNOMED-CT distribution, including
// full distributions, a snapshot or a delta.
//
// Full	The files representing each type of component contain every version of every component ever released.
// Snapshot	The files representing each type of component contain one version of every component released up to the time of the snapshot. The version of each component contained in a snapshot is the most recent version of that component at the time of the snapshot.
// Delta	The files representing each type of component contain only component versions created since the previous release. Each component version in a delta release represents either a new component or a change to an existing component.
//
// NB: The associated import functionality imports only from snapshot files.
//
// // Copyright 2018 Mark Wardle / Eldrix Ltd
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
package rf2

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"golang.org/x/text/language"
	"time"
	"unicode"
)

// A Concept represents a SNOMED-CT concept.
// The RF2 release allows multiple duplicate entries per concept identifier to permit versioning.
// As such, we have a compound primary key made up of the concept identifier and the effective time.
// Only one concept with a specified identifier will be active at any time point.
// See https://confluence.ihtsdotools.org/display/DOCRELFMT/3.2.1.+Concept+File+Specification
type Concept struct {
	ID                 snomed.Identifier // Uniquely identifies the concept.
	EffectiveTime      time.Time         // Specifies the inclusive date at which the component version's state became the then current valid state of the component
	Active             bool              // Specifies whether the concept was active or inactive from the nominal release date specified by the effectiveTime.
	ModuleID           snomed.Identifier // Identifies the concept version's module. Set to a descendant of 900000000000443000 |Module|within the metadata hierarchy.
	DefinitionStatusID snomed.Identifier // Specifies if the concept version is primitive or sufficiently defined. Set to a descendant of 900000000000444006 |Definition status|in the metadata hierarchy.
}

// Available definition statuses (at the time of writing)
// The RF2 uses SNOMED concepts to populate enums which means that conceivably, the enum
// could be extended in future releases. This rather complicates writing any inferential code.
const (
	primitive snomed.Identifier = 900000000000074008 // defines a primitive concept does not have sufficient defining relationships to computably distinguish them from more general concepts(supertypes.
	defined   snomed.Identifier = 900000000000073002 // defines a concept with a formal logic definition that is sufficient to distinguish its meaning from other similar concepts.
)

// IsSufficientlyDefined returns whether this concept has a formal logic definition that is sufficient to distinguish
// its meaning from other similar concepts.
func (c *Concept) IsSufficientlyDefined() bool {
	return c.DefinitionStatusID == defined
}

// A Description holds descriptions that describe SNOMED CT concepts.
// A description is used to give meaning to a concept and provide well-understood and standard ways of referring to a concept.
// See https://confluence.ihtsdotools.org/display/DOCRELFMT/3.2.2.+Description+File+Specification
type Description struct {
	ID               snomed.Identifier // Uniquely identifies the description.
	EffectiveTime    time.Time         // Specifies the inclusive date at which the component version's state became the then current valid state of the component
	Active           bool              // Specifies whether the state of the description was active or inactive from the nominal release date specified by the effectiveTime .
	ModuleID         snomed.Identifier // Identifies the description version's module. Set to a child of 900000000000443000 |Module|within the metadata hierarchy.
	ConceptID        snomed.Identifier // Identifies the concept to which this description applies. Set to the identifier of a concept in the 138875005 |SNOMED CT Concept| hierarchy within the Concept.
	LanguageCode     string            // Specifies the language of the description text using the two character ISO-639-1 code. Note that this specifies a language level only, not a dialect or country code.
	TypeID           snomed.Identifier // Identifies whether the description is fully specified name a synonym or other description type. This field is set to a child of 900000000000446008 |Description type|in the Metadata hierarchy.
	Term             string            // The description version's text value, represented in UTF-8 encoding.
	CaseSignificance snomed.Identifier // Identifies the concept enumeration value that represents the case significance of this description version. For example, the term may be completely case sensitive, case insensitive or initial letter case insensitive. This field will be set to a child of 900000000000447004 |Case significance|within the metadata hierarchy.
}

// Available description TypeIDs.
// The synonym defines the preferred term for the concept in the language specified.
// The Description .term represents a term that is used to represent the associated concept in the language indicated by the Description .languageCode.
// Note: The preferred term used in a given language or dialect is marked as a synonym. Preference and acceptability of a particular synonymous term is indicated by a Language refset.
const (
	definition         snomed.Identifier = 900000000000550004 // A term representing the associated concept in the language indicated by Description .languageCode.
	fullySpecifiedName snomed.Identifier = 900000000000003001 // A term unique among active descriptions in SNOMED CT that names the meaning of a concept code in a manner that is intended to be unambiguous and stable across multiple contexts.
	synonym            snomed.Identifier = 900000000000013009 // A term that is an acceptable way to express a the meaning of a SNOMED CT concept in a particular language.
)

// LanguageTag returns the language tag for this description
func (d *Description) LanguageTag() language.Tag {
	return language.Make(d.LanguageCode)
}

// IsFullySpecifiedName returns whether this is a fully specified name
func (d *Description) IsFullySpecifiedName() bool {
	return d.TypeID == fullySpecifiedName
}

// IsSynonym returns whether this is a preferred term for a language
func (d *Description) IsSynonym() bool {
	return d.TypeID == synonym
}

// IsDefinition returns whether this description is simply one of (many) alternative descriptions
func (d *Description) IsDefinition() bool {
	return d.TypeID == definition
}

// Available case significance options
const (
	entireTermCaseInsensitive     snomed.Identifier = 900000000000448009
	entireTermCaseSensitive       snomed.Identifier = 900000000000017005
	initialCharacterCaseSensitive snomed.Identifier = 900000000000020002
)

// Uncapitalized returns the term appropriately uncapitalized
// All terms are, by default, capitalized but we cannot naively assume
// it is possible to uncapitalize without checking case sensitivity.
func (d *Description) Uncapitalized() string {
	sig := d.CaseSignificance
	if sig == entireTermCaseSensitive || sig == initialCharacterCaseSensitive {
		return d.Term
	}
	for i, v := range d.Term {
		return string(unicode.ToLower(v)) + d.Term[i+1:]
	}
	return ""
}

// Relationship defines a relationship between two concepts as a type itself defined as a concept
type Relationship struct {
	ID                   snomed.Identifier // Uniquely identifies the relationship.
	EffectiveTime        time.Time         // Specifies the inclusive date at which the component version's state became the then current valid state of the component
	Active               bool              // Specifies whether the state of the relationship was active or inactive from the nominal release date specified by the effectiveTime field.
	ModuleID             snomed.Identifier // Identifies the relationship version's module. Set to a child of 900000000000443000 |Module|within the metadata hierarchy.
	SourceID             snomed.Identifier // Identifies the source concept of the relationship version. That is the concept defined by this relationship. Set to the identifier of a concept. in the Concept File.
	DestinationID        snomed.Identifier // Identifies the concept that is the destination of the relationship version.
	RelationshipGroup    int               // Groups together relationship versions that are part of a logically associated relationshipGroup. All active Relationship records with the same relationshipGroup number and sourceId are grouped in this way.
	TypeID               snomed.Identifier // Identifies the concept that represent the defining attribute (or relationship type) represented by this relationship version.
	CharacteristicTypeID snomed.Identifier // A concept enumeration value that identifies the characteristic type of the relationship version (i.e. whether the relationship version is defining, qualifying, etc.) This field is set to a descendant of 900000000000449001 |Characteristic type|in the metadata hierarchy.
	ModifierID           snomed.Identifier // Ignore. A concept enumeration value that identifies the type of Description Logic (DL) restriction (some, all, etc.).
}

// Valid types of characteristic types
const (
	additionalRelationship snomed.Identifier = 900000000000227009
	definingRelationship   snomed.Identifier = 900000000000006009 // NB:  has children inferred and stated
	inferredRelationship   snomed.Identifier = 900000000000011006 // NB: IS-A defining
	statedRelationship     snomed.Identifier = 900000000000010007 // NB: IS-A defining
	qualifyingRelationship snomed.Identifier = 900000000000225001
)

// IsAdditionalRelationship specifies whether this is a relationship to a target concept that is additional to the core
func (r *Relationship) IsAdditionalRelationship() bool {
	return r.CharacteristicTypeID == additionalRelationship
}

// IsDefiningRelationship returns whether this is a relationship to a target concept that is always necessarily true from any instance of the source concept.
func (r *Relationship) IsDefiningRelationship() bool {
	t := r.CharacteristicTypeID
	return t == definingRelationship || t == inferredRelationship || t == statedRelationship
}

// IsQualifyingRelationship An attribute-value relationship associated with a concept code to indicate to users that it may be applied to refine the meaning of the code.
// The set of qualifying relationships provide syntactically correct values that can be presented to a user for postcoordination.
// Following the introduction of the RF2 in 2012 qualifying relationships are no longer part of the standard distributed release.
// The Machine Readable Concept Model provides a more comprehensive and flexible way to identify the full set of attributes and ranges that can be applied to refine concepts in particular  domains.
func (r *Relationship) IsQualifyingRelationship() bool {
	return r.CharacteristicTypeID == qualifyingRelationship
}

// Types of Reference Set
const (
	rootRefset             snomed.Identifier = 900000000000455006 // root concept for all reference set types
	refSetDescriptorRefset snomed.Identifier = 900000000000456007 // represented by RefsetDescriptorReferenceSet
	simpleRefset           snomed.Identifier = 446609009          // represented by SimpleReferenceSet
	languageRefset         snomed.Identifier = 900000000000506000 // represented by LanguageReferenceSet
	simpleMapRefset        snomed.Identifier = 900000000000496009 // represented by SimpleMapReferenceSet
	complexMapRefset       snomed.Identifier = 447250001          // represented by ComplexMapReferenceSet
	extendedMapRefset      snomed.Identifier = 609331003          // represented by ComplexMapReferenceSet
)

// ReferenceSet support customization and enhancement of SNOMED CT content. These include representation of subsets,
// language preferences maps for or from other code systems.
// There are multiple reference set types which extend this structure
// In the specification, the referenced component ID can be a SCT identifier or a UUID which is... problematic.
// Fortunately, in concrete types of reference set ("patterns"), it is made explicit.
type ReferenceSet struct {
	ID            string            // A 128 bit unsigned Integer, uniquely identifying the reference set member.
	EffectiveTime time.Time         // Specifies the inclusive date at which this change becomes effective.
	Active        bool              // Specifies whether the member's state was active or inactive from the nominal release date specified by the effectiveTime field.
	ModuleID      snomed.Identifier // Identifies the member version's module. Set to a child of 900000000000443000 |Module| within the metadata hierarchy .
	RefsetID      snomed.Identifier // Uniquely identifies the reference set that this extension row is part of. Set to a descendant of 900000000000455006 |Reference set| within the metadata hierarchy .
}

// RefSetDescriptorReferenceSet is a type of reference set that provides information about a different reference set
// See https://confluence.ihtsdotools.org/display/DOCRELFMT/4.2.11.+Reference+Set+Descriptor
// It provides the additional structure for a given reference set.
type RefSetDescriptorReferenceSet struct {
	ReferenceSet
	ReferencedComponentID  snomed.Identifier
	AttributeDescriptionID snomed.Identifier // Specifies the name of an attribute that is used in the reference set to which this descriptor applies.
	AttributeTypeID        snomed.Identifier // Specifies the data type of this attribute in the reference set to which this descriptor applies.
	AttributeOrder         uint              // An unsigned Integer, providing an ordering for the additional attributes extending the reference set .
}

// SimpleReferenceSet is a simple reference set usable for defining subsets
// See https://confluence.ihtsdotools.org/display/DOCRELFMT/4.2.1.+Simple+Reference+Set
type SimpleReferenceSet struct {
	ReferenceSet
	ReferencedComponentID snomed.Identifier
}

// LanguageReferenceSet is a A 900000000000506000 |Language type reference set| supporting the representation of
// language and dialects preferences for the use of particular descriptions.
// "The most common use case for this type of reference set is to specify the acceptable and preferred terms
// for use within a particular country or region. However, the same type of reference set can also be used to
// represent preferences for use of descriptions in a more specific context such as a clinical specialty,
// organization or department.
//
// No more than one description of a specific description type associated with a single concept may have the acceptabilityId value 900000000000548007 |Preferred|.
// Every active concept should have one preferred synonym in each language.
// This means that a language reference set should assign the acceptabilityId  900000000000548007 |Preferred|  to one  synonym (a  description with  typeId value 900000000000013009 |synonym|) associated with each concept .
// This description is the preferred term for that concept in the specified language or dialect.
// Any  description which is not referenced by an active row in the   reference set is regarded as unacceptable (i.e. not a valid  synonym in the language or  dialect ).
// If a description becomes unacceptable, the relevant language reference set member is inactivated by adding a new row with the same id, the effectiveTime of the the change and the value active=0.
// For this reason there is no requirement for an "unacceptable" value."
// See https://confluence.ihtsdotools.org/display/DOCRELFMT/4.2.4.+Language+Reference+Set
//
type LanguageReferenceSet struct {
	ReferenceSet
	ReferencedComponentID snomed.Identifier
	AcceptabilityID       snomed.Identifier // A subtype of 900000000000511003 |Acceptability| indicating whether the description is acceptable or preferred for use in the specified language or dialect .
}

// Valid types of acceptability. If a term is not either acceptable or preferred, it is unacceptable in this language.
const (
	acceptable snomed.Identifier = 900000000000549004
	preferred  snomed.Identifier = 900000000000548007
)

// IsAcceptable returns whether the description referenced is acceptable for this concept in this language refset
func (lrs *LanguageReferenceSet) IsAcceptable() bool {
	return lrs.AcceptabilityID == acceptable
}

// IsPreferred returns whether the description referenced is the preferred for this concept in this language refset
func (lrs *LanguageReferenceSet) IsPreferred() bool {
	return lrs.AcceptabilityID == preferred
}

// IsUnacceptable returns whether the description referenced is unacceptable for this concept in this language refset
func (lrs *LanguageReferenceSet) IsUnacceptable() bool {
	return lrs.AcceptabilityID != preferred && lrs.AcceptabilityID != acceptable
}

// SimpleMapReferenceSet is a straightforward one-to-one map between SNOMED-CT concepts and another
// coding system. This is appropriate for simple maps.
// See https://confluence.ihtsdotools.org/display/DOCRELFMT/4.2.9.+Simple+Map+Reference+Set
type SimpleMapReferenceSet struct {
	ReferenceSet
	ReferencedComponentID snomed.Identifier // A reference to the SNOMED CT component to be included in the reference set.
	MapTarget             string            // The equivalent code in the other terminology, classification or code system.
}

// ComplexMapReferenceSet represents a complex one-to-many map between SNOMED-CT and another
// coding system.
// A 447250001 |Complex map type reference set|enables representation of maps where each SNOMED
// CT concept may map to one or more codes in a target scheme.
// The type of reference set supports the general set of mapping data required to enable a
// target code to be selected at run-time from a number of alternate codes. It supports
// target code selection by accommodating the inclusion of machine readable rules and/or human readable advice.
// An 609331003 |Extended map type reference set|adds an additional field to allow categorization of maps.
type ComplexMapReferenceSet struct {
	ReferenceSet
	ReferencedComponentID snomed.Identifier
	MapGroup              int               // An Integer, grouping a set of complex map records from which one may be selected as a target code.
	MapPriority           int               // Within a mapGroup, the mapPriority specifies the order in which complex map records should be checked
	MapRule               string            // A machine-readable rule, (evaluating to either 'true' or 'false' at run-time) that indicates whether this map record should be selected within its mapGroup.
	MapAdvice             string            // Human-readable advice, that may be employed by the software vendor to give an end-user advice on selection of the appropriate target code from the alternatives presented to him within the group.
	MapTarget             string            // The target code in the target terminology, classification or code system.
	CorrelationID         snomed.Identifier // A child of 447247004 |SNOMED CT source code to target map code correlation value|in the metadata hierarchy, identifying the correlation between the SNOMED CT concept and the target code.
	MapCategoryID         snomed.Identifier // Only for extended complex map refsets: Identifies the SNOMED CT concept in the metadata hierarchy which represents the MapCategory for the associated map member.
}
