package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/zyxevls/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, name, email, password, role) VALUES ($1, $2, $3, $4, $5)", user.ID, user.Name, user.Email, user.Password, user.Role)
	return err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	return &user, err
}
