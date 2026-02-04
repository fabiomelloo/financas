package repositories

import (
	"database/sql"
	"financas/internal/models"
	"fmt"
	"time"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(expense *models.Expense) error {
	query := `INSERT INTO expenses (description, amount, type, category, payer, date) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, expense.Description, expense.Amount, expense.Type, expense.Category, expense.Payer, expense.Date)
	return err
}

func (r *ExpenseRepository) FindAll() ([]models.Expense, error) {
	// Eliminar despesas deletadas
	query := `SELECT id, description, amount, type, category, payer, date FROM expenses WHERE deleted_at IS NULL`
	fmt.Println("Executing FindAll query...")
	rows, err := r.db.Query(query)
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		var dateStr string
		if err := rows.Scan(&expense.ID, &expense.Description, &expense.Amount, &expense.Type, &expense.Category, &expense.Payer, &dateStr); err != nil {
			fmt.Printf("Scan error: %v\n", err)
			return nil, err
		}
		// Parse da data vinda do SQLite
		if dateStr != "" {
			expense.Date, _ = time.Parse("2006-01-02T15:04:05Z", dateStr)
			if expense.Date.IsZero() {
				expense.Date, _ = time.Parse("2006-01-02 15:04:05", dateStr)
			}
			if expense.Date.IsZero() {
				expense.Date, _ = time.Parse("2006-01-02", dateStr)
			}
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (r *ExpenseRepository) FindByID(id int) (*models.Expense, error) {
	query := `SELECT id, description, amount, type, category, payer, date, created_at, updated_at, deleted_at FROM expenses WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var expense models.Expense
	var dateStr string
	var createdAtStr, updatedAtStr, deletedAtStr sql.NullString

	if err := row.Scan(&expense.ID, &expense.Description, &expense.Amount, &expense.Type, &expense.Category, &expense.Payer, &dateStr, &createdAtStr, &updatedAtStr, &deletedAtStr); err != nil {
		fmt.Printf("FindByID Scan Error: %v\n", err)
		return nil, err
	}

	// Parse date
	if dateStr != "" {
		expense.Date, _ = time.Parse("2006-01-02T15:04:05Z", dateStr)
		if expense.Date.IsZero() {
			expense.Date, _ = time.Parse("2006-01-02 15:04:05", dateStr)
		}
		if expense.Date.IsZero() {
			expense.Date, _ = time.Parse("2006-01-02", dateStr)
		}
	}
	return &expense, nil
}

func (r *ExpenseRepository) Update(expense *models.Expense) error {
	query := `UPDATE expenses SET description = ?, amount = ?, type = ?, category = ?, payer = ?, date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, expense.Description, expense.Amount, expense.Type, expense.Category, expense.Payer, expense.Date, expense.ID)
	return err
}

func (r *ExpenseRepository) Delete(id int) error {
	query := `UPDATE expenses SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// GetSummary retorna métricas de resumo total
func (r *ExpenseRepository) GetSummary() (float64, float64, float64, error) {
	query := `SELECT 
		COALESCE(SUM(CASE WHEN type = 'receita' THEN amount ELSE 0 END), 0) as total_income,
		COALESCE(SUM(CASE WHEN type = 'despesa' THEN amount ELSE 0 END), 0) as total_expense
	FROM expenses WHERE deleted_at IS NULL`

	var income, expense float64
	err := r.db.QueryRow(query).Scan(&income, &expense)
	if err != nil {
		return 0, 0, 0, err
	}
	balance := income - expense
	return income, expense, balance, nil
}

// CategoryMetric representa estatísticas de agrupamento
type CategoryMetric struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Type     string  `json:"type"`
}

// GetCategoryBreakdown returns expenses grouped by category
func (r *ExpenseRepository) GetCategoryBreakdown() ([]CategoryMetric, error) {
	query := `SELECT category, type, SUM(amount) as total 
			  FROM expenses 
			  WHERE deleted_at IS NULL 
			  GROUP BY category, type 
			  ORDER BY total DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []CategoryMetric
	for rows.Next() {
		var m CategoryMetric
		if err := rows.Scan(&m.Category, &m.Type, &m.Total); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// MonthlyMetric representa agregação mensal
type MonthlyMetric struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

// GetMonthlyBreakdown returns expenses grouped by month
func (r *ExpenseRepository) GetMonthlyBreakdown() ([]MonthlyMetric, error) {
	query := `SELECT 
		substr(date, 1, 7) as month,
		COALESCE(SUM(CASE WHEN type = 'receita' THEN amount ELSE 0 END), 0) as income,
		COALESCE(SUM(CASE WHEN type = 'despesa' THEN amount ELSE 0 END), 0) as expense
	FROM expenses 
	WHERE deleted_at IS NULL 
	GROUP BY month 
	ORDER BY month DESC 
	LIMIT 12`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []MonthlyMetric
	for rows.Next() {
		var m MonthlyMetric
		if err := rows.Scan(&m.Month, &m.Income, &m.Expense); err != nil {
			return nil, err
		}
		m.Balance = m.Income - m.Expense
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// TypeMetric representa totais por tipo
type TypeMetric struct {
	Type  string  `json:"type"`
	Total float64 `json:"total"`
	Count int     `json:"count"`
}

// GetTypeBreakdown returns totals by type
func (r *ExpenseRepository) GetTypeBreakdown() ([]TypeMetric, error) {
	query := `SELECT 
		type,
		SUM(amount) as total,
		COUNT(*) as count
	FROM expenses 
	WHERE deleted_at IS NULL 
	GROUP BY type 
	ORDER BY total DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []TypeMetric
	for rows.Next() {
		var m TypeMetric
		if err := rows.Scan(&m.Type, &m.Total, &m.Count); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// GetTopExpenses retorna as top N despesas
func (r *ExpenseRepository) GetTopExpenses(limit int) ([]models.Expense, error) {
	query := `SELECT id, description, amount, type, category, payer, date 
			  FROM expenses 
			  WHERE deleted_at IS NULL 
			  ORDER BY amount DESC 
			  LIMIT ?`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		var dateStr string

		if err := rows.Scan(&expense.ID, &expense.Description, &expense.Amount, &expense.Type, &expense.Category, &expense.Payer, &dateStr); err != nil {
			return nil, err
		}
		if dateStr != "" {
			expense.Date, _ = time.Parse("2006-01-02", dateStr)
			if expense.Date.IsZero() {
				expense.Date, _ = time.Parse("2006-01-02 15:04:05", dateStr)
			}
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}
