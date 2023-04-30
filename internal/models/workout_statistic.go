package models

import "time"

type WorkoutStatistic struct {
	baseEntity
	UserWorkout   *UserWorkout
	Workout       *Workout
	WorkoutDate   *time.Time
	ScheduledDate time.Time
	Comment       *string
}
