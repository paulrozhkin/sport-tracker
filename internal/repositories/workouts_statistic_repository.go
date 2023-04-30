package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"go.uber.org/zap"
)

type WorkoutsStatisticRepository struct {
	store *infrastructure.Store
	log   *zap.SugaredLogger
}

func NewWorkoutsStatisticRepository(store *infrastructure.Store,
	logger *zap.SugaredLogger) (*WorkoutsStatisticRepository, error) {
	return &WorkoutsStatisticRepository{store: store, log: logger}, nil
}

// GetShortWorkoutsStatisticByUser Get workouts statistic for user (only id and workout and schedule times)
func (wr *WorkoutsStatisticRepository) GetShortWorkoutsStatisticByUser(userId string) ([]*models.WorkoutStatistic, error) {
	query := `SELECT ws.id, scheduled_date, workout_date
				FROM workouts_statistic as ws
						 JOIN user_workouts uw on ws.user_workout = uw.id
				WHERE uw.user_id = $1
				ORDER BY scheduled_date`
	rows, err := wr.store.Pool.Query(context.Background(), query, userId)
	if err != nil {
		wr.log.Errorf("Failed to get workouts statistic due to: %v", err)
		return nil, err
	}
	var result []*models.WorkoutStatistic
	for rows.Next() {
		workout, rowScanErr := shortRowToWorkoutStatistic(rows)
		if rowScanErr != nil {
			wr.log.Errorf("Failed to scan workouts due to: %v", rowScanErr)
			continue
		}
		result = append(result, workout)
	}
	return result, nil
}

// GetWorkoutStatisticById Get workoutsStatistic  by id
func (wr *WorkoutsStatisticRepository) GetWorkoutStatisticById(id string) (*models.WorkoutStatistic, error) {
	query := `SELECT id, created, updated, user_workout, workout, workout_date, scheduled_date, comment
					FROM workouts_statistic WHERE id=$1`
	row := wr.store.Pool.QueryRow(context.Background(), query, id)
	result, err := rowToWorkoutStatistic(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("workoutsStatistic", id)
	} else if err != nil {
		wr.log.Errorf("Failed to get workoutsStatistic by id %s due to: %v", id, err)
		return nil, err
	}
	return result, nil
}

func (wr *WorkoutsStatisticRepository) UpdateWorkoutsStatistic(workoutsStatistic models.WorkoutStatistic) (*models.WorkoutStatistic, error) {
	workoutsStatistic.FillForUpdate()
	query := `UPDATE workouts_statistic SET updated=$2, workout_date=$3, comment=$4 WHERE id=$1`
	_, err := wr.store.Pool.Exec(context.Background(), query, workoutsStatistic.Id, workoutsStatistic.Updated,
		workoutsStatistic.WorkoutDate, workoutsStatistic.Comment)
	if err != nil {
		wr.log.Error("Failed to update workoutsStatistic", err)
		return nil, err
	}
	return wr.GetWorkoutStatisticById(workoutsStatistic.Id)
}

// // CreateWorkoutStatistic Create new workoutsStatistic
//
//	func (er *WorkoutsStatisticRepository) CreateWorkoutStatistic(workoutStatistic models.WorkoutStatistic) (*models.WorkoutStatistic, error) {
//		workoutStatistic.FillForCreate()
//		query := `INSERT INTO user_workouts (id, created, updated, user_id, workout_plan, active, schedule)
//					VALUES ($1, $2, $3, $4, $5, $6, $7)`
//		_, err := er.store.Pool.Exec(context.Background(), query, workoutStatistic.Id, workoutStatistic.Created, workoutStatistic.Updated,
//			workoutsStatistic.UserId, workoutsStatistic.WorkoutPlan.Id, workoutsStatistic.Active, workoutsStatistic.Schedule)
//		if err != nil {
//			er.log.Error("Failed to create workoutsStatistic", err)
//			return nil, err
//		}
//		return er.GetWorkoutsStatisticById(workoutsStatistic.Id)
//	}
//

//
// // GetActiveWorkoutsStatistic Get active workoutsStatistic for userId
//
//	func (er *WorkoutsStatisticRepository) GetActiveWorkoutsStatistic(userId string) (*models.WorkoutStatistic, error) {
//		query := `SELECT id, created, updated, user_id, workout_plan, active, schedule
//					FROM user_workouts WHERE user_id=$1 and active=true`
//		row := er.store.Pool.QueryRow(context.Background(), query, userId)
//		result, err := rowToWorkoutsStatistic(row)
//		if err != nil && errors.Is(pgx.ErrNoRows, err) {
//			return nil, models.NewNotFoundError("workoutsStatistic", userId, "user_id")
//		} else if err != nil {
//			er.log.Errorf("Failed to get active workoutsStatistic by user id %s due to: %v", userId, err)
//			return nil, err
//		}
//		return result, nil
//	}
//

func shortRowToWorkoutStatistic(row pgx.Row) (*models.WorkoutStatistic, error) {
	workoutsStatistic := &models.WorkoutStatistic{}
	err := row.Scan(&workoutsStatistic.Id, &workoutsStatistic.ScheduledDate,
		&workoutsStatistic.WorkoutDate)
	if err != nil {
		return nil, err
	}
	return workoutsStatistic, nil
}

func rowToWorkoutStatistic(row pgx.Row) (*models.WorkoutStatistic, error) {
	workoutsStatistic := &models.WorkoutStatistic{
		UserWorkout: new(models.UserWorkout),
		Workout:     new(models.Workout),
	}
	err := row.Scan(&workoutsStatistic.Id,
		&workoutsStatistic.Created,
		&workoutsStatistic.Updated,
		&workoutsStatistic.UserWorkout.Id,
		&workoutsStatistic.Workout.Id,
		&workoutsStatistic.WorkoutDate,
		&workoutsStatistic.ScheduledDate,
		&workoutsStatistic.Comment)
	if err != nil {
		return nil, err
	}
	return workoutsStatistic, nil
}
