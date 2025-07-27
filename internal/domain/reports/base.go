package reports

import (
	"billing/internal/mng"
	"billing/internal/psql"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const mnCollection = "reports"

type Reports struct {
	pgPool     *psql.Pool
	mnDatabase *mng.Database
}

func New(pgPool *psql.Pool, mnDatabase *mng.Database) *Reports {
	return &Reports{pgPool: pgPool, mnDatabase: mnDatabase}
}

type MnReport struct {
	Id                bson.ObjectID `bson:"_id"`
	ReportId          string        `bson:"report_id"`
	UserId            int           `bson:"user_id"`
	ClientGeneratedId string        `bson:"client_generated_id"`
	IsPurchased       bool          `bson:"is_purchased"`
}
