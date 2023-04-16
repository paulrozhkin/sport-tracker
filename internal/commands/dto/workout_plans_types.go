package dto

type WorkoutPlanCreateModel struct {
	Name             string   `json:"name"`
	ShortDescription *string  `json:"shortDescription"`
	Repeatable       bool     `json:"repeatable"`
	Workouts         []string `json:"workouts"`
}

type WorkoutPlanFullModel struct {
	Id               string               `json:"id"`
	Name             string               `json:"name"`
	ShortDescription *string              `json:"shortDescription"`
	Repeatable       bool                 `json:"repeatable"`
	Workouts         []*WorkoutShortModel `json:"workouts"`
}

type WorkoutPlanShortModel struct {
	Id               string  `json:"id"`
	Name             string  `json:"name"`
	ShortDescription *string `json:"shortDescription"`
	Repeatable       bool    `json:"repeatable"`
}
