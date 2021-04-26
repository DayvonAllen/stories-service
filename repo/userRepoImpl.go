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

var dbConnection = database.GetInstance()

func (u UserRepoImpl) Create(user *domain.User) error {
	cur, err := dbConnection.Collection("users").Find(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"username": user.Username},
		},

	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		_, err = dbConnection.Collection("users").InsertOne(context.TODO(), &user)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		return nil
	}

	return fmt.Errorf("user already exists")
}

func (u UserRepoImpl) FindByUsername(username string) (*domain.User, error) {
	err := dbConnection.Collection("users").FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.user)

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

	database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}