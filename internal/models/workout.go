package models

type Workout struct {
	baseEntity
	CustomName        *string
	CustomDescription *string
	Owner             string
	Complex           []*Exercise
}
