package storage

import (
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/require"
	"github.com/ulibaysya/caselya/internal/modules/container"
	"github.com/ulibaysya/caselya/internal/modules/item"
)

func TestFindContainer_Success(t *testing.T) {
	db, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	containerInstance := container.Container{
		ID:    12,
		Price: 12345,
		Items: []item.WithWeight{
			{
				ID:     1753400,
				Weight: 38028623,
			},
			{
				ID:     33662462,
				Weight: 912811,
			},
			{
				ID:     40976675,
				Weight: 4819904,
			},
			{
				ID:     24568794,
				Weight: 42440173,
			},
			{
				ID:     44054355,
				Weight: 339763146,
			},
		},
	}

	const id = int64(123890)
	rows := pgxmock.NewRows([]string{"id", "price", "items_id", "weight"}).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[0].ID, containerInstance.Items[0].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[1].ID, containerInstance.Items[1].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[2].ID, containerInstance.Items[2].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[3].ID, containerInstance.Items[3].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[4].ID, containerInstance.Items[4].Weight)
	const sql = `
		SELECT c.id, c.price, ci.items_id, ci.weight
		FROM containers AS c
		JOIN containers_items AS ci
		ON c.id = ci.containers_id
		WHERE c.id = $1
		ORDER BY ci.weight, ci.items_id;
	`
	db.ExpectQuery(sql).
		WillReturnRows(rows).
		WithArgs(id).
		RowsWillBeClosed().
		Times(1)

	ctx := t.Context()
	gotContainer, gotErr := New(db).FindContainer(ctx, id)

	require.NoError(t, db.ExpectationsWereMet(), "db expectations")
	require.NoError(t, gotErr)
	require.Equal(t, containerInstance, gotContainer)
}

func TestFindContainer_FailNotFound(t *testing.T) {
	db, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	const id = int64(123890)
	rows := pgxmock.NewRows([]string{"id", "price", "items_id", "weight"})
	const sql = `
		SELECT c.id, c.price, ci.items_id, ci.weight
		FROM containers AS c
		JOIN containers_items AS ci
		ON c.id = ci.containers_id
		WHERE c.id = $1
		ORDER BY ci.weight, ci.items_id;
	`
	db.ExpectQuery(sql).
		WillReturnRows(rows).
		WithArgs(id).
		RowsWillBeClosed().
		Times(1)

	ctx := t.Context()
	gotContainer, gotErr := New(db).FindContainer(ctx, id)

	require.NoError(t, db.ExpectationsWereMet(), "db expectations")
	require.ErrorIs(t, gotErr, container.ErrNotFound)
	require.Equal(t, container.Container{}, gotContainer)
}

func TestFindContainer_FailScanning(t *testing.T) {
	db, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	containerInstance := container.Container{
		ID:    12,
		Price: 12345,
		Items: []item.WithWeight{
			{
				ID:     1753400,
				Weight: 38028623,
			},
			{
				ID:     33662462,
				Weight: 912811,
			},
			{
				ID:     40976675,
				Weight: 4819904,
			},
			{
				ID:     24568794,
				Weight: 42440173,
			},
			{
				ID:     44054355,
				Weight: 339763146,
			},
		},
	}

	const id = int64(123890)
	someErr := errors.New("some error")
	rows := pgxmock.NewRows([]string{"id", "price", "items_id", "weight"}).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[0].ID, containerInstance.Items[0].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[1].ID, containerInstance.Items[1].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[2].ID, containerInstance.Items[2].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[3].ID, containerInstance.Items[3].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[4].ID, containerInstance.Items[4].Weight).
		RowError(3, someErr)
	const sql = `
		SELECT c.id, c.price, ci.items_id, ci.weight
		FROM containers AS c
		JOIN containers_items AS ci
		ON c.id = ci.containers_id
		WHERE c.id = $1
		ORDER BY ci.weight, ci.items_id;
	`
	db.ExpectQuery(sql).
		WillReturnRows(rows).
		WithArgs(id).
		RowsWillBeClosed().
		Times(1)

	ctx := t.Context()
	gotContainer, gotErr := New(db).FindContainer(ctx, id)

	require.NoError(t, db.ExpectationsWereMet(), "db expectations")
	require.ErrorIs(t, gotErr, someErr)
	require.Equal(t, container.Container{}, gotContainer)
}

func TestFindContainer_FailClosingRows(t *testing.T) {
	db, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	containerInstance := container.Container{
		ID:    12,
		Price: 12345,
		Items: []item.WithWeight{
			{
				ID:     1753400,
				Weight: 38028623,
			},
			{
				ID:     33662462,
				Weight: 912811,
			},
			{
				ID:     40976675,
				Weight: 4819904,
			},
			{
				ID:     24568794,
				Weight: 42440173,
			},
			{
				ID:     44054355,
				Weight: 339763146,
			},
		},
	}

	const id = int64(123890)
	someErr := errors.New("some error")
	rows := pgxmock.NewRows([]string{"id", "price", "items_id", "weight"}).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[0].ID, containerInstance.Items[0].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[1].ID, containerInstance.Items[1].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[2].ID, containerInstance.Items[2].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[3].ID, containerInstance.Items[3].Weight).
		AddRow(containerInstance.ID, containerInstance.Price, containerInstance.Items[4].ID, containerInstance.Items[4].Weight).
		CloseError(someErr)
	const sql = `
		SELECT c.id, c.price, ci.items_id, ci.weight
		FROM containers AS c
		JOIN containers_items AS ci
		ON c.id = ci.containers_id
		WHERE c.id = $1
		ORDER BY ci.weight, ci.items_id;
	`
	db.ExpectQuery(sql).
		WillReturnRows(rows).
		WithArgs(id).
		RowsWillBeClosed().
		Times(1)

	ctx := t.Context()
	gotContainer, gotErr := New(db).FindContainer(ctx, id)

	require.NoError(t, db.ExpectationsWereMet(), "db expectations")
	require.ErrorIs(t, gotErr, someErr)
	require.Equal(t, container.Container{}, gotContainer)
}


func TestFindContainer_FailQuerying(t *testing.T) {
	db, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	const id = int64(123890)
	someErr := errors.New("some error")
	const sql = `
		SELECT c.id, c.price, ci.items_id, ci.weight
		FROM containers AS c
		JOIN containers_items AS ci
		ON c.id = ci.containers_id
		WHERE c.id = $1
		ORDER BY ci.weight, ci.items_id;
	`
	db.ExpectQuery(sql).
		WillReturnError(someErr).
		WithArgs(id).
		Times(1)

	ctx := t.Context()
	gotContainer, gotErr := New(db).FindContainer(ctx, id)

	require.NoError(t, db.ExpectationsWereMet(), "db expectations")
	require.ErrorIs(t, gotErr, someErr)
	require.Equal(t, container.Container{}, gotContainer)
}
