package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ulibaysya/caselya/internal/modules/container"
	"github.com/ulibaysya/caselya/internal/modules/item"
)

func (s *Storage) FindContainer(ctx context.Context, id int64) (container.Container, error) {
	const sql = `
		SELECT c.id, c.price, ci.items_id, ci.weight
		FROM containers AS c
		JOIN containers_items AS ci
		ON c.id = ci.containers_id
		WHERE c.id = $1
		ORDER BY ci.weight, ci.items_id;
	`

	queryRows, err := s.db.Query(ctx, sql, id)
	if err != nil {
		return container.Container{}, fmt.Errorf("querying container: %w", err)
	}

	type resultedRow struct {
		containerID    int64
		containerPrice int64
		itemID         int64
		itemWeight     int64
	}

	resultedRows, err := pgx.CollectRows(queryRows, func(row pgx.CollectableRow) (resultedRow, error) {
		newRow := resultedRow{}
		err := row.Scan(&newRow.containerID, &newRow.containerPrice, &newRow.itemID, &newRow.itemWeight)
		return newRow, err
	})
	if err != nil {
		return container.Container{}, fmt.Errorf("collecting rows: %w", err)
	}

	if len(resultedRows) == 0 {
		return container.Container{}, container.ErrNotFound
	}

	items := make([]item.WithWeight, len(resultedRows))
	for i := range resultedRows {
		items[i].ID = resultedRows[i].itemID
		items[i].Weight = resultedRows[i].itemWeight
	}

	return container.Container{
		ID:    resultedRows[0].containerID,
		Price: resultedRows[0].containerPrice,
		Items: items,
	}, nil
}
