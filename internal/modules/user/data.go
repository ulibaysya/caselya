package user

import "context"

type DataChanger interface {
	ChangeBalance(ctx context.Context, delta int64) (int64, error)
	AddItems(ctx context.Context, ids []int64) ([]int64, error)
}
