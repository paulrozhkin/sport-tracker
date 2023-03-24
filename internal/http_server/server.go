package http_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/paulrozhkin/sport-tracker/config"
	"github.com/paulrozhkin/sport-tracker/internal/http_server/middlewares"
	"github.com/paulrozhkin/sport-tracker/internal/http_server/routes"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

func NewHTTPServer(lc fx.Lifecycle,
	cfg *config.Configuration,
	handler http.Handler,
	log *zap.SugaredLogger) *http.Server {
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: serverAddr, Handler: handler}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			go func() {
				log.Infof("Server started on addr: %s", srv.Addr)
				serveError := srv.Serve(ln)
				if serveError != nil {
					if errors.Is(http.ErrServerClosed, serveError) {
						log.Info("Server stopped")
					} else {
						log.Errorf("Server failed due to %v", zap.Error(serveError))
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func NewServerRoute(routes []routes.Route,
	logger *zap.SugaredLogger,
	config *config.Configuration,
	authMiddleware *middlewares.AuthAuthMiddleware) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	timeout := time.Second * time.Duration(config.Server.RequestTimeoutSeconds)
	r.Use(middleware.Timeout(timeout))
	r.Use(authMiddleware.Handle)
	for _, route := range routes {
		handler := &HttpHandler{routeHandler: route, logger: logger}
		r.Method(route.Method(), route.Pattern(), handler)
	}
	return r
}
