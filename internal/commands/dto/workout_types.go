package dto

type WorkoutCreateModel struct {
	CustomName        *string  `json:"customName"`
	CustomDescription *string  `json:"customDescription"`
	Complex           []string `json:"complex"`
}

type WorkoutFullModel struct {
	Id                string               `json:"id"`
	CustomName        *string              `json:"customName"`
	CustomDescription *string              `json:"customDescription"`
	Complex           []*ExerciseFullModel `json:"complex"`
}

type WorkoutShortModel struct {
	Id                string  `json:"id"`
	CustomName        *string `json:"customName"`
	CustomDescription *string `json:"customDescription"`
}
