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
	"golang.org/x/text/language"
	"os"
	"testing"

	"github.com/wardle/go-terminology/snomed"
)

const (
	dbFilename = "../snomed.db" // real, live database
)

func setUp(tb testing.TB) *Svc {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) { // skip these tests if no working live snomed db
		tb.Skipf("Skipping tests against a live database. To run, create a database named %s", dbFilename)
	}
	svc, err := NewService(dbFilename, true)
	if err != nil {
		tb.Fatal(err)
	}
	return svc
}
func TestExport(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	ms, err := svc.GetConcept(24700007)
	if err != nil {
		t.Fatal(err)
	}
	tags, _, _ := language.ParseAcceptLanguage("en-GB")
	d := svc.MustGetPreferredSynonym(ms, tags)
	ed := snomed.ExtendedDescription{}
	initialiseExtendedFromConcept(svc, &ed, ms)
	initialiseExtendedFromDescription(svc, &ed, d)

	if ed.PreferredDescription.GetTerm() != "Multiple sclerosis" {
		t.Error("Failed to export preferred term for multiple sclerosis")
	}
}
