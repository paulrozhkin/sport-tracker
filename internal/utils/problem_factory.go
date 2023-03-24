package utils

import (
	"encoding/json"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"net/http"
)

func WriteProblemToResponse(w http.ResponseWriter, problem dto.ProblemDetails) error {
	w.WriteHeader(problem.Status)
	responseBytes, marshalError := json.MarshalIndent(problem, "", "  ")
	if marshalError != nil {
		return marshalError
	}
	_, sendResponseErr := w.Write(responseBytes)
	if sendResponseErr != nil {
		return sendResponseErr
	}
	return nil
}

func CreateProblemFromError(r *http.Request, err error) dto.ProblemDetails {
	switch customErr := err.(type) {
	case *models.ValidationError:
		problem := createProblemFromRequestWithCustoms(r, "", "Invalid params in request", http.StatusBadRequest)
		problem.InvalidParams = customErr.Errors
		return problem
	case *models.NotFoundError:
		problem := createProblemFromRequestWithCustoms(r, "Entity not found", customErr.Error(), http.StatusNotFound)
		return problem
	case *models.AlreadyExistError:
		problem := createProblemFromRequestWithCustoms(r, "Entity already exist", customErr.Error(), http.StatusConflict)
		return problem
	default:
		return createProblemFromRequestWithCustoms(r, "", err.Error(), http.StatusInternalServerError)
	}
}

func CreateProblemFromRequest(r *http.Request, status int) dto.ProblemDetails {
	return createProblemFromRequestWithCustoms(r, "", "", status)
}

func createProblemFromRequestWithCustoms(r *http.Request, title, detail string, status int) dto.ProblemDetails {
	if title == "" {
		title = http.StatusText(status)
	}
	return dto.ProblemDetails{
		Type:     "",
		Status:   status,
		Detail:   detail,
		Instance: r.RequestURI,
		Title:    title,
	}
}
