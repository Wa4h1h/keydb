package store

import (
	"sync"
	"time"

	"github.com/Wa4h1h/keydb/internal/utils"
)

type Store interface {
	AddItem(key string, item *Item)
	GetItem(key string) (*Item, error)
	DeleteItem(key string) (string, error)
}

type Item struct {
	TTL       int64
	CreatedAt time.Time
	Value     string
}

type MemStore struct {
	Ms                          map[string]*Item
	ttlBackgroundWorkerInterval int
	sync.RWMutex
}

func NewMemStore(ttlBackgroundWorkerInterval int) Store {
	ms := &MemStore{
		Ms:                          make(map[string]*Item),
		ttlBackgroundWorkerInterval: ttlBackgroundWorkerInterval,
	}

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
	if ok {
		return "", utils.ErrItemNotRemoved
	}

	return "OK\n", nil
}

func (ms *MemStore) BackgroundTTLWorker() {
	for {
		ms.Lock()

		for key, value := range ms.Ms {
			diff := value.CreatedAt.
				Add(time.Duration(value.TTL) * time.Second).
				Sub(time.Now())

			if diff < 0 {
				delete(ms.Ms, key)
			}
		}

		ms.Unlock()

		time.Sleep(time.Duration(ms.ttlBackgroundWorkerInterval) * time.Minute)
	}
}
