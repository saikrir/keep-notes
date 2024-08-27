package main

import (
	"runtime"

	"github.com/saikrir/keep-notes/internal/datastore"
	"github.com/saikrir/keep-notes/internal/env"
	"github.com/saikrir/keep-notes/internal/logger"
	"github.com/saikrir/keep-notes/internal/service"
	"github.com/saikrir/keep-notes/internal/transport/api"
)

const ApiRootContext = "/v1/notesvc"

func Run() error {
	logger.Info("RUNNING ON ", runtime.GOOS, " Architecutre ", runtime.GOARCH)
	dataStore, err := datastore.NewSQLliteStore(true)
	if err != nil {
		logger.Error("Failed to Connect to DB ", err)
		return err
	}

	noteSvc := service.NewUserNotesService(dataStore)
	handler := api.NewHandler("/v1/notesvc", env.GetEnvValAsNumber("API_PORT"), noteSvc)
	return handler.Serve()
}

func main() {
	if err := Run(); err != nil {
		logger.Error("App will halt ", err)
		panic(err.Error())
	}
}
