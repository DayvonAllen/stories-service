package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"example.com/app/helper"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"strconv"
	"sync"
	"time"
)

type StoryRepoImpl struct {
	Story             domain.Story
	StoryDto          domain.StoryDto
	FeaturedStoryList []domain.FeaturedStoryDto
	StoryList         []domain.Story
	StoryDtoList      []domain.StoryDto
}

func (s StoryRepoImpl) Create(story *domain.CreateStoryDto) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	story.Id = primitive.NewObjectID()

	_, err := conn.StoryCollection.InsertOne(context.TODO(), &story)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	newStory := new(domain.Story)

	newStory.Id = story.Id
	newStory.Title = story.Title
	newStory.AuthorUsername = story.AuthorUsername

	go func() {
		err := SendKafkaMessage(newStory, 201)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}

		event := new(domain.Event)
		event.Action = "create story"
		event.Target = newStory.Id.String()
		event.ResourceId = newStory.Id
		event.ActorUsername = newStory.AuthorUsername
		event.Message = newStory.AuthorUsername + " created a story, Id:" + newStory.Id.String()
		err = SendEventMessage(event, 0)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return nil
}

func (s StoryRepoImpl) UpdateById(id primitive.ObjectID, newContent string, newTitle string, username string, tags *[]domain.Tag, updated bool) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	filter := bson.D{{"_id", id}, {"authorUsername", username}}
	update := bson.D{{"$set",
		bson.D{{"content", newContent},
			{"title", newTitle},
			{"updatedAt", time.Now()},
			{"tags", tags},
			{"updated", updated},
		},
	}}

	_, err := conn.StoryCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return fmt.Errorf("you can't update a story you didn't write")
	}

	story := new(domain.Story)

	story.Id = id
	story.Title = newTitle

	go func() {
		err := SendKafkaMessage(story, 200)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return nil
}

