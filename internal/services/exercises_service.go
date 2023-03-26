package services

import (
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
)

type ExercisesService struct {
	exerciseRepository *repositories.ExercisesRepository
}

func NewExercisesService(exerciseRepository *repositories.ExercisesRepository) (*ExercisesService, error) {
	return &ExercisesService{exerciseRepository: exerciseRepository}, nil
}

func (us *ExercisesService) CreateExercise(exercise models.Exercise) (*models.Exercise, error) {
	validationError := new(models.ValidationError)
	if exercise.Name == "" {
		return nil, fmt.Errorf("exercise name %s in CreateExercise", models.ArgumentNullOrEmptyError)
	}
	if exercise.Owner == "" {
		return nil, fmt.Errorf("exercise owner %s in CreateExercise", models.ArgumentNullOrEmptyError)
	}
	if validationError.HasErrors() {
		return nil, validationError
	}
	return us.exerciseRepository.CreateExercise(exercise)
}

func (us *ExercisesService) GetExerciseById(exerciseId string) (*models.Exercise, error) {
	return us.exerciseRepository.GetExerciseById(exerciseId)
}

func (us *ExercisesService) GetExercises() ([]*models.Exercise, error) {
	return us.exerciseRepository.GetExercises()
}
