package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepoImpl struct {
	user domain.User
}

func (u UserRepoImpl) Create(user *domain.User) error {
	cur, err := database.GetInstance().UserCollection.Find(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"username": user.Username},
		},

	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		_, err = database.GetInstance().UserCollection.InsertOne(context.TODO(), &user)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		return nil
	}

	return fmt.Errorf("user already exists")
}

func (u UserRepoImpl) FindByUsername(username string) (*domain.User, error) {
	err := database.GetInstance().UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.user)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return  nil, fmt.Errorf("error processing data")
	}

	return &u.user, nil
}

func (u UserRepoImpl) UpdateByID(user *domain.User) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", user}}

	database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}