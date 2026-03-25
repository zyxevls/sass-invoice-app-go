package usecase

import (
	"github.com/google/uuid"
	"github.com/zyxevls/internal/domain"
	"github.com/zyxevls/internal/helpers"
	"github.com/zyxevls/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(name, email, password, role string) error
	Login(email, password string) (string, error)
}

type authUsecase struct {
	repo repository.UserRepository
}

func NewAuthUsecase(r repository.UserRepository) AuthUsecase {
	return &authUsecase{r}
}

func (u *authUsecase) Register(name, email, password, role string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := &domain.User{
		ID:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: string(hash),
		Role:     role,
	}

	return u.repo.Create(user)
}

func (u *authUsecase) Login(email, password string) (string, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}
	return helpers.GenerateToken(user.ID)
}
