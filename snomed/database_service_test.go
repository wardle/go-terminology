package snomed

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

func setUp(tb testing.TB) (db *sql.DB, dbs *DatabaseService) {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open(dbDriver, dbinfo)
	if err != nil {
		tb.Fatal(err)
	}
	return db, NewDatabaseService(db)
}
func shutDown(db *sql.DB) {
	db.Close()
}

func TestMultipleSclerosis(t *testing.T) {
	db, snomed := setUp(t)
	ms, err := snomed.FetchConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	if ms.FullySpecifiedName != "Multiple sclerosis (disorder)" {
		t.Error("Incorrect concept.")
	}
	children, err := snomed.FetchRecursiveChildren(ms)
	if err != nil {
		t.Fatal(err)
	}
	for _, child := range children {
		if child.IsA(ms.ConceptID) == false {
			t.Errorf("Concept %s not correctly identified as type of %s", child, ms)
		}
	}
	msRelations, err := snomed.FetchParentRelationships(ms)
	if err != nil {
		t.Fatal(err)
	}
	for _, relation := range msRelations {
		_, _, _, err := snomed.ConceptsForRelationship(relation)
		if err != nil {
			t.Fatal(err)
		}
	}
	kinds, err := snomed.GetParents(ms)
	var isDemyelination = false
	for _, kind := range kinds {
		if kind.ConceptID == 6118003 {
			isDemyelination = true
			children, err := snomed.GetChildren(kind)
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
	parents, err := snomed.GetAllParents(ms)
	if err != nil {
		t.Fatal(err)
	}
	if len(parents) == 0 {
		t.Error("Invalid number of parent concepts for an individual concept")
	}
	shutDown(db)
}

func TestDescriptions(t *testing.T) {
	db, snomed := setUp(t)
	ms, err := snomed.FetchConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	desc, err := snomed.GetPreferredDescription(ms)
	if err != nil {
		t.Fatal(err)
	}
	if desc.Term != "Multiple sclerosis" {
		t.Fatal("Did not find correct synonym for multiple sclerosis concept")
	}
	shutDown(db)
}
func TestInvalidIdentifier(t *testing.T) {
	db, snomed := setUp(t)
	_, err := snomed.FetchConcept(0)
	if err == nil {
		t.Fatal("Should throw an error if a concept is not found.")
	}
	shutDown(db)
}

func TestMultipleFetch(t *testing.T) {
	db, snomed := setUp(t)
	msAndPd, err := snomed.FetchConcepts(24700007, 49049000)
	if err != nil {
		t.Fatal(err)
	}
	if len(msAndPd) != 2 {
		t.Fatal("Did not correctly fetch multiple sclerosis and Parkinson's disease!")
	}
	shutDown(db)
}

func TestRoot(t *testing.T) {
	db, snomed := setUp(t)
	root, err := snomed.FetchConcept(138875005)
	if err != nil {
		t.Fatal(err)
	}
	rootParents, err := snomed.GetAllParents(root)
	if err != nil {
		t.Error(err)
	}
	if len(rootParents) != 0 {
		t.Error("Invalid number of parent concepts for root concept")
	}
	shutDown(db)
}

func TestPathsToRoot(t *testing.T) {
	db, snomed := setUp(t)
	ms, err := snomed.FetchConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	paths, err := snomed.PathsToRoot(ms)
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
	shutDown(db)
}

func BenchmarkPathsToRoot(b *testing.B) {
	db, snomed := setUp(b)
	ms, err := snomed.FetchConcept(24700007)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		snomed.pathsToRoot(ms) // testing non-cache version obviously
	}
	shutDown(db)
}

func TestGenericise(t *testing.T) {
	db, snomed := setUp(t)
	ms, err := snomed.FetchConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}

	cnsType, err := snomed.GenericiseToRoot(ms, SctCentralNervousSystemDisease) // what type of CNS disease is this?
	if err != nil {
		t.Fatal(err)
	}

	if cnsType.ConceptID != 6118003 {
		t.Errorf("Multiple sclerosis not correctly genericised to a demyelinating disorder of the central nervous system")
	}
	shutDown(db)

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
