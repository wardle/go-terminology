package mcqs

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"bitbucket.org/wardle/go-snomed/snomed"
)

const (
	// SctDiagnosisRoot is the root concept of all diagnoses
	SctDiagnosisRoot = 64572001
	// SctFinding is the root concept of all clinical observations
	SctFinding = 404684003

	diagnosticConceptsFilename = "sct_diagnoses.csv"
	findingsConceptsFilename   = "sct_findings.csv"
	// a SQL statement to fetch all concepts belonging to a particular branch and return details with the cached parent concepts
	sqlSelectConceptFmt = `select concept_id, fully_specified_name, concept_status_code,
		string_agg(parent_concept_id::text,',') as parents
		from t_concept, t_cached_parent_concepts 
		where 
		child_concept_id=concept_id 
		and
		concept_id in (select child_concept_id from t_cached_parent_concepts where parent_concept_id=%d)
		group by concept_id`
)

// SnomedDataset is an abstract SNOMED-CT store that allows simple fetch.
type SnomedDataset interface {
	GetConcept(conceptID int) (*snomed.Concept, error)                     // fetch the concept specified
	Close()                                                                // close anything opened
	FetchRecursiveChildConcepts(root int) (map[int]*snomed.Concept, error) // fetch recursive child concepts from a root
}

// SimpleSnomedDataset is a crude unsophisticated in-memory representation of a SnomedDataset using Go's built in map
type SimpleSnomedDataset struct {
	path      string // path
	Diagnoses map[int]*snomed.Concept
	Problems  map[int]*snomed.Concept
}

// GetConcept is a crude in-memory cache tracking only diagnostic and observation type concepts
func (ssd SimpleSnomedDataset) GetConcept(conceptID int) (*snomed.Concept, error) {
	concept := ssd.Diagnoses[conceptID]
	if concept == nil {
		concept = ssd.Problems[conceptID]
		if concept == nil {
			return nil, fmt.Errorf("Could not find concept with identifier: %d. Is it a diagnosis or problem?", conceptID)
		}
	}
	return concept, nil
}

// FetchRecursiveChildConcepts is a crude implementation faking a future more sophisticated service...
func (ssd SimpleSnomedDataset) FetchRecursiveChildConcepts(root int) (map[int]*snomed.Concept, error) {
	switch {
	case root == SctDiagnosisRoot:
		return ssd.Diagnoses, nil
	case root == SctFinding:
		return ssd.Problems, nil
	default:
		return nil, fmt.Errorf("Error: not implemented: child concepts from root: %d", root)
	}
}

// Close is a NOP in this implementation
func (ssd SimpleSnomedDataset) Close() {
}

// NewSnomedDataset creates new file-based SNOMED data files at the path specified
func NewSnomedDataset(db *sql.DB, path string) (SnomedDataset, error) {
	error := os.MkdirAll(path, os.ModeDir)
	if error != nil {
		return nil, error
	}
	ds := &SimpleSnomedDataset{path: path}
	ds.Diagnoses, error = writeConceptsCsv(db, SctDiagnosisRoot, filepath.Join(path, diagnosticConceptsFilename))
	if error != nil {
		return nil, error
	}
	ds.Problems, error = writeConceptsCsv(db, SctFinding, filepath.Join(path, findingsConceptsFilename))
	if error != nil {
		return nil, error
	}
	return ds, nil
}

// OpenSnomedDataset opens an existing SNOMED-CT dataset
func OpenSnomedDataset(path string) (SnomedDataset, error) {
	ds := &SimpleSnomedDataset{path: path}
	var err error
	ds.Diagnoses, err = readConceptsCsv(filepath.Join(path, diagnosticConceptsFilename))
	if err != nil {
		return nil, err
	}
	ds.Problems, err = readConceptsCsv(filepath.Join(path, findingsConceptsFilename))
	return ds, err
}

// perform a query for all concepts within part of the IS-A hierarchy and write results to filename as csv
func writeConceptsCsv(db *sql.DB, root int, filename string) (map[int]*snomed.Concept, error) {
	concepts, err := fetchConcepts(db, root)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = writeToCsv(file, concepts)
	return concepts, err
}

// read csv file and turn into an in-memory concept map
func readConceptsCsv(filename string) (map[int]*snomed.Concept, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readFromCsv(f)
}

// conceptToCsv serialises a concept as a slice of strings
func conceptToCsv(c *snomed.Concept) []string {
	record := make([]string, 5)
	record[0] = strconv.Itoa(c.ConceptID)
	record[1] = c.FullySpecifiedName
	record[2] = strconv.Itoa(c.Status.Code)
	record[3] = listItoA(c.Parents)
	return record
}

// conceptFromCsv deserialises a concept from a slice of strings
func conceptFromCsv(row []string) (*snomed.Concept, error) {
	conceptID, err := strconv.Atoi(row[0])
	if err != nil {
		return nil, err
	}
	fullySpecifiedName := row[1]
	statusCode, err := strconv.Atoi(row[2])
	if err != nil {
		return nil, err
	}
	parents := listAtoi(row[3])
	return snomed.NewConcept(conceptID, fullySpecifiedName, statusCode, parents)
}

// write concepts to the writer in our proprietary CSV format
func writeToCsv(w io.Writer, concepts map[int]*snomed.Concept) error {
	w2 := csv.NewWriter(w)
	for _, concept := range concepts {
		record := conceptToCsv(concept)
		csvError := w2.Write(record)
		if csvError != nil {
			return csvError
		}
		w2.Flush()
	}
	return nil
}

// read concepts in our own proprietary CSV format
func readFromCsv(r io.Reader) (map[int]*snomed.Concept, error) {
	r2 := csv.NewReader(r)
	concepts := make(map[int]*snomed.Concept) // TODO: probably should be some abstract OO thing eventually..
	for {
		row, err := r2.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return concepts, err
		}
		concept, err := conceptFromCsv(row)
		if err == nil {
			concepts[concept.ConceptID] = concept
		}
	}
}

// convert a comma-delimited string containing integers into a slice of integers
func listAtoi(list string) []int {
	slist := strings.Split(strings.Replace(list, " ", "", -1), ",")
	r := make([]int, 0)
	for _, s := range slist {
		v, err := strconv.Atoi(s)
		if err == nil {
			r = append(r, v)
		}
	}
	return r
}

// convert a slice of integers into a comma-delimited string
func listItoA(list []int) string {
	r := make([]string, 0)
	for _, i := range list {
		s := strconv.Itoa(i)
		r = append(r, s)
	}
	return strings.Join(r, ",")
}

// fetch all of the concepts within SNOMED-CT beneath the given root
func fetchConcepts(db *sql.DB, root int) (map[int]*snomed.Concept, error) {
	sql := fmt.Sprintf(sqlSelectConceptFmt, root)
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	concepts := make(map[int]*snomed.Concept)
	var conceptID int
	var fullySpecifiedName string
	var conceptStatusCode int
	var parents string
	for rows.Next() {
		err = rows.Scan(&conceptID, &fullySpecifiedName, &conceptStatusCode, &parents)
		if err != nil {
			return nil, err
		}
		concept, err := snomed.NewConcept(conceptID, fullySpecifiedName, conceptStatusCode, listAtoi(parents))
		if err != nil {
			return nil, err
		}
		concepts[conceptID] = concept
	}
	return concepts, nil
}
