package service

import (
	"context"
	"database/sql"
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
	GetAllRows(context.Context) ([]UserNote, error)
}

var (
	ErrNoNotesFound         = errors.New("no userNotes were found for given ID")
	ErrNoSearchResultsFound = errors.New("no userNotes were found that matched the given criteria")
	ErrNotImplemented       = errors.New("this functionality is currently not implemented")
	ErrFindingNote          = errors.New("failed to find userNote")
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

	note, err := userNoteSvc.store.GetNote(ctx, ID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			logger.Error("No Rows found for NoteId ", ID)
			return UserNote{}, ErrNoNotesFound
		}
		return UserNote{}, err
	}

	return note, nil
}
func (userNoteSvc *UserNotesService) SearchNotes(ctx context.Context, searchTxt string) ([]UserNote, error) {

	results, err := userNoteSvc.store.SearchNote(ctx, searchTxt)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			logger.Error("No Rows found for Search Criteria ", searchTxt)
			return nil, ErrNoSearchResultsFound
		}
		return results, err
	}

	return results, nil
}

func (userNoteSvc *UserNotesService) GetAllNotes(ctx context.Context) ([]UserNote, error) {
	logger.Info("Will get all rows ")
	results, err := userNoteSvc.store.GetAllRows(ctx)
	if err != nil {
		logger.Error("No Rows found in DB ", err)
		if errors.Is(sql.ErrNoRows, err) {
			logger.Error("No Rows found in DB ", err)
			return nil, ErrNoSearchResultsFound
		}
		return results, err
	}

	return results, nil
}

func (userNoteSvc *UserNotesService) NewNote(ctx context.Context, userNote UserNote) (UserNote, error) {
	return userNoteSvc.store.CreateNote(ctx, userNote)
}

func (userNoteSvc *UserNotesService) UpdateNote(ctx context.Context, ID string, userNote UserNote) (UserNote, error) {

	if _, err := userNoteSvc.findExistingRow(ctx, ID); err != nil {
		return UserNote{}, err
	}

	return userNoteSvc.store.UpdateNote(ctx, ID, userNote)
}

func (userNoteSvc *UserNotesService) RemoveNote(ctx context.Context, ID string) (UserNote, error) {

	if _, err := userNoteSvc.findExistingRow(ctx, ID); err != nil {
		return UserNote{}, err
	}
	return userNoteSvc.store.DeleteNote(ctx, ID)
}

func (userNoteSvc *UserNotesService) findExistingRow(ctx context.Context, ID string) (UserNote, error) {

	var (
		aUserNote UserNote
		err       error
	)

	if aUserNote, err = userNoteSvc.store.GetNote(ctx, ID); err != nil {
		logger.Error("Failed to find existing row to Update ", err)
		if errors.Is(err, sql.ErrNoRows) {
			return aUserNote, ErrNoNotesFound
		}
		return aUserNote, err
	}
	return aUserNote, nil
}
