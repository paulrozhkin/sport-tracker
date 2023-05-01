package services

import (
	"errors"
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
	"go.uber.org/zap"
)

type UserWorkoutsService struct {
	userWorkoutsRepository *repositories.UserWorkoutsRepository
	logger                 *zap.SugaredLogger
	calendarGenerator      *UserWorkoutsCalendarGenerator
	calendarService        *UserWorkoutsCalendarService
}

func NewUserWorkoutsService(logger *zap.SugaredLogger,
	userWorkoutsRepository *repositories.UserWorkoutsRepository,
	calendarGenerator *UserWorkoutsCalendarGenerator,
	calendarService *UserWorkoutsCalendarService) (*UserWorkoutsService, error) {
	return &UserWorkoutsService{logger: logger,
		userWorkoutsRepository: userWorkoutsRepository,
		calendarGenerator:      calendarGenerator,
		calendarService:        calendarService}, nil
}

func (us *UserWorkoutsService) CreateUserWorkout(userWorkout models.UserWorkout) (*models.UserWorkout, error) {
	if userWorkout.WorkoutPlan == nil || (userWorkout.WorkoutPlan != nil && userWorkout.WorkoutPlan.Id == "") {
		return nil, fmt.Errorf("userWorkout workoutPlan %s in CreateUserWorkout", models.ArgumentNullOrEmptyError)
	}
	if userWorkout.UserId == "" {
		return nil, fmt.Errorf("userWorkout userId %s in CreateUserWorkout", models.ArgumentNullOrEmptyError)
	}
	_, err := us.DeactivateWorkoutForUser(userWorkout.UserId)
	if err != nil {
		return nil, err
	}
	us.logger.Infof("New active workout for user %s", userWorkout.UserId)
	userWorkout.Active = true
	newItem, err := us.userWorkoutsRepository.CreateUserWorkout(userWorkout)
	if err != nil {
		return nil, err
	}
	err = us.calendarGenerator.generateWorkoutsOfDay()
	if err != nil {
		return nil, err
	}
	return newItem, nil
}

func (us *UserWorkoutsService) GetActiveRepeatableUserWorkouts() ([]*models.UserWorkout, error) {
	activeWorkouts, err := us.userWorkoutsRepository.GetActiveRepeatableUserWorkouts()
	if err != nil {
		return nil, err
	}
	return activeWorkouts, err
}

func (us *UserWorkoutsService) DeactivateWorkoutForUser(userId string) (*models.UserWorkout, error) {
	if userId == "" {
		return nil, fmt.Errorf("userWorkout userId %s in DeactivateWorkoutForUser", models.ArgumentNullOrEmptyError)
	}
	activeUserWorkout, err := us.userWorkoutsRepository.GetActiveUserWorkout(userId)
	errType := new(models.NotFoundError)
	if errors.As(err, &errType) {
		us.logger.Infof("No active workout for user %s", userId)
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	us.logger.Infof("Deactive previous workout %s for user %s", activeUserWorkout.Id, userId)
	activeUserWorkout.Active = false
	updated, err := us.userWorkoutsRepository.UpdateUserWorkout(*activeUserWorkout)
	if err != nil {
		return nil, err
	}
	err = us.calendarService.DeleteScheduledWorkoutsForUserWorkout(updated.Id)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
