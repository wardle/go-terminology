package snomed

import (
	"database/sql"
	"fmt"
	"math/rand"
)

// define simple connection settings statically for this hack. TODO: move to configuration file
const (
	dbDriver                   = "postgres"
	dbUser                     = "mark"
	dbPassword                 = ""
	dbName                     = "rsdb"
	sqlSelectConceptFmt        = "select concept_id, fully_specified_name, concept_status_code from t_concept where concept_id in (select child_concept_id from t_cached_parent_concepts where parent_concept_id=%d)"
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
		features := make([]*Concept, numFeatures)
		for i := 0; i < numFeatures; i++ {
			features[i] = clinicalFindings[rand.Intn(nClinicalFindings)]
		}
		fmt.Println(diagnosis.FullySpecifiedName)
		fmt.Println(features)
	}
	//fmt.Printf("Status: %s", snomed.LookupStatus(1))

}

// fetch all of the concepts within SNOMED-CT beneath the given root
func fetchConcepts(db *sql.DB, root int) map[int]*Concept {
	sql := fmt.Sprintf(sqlSelectConceptFmt, root)
	rows, err := db.Query(sql)
	checkErr(err)
	concepts := make(map[int]*Concept)
	for rows.Next() {
		var conceptID int
		var fullySpecifiedName string
		var conceptStatusCode int
		err = rows.Scan(&conceptID, &fullySpecifiedName, &conceptStatusCode)
		checkErr(err)
		concept, err := CreateConcept(conceptID, fullySpecifiedName, conceptStatusCode, nil)
		if err != nil {
			panic(err)
		}
		concepts[conceptID] = concept
	}
	return concepts
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
