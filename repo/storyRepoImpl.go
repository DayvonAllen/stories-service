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
)

type StoryRepoImpl struct {
	Story domain.Story
	StoryDto domain.StoryDto
	StoryList []domain.Story
	StoryDtoList []domain.StoryDto
}

func (s StoryRepoImpl) Create(story *domain.Story) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.StoriesCollection.InsertOne(context.TODO(), &story)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (s StoryRepoImpl) UpdateByID(id primitive.ObjectID, newContent string) (*domain.Story, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"content", newContent}}}}

	err := conn.StoriesCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&s.Story)

	if err != nil {
		return nil, err
	}

	return &s.Story, nil
}

func (s StoryRepoImpl) FindAll() (*[]domain.StoryDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	//findOptions := options.FindOptions{}
	//perPage := 10
	//pageNumber, err := strconv.Atoi(page)

	//if err != nil {
	//	return nil, fmt.Errorf("page must be a number")
	//}
	//findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	//findOptions.SetLimit(int64(perPage))

	// Get all tags
	cur, err := conn.StoriesCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &s.StoryDtoList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &s.StoryDtoList, nil
}

func (s StoryRepoImpl) FindById(storyID primitive.ObjectID) (*domain.StoryDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.StoriesCollection.FindOne(context.TODO(), bson.D{{"_id", storyID}}).Decode(&s.StoryDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return  nil, fmt.Errorf("error processing data")
	}

	return &s.StoryDto, nil
}

func (s StoryRepoImpl) DeleteByID(id primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.StoriesCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func NewStoryRepoImpl() StoryRepoImpl {
	var storyRepoImpl StoryRepoImpl

	return storyRepoImpl
}

