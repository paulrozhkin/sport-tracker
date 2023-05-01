package models

import "time"

type UserWorkout struct {
	baseEntity
	UserId      string
	WorkoutPlan *WorkoutPlan
	Active      bool
	Schedule    []time.Weekday
}
