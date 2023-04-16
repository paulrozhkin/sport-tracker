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

type WorkoutPlansRepository struct {
	store              *infrastructure.Store
	log                *zap.SugaredLogger
	workoutsRepository *WorkoutsRepository
}

func NewWorkoutPlansRepository(store *infrastructure.Store,
	logger *zap.SugaredLogger,
	workoutsRepository *WorkoutsRepository) (*WorkoutPlansRepository, error) {
	return &WorkoutPlansRepository{store: store, log: logger, workoutsRepository: workoutsRepository}, nil
}

// CreateWorkoutPlan Create new workoutPlan
func (wpr *WorkoutPlansRepository) CreateWorkoutPlan(workoutPlan models.WorkoutPlan) (*models.WorkoutPlan, error) {
	workoutPlan.FillForCreate()
	query := `INSERT INTO workout_plans (id, created, updated, name, short_description, owner, repeatable, workouts) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	var workoutIds []string
	if workoutPlan.Workouts != nil {
		for _, workout := range workoutPlan.Workouts {
			workoutIds = append(workoutIds, workout.Id)
		}
	}
	_, err := wpr.store.Pool.Exec(context.Background(), query, workoutPlan.Id, workoutPlan.Created, workoutPlan.Updated,
		workoutPlan.Name, workoutPlan.ShortDescription, workoutPlan.Owner, workoutPlan.Repeatable, workoutIds)
	if err != nil {
		wpr.log.Error("Failed to create workoutPlan", err)
		return nil, err
	}
	return wpr.GetWorkoutPlanById(workoutPlan.Id)
}

// UpdateWorkoutPlan Update workoutPlan
func (wpr *WorkoutPlansRepository) UpdateWorkoutPlan(workoutPlan models.WorkoutPlan) (*models.WorkoutPlan, error) {
	workoutPlan.FillForUpdate()
	query := `UPDATE workout_plans SET updated=$2, name=$3, short_description=$4, repeatable=$5, workouts=$6 WHERE id=$1`
	var workoutIds []string
	if workoutPlan.Workouts != nil {
		for _, workout := range workoutPlan.Workouts {
			workoutIds = append(workoutIds, workout.Id)
		}
	}
	_, err := wpr.store.Pool.Exec(context.Background(), query, workoutPlan.Id, workoutPlan.Updated,
		workoutPlan.Name, workoutPlan.ShortDescription, workoutPlan.Repeatable, workoutIds)
	if err != nil {
		wpr.log.Error("Failed to update workoutPlan", err)
		return nil, err
	}
	return wpr.GetWorkoutPlanById(workoutPlan.Id)
}

// GetWorkoutPlanById Get workoutPlan by id with filled complex
func (wpr *WorkoutPlansRepository) GetWorkoutPlanById(id string) (*models.WorkoutPlan, error) {
	query := `SELECT id, created, updated, name, short_description, owner, repeatable, workouts
				FROM workout_plans WHERE id=$1`
	row := wpr.store.Pool.QueryRow(context.Background(), query, id)
	workoutPlan, err := rowToWorkoutPlan(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("workoutPlan", id)
	} else if err != nil {
		wpr.log.Errorf("Failed to get workoutPlan by id %s due to: %v", id, err)
		return nil, err
	}
	for i, workout := range workoutPlan.Workouts {
		fullModel, getErr := wpr.workoutsRepository.GetWorkoutsByIdWithoutComplex(workout.Id)
		if getErr != nil {
			wpr.log.Errorf("Failed to get workout with id %s from workoutPlan %s due to: %v", workout.Id, id, err)
			return nil, getErr
		}
		workoutPlan.Workouts[i] = fullModel
	}
	return workoutPlan, nil
}

// GetWorkoutPlansByIdWithoutComplex Get workoutPlan  by id without filled complex (only ids)
func (wpr *WorkoutPlansRepository) GetWorkoutPlansByIdWithoutComplex(id string) (*models.WorkoutPlan, error) {
	query := `SELECT id, created, updated, name, short_description, owner, repeatable, workouts
				FROM workout_plans WHERE id=$1`
	row := wpr.store.Pool.QueryRow(context.Background(), query, id)
	result, err := rowToWorkoutPlan(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("workoutPlan", id)
	} else if err != nil {
		wpr.log.Errorf("Failed to get workoutPlan by id %s due to: %v", id, err)
		return nil, err
	}
	return result, nil
}

// GetWorkoutPlans Get workoutPlans without filled complex (only ids)
func (wpr *WorkoutPlansRepository) GetWorkoutPlans() ([]*models.WorkoutPlan, error) {
	query := `SELECT id, created, updated, name, short_description, owner, repeatable, workouts
				FROM workout_plans`
	rows, err := wpr.store.Pool.Query(context.Background(), query)
	if err != nil {
		wpr.log.Errorf("Failed to get workoutPlans due to: %v", err)
		return nil, err
	}
	var result []*models.WorkoutPlan
	for rows.Next() {
		workoutPlan, rowScanErr := rowToWorkoutPlan(rows)
		if rowScanErr != nil {
			wpr.log.Errorf("Failed to scan workoutPlans due to: %v", rowScanErr)
			continue
		}
		result = append(result, workoutPlan)
	}
	return result, nil
}

func (wpr *WorkoutPlansRepository) DeleteWorkoutPlanById(id string) error {
	query := `DELETE FROM workout_plans WHERE id = $1;`
	res, err := wpr.store.Pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return models.NewNotFoundByIdError("workoutPlan", id)
	}
	return nil
}

func rowToWorkoutPlan(row pgx.Row) (*models.WorkoutPlan, error) {
	workoutPlan := &models.WorkoutPlan{}
	shortDescription := sql.NullString{}
	var workoutPlanComplex []string
	err := row.Scan(&workoutPlan.Id, &workoutPlan.Created,
		&workoutPlan.Updated, &workoutPlan.Name,
		&shortDescription, &workoutPlan.Owner,
		&workoutPlan.Repeatable, &workoutPlanComplex)
	if err != nil {
		return nil, err
	}
	if shortDescription.Valid {
		workoutPlan.ShortDescription = &shortDescription.String
	}
	for _, id := range workoutPlanComplex {
		item := new(models.Workout)
		item.Id = id
		workoutPlan.Workouts = append(workoutPlan.Workouts, item)
	}
	return workoutPlan, nil
}
