package reports

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Page []MnReport

func (r *Reports) List(userId string, limit, offset int64) (*Page, error) {
	var page Page

	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return &page, err
	}

	filter := bson.M{"user_id": uid}
	opts := options.Find().SetSort(bson.M{"_id": 1}).SetLimit(limit).SetSkip(offset)

	cursor, err := r.mnDatabase.Collection(mnCollection).Find(context.Background(), filter, opts)
	if err != nil {
		return &page, err
	}

	for cursor.Next(context.Background()) {
		var r MnReport
		err = cursor.Decode(&r)
		if err != nil {
			return &page, err
		}

		page = append(page, r)
	}

	return &page, nil
}
