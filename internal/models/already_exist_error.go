package models

import "fmt"

type AlreadyExistError struct {
	entity string
	value  string
}

func NewAlreadyExistError(entity, value string) *AlreadyExistError {
	return &AlreadyExistError{entity: entity, value: value}
}

func (v *AlreadyExistError) Error() string {
	return fmt.Sprintf("Entity '%s' with value '%s' already exist", v.entity, v.value)
}
