package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"time"
)

type SportTrackerApiGatewayServer struct {
	store  *infrastructure.Store
	router *chi.Mux
}

func NewHTTPServer() *SportTrackerApiGatewayServer {
	server := &SportTrackerApiGatewayServer{}
	server.createRoute()
	return server
}

func (s *SportTrackerApiGatewayServer) createRoute() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	//r.Get("/test", authController.Auth)
}
