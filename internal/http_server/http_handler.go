package http_server

import (
	"encoding/json"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/http_server/routes"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type HttpHandler struct {
	http.Handler
	logger       *zap.SugaredLogger
	routeHandler routes.Route
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	command, createCommandError := h.routeHandler.NewRouteExecutor()
	if createCommandError != nil {
		h.logger.Error("Can't create command", zap.Error(createCommandError))
		h.writeProblemToResponse(w,
			createProblemFromError("", "", r, createCommandError))
		return
	}

	commandContext := command.GetCommandContext()
	if commandContext == nil {
		h.logger.Error("Can't create command context")
		h.writeProblemToResponse(w,
			createProblemFromRequest("", "Can't create command context", http.StatusInternalServerError, r))
		return
	}

	if commandContext.CommandContent != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Error("Failed to read request", zap.Error(err))
			h.writeProblemToResponse(w,
				createProblemFromError("", "", r, err))
			return
		}
		parseError := json.Unmarshal(body, commandContext.CommandContent)
		if parseError != nil {
			h.logger.Error("Failed to unmarshal body", zap.Error(parseError))
			h.writeProblemToResponse(w,
				createProblemFromError("", "", r, parseError))
			return
		}
	}

	validationError := command.Validate()
	if validationError != nil {
		h.logger.Error("Validation error", zap.Error(validationError))
		h.writeProblemToResponse(w,
			createProblemFromError("", "", r, validationError))
		return
	}

	response, executionError := command.Execute()
	if executionError != nil {
		h.logger.Error("Failed to execute request", zap.Error(executionError))
		h.writeProblemToResponse(w,
			createProblemFromError("", "", r, executionError))
		return
	}

	responseBytes, marshalError := json.MarshalIndent(response, "", "  ")
	if marshalError != nil {
		h.logger.Error("Can't marshal response", zap.Error(marshalError))
		h.writeProblemToResponse(w,
			createProblemFromError("", "", r, marshalError))
		return
	}

	_, sendResponseErr := w.Write(responseBytes)
	if sendResponseErr != nil {
		h.logger.Error("Failed to send response", zap.Error(sendResponseErr))
		h.writeProblemToResponse(w,
			createProblemFromError("", "", r, sendResponseErr))
		return
	}
}

func (h *HttpHandler) writeProblemToResponse(w http.ResponseWriter, problem dto.ProblemDetails) {
	w.WriteHeader(problem.Status)
	responseBytes, marshalError := json.MarshalIndent(problem, "", "  ")
	if marshalError != nil {
		h.logger.Error("Can't marshal error response", zap.Error(marshalError))
		return
	}
	_, sendResponseErr := w.Write(responseBytes)
	if sendResponseErr != nil {
		h.logger.Error("Failed to send error response", zap.Error(sendResponseErr))
	}
}

func createProblemFromError(title, detail string, r *http.Request, err error) dto.ProblemDetails {
	switch customErr := err.(type) {
	case *models.ValidationError:
		problem := createProblemFromRequest(title, detail, http.StatusBadRequest, r)
		problem.InvalidParams = customErr.Errors
		return problem
	case *models.NotFoundError:
		problem := createProblemFromRequest("Entity not found", customErr.Error(), http.StatusNotFound, r)
		return problem
	case *models.AlreadyExistError:
		problem := createProblemFromRequest("Entity already exist", customErr.Error(), http.StatusConflict, r)
		return problem
	default:
		if detail == "" {
			detail = err.Error()
		}
		return createProblemFromRequest(title, detail, http.StatusInternalServerError, r)
	}
}

func createProblemFromRequest(title, detail string, status int, r *http.Request) dto.ProblemDetails {
	if title == "" {
		if status == http.StatusInternalServerError {
			title = "Internal server error"
		} else if status == http.StatusNotFound {
			title = "Entity not found"
		} else if status == http.StatusBadRequest {
			title = "BadRequest"
		}
	}
	return dto.ProblemDetails{
		Type:     "",
		Status:   status,
		Detail:   detail,
		Instance: r.RequestURI,
		Title:    title,
	}
}
