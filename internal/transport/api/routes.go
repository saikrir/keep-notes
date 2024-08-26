package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/saikrir/keep-notes/internal/logger"
	"github.com/saikrir/keep-notes/internal/service"
)

type Note struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	Status      string `json:"status"`
}

func (h *Handler) GetAllNotes(resp http.ResponseWriter, r *http.Request) {
	logger.Info("Will habndle")
	allRows, err := h.Service.GetAllNotes(r.Context())
	if err != nil {
		logger.Error("Failed to Get All Notes ", err)
		resp.WriteHeader(http.StatusInternalServerError)
	}
	mustMarshall(allRows, resp)
}

func ToNote(userNote service.UserNote) Note {
	return Note{
		ID:          userNote.ID,
		Description: userNote.Description,
		CreatedAt:   userNote.CreatedAt.Format(time.RFC3339),
		Status:      userNote.Status,
	}
}

func ToUserNote(note Note) service.UserNote {
	return service.UserNote{
		ID:          note.ID,
		Description: note.Description,
		CreatedAt:   time.Now(),
		Status:      note.Status,
	}
}

func mustMarshall(payload any, resp http.ResponseWriter) {
	if err := json.NewEncoder(resp).Encode(payload); err != nil {
		logger.Error("failed to marshall json ", err)
		resp.WriteHeader(http.StatusInternalServerError)
	}
}
