package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Story todo validate struct
type Story struct {
	Id primitive.ObjectID        `bson:"_id" json:"id"`
	Title string                 `bson:"title" json:"title"`
	Content string               `bson:"content" json:"content"`
	AuthorUsername string  		 `bson:"authorUsername" json:"authorUsername"`
	Likes  []string  `bson:"likes" json:"likes"`
	Dislikes []string  `bson:"dislikes" json:"dislikes"`
	LikeCount int                `bson:"likeCount" json:"likeCount"`
	DislikeCount int             `bson:"dislikeCount" json:"dislikeCount"`
	FlagCount []primitive.ObjectID	`bson:"flagCount" json:"-"`
	Score int                    `bson:"score" json:"-"`
	Tags []Tag	 			 `bson:"tags" json:"tags"`
	Comments []primitive.ObjectID `bson:"comments" json:"comments"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type StoryDto struct {
	Id primitive.ObjectID        `bson:"_id" json:"id"`
	Title string                 `json:"title"`
	Content string               `json:"content"`
	AuthorUsername string        `json:"authorUsername"`
	Likes  int                   `json:"likes"`
	Dislikes int                 `json:"dislikes"`
	Tags []Tag	                 `json:"tags"`
	Comments []CommentDto           `json:"comments"`
	CurrentUserLiked bool        `json:"currentUserLiked"`
	CurrentUserDisLiked bool     `json:"currentUserDisLiked"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

type CreateStoryDto struct {
	Title string                 `bson:"title" json:"title"`
	Content string               `bson:"content" json:"content"`
	AuthorUsername string        `bson:"authorUsername" json:"-"`
	Likes  []string  `bson:"likes" json:"-"`
	Dislikes []string  `bson:"dislikes" json:"-"`
	Tags []Tag	                 `json:"tags"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}

type UpdateStoryDto struct {
	Title string                 `bson:"title" json:"title"`
	Content string               `bson:"content" json:"content"`
	Tags []Tag	                 `json:"tags"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}