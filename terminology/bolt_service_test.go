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
package terminology

import (
	"github.com/wardle/go-terminology/snomed"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	dbFilename = "bolt-tests.db"
)

func TestStore(t *testing.T) {
	bolt, err := newBoltService(dbFilename, false)
	if err != nil {
		t.Fatal(err)
	}
	d, err := time.Parse("20060102", "20170701")
	if err != nil {
		t.Fatal(err)
	}
	c1 := &snomed.Concept{ID: 24700007, EffectiveTime: d, Active: true, ModuleID: 0, DefinitionStatusID: 900000000000073002}
	bolt.Put([]*snomed.Concept{c1})
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
