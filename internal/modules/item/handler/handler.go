package handler

import (
	"context"

	v1 "github.com/ulibaysya/caselya/api/caselya/hub/v1"
	"github.com/ulibaysya/caselya/internal/modules/item"
	"google.golang.org/grpc"
)

type Service interface {
	CreateItem(ctx context.Context, name string, price int64, tags []item.Tag) (int64, error)
}

type Handler struct {
	srvc Service

	v1.UnimplementedItemServiceServer
}

func Register(s grpc.ServiceRegistrar, srvc Service) {
	v1.RegisterItemServiceServer(s, &Handler{
		srvc: srvc,
	})
}
