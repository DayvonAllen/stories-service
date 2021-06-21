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

type StoryRepoImpl struct {
	Story domain.Story
	StoryDto domain.StoryDto
	StoryList []domain.Story
	StoryDtoList []domain.StoryDto
}

func (s StoryRepoImpl) Create(story *domain.Story) error {
	_, err := database.GetInstance().StoriesCollection.InsertOne(context.TODO(), &story)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (s StoryRepoImpl) UpdateByID(id primitive.ObjectID, newContent string) (*domain.Story, error) {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"content", newContent}}}}

	err := database.GetInstance().StoriesCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&s.Story)

	if err != nil {
		return nil, err
	}

	return &s.Story, nil
}

func (s StoryRepoImpl) FindAll() (*[]domain.Story, error) {
	// Get all tags
	cur, err := database.GetInstance().StoriesCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.Story
		err = cur.Decode(&elem)

		if err != nil {
			return nil, fmt.Errorf("error processing data")
		}

		s.StoryList = append(s.StoryList, elem)
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

	return &s.StoryList, nil
}

func (s StoryRepoImpl) FindById(storyID primitive.ObjectID) (*domain.Story, error) {
	err := database.GetInstance().StoriesCollection.FindOne(context.TODO(), bson.D{{"_id", storyID}}).Decode(&s.Story)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return  nil, fmt.Errorf("error processing data")
	}

	return &s.Story, nil
}

func (s StoryRepoImpl) DeleteByID(id primitive.ObjectID) error {
	_, err := database.GetInstance().StoriesCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func NewStoryRepoImpl() StoryRepoImpl {
	var storyRepoImpl StoryRepoImpl

	return storyRepoImpl
}

