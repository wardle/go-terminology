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
	"github.com/gogo/protobuf/io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/wardle/go-terminology/snomed"
)

// Export exports all descriptions in delimited protobuf format to the command line.
func (svc *Svc) Export() error {
	w := io.NewDelimitedWriter(os.Stdout)
	defer w.Close()

	in := getConcepts(svc)
	cpu := runtime.NumCPU()
	processors := make([]<-chan snomed.ExtendedDescription, cpu)
	for i := 0; i < cpu; i++ {
		processors[i] = convert(svc, in)
	}
	count := 0
	start := time.Now()
	for ed := range merge(processors...) {
		w.WriteMsg(&ed)
		count++
		if count%10000 == 0 {
			elapsed := time.Since(start)
			fmt.Fprintf(os.Stderr, "\rProcessed %d descriptions in %s. Mean time per description: %s...", count, elapsed, elapsed/time.Duration(count))
		}
	}
	fmt.Fprintf(os.Stderr, "\nProcessed total: %d descriptions in %s.\n", count, time.Since(start))
	return nil
}

// get all concepts
func getConcepts(svc *Svc) <-chan snomed.Concept {
	out := make(chan snomed.Concept)
	go func() {
		svc.Iterate(func(concept *snomed.Concept) error {
			out <- *concept
			return nil
		})
		close(out)
	}()
	return out
}

// convert takes a feed of concepts and turns them into an extended descriptions
func convert(svc *Svc, in <-chan snomed.Concept) <-chan snomed.ExtendedDescription {
	out := make(chan snomed.ExtendedDescription)
	go func() {
		for concept := range in {
			ed := snomed.ExtendedDescription{}
			err := initialiseExtendedFromConcept(svc, &ed, &concept)
			if err != nil {
				panic(err)
			}
			descs, err := svc.GetDescriptions(&concept)
			if err != nil {
				panic(err)
			}
			for _, d := range descs {
				ded := ed // make a copy
				err = initialiseExtendedFromDescription(svc, &ded, d)
				if err != nil {
					panic(err)
				}
				out <- ded
			}
		}
		close(out)
	}()
	return out
}

// merge multiple channels of work into a single result channel that
// can be processed serially. From https://blog.golang.org/pipelines
func merge(cs ...<-chan snomed.ExtendedDescription) <-chan snomed.ExtendedDescription {
	var wg sync.WaitGroup
	out := make(chan snomed.ExtendedDescription)
	output := func(c <-chan snomed.ExtendedDescription) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() { // close out once all goroutines are done
		wg.Wait()
		close(out)
	}()
	return out
}

func initialiseExtendedFromConcept(svc *Svc, ed *snomed.ExtendedDescription, c *snomed.Concept) error {
	ed.Concept = c
	preferred, err := svc.GetPreferredSynonym(c, BritishEnglish.LanguageReferenceSetIdentifier())
	if err != nil {
		// return nil, err // TODO: change API as not really an error or need to fallback to US English?
	}
	ed.PreferredDescription = preferred

	allParents, err := svc.GetAllParentIDs(c)
	if err != nil {
		return err
	}
	ed.RecursiveParentIds = allParents
	directParents, err := svc.GetParentIDsOfKind(c, snomed.IsAConceptID)
	if err != nil {
		return err
	}
	ed.DirectParentIds = directParents
	conceptRefsets, err := svc.GetReferenceSets(c.Id) // get reference sets for concept
	if err != nil {
		return err
	}
	ed.ConceptRefsets = conceptRefsets
	return nil
}

// TODO: pass language as a parameter rather than hard-coding British English
func initialiseExtendedFromDescription(svc *Svc, ed *snomed.ExtendedDescription, d *snomed.Description) error {
	ed.Description = d
	descRefsets, err := svc.GetReferenceSets(d.Id) // reference sets for description
	if err != nil {
		return err
	}
	ed.DescriptionRefsets = descRefsets
	return nil
}
