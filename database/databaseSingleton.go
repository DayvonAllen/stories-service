package database

import (
	"context"
	"example.com/app/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type Connection struct {
	*mongo.Client
	UserCollection *mongo.Collection
	StoriesCollection *mongo.Collection
	TagsCollection *mongo.Collection
	CommentsCollection *mongo.Collection
	LikesCollection *mongo.Collection
	DislikesCollection *mongo.Collection
	*mongo.Database
}

var dbConnection *Connection
var once sync.Once

// GetInstance creates one instance and always returns that one instance
func GetInstance() *Connection {
	// only executes this once
	once.Do(func() {
		err := connectToDB()
		if err != nil {
			panic(err)
		}
	})
	return dbConnection
}

func connectToDB() error {
	p := config.Config("DB_PORT")
	n := config.Config("DB_NAME")
	h := config.Config("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(n + h + p))
	if err != nil { return err }

	// create database
	db := client.Database("stories-service")

	// create collection
	userCollection := db.Collection("users")
	storiesCollection := db.Collection("stories")
	tagsCollection := db.Collection("tags")
	commentsCollection := db.Collection("comments")
	likesCollection := db.Collection("likes")
	dislikesCollection := db.Collection("dislikes")

	dbConnection = &Connection{client, userCollection, storiesCollection,
		tagsCollection, commentsCollection,
		likesCollection, dislikesCollection,
		db}

	return nil
}