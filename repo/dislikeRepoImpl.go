package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DisLikeRepoImpl struct {
	DisLike domain.Dislike
	DisLikeList []domain.Dislike
}

func (d DisLikeRepoImpl) CreateDisLikeForStory(username string, storyId primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	// determine whether user already liked the story or not
	story := new(domain.Story)
	err := conn.StoriesCollection.FindOne(context.TODO(), bson.D{{"_id", storyId}}).Decode(&story)
	userLiked := false
	likeList := make([]domain.Like, 0, story.LikeCount)
	like := new(domain.Like)
	if err != nil {
		return err
	}

	query := bson.M{"authorUsername": bson.M{"$in": story.Likes}}

	// Get all users
	cur, err := conn.StoriesCollection.Find(context.TODO(), query)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.Like
		err = cur.Decode(&elem)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		if username != elem.AuthorUsername {
			likeList = append(likeList, elem)
		} else {
			userLiked = true
		}
	}

	err = cur.Err()

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	// if user liked the comment remove the like from the likes arr
	if userLiked {
		opts := options.FindOneAndUpdate().SetUpsert(true)
		filter := bson.D{{"authorUsername", username}}
		update := bson.D{{"$set", bson.D{{"likes", &likeList}}}}

		err = conn.LikesCollection.FindOneAndUpdate(context.TODO(),
			filter, update, opts).Decode(&like)

		if err != nil {
			return err
		}
	}

	// if user liked the story  remove the like from the likes arr
	d.DisLike.Id = primitive.NewObjectID()
	d.DisLike.AuthorUsername = username
	d.DisLike.CreatedAt = time.Now()

	_, err = conn.DislikesCollection.InsertOne(context.TODO(), &d.DisLike)

	if err != nil {
		return err
	}

	filter := bson.D{{"_id", storyId}}
	update := bson.M{"$push": bson.M{"dislikes": username}}

	_, err = conn.StoriesCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (d DisLikeRepoImpl) CreateDisLikeForComment(username string, commentId primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	// determine whether user already liked the comment or not
	comment := new(domain.Comment)
	err := conn.CommentsCollection.FindOne(context.TODO(), bson.D{{"_id", commentId}}).Decode(&comment)
	userLiked := false
	likeList := make([]domain.Like, 0, comment.LikeCount)
	like := new(domain.Like)
	if err != nil {
		return err
	}

	query := bson.M{"authorUsername": bson.M{"$in": comment.Likes}}

	// Get all users
	cur, err := conn.CommentsCollection.Find(context.TODO(), query)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.Like
		err = cur.Decode(&elem)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		if username != elem.AuthorUsername {
			likeList = append(likeList, elem)
		} else {
			userLiked = true
		}
	}

	err = cur.Err()

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	// if user liked the comment remove the like from the likes arr
	if userLiked {
		opts := options.FindOneAndUpdate().SetUpsert(true)
		filter := bson.D{{"authorUsername", username}}
		update := bson.D{{"$set", bson.D{{"likes", &likeList}}}}

		err = conn.LikesCollection.FindOneAndUpdate(context.TODO(),
			filter, update, opts).Decode(&like)

		if err != nil {
			return err
		}
	}

	// save dislike
	d.DisLike.Id = primitive.NewObjectID()
	d.DisLike.AuthorUsername = username
	d.DisLike.CreatedAt = time.Now()

	_, err = conn.DislikesCollection.InsertOne(context.TODO(), &d.DisLike)

	if err != nil {
		return err
	}

	filter := bson.D{{"_id", commentId}}
	update := bson.M{"$push": bson.M{"dislikes": username}}

	_, err = conn.CommentsCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (d DisLikeRepoImpl) DeleteByUsername(username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	_, err := conn.DislikesCollection.DeleteOne(context.TODO(), bson.D{{"authorUsername", username}})
	if err != nil {
		return err
	}
	return nil
}

func NewDisLikeRepoImpl() DisLikeRepoImpl {
	var disLikeRepoImpl DisLikeRepoImpl

	return disLikeRepoImpl
}