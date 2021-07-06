package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type CommentRepoImpl struct {
	Comment domain.Comment
	CommentList []domain.Comment
	CommentDtoList []domain.CommentDto
}

func (c CommentRepoImpl) Create(comment *domain.Comment) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.CommentsCollection.InsertOne(context.TODO(), &comment)

	if err != nil {
		return err
	}

	return  nil
}

func (c CommentRepoImpl) FindById(storyID primitive.ObjectID) (*domain.Comment, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.CommentsCollection.FindOne(context.TODO(), bson.D{{"_id", storyID}}).Decode(&c.Comment)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return  nil, fmt.Errorf("error processing data")
	}

	return &c.Comment, nil
}

func (c CommentRepoImpl) FindAllCommentsByStoryId(storyID primitive.ObjectID) (*[]domain.CommentDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.CommentsCollection.Find(context.TODO(), bson.D{{"storyId", storyID}})

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return  nil, fmt.Errorf("error processing data")
	}

	if err = cur.All(context.TODO(), &c.CommentDtoList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &c.CommentDtoList, nil
}

func (c CommentRepoImpl) UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time) (*domain.Comment, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"content", newContent}, {"edited", edited},
		{"updatedTime", updatedTime}}}}

	err := conn.CommentsCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&c.Comment)

	if err != nil {
		return nil, err
	}

	return &c.Comment, nil
}

func (c CommentRepoImpl) DeleteById(id primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.CommentsCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func NewCommentRepoImpl() CommentRepoImpl {
	var commentRepoImpl CommentRepoImpl

	return commentRepoImpl
}
