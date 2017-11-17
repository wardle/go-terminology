package mcqs

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"fmt"
	"strings"
)

// Question is made up of a scenario containing findings, a lead-in, five possible answers and one single best answer.
type Question struct {
	Age             int
	Sex             Sex
	Findings        []*ClinicalFinding // a list of clinical findings derived from the scenario (stem)
	LeadIn          LeadIn             // the question based on the stem
	PossibleAnswers []*snomed.Concept  // the (usually five) answers, all of which could be correct but only one is the best answer
	CorrectAnswer   *snomed.Concept    // the single best answer
}

func (q Question) String() string {
	findings := make([]string, 0)
	for _, finding := range q.Findings {
		findings = append(findings, finding.String())
	}
	return fmt.Sprintf("[%s] --> %s", strings.Join(findings, ", "), q.CorrectAnswer.FullySpecifiedName)

}

// Sex of patient in the question
type Sex int

// Possible values for Sex
const (
	Male   Sex = 1
	Female Sex = 2
)

// LeadIn is the question asked after the scenario.
type LeadIn int

// We currently support only a single type of lead-in, as the focus here is on diagnosis only
// Interestingly however, the answers to next tests, treatment or action usually depend on
// multi-step problem-solving from clinical scenario -> likely diagnosis -> treatment.
// As such, training on diagnoses and then later training on treatments may be a sensible approach
const (
	WhatIsDiagnosis LeadIn = iota // "What is the most likely diagnosis?"
	WhatIsTest                    // "What is the most appropriate next diagnostic test?"
	WhatIsTreatment               // "What is the most appropriate treatment?"
	WhatIsAction                  // "What is the most appropriate action?"
)

// Duration reflects the temporal course of a clinical finding
// This is simply to make our fake questions seem a bit more real
type Duration int

// Valid types of Duration
const (
	Unknown  Duration = iota // symptom onset is unknown / not specified
	Acute                    // the symptom came on acutely
	Subacute                 // the symptom came on subacutely
	Chronic                  // the symptom has been chronic
	Episodic                 // the symptom has been intermittent or episodic
)

func (d Duration) String() string {
	switch {
	case d == Acute:
		return "Acute"
	case d == Subacute:
		return "Subacute"
	case d == Chronic:
		return "Chronic"
	case d == Episodic:
		return "Episodic"
	default:
		return "Unknown"
	}
}

// ClinicalFinding combines a clinical finding SNOMED-CT concept and a duration
// e.g. acute chest pain
type ClinicalFinding struct {
	Concept  *snomed.Concept
	Duration Duration
}

func (cf ClinicalFinding) String() string {
	return cf.Duration.String() + " " + cf.Concept.FullySpecifiedName
}
