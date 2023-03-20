package main

import (
	"encoding/json"
	"fmt"
	"github.com/paulrozhkin/sport-tracker/config"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"go.uber.org/zap"
	"log"
	"net/http"
)

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

	_, err = infrastructure.CreateAndMigrate(&conf.Database)
	if err != nil {
		zap.S().Fatalf("Connection fail due to: %s", err)
	}

	server := NewHTTPServer()

	zap.S().Info("Server started on port: ", conf.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf("localhost:%d", conf.Server.Port), server.router)
	if err != nil {
		zap.S().Error("Server failed due to %v", err)

		panic("Fatal error due to: " + err.Error())
	}
}
