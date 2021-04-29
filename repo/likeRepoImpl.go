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

type LikeRepoImpl struct {
	Like domain.Like
	LikeList []domain.Like
}

func (l LikeRepoImpl) CreateLikeForStory(username string, storyId primitive.ObjectID) error {
	// determine whether user already disliked the story or not
	story := new(domain.Story)
	err := database.GetInstance().StoriesCollection.FindOne(context.TODO(), bson.D{{"_id", storyId}}).Decode(&story)
	userDisLiked := false
	disLikeList := make([]domain.Dislike, 0, story.DislikeCount)
	disLike := new(domain.Dislike)
	if err != nil {
		return err
	}

	query := bson.M{"username": bson.M{"$in": story.Dislikes}}

	// Get all users
	cur, err := database.GetInstance().StoriesCollection.Find(context.TODO(), query)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.Dislike
		err = cur.Decode(&elem)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		if username != elem.AuthorUsername {
			disLikeList = append(disLikeList, elem)
		} else {
			userDisLiked = true
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

	// if user disliked the comment remove the dislike from the dislike arr
	if userDisLiked {
		opts := options.FindOneAndUpdate().SetUpsert(true)
		filter := bson.D{{"username", username}}
		update := bson.D{{"$set", bson.D{{"dislikes", &disLikeList}}}}

		err = database.GetInstance().DislikesCollection.FindOneAndUpdate(context.TODO(),
			filter, update, opts).Decode(&disLike)

		if err != nil {
			return err
		}
	}

	// if user liked the story remove the like from the likes arr
	l.Like.Id = primitive.NewObjectID()
	l.Like.AuthorUsername = username
	l.Like.CreatedAt = time.Now()

	_, err = database.GetInstance().LikesCollection.InsertOne(context.TODO(), &l.Like)

	if err != nil {
		return err
	}

	filter := bson.D{{"_id", storyId}}
	update := bson.M{"$push": bson.M{"likes": username}}

	_, err = database.GetInstance().StoriesCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (l LikeRepoImpl) CreateLikeForComment(username string, commentId primitive.ObjectID) error {
	// determine whether user already disliked the story or not
	comment := new(domain.Comment)
	err := database.GetInstance().CommentsCollection.FindOne(context.TODO(), bson.D{{"_id", commentId}}).Decode(&comment)
	userDisLiked := false
	disLikeList := make([]domain.Dislike, 0, comment.DislikeCount)
	disLike := new(domain.Dislike)
	if err != nil {
		return err
	}

	query := bson.M{"username": bson.M{"$in": comment.Dislikes}}

	// Get all users
	cur, err := database.GetInstance().CommentsCollection.Find(context.TODO(), query)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.Dislike
		err = cur.Decode(&elem)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		if username != elem.AuthorUsername {
			disLikeList = append(disLikeList, elem)
		} else {
			userDisLiked = true
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

	// if user disliked the comment remove the dislike from the dislike arr
	if userDisLiked {
		opts := options.FindOneAndUpdate().SetUpsert(true)
		filter := bson.D{{"username", username}}
		update := bson.D{{"$set", bson.D{{"dislikes", &disLikeList}}}}

		err = database.GetInstance().DislikesCollection.FindOneAndUpdate(context.TODO(),
			filter, update, opts).Decode(&disLike)

		if err != nil {
			return err
		}
	}

	// if user liked the comment remove the like from the likes arr
	l.Like.Id = primitive.NewObjectID()
	l.Like.AuthorUsername = username
	l.Like.CreatedAt = time.Now()

	_, err = database.GetInstance().LikesCollection.InsertOne(context.TODO(), &l.Like)

	if err != nil {
		return err
	}

	filter := bson.D{{"_id", commentId}}
	update := bson.M{"$push": bson.M{"likes": username}}

	_, err = database.GetInstance().CommentsCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (l LikeRepoImpl) DeleteByUsername(username string) error {
	_, err := database.GetInstance().CommentsCollection.DeleteOne(context.TODO(), bson.D{{"authorUsername", username}})
	if err != nil {
		return err
	}
	return nil
}

func NewLikeRepoImpl() LikeRepoImpl {
	var likeRepoImpl LikeRepoImpl

	return likeRepoImpl
}