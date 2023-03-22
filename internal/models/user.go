package models

const (
	UserGenderMale   UserGender = iota
	UserGenderFemale UserGender = iota
)

type UserGender int

type User struct {
	entity
	username string
	password string
	name     string
	gender   UserGender
	height   int
}
