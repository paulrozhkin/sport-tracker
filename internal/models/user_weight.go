package models

import "time"

type UserWeight struct {
	baseEntity
	User   string
	Weight float32
	Date   time.Time
}
