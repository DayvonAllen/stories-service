package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Story todo validate struct
type Story struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	Title          string             `bson:"title" json:"title"`
	Content        string             `bson:"content" json:"content"`
	Preview        string             `bson:"preview" json:"preview"`
	AuthorUsername string             `bson:"authorUsername" json:"authorUsername"`
	Likes          []string           `bson:"likes" json:"-"`
	Dislikes       []string           `bson:"dislikes" json:"-"`
	LikeCount      int                `bson:"likeCount" json:"likeCount"`
	DislikeCount   int                `bson:"dislikeCount" json:"dislikeCount"`
	Score          int                `bson:"score" json:"-"`
	Tags           []Tag              `bson:"tags" json:"tags"`
	Updated        bool               `bson:"updated" json:"updated"`
	CreatedAt      time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"-"`
	CreatedDate    string             `bson:"createdDate" json:"createdDate"`
	UpdatedDate    string             `bson:"updatedDate" json:"updatedDate"`
}

type StoryDto struct {
	Id                  primitive.ObjectID `bson:"_id" json:"id"`
	Title               string             `json:"title"`
	Content             string             `json:"content"`
	AuthorUsername      string             `json:"authorUsername"`
	Preview             string             `json:"preview"`
	Likes               []string           `bson:"likes" json:"-"`
	Dislikes            []string           `bson:"dislikes" json:"-"`
	LikeCount           int                `json:"likes"`
	DislikeCount        int                `json:"dislikes"`
	Tags                []Tag              `json:"tags"`
	Comments            *[]CommentDto      `json:"comments"`
	CurrentUserLiked    bool               `json:"currentUserLiked"`
	CurrentUserDisLiked bool               `json:"currentUserDisLiked"`
	Updated             bool               `json:"updated"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
	CreatedDate         string             `json:"createdDate"`
	UpdatedDate         string             `json:"updatedDate"`
}

type FeaturedStoryDto struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	Title          string             `json:"title"`
	Content        string             `json:"content"`
	AuthorUsername string             `json:"authorUsername"`
	Preview        string             `json:"preview"`
	LikeCount      int                `json:"likes"`
	DislikeCount   int                `json:"dislikes"`
	Tags           []Tag              `json:"tags"`
	Updated        bool               `json:"updated"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
	CreatedDate    string             `json:"createdDate"`
	UpdatedDate    string             `json:"updatedDate"`
}

type CreateStoryDto struct {
	Id             primitive.ObjectID `bson:"_id" json:"-"`
	Title          string    `bson:"title" json:"title"`
	Content        string    `bson:"content" json:"content"`
	AuthorUsername string    `bson:"authorUsername" json:"-"`
	Preview        string    `bson:"preview" json:"-"`
	Likes          []string  `bson:"likes" json:"-"`
	Dislikes       []string  `bson:"dislikes" json:"-"`
	Tags           []Tag     `json:"tags"`
	Updated        bool      `bson:"updated" json:"-"`
	CreatedAt      time.Time `bson:"createdAt" json:"-"`
	UpdatedAt      time.Time `bson:"updatedAt" json:"-"`
	CreatedDate    string    `bson:"createdDate" json:"-"`
	UpdatedDate    string    `bson:"updatedDate" json:"-"`
}


type UpdateStoryDto struct {
	Title       string    `bson:"title" json:"title"`
	Content     string    `bson:"content" json:"content"`
	Preview     string    `bson:"preview" json:"-"`
	Tags        []Tag     `json:"tags"`
	Updated     bool      `bson:"updated" json:"-"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"-"`
	UpdatedDate string    `bson:"updatedDate" json:"-"`
}
