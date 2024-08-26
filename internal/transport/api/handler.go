package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/saikrir/keep-notes/internal/logger"
	"github.com/saikrir/keep-notes/internal/service"
)

type NotesService interface {
	FindNote(context.Context, string) (service.UserNote, error)
	SearchNotes(context.Context, string) ([]service.UserNote, error)
	NewNote(context.Context, service.UserNote) (service.UserNote, error)
	UpdateNote(context.Context, string, service.UserNote) (service.UserNote, error)
	RemoveNote(context.Context, string) (service.UserNote, error)
	GetAllNotes(context.Context) ([]service.UserNote, error)
}

type Handler struct {
	ApiRouter   *http.ServeMux
	Server      *http.Server
	Service     NotesService
	RootContext string
}

func NewHandler(rootContext string, port int, service NotesService) *Handler {
	h := &Handler{
		Service:     service,
		RootContext: rootContext,
	}

	h.ApiRouter = http.NewServeMux()
	rootContextMux := http.NewServeMux()
	rootContextMux.Handle(fmt.Sprintf("%s/", rootContext), http.StripPrefix(rootContext, h.ApiRouter))

	h.mapRoutes()

	h.Server = &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", port),
		Handler:           rootContextMux,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	return h
}

func (h *Handler) mapRoutes() {
	h.ApiRouter.HandleFunc("GET /notes", h.GetAllNotes)
}

func (h *Handler) Serve() error {

	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			logger.Error("Server will stop ", err)
		}
	}()

	// These 3 lines, capture CTRL+C and pass control to the app
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)
	logger.Info("Server has shutdown")
	return nil
}
