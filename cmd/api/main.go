package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/paulrozhkin/sport-tracker/config"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type sportTrackerApiGatewayServer struct {
	store *infrastructure.Store
}

func newServer(configurations *config.Configurations) *sportTrackerApiGatewayServer {
	store, err := infrastructure.CreateAndMigrate(&configurations.Database)
	if err != nil {
		zap.S().Fatalf("Connection fail due to: %s", err)
	}
	return &sportTrackerApiGatewayServer{store: store}
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	conf, err := config.LoadConfigurations()
	if err != nil {
		zap.S().Fatal("Failed to load config due to:", err)
	}

	confString, _ := json.MarshalIndent(conf, "", " ")
	zap.S().Info("Configuration:\n", string(confString))
	server := newServer(conf)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/test", server.testHandler)

	zap.S().Info("Server started on port: ", conf.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf("localhost:%d", conf.Server.Port), r)
	if err != nil {
		zap.S().Error("Server failed due to %v", err)

		panic("Fatal error due to: " + err.Error())
	}
}

func (ts *sportTrackerApiGatewayServer) testHandler(w http.ResponseWriter, req *http.Request) {
	zap.S().Info("test")
	w.Write([]byte("hi"))
}
