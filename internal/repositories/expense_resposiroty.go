package repositories

import (
	"database/sql"
	"financas/internal/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(expense *models.Expense) error {
	query := `INSERT INTO expenses (description, amount, type, category, date) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, expense.Description, expense.Amount, expense.Type, expense.Category, expense.Date)
	return err
}

func (r *ExpenseRepository) FindAll() ([]models.Expense, error) {
	query := `SELECT * FROM expenses`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(&expense.ID, &expense.Description, &expense.Amount, &expense.Type, &expense.Category, &expense.Date, &expense.CreatedAt, &expense.UpdatedAt, &expense.DeletedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (r *ExpenseRepository) FindByID(id int) (*models.Expense, error) {
	query := `SELECT * FROM expenses WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var expense models.Expense
	if err := row.Scan(&expense.ID, &expense.Description, &expense.Amount, &expense.Type, &expense.Category, &expense.Date, &expense.CreatedAt, &expense.UpdatedAt, &expense.DeletedAt); err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) Update(expense *models.Expense) error {
	query := `UPDATE expenses SET description = ?, amount = ?, type = ?, category = ?, date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, expense.Description, expense.Amount, expense.Type, expense.Category, expense.Date, expense.ID)
	return err
}

func (r *ExpenseRepository) Delete(id int) error {
	query := `UPDATE expenses SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
