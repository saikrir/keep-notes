package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/saikrir/keep-notes/internal/logger"
	go_ora "github.com/sijms/go-ora/v2"
)

type Database struct {
	Client *sqlx.DB
}

func NewDatabase() (*Database, error) {

	connStr := go_ora.BuildUrl("localhost", 1521, "XEPDB1", "APP_USER", "tempid1", nil)
	logger.Info("DB Str", connStr)

	conn, err := sqlx.Connect("oracle", connStr)

	if err != nil {
		logger.Error("Failed to connect to db ", err)
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		logger.Error("Failed to Ping db ", err)
		return nil, err
	}

	logger.Info("Connect to db")

	temp := make(map[string]any)
	if err = conn.Select(temp, "SELECT * FROM T_USER_NOTES"); err != nil {
		logger.Error("Failed to run query ", err)
	}

	logger.Info("DT ", temp)

	return &Database{Client: conn}, nil
}
