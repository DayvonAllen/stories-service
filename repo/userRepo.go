package repo

import (
	"example.com/app/domain"
)

type UserRepo interface {
	Create(*domain.User) error
	FindByUsername(string) (*domain.User, error)
	UpdateByID(domain.User)  error
}

