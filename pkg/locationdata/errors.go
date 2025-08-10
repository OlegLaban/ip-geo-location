package locationdata

import "errors"

var (
	ErrCantGetDataFromCache = errors.New("can`t get data from cache")
	ErrCantGetIP = errors.New("can`t get ip")
	ErrCantGetGeoData = errors.New("can`t get geo data")
	ErrCantDecodeGeoData = errors.New("can`t decode geo data")
	ErrGetImageViaHttp  = errors.New("can`t get image of flag via http")
)