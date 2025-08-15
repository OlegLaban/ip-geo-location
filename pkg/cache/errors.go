package cache

import "errors"

var (
	ErrGetDataFromCallback  = errors.New("can`t get data from callback")
	ErrSetDatToCache        = errors.New("ca`t set data to cache")
	ErrCantGetDataFromCache = errors.New("can`t get data from cache")
)
