package main

import (
	"runtime"

	"github.com/saikrir/keep-notes/internal/database"
	"github.com/saikrir/keep-notes/internal/logger"
)

func Run() error {
	logger.Info("RUNNING ON ", runtime.GOOS, " Architecutre ", runtime.GOARCH)
	database.NewDatabase()
	return nil
}

func main() {
	Run()
}
