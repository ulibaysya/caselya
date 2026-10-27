package storage

import (
	"testing"

	"github.com/ulibaysya/caselya/internal/modules/item"
)

func TestBuildTagsInsertSQL(t *testing.T) {
	tags := []item.Tag{
		{
			Key:    "some-key",
			Values: []string{"val1", "val2"},
		},
		{
			Key:    "another-key",
			Values: []string{"val3", "val4"},
		},
		{
			Key:    "third-key",
			Values: []string{"val5", "val6"},
		},
	}

	sql, args := buildTagsInsertSQL(12, tags)
	t.Logf("%+v \n %+v", sql, args)
}
