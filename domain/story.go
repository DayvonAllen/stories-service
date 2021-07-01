package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Story struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Title string                 `bson:"title" json:"title"`
	Content string               `bson:"content" json:"content"`
	AuthorId primitive.ObjectID  `bson:"authorId" json:"authorId"`
	Likes  []primitive.ObjectID   `bson:"likes" json:"likes"`
	Dislikes []primitive.ObjectID  `bson:"dislikes" json:"dislikes"`
	LikeCount int                `bson:"likeCount" json:"likeCount"`
	DislikeCount int             `bson:"dislikeCount" json:"dislikeCount"`
	FlagCount []primitive.ObjectID	`bson:"flagCount" json:"-"`
	Score int                    `bson:"score" json:"score"`
	Tags []primitive.ObjectID	 `bson:"tags" json:"tags"`
	Comments []primitive.ObjectID `bson:"comments" json:"comments"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}

type StoryDto struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Title string                 `bson:"title" json:"title"`
	Content string               `bson:"content" json:"content"`
	AuthorId primitive.ObjectID  `bson:"authorId" json:"authorId"`
	Likes  []primitive.ObjectID   `bson:"likes" json:"likes"`
	Dislikes []primitive.ObjectID  `bson:"dislikes" json:"dislikes"`
	LikeCount int                `bson:"likeCount" json:"likeCount"`
	DislikeCount int             `bson:"dislikeCount" json:"dislikeCount"`
	FlagCount []primitive.ObjectID	`bson:"flagCount" json:"-"`
	Score int                    `bson:"score" json:"score"`
	Tags []primitive.ObjectID	 `bson:"tags" json:"tags"`
	Comments []primitive.ObjectID `bson:"comments" json:"comments"`
	CurrentUserLiked bool        `bson:"-" json:"currentUserLiked"`
	CurrentUserDisLiked bool        `bson:"-" json:"currentUserDisLiked"`
}