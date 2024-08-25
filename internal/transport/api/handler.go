package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/saikrir/keep-notes/internal/service"
)

type NotesService interface {
	FindNote(context.Context, string) (service.UserNote, error)
	SearchNotes(context.Context, string) ([]service.UserNote, error)
	NewNote(context.Context, service.UserNote) (service.UserNote, error)
	UpdateNote(context.Context, string, service.UserNote) (service.UserNote, error)
	RemoveNote(context.Context, string) (service.UserNote, error)
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
	h.ApiRouter.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
}

func (h *Handler) Serve() error {
	if err := h.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
