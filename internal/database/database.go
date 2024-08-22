package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/saikrir/keep-notes/internal/logger"
	"github.com/saikrir/keep-notes/internal/service"
)

type UserNoteRow struct {
	ID          sql.NullString `db:"ID"`
	Description sql.NullString `db:"DESCRIPTION"`
	CreatedAt   sql.NullString `db:"CREATED_AT"`
	Status      sql.NullString `db:"STATUS"`
}

type Database struct {
	Client *sqlx.DB
}

func NewDatabase() (*Database, error) {
	db, err := sqlx.Connect("sqlite3", "keepnotes.db")
	if err != nil {
		logger.Error("Failed to connect to database ", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Error("Failed to ping database ", err)
		return nil, err
	}

	logger.Info("Connected to Database successfully")
	return &Database{Client: db}, nil
}

func (db *Database) InitSchema() error {
	var schema = `
			DROP TABLE IF EXISTS T_USER_NOTES;
			CREATE TABLE T_USER_NOTES (
				ID INTEGER PRIMARY KEY AUTOINCREMENT,
		    	DESCRIPTION VARCHAR(80)  NOT NULL,
			    CREATED_AT TEXT  DEFAULT CURRENT_TIMESTAMP,
				STATUS  VARCHAR(10) DEFAULT 'ACTIVE'
			); `

	if _, err := db.Client.Exec(schema, nil); err != nil {
		logger.Error("Failed to execute Schema ", err)
		return err
	}
	return nil
}

func (db *Database) CreateNote(ctx context.Context, note service.UserNote) (service.UserNote, error) {
	insertSQL := "INSERT INTO T_USER_NOTES(DESCRIPTION) values(:description) "
	userNoteRow := ToUserNoteRow(note)

	var (
		txn *sqlx.Tx
		err error
	)

	txn = db.Client.MustBegin()
	db.Client.MustExecContext(ctx, insertSQL, userNoteRow.Description.String)
	if err = txn.Commit(); err != nil {
		logger.Error("Failed to commit txn")
	}

	logger.Info("Transaction committed")
	return note, nil
}

func (db *Database) UpdateNote(ctx context.Context, ID string, existingNote service.UserNote) (service.UserNote, error) {
	updateSQL := "UPDATE T_USER_NOTES SET DESCRIPTION = :description, STATUS = :status where ID=:id"
	exisitingRow := ToUserNoteRow(existingNote)
	exisitingRow.ID = sql.NullString{String: ID, Valid: true}

	_, err := db.Client.NamedExecContext(ctx, updateSQL, existingNote)

	if err != nil {
		logger.Error("Failed to update ", err)
		return service.UserNote{}, err
	}

	return ToUserNote(exisitingRow), nil
}

func (db *Database) DeleteNote(ctx context.Context, ID string) (service.UserNote, error) {
	deleteSQL := "DELETE FROM T_USER_NOTES where ID = $1"

	var (
		existingRow service.UserNote
		err         error
	)

	if existingRow, err = db.GetNote(ctx, ID); err != nil {
		logger.Error("Failed to find existing row to delete ", err)
		return existingRow, err
	}

	if _, err := db.Client.ExecContext(ctx, deleteSQL, ID); err != nil {
		logger.Error("Failed to delete ", err)
		return existingRow, err
	}

	return existingRow, nil
}

func (db *Database) GetNote(ctx context.Context, noteId string) (service.UserNote, error) {
	var userNoteRow UserNoteRow

	selectSQL := "SELECT ID, DESCRIPTION, CREATED_AT, STATUS FROM T_USER_NOTES WHERE ID = $1"
	row := db.Client.QueryRowContext(ctx, selectSQL, noteId)

	if err := row.Scan(&userNoteRow.ID, &userNoteRow.Description, &userNoteRow.CreatedAt, &userNoteRow.Status); err != nil {
		logger.Error("Failed to scan row ", err)
		return service.UserNote{}, err
	}
	return ToUserNote(userNoteRow), nil
}

func (db *Database) SearchNote(ctx context.Context, searchTxt string) ([]service.UserNote, error) {
	var (
		searchResults []service.UserNote
		returnRows    []UserNoteRow
		err           error
	)

	searchSQL := "SELECT ID, DESCRIPTION, CREATED_AT, STATUS FROM T_USER_NOTES WHERE lower(DESCRIPTION) LIKE $1"
	if err = db.Client.SelectContext(ctx, &returnRows, searchSQL, "%"+searchTxt+"%"); err != nil {
		logger.Error("Result Err", err)
		return nil, err
	}

	for _, row := range returnRows {
		searchResults = append(searchResults, ToUserNote(row))
	}
	return searchResults, nil

}

func (db *Database) GetAllRows(ctx context.Context) ([]service.UserNote, error) {
	selectSQL := "SELECT ID, DESCRIPTION, CREATED_AT, STATUS FROM T_USER_NOTES"
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
	return retResults, nil
}

func ToUserNoteRow(note service.UserNote) UserNoteRow {
	return UserNoteRow{
		ID:          sql.NullString{String: note.ID, Valid: true},
		Description: sql.NullString{String: note.Description, Valid: true},
		CreatedAt:   sql.NullString{String: note.CreatedAt, Valid: true},
		Status:      sql.NullString{String: note.Status, Valid: true},
	}
}

func ToUserNote(noteRow UserNoteRow) service.UserNote {
	return service.UserNote{
		ID:          noteRow.ID.String,
		Description: noteRow.Description.String,
		CreatedAt:   noteRow.CreatedAt.String,
		Status:      noteRow.Status.String,
	}
}
