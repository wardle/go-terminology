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
	ID   string // description ID
	Term string // the term itself

	RecursiveParents   []string
	DirectParents      []string
	ConceptRefsets     []string
	DescriptionRefsets []string
	ConceptActive      bool
	DescriptionActive  bool
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
	boolMapping := bleve.NewBooleanFieldMapping()
	boolMapping.IncludeInAll = false
	boolMapping.Store = false

	documentMapping.AddFieldMappingsAt("RecursiveParents", keywordMapping)
	documentMapping.AddFieldMappingsAt("DirectParents", keywordMapping)
	documentMapping.AddFieldMappingsAt("ConceptRefsets", keywordMapping)
	documentMapping.AddFieldMappingsAt("DescriptionRefsets", keywordMapping)
	documentMapping.AddFieldMappingsAt("DirectParents", keywordMapping)
	documentMapping.AddFieldMappingsAt("ConceptActive", boolMapping)
	documentMapping.AddFieldMappingsAt("DescriptionActive", boolMapping)

	index, err = bleve.NewUsing(path, indexMapping, scorch.Name, scorch.Name, nil)
	return &bleveService{index: index}, err
}

func (bs *bleveService) Statistics() (uint64, error) {
	return bs.index.DocCount()
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
		for _, id := range ed.GetAllParentIds() {
			docs[i].RecursiveParents = append(docs[i].RecursiveParents, strconv.FormatInt(id, 10))
		}
		for _, id := range ed.GetDirectParentIds() {
			docs[i].DirectParents = append(docs[i].DirectParents, strconv.FormatInt(id, 10))
		}
		for _, id := range ed.GetConceptRefsets() {
			docs[i].ConceptRefsets = append(docs[i].ConceptRefsets, strconv.FormatInt(id, 10))
		}
		for _, id := range ed.GetDescriptionRefsets() {
			docs[i].DescriptionRefsets = append(docs[i].DescriptionRefsets, strconv.FormatInt(id, 10))
		}
		docs[i].ConceptActive = ed.GetConcept().GetActive()
		docs[i].DescriptionActive = ed.GetDescription().GetActive()
		if err := batch.Index(docs[i].ID, &docs[i]); err != nil {
			return err
		}
	}
	return bs.index.Batch(batch)
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
	s := strings.TrimSpace(sr.S)
	query := bleve.NewConjunctionQuery()
	for _, token := range strings.Split(s, " ") {
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
	if len(sr.GetIsA()) > 0 {
		qs := bleve.NewDisjunctionQuery()
		for _, id := range sr.GetIsA() {
			q := bleve.NewTermQuery(strconv.FormatInt(id, 10))
			q.SetField("RecursiveParents")
			qs.AddQuery(q)
		}
		query.AddQuery(qs)
	}
	if len(sr.GetDirectParents()) > 0 {
		qs := bleve.NewDisjunctionQuery()
		for _, id := range sr.GetDirectParents() {
			q := bleve.NewTermQuery(strconv.FormatInt(id, 10))
			q.SetField("DirectParents")
			qs.AddQuery(q)
		}
		query.AddQuery(qs)
	}
	if len(sr.GetConceptRefsets()) > 0 {
		qs := bleve.NewDisjunctionQuery()
		for _, id := range sr.GetConceptRefsets() {
			q := bleve.NewTermQuery(strconv.FormatInt(id, 10))
			q.SetField("ConceptRefsets")
			qs.AddQuery(q)
		}
		query.AddQuery(qs)
	}
	if len(sr.GetDescriptionRefsets()) > 0 {
		qs := bleve.NewDisjunctionQuery()
		for _, id := range sr.GetDescriptionRefsets() {
			q := bleve.NewTermQuery(strconv.FormatInt(id, 64))
			q.SetField("DescriptionRefsets")
			qs.AddQuery(q)
		}
		query.AddQuery(qs)
	}
	if !sr.GetIncludeInactive() { // unless explicitly requested, include only active concepts
		q := bleve.NewBoolFieldQuery(true)
		q.SetField("ConceptActive")
		query.AddQuery(q)
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
