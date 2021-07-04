package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryRepo interface {
	 Create(story *domain.CreateStoryDto) error
	 UpdateById(primitive.ObjectID, string, string, string, *[]domain.Tag) error
	 FindAll(string) (*[]domain.Story, error)
	 FeaturedStories() (*[]domain.FeaturedStoryDto, error)
	 LikeStoryById(primitive.ObjectID, string) error
	 DisLikeStoryById(primitive.ObjectID, string) error
	 FindById(primitive.ObjectID, string) (*domain.StoryDto, error)
	 DeleteById(primitive.ObjectID, string) error
}
