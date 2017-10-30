package mcqs

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"bytes"
	"reflect"
	"testing"
)

// test roundtripping from CSV and back
func TestCsvToFromStrings(t *testing.T) {
	concept1, err := snomed.NewConcept(1, "Wibble", 0, []int{1, 2, 3, 4, 5})
	if err != nil {
		t.Errorf("Failed to create concept: %s", err)
	}
	csv := conceptToCsv(concept1)
	concept2, err := conceptFromCsv(csv)
	if err != nil {
		t.Fatalf("Failed to roundtrip concept to []string and back: %s", err)
	}
	if reflect.DeepEqual(concept1, concept2) == false {
		t.Error("Failed to roundtrip concept to []string and back")
	}
}

// tests writing and reading in CSV format to an abstract I/O buffer
func TestWriteReadCsv(t *testing.T) {
	concept1, err := snomed.NewConcept(1, "Wibble", 0, []int{1, 2, 3, 4, 5})
	concepts := make(map[snomed.Identifier]*snomed.Concept)
	concepts[1] = concept1
	if err != nil {
		t.Errorf("Failed to create concept: %s", err)
	}
	var buffer bytes.Buffer
	writeToCsv(&buffer, concepts)
	concepts2, err := readFromCsv(&buffer)
	if err != nil {
		t.Errorf("Failed to read concepts from buffer: %s", err)
	}
	concept2 := concepts2[1]
	if reflect.DeepEqual(concept1, concept2) == false {
		t.Error("Failed to roundtrip concept to []string and back")
	}
}
