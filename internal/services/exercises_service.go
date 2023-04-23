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
	// TODO: validate recursion and hierarchical level exercise
	return us.exerciseRepository.CreateExercise(exercise)
}

func (us *ExercisesService) UpdateExercise(exercise models.Exercise) (*models.Exercise, error) {
	if exercise.Id == "" {
		return nil, fmt.Errorf("exercise id %s in UpdateExercise", models.ArgumentNullOrEmptyError)
	}
	if exercise.Name == "" {
		return nil, fmt.Errorf("exercise name %s in UpdateExercise", models.ArgumentNullOrEmptyError)
	}
	originalExercise, err := us.exerciseRepository.GetExercisesByIdWithoutComplex(exercise.Id)
	if err != nil {
		return nil, err
	}
	if originalExercise.Owner != exercise.Owner {
		return nil, models.NewNoRightsOnEntityError("exercise", exercise.Id)
	}
	return us.exerciseRepository.UpdateExercise(exercise)
}

func (us *ExercisesService) GetExerciseById(exerciseId string) (*models.Exercise, error) {
	if exerciseId == "" {
		return nil, fmt.Errorf("exercise name %s in GetExerciseById", models.ArgumentNullOrEmptyError)
	}
	return us.exerciseRepository.GetExerciseById(exerciseId)
}

func (us *ExercisesService) GetExercises() ([]*models.Exercise, error) {
	return us.exerciseRepository.GetExercises()
}

func (us *ExercisesService) DeleteExerciseById(exerciseId string) error {
	if exerciseId == "" {
		return fmt.Errorf("exercise name %s in DeleteExerciseById", models.ArgumentNullOrEmptyError)
	}
	return us.exerciseRepository.DeleteExerciseById(exerciseId)
}
