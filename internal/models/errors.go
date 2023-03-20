package models

import (
	"encoding/json"
	"go.uber.org/zap"
)

type ParamError struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type ValidationError struct {
	Errors []*ParamError
}

func (v *ValidationError) AddError(fieldName string, validationErr error) {
	v.Errors = append(v.Errors, &ParamError{Name: fieldName, Reason: validationErr.Error()})
}

func (v *ValidationError) Error() string {
	bytes, err := json.MarshalIndent(v.Errors, "", "\t")
	if err != nil {
		zap.S().Errorf("Failed conert ValidationError %v to json due to: %v", v, err)
		return err.Error()
	}
	return string(bytes)
}
