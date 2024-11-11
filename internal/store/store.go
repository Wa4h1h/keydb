package store

import (
	"sync"
	"time"

	"github.com/Wa4h1h/memdb/internal/utils"
)

type Store interface {
	AddItem(key string, item *Item)
	GetItem(key string) (*Item, error)
	DeleteItem(key string) (string, error)
}

type Item struct {
	TTL       int64
	Value     string
	CreatedAt time.Time
}

type MemStore struct {
	Ms map[string]*Item
	sync.RWMutex
}

func NewMemStore() *MemStore {
	ms := &MemStore{Ms: make(map[string]*Item)}

	go ms.BackgroundTTLWorker()

	return ms
}

func (ms *MemStore) AddItem(key string, item *Item) {
	ms.Lock()
	defer ms.Unlock()

	ms.Ms[key] = item
}

func (ms *MemStore) GetItem(key string) (*Item, error) {
	ms.RLock()
	defer ms.RUnlock()

	val, ok := ms.Ms[key]
	if !ok {
		return nil, utils.ErrNotFoundItem
	}

	return val, nil
}

func (ms *MemStore) DeleteItem(key string) (string, error) {
	ms.Lock()
	defer ms.Unlock()

	delete(ms.Ms, key)

	_, ok := ms.Ms[key]
	if !ok || len(ms.Ms) == 0 {
		return "", utils.ErrItemNotRemoved
	}

	return "OK\n", nil
}

func (ms *MemStore) BackgroundTTLWorker() {
	ms.Lock()
	defer ms.Unlock()

	for key, value := range ms.Ms {
		diff := value.CreatedAt.
			Add(time.Duration(value.TTL) * time.Second).
			Sub(time.Now())

		if diff < 0 {
			delete(ms.Ms, key)
		}
	}
}
