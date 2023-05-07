package services

import (
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type UsersService struct {
	userRepository *repositories.UsersRepository
}

func NewUsersService(userRepository *repositories.UsersRepository) (*UsersService, error) {
	return &UsersService{userRepository: userRepository}, nil
}

func (us *UsersService) CreateUser(user models.User) (*models.User, error) {
	if user.Username == "" {
		return nil, fmt.Errorf("username %s in CreateExercise", models.ArgumentNullOrEmptyError)
	}
	if user.Password == "" {
		return nil, fmt.Errorf("user password %s in CreateExercise", models.ArgumentNullOrEmptyError)
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	return us.userRepository.CreateUser(user)
}

func (us *UsersService) GetUserByUsername(username string) (*models.User, error) {
	return us.userRepository.GetUserByUsername(username)
}
