package mcqs

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/wardle/go-snomed/snomed"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// define simple connection settings statically for this hack. TODO: move to configuration file
const (
	dbDriver   = "postgres"
	dbUser     = "mark"
	dbPassword = ""
	dbName     = "rsdb"
	// a SQL statement to fetch all concepts belonging to a particular branch and return details and the cached parent concepts
	sqlSelectConceptFmt = `select concept_id, fully_specified_name, concept_status_code,
	string_agg(parent_concept_id::text,',') as parents
	from t_concept, t_cached_parent_concepts 
	where 
	child_concept_id=concept_id 
	and
	concept_id in (select child_concept_id from t_cached_parent_concepts where parent_concept_id=%d)
	group by concept_id`
	sctDiagnosisRoot           = 64572001  // root concept of all diagnoses
	sctClinicalObservationRoot = 250171008 // root concept of all clinical observations
)

// GenerateSnomedCT generates a set of useful intermediate files containing pre-computed and cached
// SNOMED-CT data, useful for sharing with others without depending on a dedicated terminology server.
func GenerateSnomedCT() {
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	spinner.Start()
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open(dbDriver, dbinfo)
	checkErr(err)
	defer db.Close()

	writeConceptsCsv(db, sctDiagnosisRoot, "sct_diagnoses.csv")
	spinner.Reverse()
	writeConceptsCsv(db, sctClinicalObservationRoot, "sct_findings.csv")
	spinner.Stop()
}

// perform a query for all concepts within part of the IS-A hierarchy and write results to filename as csv
func writeConceptsCsv(db *sql.DB, root int, filename string) {
	concepts := fetchConcepts(db, root)
	file, err := os.Create(filename)
	defer file.Close()
	checkErr(err)
	writeToCsv(file, concepts)

}

func (c snomed.Concept) toCsv() []string {
	record := make([]string, 5)
	record[0] = strconv.Itoa(concept.ConceptID)
	record[1] = concept.FullySpecifiedName
	record[2] = concept.Status.Title
	record[3] = listItoA(concept.Parents)
	parentFsns := make([]string, 0)
	for _, parentID := range concept.Parents {
		parent := concepts[parentID]
		if parent != nil {
			parentFsns = append(parentFsns, parent.FullySpecifiedName)
		}
	}
	record[4] = strings.Join(parentFsns, ",")
	return record
}

// write the concepts to the writer as a CSV file
func writeToCsv(w io.Writer, concepts map[int]*snomed.Concept) {
	w2 := csv.NewWriter(w)
	for _, concept := range concepts {
		record := concept.toCsv()
		csvError := w2.Write(record)
		if csvError != nil {
			log.Fatalf("Failed to write CSV file: %s", csvError)
		}
		w2.Flush()
	}
	//fmt.Printf("Status: %s", snomed.LookupStatus(1))
}

// fetch all of the concepts within SNOMED-CT beneath the given root
func fetchConcepts(db *sql.DB, root int) map[int]*snomed.Concept {
	sql := fmt.Sprintf(sqlSelectConceptFmt, root)
	rows, err := db.Query(sql)
	checkErr(err)
	concepts := make(map[int]*snomed.Concept)
	for rows.Next() {
		var conceptID int
		var fullySpecifiedName string
		var conceptStatusCode int
		var parents string
		err = rows.Scan(&conceptID, &fullySpecifiedName, &conceptStatusCode, &parents)
		checkErr(err)
		concept, err := snomed.CreateConcept(conceptID, fullySpecifiedName, conceptStatusCode, listAtoi(parents))
		if err != nil {
			panic(err)
		}
		concepts[conceptID] = concept
	}
	return concepts
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

// panic if there is an error
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
