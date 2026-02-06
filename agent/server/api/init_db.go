package api

import (
	"database/sql"

	"github.com/RZhurakovskiy/agent/server/db"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := sqlDB.Exec(db.SchemaSQL); err != nil {
		sqlDB.Close()
		return nil, err
	}
	return sqlDB, nil
}
