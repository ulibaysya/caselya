package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/ulibaysya/caselya/internal/modules/item"
)

func (s *Storage) CreateItem(ctx context.Context, name string, price int64, tags []item.Tag) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("beginning tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if len(tags) != 0 {
		for _, tag := range tags {
			if len(tag.Values) == 0 {
				// TODO return specific wrapped error
				return 0, fmt.Errorf("tag %v has 0 values", tag.Key)
			}
		}
	}

	var itemID int64
	const itemInsertSQL = `
		INSERT INTO items 
		DEFAULT VALUES 
		RETURNING id;
	`
	if err = tx.QueryRow(ctx, itemInsertSQL).Scan(&itemID); err != nil {
		return 0, fmt.Errorf("inserting item: %w", err)
	}

	if len(tags) != 0 {
		tagsInsertSQL, args := buildTagsInsertSQL(itemID, tags)
		if _, err = tx.Exec(ctx, tagsInsertSQL, args...); err != nil {
			// TODO check constraint errors (key not found, value not found)
			return 0, fmt.Errorf("inserting tags: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commiting tx: %w", err)
	}

	return itemID, nil
}

func buildTagsInsertSQL(itemID int64, tags []item.Tag) (string, []any) {
	b := &strings.Builder{}

	b.WriteString(`
		INSERT INTO item_tag_relations 
		SELECT $1, v.id 
		FROM item_tag_values AS v 
		JOIN item_tag_keys AS k ON v.tag_key_id = k.id 
		JOIN (
			VALUES 
				($2, $3)
	`)

	args := []any{itemID, tags[0].Key, tags[0].Values[0]}

	latestKeyIndex := 2
	latestValueIndex := latestKeyIndex + 2
	for _, value := range tags[0].Values[1:] {
		fmt.Fprintf(b, `,($%d, $%d)`, latestKeyIndex, latestValueIndex)
		args = append(args, value)
		latestValueIndex += 1
	}
	latestKeyIndex = latestValueIndex

	for _, tag := range tags[1:] {
		args = append(args, tag.Key)
		latestValueIndex := latestKeyIndex + 1
		for _, value := range tag.Values {
			fmt.Fprintf(b, `,($%d, $%d)`, latestKeyIndex, latestValueIndex)
			args = append(args, value)
			latestValueIndex += 1
		}
		latestKeyIndex = latestValueIndex
	}

	b.WriteString(`
		) AS tags(key_name, value_name)
		ON tags.key_name = k.name
		AND tags.value_name = v.name;
	`)

	return b.String(), args
}
