package service

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ulibaysya/caselya/internal/modules/container"
	"github.com/ulibaysya/caselya/internal/modules/item"
	"github.com/ulibaysya/caselya/internal/modules/user"
	pkgslog "github.com/ulibaysya/caselya/pkg/slog"
)

func TestOpeningContainers_Success(t *testing.T) {
	containerStorage := NewMockContainerStorage(t)
	containerID := int64(589)
	containerInstance := container.Container{
		ID:    containerID,
		Price: 100,
		Items: []item.WithWeight{
			{ID: 45669957, Weight: 9},
			{ID: 5501121, Weight: 2},
			{ID: 37675023, Weight: 3},
			{ID: 1510940, Weight: 4},
			{ID: 20735405, Weight: 5},
		},
	}
	containerStorage.
		On("FindContainer", mock.Anything, containerID).
		Return(containerInstance, nil).
		Once()

	userDataChanger := NewMockDataChanger(t)
	const n = 31
	expectedBalanceDelta := -(n * containerInstance.Price)
	userDataChanger.
		On("ChangeBalance", mock.Anything, expectedBalanceDelta).
		Return(int64(50), nil).
		Once()

	wonItemsIndexes := make([]int, 0, n)
	for range n {
		wonItemsIndexes = append(wonItemsIndexes, rand.IntN(len(containerInstance.Items)))
	}
	itemPicker := NewMockItemPicker(t)
	itemPicker.
		On("Pick", mock.Anything, containerInstance.Items, n).
		Return(wonItemsIndexes, nil).
		Once()

	wonItems := make([]int64, 0, len(wonItemsIndexes))
	for _, i := range wonItemsIndexes {
		wonItems = append(wonItems, containerInstance.Items[i].ID)
	}
	inventoryItemsIDs := make([]int64, 0, len(wonItems))
	for range len(wonItems) {
		inventoryItemsIDs = append(inventoryItemsIDs, rand.Int64())
	}
	userDataChanger.
		On("AddItems", mock.Anything, wonItems).
		Return(inventoryItemsIDs, nil).
		Once()

	userStorage := NewMockUserStorage(t)
	userID := int64(904)
	userStorage.
		On("ChangeData", mock.Anything, userID, mock.Anything).
		Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			changeFunc := args.Get(2).(func(ctx context.Context, dc user.DataChanger) error)
			err := changeFunc(ctx, userDataChanger)
			require.NoError(t, err)
		}).
		Return(nil).
		Once()

	ctx := t.Context()
	gotOpeningResults, gotErr := New(pkgslog.DiscardLogger, containerStorage, userStorage, itemPicker).OpenContainers(ctx, userID, containerID, n)

	require.NoError(t, gotErr)
	require.Len(t, gotOpeningResults, n)
	for i, gotOpeningResult := range gotOpeningResults {
		require.Equal(t, container.OpeningResult{
			InventoryItemID: inventoryItemsIDs[i],
			ItemID:          wonItems[i],
		}, gotOpeningResult)
	}
}

func TestOpeningContainers_ErrorAdding(t *testing.T) {
	containerStorage := NewMockContainerStorage(t)
	containerID := int64(589)
	containerInstance := container.Container{
		ID:    containerID,
		Price: 100,
		Items: []item.WithWeight{
			{ID: 45669957, Weight: 9},
			{ID: 5501121, Weight: 2},
			{ID: 37675023, Weight: 3},
			{ID: 1510940, Weight: 4},
			{ID: 20735405, Weight: 5},
		},
	}
	containerStorage.
		On("FindContainer", mock.Anything, containerID).
		Return(containerInstance, nil).
		Once()

	userDataChanger := NewMockDataChanger(t)
	const n = 31
	expectedBalanceDelta := -(n * containerInstance.Price)
	userDataChanger.
		On("ChangeBalance", mock.Anything, expectedBalanceDelta).
		Return(int64(50), nil).
		Once()

	wonItemsIndexes := make([]int, 0, n)
	for range n {
		wonItemsIndexes = append(wonItemsIndexes, rand.IntN(len(containerInstance.Items)))
	}
	itemPicker := NewMockItemPicker(t)
	itemPicker.
		On("Pick", mock.Anything, containerInstance.Items, n).
		Return(wonItemsIndexes, nil).
		Once()

	wonItems := make([]int64, 0, len(wonItemsIndexes))
	for _, i := range wonItemsIndexes {
		wonItems = append(wonItems, containerInstance.Items[i].ID)
	}
	someErr := fmt.Errorf("some error")
	userDataChanger.
		On("AddItems", mock.Anything, wonItems).
		Return(nil, someErr).
		Once()

	userStorage := NewMockUserStorage(t)
	userID := int64(904)
	userStorage.
		On("ChangeData", mock.Anything, userID, mock.Anything).
		Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			changeFunc := args.Get(2).(func(ctx context.Context, dc user.DataChanger) error)
			err := changeFunc(ctx, userDataChanger)
			require.ErrorIs(t, err, someErr)
		}).
		Return(someErr).
		Once()

	ctx := t.Context()
	gotOpeningResults, gotErr := New(pkgslog.DiscardLogger, containerStorage, userStorage, itemPicker).OpenContainers(ctx, userID, containerID, n)

	require.ErrorIs(t, gotErr, someErr)
	require.Len(t, gotOpeningResults, 0)
	require.Equal(t, []container.OpeningResult(nil), gotOpeningResults)
}

