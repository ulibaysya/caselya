package container

import "github.com/ulibaysya/caselya/internal/modules/item"

type Container struct {
	ID    int64
	Price int64
	Items []item.WithWeight
}
