package dto

type WorkoutsCalendarModel struct {
	History  []*WorkoutStatisticModel `json:"history"`
	Current  *WorkoutStatisticModel   `json:"current"`
	Upcoming []*WorkoutStatisticModel `json:"upcoming"`
}

type WorkoutStatisticShortModel struct {
	Id            string    `json:"id"`
	ScheduledDate JsonDate  `json:"scheduledDate"`
	WorkoutDate   *JsonDate `json:"workoutDate,omitempty"`
}

type WorkoutStatisticModel struct {
	Id            string            `json:"id"`
	ScheduledDate JsonDate          `json:"scheduledDate"`
	WorkoutDate   *JsonDate         `json:"workoutDate,omitempty"`
	Workout       *WorkoutFullModel `json:"workout,omitempty"`
	Comment       *string           `json:"comment,omitempty"`
}

type ConfirmVisitModel struct {
	Comment     *string   `json:"comment,omitempty"`
	WorkoutDate *JsonDate `json:"workoutDate,omitempty"`
}
