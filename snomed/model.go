// Package snomed defines the specification for SNOMED-CT releases in the RF2 format.
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
//go:generate protoc -I../vendor/terminology/protos --go_out=plugins=gprc:. ../vendor/terminology/protos/snomed.proto
//go:generate protoc -I../vendor/terminology/protos -I../vendor/terminology/vendor/googleapis --go_out=plugins=grpc:. ../vendor/terminology/protos/server.proto
//go:generate protoc -I../vendor/terminology/protos -I../vendor/terminology/vendor/googleapis --grpc-gateway_out=logtostderr=true:. ../vendor/terminology/protos/server.proto
//
// ***************************************************************************
//    Copyright 2018 Mark Wardle / Eldrix Ltd
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
package snomed

import (
	"golang.org/x/text/language"
	"unicode"
)

// DefinitionStatusID is an indication as to whether a concept is fully defined or primitive.
// A concept may be primitive if
// a) it is inherently impossible to define the concept by defining relationships with other concepts
// b) because  attributes that would distinguish the concept from another concept are not (yet) in the concept model
// c) because the concept is not in a directly clinical domain but in one of the supporting domains used to define clinical concepts
// d) because modeling of defining relationships is a continuing process and in same cases is incomplete
type DefinitionStatusID int64

// Available definition statuses (at the time of writing)
// The RF2 uses SNOMED concepts to populate enums which means that conceivably, the enum
// could be extended in future releases. This rather complicates writing any inferential code.
const (
	Primitive DefinitionStatusID = 900000000000074008 // defines a primitive concept does not have sufficient defining relationships to computably distinguish them from more general concepts(supertypes.
	Defined   DefinitionStatusID = 900000000000073002 // defines a concept with a formal logic definition that is sufficient to distinguish its meaning from other similar concepts.
)

// IsSufficientlyDefined returns whether this concept has a formal logic definition that is sufficient to distinguish
// its meaning from other similar concepts.
func (c *Concept) IsSufficientlyDefined() bool {
	return c.DefinitionStatusId == int64(Defined)
}

// DescriptionTypeID gives the type this description represents for the concept.
type DescriptionTypeID int64

// Available description TypeIDs.
// The synonym defines the preferred term for the concept in the language specified.
// The Description .term represents a term that is used to represent the associated concept in the language indicated by the Description .languageCode.
// Note: The preferred term used in a given language or dialect is marked as a synonym. Preference and acceptability of a particular synonymous term is indicated by a Language refset.
const (
	Definition         DescriptionTypeID = 900000000000550004 // A term representing the associated concept in the language indicated by Description .languageCode.
	FullySpecifiedName DescriptionTypeID = 900000000000003001 // A term unique among active descriptions in SNOMED CT that names the meaning of a concept code in a manner that is intended to be unambiguous and stable across multiple contexts.
	Synonym            DescriptionTypeID = 900000000000013009 // A term that is an acceptable way to express a the meaning of a SNOMED CT concept in a particular language.
)

// LanguageTag returns the language tag for this description
func (d *Description) LanguageTag() language.Tag {
	return language.Make(d.LanguageCode)
}

// IsFullySpecifiedName returns whether this is a fully specified name
func (d *Description) IsFullySpecifiedName() bool {
	return d.TypeId == int64(FullySpecifiedName)
}

// IsSynonym returns whether this is a preferred term for a language
func (d *Description) IsSynonym() bool {
	return d.TypeId == int64(Synonym)
}

// IsDefinition returns whether this description is simply one of (many) alternative descriptions
func (d *Description) IsDefinition() bool {
	return d.TypeId == int64(Definition)
}

// CaseSignificanceID provides information about case significance for the description
type CaseSignificanceID int64

// Available case significance options
const (
	EntireTermCaseInsensitive     CaseSignificanceID = 900000000000448009
	EntireTermCaseSensitive       CaseSignificanceID = 900000000000017005
	InitialCharacterCaseSensitive CaseSignificanceID = 900000000000020002
)

// Uncapitalized returns the term appropriately uncapitalized
// All terms are, by default, capitalized but we cannot naively assume
// it is possible to uncapitalize without checking case sensitivity.
func (d *Description) Uncapitalized() string {
	sig := d.CaseSignificance
	if sig == int64(EntireTermCaseSensitive) || sig == int64(InitialCharacterCaseSensitive) {
		return d.Term
	}
	for i, v := range d.Term {
		return string(unicode.ToLower(v)) + d.Term[i+1:]
	}
	return ""
}

// Valid types of characteristic types
const (
	additionalRelationship int64 = 900000000000227009
	definingRelationship   int64 = 900000000000006009 // NB:  has children inferred and stated
	inferredRelationship   int64 = 900000000000011006 // NB: IS-A defining
	statedRelationship     int64 = 900000000000010007 // NB: IS-A defining
	qualifyingRelationship int64 = 900000000000225001
)

// IsAdditionalRelationship specifies whether this is a relationship to a target concept that is additional to the core
func (r *Relationship) IsAdditionalRelationship() bool {
	return r.CharacteristicTypeId == additionalRelationship
}

// IsDefiningRelationship returns whether this is a relationship to a target concept that is always necessarily true from any instance of the source concept.
func (r *Relationship) IsDefiningRelationship() bool {
	t := r.CharacteristicTypeId
	return t == definingRelationship || t == inferredRelationship || t == statedRelationship
}

// IsQualifyingRelationship An attribute-value relationship associated with a concept code to indicate to users that it may be applied to refine the meaning of the code.
// The set of qualifying relationships provide syntactically correct values that can be presented to a user for postcoordination.
// Following the introduction of the RF2 in 2012 qualifying relationships are no longer part of the standard distributed release.
// The Machine Readable Concept Model provides a more comprehensive and flexible way to identify the full set of attributes and ranges that can be applied to refine concepts in particular  domains.
func (r *Relationship) IsQualifyingRelationship() bool {
	return r.CharacteristicTypeId == qualifyingRelationship
}

// Types of Reference Set
const (
	rootRefset             int64 = 900000000000455006 // root concept for all reference set types
	refSetDescriptorRefset int64 = 900000000000456007 // represented by RefsetDescriptorReferenceSet
	simpleRefset           int64 = 446609009          // represented by SimpleReferenceSet
	languageRefset         int64 = 900000000000506000 // represented by LanguageReferenceSet
	simpleMapRefset        int64 = 900000000000496009 // represented by SimpleMapReferenceSet
	complexMapRefset       int64 = 447250001          // represented by ComplexMapReferenceSet
	extendedMapRefset      int64 = 609331003          // represented by ComplexMapReferenceSet
)

// Valid types of acceptability. If a term is not either acceptable or preferred, it is unacceptable in this language.
const (
	acceptable int64 = 900000000000549004
	preferred  int64 = 900000000000548007
)

// IsAcceptable returns whether the description referenced is acceptable for this concept in this language refset
func (lrs *LanguageReferenceSet) IsAcceptable() bool {
	if lrs != nil {
		return lrs.AcceptabilityId == acceptable
	}
	return false
}

// IsPreferred returns whether the description referenced is the preferred for this concept in this language refset
func (lrs *LanguageReferenceSet) IsPreferred() bool {
	if lrs != nil {
		return lrs.AcceptabilityId == preferred
	}
	return false
}

// IsUnacceptable returns whether the description referenced is unacceptable for this concept in this language refset
func (lrs *LanguageReferenceSet) IsUnacceptable() bool {
	if lrs != nil {
		return lrs.AcceptabilityId != preferred && lrs.AcceptabilityId != acceptable
	}
	return true
}
