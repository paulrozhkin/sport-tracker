package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorMessageInJsonForValidationError(t *testing.T) {
	validationError := new(ValidationError)
	validationError.AddError("someField1", errors.New("some error 1"))
	validationError.AddError("someField2", errors.New("some error 2"))
	expectedJson := "[{\"name\":\"someField1\",\"reason\":\"some error 1\"},{\"name\":\"someField2\",\"reason\":\"some error 2\"}]"
	buffer := new(bytes.Buffer)
	err := json.Compact(buffer, []byte(validationError.Error()))
	assert.NoError(t, err)
	assert.Equal(t, expectedJson, buffer.String())
}
