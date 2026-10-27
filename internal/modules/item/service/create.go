package service

import (
	"context"

	"github.com/ulibaysya/caselya/internal/modules/item"
)

func (s *Service) CreateItem(ctx context.Context, name string, price int64, tags []item.Tag) (int64, error) {
	return s.strg.CreateItem(ctx, name, price, tags)
}
