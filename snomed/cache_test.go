package snomed

import (
	"testing"
)

func TestSimple(t *testing.T) {
	var cache = NewCache()
	cache.Put(1, "Hello")
	v, ok := cache.Get(1)
	if !ok || v != "Hello" {
		t.Error("Failed simple put/get.")
	}
	v, ok = cache.Get(2)
	if ok {
		t.Error("Cache thinks it has a value when it shouldn't")
	}
	v, err := cache.GetOrElse(2, func(id int) (interface{}, error) {
		return 3, nil
	})
	v, ok = cache.Get(2)
	if !ok || v != 3 {
		t.Error("Value not correctly cached.")
	}
	if err != nil {
		t.Error(err)
	}
}
