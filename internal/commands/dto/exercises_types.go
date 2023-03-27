package dto

type ExerciseCreateModel struct {
	Name             string   `json:"name"`
	ShortDescription *string  `json:"shortDescription"`
	Complex          []string `json:"complex"`
}

type ExerciseFullModel struct {
	Id               string               `json:"id"`
	Name             string               `json:"name"`
	ShortDescription *string              `json:"shortDescription"`
	Complex          []*ExerciseFullModel `json:"complex"`
}

type ExerciseShortModel struct {
	Id               string  `json:"id"`
	Name             string  `json:"name"`
	ShortDescription *string `json:"shortDescription"`
}
