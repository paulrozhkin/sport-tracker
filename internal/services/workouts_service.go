package services

import (
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
)

type WorkoutsService struct {
	workoutRepository *repositories.WorkoutsRepository
}

func NewWorkoutsService(workoutRepository *repositories.WorkoutsRepository) (*WorkoutsService, error) {
	return &WorkoutsService{workoutRepository: workoutRepository}, nil
}

func (ws *WorkoutsService) CreateWorkout(workout models.Workout) (*models.Workout, error) {
	if workout.Owner == "" {
		return nil, fmt.Errorf("workout owner %s in CreateWorkout", models.ArgumentNullOrEmptyError)
	}
	return ws.workoutRepository.CreateWorkout(workout)
}

func (ws *WorkoutsService) UpdateWorkout(workout models.Workout) (*models.Workout, error) {
	if workout.Id == "" {
		return nil, fmt.Errorf("workout id %s in UpdateWorkout", models.ArgumentNullOrEmptyError)
	}
	originalWorkout, err := ws.workoutRepository.GetWorkoutsByIdWithoutComplex(workout.Id)
	if err != nil {
		return nil, err
	}
	if originalWorkout.Owner != workout.Owner {
		return nil, models.NewNoRightsOnEntityError("workout", workout.Id)
	}
	return ws.workoutRepository.UpdateWorkout(workout)
}

func (ws *WorkoutsService) GetWorkoutById(workoutId string) (*models.Workout, error) {
	if workoutId == "" {
		return nil, fmt.Errorf("workout name %s in GetWorkoutById", models.ArgumentNullOrEmptyError)
	}
	return ws.workoutRepository.GetWorkoutById(workoutId)
}

func (ws *WorkoutsService) GetWorkouts() ([]*models.Workout, error) {
	return ws.workoutRepository.GetWorkouts()
}

func (ws *WorkoutsService) DeleteWorkoutById(workoutId string) error {
	if workoutId == "" {
		return fmt.Errorf("workout name %s in DeleteWorkoutById", models.ArgumentNullOrEmptyError)
	}
	return ws.workoutRepository.DeleteWorkoutById(workoutId)
}
