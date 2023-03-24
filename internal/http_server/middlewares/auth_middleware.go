package middlewares

import (
	"context"
	"errors"
	"github.com/docker/distribution/registry/auth"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	ClaimsContextKey    = "claims"
	AuthorizationHeader = "Authorization"
)

type AuthAuthMiddleware struct {
	tokenService *services.TokenService
	logger       *zap.SugaredLogger
}

func NewAuthMiddleware(tokenService *services.TokenService, logger *zap.SugaredLogger) (*AuthAuthMiddleware, error) {
	return &AuthAuthMiddleware{tokenService: tokenService, logger: logger}, nil
}

func (a *AuthAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: need refactoring. Maybe move http_handler or use chi features
		if r.RequestURI == "/auth" || r.RequestURI == "/register" {
			next.ServeHTTP(w, r.WithContext(r.Context()))
			return
		}
		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			a.sendErrorResponseAndLogError(w, utils.CreateProblemFromRequest(r, http.StatusUnauthorized))
			return
		}
		headersParts := strings.Split(authHeader, " ")
		if len(headersParts) != 2 {
			a.sendErrorResponseAndLogError(w, utils.CreateProblemFromRequest(r, http.StatusUnauthorized))
			return
		}
		claims, err := a.tokenService.ParseToken(headersParts[1])
		if err != nil {
			if errors.Is(auth.ErrInvalidCredential, err) {
				a.sendErrorResponseAndLogError(w, utils.CreateProblemFromRequest(r, http.StatusUnauthorized))
			} else {
				a.sendErrorResponseAndLogError(w, utils.CreateProblemFromError(r, err))
			}
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthAuthMiddleware) sendErrorResponseAndLogError(w http.ResponseWriter, problemDetails dto.ProblemDetails) {
	err := utils.WriteProblemToResponse(w, problemDetails)
	if err != nil {
		a.logger.Error("Failed to send error response", zap.Error(err))
	}
}
