package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ulibaysya/caselya/internal/modules/user"
)

func (s *Storage) ChangeBalance(ctx context.Context, id int64, delta int64) (int64, error) {
	const sql = `
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance;
	`

	var balanceAfterChange int64
	err := s.db.QueryRow(ctx, sql, delta, id).Scan(&balanceAfterChange)
	if err != nil {
		if isInsufficientBalanceError(err) {
			return 0, user.ErrInsufficientBalance
		}
		return 0, fmt.Errorf("updating user balance: %w", err)
	}

	return balanceAfterChange, nil
}

func isInsufficientBalanceError(err error) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}
	if pgErr.Code != pgerrcode.CheckViolation {
		return false
	}
	const checkConstraintName = "users_balance_check"
	return pgErr.ConstraintName == checkConstraintName
}
