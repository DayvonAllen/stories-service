package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"example.com/app/helper"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"sync"

	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepoImpl struct {
	user        domain.User
	currentUser domain.CurrentUserProfile
	viewedUser  domain.ViewUserProfile
}

func (u UserRepoImpl) Create(user *domain.User) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.UserCollection.Find(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"username": user.Username},
		},
	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		_, err = conn.UserCollection.InsertOne(context.TODO(), &user)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		return nil
	}

	return fmt.Errorf("user already exists")
}

func (u UserRepoImpl) GetCurrentUserProfile(username string) (*domain.CurrentUserProfile, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.currentUser)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	stories, err := StoryRepoImpl{}.FindAllByUsername(username)

	if err != nil {
		return nil, err
	}

	u.currentUser.Posts = *stories

	return &u.currentUser, nil
}

func (u UserRepoImpl) GetUserProfile(username string) (*domain.ViewUserProfile, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.viewedUser)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	if u.viewedUser.ProfileIsViewable == false {
		return nil, fmt.Errorf("cannot view user")
	}

	if !u.viewedUser.DisplayFollowerCount {
		u.viewedUser.FollowerCount = -1
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		stories, err := StoryRepoImpl{}.FindAllByUsername(username)

		if err != nil {
			panic(err)
		}

		u.viewedUser.Posts = *stories

		return
	}()

	go func() {
		defer wg.Done()
		u.viewedUser.IsFollowing = helper.CurrentUserInteraction(u.viewedUser.Followers, username)
		return
	}()

	wg.Wait()
	return &u.viewedUser, nil
}

func (u UserRepoImpl) UpdateByID(user *domain.User) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", user}}

	conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}

func (u UserRepoImpl) DeleteByID(user *domain.User) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := conn.StartSession()

	if err != nil {
		panic(err)
	}

	defer session.EndSession(context.Background())

	// execute this code in a logical transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			_, err := conn.UserCollection.DeleteOne(context.TODO(), bson.D{{"_id", user.Id}})

			if err != nil {
				panic(err)
			}

			return
		}()

		go func() {
			defer wg.Done()
			_, err = conn.StoryCollection.DeleteMany(context.TODO(), bson.D{{"authorUsername", user.Username}})

			if err != nil {
				panic(err)

			}
			return
		}()

		wg.Wait()

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to delete user")
	}

	return nil
}

func NewUserRepoImpl() UserRepoImpl {
	var userRepoImpl UserRepoImpl

	return userRepoImpl
}
