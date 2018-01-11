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

package db

import (
	"fmt"
	"log"
	"os"

	"github.com/wardle/go-terminology/snomed"
)

// PerformImport performs import of SNOMED-CT structures from the root specified.
// This automatically clears the precomputations, if they exist, but does
// not run precomputations at the end as the user may run multiple individual imports
// from multiple SNOMED-CT distributions before finally running precomputations
// at the end of multiple imports.
func (sct *Snomed) PerformImport(root string) {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile) // for future use
	importer := snomed.NewImporter(logger)
	concepts, descriptions, relationships := 0, 0, 0
	importer.SetConceptHandler(func(c []*snomed.Concept) {
		concepts = concepts + len(c)
		err := sct.PutConcepts(c)
		if err != nil {
			logger.Printf("error importing concept : %v", err)
		}
	})
	importer.SetDescriptionHandler(func(d []*snomed.Description) {
		descriptions += len(d)
		err := sct.PutDescriptions(d)
		if err != nil {
			logger.Printf("error importing description : %v", err)
		}
	})
	importer.SetRelationshipHandler(func(r []*snomed.Relationship) {
		relationships += len(r)
		err := sct.PutRelationships(r)
		if err != nil {
			logger.Printf("error importing relationship : %v", err)
		}
	})
	sct.ClearPrecomputations()
	err := importer.ImportFiles(root)
	if err != nil {
		log.Fatalf("Could not import files: %v", err)
	}
	fmt.Printf("Imported %d concepts, %d descriptions and %d relationships\n", concepts, descriptions, relationships)
}

// ClearPrecomputations clears all precached precomputations
func (sct *Snomed) ClearPrecomputations() {
	// TODO(mw):implement
}

// PerformPrecomputations performs precomputations caching the results
func (sct *Snomed) PerformPrecomputations() {
	// TODO(mw):implement
}
