package store

import (
	"sync"

	"github.com/Wa4h1h/memdb/internal/utils"
)

type Store interface {
	AddItem(key string, item *Item)
	GetItem(key string) (*Item, error)
	DeleteItem(key string) (string, error)
}

type Item struct {
	Value string
	TTL   int64
}

type MemStore struct {
	Ms map[string]*Item
	m  sync.RWMutex
}

func NewMemStore() *MemStore {
	return &MemStore{Ms: make(map[string]*Item)}
}

func (d *MemStore) AddItem(key string, item *Item) {
	d.m.Lock()
	defer d.m.Unlock()

	d.Ms[key] = item
}

func (d *MemStore) GetItem(key string) (*Item, error) {
	d.m.RLock()
	defer d.m.RUnlock()

	val, ok := d.Ms[key]
	if !ok {
		return nil, utils.ErrNotFoundItem
	}

	return val, nil
}

func (d *MemStore) DeleteItem(key string) (string, error) {
	d.m.Lock()
	defer d.m.Unlock()

	delete(d.Ms, key)

	_, ok := d.Ms[key]
	if !ok || len(d.Ms) == 0 {
		return "", utils.ErrItemNotRemoved
	}

	return "OK\n", nil
}
