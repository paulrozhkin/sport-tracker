package services

import (
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
	"time"
)

const (
	MaxWeight = 3000
	MinWeight = 1
)

type UserStatisticService struct {
	userStatisticRepository *repositories.UserWeightMeasurementRepository
}

func NewUserStatisticService(userStatisticRepository *repositories.UserWeightMeasurementRepository) (*UserStatisticService, error) {
	return &UserStatisticService{userStatisticRepository: userStatisticRepository}, nil
}

func (us *UserStatisticService) GetGeneralStatisticForUser(userId string) (*models.UserStatistic, error) {
	generalStatistic := new(models.UserStatistic)
	userWeights, err := us.userStatisticRepository.GetUserWeightMeasurementByUserId(userId)
	if err != nil {
		return nil, err
	}
	generalStatistic.Weight = userWeights
	today := time.Now().UTC()
	firstDayOfMonth := utils.TruncateToMonth(today)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	firstDayOfYear := utils.TruncateToYear(today)
	lastDayOfYear := firstDayOfYear.AddDate(1, 0, -1)
	generalStatistic.WorkoutsPerMonth, err = us.userStatisticRepository.GetWorkoutVisitByDateForUser(userId,
		firstDayOfMonth, lastDayOfMonth)
	if err != nil {
		return nil, err
	}
	generalStatistic.WorkoutsPerYear, err = us.userStatisticRepository.GetWorkoutVisitByDateForUser(userId,
		firstDayOfYear, lastDayOfYear)
	if err != nil {
		return nil, err
	}
	return generalStatistic, nil
}

func (us *UserStatisticService) CreateUserWeightMeasurement(userWeight models.UserWeight) (*models.UserWeight, error) {
	if userWeight.Weight >= MaxWeight || userWeight.Weight <= MinWeight {
		return nil, fmt.Errorf("userWeight weight is invalid (%f) in UpdateUserWeightMeasurement", userWeight.Weight)
	}
	userWeight.Date = utils.TruncateToDay(userWeight.Date)
	return us.userStatisticRepository.CreateUserWeightMeasurement(userWeight)
}

func (us *UserStatisticService) UpdateUserWeightMeasurement(userWeight models.UserWeight) (*models.UserWeight, error) {
	if userWeight.Id == "" {
		return nil, fmt.Errorf("userWeight id %s in UpdateUserWeightMeasurement", models.ArgumentNullOrEmptyError)
	}
	if userWeight.Weight >= MaxWeight || userWeight.Weight <= MinWeight {
		return nil, fmt.Errorf("userWeight weight is invalid (%f) in UpdateUserWeightMeasurement", userWeight.Weight)
	}
	originalUserStatistic, err := us.userStatisticRepository.GetUserWeightMeasurement(userWeight.Id)
	if err != nil {
		return nil, err
	}
	if originalUserStatistic.User != userWeight.User {
		return nil, models.NewNoRightsOnEntityError("user weight", userWeight.Id)
	}
	userWeight.Date = utils.TruncateToDay(userWeight.Date)
	return us.userStatisticRepository.UpdateUserWeightMeasurement(userWeight)
}

func (us *UserStatisticService) DeleteWeightMeasurementById(weightMeasurementId string) error {
	if weightMeasurementId == "" {
		return fmt.Errorf("weightMeasurementId %s in DeleteWeightMeasurementById", models.ArgumentNullOrEmptyError)
	}
	return us.userStatisticRepository.DeleteUserWeightMeasurementById(weightMeasurementId)
}
