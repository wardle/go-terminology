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
func (svc *Svc) PerformImport(root string) {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	concepts, descriptions, relationships, refsets := 0, 0, 0, 0
	var err error
	importer := snomed.NewImporter(logger, func(o interface{}) {
		err = svc.Put(o)
		if err != nil {
			logger.Printf("error importing : %v", err)
		} else {
			switch o.(type) {
			case []*snomed.Concept:
				concepts += len(o.([]*snomed.Concept))
			case []*snomed.Description:
				descriptions += len(o.([]*snomed.Description))
			case []*snomed.Relationship:
				relationships += len(o.([]*snomed.Relationship))
			case []*snomed.ReferenceSetItem:
				refsets += len(o.([]*snomed.ReferenceSetItem))
			}
		}
	})
	svc.ClearPrecomputations()
	err = importer.ImportFiles(root)
	if err != nil {
		log.Fatalf("Could not import files: %v", err)
	}
	fmt.Printf("Imported %d concepts, %d descriptions, %d relationships and %d refsets\n", concepts, descriptions, relationships, refsets)
}

// ClearPrecomputations clears all precached precomputations
func (svc *Svc) ClearPrecomputations() {
	// TODO(mw):implement
}

// PerformPrecomputations performs precomputations caching the results
func (svc *Svc) PerformPrecomputations() {
	// TODO(mw):implement
}
