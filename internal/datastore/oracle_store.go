package datastore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/saikrir/keep-notes/internal/env"
	"github.com/saikrir/keep-notes/internal/logger"
	"github.com/saikrir/keep-notes/internal/service"

	go_ora "github.com/sijms/go-ora/v2"
)

type UserNoteRow struct {
	ID          sql.NullString `db:"ID"`
	Description sql.NullString `db:"DESCRIPTION"`
	CreatedAt   sql.NullTime   `db:"CREATED_AT"`
	Status      sql.NullString `db:"STATUS"`
}

type OracleStore struct {
	Client *sqlx.DB
}

func NewOracleStore() (*OracleStore, error) {

	// Need following env vars
	// $DB_HOST, $DB_PORT, $DB_NAME, $DB_USER, $DB_PASS\
	connStr := go_ora.BuildUrl(
		env.GetEnvValAsString("DB_HOST"),
		env.GetEnvValAsNumber("DB_PORT"),
		env.GetEnvValAsString("DB_NAME"),
		env.GetEnvValAsString("DB_USER"),
		env.GetEnvValAsString("DB_PASS"),
		nil,
	)

	logger.Debug("DB Str", connStr)

	conn, err := sqlx.Connect("oracle", connStr)

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to db with conn str %s", connStr), err)
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		logger.Error(fmt.Sprintf("Failed to ping db with conn str %s", connStr), err)
		return nil, err
	}

	logger.Debug("Connect to db")

	return &OracleStore{Client: conn}, nil
}

func (db *OracleStore) CreateNote(ctx context.Context, note service.UserNote) (service.UserNote, error) {

	var (
		txn    *sqlx.Tx
		err    error
		result sql.Result
	)

	insertSQL := fmt.Sprintf("INSERT INTO %s.T_USER_NOTES(DESCRIPTION) values(:1) ", env.GetEnvValAsString("DB_USER"))
	userNoteRow := ToUserNoteRow(note)

	txn = db.Client.MustBegin()
	result, err = db.Client.ExecContext(ctx, insertSQL, userNoteRow.Description.String)
	if err != nil {
		logger.Error("Failed to execute insert statement ", err)
		return service.UserNote{}, fmt.Errorf("error occured when creating a new record %w", err)
	}

	if err = txn.Commit(); err != nil {
		logger.Error("Failed to commit txn ", err)
		return service.UserNote{}, fmt.Errorf("error occured when creating a new record %w", err)
	}
	lastID, err := result.RowsAffected()
	if err != nil {
		logger.Error("Failed to get last Insert Id")
	}
	logger.Info("Transaction committed, new Row created with ID ", lastID)
	return note, nil
}

func (db *OracleStore) UpdateNote(ctx context.Context, ID string, existingNote service.UserNote) (service.UserNote, error) {
	var (
		txn    *sqlx.Tx
		err    error
		result sql.Result
	)

	updateSQL := fmt.Sprintf("UPDATE %s.T_USER_NOTES SET DESCRIPTION = :1, STATUS = :2 where ID=:3", env.GetEnvValAsString("DB_USER"))
	exisitingRow := ToUserNoteRow(existingNote)
	exisitingRow.ID = sql.NullString{String: ID, Valid: true}

	txn = db.Client.MustBegin()

	result, err = db.Client.ExecContext(ctx, updateSQL, existingNote.Description, exisitingRow.Status, ID)

	if err != nil {
		logger.Error("Failed to execute Update statement ", err)
		return service.UserNote{}, fmt.Errorf("error occured when updating new record %w", err)
	}

	if err = txn.Commit(); err != nil {
		logger.Error("Failed to commit txn ", err)
		return service.UserNote{}, fmt.Errorf("error occured when creating a new record %w", err)
	}

	nRowsAffected, _ := result.RowsAffected()

	logger.Info(nRowsAffected, " Rows were updated ")

	return ToUserNote(exisitingRow), nil
}

