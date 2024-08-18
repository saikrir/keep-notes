package main

import (
	"fmt"

	"github.com/saikrir/keep-notes/internal/logger"
)

func Run() error {
	logger.Info("Hello", "World")
	return nil
}

func main() {
	Run()
	fmt.Println("Hello World")
}
