package services

import (
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
	"go.uber.org/zap"
	"time"
)

type UserWorkoutsCalendarService struct {
	workoutsStatisticRepository *repositories.WorkoutsStatisticRepository
	logger                      *zap.SugaredLogger
	workoutsService             *WorkoutsService
	userWorkoutsRepository      *repositories.UserWorkoutsRepository
}

func NewUserWorkoutsCalendarService(logger *zap.SugaredLogger,
	workoutsStatisticRepository *repositories.WorkoutsStatisticRepository,
	workoutsService *WorkoutsService,
	userWorkoutsRepository *repositories.UserWorkoutsRepository) (*UserWorkoutsCalendarService, error) {
	return &UserWorkoutsCalendarService{logger: logger,
		workoutsStatisticRepository: workoutsStatisticRepository,
		workoutsService:             workoutsService,
		userWorkoutsRepository:      userWorkoutsRepository}, nil
}

func (us *UserWorkoutsCalendarService) GetWorkoutStatisticById(statisticId string) (*models.WorkoutStatistic, error) {
	workoutStatisticInfo, err := us.workoutsStatisticRepository.GetWorkoutStatisticById(statisticId)
	if err != nil {
		return nil, err
	}
	workoutStatisticInfo.Workout, err = us.workoutsService.GetWorkoutById(workoutStatisticInfo.Workout.Id)
	if err != nil {
		return nil, err
	}
	return workoutStatisticInfo, nil
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
	calendar.History = make([]*models.WorkoutStatistic, 0)
	calendar.Upcoming = make([]*models.WorkoutStatistic, 0)
	currentDayUnix := utils.GetTodayUtc().Unix()
	var nextWorkout *models.WorkoutStatistic
	for _, statisticItem := range statistic {
		dateToCompare := statisticItem.ScheduledDate
		if statisticItem.WorkoutDate != nil {
			dateToCompare = *statisticItem.WorkoutDate
		}
		if currentDayUnix > dateToCompare.Unix() {
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

func (us *UserWorkoutsCalendarService) ConfirmVisit(confirmModel models.ConfirmVisit) (*models.WorkoutStatistic, error) {
	if confirmModel.WorkoutVisitId == "" {
		return nil, fmt.Errorf("userId %s in GetCalendarForUser", models.ArgumentNullOrEmptyError)
	}
	statistic, err := us.workoutsStatisticRepository.GetWorkoutStatisticById(confirmModel.WorkoutVisitId)
	if err != nil {
		return nil, err
	}
	statistic.Comment = confirmModel.Comment
	statistic.WorkoutDate = confirmModel.WorkoutDate
	workout, err := us.workoutsStatisticRepository.UpdateWorkoutsStatistic(*statistic)
	if err != nil {
		return nil, err
	}
	workout.Workout, err = us.workoutsService.GetWorkoutById(workout.Workout.Id)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (us *UserWorkoutsCalendarService) DeleteScheduledWorkoutsForUserWorkout(userWorkoutId string) error {
	return us.workoutsStatisticRepository.DeleteScheduledWorkoutsForUserWorkout(userWorkoutId)
}

func (us *UserWorkoutsCalendarService) CreateCalendarForRepeatableActiveUserWorkouts() error {
	us.logger.Info("Started generation of new calendar")
	// Get all active repeatable workouts
	activeWorkouts, err := us.userWorkoutsRepository.GetActiveRepeatableUserWorkouts()
	if err != nil {
		return err
	}
	return us.createCalendarForUserWorkouts(activeWorkouts)
}

func (us *UserWorkoutsCalendarService) CreateCalendarForUserWorkout(userWorkout *models.UserWorkout) error {
	return us.createCalendarForUserWorkouts([]*models.UserWorkout{userWorkout})
}

func (us *UserWorkoutsCalendarService) createCalendarForUserWorkouts(workouts []*models.UserWorkout) error {
	var userWorkoutIds []string
	for _, item := range workouts {
		userWorkoutIds = append(userWorkoutIds, item.Id)
	}
	// Get last scheduled time for every active user workout
	lastSchedules, err := us.workoutsStatisticRepository.GetLastWorkoutStatisticForUserWorkouts(userWorkoutIds)
	if err != nil {
		return err
	}
	_ = lastSchedules
	var newStatistic []*models.WorkoutStatistic
	// Add two week for generation
	twoWeek := utils.GetTodayUtc().AddDate(0, 0, 7*2)
	for _, item := range workouts {
		var lastSchedule time.Time
		var lastWorkout *models.Workout
		if value, ok := lastSchedules[item.Id]; ok && value != nil {
			lastSchedule = value.ScheduledDate
			lastSchedule = lastSchedule.AddDate(0, 0, 1)
			lastWorkout = getNextWorkout(item.WorkoutPlan.Workouts, value.Workout, item.WorkoutPlan.Repeatable)
		} else {
			lastSchedule = utils.TruncateToDay(item.Created)
			lastWorkout = item.WorkoutPlan.Workouts[0]
		}

		if item.WorkoutPlan.Repeatable {
			// Generate a scheduled calendar for each next day
			for ; lastSchedule.Before(twoWeek); lastSchedule = lastSchedule.AddDate(0, 0, 1) {
				if weekdayContainsInSchedule(item.Schedule, lastSchedule.Weekday()) {
					statistic := &models.WorkoutStatistic{UserWorkout: &models.UserWorkout{}, Workout: &models.Workout{}}
					statistic.ScheduledDate = lastSchedule
					statistic.UserWorkout = item
					statistic.Workout = lastWorkout
					newStatistic = append(newStatistic, statistic)
					lastWorkout = getNextWorkout(item.WorkoutPlan.Workouts, lastWorkout, true)
				}
			}
		} else {
			for ; ; lastSchedule = lastSchedule.AddDate(0, 0, 1) {
				if weekdayContainsInSchedule(item.Schedule, lastSchedule.Weekday()) {
					if lastWorkout == nil {
						break
					}
					statistic := &models.WorkoutStatistic{UserWorkout: &models.UserWorkout{}, Workout: &models.Workout{}}
					statistic.ScheduledDate = lastSchedule
					statistic.UserWorkout = item
					statistic.Workout = lastWorkout
					newStatistic = append(newStatistic, statistic)
					lastWorkout = getNextWorkout(item.WorkoutPlan.Workouts, lastWorkout, false)
				}
			}
		}
	}
	if len(newStatistic) > 0 {
		err = us.workoutsStatisticRepository.CreateWorkoutStatistics(newStatistic)
		if err != nil {
			return err
		}
	}
	us.logger.Infof("Finished generation of new calendar. New items created %d", len(newStatistic))
	return nil
}

func weekdayContainsInSchedule(schedule []time.Weekday, compareValue time.Weekday) bool {
	for _, item := range schedule {
		if item == compareValue {
			return true
		}
	}
	return false
}

func getNextWorkout(items []*models.Workout, current *models.Workout, isRing bool) *models.Workout {
	nextIndex := 0
	for i, item := range items {
		if item.Id == current.Id {
			nextIndex = i + 1
			break
		}
	}
	if nextIndex >= len(items) {
		if !isRing {
			return nil
		}
		nextIndex = 0
	}
	return items[nextIndex]
}
