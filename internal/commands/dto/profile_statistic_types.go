package dto

type GeneralStatisticModel struct {
	WorkoutsPerMonth int                     `json:"workoutsPerMonth"`
	WorkoutsPerYear  int                     `json:"workoutsPerYear"`
	WeightStatistic  []*WeightStatisticModel `json:"weightStatistic"`
}

type WeightStatisticModel struct {
	Id     string   `json:"id"`
	Date   JsonDate `json:"date"`
	Weight float32  `json:"weight"`
}

type CreateWeightStatisticModel struct {
	Date   JsonDate `json:"date"`
	Weight float32  `json:"weight"`
}
