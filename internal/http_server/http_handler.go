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
			createProblemFromError("", "", http.StatusInternalServerError, r, createCommandError))
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
				createProblemFromError("", "", http.StatusInternalServerError, r, err))
			return
		}
		parseError := json.Unmarshal(body, commandContext.CommandContent)
		if parseError != nil {
			h.logger.Error("Failed to unmarshal body", zap.Error(parseError))
			h.writeProblemToResponse(w,
				createProblemFromError("", "", http.StatusInternalServerError, r, parseError))
			return
		}
	}

	validationError := command.Validate()
	if validationError != nil {
		h.logger.Error("Validation error", zap.Error(validationError))
		h.writeProblemToResponse(w,
			createProblemFromError("", "", http.StatusBadRequest, r, validationError))
		return
	}

	response, executionError := command.Execute()
	if executionError != nil {
		h.logger.Error("Failed to execute request", zap.Error(executionError))
		h.writeProblemToResponse(w,
			createProblemFromError("", "", http.StatusInternalServerError, r, executionError))
		return
	}

	responseBytes, marshalError := json.MarshalIndent(response, "", "  ")
	if marshalError != nil {
		h.logger.Error("Can't marshal response", zap.Error(marshalError))
		h.writeProblemToResponse(w,
			createProblemFromError("", "", http.StatusInternalServerError, r, marshalError))
		return
	}

	_, sendResponseErr := w.Write(responseBytes)
	if sendResponseErr != nil {
		h.logger.Error("Failed to send response", zap.Error(sendResponseErr))
		h.writeProblemToResponse(w,
			createProblemFromError("", "", http.StatusInternalServerError, r, sendResponseErr))
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

func createProblemFromError(title, detail string, status int, r *http.Request, err error) dto.ProblemDetails {
	if validationError, ok := err.(*models.ValidationError); ok {
		problem := createProblemFromRequest(title, detail, status, r)
		problem.InvalidParams = validationError.Errors
		return problem
	} else {
		if detail == "" {
			detail = err.Error()
		}
	}
	return createProblemFromRequest(title, detail, status, r)
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
