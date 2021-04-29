package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagRepoImpl struct {
	Tag domain.Tag
	TagList []domain.Tag
}

func (t TagRepoImpl) Create(tag *domain.Tag) error {
	_, err := database.GetInstance().TagsCollection.InsertOne(context.TODO(), &tag)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (t TagRepoImpl) CreateMany(tags *[]interface{}) error {

	_, err := database.GetInstance().TagsCollection.InsertMany(context.TODO(), *tags)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (t TagRepoImpl) FindByTagName(tagName string) (*domain.Tag, error) {

	err := database.GetInstance().TagsCollection.FindOne(context.TODO(), bson.D{{"tagName", tagName}}).Decode(&t.Tag)

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &t.Tag, nil
}

func (t TagRepoImpl) FindAll() (*[]domain.Tag, error) {
	// Get all tags
	cur, err := database.GetInstance().TagsCollection.Find(context.TODO(), bson.M{})

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

func (t TagRepoImpl) DeleteByID(id primitive.ObjectID) error {
	_, err := database.GetInstance().TagsCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func NewTagRepoImpl() TagRepoImpl {
	var tagRepoImpl TagRepoImpl

	return tagRepoImpl
}