func TestOpeningContainers_ErrorPicking(t *testing.T) {
	containerStorage := NewMockContainerStorage(t)
	containerID := int64(589)
	containerInstance := container.Container{
		ID:    containerID,
		Price: 100,
		Items: []item.WithWeight{
			{ID: 45669957, Weight: 9},
			{ID: 5501121, Weight: 2},
			{ID: 37675023, Weight: 3},
			{ID: 1510940, Weight: 4},
			{ID: 20735405, Weight: 5},
		},
	}
	containerStorage.
		On("FindContainer", mock.Anything, containerID).
		Return(containerInstance, nil).
		Once()

	userDataChanger := NewMockDataChanger(t)
	const n = 31
	expectedBalanceDelta := -(n * containerInstance.Price)
	userDataChanger.
		On("ChangeBalance", mock.Anything, expectedBalanceDelta).
		Return(int64(50), nil).
		Once()

	itemPicker := NewMockItemPicker(t)
	someErr := fmt.Errorf("some error")
	itemPicker.
		On("Pick", mock.Anything, containerInstance.Items, n).
		Return([]int(nil), someErr).
		Once()

	userStorage := NewMockUserStorage(t)
	userID := int64(904)
	userStorage.
		On("ChangeData", mock.Anything, userID, mock.Anything).
		Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			changeFunc := args.Get(2).(func(ctx context.Context, dc user.DataChanger) error)
			err := changeFunc(ctx, userDataChanger)
			require.ErrorIs(t, err, someErr)
		}).
		Return(someErr).
		Once()

	ctx := t.Context()
	gotOpeningResults, gotErr := New(pkgslog.DiscardLogger, containerStorage, userStorage, itemPicker).OpenContainers(ctx, userID, containerID, n)

	require.ErrorIs(t, gotErr, someErr)
	require.Len(t, gotOpeningResults, 0)
	require.Equal(t, []container.OpeningResult(nil), gotOpeningResults)
}

func TestOpeningContainers_ErrorChangingBalance(t *testing.T) {
	containerStorage := NewMockContainerStorage(t)
	containerID := int64(589)
	containerInstance := container.Container{
		ID:    containerID,
		Price: 100,
		Items: []item.WithWeight{
			{ID: 45669957, Weight: 9},
			{ID: 5501121, Weight: 2},
			{ID: 37675023, Weight: 3},
			{ID: 1510940, Weight: 4},
			{ID: 20735405, Weight: 5},
		},
	}
	containerStorage.
		On("FindContainer", mock.Anything, containerID).
		Return(containerInstance, nil).
		Once()

	userDataChanger := NewMockDataChanger(t)
	const n = 31
	expectedBalanceDelta := -(n * containerInstance.Price)
	someErr := fmt.Errorf("some error")
	userDataChanger.
		On("ChangeBalance", mock.Anything, expectedBalanceDelta).
		Return(int64(0), someErr).
		Once()

	userStorage := NewMockUserStorage(t)
	userID := int64(904)
	userStorage.
		On("ChangeData", mock.Anything, userID, mock.Anything).
		Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			changeFunc := args.Get(2).(func(ctx context.Context, dc user.DataChanger) error)
			err := changeFunc(ctx, userDataChanger)
			require.ErrorIs(t, err, someErr)
		}).
		Return(someErr).
		Once()

	ctx := t.Context()
	gotOpeningResults, gotErr := New(pkgslog.DiscardLogger, containerStorage, userStorage, NewMockItemPicker(t)).OpenContainers(ctx, userID, containerID, n)

	require.ErrorIs(t, gotErr, someErr)
	require.Len(t, gotOpeningResults, 0)
	require.Equal(t, []container.OpeningResult(nil), gotOpeningResults)
}

