package http_server

import (
	"encoding/json"
	"errors"
	"github.com/docker/distribution/registry/auth"
	"github.com/go-chi/chi/v5"
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/http_server/routes"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
)

type HttpHandler struct {
	http.Handler
	logger       *zap.SugaredLogger
	routeHandler routes.Route
	tokenService *services.TokenService
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle panic
	defer func() {
		if err := recover(); err != nil {
			h.logger.Errorf("Can't execute command due to panic %v", err)
			h.sendErrorResponseAndLogError(w, utils.CreateProblemFromRequest(r, http.StatusInternalServerError))
		}
	}()
	// Create command to execute
	command, createCommandError := h.routeHandler.NewRouteExecutor()
	if createCommandError != nil {
		h.logger.Error("Can't create command", zap.Error(createCommandError))
		h.sendErrorResponseAndLogError(w,
			utils.CreateProblemFromError(r, createCommandError))
		return
	}
	// Context of command with request body and url params
	commandContext := command.GetCommandContext()
	if commandContext == nil {
		h.logger.Error("Can't create command context")
		h.sendErrorResponseAndLogError(w, utils.CreateProblemFromError(r, errors.New("can't create command")))
		return
	}
	// If authorization is required, then get claims
	if command.RequireAuthorization() {
		claims, authorizationError := h.HandleAuthorization(w, r)
		if authorizationError != nil {
			return
		}
		command.SetAuthorization(claims)
	}
	// Parse body to context struct
	if commandContext.CommandContent != nil {
		parseBodyError := h.parseContextBody(r, commandContext)
		if parseBodyError != nil {
			h.sendErrorResponseAndLogError(w,
				utils.CreateProblemFromError(r, parseBodyError))
		}
	}
	// Parse url params to context map
	if commandContext.CommandParameters != nil {
		h.parseContextParams(r, commandContext)
	}
	// Validate context struct and url params
	validationError := command.Validate()
	if validationError != nil {
		h.logger.Error("Request validation error", zap.Error(validationError))
		h.sendErrorResponseAndLogError(w,
			utils.CreateProblemFromError(r, validationError))
		return
	}
	// Execute the command
	response, executionError := command.Execute()
	if executionError != nil {
		h.logger.Error("Failed to execute request", zap.Error(executionError))
		h.sendErrorResponseAndLogError(w, utils.CreateProblemFromError(r, executionError))
		return
	}
	// Send a response
	h.sendSuccessResponse(r, w, response)
}

func (h *HttpHandler) HandleAuthorization(w http.ResponseWriter, r *http.Request) (*models.Claims, error) {
	authHeader := r.Header.Get(AuthorizationHeader)
	if authHeader == "" {
		h.sendErrorResponseAndLogError(w, utils.CreateProblemFromRequest(r, http.StatusUnauthorized))
		return nil, errors.New("no authorization header in request")
	}
	headersParts := strings.Split(authHeader, " ")
	if len(headersParts) != 2 {
		h.sendErrorResponseAndLogError(w, utils.CreateProblemFromRequest(r, http.StatusUnauthorized))
		return nil, errors.New("invalid count in header")
	}
	claims, err := h.tokenService.ParseToken(headersParts[1])
	if err != nil {
		if errors.Is(auth.ErrInvalidCredential, err) {
			h.sendErrorResponseAndLogError(w, utils.CreateProblemFromRequest(r, http.StatusUnauthorized))
		} else {
			h.sendErrorResponseAndLogError(w, utils.CreateProblemFromError(r, err))
		}
		h.logger.Error()
		return nil, errors.New("can't parse token")
	}
	return claims, nil
}

func (h *HttpHandler) parseContextBody(r *http.Request, commandContext *commands.CommandContext) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read request", zap.Error(err))
		return err
	}
	parseError := json.Unmarshal(body, commandContext.CommandContent)
	if parseError != nil {
		h.logger.Error("Failed to unmarshal body", zap.Error(parseError))
		return err
	}
	return nil
}

func (h *HttpHandler) parseContextParams(r *http.Request, commandContext *commands.CommandContext) {
	for key := range commandContext.CommandParameters {
		urlParam := chi.URLParam(r, key)
		commandContext.CommandParameters[key] = urlParam
	}
}

func (h *HttpHandler) sendSuccessResponse(r *http.Request, w http.ResponseWriter, response interface{}) {
	w.WriteHeader(http.StatusOK)
	if response == nil {
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
