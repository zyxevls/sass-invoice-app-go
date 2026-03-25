package repository

import "github.com/zyxevls/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}
