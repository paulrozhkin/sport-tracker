package models

type Exercise struct {
	baseEntity
	Name             string
	ShortDescription *string
	Owner            string
	Complex          []*Exercise
}
