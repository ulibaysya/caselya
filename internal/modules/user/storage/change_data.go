package storage

import (
	"context"
	"fmt"

	"github.com/ulibaysya/caselya/internal/modules/user"
)

func (s *Storage) ChangeData(ctx context.Context, id int64, changeFunc func(ctx context.Context, dc user.DataChanger) error) (err error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err = changeFunc(ctx, &dataChanger{
		id:   id,
		strg: s,
	}); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}

type dataChanger struct {
	id   int64
	strg *Storage
}

func (c *dataChanger) ChangeBalance(ctx context.Context, delta int64) (int64, error) {
	return c.strg.ChangeBalance(ctx, c.id, delta)
}

func (s *dataChanger) AddItems(ctx context.Context, itemsIDs []int64) ([]int64, error) {
	return s.strg.AddItems(ctx, s.id, itemsIDs)
}
