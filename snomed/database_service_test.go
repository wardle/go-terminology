package snomed

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/text/language"
	"reflect"
	"strings"
	"testing"
)

// database connection parameters for testing
const (
	dbDriver   = "postgres"
	dbUser     = "mark"
	dbPassword = ""
	dbName     = "rsdb"
)

func setUp(tb testing.TB) (dbs *Snomed) {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open(dbDriver, dbinfo)
	if err != nil {
		tb.Fatal(err)
	}
	return &Snomed{Service: NewDatabaseService(db), Language: language.BritishEnglish}
}
func shutDown(dbs *Snomed) {
	dbs.Close()
}

func TestMultipleSclerosis(t *testing.T) {
	sct := setUp(t)
	ms, err := sct.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	if ms.FullySpecifiedName != "Multiple sclerosis (disorder)" {
		t.Error("Incorrect concept.")
	}
	children, err := sct.GetRecursiveChildren(ms)
	if err != nil {
		t.Fatal(err)
	}
	for _, child := range children {
		if child.IsA(ms.ConceptID) == false {
			t.Errorf("Concept %s not correctly identified as type of %s", child, ms)
		}
	}
	msRelations, err := sct.GetParentRelationships(ms)
	if err != nil {
		t.Fatal(err)
	}
	for _, relation := range msRelations {
		_, _, _, err := sct.ConceptsForRelationship(relation)
		if err != nil {
			t.Fatal(err)
		}
	}
	kinds, err := sct.GetParents(ms)
	var isDemyelination = false
	for _, kind := range kinds {
		if kind.ConceptID == 6118003 {
			isDemyelination = true
			children, err := sct.GetChildren(kind)
			if err != nil {
				t.Fatal(err)
			}
			var found = false
			for _, child := range children {
				if child.ConceptID == ms.ConceptID {
					found = true
				}
			}
			if !found {
				t.Error("Multiple sclerosis not a child of demyelinating disorder!")
			}
		}
	}
	if isDemyelination == false {
		t.Error("Multiple sclerosis not correctly identified as a demyelinating disorder")
	}
	parents, err := sct.GetAllParents(ms)
	if err != nil {
		t.Fatal(err)
	}
	if len(parents) == 0 {
		t.Error("Invalid number of parent concepts for an individual concept")
	}
	shutDown(sct)
}

func TestDescriptions(t *testing.T) {
	sct := setUp(t)
	ms, err := sct.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	desc, err := sct.GetPreferredDescription(ms)
	if err != nil {
		t.Fatal(err)
	}
	if desc.Term != "Multiple sclerosis" {
		t.Fatal("Did not find correct synonym for multiple sclerosis concept")
	}
	shutDown(sct)
}
func TestInvalidIdentifier(t *testing.T) {
	sct := setUp(t)
	_, err := sct.GetConcept(0)
	if err == nil {
		t.Fatal("Should throw an error if a concept is not found.")
	}
	shutDown(sct)
}

func TestMultipleFetch(t *testing.T) {
	sct := setUp(t)
	msAndPd, err := sct.GetConcepts(24700007, 49049000)
	if err != nil {
		t.Fatal(err)
	}
	if len(msAndPd) != 2 {
		t.Fatal("Did not correctly fetch multiple sclerosis and Parkinson's disease!")
	}
	shutDown(sct)
}

func TestRoot(t *testing.T) {
	sct := setUp(t)
	root, err := sct.GetConcept(138875005)
	if err != nil {
		t.Fatal(err)
	}
	rootParents, err := sct.GetAllParents(root)
	if err != nil {
		t.Error(err)
	}
	if len(rootParents) != 0 {
		t.Error("Invalid number of parent concepts for root concept")
	}
	shutDown(sct)
}

func TestPathsToRoot(t *testing.T) {
	sct := setUp(t)
	ms, err := sct.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	paths, err := sct.PathsToRoot(ms)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		if path[0].ConceptID != 24700007 {
			t.Error("Path doesn't include origin concept")
		}
		if path[len(path)-1].ConceptID != 138875005 {
			t.Error("Path doesn't include root concept")
		}
	}
	shutDown(sct)
}

func BenchmarkPathsToRoot(b *testing.B) {
	sct := setUp(b)
	ms, err := sct.GetConcept(24700007)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		sct.PathsToRoot(ms)
	}
	shutDown(sct)
}

func TestGenericise(t *testing.T) {
	sct := setUp(t)
	ms, err := sct.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}

	cnsType, err := sct.GenericiseToRoot(ms, SctCentralNervousSystemDisease) // what type of CNS disease is this?
	if err != nil {
		t.Fatal(err)
	}

	if cnsType.ConceptID != 6118003 {
		t.Errorf("Multiple sclerosis not correctly genericised to a demyelinating disorder of the central nervous system")
	}
	shutDown(sct)

}

func TestListAtoi(t *testing.T) {
	testAtoi(t, "123,456,789,123456789", []int{123, 456, 789, 123456789}, true)
	testAtoi(t, "123,456,789, 123456789", []int{123, 456, 789, 123456789}, true)
	testAtoi(t, "aaa,vbb,cc,123", []int{123}, false)
	testAtoi(t, "", []int{}, true)
}

// test conversion of a comma-delimited string of integers to a slice, optionally testing roundtripping
func testAtoi(t *testing.T, input string, expected []int, roundtrip bool) {
	r := ListAtoi(input)
	if reflect.DeepEqual(r, expected) == false {
		t.Errorf("Failed to parse: %s. Parsed to %v", input, r)
	}
	if roundtrip {
		v := ListItoA(expected)
		if reflect.DeepEqual(v, strings.Replace(input, " ", "", -1)) == false {
			t.Errorf("Failed to parse: %v. Parsed to %s", expected, v)
		}
	}
}
