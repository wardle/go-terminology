package terminology

import (
	"github.com/blevesearch/bleve"
)

type bleveService struct {
	path  string
	index bleve.Index
}

func newBleveService(path string, readOnly bool) (*bleveService, error) {
	return &bleveService{path: path, index: nil}, nil
}

func (bs bleveService) Search(search *SearchRequest) ([]int, error) {
	panic("Not implemented")
}

func (bs bleveService) Close() error {
	return nil
}
