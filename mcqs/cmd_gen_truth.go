package mcqs

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"encoding/json"
	"fmt"
	"math/rand"
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
func GenerateFakeTruth(db *snomed.Snomed, n int) {
	rootDiagnosis, err := db.FetchConcept(SctDiagnosisRoot)
	checkError(err)
	allDiagnoses, err := db.FetchRecursiveChildren(rootDiagnosis)
	checkError(err)
	diagnoses := allDiagnoses
	if n >= 0 {
		diagnoses = make([]*snomed.Concept, n) // use the specified number to limit to (n) diagnoses
		l := len(allDiagnoses)
		for i := 0; i < n; i++ {
			r := rand.Intn(l)
			diagnoses[i] = allDiagnoses[r]
		}
	}
	allTruth := make([]*FakeTruth, 0, len(diagnoses)+1)
	mi, err := MyocardialInfarctionTruth(db) // always prepend a "real" truth for illustrative purposes
	checkError(err)
	allTruth = append(allTruth, mi)
	for _, diag := range diagnoses {
		truth, ok := generateTruth(db, diag)
		if ok {
			allTruth = append(allTruth, truth)
		}
	}
	prevalence := make(map[snomed.Identifier]float64, 0)
	questions := make([]*Question, 0)
	for _, truth := range allTruth {
		p := 5 + int(calculatePrevalence(db, prevalence, truth.Diagnosis)*10000)*n // we'll impute for this diagnosis based on prevalence
		for i := 0; i < p; i++ {                                                   // generate number of questions commensurate with prevalence
			question := truth.ToQuestion(db)
			questions = append(questions, question)
		}
	}
	json, err := json.MarshalIndent(questions, "", "  ")
	checkError(err)
	fmt.Print(string(json))
}

func generateTruth(db *snomed.Snomed, diagnosis *snomed.Concept) (*FakeTruth, bool) {
	symptoms, err := relatedBySiteForDiagnosis(db, diagnosis)
	checkError(err)
	totalSymptoms := len(symptoms)
	if totalSymptoms > 0 {
		numSymptoms := 1 + rand.Intn(min(30, totalSymptoms))
		problems := make([]*FakeProblem, numSymptoms)
		parents, err := db.GetAllParents(diagnosis)
		checkError(err)
		for i := 0; i < numSymptoms; i++ {
			symptom := symptoms[rand.Intn(totalSymptoms-1)]
			problem := &FakeProblem{symptom, randomDuration(), rand.Float64()}
			problems[i] = problem
		}
		meanAge := randomAge()
		sd := min(meanAge, 20)
		truth := &FakeTruth{diagnosis, parents, problems, randomSexBias(), meanAge, rand.Intn(sd)}
		return truth, true
	}
	return nil, false
}

func randomDuration() Duration {
	possible := []Duration{Acute, Subacute, Chronic, Episodic}
	return possible[rand.Intn(len(possible)-1)]
}

func randomSexBias() SexBias {
	possible := []SexBias{MenOnly, FemaleOnly, NoSexBias}
	return possible[rand.Intn(len(possible)-1)]
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
	SexBias   SexBias           // does this disorder have a sex bias?
	MeanAge   int               // mean age
	StdDevAge int               // standard deviation for age for this disorder
}

// SexBias limits disorders to a gender, if appropriate
type SexBias int

// Possible options for SexBias
const (
	MenOnly SexBias = iota
	FemaleOnly
	NoSexBias
)

// RandomSex generates a random sex based on bias.
func (sb SexBias) RandomSex() Sex {
	switch {
	case sb == MenOnly:
		return Male
	case sb == FemaleOnly:
		return Female
	default:
		if rand.Float32() >= 0.5 {
			return Male
		}
		return Female
	}
}

func randomAge() int {
	if rand.Float32() > 0.01 {
		return rand.Intn(80) + 20
	}
	return rand.Intn(20)
}

func (ft FakeTruth) String() string {
	problems := make([]string, 0)
	for _, problem := range ft.Problems {
		problems = append(problems, problem.String())
	}
	return fmt.Sprintf("%s: %s", ft.Diagnosis.FullySpecifiedName, strings.Join(problems, ", "))
}

// ToQuestion creates a fake question from a fake truth by choosing a random selection of the symptoms on offer.
func (ft FakeTruth) ToQuestion(db *snomed.Snomed) *Question {
	findings := make([]*ClinicalFinding, 0)
	for _, problem := range ft.Problems {
		if problem.Probability > rand.Float64() {
			findings = append(findings, problem.ToFinding(db))
		}
	}
	age := randomAge()
	if ft.MeanAge > 0 && ft.StdDevAge > 0 {
		age = int(rand.NormFloat64()*float64(ft.StdDevAge) + float64(ft.MeanAge))
		if age < 0 {
			age = 0
		}
	}
	sex := ft.SexBias.RandomSex()
	parents, err := db.GetAllParents(ft.Diagnosis)
	checkError(err)
	return &Question{Age: age, Sex: sex, Findings: findings, LeadIn: WhatIsDiagnosis,
		PossibleAnswers: nil, CorrectAnswer: ft.Diagnosis, Parents: parents}
}

