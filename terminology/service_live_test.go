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
	"math/rand"
	"os"
	"testing"
	"time"

	"golang.org/x/text/language"

	"github.com/wardle/go-terminology/dmd"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
)

const (
	dbFilename = "../snomed.db" // real, live database
)

func setUp(tb testing.TB) *terminology.Svc {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) { // skip these tests if no working live snomed db
		tb.Skipf("Skipping tests against a live database. To run, create a database named %s", dbFilename)
	}
	svc, err := terminology.NewService(dbFilename, true)
	if err != nil {
		tb.Fatal(err)
	}
	return svc
}
func TestService(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	ms, err := svc.Concept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	parents, err := svc.AllParents(ms.Id)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range parents {
		if p.Id == 6118003 {
			found = true
		}
	}
	if !found {
		t.Fatal("Multiple sclerosis not correctly identified as a type of demyelinating disease")
	}
	if !svc.IsA(ms, 6118003) {
		t.Fatal("Multiple sclerosis not correctly identified as a type of demyelinating disease")
	}

	allChildrenIDs, err := svc.AllChildrenIDs(context.Background(), ms.Id, 500)
	if err != nil {
		t.Fatal(err)
	}
	allChildren, err := svc.Concepts(allChildrenIDs...)
	if err != nil {
		t.Fatal(err)
	}
	if len(allChildren) < 2 {
		t.Fatal("Did not correctly find many recursive children for MS")
	}
	for _, child := range allChildren {
		fsn, err := svc.FullySpecifiedName(child, []language.Tag{terminology.BritishEnglish.Tag()})
		if err != nil || fsn == nil {
			t.Fatalf("Missing FSN for concept %d : %v", child.Id, err)
		}
	}
}

func TestReferenceSets(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	ms, err := svc.Concept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	if ms.Id != 24700007 {
		t.Fatal("failed to fetch concept: multiple sclerosis")
	}
	items, err := svc.ComponentFromReferenceSet(900000000000497000, 24700007) // is multiple sclerosis in the Read crossmap?
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, item := range items {
		if sm := item.GetSimpleMap(); sm != nil {
			if sm.GetMapTarget() != "F20.." {
				t.Fatal("Multiple sclerosis should be F20..")
			}
			found = true
		}
	}
	if !found {
		t.Fatal("Multiple sclerosis not found in Read crossmap!")
	}
	exists, err := svc.IsInReferenceSet(24700007, 900000000000497000)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("Multiple sclerosis not found in Read crossmap!")
	}

	rsi, err := svc.ReferenceSetItem("d55ce305-3dcc-5723-8814-cd26486c37f7") // this is from emergency care refset - MS
	if err != nil {
		t.Fatal(err)
	}
	if rsi.GetSimple() == nil {
		t.Fatalf("expected: simple reference set. found: %v", rsi)
	}

	items2, err := svc.ComponentFromReferenceSet(991411000000109, 24700007) // emergency care diagnosis refset
	if err != nil && err != terminology.ErrNotFound {
		t.Fatal(err)
	}
	if len(items2) == 0 {
		t.Fatal("MS not in emergency care diagnosis reference set")
	}

}

func TestDrugs(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	amlodipine, err := svc.Concept(108537001)
	if err != nil {
		t.Fatal(err)
	}
	fsn, err := svc.FullySpecifiedName(amlodipine, []language.Tag{terminology.BritishEnglish.Tag()})
	if err != nil {
		t.Fatal(err)
	}
	if fsn == nil {
		descs, err := svc.Descriptions(amlodipine.Id)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("missing FSN in descriptions for %v:\n%v", amlodipine, descs)
		t.Fatalf("missing FSN for drug %v", amlodipine)
	}
	if fsn.ConceptId != 108537001 {
		t.Errorf("FSN returned for incorrect concept. expected: 108537001. got: %d", fsn.ConceptId)
	}
	if fsn.Term != "Amlodipine (substance)" {
		t.Errorf("FSN for amlodipine incorrect.")
	}
}

