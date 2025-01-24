package common

import "errors"

var (
	ErrNoItems = errors.New("items is empty")
	ErrNoStock = errors.New("some item is not in stock")
)