func (db *OracleStore) DeleteNote(ctx context.Context, ID string) (service.UserNote, error) {
	deleteSQL := fmt.Sprintf("DELETE FROM %s.T_USER_NOTES where ID = :1", env.GetEnvValAsString("DB_USER"))

	var (
		txn         *sqlx.Tx
		existingRow service.UserNote
		err         error
		result      sql.Result
	)

	txn = db.Client.MustBegin()

	result, err = db.Client.ExecContext(ctx, deleteSQL, ID)

	if err != nil {
		logger.Error("failed to execute Delete statement ", err)
		return service.UserNote{}, fmt.Errorf("error occured when deleting record %w", err)
	}

	if err = txn.Commit(); err != nil {
		logger.Error("Failed to commit txn ", err)
		return service.UserNote{}, fmt.Errorf("error occured when deleting record %w", err)
	}

	numRowsAffected, _ := result.RowsAffected()
	logger.Info(numRowsAffected, " Rows were delete")
	return existingRow, nil
}

func (db *OracleStore) GetNote(ctx context.Context, noteId string) (service.UserNote, error) {
	var userNoteRow UserNoteRow

	selectSQL := fmt.Sprintf("SELECT ID, DESCRIPTION, CREATED_AT, STATUS FROM %s.T_USER_NOTES WHERE ID = :1", env.GetEnvValAsString("DB_USER"))
	row := db.Client.QueryRowContext(ctx, selectSQL, noteId)

	if err := row.Scan(&userNoteRow.ID, &userNoteRow.Description, &userNoteRow.CreatedAt, &userNoteRow.Status); err != nil {
		logger.Error("Failed to scan row ", err)
		return service.UserNote{}, err
	}
	return ToUserNote(userNoteRow), nil
}

func (db *OracleStore) SearchNote(ctx context.Context, searchTxt string) ([]service.UserNote, error) {
	var (
		searchResults []service.UserNote
		returnRows    []UserNoteRow
		err           error
	)

	searchSQL := fmt.Sprintf("SELECT ID, DESCRIPTION, CREATED_AT, STATUS FROM %s.T_USER_NOTES WHERE lower(DESCRIPTION) LIKE $1 ", env.GetEnvValAsString("DB_USER"))
	if err = db.Client.SelectContext(ctx, &returnRows, searchSQL, "%"+searchTxt+"%"); err != nil {
		logger.Error("error executing search Query ", err)
		return nil, err
	}

	for _, row := range returnRows {
		searchResults = append(searchResults, ToUserNote(row))
	}
	return searchResults, nil

}

func (db *OracleStore) GetAllRows(ctx context.Context) ([]service.UserNote, error) {
	selectSQL := fmt.Sprintf("SELECT ID, DESCRIPTION, CREATED_AT, STATUS FROM %s.T_USER_NOTES ORDER BY CREATED_AT DESC", env.GetEnvValAsString("DB_USER"))

	var (
		err        error
		rows       []UserNoteRow
		retResults []service.UserNote
	)

	if err = db.Client.SelectContext(ctx, &rows, selectSQL); err != nil {
		logger.Error("Failed to lookup all rows ", err)
		return nil, err
	}

	for _, row := range rows {
		retResults = append(retResults, ToUserNote(row))
	}
	logger.Debug(fmt.Sprintf("Found %d rows", len(retResults)))
	return retResults, nil
}

func ToUserNoteRow(note service.UserNote) UserNoteRow {
	return UserNoteRow{
		ID:          sql.NullString{String: note.ID, Valid: true},
		Description: sql.NullString{String: note.Description, Valid: true},
		CreatedAt:   sql.NullTime{Time: note.CreatedAt, Valid: true},
		Status:      sql.NullString{String: note.Status, Valid: true},
	}
}

func ToUserNote(noteRow UserNoteRow) service.UserNote {
	return service.UserNote{
		ID:          noteRow.ID.String,
		Description: noteRow.Description.String,
		CreatedAt:   noteRow.CreatedAt.Time,
		Status:      noteRow.Status.String,
	}
}
