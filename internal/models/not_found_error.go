package models

import "fmt"

type NotFoundError struct {
	entity string
	column string
	value  string
}

func NewNotFoundError(entity, valueColumn, searchColumn string) *NotFoundError {
	return &NotFoundError{entity: entity, value: valueColumn, column: searchColumn}
}

func NewNotFoundByIdError(entity, valueColumn string) *NotFoundError {
	return NewNotFoundError(entity, valueColumn, "id")
}

func (v *NotFoundError) Error() string {
	return fmt.Sprintf("Can't find entity '%s' by '%s' = '%s'", v.entity, v.column, v.value)
}
