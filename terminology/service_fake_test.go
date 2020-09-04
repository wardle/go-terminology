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
	"context"
	"os"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
)

const (
	fakeDbFilename = "fake-tests.db" // transient fake
)

func TestStore(t *testing.T) {
	svc, err := terminology.NewService(fakeDbFilename, false)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fakeDbFilename)
	defer svc.Close()

	date, err := time.Parse("20060102", "20170701")
	if err != nil {
		t.Fatal(err)
	}
	d := timestamppb.New(date)
	c1 := &snomed.Concept{Id: 24700007, EffectiveTime: d, Active: true, ModuleId: 0, DefinitionStatusId: 900000000000073002}
	c2 := &snomed.Concept{Id: 6118003, EffectiveTime: d, Active: true, ModuleId: 0, DefinitionStatusId: 900000000000073002}
	c3 := &snomed.Concept{Id: snomed.IsA, EffectiveTime: d, Active: true}
	d1 := &snomed.Description{Id: 41398015, ConceptId: 24700007, EffectiveTime: d, Active: true, ModuleId: 0, Term: "Multiple sclerosis", TypeId: 900000000000013009}
	d2 := &snomed.Description{Id: 1223979019, ConceptId: 24700007, EffectiveTime: d, Active: true, ModuleId: 0, Term: "Disseminated sclerosis", TypeId: 900000000000013009}
	d3 := &snomed.Description{Id: 11161017, ConceptId: 6118003, EffectiveTime: d, Active: true, ModuleId: 0, Term: "Demyelinating disease", TypeId: 900000000000013009}
	d4 := &snomed.Description{Id: 181114011, ConceptId: 116680003, EffectiveTime: d, Active: true, ModuleId: 0, Term: "Is a", TypeId: 900000000000013009}
	r1 := &snomed.Relationship{Id: 1, Active: true, EffectiveTime: d, SourceId: c1.Id, DestinationId: c2.Id, TypeId: snomed.IsA}
	ctx := context.Background()
	if err := svc.Put(ctx, []*snomed.Concept{c1, c2, c3}); err != nil {
		t.Fatal(err)
	}
	if err := svc.Put(ctx, []*snomed.Description{d1, d2, d3, d4}); err != nil {
		t.Fatal(err)
	}
	if err := svc.Put(ctx, []*snomed.Relationship{r1}); err != nil {
		t.Fatal(err)
	}
	err = svc.PerformPrecomputations(ctx, 500, false)
	if err != nil {
		t.Fatal(err)
	}
	c, err := svc.Concept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(c1, c) {
		t.Fatalf("Concept not stored and retrieved correctly! Tried to store: %v, got back %v", c1, c)
	}
	_, err = svc.Concept(0)
	if err == nil {
		t.Fatal("Failed to flag unfound concept")
	}
	descriptions, err := svc.Descriptions(c.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(descriptions) != 2 {
		t.Fatalf("Returned wrong number of descriptions. expected: %d. got: %d", 2, len(descriptions))
	}

	for _, d := range descriptions {
		if d.Id != d1.Id && d.Id != d2.Id {
			t.Fatal("did not get correct descriptions back for concept")
		}
	}
	if descriptions[0].Id == descriptions[1].Id {
		t.Fatal("got back two descriptions with the same identifier")
	}
	childRels, err := svc.ChildRelationships(c1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(childRels) != 0 {
		t.Fatal("Multiple sclerosis given child concepts!")
	}
	parentRels, err := svc.ParentRelationships(c1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(parentRels) != 1 || parentRels[0].DestinationId != c2.Id {
		t.Fatal("Demyelinating disease not a parent of multiple sclerosis")
	}
	childRels, err = svc.ChildRelationships(c2.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(childRels) != 1 || childRels[0].SourceId != c1.Id {
		t.Fatal("Multiple sclerosis not a child of demyelinating disease of the CNS")
	}
	parents, err := svc.Parents(c1.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(parents) != 1 || parents[0] != c2.Id {
		t.Fatal("Demyelinating disease not a parent of multiple sclerosis")
	}
	children, err := svc.Children(c1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(children) != 0 {
		t.Fatal("Multiple sclerosis given child concepts!")
	}

}
