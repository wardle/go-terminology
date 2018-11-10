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
	"golang.org/x/text/language"
	"os"
	"time"

	"github.com/wardle/go-terminology/snomed"
)

// Export exports all descriptions in delimited protobuf format to the command line.
func (svc *Svc) Export() error {
	w := io.NewDelimitedWriter(os.Stdout)
	defer w.Close()

	count := 0
	start := time.Now()
	err := svc.Iterate(func(concept *snomed.Concept) error {
		ed := snomed.ExtendedDescription{}
		err := initialiseExtendedFromConcept(svc, &ed, concept)
		if err != nil {
			panic(err)
		}
		descs, err := svc.Descriptions(concept.Id)
		if err != nil {
			panic(err)
		}
		for _, d := range descs {
			ded := ed // make a copy
			err = initialiseExtendedFromDescription(svc, &ded, d)
			if err != nil {
				panic(err)
			}
			w.WriteMsg(&ded)
			count++
			if count%10000 == 0 {
				elapsed := time.Since(start)
				fmt.Fprintf(os.Stderr, "\rProcessed %d descriptions in %s. Mean time per description: %s...", count, elapsed, elapsed/time.Duration(count))
			}
		}
		return nil
	})
	fmt.Fprintf(os.Stderr, "\nProcessed total: %d descriptions in %s.\n", count, time.Since(start))
	return err
}

func initialiseExtendedFromConcept(svc *Svc, ed *snomed.ExtendedDescription, c *snomed.Concept) error {
	ed.Concept = c
	tags, _, _ := language.ParseAcceptLanguage("en-GB")
	ed.PreferredDescription = svc.MustGetPreferredSynonym(c, tags)

	allParents, err := svc.AllParentIDs(c)
	if err != nil {
		return err
	}
	ed.RecursiveParentIds = allParents
	directParents, err := svc.ParentIDsOfKind(c, snomed.IsA)
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

// TODO: pass language as a parameter rather than hard-coding British English
func initialiseExtendedFromDescription(svc *Svc, ed *snomed.ExtendedDescription, d *snomed.Description) error {
	ed.Description = d
	descRefsets, err := svc.ComponentReferenceSets(d.Id) // reference sets for description
	if err != nil {
		return err
	}
	ed.DescriptionRefsets = descRefsets
	return nil
}
