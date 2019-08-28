package terminology

import (
	"os"
	"testing"

	"github.com/wardle/go-terminology/snomed"
)

const fakeIndex = "fakeIndex"

func TestStore(t *testing.T) {
	defer os.RemoveAll(fakeIndex)
	bleve, err := newBleveIndex(fakeIndex, false)
	if err != nil {
		t.Fatal(err)
	}
	eds := make([]*snomed.ExtendedDescription, 0)
	fakeEd := snomed.ExtendedDescription{
		Concept: &snomed.Concept{
			Id:     24700007,
			Active: true,
		},
		Description: &snomed.Description{
			Id:           0,
			Term:         "Multiple sclerosis",
			Active:       true,
			ConceptId:    24700007,
			LanguageCode: "en",
		},
		PreferredDescription: &snomed.Description{
			Id:           0,
			Term:         "Multiple sclerosis",
			Active:       true,
			ConceptId:    24700007,
			LanguageCode: "en",
		},
		ConceptRefsets:     []int64{},
		DescriptionRefsets: []int64{},
		AllParentIds:       []int64{64572001},
	}
	eds = append(eds, &fakeEd)
	defer bleve.Close() // FILO
	if err := bleve.Index(eds); err != nil {
		t.Fatal(err)
	}
	r1, err := bleve.Search(&snomed.SearchRequest{
		S:   "mult scler",
		IsA: []int64{64572001},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(r1) != 1 {
		t.Fatalf("incorrect number of results. expected 1, got: %d", len(r1))
	}
	if r1[0] != 0 {
		t.Fatal("incorrect search result")
	}
	r2, err := bleve.Search(&snomed.SearchRequest{
		S: "parkin",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(r2) != 0 {
		t.Fatalf("incorrect number of results, expected 0, got: %d", len(r2))
	}
}
