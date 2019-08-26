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
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wardle/go-terminology/snomed"
)

// Importer manages the import of SNOMED CT components from the filesystem
type Importer struct {
	storer Storer
	snomed.ImportChannels
	batchSize                                          int // size of each batch.
	threads                                            int // number of threads importing a type of component
	verbose                                            bool
	nconcepts, ndescriptions, nrelationships, nrefsets int32
}

// Storer defines the behaviour of a service that can accept a batch of SNOMED components for storage.
type Storer interface {
	Put(context context.Context, components interface{}) error
}

// NoopStorer is a no-op storer, which does nothing with the results. Useful for profiling.
type NoopStorer struct{}

// Put a batch of SNOMED components, which is a no-op.
func (ns *NoopStorer) Put(context context.Context, components interface{}) error {
	return nil
}

// NewImporter creates a new Importer
func NewImporter(storer Storer, batchSize int, threads int, verbose bool) *Importer {
	if batchSize == 0 {
		batchSize = 500
	}
	if threads == 0 {
		threads = runtime.NumCPU()
	}
	importer := &Importer{
		storer:    storer,
		verbose:   verbose,
		batchSize: batchSize,
		threads:   threads,
	}
	return importer
}

// Import starts the import process, returning errors when done
func (im *Importer) Import(ctx context.Context, root string) {
	start := time.Now()
	im.ImportChannels = *snomed.FastImport(ctx, root, im.batchSize)
	var conceptsWg, descriptionsWg, relationshipsWg, refsetsWg sync.WaitGroup
	done := make(chan struct{})
	if im.verbose {
		go func() {
			for {
				for _, r := range `-\|/` {
					select {
					case <-ctx.Done():
						return
					case <-done:
						return
					default:
						im.progress(start, string(r)+" importing")
						time.Sleep(1 * time.Second)
					}
				}
			}
		}()
	}
	for i := 0; i < im.threads; i++ {
		conceptsWg.Add(1)
		go func() {
			defer conceptsWg.Done()
			im.importConcepts(ctx, im.Concepts)
		}()
	}
	for i := 0; i < im.threads; i++ {
		descriptionsWg.Add(1)
		go func() {
			defer descriptionsWg.Done()
			im.importDescriptions(ctx, im.Descriptions)
		}()
	}
	for i := 0; i < im.threads; i++ {
		relationshipsWg.Add(1)
		go func() {
			defer relationshipsWg.Done()
			im.importRelationships(ctx, im.Relationships)
		}()
	}
	for i := 0; i < im.threads; i++ {
		refsetsWg.Add(1)
		go func() {
			defer refsetsWg.Done()
			im.importRefsets(ctx, im.Refsets)
		}()
	}
	conceptsWg.Wait()
	descriptionsWg.Wait()
	relationshipsWg.Wait()
	refsetsWg.Wait()
	close(done)
	im.progress(start, "Import complete. Processed: ")
}

func (im *Importer) progress(start time.Time, prefix string) {
	fmt.Printf("\r%s: %s: %d concepts, %d descriptions, %d relationships and %d refset items...",
		prefix,
		time.Since(start),
		atomic.LoadInt32(&im.nconcepts),
		atomic.LoadInt32(&im.ndescriptions),
		atomic.LoadInt32(&im.nrelationships),
		atomic.LoadInt32(&im.nrefsets))
}

func (im *Importer) importConcepts(ctx context.Context, cc <-chan []*snomed.Concept) {
	for {
		select {
		case <-ctx.Done():
			return
		case batch := <-cc:
			if batch == nil {
				return
			}
			err := im.storer.Put(ctx, batch)
			if err != nil {
				panic(err)
			}
			atomic.AddInt32(&im.nconcepts, int32(len(batch)))
		}
	}
}
func (im *Importer) importDescriptions(ctx context.Context, dd <-chan []*snomed.Description) {
	for {
		select {
		case <-ctx.Done():
			return
		case batch := <-dd:
			if batch == nil {
				return
			}
			err := im.storer.Put(ctx, batch)
			if err != nil {
				panic(err)
			}
			atomic.AddInt32(&im.ndescriptions, int32(len(batch)))
		}
	}
}
func (im *Importer) importRelationships(ctx context.Context, rels <-chan []*snomed.Relationship) {
	for {
		select {
		case <-ctx.Done():
			return
		case batch := <-rels:
			if batch == nil {
				return
			}
			err := im.storer.Put(ctx, batch)
			if err != nil {
				panic(err)
			}
			atomic.AddInt32(&im.nrelationships, int32(len(batch)))
		}
	}
}

func (im *Importer) importRefsets(ctx context.Context, refsets <-chan []*snomed.ReferenceSetItem) {
	for {
		select {
		case <-ctx.Done():
			return
		case batch := <-refsets:
			if batch == nil {
				return
			}
			err := im.storer.Put(ctx, batch)
			if err != nil {
				panic(err)
			}
			atomic.AddInt32(&im.nrefsets, int32(len(batch)))
		}
	}
}
