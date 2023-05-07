package http_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/paulrozhkin/sport-tracker/config"
	"github.com/paulrozhkin/sport-tracker/internal/http_server/routes"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"moul.io/chizap"
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
	loggerRaw *zap.Logger,
	config *config.Configuration,
	tokenService *services.TokenService) http.Handler {
	r := chi.NewRouter()
	if config.Server.DisableCORS {
		logger.Info("CORS disabled")
		// Basic CORS
		// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
		corsOptions := cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}

		r.Use(cors.Handler(corsOptions))
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(chizap.New(loggerRaw, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))
	r.Use(middleware.Recoverer)

	timeout := time.Second * time.Duration(config.Server.RequestTimeoutSeconds)
	r.Use(middleware.Timeout(timeout))
	for _, route := range routes {
		handler := &HttpHandler{routeHandler: route, logger: logger, tokenService: tokenService}
		r.Method(route.Method(), route.Pattern(), handler)
	}
	return r
}
