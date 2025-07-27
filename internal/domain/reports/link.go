package reports

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Reports) LinkAnonymous(clientGeneratedId, userId string) error {
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return err
	}

	filter := bson.M{"client_generated_id": clientGeneratedId, "user_id": nil}
	update := bson.M{"$set": bson.M{"user_id": uid}}

	_, err = r.mnDatabase.Collection(mnCollection).UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
