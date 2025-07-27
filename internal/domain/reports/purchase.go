package reports

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	defaultCost                   = 100
	constraintViolationSQLCode    = "23514"
	negativeBalanceConstraintName = "balance_never_negative"
)

func (r *Reports) Purchase(reportId, userId string) (bool, bool, error) {
	paid, err := r.alreadyPaid(reportId, userId)
	if err != nil {
		return false, false, err
	}
	if paid {
		return true, false, nil
	}

	tx, err := r.pgPool.Begin(context.Background())
	if err != nil {
		return false, false, err
	}

	defer tx.Rollback(context.Background())

	uq := `UPDATE users SET balance = balance - $1 WHERE id = $2`

	_, err = tx.Exec(context.Background(), uq, defaultCost, userId)
	if err != nil && negativeBalanceError(err) {
		return false, true, err
	}
	if err != nil && !negativeBalanceError(err) {
		return false, false, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return false, false, err
	}

	return true, false, nil
}

func (r *Reports) alreadyPaid(reportId, userId string) (bool, error) {
	fq := `SELECT 1 FROM financial_operations WHERE report_id = $1 AND user_id = $2 AND state IN ('frozen', 'paid')`

	err := r.pgPool.QueryRow(context.Background(), fq, reportId, userId).Scan()
	if err != nil && err != pgx.ErrNoRows {
		return false, err
	}
	if err != nil && err == pgx.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func negativeBalanceError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}

	if pgErr.Code == constraintViolationSQLCode && pgErr.ConstraintName == negativeBalanceConstraintName {
		return true
	}

	return false
}
