package models

import (
	"time"
)

type entity struct {
	Id      string
	Updated time.Time
	Created time.Time
}
