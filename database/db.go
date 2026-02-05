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

	// Tabela de expenses (existente)
	expensesTable := `CREATE TABLE IF NOT EXISTS expenses (
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
	if _, err = db.Exec(expensesTable); err != nil {
		return nil, err
	}

	// Migration: adicionar coluna payer se não existir
	db.Exec(`ALTER TABLE expenses ADD COLUMN payer TEXT DEFAULT ''`)

	// ============================================
	// NOVAS TABELAS - Sistema de Rateio + Gamificação
	// ============================================

	// Tabela de usuários/membros da equipe
	usersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		points INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	if _, err = db.Exec(usersTable); err != nil {
		return nil, err
	}

	// Tabela de compras de lanche
	purchasesTable := `CREATE TABLE IF NOT EXISTS purchases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		amount REAL NOT NULL,
		date DATE NOT NULL,
		month TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`
	if _, err = db.Exec(purchasesTable); err != nil {
		return nil, err
	}

	// Tabela de balanço mensal
	monthlyBalancesTable := `CREATE TABLE IF NOT EXISTS monthly_balances (
		user_id INTEGER NOT NULL,
		month TEXT NOT NULL,
		total_paid REAL DEFAULT 0,
		share_value REAL DEFAULT 0,
		balance REAL DEFAULT 0,
		PRIMARY KEY (user_id, month),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`
	if _, err = db.Exec(monthlyBalancesTable); err != nil {
		return nil, err
	}

	// Tabela de conquistas (badges)
	achievementsTable := `CREATE TABLE IF NOT EXISTS achievements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		icon TEXT
	)`
	if _, err = db.Exec(achievementsTable); err != nil {
		return nil, err
	}

	// Tabela de conquistas dos usuários
	userAchievementsTable := `CREATE TABLE IF NOT EXISTS user_achievements (
		user_id INTEGER NOT NULL,
		achievement_id INTEGER NOT NULL,
		month TEXT NOT NULL,
		awarded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, achievement_id, month),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (achievement_id) REFERENCES achievements(id)
	)`
	if _, err = db.Exec(userAchievementsTable); err != nil {
		return nil, err
	}

	// Seed de conquistas padrão (insert or ignore para não duplicar)
	seedAchievements := `
		INSERT OR IGNORE INTO achievements (name, description, icon) VALUES
		('Mecenas', 'Maior crédito do mês', '🏅'),
		('Contador', 'Mais compras no mês', '🧾'),
		('Equilibrado', 'Saldo próximo de zero', '🔄'),
		('Mão Aberta', 'Maior gasto individual', '💸'),
		('Caloteiro Simpático', 'Maior débito do mês', '🐢')
	`
	db.Exec(seedAchievements)

	return db, nil
}
