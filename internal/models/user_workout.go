package models

type DaysOfWeek int64

const (
	Monday DaysOfWeek = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type UserWorkout struct {
	baseEntity
	UserId      string
	WorkoutPlan *WorkoutPlan
	Active      bool
	Schedule    []DaysOfWeek
}
