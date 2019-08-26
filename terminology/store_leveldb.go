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
	"bytes"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// levelStore is a concrete file-based database store for SNOMED-CT using goleveldb
type levelStore struct {
	db *leveldb.DB
}

type levelBatch struct {
	batch  leveldb.Batch
	store  *levelStore
	errors []error
}

func (ls *levelStore) Update(f func(Batch) error) error {
	batch := &levelBatch{
		store: ls,
	}
	err := f(batch)
	if err != nil {
		return err
	}
	if len(batch.errors) > 0 {
		return fmt.Errorf("errors on update: %v", batch.errors)
	}
	return ls.db.Write(&batch.batch, nil)
}

func (ls *levelStore) View(f func(Batch) error) error {
	batch := &levelBatch{
		store: ls,
	}
	return f(batch)
}

func (lb *levelBatch) Get(b bucket, key []byte, pb proto.Message) error {
	d, err := lb.store.db.Get(bytes.Join([][]byte{b.name(), key}, nil), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return ErrNotFound
		}
		return err
	}
	return proto.Unmarshal(d, pb)
}

func (lb *levelBatch) GetIndexEntries(b bucket, key []byte) ([][]byte, error) {
	prefix := bytes.Join([][]byte{b.name(), key}, nil)
	lp := len(prefix)
	iter := lb.store.db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	result := make([][]byte, 0)
	for iter.Next() {
		k := iter.Key()
		entry := k[lp:]
		entry2 := make([]byte, len(entry))
		copy(entry2, entry)
		result = append(result, entry2) // we have to store a copy
	}
	return result, iter.Error()
}

func (lb *levelBatch) Put(b bucket, key []byte, value proto.Message) {
	d, err := proto.Marshal(value)
	if err != nil {
		lb.errors = append(lb.errors, err)
	}
	k := bytes.Join([][]byte{b.name(), key}, nil)
	lb.batch.Put(k, d)
}

func (lb *levelBatch) AddIndexEntry(b bucket, key []byte, value []byte) {
	k := bytes.Join([][]byte{b.name(), key, value}, nil)
	lb.batch.Put(k, []byte{'.'})
}

func (lb *levelBatch) ClearIndexEntries(b bucket) {
	batch := new(leveldb.Batch)
	count := 0
	lb.Iterate(b, nil, func(k, v []byte) error {
		if count < 5000 {
			batch.Delete(k)
			count++
		} else {
			count = 0
			lb.store.db.Write(batch, nil)
			batch.Reset()
		}
		return nil
	})
	lb.store.db.Write(batch, nil)
}

func (lb *levelBatch) Iterate(b bucket, keyPrefix []byte, f func(key, value []byte) error) error {
	k := bytes.Join([][]byte{b.name(), keyPrefix}, nil)
	iter := lb.store.db.NewIterator(util.BytesPrefix(k), nil)
	defer iter.Release()

	for iter.Next() {
		if err := f(iter.Key(), iter.Value()); err != nil {
			return err
		}
	}
	return iter.Error()
}

func (ls *levelStore) Close() error {
	return ls.db.Close()
}

func newLevelService(filename string, readOnly bool) (*levelStore, error) {
	opts := opt.Options{ReadOnly: readOnly}
	db, err := leveldb.OpenFile(filename, &opts)
	if err != nil {
		return nil, err
	}
	return &levelStore{
		db: db,
	}, nil
}