// FakeProblem records a clinical finding or observation and its probability
// for an owning Diagnosis.
type FakeProblem struct {
	Problem     *snomed.Concept // problem
	Duration    Duration        // duration
	Probability float64         // probability of this problem for this condition
}

func (fp FakeProblem) String() string {
	return fmt.Sprintf("%s (%f%%)", fp.Problem.FullySpecifiedName, fp.Probability)
}

// ToFinding turns a fake problem from a fake truth into a clinical finding
func (fp FakeProblem) ToFinding(db *snomed.Snomed) *ClinicalFinding {
	parents, err := db.GetAllParents(fp.Problem)
	checkError(err)
	return &ClinicalFinding{fp.Problem, parents, fp.Duration}
}

// convenience structure to allow literal defined truth for demonstration purposes.
type explicitTruth struct {
	diagnosis snomed.Identifier
	problems  []*explicitProblem
	meanAge   int
	stdDevAge int
}

// convenience structure to allow literal defined problem for demonstration purposes.
type explicitProblem struct {
	conceptID   snomed.Identifier
	duration    Duration
	probability float64
}

// toFakeTruth converts a (usually literal defined) explicit truth into a fake truth
func (et explicitTruth) toFakeTruth(db *snomed.Snomed) (*FakeTruth, error) {
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
	return &FakeTruth{diagnosis, parents, problems, NoSexBias, et.meanAge, et.stdDevAge}, nil
}

// toFakeProblem converts a (usually literal defined) explicit problem into a fake problem
func (ep explicitProblem) toFakeProblem(db *snomed.Snomed) (*FakeProblem, error) {
	concept, err := db.FetchConcept(int(ep.conceptID))
	if err != nil {
		return nil, err
	}
	return &FakeProblem{concept, ep.duration, ep.probability}, nil
}

var myocardialInfarction = &explicitTruth{22298006,
	[]*explicitProblem{
		&explicitProblem{29857009, Acute, 0.95},  // chest pain
		&explicitProblem{267036007, Acute, 0.70}, // breathlessness
		&explicitProblem{415690000, Acute, 0.80}, // sweating
		&explicitProblem{426555006, Acute, 0.55}, // paint ot jaw
		&explicitProblem{76388001, Acute, 0.60},  // ST elevation on ECG - this will inherently say "ECG abnormal"
	}, 60, 20}

// MyocardialInfarctionTruth generates a truth for myocardial infarction for demonstration and testing purposes.
func MyocardialInfarctionTruth(db *snomed.Snomed) (*FakeTruth, error) {
	return myocardialInfarction.toFakeTruth(db)
}

// RelatedBySiteForDiagnosis is a hacky way of getting a relatively reasonable list of clinical
// findings for any arbitrary diagnosis by walking the SNOMED-CT ontology by finding site and finding
// clinical findings for that site. It isn't at all perfect, but might make it look authentic to a non-medic!
func relatedBySiteForDiagnosis(dbs *snomed.Snomed, concept *snomed.Concept) ([]*snomed.Concept, error) {
	sites, err := dbs.GetParentsOfKind(concept, snomed.FindingSite) // where is this disease?
	if err != nil {
		return nil, err
	}
	allSymptoms := make(map[snomed.Identifier]*snomed.Concept)
	thoracic, err := dbs.FetchConcept(51185008)  // high-level structure
	structures, err := dbs.GetSiblings(thoracic) // get similiar high-level structures
	structures = append(structures, thoracic)
	structures2 := snomed.SliceToMap(structures)
	genericSites := make([]*snomed.Concept, 0)
	for _, site := range sites {
		genericSite, ok := dbs.Genericise(site, structures2)
		if ok {
			genericSites = append(genericSites, genericSite)
		}
	}
	for _, site := range genericSites {
		allChildren, _ := dbs.FetchRecursiveChildren(site)
		for _, child := range allChildren {
			symptoms, err := dbs.GetChildrenOfKind(child, snomed.FindingSite)
			if err != nil {
				return nil, err
			}
			for _, symptom := range symptoms {
				if symptom.IsA(snomed.SctDisease) == false {
					allSymptoms[symptom.ConceptID] = symptom
				}
			}
		}
	}
	return snomed.MapToSlice(allSymptoms), nil
}

/*
// symptomsForDiagnosis returns some possible symptoms for an arbitrary diagnosis by
// navigating the SNOMED ontology and finding symptoms related to sites of the diagnosis in question
func symptomsForDiagnosis(db *snomed.DatabaseService, concept *snomed.Concept, minimum int) []*snomed.Concept {
	sites, err := db.GetParentsOfKind(concept, snomed.FindingSite) // where is this disease?
	results := make([]*snomed.Concept, 0)
	for _, site := range sites {
		symptoms := symptomsForSites
		results = append(results, symptoms...)
	}
}
*/
