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
// TODO: use a file-based database for tests
const (
	dbDriver   = "postgres"
	dbUser     = "mark"
	dbPassword = ""
	dbName     = "rsdb"
)

func openConnection(t *testing.T) *sql.DB {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open(dbDriver, dbinfo)
	if err != nil {
		t.Error(err)
	}
	return db
}

func TestConnection(t *testing.T) {
	db := openConnection(t)
	defer db.Close()
	snomed := NewDatabaseService(db)
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
	parents, err := snomed.GetAllParents(ms)
	if err != nil {
		t.Fatal(err)
	}
	if len(parents) == 0 {
		t.Error("Invalid number of parent concepts for an individual concept")
	}

	_, err = snomed.FetchConcept(0)
	if err == nil {
		//t.Fatal("Should throw an error if a concept is not found.")
	}
	mspd, err := snomed.FetchConcepts(24700007, 49049000)
	if err != nil {
		t.Fatal(err)
	}
	if (snomed.cache.cache[24700007] != ms) || (snomed.cache.cache[49049000] != mspd[1]) || snomed.cache.cache[26823002] != nil {
		t.Error("Concepts not cached correctly.")
	}
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
