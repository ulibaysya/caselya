package service

import (
	"context"

	"github.com/ulibaysya/caselya/internal/modules/item"
)

type RNG interface {
	Int64() int64
}

type Storage interface {
	CreateItem(ctx context.Context, name string, price int64, tags []item.Tag) (int64, error)
}

type Service struct {
	strg Storage
	rng  RNG
}

func New(strg Storage, rng RNG) *Service {
	return &Service{
		strg: strg,
		rng:  rng,
	}
}
