package service

import (
	"context"
	"errors"
	"time"

	"github.com/saikrir/keep-notes/internal/logger"
)

type UserNote struct {
	ID, Description, Status string
	CreatedAt               time.Time
}

type Store interface {
	GetNote(context.Context, string) (UserNote, error)
	CreateNote(context.Context, UserNote) (UserNote, error)
	UpdateNote(context.Context, string, UserNote) (UserNote, error)
	DeleteNote(context.Context, string) (UserNote, error)
	SearchNote(context.Context, string) ([]UserNote, error)
}

var (
	ErrNoNotesFound         = errors.New("no userNotes were found for given ID")
	ErrNoSearchResultsFound = errors.New("no userNotes were found that matched the given criteria")
	ErrNotImplemented       = errors.New("this functionality is currently not implemented")
	ErrFindingNote          = errors.New("faailed to find userNote")
	ErrCreation             = errors.New("failed to create userNote")
	ErrUpdate               = errors.New("failed to update userNote")
	ErrDelete               = errors.New("failed to delete userNote")
)

type UserNotesService struct {
	store Store
}

func NewUserNotesService(store Store) *UserNotesService {
	return &UserNotesService{
		store: store,
	}
}

func (userNoteSvc *UserNotesService) FindNote(ctx context.Context, ID string) (UserNote, error) {
	logger.Debug("will try and locate Note with Id", ID)

	userNote, err := userNoteSvc.store.GetNote(ctx, ID)
	if err != nil {
		logger.Error("Failed to find Note ", err)
		return UserNote{}, ErrFindingNote
	}

	return userNote, nil
}
func (userNoteSvc *UserNotesService) SearchNotes(ctx context.Context, searchTxt string) ([]UserNote, error) {
	return userNoteSvc.store.SearchNote(ctx, searchTxt)
}

func (userNoteSvc *UserNotesService) NewNote(ctx context.Context, userNote UserNote) (UserNote, error) {
	return userNoteSvc.store.CreateNote(ctx, userNote)
}

func (userNoteSvc *UserNotesService) UpdateNote(ctx context.Context, ID string, userNote UserNote) (UserNote, error) {
	return userNoteSvc.store.UpdateNote(ctx, ID, userNote)
}

func (userNoteSvc *UserNotesService) RemoveNote(ctx context.Context, ID string) (UserNote, error) {
	return userNoteSvc.store.DeleteNote(ctx, ID)
}
