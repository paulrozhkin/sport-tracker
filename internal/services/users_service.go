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

func (us *UsersService) CreateUser(user models.User) (*models.User, error) {
	return us.userRepository.CreateUser(user)
}

func (us *UsersService) GetUserByUsername(username string) (*models.User, error) {
	return us.userRepository.GetUserByUsername(username)
}
