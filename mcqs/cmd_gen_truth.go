package mcqs

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"fmt"
	"strings"
)

// FakeTruth is an intermediate transitional data structure used to generate
// multiple questions from that same truth. The idea is to have "myocardial infarction"
// represented by "chest pain" (95%), "breathlessness" (60%), "sweating" (40%), "ECG: ST elevation" (80%)
// (figures make up) but of course, our data will be fake but at least relatively consistent as questions
// will be generated from the same "fake truth" but with different combinations of problems.
type FakeTruth struct {
	Diagnosis *snomed.Concept
	Parents   []*snomed.Concept // convenience pointers to parents
	Problems  []*FakeProblem    // problems for this diagnosis
}

func (ft FakeTruth) String() string {
	problems := make([]string, 0)
	for _, problem := range ft.Problems {
		problems = append(problems, problem.String())
	}
	return fmt.Sprintf("%s: %s", ft.Diagnosis.FullySpecifiedName, strings.Join(problems, ", "))
}

// FakeProblem records a clinical finding or observation and its probability
// for an owning Diagnosis.
type FakeProblem struct {
	Problem     *snomed.Concept // problem
	Probability int             // probability of this problem for this condition
}

func (fp FakeProblem) String() string {
	return fmt.Sprintf("%s (%d%%)", fp.Problem.FullySpecifiedName, fp.Probability)
}

// GenerateFakeTruth takes the SNOMED-CT ontology and uses it to build a fake "truth" model
// that represents the clinical findings that are seen in each type of diagnosis.
//
// While simply generating random problems for each diagnosis might be one approach, it is incorrect as
// we have a clear subsumption IS-A hierarchy which can be used. As such, related diagnostic concepts
// should share similar clinical problems in order to generate reasonable fake data.
func GenerateFakeTruth(dataset SnomedDataset) {
	mi, err := MyocardialInfarctionTruth(dataset)
	fmt.Print(mi, err)
}

// convenience structure to allow literal defined truth for demonstration purposes.
type explicitTruth struct {
	diagnosis int
	problems  []*explicitProblem
}

// convenience structure to allow literal defined problem for demonstration purposes.
type explicitProblem struct {
	conceptID   int
	probability int
}

// toFakeTruth converts a (usually literal defined) explicit truth into a fake truth
func (et explicitTruth) toFakeTruth(dataset SnomedDataset) (*FakeTruth, error) {
	diagnosis, err := dataset.GetConcept(et.diagnosis)
	if err != nil {
		return nil, err
	}
	problems := make([]*FakeProblem, 0, len(et.problems))
	for _, p := range et.problems {
		fp, err := p.toFakeProblem(dataset)
		if err != nil {
			return nil, err
		}
		problems = append(problems, fp)
	}
	parents, err := getConcepts(dataset, false, diagnosis.Parents...)
	if err != nil {
		return nil, err
	}
	return &FakeTruth{diagnosis, parents, problems}, nil
}

// toFakeProblem converts a (usually literal defined) explicit problem into a fake problem
func (ep explicitProblem) toFakeProblem(dataset SnomedDataset) (*FakeProblem, error) {
	concept, err := dataset.GetConcept(ep.conceptID)
	if err != nil {
		return nil, err
	}
	return &FakeProblem{concept, ep.probability}, nil
}

var myocardialInfarction = &explicitTruth{22298006,
	[]*explicitProblem{
		&explicitProblem{29857009, 95},  // chest pain
		&explicitProblem{267036007, 70}, // breathlessness
		&explicitProblem{415690000, 80}, // sweating
		&explicitProblem{426555006, 55}, // paint ot jaw
		&explicitProblem{76388001, 60},  // ST elevation on ECG - this will inherently say "ECG abnormal"
	}}

// MyocardialInfarctionTruth generates a truth for myocardial infarction for demonstration and testing purposes.
func MyocardialInfarctionTruth(dataset SnomedDataset) (*FakeTruth, error) {
	return myocardialInfarction.toFakeTruth(dataset)
}

// FetchMode defines whether to ignore missing or incorrect identifiers during batch operations
type FetchMode int

// Valid types of FetchMode
const (
	Strict  FetchMode = iota // Strict FetchMode will raise an error
	Relaxed                  // Relaxed FetchMode will ignore missing or incorrect identifiers

)

// convenience function to get a list of concepts
func getConcepts(dataset SnomedDataset, strict bool, conceptIDs ...int) ([]*snomed.Concept, error) {
	result := make([]*snomed.Concept, 0, len(conceptIDs))
	for _, conceptID := range conceptIDs {
		concept, err := dataset.GetConcept(conceptID)
		if err != nil {
			var sctid = snomed.Identifier(conceptID)
			if strict || sctid.IsValid() == false || sctid.IsConcept() == false {
				return nil, err
			}
		}
		result = append(result, concept)
	}
	return result, nil
}