func (s StoryRepoImpl) FindAll(page string, newStoriesQuery bool) (*[]domain.Story, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newStoriesQuery {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	cur, err := conn.StoryCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &s.StoryList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &s.StoryList, nil
}

func (s StoryRepoImpl) FindAllByUsername(username string) (*[]domain.StoryDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.StoryCollection.Find(context.TODO(), bson.D{{"authorUsername", username}})

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

func (s StoryRepoImpl) FeaturedStories() (*[]domain.FeaturedStoryDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	findOptions := options.FindOptions{}

	findOptions.SetLimit(10)
	findOptions.SetSort(bson.D{{"score", -1}})

	cur, err := conn.StoryCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &s.FeaturedStoryList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &s.FeaturedStoryList, nil
}

func (s StoryRepoImpl) LikeStoryById(storyId primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	ctx := context.TODO()

	cur, err := conn.StoryCollection.Find(ctx, bson.D{
		{"_id", storyId}, {"likes", username},
	})

	if err != nil {
		return err
	}

	if cur.Next(ctx) {
		return fmt.Errorf("you've already liked this story")
	}

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

		filter := bson.D{{"_id", storyId}}
		update := bson.M{"$pull": bson.M{"dislikes": username}}

		res, err := conn.StoryCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.StoryCollection.FindOne(context.TODO(),
			filter).Decode(&s.Story)

		s.Story.DislikeCount = len(s.Story.Dislikes)
		s.Story.Score++

		update = bson.M{"$push": bson.M{"likes": username}, "$inc": bson.M{"likeCount": 1}, "$set": bson.D{{"dislikeCount", s.Story.DislikeCount}, {"score", s.Story.Score}}}

		filter = bson.D{{"_id", storyId}}

		_, err = conn.StoryCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		go func() {
			event := new(domain.Event)
			event.Action = "like story"
			event.Target = storyId.String()
			event.ResourceId = storyId
			event.ActorUsername = username
			event.Message = username + " liked a story with the ID:" + storyId.String()
			err = SendEventMessage(event, 0)
			if err != nil {
				fmt.Println("Error publishing...")
				return
			}
		}()

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to like story")
	}

	return nil
}

func (s StoryRepoImpl) DisLikeStoryById(storyId primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	ctx := context.TODO()

	cur, err := conn.StoryCollection.Find(ctx, bson.D{
		{"_id", storyId}, {"dislikes", username},
	})

	if err != nil {
		return err
	}

	if cur.Next(ctx) {
		return fmt.Errorf("you've already disliked this story")
	}

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

		filter := bson.D{{"_id", storyId}}
		update := bson.M{"$pull": bson.M{"likes": username}}

		res, err := conn.StoryCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.StoryCollection.FindOne(context.TODO(),
			filter).Decode(&s.Story)

		s.Story.LikeCount = len(s.Story.Likes)
		s.Story.Score--

		update = bson.M{"$push": bson.M{"dislikes": username}, "$inc": bson.M{"dislikeCount": 1}, "$set": bson.D{{"likeCount", s.Story.LikeCount}, {"score", s.Story.Score}}}

		filter = bson.D{{"_id", storyId}}

		_, err = conn.StoryCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to dislike story")
	}

	go func() {
		event := new(domain.Event)
		event.Action = "dislike story"
		event.Target = storyId.String()
		event.ResourceId = storyId
		event.ActorUsername = username
		event.Message = username + " disliked a story with the ID:" + storyId.String()
		err = SendEventMessage(event, 0)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return nil
}

func (s StoryRepoImpl) FindById(storyID primitive.ObjectID, username string) (*domain.StoryDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.StoryCollection.FindOne(context.TODO(), bson.D{{"_id", storyID}}).Decode(&s.StoryDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.StoryDto.CurrentUserLiked = helper.CurrentUserInteraction(s.StoryDto.Likes, username)

		if !s.StoryDto.CurrentUserLiked {
			s.StoryDto.CurrentUserDisLiked = helper.CurrentUserInteraction(s.StoryDto.Dislikes, username)
		}
		return
	}()

	go func() {
		defer wg.Done()
		s.StoryDto.Comments, err = CommentRepoImpl{}.FindAllCommentsByResourceId(s.StoryDto.Id, username)
		if err != nil {
			panic(err)
		}
		return
	}()

	wg.Wait()

	go func() {
		event := new(domain.Event)
		event.Action = "story viewed"
		event.Target = storyID.String()
		event.ResourceId = storyID
		event.ActorUsername = username
		event.Message = username + " viewed story with the ID:" + storyID.String()
		err = SendEventMessage(event, 0)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return &s.StoryDto, nil
}

func (s StoryRepoImpl) UpdateFlagCount(flag *domain.Flag) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.FlagCollection.Find(context.TODO(), bson.M{
		"$and": []interface{}{
			bson.M{"flaggerID": flag.FlaggerID},
			bson.M{"flaggedResource": flag.FlaggedResource},
		},
	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		flag.Id = primitive.NewObjectID()
		_, err = conn.FlagCollection.InsertOne(context.TODO(), &flag)

		return nil
	}

	go func() {
		event := new(domain.Event)
		event.Action = "flag story"
		event.Target = flag.FlaggedResource.String()
		event.ResourceId = flag.FlaggedResource
		event.ActorUsername = flag.FlaggerID.String()
		event.Message = "story flagged"
		err = SendEventMessage(event, 0)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()


	return fmt.Errorf("you've already flagged this story")
}

func (s StoryRepoImpl) DeleteById(id primitive.ObjectID, username string) error {
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
		wg.Add(4)

		go func() {
			defer wg.Done()
			res, err := conn.StoryCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}, {"authorUsername", username}})

			if err != nil {
				panic(err)
			}

			if res.DeletedCount == 0 {
				panic(fmt.Errorf("failed to delete story"))
			}
			return
		}()

		go func() {
			defer wg.Done()
			_, err = conn.FlagCollection.DeleteMany(context.TODO(), bson.D{{"flaggedResource", id}})

			if err != nil {
				panic(err)
			}
			return
		}()

		go func() {
			defer wg.Done()
			_, err = conn.ReadLaterCollection.DeleteMany(context.TODO(), bson.D{{"story._id", id}})

			if err != nil {
				panic(err)
			}
			return
		}()

		go func() {
			defer wg.Done()
			err = CommentRepoImpl{}.DeleteManyById(id, username)

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
		return err
	}

	story := new(domain.Story)
	story.Id = id

	go func() {
		err := SendKafkaMessage(story, 204)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	go func() {
		event := new(domain.Event)
		event.Action = "delete story"
		event.Target = id.String()
		event.ResourceId = id
		event.ActorUsername = username
		event.Message = "story deleted"
		err = SendEventMessage(event, 0)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	return nil
}

func NewStoryRepoImpl() StoryRepoImpl {
	var storyRepoImpl StoryRepoImpl

	return storyRepoImpl
}
