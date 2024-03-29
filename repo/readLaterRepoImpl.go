package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type ReadLaterRepoImpl struct {
	ReadLater             domain.ReadLater
	ReadLaterList             []domain.ReadLater
	ReadLaterDto          domain.ReadLaterDto
}

func (r ReadLaterRepoImpl) Create(username string, storyId primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	story := new(domain.StoryDto)
	err := conn.StoryCollection.FindOne(context.TODO(), bson.D{{"_id", storyId}}).Decode(&story)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("error processing data")
	}

	err = conn.ReadLaterCollection.FindOne(context.TODO(), bson.D{{"story._id", story.Id}, {"username", username}}).Decode(&story)

	if err != nil {
		r.ReadLater.CreatedAt = time.Now()
		r.ReadLater.UpdatedAt = time.Now()
		r.ReadLater.Username = username
		r.ReadLater.Story = *story
		r.ReadLater.Id = primitive.NewObjectID()

		_, err = conn.ReadLaterCollection.InsertOne(context.TODO(), &r.ReadLater)

		if err != nil {
			return fmt.Errorf("error processing data")
		}
		return nil
	}

	go func() {
		event := new(domain.Event)
		event.Action = "add to read-later"
		event.Target = storyId.String()
		event.ResourceId = storyId
		event.ActorUsername = username
		event.Message = username + " story added to read later list with the ID:" + storyId.String()
		err = SendEventMessage(event, 0)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return fmt.Errorf("you already added this story to your read later list")
}

func (r ReadLaterRepoImpl) GetByUsername(username string) (*domain.ReadLaterDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.ReadLaterCollection.Find(context.TODO(), bson.D{{"username", username}})

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &r.ReadLaterList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	r.ReadLaterDto.ReadLaterItems = r.ReadLaterList

	return &r.ReadLaterDto, nil
}

func (r ReadLaterRepoImpl) Delete(id primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	res, err := conn.ReadLaterCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}, {"username", username}})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("you can't delete this item")
	}

	return nil
}

func NewReadLaterRepoImpl() ReadLaterRepoImpl {
	var readLaterRepoImpl ReadLaterRepoImpl

	return readLaterRepoImpl
}

