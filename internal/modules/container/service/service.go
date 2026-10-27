package service

import (
	"context"
	"log/slog"

	"github.com/ulibaysya/caselya/internal/modules/container"
	"github.com/ulibaysya/caselya/internal/modules/item"
	"github.com/ulibaysya/caselya/internal/modules/user"
)

type ContainerStorage interface {
	FindContainer(ctx context.Context, id int64) (container.Container, error)
}

type UserStorage interface {
	ChangeData(ctx context.Context, id int64, changeFunc func(ctx context.Context, dc user.DataChanger) error) error
}

type ItemPicker interface {
	Pick(ctx context.Context, items []item.WithWeight, count int) ([]int, error)
}

type Service struct {
	logger           *slog.Logger
	containerStorage ContainerStorage
	userStorage      UserStorage
	itemPicker       ItemPicker
}

func New(logger *slog.Logger, containerStorage ContainerStorage, userStorage UserStorage, itemPicker ItemPicker) *Service {
	return &Service{
		logger:           logger,
		containerStorage: containerStorage,
		userStorage:      userStorage,
		itemPicker:       itemPicker,
	}
}
