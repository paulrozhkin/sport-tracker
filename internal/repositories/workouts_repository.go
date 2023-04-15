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

type WorkoutsRepository struct {
	store               *infrastructure.Store
	log                 *zap.SugaredLogger
	exercisesRepository *ExercisesRepository
}

func NewWorkoutsRepository(store *infrastructure.Store,
	logger *zap.SugaredLogger,
	exercisesRepository *ExercisesRepository) (*WorkoutsRepository, error) {
	return &WorkoutsRepository{store: store, log: logger, exercisesRepository: exercisesRepository}, nil
}

// CreateWorkout Create new workout
func (wr *WorkoutsRepository) CreateWorkout(workout models.Workout) (*models.Workout, error) {
	workout.FillForCreate()
	query := `INSERT INTO workouts (id, created, updated, custom_name, custom_description, owner, complex) 
				VALUES ($1, $2, $3, $4, $5, $6, $7)`
	var complexIds []string
	if workout.Complex != nil {
		for _, complexWorkout := range workout.Complex {
			complexIds = append(complexIds, complexWorkout.Id)
		}
	}
	_, err := wr.store.Pool.Exec(context.Background(), query, workout.Id, workout.Created, workout.Updated,
		workout.CustomName, workout.CustomDescription, workout.Owner, complexIds)
	if err != nil {
		wr.log.Error("Failed to create workout", err)
		return nil, err
	}
	return wr.GetWorkoutById(workout.Id)
}

func (wr *WorkoutsRepository) UpdateWorkout(workout models.Workout) (*models.Workout, error) {
	workout.FillForUpdate()
	query := `UPDATE workouts SET updated=$2, custom_name=$3, custom_description=$4, complex=$5 WHERE id=$1`
	var complexIds []string
	if workout.Complex != nil {
		for _, complexWorkout := range workout.Complex {
			complexIds = append(complexIds, complexWorkout.Id)
		}
	}
	_, err := wr.store.Pool.Exec(context.Background(), query, workout.Id, workout.Updated,
		workout.CustomName, workout.CustomDescription, complexIds)
	if err != nil {
		wr.log.Error("Failed to update workout", err)
		return nil, err
	}
	return wr.GetWorkoutById(workout.Id)
}

// GetWorkoutById Get workout by id with filled complex
func (wr *WorkoutsRepository) GetWorkoutById(id string) (*models.Workout, error) {
	query := `SELECT id, created, updated, custom_name, custom_description, owner, complex
				FROM workouts WHERE id=$1`
	row := wr.store.Pool.QueryRow(context.Background(), query, id)
	workout, err := rowToWorkout(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("workout", id)
	} else if err != nil {
		wr.log.Errorf("Failed to get workout by id %s due to: %v", id, err)
		return nil, err
	}
	for i, exercise := range workout.Complex {
		fullExercise, exerciseErr := wr.exercisesRepository.GetExerciseById(exercise.Id)
		if exerciseErr != nil {
			wr.log.Errorf("Failed to get exercise with id %s from workout %s due to: %v", exercise.Id, id, err)
			return nil, exerciseErr
		}
		workout.Complex[i] = fullExercise
	}
	return workout, nil
}

// GetWorkoutsByIdWithoutComplex Get workout  by id without filled complex (only ids)
func (wr *WorkoutsRepository) GetWorkoutsByIdWithoutComplex(id string) (*models.Workout, error) {
	query := `SELECT id, created, updated, custom_name, custom_description, owner, complex
				FROM workouts WHERE id=$1`
	row := wr.store.Pool.QueryRow(context.Background(), query, id)
	result, err := rowToWorkout(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("workout", id)
	} else if err != nil {
		wr.log.Errorf("Failed to get workout by id %s due to: %v", id, err)
		return nil, err
	}
	return result, nil
}

// GetWorkouts Get workouts without filled complex (only ids)
func (wr *WorkoutsRepository) GetWorkouts() ([]*models.Workout, error) {
	query := `SELECT id, created, updated, custom_name, custom_description, owner, complex
				FROM workouts`
	rows, err := wr.store.Pool.Query(context.Background(), query)
	if err != nil {
		wr.log.Errorf("Failed to get workouts due to: %v", err)
		return nil, err
	}
	var result []*models.Workout
	for rows.Next() {
		workout, rowScanErr := rowToWorkout(rows)
		if rowScanErr != nil {
			wr.log.Errorf("Failed to scan workouts due to: %v", rowScanErr)
			continue
		}
		result = append(result, workout)
	}
	return result, nil
}

func (wr *WorkoutsRepository) DeleteWorkoutById(id string) error {
	query := `DELETE FROM workouts WHERE id = $1;`
	res, err := wr.store.Pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return models.NewNotFoundByIdError("workout", id)
	}
	return nil
}

func rowToWorkout(row pgx.Row) (*models.Workout, error) {
	workout := &models.Workout{}
	customName := sql.NullString{}
	customDescription := sql.NullString{}
	var workoutComplex []string
	err := row.Scan(&workout.Id, &workout.Created,
		&workout.Updated, &customName,
		&customDescription, &workout.Owner, &workoutComplex)
	if err != nil {
		return nil, err
	}
	if customName.Valid {
		workout.CustomName = &customName.String
	}
	if customDescription.Valid {
		workout.CustomDescription = &customDescription.String
	}
	for _, id := range workoutComplex {
		exercise := new(models.Exercise)
		exercise.Id = id
		workout.Complex = append(workout.Complex, exercise)
	}
	return workout, nil
}
