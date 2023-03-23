package models

const (
	UserGenderMale   UserGender = iota
	UserGenderFemale UserGender = iota
)

type UserGender int

type User struct {
	baseEntity
	Username string
	Password string
	Name     string
	Gender   UserGender
	Height   int
}
