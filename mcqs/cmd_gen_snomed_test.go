package mcqs

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"reflect"
	"strings"
	"testing"
)

func TestListAtoi(t *testing.T) {
	testAtoi(t, "123,456,789,123456789", []int{123, 456, 789, 123456789}, true)
	testAtoi(t, "123,456,789, 123456789", []int{123, 456, 789, 123456789}, true)
	testAtoi(t, "aaa,vbb,cc,123", []int{123}, false)
}

// test conversion of a comma-delimited string of integers to a slice, optionally testing roundtripping
func testAtoi(t *testing.T, input string, expected []int, roundtrip bool) {
	r := listAtoi(input)
	if reflect.DeepEqual(r, expected) == false {
		t.Errorf("Failed to parse: %s. Parsed to %v", input, r)
	}
	if roundtrip {
		v := listItoA(expected)
		if reflect.DeepEqual(v, strings.Replace(input, " ", "", -1)) == false {
			t.Errorf("Failed to parse: %v. Parsed to %s", expected, v)
		}
	}
}

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

func TestWriteReadCsv(t *testing.T) {
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
