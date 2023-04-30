package models

import "time"

type ConfirmVisit struct {
	WorkoutVisitId string
	Comment        *string
	WorkoutDate    *time.Time
}
