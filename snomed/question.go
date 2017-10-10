package snomed

// Question is made up of a scenario, a lead-in, five possible answers and one single best answer.
type Question struct {
	Scenario Scenario
}

// Scenario consists of a list of clinical findings
type Scenario struct {
	Concepts []Concept
}

// Duration reflects the temporal course of a clinical finding
type Duration int

const (
	Acute Duration = iota
	Subacute
	Chronic
)

// ClinicalFinding combines a clinical finding SNOMED-CT concept and a duration
// e.g. acute chest pain
type ClinicalFinding struct {
	Concept  Concept
	Duration Duration
}
