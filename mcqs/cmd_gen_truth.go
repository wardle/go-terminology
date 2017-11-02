package mcqs

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"fmt"
	"strings"
)

const (
	// SctDiagnosisRoot is the root concept of all diagnoses
	SctDiagnosisRoot = 64572001
	// SctFinding is the root concept of all clinical observations
	SctFinding = 404684003
)

// GenerateFakeTruth takes the SNOMED-CT ontology and uses it to build a fake "truth" model
// that represents the clinical findings that are seen in each type of diagnosis.
//
// While simply generating random problems for each diagnosis might be one approach, it is incorrect as
// we have a clear subsumption IS-A hierarchy which can be used. As such, related diagnostic concepts
// should share similar clinical problems in order to generate reasonable fake data.
func GenerateFakeTruth(db *snomed.DatabaseService) {
	rootDiagnosis, err := db.FetchConcept(SctDiagnosisRoot)
	checkError(err)
	allDiagnoses, err := db.FetchRecursiveChildren(rootDiagnosis)
	fmt.Printf("Fetched %d diagnoses....\n", len(allDiagnoses))
	mi, err := MyocardialInfarctionTruth(db)
	if err != nil {
		panic(err)
	}
	fmt.Print(mi)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

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
	Duration    Duration        // duration
	Probability int             // probability of this problem for this condition
}

func (fp FakeProblem) String() string {
	return fmt.Sprintf("%s (%d%%)", fp.Problem.FullySpecifiedName, fp.Probability)
}

// convenience structure to allow literal defined truth for demonstration purposes.
type explicitTruth struct {
	diagnosis snomed.Identifier
	problems  []*explicitProblem
}

// convenience structure to allow literal defined problem for demonstration purposes.
type explicitProblem struct {
	conceptID   snomed.Identifier
	duration    Duration
	probability int
}

// toFakeTruth converts a (usually literal defined) explicit truth into a fake truth
func (et explicitTruth) toFakeTruth(db *snomed.DatabaseService) (*FakeTruth, error) {
	diagnosis, err := db.FetchConcept(int(et.diagnosis))
	if err != nil {
		return nil, err
	}
	problems := make([]*FakeProblem, 0, len(et.problems))
	for _, p := range et.problems {
		fp, err := p.toFakeProblem(db)
		if err != nil {
			return nil, err
		}
		problems = append(problems, fp)
	}
	parents, err := db.GetAllParents(diagnosis)
	if err != nil {
		return nil, err
	}
	return &FakeTruth{diagnosis, parents, problems}, nil
}

// toFakeProblem converts a (usually literal defined) explicit problem into a fake problem
func (ep explicitProblem) toFakeProblem(db *snomed.DatabaseService) (*FakeProblem, error) {
	concept, err := db.FetchConcept(int(ep.conceptID))
	if err != nil {
		return nil, err
	}
	return &FakeProblem{concept, ep.duration, ep.probability}, nil
}

var myocardialInfarction = &explicitTruth{22298006,
	[]*explicitProblem{
		&explicitProblem{29857009, Acute, 95},  // chest pain
		&explicitProblem{267036007, Acute, 70}, // breathlessness
		&explicitProblem{415690000, Acute, 80}, // sweating
		&explicitProblem{426555006, Acute, 55}, // paint ot jaw
		&explicitProblem{76388001, Acute, 60},  // ST elevation on ECG - this will inherently say "ECG abnormal"
	}}

// MyocardialInfarctionTruth generates a truth for myocardial infarction for demonstration and testing purposes.
func MyocardialInfarctionTruth(db *snomed.DatabaseService) (*FakeTruth, error) {
	return myocardialInfarction.toFakeTruth(db)
}
