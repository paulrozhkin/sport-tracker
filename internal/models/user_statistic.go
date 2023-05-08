package models

type UserStatistic struct {
	baseEntity
	WorkoutsPerMonth int
	WorkoutsPerYear  int
	Weight           []*UserWeight
}
