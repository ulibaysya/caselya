package storage

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/require"
	"github.com/ulibaysya/caselya/internal/modules/user"
)

func TestUserDataChange_Success(t *testing.T) {
	db, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	db.ExpectBegin().Times(1)
	db.ExpectCommit().Times(1)
	db.ExpectRollback().Times(1)

	ctx := t.Context()
	const id = int64(17289)
	gotErr := New(db).ChangeData(ctx, id, func(ctx context.Context, dc user.DataChanger) error {
		return nil
	})

	require.NoError(t, gotErr)
}
