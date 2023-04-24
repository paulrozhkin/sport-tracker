package dto

type WorkoutsCalendarModel struct {
	History  []*WorkoutStatisticModel `json:"history"`
	Current  *WorkoutStatisticModel   `json:"current"`
	Upcoming []*WorkoutStatisticModel `json:"upcoming"`
}

type WorkoutStatisticModel struct {
	Id            string            `json:"id"`
	ScheduledDate JsonDate          `json:"scheduledDate"`
	WorkoutDate   *JsonDate         `json:"workoutDate,omitempty"`
	Workout       *WorkoutFullModel `json:"workout,omitempty"`
}
