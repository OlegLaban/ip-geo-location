package cache

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/OlegLaban/geo-flag/pkg/locationdata"
)

type CacheService struct {
	storage map[string][]byte
	logger *slog.Logger
}

func NewCacheService(logger *slog.Logger) *CacheService {
	return &CacheService{storage: make(map[string][]byte), logger: logger}
}

func (cs *CacheService) Set(key string, data []byte) error {
	cs.storage[key] = data
	cs.logger.Debug(fmt.Sprintf("data was seted success key - %s, len - %d", key, len(data)))
	return nil
}

func (cs *CacheService) Get(key string) ([]byte, error) {
	data, ok := cs.storage[key]
	if !ok {
		cs.logger.Debug(fmt.Sprintf("key %s not found in cache", key))
		return []byte{}, errors.Join(locationdata.ErrCantGetDataFromCache, fmt.Errorf("not found data via key %s", key))
	}
	cs.logger.Debug(fmt.Sprintf("key was got from cache key - %s, len - %d", key, len(data)))

	return data, nil
}

func (cs *CacheService) GetWithCallback(key string, callback func () ([]byte, error)) ([]byte, error) {
	if !cs.Has(key) {
		cs.logger.Debug(fmt.Sprintf("key - %s, not found, try callback", key))
		val, err := callback()
		if err != nil {
			cs.logger.Error("can`t get data for cache from callback", "err", err)
			return []byte{}, errors.Join(ErrGetDataFromCallback, err)
		}
		cs.Set(key, val)
	}

	return cs.Get(key)
}

func (cs *CacheService) Delete(key string) error {
	cs.logger.Debug(fmt.Sprintf("try delete key - %s from cache", key))
	delete(cs.storage, key)

	return nil
}

func (cs *CacheService) Has(key string) bool {
	cs.logger.Debug(fmt.Sprintf("try verify key - %s in cache", key))
	_, ok := cs.storage[key]
	return ok
}