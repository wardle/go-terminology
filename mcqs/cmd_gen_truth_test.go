package mcqs

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/text/language"
	"testing"
)

// database connection parameters for testing
const (
	dbDriver   = "postgres"
	dbUser     = "mark"
	dbPassword = ""
	dbName     = "rsdb"
)

func setUp(t *testing.T) (*sql.DB, *snomed.Snomed) {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open(dbDriver, dbinfo)
	if err != nil {
		t.Fatal(err)
	}
	return db, &snomed.Snomed{Service: snomed.NewDatabaseService(db), Language: language.BritishEnglish}
}

func shutDown(db *sql.DB, service *snomed.Snomed) {
	db.Close()
}

const ()

func TestPossibleSymptoms(t *testing.T) {
	db, snomed := setUp(t)
	mi, err := snomed.FetchConcept(22298006) // myocardial infarction
	if err != nil {
		t.Fatal(err)
	}
	symptoms, err := relatedBySiteForDiagnosis(snomed, mi)
	if err != nil {
		t.Fatal(err)
	}
	foundChestPain := false
	for _, sx := range symptoms {
		if sx.ConceptID.AsInteger() == 279019008 { // did we identify crushing chest pain as possible symptom?
			foundChestPain = true
		}
	}
	if !foundChestPain {
		t.Fatal("Crushing chest pain not correctly identified as a possible symptom in MI")
	}
	shutDown(db, snomed)
}
