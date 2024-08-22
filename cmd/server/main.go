package main

import (
	"runtime"

	"github.com/saikrir/keep-notes/internal/database"
	"github.com/saikrir/keep-notes/internal/logger"
	"github.com/saikrir/keep-notes/internal/service"
)

func Run() error {
	logger.Info("RUNNING ON ", runtime.GOOS, " Architecutre ", runtime.GOARCH)
	db, err := database.NewDatabase()
	if err != nil {
		panic(err.Error())
	}
	if err = db.InitSchema(); err != nil {
		panic(err.Error())
	}

	noteSvc := service.NewUserNotesService(db)
	logger.Info("Service initialized ", noteSvc)
	return nil
}

func main() {
	Run()
}
