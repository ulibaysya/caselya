package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) AddItems(ctx context.Context, userID int64, ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	sql, args := buildAddItemsSQL(userID, ids)
	queryRows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("adding items to the user: %w", err)
	}
	type resultedRow struct {
		InventoryItemID int64
		ItemID          int64
	}
	resultedRows, err := pgx.CollectRows(queryRows, func(row pgx.CollectableRow) (resultedRow, error) {
		newRow := resultedRow{}
		err := row.Scan(&newRow.InventoryItemID, &newRow.ItemID)
		return newRow, err
	})
	if err != nil {
		return nil, fmt.Errorf("collecting rows: %w", err)
	}
	if len(resultedRows) != len(ids) {
		return nil, fmt.Errorf("unexpected number of collected rows: got %d, expected %d", len(resultedRows), len(ids))
	}

	inventoryIDsByItem := map[int64][]int64{}
	for _, i := range resultedRows {
		inventoryIDsByItem[i.ItemID] = append(inventoryIDsByItem[i.ItemID], i.InventoryItemID)
	}

	inventoryItemsIDs := make([]int64, 0, len(ids))
	for _, i := range ids {
		matchedInventoryItemsIDs, ok := inventoryIDsByItem[i]
		if !ok || len(matchedInventoryItemsIDs) == 0 {
			return nil, fmt.Errorf("can't find added items for item %d", i)
		}
		inventoryItemsIDs = append(inventoryItemsIDs, matchedInventoryItemsIDs[0])
		inventoryIDsByItem[i] = matchedInventoryItemsIDs[1:]
	}

	return inventoryItemsIDs, nil
}

func buildAddItemsSQL(userID int64, itemsIDs []int64) (string, []any) {
	b := &strings.Builder{}

	b.WriteString("INSERT INTO users_items (users_id, items_id) VALUES ($1, $2)")
	args := []any{userID, itemsIDs[0]}

	argIndex := 3
	for _, i := range itemsIDs[1:] {
		fmt.Fprintf(b, ", ($1, $%d)", argIndex)
		args = append(args, i)
		argIndex += 1
	}

	b.WriteString(" RETURNING id, items_id;")

	return b.String(), args
}
