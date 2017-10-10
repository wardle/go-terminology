package mcqs

import (
	"database/sql"
	"fmt"
	"github.com/wardle/go-snomed/snomed"
	"math/rand"
	"strconv"
	"strings"
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
	group by concept_id limit 10`
	sctDiagnosisRoot           = 64572001  // root concept of all diagnoses
	sctClinicalObservationRoot = 250171008 // root concept of all clinical observations
)

// GenerateSnomedCT generates a set of useful intermediate files containing pre-computed and cached
// SNOMED-CT data, useful for sharing with others without depending on a dedicated terminology server.
func GenerateSnomedCT() {
	fmt.Println("Connecting to database..")
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open(dbDriver, dbinfo)
	checkErr(err)
	defer db.Close()

	fmt.Println("# Querying")
	diagnoses := fetchConcepts(db, sctDiagnosisRoot)
	clinicalFindings := fetchConcepts(db, sctClinicalObservationRoot)
	nClinicalFindings := len(clinicalFindings)
	fmt.Printf("Number of diagnoses: %d\n", len(diagnoses))
	fmt.Printf("Number of clinical observations: %d\n", nClinicalFindings)

	fmt.Println("Writing csv...")
	//w := csv.NewWriter(os.Stdout)
	// randomly select a few clinical observations for each diagnosis.
	// this is fake data obviously!
	for _, diagnosis := range diagnoses {
		numFeatures := rand.Intn(10)
		features := make([]*snomed.Concept, numFeatures)
		for i := 0; i < numFeatures; i++ {
			features[i] = clinicalFindings[rand.Intn(nClinicalFindings)]
		}
		fmt.Println(diagnosis.FullySpecifiedName)
		fmt.Println(features)
	}
	//fmt.Printf("Status: %s", snomed.LookupStatus(1))

}

// fetch all of the concepts within SNOMED-CT beneath the given root
func fetchConcepts(db *sql.DB, root int) map[int]*snomed.Concept {
	sql := fmt.Sprintf(sqlSelectConceptFmt, root)
	fmt.Print(sql)
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

// convert a comma-delimited list into a slice of integers
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
