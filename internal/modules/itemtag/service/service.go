package service

import (
	"context"
	"log/slog"
)

type Storage interface {
	CreateKey(ctx context.Context, name string, public bool, values []string) error
}

type Service struct {
	logger *slog.Logger
	strg   Storage
}

func New(logger *slog.Logger, strg Storage) *Service {
	return &Service{
		logger: logger,
		strg:   strg,
	}
}
