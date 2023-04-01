package models

import "fmt"

type NoRightsOnEntityError struct {
	entity string
	value  string
}

func NewNoRightsOnEntityError(entity, value string) *NoRightsOnEntityError {
	return &NoRightsOnEntityError{entity: entity, value: value}
}

func (v *NoRightsOnEntityError) Error() string {
	return fmt.Sprintf("User does not have rights to entity '%s' with id '%s'.", v.entity, v.value)
}
