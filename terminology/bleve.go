package terminology

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/wardle/go-terminology/snomed"
)

// bleveService encapsulates the bleve search functionality
type bleveService struct {
	index bleve.Index
}

// document is the document indexed by bleve for fast free text search and autocompletion
type document struct {
	ID       string   // description ID
	Term     string   // the term itself
	Keywords []string // list of keywords permitting faceted search
}

// NewBleveIndex creates or opens a bleve index at the location specified.
func newBleveIndex(path string, readOnly bool) (*bleveService, error) {
	config := map[string]interface{}{
		"read_only": readOnly,
	}
	index, err := bleve.OpenUsing(path, config)
	if err == nil {
		return &bleveService{index: index}, err
	}
	if err != bleve.ErrorIndexPathDoesNotExist {
		return nil, err
	}

	if readOnly {
		return nil, fmt.Errorf("cannot open index in read-only mode: index doesn't exist at %s", path)
	}
	indexMapping := bleve.NewIndexMapping()
	documentMapping := bleve.NewDocumentMapping() // index only a single type of document
	indexMapping.AddDocumentMapping("document", documentMapping)
	indexMapping.DefaultType = "document"
	// the id
	idMapping := bleve.NewTextFieldMapping()
	idMapping.IncludeInAll = false
	idMapping.IncludeTermVectors = false
	idMapping.Store = true
	idMapping.Analyzer = keyword.Name

	// the term
	termMapping := bleve.NewTextFieldMapping()
	termMapping.Analyzer = "en"
	termMapping.Store = false
	documentMapping.AddFieldMappingsAt("Term", termMapping)

	// the keywords
	keywordMapping := bleve.NewTextFieldMapping()
	keywordMapping.Analyzer = keyword.Name
	keywordMapping.Store = false
	keywordMapping.IncludeInAll = false
	keywordMapping.IncludeTermVectors = false
	documentMapping.AddFieldMappingsAt("Keywords", keywordMapping)

	index, err = bleve.NewUsing(path, indexMapping, scorch.Name, scorch.Name, nil)
	return &bleveService{index: index}, err
}

func (bs *bleveService) Index(eds []*snomed.ExtendedDescription) error {
	batch := bs.index.NewBatch()
	docs := make([]document, len(eds))
	for i, ed := range eds {
		if ed.GetDescription().IsFullySpecifiedName() { // always omit FSN from the index
			continue
		}
		docs[i].Term = ed.GetDescription().GetTerm()
		docs[i].ID = strconv.FormatInt(ed.GetDescription().GetId(), 10)
		var kws = keywords{
			recursiveParents:   ed.GetRecursiveParentIds(),
			directParents:      ed.GetDirectParentIds(),
			conceptRefsets:     ed.GetConceptRefsets(),
			descriptionRefsets: ed.GetDescriptionRefsets(),
			conceptActive:      ed.GetConcept().GetActive(),
			descriptionActive:  ed.GetDescription().GetActive(),
		}
		docs[i].Keywords = kws.toKeywords()
		if err := batch.Index(docs[i].ID, &docs[i]); err != nil {
			return err
		}
	}
	return bs.index.Batch(batch)
}

type keywords struct {
	recursiveParents   []int64
	directParents      []int64
	conceptRefsets     []int64
	descriptionRefsets []int64
	conceptActive      bool
	descriptionActive  bool
}

func (kw keywords) toKeywords() []string {
	words := make([]string, 0)
	writeIdentifiers(&words, "rp", kw.recursiveParents)
	writeIdentifiers(&words, "dp", kw.directParents)
	writeIdentifiers(&words, "cr", kw.conceptRefsets)
	writeIdentifiers(&words, "dr", kw.descriptionRefsets)
	if kw.conceptActive {
		words = append(words, "ca")
	}
	if kw.descriptionActive {
		words = append(words, "da")
	}
	return words
}

func writeIdentifiers(words *[]string, prefix string, ids []int64) {
	var sb strings.Builder
	for _, id := range ids {
		sb.WriteString(prefix)
		sb.WriteString(strconv.FormatInt(id, 10))
		*words = append(*words, sb.String())
		sb.Reset()
	}
}

func (bs *bleveService) Search(sr *snomed.SearchRequest) ([]int64, error) {

	if len(sr.GetIsA()) == 0 {
		sr.IsA = []int64{138875005}
	}
	if sr.MaximumHits == 0 {
		sr.MaximumHits = 100
	}
	if sr.S == "" {
		// TODO: implement list of recursive children, up to a maximum (useful for drop-downs)
		return nil, fmt.Errorf("No search string in request")
	}

	query := bleve.NewConjunctionQuery()
	for _, token := range strings.Split(sr.S, " ") {
		tokenQuery := bleve.NewMatchQuery(token)
		tokenQuery.SetField("Term")
		if len(token) < 3 {
			query.AddQuery(tokenQuery)
			continue
		}
		termQuery := bleve.NewDisjunctionQuery()
		termQuery.AddQuery(tokenQuery)
		prefixQuery := bleve.NewPrefixQuery(token)
		prefixQuery.SetField("Term")
		termQuery.AddQuery(prefixQuery)
		if sr.Fuzzy == snomed.SearchRequest_ALWAYS_FUZZY {
			fuzzyQuery := bleve.NewFuzzyQuery(token)
			fuzzyQuery.SetField("Term")
			fuzzyQuery.SetFuzziness(2)
			termQuery.AddQuery(fuzzyQuery)
		}
		query.AddQuery(termQuery)
	}
	var kws = keywords{
		recursiveParents:   sr.GetIsA(),
		directParents:      sr.GetDirectParents(),
		conceptRefsets:     sr.GetConceptRefsets(),
		descriptionRefsets: sr.GetDescriptionRefsets(),
		conceptActive:      !sr.GetIncludeInactive(),
	}
	keywords := kws.toKeywords()
	if len(keywords) > 0 {
		keywordsQuery := bleve.NewConjunctionQuery()
		for _, keyword := range keywords {
			kq := bleve.NewTermQuery(keyword)
			kq.SetField("Keywords")
			keywordsQuery.AddQuery(kq)
		}
		query.AddQuery(keywordsQuery)
	}
	req := bleve.NewSearchRequest(query)
	req.Size = int(sr.GetMaximumHits())
	result, err := bs.index.Search(req)
	if err != nil {
		return nil, err
	}

	results := make([]int64, 0)
	for _, hit := range result.Hits {
		id, err := strconv.ParseInt(hit.ID, 10, 64)
		if err != nil {
			return nil, err
		}
		results = append(results, id)
	}

	// perform fallback if no hits, and if requested.
	if (len(results) == 0) && (sr.Fuzzy == snomed.SearchRequest_FALLBACK_FUZZY) {
		sr.Fuzzy = snomed.SearchRequest_ALWAYS_FUZZY
		return bs.Search(sr)
	}
	return results, nil
}

func (bs *bleveService) Close() error {
	bs.index.Close()
	return nil
}
