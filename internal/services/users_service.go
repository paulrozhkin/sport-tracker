package services

import (
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
)

type UsersService struct {
	userRepository *repositories.UsersRepository
}

func NewUserService(userRepository *repositories.UsersRepository) (*UsersService, error) {
	return &UsersService{userRepository: userRepository}, nil
}

func (us *UsersService) CreateUser(credentials models.Credentials) (*models.User, error) {
	return us.userRepository.CreateUser(credentials)
}
