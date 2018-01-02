package db_test

import (
	"bitbucket.org/wardle/go-snomed/db"
	"golang.org/x/text/language"
	//"bitbucket.org/wardle/go-snomed/rf2"
	"os"
	"testing"
)

const (
	dbFilename = "../snomed.db"
)

func hasLiveDatabase(t *testing.T) {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) { // skip these tests if no working live snomed db
		t.Skipf("Skipping tests against a live database. To run, create a database named %s", dbFilename)
	}
}

func TestService(t *testing.T) {
	hasLiveDatabase(t)
	bolt, err := db.NewBoltService(dbFilename)
	if err != nil {
		t.Fatal(err)
	}
	snomed := &db.Snomed{Service: bolt, Language: language.BritishEnglish}
	ms, err := snomed.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	parents, err := snomed.GetAllParents(ms)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range parents {
		if p.ID == 6118003 {
			found = true
		}
	}
	if !found {
		t.Fatal("Multiple sclerosis not correctly identified as a type of demyelinating disease")
	}
	if !snomed.IsA(ms, 6118003) {
		t.Fatal("Multiple sclerosis not correctly identified as a type of demyelinating disease")
	}
}
