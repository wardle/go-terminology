package snomed

import (
	"testing"
)

func TestIdentifiers(t *testing.T) {
	testIdentifier(t, 24700007, true, false, false)  // multiple sclerosis concept
	testIdentifier(t, 123037004, true, false, false) // body structure concept
	testIdentifier(t, 724699017, false, true, false) // body structure description
	testIdentifier(t, 1399025, false, false, true)   // a relationship
	testIdentifier(t, 24700001, false, false, false) // invalid concept
}

func testIdentifier(t *testing.T, id Identifier, concept bool, description bool, relationship bool) {
	if concept || description || relationship {
		if id.IsValid() == false {
			t.Errorf("Identifier %d incorrectly identified as invalid.", id)
		}
		if id.IsConcept() != concept {
			t.Errorf("Identifier %d not correctly identified as a concept.", id)
		}
		if id.IsDescription() != description {
			t.Errorf("Identifier %d misidentified as a description", id)
		}
		if id.isRelationship() != relationship {
			t.Errorf("Identifier %d misidentified as a relationship", id)
		}
	} else {
		if id.IsValid() {
			t.Errorf("Identifier %d incorrectly identified as valid", id)
		}
	}
}
