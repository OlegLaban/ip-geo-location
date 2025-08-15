package ipservice

import "errors"

var (
	ErrGetIP      = errors.New("can`t get ip")
	ErrDecodeBody = errors.New("can`t decode body for get ip")
)
