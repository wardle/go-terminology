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
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/gogo/protobuf/io"
	"golang.org/x/text/language"

	"github.com/wardle/go-terminology/snomed"
)

// Export exports all descriptions in delimited protobuf format to the command line.
func (svc *Svc) Export(lang string) error {
	w := io.NewDelimitedWriter(os.Stdout)
	defer w.Close()
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		return err
	}
	count := 0
	start := time.Now()
	eds := svc.iterateExtendedDescriptions(context.Background(), tags)
	for ed := range eds {
		w.WriteMsg(ed)
		count++
		if count%10000 == 0 {
			elapsed := time.Since(start)
			fmt.Fprintf(os.Stderr, "\rProcessed %d descriptions in %s. Mean time per description: %s...", count, elapsed, elapsed/time.Duration(count))
		}
	}
	fmt.Fprintf(os.Stderr, "\nProcessed total: %d descriptions in %s.\n", count, time.Since(start))
	return err
}

func (svc *Svc) iterateExtendedDescriptions(ctx context.Context, tags []language.Tag) <-chan *snomed.ExtendedDescription {
	conceptc := svc.IterateConcepts(ctx)
	resultc := make(chan *snomed.ExtendedDescription)
	go func() {
		defer close(resultc)
		var wg sync.WaitGroup
		for i := 0; i < runtime.NumCPU(); i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					case concept := <-conceptc:
						if concept == nil {
							return
						}
						svc.makeExtendedDescriptions(ctx, concept, tags, resultc)
					}
				}
			}()
		}
		wg.Wait() // wait, and then close the channel
	}()
	return resultc
}

func (svc *Svc) makeExtendedDescriptions(ctx context.Context, concept *snomed.Concept, tags []language.Tag, resultc chan<- *snomed.ExtendedDescription) {
	ed := snomed.ExtendedDescription{}
	err := initialiseExtendedFromConcept(svc, &ed, concept, tags)
	if err != nil {
		panic(err)
	}
	descs, err := svc.Descriptions(concept.Id)
	if err != nil {
		panic(err)
	}
	for _, d := range descs {
		ded := ed // make a copy
		if err = initialiseExtendedFromDescription(svc, &ded, d); err != nil {
			panic(err)
		}
		select {
		case <-ctx.Done():
			return
		case resultc <- &ded:
		}
	}
}

func initialiseExtendedFromConcept(svc *Svc, ed *snomed.ExtendedDescription, c *snomed.Concept, tags []language.Tag) error {
	ed.Concept = c
	ed.PreferredDescription = svc.MustGetPreferredSynonym(c.Id, tags)
	allParents, err := svc.AllParentIDs(c.Id)
	if err != nil {
		return err
	}
	ed.AllParentIds = allParents
	directParents, err := svc.ParentIDsOfKind(c.Id, snomed.IsA)
	if err != nil {
		return err
	}
	ed.DirectParentIds = directParents
	conceptRefsets, err := svc.ComponentReferenceSets(c.Id) // get reference sets for concept
	if err != nil {
		return err
	}
	ed.ConceptRefsets = conceptRefsets
	return nil
}

func initialiseExtendedFromDescription(svc *Svc, ed *snomed.ExtendedDescription, d *snomed.Description) error {
	ed.Description = d
	descRefsets, err := svc.ComponentReferenceSets(d.Id) // reference sets for description
	if err != nil {
		return err
	}
	ed.DescriptionRefsets = descRefsets
	return nil
}
