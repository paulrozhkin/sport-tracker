package services

import (
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
)

type UserWorkoutsService struct {
	userWorkoutsRepository *repositories.UserWorkoutsRepository
}

func NewUserWorkoutsService(userWorkoutsRepository *repositories.UserWorkoutsRepository) (*UserWorkoutsService, error) {
	return &UserWorkoutsService{userWorkoutsRepository: userWorkoutsRepository}, nil
}

func (us *UserWorkoutsService) CreateUserWorkout(userWorkout models.UserWorkout) (*models.UserWorkout, error) {
	if userWorkout.WorkoutPlan == nil || (userWorkout.WorkoutPlan != nil && userWorkout.WorkoutPlan.Id == "") {
		return nil, fmt.Errorf("userWorkout workoutPlan %s in CreateUserWorkout", models.ArgumentNullOrEmptyError)
	}
	if userWorkout.UserId == "" {
		return nil, fmt.Errorf("userWorkout userId %s in CreateUserWorkout", models.ArgumentNullOrEmptyError)
	}
	return us.userWorkoutsRepository.CreateUserWorkout(userWorkout)
}
