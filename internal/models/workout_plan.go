package models

type WorkoutPlan struct {
	baseEntity
	Name             string
	ShortDescription *string
	Owner            string
	Repeatable       bool
	Workouts         []*Workout
}
