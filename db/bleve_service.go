package db

import (
	"github.com/blevesearch/bleve"
)

type BleveService struct {
	path  string
	index bleve.Index
}

func NewBleveService(path string, readOnly bool) (*BleveService, error) {
	return &BleveService{path: path, index: nil}, nil
}

func (bs BleveService) Search(search *SearchRequest) ([]int, error) {
	panic("Not implemented")
}

func (bs BleveService) Close() error {
	return nil
}