func TestIterator(t *testing.T) {
	duration := 200 * time.Millisecond
	if testing.Short() {
		duration = duration / 10
	}
	svc := setUp(t)
	defer svc.Close()
	count := 0
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	conceptc := svc.IterateConcepts(ctx)
	for c := range conceptc {
		if c.Err != nil {
			t.Fatal(c.Err)
		}
		count++
	}
	if count == 0 {
		t.Errorf("Did not iterate appropriately")
	}
	t.Logf("Iterated across %d concepts", count)
}

func BenchmarkGetConceptAndDescriptions(b *testing.B) {
	svc := setUp(b)
	defer svc.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms, err := svc.Concept(24700007)
		if err != nil {
			b.Fatal(err)
		}
		_, err = svc.Descriptions(ms.Id)
		if err != nil {
			b.Fatal(err)
		}
		d, err := svc.PreferredSynonym(ms.Id, []language.Tag{terminology.BritishEnglish.Tag()})
		if err != nil {
			b.Fatal(err)
		}
		if d == nil {
			b.Fatalf("missing synonym for %v", ms)
		}
	}
}

func BenchmarkExtendedConcept(b *testing.B) {
	tags := []language.Tag{terminology.BritishEnglish.Tag()}
	svc := setUp(b)
	defer svc.Close()
	vtms, err := svc.ReferenceSetComponents(dmd.VtmReferenceSet)
	if err != nil {
		b.Fatal(err)
	}
	nvtms := len(vtms)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := rand.Intn(nvtms)
		for vtmID := range vtms {
			if r == 0 {
				vtm, err := svc.ExtendedConcept(vtmID, tags)
				if err != nil {
					b.Fatal(err)
				}
				if vtm.Concept.Id != vtmID {
					b.Fatalf("extended concept identifier does not match requested. expected: %d, got: %v", vtmID, vtm)
				}
			}
			r--
		}
	}

}

func BenchmarkIsA(b *testing.B) {
	svc := setUp(b)
	defer svc.Close()
	ms, err := svc.Concept(24700007)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	var (
		isDemyelinating  bool
		isPharmaceutical bool
	)
	for i := 0; i < b.N; i++ {
		isDemyelinating = svc.IsA(ms, 6118003)
		isPharmaceutical = svc.IsA(ms, 373873005)
	}
	if isDemyelinating == false || isPharmaceutical == true {
		b.Fatal("MS misclassified using IS-A hierarchy")
	}
}

func BenchmarkAllChildren(b *testing.B) {
	svc := setUp(b)
	defer svc.Close()
	b.ResetTimer()
	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		allChildren, err := svc.AllChildrenIDs(ctx, 64572001, 1000000) // degenerative disease  > 75,000 children should exist in ontology
		if err != nil {
			b.Fatal(err)
		}
		if len(allChildren) == 0 {
			b.Fatal("no children found")
		}
	}
}

