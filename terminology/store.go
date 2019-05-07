package terminology

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wardle/go-terminology/snomed"
	"strings"
)

type bucket int

const (
	bkConcepts      bucket = iota // concepts, keyed by SCTID (uint64)
	bkDescriptions                // descriptions, keyed by SCTID (uint64)
	bkRelationships               // relationships, keyed by SCTID (uint64)
	bkRefsetItems                 // refset items, keyed by their uuid (string)

	ixConceptDescriptions        // key: concept_id-description_id
	ixConceptParentRelationships // key: concept_id-relationship_id
	ixConceptChildRelationships  // key: concept_id-relationship_id
	ixComponentReferenceSets     // key: component_id-refset_id

	ixReferenceSetComponentItems // key: refset_id-component_id-reference_set_item_id
	ixRefsetTargetItems          // key: refset_id-target_code-SPACE-reference_set_item_id

	ixReferenceSets // key: refset_id

	ixConceptRecursiveParents // keys: concept_id-parent_id

)

var bucketNames = [...][]byte{
	[]byte("con"), // key: sct_id value: concept
	[]byte("des"), // key: sct_id value: description
	[]byte("rel"), // key: sct_id value: relationship
	[]byte("ref"), // key: uuid value: component

	[]byte("cds"),
	[]byte("cpr"),
	[]byte("ccr"),
	[]byte("crs"),

	[]byte("rci"),
	[]byte("rti"),
	[]byte("rfs"),

	[]byte("crp"), // key: concept_id-parent_id
}

func (b bucket) name() []byte {
	return bucketNames[b]
}

func compoundKey(keys ...[]byte) []byte {
	return bytes.Join(keys, nil)
}

// Batch represents an abstract batch operation against the KV store
type Batch interface {
	// Get an object from the specified bucket with the specified key
	Get(b bucket, key []byte, value proto.Message) error

	// Put and object into the specified bucket with the specified key
	Put(b bucket, key []byte, value proto.Message)

	// Add an index entry for the specified bucket and key
	AddIndexEntry(b bucket, key []byte, value []byte)

	// Get all index entries for the specified bucket and key
	GetIndexEntries(b bucket, key []byte) ([][]byte, error)

	// Clear all index entries for the specified bucket and key
	ClearIndexEntries(b bucket, key []byte) error

	// Iterate iterates through a bucket
	Iterate(b bucket, keyPrefix []byte, f func(key, value []byte) error) error
}

// Store is an abstract key-value store divided into logical buckets of information.
type Store interface {
	// View creates a read-only transaction
	View(func(Batch) error) error

	// Update creates a read and write transaction
	Update(func(Batch) error) error

	// Close closes any resources associated with the key-value store
	Close() error
}

// ErrDatabaseNotInitialised is the error when database not properly initialised
var ErrDatabaseNotInitialised = errors.New("database not initialised")

// ErrNotFound is the error when something isn't found
var ErrNotFound = errors.New("Not found")

// Search represents the backend opaque abstract SNOMED-CT search service.
type Search interface {
	Index(eds []*snomed.ExtendedDescription) error
	Search(sr *snomed.SearchRequest) ([]int64, error) //TODO: rename autocomplete
	Close() error
}

// Statistics on the persistence store
type Statistics struct {
	concepts      int
	descriptions  int
	relationships int
	refsetItems   int
	refsets       []string
}

func (st Statistics) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Number of concepts: %d\n", st.concepts))
	b.WriteString(fmt.Sprintf("Number of descriptions: %d\n", st.descriptions))
	b.WriteString(fmt.Sprintf("Number of relationships: %d\n", st.relationships))
	b.WriteString(fmt.Sprintf("Number of reference set items: %d\n", st.refsetItems))
	b.WriteString(fmt.Sprintf("Number of installed refsets: %d:\n", len(st.refsets)))

	for _, s := range st.refsets {
		b.WriteString(fmt.Sprintf("  Installed refset: %s\n", s))
	}
	return b.String()
}