package usecase

import (
	"github.com/google/uuid"
	"github.com/zyxevls/internal/domain"
	"github.com/zyxevls/internal/helpers"
	"github.com/zyxevls/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(email, password string) error
	Login(email, password string) (string, error)
}

type authUsecase struct {
	repo repository.UserRepository
}

func NewAuthUsecase(r repository.UserRepository) AuthUsecase {
	return &authUsecase{r}
}

func (u *authUsecase) Register(email, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := &domain.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hash),
	}

	return u.repo.Create(user)
}

func (u *authUsecase) Login(email, password string) (string, error) {
	user, _ := u.repo.FindByEmail(email)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}
	return helpers.GenerateToken(user.ID)
}