func BenchmarkAllChildrenShort(b *testing.B) {
	svc := setUp(b)
	defer svc.Close()
	b.ResetTimer()
	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		allChildren, err := svc.AllChildrenIDs(ctx, 24700007, 1000000) // ms - only 15 or so children
		if err != nil {
			b.Fatal(err)
		}
		if len(allChildren) == 0 {
			b.Fatal("no children found")
		}
	}
}
func TestLocalisation(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	appendicectomy, err := svc.Concept(80146002)
	if err != nil {
		t.Fatal(err)
	}
	d1, err := svc.PreferredSynonym(appendicectomy.Id, []language.Tag{terminology.BritishEnglish.Tag()})
	if err != nil {
		t.Fatal(err)
	}
	if d1 == nil {
		t.Fatalf("missing preferred synonym for %v in british english", appendicectomy)
	}
	d2, err := svc.PreferredSynonym(appendicectomy.Id, []language.Tag{terminology.AmericanEnglish.Tag()})
	if err != nil {
		t.Fatal(err)
	}
	if d2 == nil {
		t.Fatalf("missing preferred synonym for %v in american english", appendicectomy)
	}
	if d1.Term != "Appendicectomy" {
		t.Fatalf("%s is not the correct British term for Appendicectomy", d1.Term)
	}
	if d2.Term != "Appendectomy" {
		t.Fatalf("%s is not the correct British term for Appendicectomy", d2.Term)
	}
	fsn1 := svc.MustGetFullySpecifiedName(appendicectomy, []language.Tag{terminology.BritishEnglish.Tag()})
	fsn2 := svc.MustGetFullySpecifiedName(appendicectomy, []language.Tag{terminology.AmericanEnglish.Tag()})
	if fsn1.Term != fsn2.Term {
		t.Fatalf("fsn for appendicectomy appears to be different for British and American English: %s vs %s", fsn1.Term, fsn2.Term)
	}
	// request for a language not installed
	d3, err := svc.PreferredSynonym(appendicectomy.Id, []language.Tag{language.Swahili})
	if err != nil {
		t.Fatal(err)
	}
	if d3 == nil {
		t.Fatal("did not appropriately fallback for uninstalled language request")
	}
	if d3.Id != d2.Id {
		t.Fatalf("did not fallback to American English for uninstalled language. got: %v", d3)
	}
}

func TestGenericisation(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	adem, err := svc.Concept(83942000) // acute disseminated encephalomyelitis
	if err != nil {
		t.Fatal(err)
	}
	refsetItems, err := svc.ReferenceSetComponents(991411000000109) // emergency care reference set
	encephalitis, err := svc.GenericiseToBest(adem.Id, refsetItems)
	if err != nil {
		t.Fatalf("Could not map ADEM to the emergency care reference set: %s", err)
	}
	if encephalitis != 45170000 {
		t.Fatalf("Did not map ADEM to encephalitis but to %s", svc.MustGetPreferredSynonym(encephalitis, []language.Tag{terminology.BritishEnglish.Tag()}).Term)
	}
}
func TestGenericisation2(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	firstMI, err := svc.Concept(394710008) // some weird old specific type of MI - "firstMI"
	if err != nil {
		t.Fatal(err)
	}
	refsetItems, err := svc.ReferenceSetComponents(991411000000109) // emergency care reference set
	mi, err := svc.GenericiseToBest(firstMI.Id, refsetItems)
	if err != nil {
		t.Fatalf("Could not map 'First MI' to the emergency care reference set: %s", err)
	}
	if mi != 22298006 {
		t.Fatalf("Did not map 'First MI' to 'MI' but to %s", svc.MustGetPreferredSynonym(mi, []language.Tag{terminology.BritishEnglish.Tag()}).Term)
	}
}
func TestRefinements(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	tags := []language.Tag{terminology.BritishEnglish.Tag()}
	response, err := svc.Refinements(60404007, 20, tags) // cerebral abscess
	if err != nil {
		t.Error(err)
	}
	refinements := response.GetRefinements()
	if len(refinements) != 2 {
		t.Errorf("expected two refinements, got: %d:\n%v", len(refinements), refinements)
	}
	if refinements[0].GetRootValue().GetConceptId() != 116676008 && refinements[1].GetRootValue().GetConceptId() != 83678007 {
		t.Errorf("did not correctly identify that cerebral abscess can be refined by abscess morphology and finding site within cerebral structure. got:%v", refinements)
	}
}

func TestSearch(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	tags := []language.Tag{terminology.BritishEnglish.Tag()}
	request := &snomed.SearchRequest{
		S:     "amlodipine",
		Fuzzy: snomed.SearchRequest_ALWAYS_FUZZY,
		IsA:   []int64{10363601000001109},
	}
	response, err := svc.Search(request, tags)
	if err != nil {
		t.Error(err)
	}
	if len(response.Items) == 0 {
		t.Errorf("search for amlodipine:no results")
	}
	if response.Items[0].ConceptId != 108537001 {
		t.Errorf("Did not return amlodipine 108537001 as first result. got: %v", response.Items[0])
	}

}
