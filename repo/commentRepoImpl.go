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
)

type CommentRepoImpl struct {
	Comment domain.Comment
	CommentList []domain.Comment
}

func (c CommentRepoImpl) Create(comment *domain.Comment) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	story := new(domain.Story)
	err := conn.StoriesCollection.FindOne(context.TODO(), bson.D{{"_id", comment.StoryId}}).Decode(story)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("story not found")
		}
		return err
	}

	filter := bson.D{{"_id", story.Id}}
	update := bson.M{"$push": bson.M{"comments": comment.Id}}

	_, err = conn.StoriesCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return err
	}

	_, err = conn.CommentsCollection.InsertOne(context.TODO(), &comment)

	if err != nil {
		return err
	}

	return  nil
}

func (c CommentRepoImpl) FindById(commentID primitive.ObjectID) (*domain.Comment, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.CommentsCollection.FindOne(context.TODO(), bson.D{{"_id", commentID}}).Decode(&c.Comment)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return  nil, fmt.Errorf("error processing data")
	}

	return &c.Comment, nil
}

func (c CommentRepoImpl) FindAllCommentsByStoryId(storyID primitive.ObjectID) (*[]domain.Comment, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	story := new(domain.Story)
	err := conn.StoriesCollection.FindOne(context.TODO(), bson.D{{"_id", storyID}}).Decode(story)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return  nil, fmt.Errorf("error processing data")
	}

	query := bson.M{"_id": bson.M{"$in": story.Comments}}

	cur, err := conn.CommentsCollection.Find(context.TODO(), query)

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.Comment
		err = cur.Decode(&elem)

		if err != nil {
			return nil, fmt.Errorf("error processing data")
		}

		c.CommentList = append(c.CommentList, elem)
	}

	err = cur.Err()

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &c.CommentList, nil
}

func (c CommentRepoImpl) UpdateById(id primitive.ObjectID, newContent string) (*domain.Comment, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"content", newContent}}}}

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
