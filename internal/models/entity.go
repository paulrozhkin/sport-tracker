package models

import (
	"github.com/google/uuid"
	"time"
)

type baseEntity struct {
	Id      string
	Updated time.Time
	Created time.Time
}

func (entity *baseEntity) FillForCreate() {
	entity.Id = uuid.New().String()
	entity.Created = time.Now().UTC()
	entity.Updated = time.Now().UTC()
}

func (entity *baseEntity) FillForUpdate() {
	entity.Updated = time.Now().UTC()
}
