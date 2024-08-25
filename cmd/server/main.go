package main

import (
	"context"
	"runtime"

	"github.com/saikrir/keep-notes/internal/datastore"
	"github.com/saikrir/keep-notes/internal/logger"
	"github.com/saikrir/keep-notes/internal/service"
)

func Run() error {
	logger.Info("RUNNING ON ", runtime.GOOS, " Architecutre ", runtime.GOARCH)
	db, err := datastore.NewOracleStore()

	if err != nil {
		logger.Error("Failed to Connect to DB ", err)
		panic(err.Error())
	}

	appNote := service.UserNote{
		Description: "Sample",
		Status:      "Active",
	}

	if _, err := db.CreateNote(context.Background(), appNote); err != nil {
		logger.Error("Failed to create row ", err)
	}

	return nil
}

func main() {
	Run()
}
