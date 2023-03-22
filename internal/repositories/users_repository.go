package repositories

import (
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/models"
)

type UsersRepository struct {
	store *infrastructure.Store
}

func NewUsersRepository(store *infrastructure.Store) (*UsersRepository, error) {
	return &UsersRepository{store: store}, nil
}

func (*UsersRepository) CreateUser(credentials models.Credentials) (*models.User, error) {
	return nil, fmt.Errorf("not implemented")
}