func TestOpeningContainers_ErrorChangingUserData(t *testing.T) {
	containerStorage := NewMockContainerStorage(t)
	containerID := int64(589)
	containerInstance := container.Container{
		ID:    containerID,
		Price: 100,
		Items: []item.WithWeight{
			{ID: 45669957, Weight: 9},
			{ID: 5501121, Weight: 2},
			{ID: 37675023, Weight: 3},
			{ID: 1510940, Weight: 4},
			{ID: 20735405, Weight: 5},
		},
	}
	containerStorage.
		On("FindContainer", mock.Anything, containerID).
		Return(containerInstance, nil).
		Once()

	userStorage := NewMockUserStorage(t)
	userID := int64(904)
	someErr := fmt.Errorf("some error")
	userStorage.
		On("ChangeData", mock.Anything, userID, mock.Anything).
		Return(someErr).
		Once()

	ctx := t.Context()
	const n = 32
	gotOpeningResults, gotErr := New(pkgslog.DiscardLogger, containerStorage, userStorage, NewMockItemPicker(t)).OpenContainers(ctx, userID, containerID, n)

	require.ErrorIs(t, gotErr, someErr)
	require.Len(t, gotOpeningResults, 0)
	require.Equal(t, []container.OpeningResult(nil), gotOpeningResults)
}

func TestOpeningContainers_ErrorFindingContainer(t *testing.T) {
	containerStorage := NewMockContainerStorage(t)
	containerID := int64(589)
	someErr := fmt.Errorf("some error")
	containerStorage.
		On("FindContainer", mock.Anything, containerID).
		Return(container.Container{}, someErr).
		Once()

	ctx := t.Context()
	const n = 32
	const userID = 12312
	gotOpeningResults, gotErr := New(pkgslog.DiscardLogger, containerStorage, NewMockUserStorage(t), NewMockItemPicker(t)).OpenContainers(ctx, userID, containerID, n)

	require.ErrorIs(t, gotErr, someErr)
	require.Len(t, gotOpeningResults, 0)
	require.Equal(t, []container.OpeningResult(nil), gotOpeningResults)
}

func TestCalculateOpeningPrice(t *testing.T) {
	tests := []struct {
		name           string
		n              int
		price          int64
		expectedPrice  int64
		expectedErrMsg string
	}{
		{
			name:           "NegPrice",
			n:              1,
			price:          -1,
			expectedErrMsg: "negative container price: -1",
		},
		{
			name:  "zeroPrice",
			n:     1,
			price: 0,
		},
		{
			name:           "Overflow",
			n:              2,
			price:          math.MaxInt64,
			expectedErrMsg: fmt.Sprintf("%d openings and %d container price overflow int64", 2, math.MaxInt64),
		},
		{
			name:          "MaxPossiblePrice",
			n:             1,
			price:         math.MaxInt64,
			expectedPrice: math.MaxInt64,
		},
		{
			name:          "MaxPossibleN",
			n:             100,
			price:         321,
			expectedPrice: 32100,
		},
		{
			name:          "Successful",
			n:             7,
			price:         321,
			expectedPrice: 2247,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotPrice, gotErr := calculateOpeningPrice(test.n, test.price)

			if gotErr != nil {
				if test.expectedErrMsg == "" {
					t.Fatalf("expected nil error, got non-nil %v", gotErr)
				}
				if msg := gotErr.Error(); msg != test.expectedErrMsg {
					t.Fatalf("expected error msg: %v\ngot: %v", test.expectedErrMsg, msg)
				}
			} else if test.expectedErrMsg != "" {
				t.Fatalf("expected non-nil error, got nil")
			}

			if test.expectedPrice != gotPrice {
				t.Fatalf("expected price %v, got %v", test.expectedPrice, gotPrice)
			}
		})
	}
}
