package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/ulibaysya/caselya/internal/modules/container"
	"github.com/ulibaysya/caselya/internal/modules/user"
	pkgslog "github.com/ulibaysya/caselya/pkg/slog"
)

func (s *Service) OpenContainers(ctx context.Context, userID int64, containerID int64, n int) ([]container.OpeningResult, error) {
	const maxN = 100
	if n <= 0 || n > maxN {
		return nil, fmt.Errorf("%w: must be between 1 and %d, got %d", container.ErrInvalidOpeningCount, maxN, n)
	}

	contToOpen, err := s.containerStorage.FindContainer(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("finding container: %w", err)
	}
	s.logger.LogAttrs(ctx, slog.LevelDebug, "found container",
		slog.Int64("id", contToOpen.ID),
		slog.Int64("price", contToOpen.Price),
		slog.Any("items", contToOpen.Items),
	)

	openingPrice, err := calculateOpeningPrice(n, contToOpen.Price)
	if err != nil {
		s.logger.LogAttrs(ctx, slog.LevelWarn, "failed counting opening price",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("%w: %w", container.ErrInvalidOpeningPrice, err)
	}
	s.logger.LogAttrs(ctx, slog.LevelDebug, "counted opening price",
		slog.Int64("openingPrice", openingPrice),
	)

	var wonItems []int64
	var inventoryItems []int64
	if err = s.userStorage.ChangeData(ctx, userID, func(ctx context.Context, dc user.DataChanger) error {
		balanceDelta := -openingPrice
		newBalance, err := dc.ChangeBalance(ctx, balanceDelta)
		if err != nil {
			return fmt.Errorf("changing user balance: %w", err)
		}
		s.logger.LogAttrs(ctx, slog.LevelDebug, "successfully changed user balance",
			slog.Int64("newBalance", newBalance),
		)

		wonItemsIndexes, err := s.itemPicker.Pick(ctx, contToOpen.Items, n)
		if err != nil {
			return fmt.Errorf("picking items from container: %w", err)
		}
		wonItems = make([]int64, 0, len(wonItemsIndexes))
		for _, i := range wonItemsIndexes {
			wonItems = append(wonItems, contToOpen.Items[i].ID)
		}
		s.logger.LogAttrs(ctx, slog.LevelDebug, "successfully picked items from container",
			slog.Any("indexesInContainer", wonItemsIndexes),
			slog.Any("itemsIDs", wonItems),
		)

		inventoryItems, err = dc.AddItems(ctx, wonItems)
		if err != nil {
			s.logger.LogAttrs(ctx, slog.LevelError, "got error when adding items to the user",
				pkgslog.Error("error", err),
			)
			return fmt.Errorf("adding items to user: %w", err)
		}
		if len(wonItems) != len(inventoryItems) {
			return fmt.Errorf("adding items to user: expected %d inventory item ids, got %d", len(wonItems), len(inventoryItems))
		}
		s.logger.LogAttrs(ctx, slog.LevelDebug, "successfully added picked items to user",
			slog.Any("inventoryItemsIDs", inventoryItems),
		)

		return nil
	}); err != nil {
		return nil, err
	}
	openingResults := make([]container.OpeningResult, 0, len(inventoryItems))
	for i := range inventoryItems {
		openingResults = append(openingResults, container.OpeningResult{
			ItemID:          wonItems[i],
			InventoryItemID: inventoryItems[i],
		})
	}

	s.logger.LogAttrs(ctx, slog.LevelDebug, "successfully performed opening",
		slog.Any("openingResults", openingResults),
	)

	return openingResults, nil
}

func calculateOpeningPrice(n int, price int64) (int64, error) {
	if price < 0 {
		return 0, fmt.Errorf("negative container price: %d", price)
	}
	if price == 0 {
		return 0, nil
	}

	if int64(n) > math.MaxInt64/price {
		return 0, fmt.Errorf("%d openings and %d container price overflow int64", n, price)
	}

	return int64(n) * price, nil
}
