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
		payer TEXT DEFAULT '',
		date DATE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME DEFAULT NULL
	)`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	// Migration simplificada: tentar adicionar a coluna se ela não existir
	// SQLite não suporta "IF NOT EXISTS" em ADD COLUMN nativamente de forma simples em todas as versões,
	// mas podemos tentar rodar o ALTER e ignorar erro de coluna duplicada.
	migration := `ALTER TABLE expenses ADD COLUMN payer TEXT DEFAULT ''`
	db.Exec(migration) // Ignoramos erro aqui propositalmente (ex: coluna já existe)

	return db, nil
}
