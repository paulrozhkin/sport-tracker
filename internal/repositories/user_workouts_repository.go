package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"go.uber.org/zap"
)

type UserWorkoutsRepository struct {
	store *infrastructure.Store
	log   *zap.SugaredLogger
}

func NewUserWorkoutsRepository(store *infrastructure.Store,
	logger *zap.SugaredLogger) (*UserWorkoutsRepository, error) {
	return &UserWorkoutsRepository{store: store, log: logger}, nil
}

// CreateUserWorkout Create new userWorkout
func (er *UserWorkoutsRepository) CreateUserWorkout(userWorkout models.UserWorkout) (*models.UserWorkout, error) {
	userWorkout.FillForCreate()
	query := `INSERT INTO user_workouts (id, created, updated, user_id, workout_plan, active, schedule) 
				VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := er.store.Pool.Exec(context.Background(), query, userWorkout.Id, userWorkout.Created, userWorkout.Updated,
		userWorkout.UserId, userWorkout.WorkoutPlan.Id, userWorkout.Active, userWorkout.Schedule)
	if err != nil {
		er.log.Error("Failed to create userWorkout", err)
		return nil, err
	}
	return er.GetUserWorkoutById(userWorkout.Id)
}

func (er *UserWorkoutsRepository) UpdateUserWorkout(userWorkout models.UserWorkout) (*models.UserWorkout, error) {
	userWorkout.FillForUpdate()
	query := `UPDATE user_workouts SET updated=$2, active=$3, schedule=$4 WHERE id=$1`
	_, err := er.store.Pool.Exec(context.Background(), query, userWorkout.Id, userWorkout.Updated,
		userWorkout.Active, userWorkout.Schedule)
	if err != nil {
		er.log.Error("Failed to update userWorkout", err)
		return nil, err
	}
	return er.GetUserWorkoutById(userWorkout.Id)
}

// GetActiveUserWorkout Get active userWorkout for userId
// TODO: it's not repository pattern, it's DAO. Refactoring later
func (er *UserWorkoutsRepository) GetActiveUserWorkout(userId string) (*models.UserWorkout, error) {
	query := `SELECT id, created, updated, user_id, workout_plan, active, schedule
				FROM user_workouts WHERE user_id=$1 and active=true`
	row := er.store.Pool.QueryRow(context.Background(), query, userId)
	result, err := rowToUserWorkout(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundError("userWorkout", userId, "user_id")
	} else if err != nil {
		er.log.Errorf("Failed to get active userWorkout by user id %s due to: %v", userId, err)
		return nil, err
	}
	return result, nil
}

func (er *UserWorkoutsRepository) GetActiveRepeatableUserWorkouts() ([]*models.UserWorkout, error) {
	query := `SELECT uw.id, uw.created, user_id, workout_plan, schedule, wp.workouts
				FROM user_workouts as uw
						 JOIN workout_plans wp on uw.workout_plan = wp.id
				WHERE active = true
				  and repeatable = true`
	rows, err := er.store.Pool.Query(context.Background(), query)
	if err != nil {
		er.log.Errorf("Failed to get user workouts due to: %v", err)
		return nil, err
	}
	var result []*models.UserWorkout
	for rows.Next() {
		exercise, rowScanErr := activeRepeatableRowToUserWorkout(rows)
		if rowScanErr != nil {
			er.log.Errorf("Failed to scan user workout due to: %v", rowScanErr)
			continue
		}
		result = append(result, exercise)
	}
	return result, nil
}

// GetUserWorkoutById Get userWorkout  by id
func (er *UserWorkoutsRepository) GetUserWorkoutById(id string) (*models.UserWorkout, error) {
	query := `SELECT id, created, updated, user_id, workout_plan, active, schedule
				FROM user_workouts WHERE id=$1`
	row := er.store.Pool.QueryRow(context.Background(), query, id)
	result, err := rowToUserWorkout(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("userWorkout", id)
	} else if err != nil {
		er.log.Errorf("Failed to get userWorkout by id %s due to: %v", id, err)
		return nil, err
	}
	return result, nil
}

func activeRepeatableRowToUserWorkout(row pgx.Row) (*models.UserWorkout, error) {
	userWorkout := &models.UserWorkout{WorkoutPlan: &models.WorkoutPlan{}}
	var workoutIds []string
	err := row.Scan(&userWorkout.Id, &userWorkout.Created,
		&userWorkout.UserId,
		&userWorkout.WorkoutPlan.Id, &userWorkout.Schedule,
		&workoutIds)
	for _, item := range workoutIds {
		workout := &models.Workout{}
		workout.Id = item
		userWorkout.WorkoutPlan.Workouts = append(userWorkout.WorkoutPlan.Workouts, workout)
	}
	if err != nil {
		return nil, err
	}
	return userWorkout, nil
}

func rowToUserWorkout(row pgx.Row) (*models.UserWorkout, error) {
	userWorkout := &models.UserWorkout{WorkoutPlan: &models.WorkoutPlan{}}
	err := row.Scan(&userWorkout.Id, &userWorkout.Created,
		&userWorkout.Updated, &userWorkout.UserId,
		&userWorkout.WorkoutPlan.Id, &userWorkout.Active,
		&userWorkout.Schedule)
	if err != nil {
		return nil, err
	}
	return userWorkout, nil
}
