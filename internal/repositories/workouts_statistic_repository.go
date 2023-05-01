package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
	"go.uber.org/zap"
	"time"
)

type WorkoutsStatisticRepository struct {
	store *infrastructure.Store
	log   *zap.SugaredLogger
}

func NewWorkoutsStatisticRepository(store *infrastructure.Store,
	logger *zap.SugaredLogger) (*WorkoutsStatisticRepository, error) {
	return &WorkoutsStatisticRepository{store: store, log: logger}, nil
}

// CreateWorkoutStatistics Create new workoutsStatistics
func (wr *WorkoutsStatisticRepository) CreateWorkoutStatistics(workoutStatistics []*models.WorkoutStatistic) error {
	batch := &pgx.Batch{}
	for _, workoutStatistic := range workoutStatistics {
		workoutStatistic.FillForCreate()
		query := `INSERT INTO workouts_statistic (id, created, updated, user_workout, workout, scheduled_date)
					VALUES ($1, $2, $3, $4, $5, $6)`
		batch.Queue(query, workoutStatistic.Id, workoutStatistic.Created, workoutStatistic.Updated,
			workoutStatistic.UserWorkout.Id, workoutStatistic.Workout.Id,
			workoutStatistic.ScheduledDate)
	}
	br := wr.store.Pool.SendBatch(context.Background(), batch)
	_, err := br.Exec()
	if err != nil {
		wr.log.Error("Failed to create workoutsStatistics", err)
		return err
	}
	return nil
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

func (wr *WorkoutsStatisticRepository) GetLastWorkoutStatisticForUserWorkouts(userWorkoutsIds []string) (map[string]*models.WorkoutStatistic, error) {
	query := `SELECT DISTINCT ON (user_workout) user_workout, scheduled_date, workout
				FROM workouts_statistic as uw
				WHERE user_workout = ANY ($1)
				ORDER BY user_workout, scheduled_date DESC`
	rows, err := wr.store.Pool.Query(context.Background(), query, userWorkoutsIds)
	if err != nil {
		wr.log.Errorf("Failed to get workouts statistic due to: %v", err)
		return nil, err
	}
	result := make(map[string]*models.WorkoutStatistic)
	for rows.Next() {
		var id string
		var scheduledTime time.Time
		var workoutId string
		rowScanErr := rows.Scan(&id,
			&scheduledTime,
			&workoutId)
		if rowScanErr != nil {
			wr.log.Errorf("Failed to scan workouts due to: %v", rowScanErr)
			continue
		}
		result[id] = &models.WorkoutStatistic{ScheduledDate: scheduledTime, Workout: &models.Workout{}}
		result[id].Workout.Id = workoutId
	}
	// append programmatically values that don't exist in db
	for _, item := range userWorkoutsIds {
		if _, ok := result[item]; !ok {
			result[item] = nil
		}
	}
	return result, nil
}

func (wr *WorkoutsStatisticRepository) DeleteScheduledWorkoutsForUserWorkout(userWorkoutId string) error {
	today := utils.GetTodayUtc()
	query := `DELETE FROM workouts_statistic WHERE user_workout=$1 and workout_date is NULL and scheduled_date>=$2`
	_, err := wr.store.Pool.Exec(context.Background(), query, userWorkoutId, today)
	if err != nil {
		return err
	}
	return nil
}

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
