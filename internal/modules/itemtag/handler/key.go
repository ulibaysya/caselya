package handler

import (
	"context"
	"errors"

	v1 "github.com/ulibaysya/caselya/api/caselya/hub/v1"
	"github.com/ulibaysya/caselya/internal/modules/itemtag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateKey(ctx context.Context, req *v1.CreateKeyRequest) (*v1.CreateKeyResponse, error) {
	if err := h.srvc.CreateKey(ctx, req.GetName(), req.GetPublic(), req.GetValues()); err != nil {
		if errors.Is(err, itemtag.ErrKeyAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.CreateKeyResponse{}, nil
}
