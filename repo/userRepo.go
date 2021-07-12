package repo

import (
	"example.com/app/domain"
)

type UserRepo interface {
	Create(*domain.User) error
	GetCurrentUserProfile(string) (*domain.CurrentUserProfile, error)
	GetUserProfile(string, string) (*domain.ViewUserProfile, error)
	UpdateByID(*domain.User) error
	DeleteByID(*domain.User) error
}
