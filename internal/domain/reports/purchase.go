package reports

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	defaultCost = 100

	constraintViolationUniqueCode = "23505"
	constraintViolationCheckCode  = "23514"
	negativeBalanceConstraintName = "balance_never_negative"

	financialOperationInitState = "frozen"
	financialOperationPaidState = "paid"
	financialOperationRetState  = "returned"

	mnCollection = "reports"
)

type MnReport struct {
	Id                bson.ObjectID `bson:"_id"`
	ReportId          string        `bson:"report_id"`
	UserId            int           `bson:"user_id"`
	ClientGeneratedId string        `bson:"client_generated_id"`
	IsPurchased       bool          `bson:"is_purchased"`
}

func (r *Reports) Purchase(reportId, userId string) (bool, bool, error) {

	paid, err := r.alreadyPaid(reportId, userId)
	if err != nil {
		return false, false, err
	}
	if paid {
		return true, false, nil
	}

	err = r.createPayment(reportId, userId)
	if err == nil {
		err = r.markPurchased(reportId, userId)
		if err == nil {
			err = r.markPaid(reportId, userId)
			if err != nil {
				return false, false, err
			}

			return true, false, nil
		}

		err = r.returnPayment(reportId, userId)
		if err != nil {
			return false, false, err
		}

		return false, false, nil
	}
	if err != nil && negativeBalanceError(err) {
		return false, true, err
	}
	if err != nil && !duplicateTupleError(err) {
		return false, false, err
	}

	err = r.rewritePayment(reportId, userId)
	if err == nil {
		err = r.markPurchased(reportId, userId)
		if err == nil {
			err = r.markPaid(reportId, userId)
			if err != nil {
				return false, false, err
			}

			return true, false, nil
		}

		err = r.returnPayment(reportId, userId)
		if err != nil {
			return false, false, err
		}

		return false, false, nil
	}
	if err != nil && negativeBalanceError(err) {
		return false, true, err
	}

	return false, false, err
}

func (r *Reports) alreadyPaid(reportId, userId string) (bool, error) {

	fq := `SELECT EXISTS (SELECT 1 FROM financial_operations WHERE report_id = $1 AND user_id = $2 AND state IN ($3, $4))`

	var b bool
	err := r.pgPool.QueryRow(context.Background(), fq, reportId, userId, financialOperationInitState, financialOperationPaidState).Scan(&b)
	if err != nil {
		return false, err
	}

	return b, nil
}

func (r *Reports) createPayment(reportId, userId string) error {
	tx, err := r.pgPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	uq := `UPDATE users SET balance = balance - $1 WHERE id = $2`

	_, err = tx.Exec(context.Background(), uq, defaultCost, userId)
	if err != nil {
		return err
	}

	fq := `INSERT INTO financial_operations (report_id, user_id, cost, state) VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(context.Background(), fq, reportId, userId, defaultCost, financialOperationInitState)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Reports) rewritePayment(reportId, userId string) error {
	tx, err := r.pgPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	uq := `UPDATE users SET balance = balance - $1 WHERE id = $2`

	_, err = tx.Exec(context.Background(), uq, defaultCost, userId)
	if err != nil {
		return err
	}

	fuq := `UPDATE financial_operations SET cost = $1, state = $2 WHERE report_id = $3 AND user_id = $4 AND state <> $2`

	_, err = tx.Exec(context.Background(), fuq, defaultCost, financialOperationInitState, reportId, userId)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Reports) returnPayment(reportId, userId string) error {
	tx, err := r.pgPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	// NOTE: возвращаю фиксированную стоимость
	uq := `UPDATE users SET balance = balance + $1 WHERE id = $2`

	_, err = tx.Exec(context.Background(), uq, defaultCost, userId)
	if err != nil {
		return err
	}

	fuq := `UPDATE financial_operations SET state = $1 WHERE report_id = $2 AND user_id = $3 AND state <> $1`

	_, err = tx.Exec(context.Background(), fuq, financialOperationRetState, reportId, userId)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Reports) markPaid(reportId, userId string) error {
	fuq := `UPDATE financial_operations SET state = $1 WHERE report_id = $2 AND user_id = $3 AND state = $4`

	_, err := r.pgPool.Exec(context.Background(), fuq, financialOperationPaidState, reportId, userId, financialOperationInitState)
	if err != nil {
		return err
	}

	return nil
}

func (r *Reports) markPurchased(reportId, userId string) error {
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return err
	}

	var result MnReport

	opts := options.FindOneAndUpdate().SetUpsert(false)
	filter := bson.M{"report_id": reportId, "user_id": uid}
	update := bson.M{"$set": bson.M{"is_purchased": true}}

	err = r.mnDatabase.Collection(mnCollection).FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}

func negativeBalanceError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}

	if pgErr.Code == constraintViolationCheckCode && pgErr.ConstraintName == negativeBalanceConstraintName {
		return true
	}

	return false
}

func duplicateTupleError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}

	if pgErr.Code == constraintViolationUniqueCode {
		return true
	}

	return false
}
