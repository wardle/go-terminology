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

package terminology_test

import (
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"os"
	"testing"
	"time"

	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
)

const (
	fakeDbFilename = "bolt-tests.db" // transient fake
)

func TestStore(t *testing.T) {
	svc, err := terminology.NewService(fakeDbFilename, false)
	if err != nil {
		t.Fatal(err)
	}
	defer svc.Close()
	defer os.RemoveAll(fakeDbFilename)

	date, err := time.Parse("20060102", "20170701")
	if err != nil {
		t.Fatal(err)
	}
	d, err := ptypes.TimestampProto(date)
	if err != nil {
		t.Fatal(err)
	}
	c1 := &snomed.Concept{Id: 24700007, EffectiveTime: d, Active: true, ModuleId: 0, DefinitionStatusId: 900000000000073002}
	c2 := &snomed.Concept{Id: 6118003, EffectiveTime: d, Active: true, ModuleId: 0, DefinitionStatusId: 900000000000073002}
	c3 := &snomed.Concept{Id: snomed.IsA, EffectiveTime: d, Active: true}
	d1 := &snomed.Description{Id: 41398015, ConceptId: 24700007, EffectiveTime: d, Active: true, ModuleId: 0, Term: "Multiple sclerosis"}
	d2 := &snomed.Description{Id: 1223979019, ConceptId: 24700007, EffectiveTime: d, Active: true, ModuleId: 0, Term: "Disseminated sclerosis"}
	d3 := &snomed.Description{Id: 11161017, ConceptId: 6118003, EffectiveTime: d, Active: true, ModuleId: 0, Term: "Demyelinating disease"}
	r1 := &snomed.Relationship{Id: 1, Active: true, EffectiveTime: d, SourceId: c1.Id, DestinationId: c2.Id, TypeId: snomed.IsA}
	svc.Put([]*snomed.Concept{c1, c2, c3})
	svc.Put([]*snomed.Description{d1, d2, d3})
	svc.Put([]*snomed.Relationship{r1})
	c, err := svc.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(c1, c) {
		t.Fatal("Concept not stored and retrieved correctly!")
	}
	_, err = svc.GetConcept(0)
	if err == nil {
		t.Fatal("Failed to flag unfound concept")
	}
	descriptions, err := svc.GetDescriptions(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(descriptions) != 2 {
		t.Fatal("Returned wrong number of descriptions")
	}

	for _, d := range descriptions {
		if d.Id != d1.Id && d.Id != d2.Id {
			t.Fatal("did not get correct descriptions back for concept")
		}
	}
	childRels, err := svc.GetChildRelationships(c1)
	if err != nil {
		t.Fatal(err)
	}
	if len(childRels) != 0 {
		t.Fatal("Multiple sclerosis given child concepts!")
	}
	parentRels, err := svc.GetParentRelationships(c1)
	if err != nil {
		t.Fatal(err)
	}
	if len(parentRels) != 1 || parentRels[0].DestinationId != c2.Id {
		t.Fatal("Demyelinating disease not a parent of multiple sclerosis")
	}
	childRels, err = svc.GetChildRelationships(c2)
	if len(childRels) != 1 || childRels[0].SourceId != c1.Id {
		t.Fatal("Multiple sclerosis not a child of demyelinating disease of the CNS")
	}
	parents, err := svc.GetParents(c1)
	if len(parents) != 1 || parents[0].Id != c2.Id {
		t.Fatal("Demyelinating disease not a parent of multiple sclerosis")
	}
	children, err := svc.GetChildren(c1)
	if len(children) != 0 {
		t.Fatal("Multiple sclerosis given child concepts!")
	}

}
