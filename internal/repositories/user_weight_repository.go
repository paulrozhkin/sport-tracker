package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"go.uber.org/zap"
	"time"
)

type UserWeightMeasurementRepository struct {
	store *infrastructure.Store
	log   *zap.SugaredLogger
}

func NewUserWeightMeasurementRepository(store *infrastructure.Store,
	logger *zap.SugaredLogger) (*UserWeightMeasurementRepository, error) {
	return &UserWeightMeasurementRepository{store: store, log: logger}, nil
}

// CreateUserWeightMeasurement Create new userStatistic
func (er *UserWeightMeasurementRepository) CreateUserWeightMeasurement(weight models.UserWeight) (*models.UserWeight, error) {
	weight.FillForCreate()
	query := `INSERT INTO weight_statistic (id, created, updated, weight, date, user_id) 
				VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := er.store.Pool.Exec(context.Background(), query, weight.Id, weight.Created, weight.Updated,
		weight.Weight, weight.Date, weight.User)
	if err != nil {
		er.log.Error("Failed to create user weight", err)
		return nil, err
	}
	return er.GetUserWeightMeasurement(weight.Id)
}

func (er *UserWeightMeasurementRepository) UpdateUserWeightMeasurement(userStatistic models.UserWeight) (*models.UserWeight, error) {
	userStatistic.FillForUpdate()
	query := `UPDATE weight_statistic SET updated=$2, weight=$3, date=$4 WHERE id=$1`
	_, err := er.store.Pool.Exec(context.Background(), query, userStatistic.Id, userStatistic.Updated,
		userStatistic.Weight, userStatistic.Date)
	if err != nil {
		er.log.Error("Failed to update user weight", err)
		return nil, err
	}
	return er.GetUserWeightMeasurement(userStatistic.Id)
}

func (er *UserWeightMeasurementRepository) GetUserWeightMeasurement(id string) (*models.UserWeight, error) {
	query := `SELECT id, created, updated, weight, date, user_id
				FROM weight_statistic WHERE id=$1`
	row := er.store.Pool.QueryRow(context.Background(), query, id)
	result, err := rowToUserWeightMeasurement(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("user weight", id)
	} else if err != nil {
		er.log.Errorf("Failed to get user weight by id %s due to: %v", id, err)
		return nil, err
	}
	return result, nil
}

func (er *UserWeightMeasurementRepository) GetUserWeightMeasurementByUserId(userId string) ([]*models.UserWeight, error) {
	query := `SELECT id, created, updated, weight, date, user_id
				FROM weight_statistic WHERE user_id=$1 ORDER BY date`
	rows, err := er.store.Pool.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	result := make([]*models.UserWeight, 0)
	for rows.Next() {
		weight, err := rowToUserWeightMeasurement(rows)
		if err != nil {
			continue
		}
		result = append(result, weight)
	}
	return result, nil
}

func (er *UserWeightMeasurementRepository) DeleteUserWeightMeasurementById(id string) error {
	query := `DELETE FROM weight_statistic WHERE id = $1;`
	res, err := er.store.Pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return models.NewNotFoundByIdError("user weight", id)
	}
	return nil
}

func (er *UserWeightMeasurementRepository) GetWorkoutVisitByDateForUser(userId string, from time.Time, to time.Time) (int, error) {
	query := `SELECT COUNT(*)
				FROM workouts_statistic as ws
						 JOIN user_workouts uw on ws.user_workout = uw.id
				WHERE uw.user_id = $1
				  AND ws.workout_date IS NOT NULL
				  AND ws.workout_date >= $2 
				  AND ws.workout_date <= $3`
	var count int64
	err := er.store.Pool.QueryRow(context.Background(), query, userId, from, to).Scan(&count)
	if err != nil {
		return 0, err
	}
	return int(count), err
}

func rowToUserWeightMeasurement(row pgx.Row) (*models.UserWeight, error) {
	userStatistic := &models.UserWeight{}
	err := row.Scan(&userStatistic.Id, &userStatistic.Created,
		&userStatistic.Updated, &userStatistic.Weight,
		&userStatistic.Date, &userStatistic.User)
	if err != nil {
		return nil, err
	}
	return userStatistic, nil
}
