package container

import "errors"

var (
	ErrInvalidOpeningCount = errors.New("invalid opening count")
	ErrInvalidOpeningPrice = errors.New("invalid opening price")
	ErrNotFound            = errors.New("not found")
)
