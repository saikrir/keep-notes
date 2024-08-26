package api

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	allNotes, err := h.Service.GetAllNotes(r.Context())
	if err != nil {
		logger.Error("Failed to Get All Notes ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	retNotes := make([]Note, len(allNotes))
	for i, userNote := range allNotes {
		retNotes[i] = ToNote(userNote)
	}
	mustMarshall(retNotes, w)
}

func (h *Handler) GetNoteById(w http.ResponseWriter, r *http.Request) {

	var (
		err      error
		userNote service.UserNote
	)

	noteId := r.PathValue("noteId")
	if userNote, err = h.Service.FindNote(r.Context(), noteId); err != nil {
		if errors.Is(err, service.ErrNoNotesFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mustMarshall(ToNote(userNote), w)
}

func (h *Handler) PostNote(w http.ResponseWriter, r *http.Request) {

	var (
		aNote    Note
		err      error
		userNote service.UserNote
	)
	if err := json.NewDecoder(r.Body).Decode(&aNote); err != nil {
		logger.Error("Failed read request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	userNote = ToUserNote(aNote)

	if userNote, err = h.Service.NewNote(r.Context(), userNote); err != nil {
		logger.Error("Failed create note ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mustMarshall(userNote, w)
}

func (h *Handler) PutNote(w http.ResponseWriter, r *http.Request) {
	var (
		aNote    Note
		err      error
		userNote service.UserNote
	)
	if err := json.NewDecoder(r.Body).Decode(&aNote); err != nil {
		logger.Error("Failed read request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	noteId := r.PathValue("noteId")
	userNote = ToUserNote(aNote)

	if userNote, err = h.Service.UpdateNote(r.Context(), noteId, userNote); err != nil {

		if errors.Is(err, service.ErrNoNotesFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		logger.Error("Failed Update note ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mustMarshall(userNote, w)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		userNote service.UserNote
	)

	noteId := r.PathValue("noteId")

	if userNote, err = h.Service.RemoveNote(r.Context(), noteId); err != nil {

		if errors.Is(err, service.ErrNoNotesFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		logger.Error("Failed Delete note ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mustMarshall(userNote, w)
}

func (h *Handler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	resp := make(map[string]string)
	resp["message"] = fmt.Sprintf("%s is not supported by the api ", r.RequestURI)
	mustMarshall(resp, w)
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

func mustMarshall(payload any, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error("failed to marshall json ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
