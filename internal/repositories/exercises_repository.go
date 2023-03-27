package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"go.uber.org/zap"
)

type ExercisesRepository struct {
	store *infrastructure.Store
	log   *zap.SugaredLogger
}

func NewExercisesRepository(store *infrastructure.Store,
	logger *zap.SugaredLogger) (*ExercisesRepository, error) {
	return &ExercisesRepository{store: store, log: logger}, nil
}

// CreateExercise Create new exercise
func (er *ExercisesRepository) CreateExercise(exercise models.Exercise) (*models.Exercise, error) {
	exercise.FillForCreate()
	query := `INSERT INTO exercises (id, created, updated, name, short_description, owner, complex) 
				VALUES ($1, $2, $3, $4, $5, $6, $7)`
	var complexIds []string
	if exercise.Complex != nil {
		for _, complexExercise := range exercise.Complex {
			complexIds = append(complexIds, complexExercise.Id)
		}
	}
	_, err := er.store.Pool.Exec(context.Background(), query, exercise.Id, exercise.Created, exercise.Updated,
		exercise.Name, exercise.ShortDescription, exercise.Owner, complexIds)
	if err != nil {
		er.log.Error("Failed to create exercise", err)
		return nil, err
	}
	return er.GetExerciseById(exercise.Id)
}

// GetExerciseById Get exercise by id with filled complex
func (er *ExercisesRepository) GetExerciseById(id string) (*models.Exercise, error) {
	query := `SELECT id, created, updated, name, short_description, owner, complex
				FROM exercises WHERE id=$1`
	row := er.store.Pool.QueryRow(context.Background(), query, id)
	exercise, err := rowToExercise(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("exercise", id)
	} else if err != nil {
		er.log.Errorf("Failed to get user by id %s due to: %v", id, err)
		return nil, err
	}
	return exercise, nil
}

// GetExercises Get exercise without filled complex (only ids)
func (er *ExercisesRepository) GetExercises() ([]*models.Exercise, error) {
	query := `SELECT id, created, updated, name, short_description, owner, complex
				FROM exercises`
	rows, err := er.store.Pool.Query(context.Background(), query)
	if err != nil {
		er.log.Errorf("Failed to get exercises due to: %v", err)
		return nil, err
	}
	var result []*models.Exercise
	for rows.Next() {
		exercise, rowScanErr := rowToExercise(rows)
		if rowScanErr != nil {
			er.log.Errorf("Failed to scan exercises due to: %v", rowScanErr)
			continue
		}
		result = append(result, exercise)
	}
	return result, nil
}

func (er *ExercisesRepository) fillComplex(exercise *models.Exercise) {
	for _, exercise := range exercise.Complex {

	}
}

func rowToExercise(row pgx.Row) (*models.Exercise, error) {
	exercise := &models.Exercise{}
	shortDescription := sql.NullString{}
	var exerciseComplex []string
	err := row.Scan(&exercise.Id, &exercise.Created,
		&exercise.Updated, &exercise.Name,
		&shortDescription, &exercise.Owner, &exerciseComplex)
	if err != nil {
		return nil, err
	}
	if shortDescription.Valid {
		exercise.ShortDescription = &shortDescription.String
	}
	for _, id := range exerciseComplex {
		internalExercise := new(models.Exercise)
		internalExercise.Id = id
		exercise.Complex = append(exercise.Complex, internalExercise)
	}
	return exercise, nil
}
