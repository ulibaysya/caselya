package handler

import (
	"context"

	v1 "github.com/ulibaysya/caselya/api/caselya/hub/v1"
	"google.golang.org/grpc"
)

type Service interface {
	CreateKey(ctx context.Context, name string, public bool, values []string) error
}

type Handler struct {
	srvc Service

	v1.UnimplementedItemTagServiceServer
}

func Register(s grpc.ServiceRegistrar, srvc Service) {
	v1.RegisterItemTagServiceServer(s, &Handler{
		srvc: srvc,
	})
}
