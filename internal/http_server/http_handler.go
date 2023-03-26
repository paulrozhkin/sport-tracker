package http_server

import (
	"encoding/json"
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/http_server/routes"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
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
		h.sendErrorResponseAndLogError(w,
			utils.CreateProblemFromError(r, createCommandError))
		return
	}
	commandContext := command.GetCommandContext()
	if commandContext == nil {
		h.logger.Error("Can't create command context")
		h.sendErrorResponseAndLogError(w, utils.CreateProblemFromError(r, errors.New("can't create command")))
		return
	}

	if commandContext.CommandContent != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Error("Failed to read request", zap.Error(err))
			h.sendErrorResponseAndLogError(w,
				utils.CreateProblemFromError(r, err))
			return
		}
		parseError := json.Unmarshal(body, commandContext.CommandContent)
		if parseError != nil {
			h.logger.Error("Failed to unmarshal body", zap.Error(parseError))
			h.sendErrorResponseAndLogError(w,
				utils.CreateProblemFromError(r, parseError))
			return
		}
	}

	validationError := command.Validate()
	if validationError != nil {
		h.logger.Error("Request validation error", zap.Error(validationError))
		h.sendErrorResponseAndLogError(w,
			utils.CreateProblemFromError(r, validationError))
		return
	}

	response, executionError := command.Execute()
	if executionError != nil {
		h.logger.Error("Failed to execute request", zap.Error(executionError))
		h.sendErrorResponseAndLogError(w, utils.CreateProblemFromError(r, executionError))
		return
	}

	responseBytes, marshalError := json.MarshalIndent(response, "", "  ")
	if marshalError != nil {
		h.logger.Error("Can't marshal response", zap.Error(marshalError))
		h.sendErrorResponseAndLogError(w,
			utils.CreateProblemFromError(r, marshalError))
		return
	}

	_, sendResponseErr := w.Write(responseBytes)
	if sendResponseErr != nil {
		h.logger.Error("Failed to send response", zap.Error(sendResponseErr))
		h.sendErrorResponseAndLogError(w,
			utils.CreateProblemFromError(r, sendResponseErr))
		return
	}
}

func (h *HttpHandler) sendErrorResponseAndLogError(w http.ResponseWriter, problemDetails dto.ProblemDetails) {
	err := utils.WriteProblemToResponse(w, problemDetails)
	if err != nil {
		h.logger.Error("Failed to send error response", zap.Error(err))
	}
}
