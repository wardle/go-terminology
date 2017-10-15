package mcqs

import (
	"bitbucket.org/wardle/go-snomed/snomed"
)

type FakeTruth struct {
	Diagnosis *snomed.Concept
	Parents   []*snomed.Concept // convenience pointers to parents
	Problems  []*snomed.Concept // problems for this diagnosis
}

// GenerateFakeTruth takes the SNOMED-CT ontology and uses it to build a fake "truth" model
// that represents the clinical findings that are seen in each type of diagnosis.
//
// While simply generating random problems for each diagnosis might be one approach, it is incorrect as
// we have a clear subsumption IS-A hierarchy which can be used. As such, related diagnostic concepts
// should share similar clinical problems in order to generate reasonable fake data.
func GenerateFakeTruth(dataset SnomedDataset) {
}
