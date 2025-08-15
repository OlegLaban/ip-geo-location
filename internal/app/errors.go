package app

import "errors"

var (
	ErrLoadGeoData = errors.New("get country data err")
	ErrLoadFlag    = errors.New("flag generation err")
)
