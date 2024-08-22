package service

import (
	"context"
	"errors"

	"github.com/saikrir/keep-notes/internal/logger"
)

type UserNote struct {
	ID, Description, Status string
	CreatedAt               string
}

type Store interface {
	GetNote(context.Context, string) (UserNote, error)
	CreateNote(context.Context, UserNote) (UserNote, error)
	UpdateNote(context.Context, string, UserNote) (UserNote, error)
	DeleteNote(context.Context, string) (UserNote, error)
	SearchNote(context.Context, string) ([]UserNote, error)
}

var (
	ErrNoNotesFound         = errors.New("No userNotes were found for given ID")
	ErrNoSearchResultsFound = errors.New("No userNotes were found that matched the given criteria")
	ErrNotImplemented       = errors.New("This functionality is currently not implemented")
	ErrFindingNote          = errors.New("Failed to find userNote")
	ErrCreation             = errors.New("Failed to create userNote")
	ErrUpdate               = errors.New("Failed to update userNote")
	ErrDelete               = errors.New("Failed to delete userNote")
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
	return nil, ErrNotImplemented
}

func (userNoteSvc *UserNotesService) NewNote(ctx context.Context, userNote UserNote) (UserNote, error) {
	return UserNote{}, ErrNotImplemented
}

func (userNoteSvc *UserNotesService) UpdateNote(ctx context.Context, ID string, userNote UserNote) (UserNote, error) {
	return UserNote{}, ErrNotImplemented
}

func (userNoteSvc *UserNotesService) RemoveNote(ctx context.Context, ID string) (UserNote, error) {
	return UserNote{}, ErrNotImplemented
}
