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
	"github.com/wardle/go-terminology/snomed"
	"fmt"
	"log"
	"os"
)

// PerformImport performs import of SNOMED-CT structures from the root specified.
// This automatically clears the precomputations, if they exist, but does
// not run precomputations at the end as the user may run multiple individual imports
// from multiple SNOMED-CT distributions before finally running precomputations
// at the end of multiple imports.
func PerformImport(bolt *BoltService, root string) {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile) // for future use
	importer := snomed.NewImporter(logger)
	concepts, descriptions, relationships := 0, 0, 0
	importer.SetConceptHandler(func(c []*snomed.Concept) {
		concepts = concepts + len(c)
		bolt.PutConcepts(c...)
	})
	importer.SetDescriptionHandler(func(d []*snomed.Description) {
		descriptions += len(d)
		bolt.PutDescriptions(d...)
	})
	importer.SetRelationshipHandler(func(r []*snomed.Relationship) {
		relationships += len(r)
		bolt.PutRelationships(r...)
	})
	ClearPrecomputations(bolt)
	err := importer.ImportFiles(root)
	if err != nil {
		log.Fatalf("Could not import files: %v", err)
	}
	fmt.Printf("Imported %d concepts, %d descriptions and %d relationships\n", concepts, descriptions, relationships)
}

// ClearPrecomputations clears all precached precomputations
func ClearPrecomputations(bolt *BoltService) {
	// TODO(mw):implement
}

// PerformPrecomputations performs precomputations caching the results
func PerformPrecomputations(bolt *BoltService) {
	// TODO(mw):implement
}
