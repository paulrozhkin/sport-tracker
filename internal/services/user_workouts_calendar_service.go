package services

import (
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
	"go.uber.org/zap"
)

type UserWorkoutsCalendarService struct {
	workoutsStatisticRepository *repositories.WorkoutsStatisticRepository
	logger                      *zap.SugaredLogger
	workoutsService             *WorkoutsService
}

func NewUserWorkoutsCalendarService(logger *zap.SugaredLogger,
	workoutsStatisticRepository *repositories.WorkoutsStatisticRepository,
	workoutsService *WorkoutsService) (*UserWorkoutsCalendarService, error) {
	return &UserWorkoutsCalendarService{logger: logger,
		workoutsStatisticRepository: workoutsStatisticRepository,
		workoutsService:             workoutsService}, nil
}

func (us *UserWorkoutsCalendarService) GetCalendarForUser(userId string) (*models.WorkoutsCalendar, error) {
	if userId == "" {
		return nil, fmt.Errorf("userId %s in GetCalendarForUser", models.ArgumentNullOrEmptyError)
	}

	statistic, err := us.workoutsStatisticRepository.GetShortWorkoutsStatisticByUser(userId)
	if err != nil {
		return nil, err
	}
	us.logger.Infof("Get %d items in statistic", len(statistic))

	calendar := new(models.WorkoutsCalendar)
	currentDayUnix := utils.GetTodayUtc().Unix()
	var nextWorkout *models.WorkoutStatistic
	for _, statisticItem := range statistic {
		dateToCompare := statisticItem.ScheduledDate
		if statisticItem.WorkoutDate != nil {
			dateToCompare = *statisticItem.WorkoutDate
		}
		if currentDayUnix >= dateToCompare.Unix() {
			calendar.History = append(calendar.History, statisticItem)
		} else {
			if nextWorkout == nil {
				nextWorkout = statisticItem
			} else {
				calendar.Upcoming = append(calendar.Upcoming, statisticItem)
			}
		}
	}

	if nextWorkout != nil {
		workoutStatisticInfo, err := us.workoutsStatisticRepository.GetWorkoutStatisticById(nextWorkout.Id)
		if err != nil {
			return nil, err
		}
		workoutStatisticInfo.Workout, err = us.workoutsService.GetWorkoutById(workoutStatisticInfo.Workout.Id)
		if err != nil {
			return nil, err
		}
		calendar.Current = workoutStatisticInfo
	}

	return calendar, nil
}
