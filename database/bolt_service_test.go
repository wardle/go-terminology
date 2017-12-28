package database

import (
	"os"
	"testing"
	"bitbucket.org/wardle/go-snomed/snomed"
)

const (
	dbFilename = "bolt-tests.db"
)

func TestStore(t *testing.T) {
	bolt, err := NewBoltService(dbFilename)
	if err != nil {
		t.Fatal(err)
	}
	c1 := &snomed.Concept{ConceptID: 24700007, FullySpecifiedName: "Multiple sclerosis", Parents: nil, Status: snomed.Current.AsStatus()}
	bolt.PutConcepts(c1)
	c2, err := bolt.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	if c1.ConceptID != c2.ConceptID || c1.FullySpecifiedName != c2.FullySpecifiedName {
		t.Fatal("Concept not stored and retrieved correctly!")
	}
	c3, err := bolt.GetConcept(0)
	if c3 != nil && err != nil {
		t.Fatal("Failed to flag unfound concept")
	}
	os.Remove(dbFilename)
}
