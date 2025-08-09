package cache

import (
	"errors"
	"fmt"

	"github.com/OlegLaban/geo-flag/pkg/locationdata"
)

type CacheService struct {
	storage map[string][]byte
}

func NewCacheService() *CacheService {
	return &CacheService{storage: make(map[string][]byte)}
}

func (cs *CacheService) Set(key string, data []byte) error {
	cs.storage[key] = data

	return nil
}

func (cs *CacheService) Get(key string) ([]byte, error) {
	data, ok := cs.storage[key]
	if !ok {
		return []byte{}, errors.Join(locationdata.ErrCantGetDataFromCache, fmt.Errorf("not found data via key %s", key))
	}

	return data, nil
}

func (cs *CacheService) GetWithCallback(key string, callback func () ([]byte, error)) ([]byte, error) {
	if !cs.Has(key) {
		val, err := callback()
		if err != nil {
			return []byte{}, errors.Join(ErrGetDataFromCallback, err)
		}
		cs.Set(key, val)
	}

	return cs.Get(key)
}

func (cs *CacheService) Delete(key string) error {
	delete(cs.storage, key)

	return nil
}

func (cs *CacheService) Has(key string) bool {
	_, ok := cs.storage[key]
	return ok
}