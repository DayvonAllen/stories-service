package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TagRepoImpl struct {
	Tag domain.Tag
	TagList []domain.Tag
}

func (t TagRepoImpl) Create(tag *domain.Tag) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.TagsCollection.FindOne(context.TODO(), bson.D{{"tagName", tag.TagName}}).Decode(&t.Tag)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("not found")
		}
		return err
	}
	_, err = conn.TagsCollection.InsertOne(context.TODO(), &tag)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (t TagRepoImpl) CreateMany(tags *[]interface{}) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.TagsCollection.InsertMany(context.TODO(), *tags)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (t TagRepoImpl) FindByTagName(tagName string) (*domain.Tag, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.TagsCollection.FindOne(context.TODO(), bson.D{{"tagName", tagName}}).Decode(&t.Tag)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("not found")
		}
		return nil, err
	}

	return &t.Tag, nil
}

func (t TagRepoImpl) FindAll() (*[]domain.Tag, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	// Get all tags
	cur, err := conn.TagsCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.Tag
		err = cur.Decode(&elem)

		if err != nil {
			return nil, fmt.Errorf("error processing data")
		}

		t.TagList = append(t.TagList, elem)
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

	return &t.TagList, nil
}

func (t TagRepoImpl) DeleteById(id primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.TagsCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func NewTagRepoImpl() TagRepoImpl {
	var tagRepoImpl TagRepoImpl

	return tagRepoImpl
}