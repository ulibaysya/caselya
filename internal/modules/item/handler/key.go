package handler

import (
	"context"

	v1 "github.com/ulibaysya/caselya/api/caselya/hub/v1"
	"github.com/ulibaysya/caselya/internal/modules/item"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateItem(ctx context.Context, req *v1.CreateItemRequest) (*v1.CreateItemResponse, error) {
	tags := make([]item.Tag, 0, len(req.GetTags()))
	for _, reqTag := range req.GetTags() {
		tags = append(tags, item.Tag{
			Key:    reqTag.Key,
			Values: reqTag.GetValues(),
		})
	}
	id, err := h.srvc.CreateItem(ctx, req.GetName(), req.GetPrice(), tags)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.CreateItemResponse{
		Id: id,
	}, nil
}
