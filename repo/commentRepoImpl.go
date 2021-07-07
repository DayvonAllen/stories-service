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
	"time"
)

type CommentRepoImpl struct {
	Comment        domain.Comment
	CommentDto     domain.CommentDto
	CommentList    []domain.Comment
	CommentDtoList []domain.CommentDto
}

func (c CommentRepoImpl) Create(comment *domain.Comment, mongoCollection *mongo.Collection, conn *database.Connection, dbType string) error {
	if dbType == "comment" {
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

		callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

			obj := new(domain.Comment)
			err = mongoCollection.FindOne(context.TODO(), bson.D{{"_id", comment.ResourceId}}).Decode(&obj)

			if err != nil {
				return nil, fmt.Errorf("resource not found")
			}

			_, err = conn.CommentsCollection.InsertOne(context.TODO(), &comment)

			if err != nil {
				return nil, fmt.Errorf("error!!!")
			}

			err = mongoCollection.FindOne(context.TODO(), bson.D{{"_id", comment.Id}}).Decode(&c.CommentDto)

			obj.Replies = append(obj.Replies, c.CommentDto)

			filter := bson.D{{"_id", obj.Id}}
			update := bson.D{{"$set", bson.D{{"replies", obj.Replies}}}}

			_, err = conn.CommentsCollection.UpdateOne(context.TODO(), filter, update)

			return nil, err
		}

		_, err = session.WithTransaction(context.Background(), callback, txnOpts)

		if err != nil {
			return fmt.Errorf("failed to create comment")
		}

		return nil

	}

	obj := new(domain.Story)

	err := mongoCollection.FindOne(context.TODO(), bson.D{{"_id", comment.ResourceId}}).Decode(&obj)

	if err != nil {
		return fmt.Errorf("resource not found")
	}

	_, err = conn.CommentsCollection.InsertOne(context.TODO(), &comment)

	if err != nil {
		return err
	}

	return nil
}

func (c CommentRepoImpl) FindAllCommentsByResourceId(resourceID primitive.ObjectID, username string) (*[]domain.CommentDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.CommentsCollection.Find(context.TODO(), bson.D{{"resourceId", resourceID}})

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	if err = cur.All(context.TODO(), &c.CommentDtoList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	comments := make([]domain.CommentDto, 0, len(c.CommentDtoList))
	for _, v := range  c.CommentDtoList {
		v.CurrentUserLiked = helper.CurrentUserStoryInteraction(v.Likes, username)
		fmt.Println(v)

		if !v.CurrentUserLiked {
			v.CurrentUserDisLiked = helper.CurrentUserStoryInteraction(v.Dislikes, username)
		}
		comments = append(comments, v)
	}
	return &comments, nil
}

func (c CommentRepoImpl) UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) (*domain.Comment, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}, {"authorUsername", username}}
	update := bson.D{{"$set", bson.D{{"content", newContent}, {"edited", edited},
		{"updatedTime", updatedTime}}}}

	err := conn.CommentsCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&c.Comment)

	if err != nil {
		return nil, fmt.Errorf("cannot update comment that you didn't write")
	}

	return &c.Comment, nil
}

func (c CommentRepoImpl) LikeCommentById(commentId primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	ctx := context.TODO()

	cur, err := conn.CommentsCollection.Find(ctx, bson.D{
		{"_id", commentId}, {"likes", username},
	})

	if err != nil {
		return err
	}

	if cur.Next(ctx) {
		return fmt.Errorf("you've already liked this comment")
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

		filter := bson.D{{"_id", commentId}}
		update := bson.M{"$pull": bson.M{"dislikes": username}}

		res, err := conn.CommentsCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.CommentsCollection.FindOne(context.TODO(),
			filter).Decode(&c.Comment)

		c.Comment.DislikeCount = len(c.Comment.Dislikes)

		update = bson.M{"$push": bson.M{"likes": username}, "$inc": bson.M{"likeCount": 1}, "$set": bson.D{{"dislikeCount", c.Comment.DislikeCount}}}

		filter = bson.D{{"_id", commentId}}

		_, err = conn.CommentsCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to like comment")
	}

	return nil
}

func (c CommentRepoImpl) DisLikeCommentById(commentId primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	ctx := context.TODO()

	cur, err := conn.CommentsCollection.Find(ctx, bson.D{
		{"_id", commentId}, {"dislikes", username},
	})

	if err != nil {
		return err
	}

	if cur.Next(ctx) {
		return fmt.Errorf("you've already disliked this comment")
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

		filter := bson.D{{"_id", commentId}}
		update := bson.M{"$pull": bson.M{"likes": username}}

		res, err := conn.CommentsCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.CommentsCollection.FindOne(context.TODO(),
			filter).Decode(&c.Comment)

		c.Comment.LikeCount = len(c.Comment.Likes)

		update = bson.M{"$push": bson.M{"dislikes": username}, "$inc": bson.M{"dislikeCount": 1}, "$set": bson.D{{"likeCount", c.Comment.LikeCount}}}

		filter = bson.D{{"_id", commentId}}

		_, err = conn.CommentsCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to dislike comment")
	}

	return nil
}

func (c CommentRepoImpl) UpdateFlagCount(flag *domain.Flag) error {
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

	return fmt.Errorf("you've already flagged this comment")
}

func (c CommentRepoImpl) DeleteById(id primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	res, err := conn.CommentsCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}, {"authorUsername", username}})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("you can't delete a comment that you didn't create")
	}

	return nil
}

func NewCommentRepoImpl() CommentRepoImpl {
	var commentRepoImpl CommentRepoImpl

	return commentRepoImpl
}
