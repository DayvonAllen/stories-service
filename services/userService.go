package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
)

type UserService interface {
	GetCurrentUserProfile(string) (*domain.CurrentUserProfile, error)
	GetUserProfile(string) (*domain.ViewUserProfile, error)
}

type DefaultUserService struct {
	repo repo.UserRepo
}

func (s DefaultUserService) GetCurrentUserProfile(username string) (*domain.CurrentUserProfile, error) {
	currentUser, err := s.repo.GetCurrentUserProfile(username)
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

func (s DefaultUserService) GetUserProfile(username string) (*domain.ViewUserProfile, error) {
	currentUser, err := s.repo.GetUserProfile(username)
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

func NewUserService(repository repo.UserRepo) DefaultUserService {
	return DefaultUserService{repository}
}