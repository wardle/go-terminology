// Copyright 2018 Mark Wardle / Eldrix Ltd
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

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

func testIdentifier(t *testing.T, identifier int64, concept bool, description bool, relationship bool) {
	id := Identifier(identifier)
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
