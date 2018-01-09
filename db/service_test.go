package db_test

import (
	"bitbucket.org/wardle/go-snomed/db"
	"golang.org/x/text/language"
	"os"
	"testing"
)

const (
	dbFilename = "../snomed.db"
)

func setUp(tb testing.TB) *db.Snomed {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) { // skip these tests if no working live snomed db
		tb.Skipf("Skipping tests against a live database. To run, create a database named %s", dbFilename)
	}
	bolt, err := db.NewBoltService(dbFilename, false)
	bleve := db.BleveService{}
	if err != nil {
		tb.Fatal(err)
	}
	return &db.Snomed{Store: bolt, Search: bleve, Language: language.BritishEnglish}
}

func TestService(t *testing.T) {
	snomed := setUp(t)
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

func BenchmarkGetConceptAndDescriptions(b *testing.B) {
	snomed := setUp(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms, err := snomed.GetConcept(24700007)
		if err != nil {
			b.Fatal(err)
		}
		_, err = snomed.GetDescriptions(ms)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIsA(b *testing.B) {
	snomed := setUp(b)
	ms, err := snomed.GetConcept(24700007)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	var (
		isDemyelinating  bool
		isPharmaceutical bool
	)
	for i := 0; i < b.N; i++ {
		isDemyelinating = snomed.IsA(ms, 6118003)
		isPharmaceutical = snomed.IsA(ms, 373873005)
	}
	if isDemyelinating == false || isPharmaceutical == true {
		b.Fatal("MS misclassified using IS-A hierarchy")
	}
}
