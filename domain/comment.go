package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	StoryId        primitive.ObjectID `bson:"storyId" json:"storyId"`
	Content string  `bson:"content" json:"content"`
	AuthorUsername string  `bson:"authorUsername" json:"authorUsername"`
	Likes  []primitive.ObjectID   `bson:"likes" json:"-"`
	Dislikes []primitive.ObjectID  `bson:"dislikes" json:"-"`
	LikeCount int                `bson:"likeCount" json:"likeCount"`
	DislikeCount int             `bson:"dislikeCount" json:"dislikeCount"`
	FlagCount []primitive.ObjectID	`bson:"flagCount" json:"-"`
	Replies []primitive.ObjectID `bson:"replies" json:"replies"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}

type CommentDto struct {
	Content string  		`bson:"content" json:"content"`
	AuthorUsername string  `bson:"authorUsername" json:"authorUsername"`
	LikeCount int                `bson:"likeCount" json:"likeCount"`
	DislikeCount int             `bson:"dislikeCount" json:"dislikeCount"`
	Replies []CommentDto 		`bson:"replies" json:"replies"`
	CurrentUserLiked bool        `bson:"currentUserLiked" json:"currentUserLiked"`
	CurrentUserDisLiked bool        `bson:"currentUserDisLiked" json:"currentUserDisLiked"`
}
