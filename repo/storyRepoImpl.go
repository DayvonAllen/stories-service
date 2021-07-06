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
	"time"
)

type StoryRepoImpl struct {
	Story        domain.Story
	StoryDto     domain.StoryDto
	FeaturedStoryList    []domain.FeaturedStoryDto
	StoryList    []domain.Story
	StoryDtoList []domain.StoryDto
}

func (s StoryRepoImpl) Create(story *domain.CreateStoryDto) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.StoriesCollection.InsertOne(context.TODO(), &story)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (s StoryRepoImpl) UpdateById(id primitive.ObjectID, newContent string, newTitle string, username string, tags *[]domain.Tag, updated bool)  error {
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

	_, err := conn.StoriesCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return fmt.Errorf("you can't update a story you didn't write")
	}

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

	cur, err := conn.StoriesCollection.Find(context.TODO(), bson.M{}, &findOptions)

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

func (s StoryRepoImpl) FeaturedStories() (*[]domain.FeaturedStoryDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	findOptions := options.FindOptions{}

	findOptions.SetLimit(10)
	findOptions.SetSort(bson.D{{"score", -1}})

	cur, err := conn.StoriesCollection.Find(context.TODO(), bson.M{}, &findOptions)

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

	cur, err := conn.StoriesCollection.Find(ctx, bson.D{
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

		res, err := conn.StoriesCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.StoriesCollection.FindOne(context.TODO(),
			filter).Decode(&s.Story)

		s.Story.DislikeCount = len(s.Story.Dislikes)
		s.Story.Score++

		update = bson.M{"$push": bson.M{"likes": username}, "$inc": bson.M{"likeCount": 1},  "$set": bson.D{{"dislikeCount", s.Story.DislikeCount}, {"score", s.Story.Score}}}

		filter = bson.D{{"_id", storyId}}

		_, err = conn.StoriesCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

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

	cur, err := conn.StoriesCollection.Find(ctx, bson.D{
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

		res, err := conn.StoriesCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.StoriesCollection.FindOne(context.TODO(),
			filter).Decode(&s.Story)

		s.Story.LikeCount = len(s.Story.Likes)
		s.Story.Score--

		update = bson.M{"$push": bson.M{"dislikes": username}, "$inc": bson.M{"dislikeCount": 1},  "$set": bson.D{{"likeCount", s.Story.LikeCount}, {"score", s.Story.Score}}}

		filter = bson.D{{"_id", storyId}}

		_, err = conn.StoriesCollection.UpdateOne(context.TODO(),
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

	return nil
}

func (s StoryRepoImpl) FindById(storyID primitive.ObjectID, username string) (*domain.StoryDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.StoriesCollection.FindOne(context.TODO(), bson.D{{"_id", storyID}}).Decode(&s.StoryDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	s.StoryDto.CurrentUserLiked = helper.CurrentUserStoryInteraction(s.StoryDto.Likes, username)

	if !s.StoryDto.CurrentUserLiked {
		s.StoryDto.CurrentUserDisLiked = helper.CurrentUserStoryInteraction(s.StoryDto.Dislikes, username)
	}

	s.StoryDto.Comments, err = CommentRepoImpl{}.FindAllCommentsByStoryId(s.StoryDto.Id)

	if err != nil {
		return nil, err
	}

	return &s.StoryDto, nil
}

func (s StoryRepoImpl) DeleteById(id primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	res, err := conn.StoriesCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}, {"authorUsername", username}})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("you can't delete a story that you didn't create")
	}

	return nil
}

func NewStoryRepoImpl() StoryRepoImpl {
	var storyRepoImpl StoryRepoImpl

	return storyRepoImpl
}
