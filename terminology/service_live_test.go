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
	"fmt"
	"os"
	"testing"

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
	svc, err := terminology.NewService(dbFilename, false)
	if err != nil {
		tb.Fatal(err)
	}
	return svc
}
func TestService(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	ms, err := svc.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	parents, err := svc.GetAllParents(ms)
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

	allChildrenIDs, err := svc.GetAllChildrenIDs(ms)
	if err != nil {
		t.Fatal(err)
	}
	allChildren, err := svc.GetConcepts(allChildrenIDs...)
	if err != nil {
		t.Fatal(err)
	}
	if len(allChildren) < 2 {
		t.Fatal("Did not correctly find many recursive children for MS")
	}
	for _, child := range allChildren {
		fsn, err := svc.GetFullySpecifiedName(child, terminology.BritishEnglish.LanguageReferenceSetIdentifier())
		if err != nil || fsn == nil {
			t.Fatalf("Missing FSN for concept %d : %v", child.Id, err)
		}
	}
}

func TestIterator(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	concepts, descriptions := 0, 0
	finished := fmt.Errorf("Finished")
	err := svc.Iterate(func(concept *snomed.Concept) error {
		concepts++
		descs, err := svc.GetDescriptions(concept)
		if err != nil {
			return err
		}
		descriptions += len(descs)
		if concepts == 5000 {
			return finished
		}
		return nil
	})
	if err != nil && err != finished {
		t.Fatal(err)
	}
	t.Logf("Iterated across %d concepts and %d descriptions", concepts, descriptions)
}

func BenchmarkGetConceptAndDescriptions(b *testing.B) {
	svc := setUp(b)
	defer svc.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms, err := svc.GetConcept(24700007)
		if err != nil {
			b.Fatal(err)
		}
		_, err = svc.GetDescriptions(ms)
		if err != nil {
			b.Fatal(err)
		}
		_, err = svc.GetPreferredSynonym(ms, terminology.BritishEnglish.LanguageReferenceSetIdentifier())
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIsA(b *testing.B) {
	svc := setUp(b)
	defer svc.Close()
	ms, err := svc.GetConcept(24700007)
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

func TestStatistics(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	_, err := svc.GetStatistics()
	if err != nil {
		t.Fatal(err)
	}

}

func TestLocalisation(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	appendicectomy, err := svc.GetConcept(80146002)
	if err != nil {
		t.Fatal(err)
	}
	d1, err := svc.GetPreferredSynonym(appendicectomy, terminology.BritishEnglish.LanguageReferenceSetIdentifier())
	if err != nil {
		t.Fatal(err)
	}
	d2, err := svc.GetPreferredSynonym(appendicectomy, terminology.AmericanEnglish.LanguageReferenceSetIdentifier())
	if err != nil {
		t.Fatal(err)
	}
	if d1.Term != "Appendicectomy" {
		t.Fatalf("%s is not the correct British term for Appendicectomy", d1.Term)
	}
	if d2.Term != "Appendectomy" {
		t.Fatalf("%s is not the correct British term for Appendicectomy", d2.Term)
	}
	fsn1 := svc.MustGetFullySpecifiedName(appendicectomy, terminology.BritishEnglish.LanguageReferenceSetIdentifier())
	fsn2 := svc.MustGetFullySpecifiedName(appendicectomy, terminology.AmericanEnglish.LanguageReferenceSetIdentifier())
	if fsn1.Term != fsn2.Term {
		t.Fatalf("fsn for appendicectomy appears to be different for British and American English: %s vs %s", fsn1.Term, fsn2.Term)
	}
}

func TestGenericisation(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	adem, err := svc.GetConcept(83942000) // acute disseminated encephalomyelitis
	if err != nil {
		t.Fatal(err)
	}
	path, err := svc.LongestPathToRoot(adem)
	if err != nil {
		t.Fatal(err)
	}
	refset, err := svc.GetReferenceSet(991411000000109) // emergency care reference set
	encephalitis, ok := svc.GenericiseTo(adem, refset)
	if !ok {
		t.Fatal("Could not map ADEM to the emergency care reference set")
	}
	if encephalitis.Id != 45170000 {
		t.Fatalf("Did not map ADEM to encephalitis but to %s", svc.MustGetPreferredSynonym(encephalitis, terminology.BritishEnglish.LanguageReferenceSetIdentifier()).Term)
	}
	debugPath(svc, path)
}

func debugPath(svc *terminology.Svc, path []*snomed.Concept) {
	for _, c := range path {
		fmt.Printf("%s(%d)--", svc.MustGetPreferredSynonym(c, terminology.BritishEnglish.LanguageReferenceSetIdentifier()).Term, c.Id)
	}
}
