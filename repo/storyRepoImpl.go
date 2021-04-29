package repo

import "example.com/app/domain"

type StoryRepoImpl struct {
	Story domain.Story
	StoryDto domain.StoryDto
	StoryList []domain.Story
	StoryDtoList []domain.StoryDto
}
