package main

import (
	"runtime"

	"github.com/saikrir/keep-notes/internal/logger"
)

func Run() error {
	logger.Info("RUNNING ON ", runtime.GOOS, " Architecutre ", runtime.GOARCH)
	return nil
}

func main() {
	Run()
}
