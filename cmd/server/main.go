package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"

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
	db.InitSchema()
	note := service.UserNote{Description: "Sample Note2"}
	if err := db.CreateNote(context.Background(), note); err != nil {
		logger.Error("Failed to insert Row ", err)
	}
	logger.Info("ROWS Inserted ")

	row, err := db.GetNote(context.Background(), "1")
	fmt.Println("Qurty ", row)

	row.Description = "DuFFY"

	results, err := db.SearchNotes(context.Background(), strings.ToLower("Sam"))
	logger.Info("Results ", results)

	return nil
}

func main() {
	Run()
}
