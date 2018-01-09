package db

import (
	"bitbucket.org/wardle/go-snomed/rf2"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	dbFilename = "bolt-tests.db"
)

func TestStore(t *testing.T) {
	bolt, err := NewBoltService(dbFilename, false)
	if err != nil {
		t.Fatal(err)
	}
	d, err := time.Parse("20060102", "20170701")
	if err != nil {
		t.Fatal(err)
	}
	c1 := &rf2.Concept{ID: 24700007, EffectiveTime: d, Active: true, ModuleID: 0, DefinitionStatusID: 900000000000073002}
	bolt.PutConcepts(c1)
	c2, err := bolt.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c1, c2) {
		t.Fatal("Concept not stored and retrieved correctly!")
	}
	c3, err := bolt.GetConcept(0)
	if c3 != nil && err != nil {
		t.Fatal("Failed to flag unfound concept")
	}
	os.Remove(dbFilename)
}
