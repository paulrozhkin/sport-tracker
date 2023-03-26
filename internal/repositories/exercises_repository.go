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

func (er *ExercisesRepository) CreateExercise(exercise models.Exercise) (*models.Exercise, error) {
	exercise.FillForCreate()
	query := `INSERT INTO exercises (id, created, updated, name, short_description, owner, complex) 
				VALUES ($1, $2, $3, $4, $5, $6, $7)`
	complexIds := make([]string, len(exercise.Complex))
	for _, complexExercise := range exercise.Complex {
		complexIds = append(complexIds, complexExercise.Id)
	}
	_, err := er.store.Pool.Exec(context.Background(), query, exercise.Id, exercise.Created, exercise.Updated,
		exercise.Name, exercise.ShortDescription, exercise.Owner, complexIds)
	if err != nil {
		er.log.Error("Failed to create exercise", err)
		return nil, err
	}
	return er.GetExerciseById(exercise.Id)
}

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

//	func (ur *UsersRepository) GetUserByUsername(username string) (*models.User, error) {
//		query := "SELECT id, created, updated, username, password, name, gender, height  FROM users WHERE username=$1"
//		row := ur.store.Pool.QueryRow(context.Background(), query, username)
//		user, err := rowToUser(row)
//		if err != nil && errors.Is(pgx.ErrNoRows, err) {
//			return nil, models.NewNotFoundError("user", username, "username")
//		} else if err != nil {
//			ur.log.Errorf("Failed to get user by username %s due to: %v", username, err)
//			return nil, err
//		}
//		return user, nil
//	}
func rowToExercise(row pgx.Row) (*models.Exercise, error) {
	exercise := &models.Exercise{}
	shortDescription := sql.NullString{}
	err := row.Scan(&exercise.Id, &exercise.Created,
		&exercise.Updated, &exercise.Name,
		shortDescription, &exercise.Owner, &exercise.Complex)
	if err != nil {
		return nil, err
	}
	if shortDescription.Valid {
		exercise.ShortDescription = &shortDescription.String
	}
	return exercise, nil
}
