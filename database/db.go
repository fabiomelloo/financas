package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Connect() (*sql.DB, error) {

	db, err := sql.Open("sqlite", "./financas.db")
	if err != nil {
		return nil, err
	}

	query := `CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		amount REAL NOT NULL,
		type TEXT NOT NULL,
		category TEXT NOT NULL,
		date DATE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}
