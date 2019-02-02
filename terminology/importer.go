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
	"github.com/wardle/go-terminology/snomed"
	"log"
	"os"
	"strings"
	"time"
)

type timedTask struct {
	time  time.Duration
	count int
}

func (ts *timedTask) String() string {
	return fmt.Sprintf("%s", time.Duration(int(ts.time)/ts.count))
}

type timedTasks struct {
	tasks map[string]*timedTask
}

func newTimedTasks() *timedTasks {
	tt := new(timedTasks)
	tt.tasks = make(map[string]*timedTask)
	return tt
}

func (t *timedTasks) String() string {
	var b strings.Builder
	for name, ts := range t.tasks {
		b.WriteString(name)
		b.WriteString(": ")
		b.WriteString(ts.String())
		b.WriteString(". ")
	}
	return strings.TrimSpace(b.String())
}

func (t *timedTasks) recordTime(name string, d time.Duration, n int) {
	var s *timedTask
	var ok bool
	if s, ok = t.tasks[name]; !ok {
		s = new(timedTask)
		t.tasks[name] = s
	}
	s.time += d
	s.count += n
}

// PerformImport performs import of SNOMED-CT structures from the root specified.
// This automatically clears the precomputations, if they exist, but does
// not run precomputations at the end as the user may run multiple individual imports
// from multiple SNOMED-CT distributions before finally running precomputations
// at the end of multiple imports.
func (svc *Svc) PerformImport(root string, verbose bool) {
	logger := log.New(os.Stdout, "import: ", log.Lshortfile)
	concepts, descriptions, relationships, refsets := 0, 0, 0, 0
	var err error
	tt := newTimedTasks()
	start := time.Now()
	tick := start
	importer := snomed.NewImporter(logger, func(o interface{}) {
		batchStart := time.Now()
		err = svc.Put(o)
		duration := time.Since(batchStart)
		if err != nil {
			logger.Printf("error importing : %v", err)
		} else {
			switch o.(type) {
			case []*snomed.Concept:
				concepts += len(o.([]*snomed.Concept))
				tt.recordTime("concepts", duration, len(o.([]*snomed.Concept)))
			case []*snomed.Description:
				descriptions += len(o.([]*snomed.Description))
				tt.recordTime("descriptions", duration, len(o.([]*snomed.Description)))
			case []*snomed.Relationship:
				relationships += len(o.([]*snomed.Relationship))
				tt.recordTime("relationships", duration, len(o.([]*snomed.Relationship)))
			case []*snomed.ReferenceSetItem:
				refsets += len(o.([]*snomed.ReferenceSetItem))
				tt.recordTime("reference set items", duration, len(o.([]*snomed.ReferenceSetItem)))
			}
		}
		if verbose {
			tickDuration := time.Since(tick).Seconds()
			if tickDuration > 10 {
				tick = time.Now()
				logger.Printf("%s : %v - %d concepts, %d descriptions, %d relationships and %d refset items...\n",
					time.Since(start), tt, concepts, descriptions, relationships, refsets)
			}
		}
	})
	svc.ClearPrecomputations()
	err = importer.ImportFiles(root)
	if err != nil {
		log.Fatalf("Could not import files: %v", err)
	}
	fmt.Printf("Complete; imported %d concepts, %d descriptions, %d relationships and %d refset items\n", concepts, descriptions, relationships, refsets)
	fmt.Printf("Duration : %v\n", tt)
}
