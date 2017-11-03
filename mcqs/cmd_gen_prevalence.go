package mcqs

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

// GeneratePrevalence generates fake prevalence data for each diagnostic concept.
// A real prevalence dataset would be based on mapping ICD-10 codes from inpatients or Read codes
// from primary care. Prevalence data would be different in different contexts.
func GeneratePrevalence(db *snomed.DatabaseService, number int) {
	rootDiagnosis, err := db.FetchConcept(SctDiagnosisRoot)
	checkError(err)
	allDiagnoses, err := db.FetchRecursiveChildren(rootDiagnosis)
	if number > 0 {
		allDiagnoses = allDiagnoses[0:min(number, len(allDiagnoses))]
	}
	results := make(map[snomed.Identifier]float64)
	writer := csv.NewWriter(os.Stdout)
	line := []string{"diagnosis", "concept_id", "prevalence", "parents"}
	writer.Write(line)
	for _, diagnosis := range allDiagnoses {
		prevalence := calculatePrevalence(db, results, diagnosis)
		line[0] = diagnosis.FullySpecifiedName
		line[1] = strconv.Itoa(int(diagnosis.ConceptID))
		line[2] = strconv.FormatFloat(prevalence, 'f', -1, 64)
		line[3] = snomed.ListItoA(diagnosis.Parents)
		writer.Write(line)
	}
	writer.Flush()
}

func printDiagnosisAndPrevalence(db *snomed.DatabaseService, concept *snomed.Concept, prevalence float64) {
	fmt.Printf("L%d: %s (%d) -- %f ", calculateLevelInHierarchy(db, concept), concept.FullySpecifiedName, concept.ConceptID, prevalence)
}

// calculatePrevalence creates a fake prevalence for a concept but ensures that, while random, that
// the prevalence data is internally consistent. It does this by recursively walking the IS-A hierarchy
// and ensuring that the total prevalence of a diagnosis and its siblings adds up to the total of the
// prevalence of the parents. When there are no parents (at the top of the diagnosis tree), then the
// prevalence is randomly distributed between the top siblings.
func calculatePrevalence(db *snomed.DatabaseService, results map[snomed.Identifier]float64, diagnosis *snomed.Concept) float64 {
	prevalence := results[diagnosis.ConceptID]
	if prevalence > 0 {
		return prevalence
	}
	if diagnosis.ConceptID == SctDiagnosisRoot {
		prevalence = 1.0
	} else {
		parents, err := db.GetParents(diagnosis)
		if err != nil {
			panic(err)
		}
		var parentPrevalence = 1.0
		for _, parent := range parents {
			parentPrevalence = math.Min(parentPrevalence, calculatePrevalence(db, results, parent))
		}
		prevalence = rand.Float64()*0.5 + 0.5
		if parentPrevalence > 0 {
			prevalence = prevalence * parentPrevalence / float64(calculateLevelInHierarchy(db, diagnosis))
		}
	}
	results[diagnosis.ConceptID] = prevalence
	return prevalence
}

// calculateLevelInHierarchy determines the number of steps away from the root SNOMED-CT concept
func calculateLevelInHierarchy(db *snomed.DatabaseService, concept *snomed.Concept) int {
	parents, err := db.GetParents(concept)
	if err != nil {
		log.Fatal(err)
	}
	if len(parents) == 0 {
		return 0
	}
	level := math.MaxInt8
	for _, parent := range parents {
		level = min(level, calculateLevelInHierarchy(db, parent))
	}
	return level + 1
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
