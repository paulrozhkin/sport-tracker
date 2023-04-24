package models

type WorkoutsCalendar struct {
	// Short info about history (only ids, scheduled date and workout date)
	History []*WorkoutStatistic
	// Full info about current workout
	Current *WorkoutStatistic
	// Short info about upcoming workouts (only ids and scheduled date)
	Upcoming []*WorkoutStatistic
}
