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

package snomed

// Common concept identifiers, predominantly for the SNOMED-CT metadata model.
const (
	// The SNOMED CT root concept
	Root Identifier = 138875005

	// IsAConceptID represents the relationship type, IS-A; the commonest type of relationship
	IsA = 116680003

	// Top level concepts
	BodyStructure                   = 123037004
	ClinicalFinding                 = 404684003
	EnvironmentGeographicLocation   = 308916002
	Event                           = 308916002
	ObservableEntity                = 363787002
	Organism                        = 410607006
	PharmaceuticalBiologicalProduct = 373873005
	PhysicalForce                   = 78621006
	PhysicalObject                  = 260787004
	Procedure                       = 71388002
	QualifierValue                  = 362981000
	RecordArtefact                  = 419891008
	SituationWithExplicitContext    = 243796009
	SocialContext                   = 243796009
	Specimen                        = 123038009
	StagingAndScales                = 254291000
	Substance                       = 254291000
	LinkageConcept                  = 106237007

	// Special concepts
	SpecialConcept      = 370115009
	NavigationalConcept = 363743006

	// RelationshipType concepts
	Attribute                    = 246061005
	ConceptModelAttribute        = 410662002
	Access                       = 260507000
	AssociatedFinding            = 246090004
	AssociatedMorphology         = 116676008
	AssociatedProcedure          = 363589002
	AssociatedWith               = 47429007
	After                        = 255234002
	CausativeAgent               = 246075003
	DueTo                        = 42752001
	ClinicalCourse               = 263502005
	Component                    = 246093002
	DirectSubstance              = 363701004
	Episodicity                  = 246456000
	FindingContext               = 408729009
	FindingInformer              = 419066007
	FindingMethod                = 418775008
	FindingSite                  = 363698007
	HasActiveIngredient          = 127489000
	HasDefinitionalManifestation = 363705008
	HasDoseForm                  = 411116001
	HasFocus                     = 363702006
	HasIntent                    = 363703001
	HasInterpretation            = 363713009
	HasSpecimen                  = 116686009
	Interprets                   = 363714003
	Laterality                   = 272741003
	MeasurementMethod            = 370129005
	Method                       = 260686004
	Occurrence                   = 246454002
	PartOf                       = 123005000
	PathologicalProcess          = 370135005
	Priority                     = 260870009
	ProcedureContext             = 408730004
	ProcedureDevice              = 405815000
	DirectDevice                 = 363699004
	IndirectDevice               = 363710007
	UsingDevice                  = 424226004
	UsingAccessDevice            = 425391005
	ProcedureMorphology          = 405816004
	DirectMorphology             = 363700003
	IndirectMorphology           = 363709002
	ProcedureSite                = 363704007
	ProcedureSiteDirect          = 405813007
	ProcedureSiteIndirect        = 405814001
	Property                     = 370130000
	RecipientCategory            = 370131001
	RevisionStatus               = 246513007
	RouteOfAdministration        = 410675002
	ScaleType                    = 370132008
	Severity                     = 246112005
	SpecimenProcedure            = 118171006
	SpecimenSourceIdentity       = 118170007
	SpecimenSourceMorphology     = 118168003
	SpecimenSourceTopography     = 118169006
	SpecimenSubstance            = 370133003
	SubjectOfInformation         = 131195008
	SubjectRelationshipContext   = 408732007
	SurgicalApproach             = 424876005
	TemporalContext              = 408731000
	TimeAspect                   = 370134009
	UsingEnergy                  = 424244007
	UsingSubstance               = 42436100
)
