package repositories

import (
	"database/sql"
	"financas/internal/models"
	"time"
)

type PurchaseRepository struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

// Create insere uma nova compra de lanche
func (r *PurchaseRepository) Create(purchase *models.Purchase) error {
	// Extrair mês da data (formato "2026-02")
	purchase.Month = purchase.Date.Format("2006-01")

	query := `INSERT INTO purchases (user_id, amount, date, month) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, purchase.UserID, purchase.Amount, purchase.Date.Format("2006-01-02"), purchase.Month)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	purchase.ID = int(id)
	return nil
}

// FindAll retorna todas as compras com nome do usuário
func (r *PurchaseRepository) FindAll() ([]models.Purchase, error) {
	query := `
		SELECT p.id, p.user_id, u.name, p.amount, p.date, p.month, p.created_at
		FROM purchases p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.date DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var p models.Purchase
		var dateStr, createdAtStr string
		if err := rows.Scan(&p.ID, &p.UserID, &p.UserName, &p.Amount, &dateStr, &p.Month, &createdAtStr); err != nil {
			return nil, err
		}
		p.Date, _ = time.Parse("2006-01-02", dateStr)
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		purchases = append(purchases, p)
	}
	return purchases, nil
}

// FindByMonth retorna compras de um mês específico
func (r *PurchaseRepository) FindByMonth(month string) ([]models.Purchase, error) {
	query := `
		SELECT p.id, p.user_id, u.name, p.amount, p.date, p.month, p.created_at
		FROM purchases p
		JOIN users u ON p.user_id = u.id
		WHERE p.month = ?
		ORDER BY p.date DESC
	`
	rows, err := r.db.Query(query, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var p models.Purchase
		var dateStr, createdAtStr string
		if err := rows.Scan(&p.ID, &p.UserID, &p.UserName, &p.Amount, &dateStr, &p.Month, &createdAtStr); err != nil {
			return nil, err
		}
		p.Date, _ = time.Parse("2006-01-02", dateStr)
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		purchases = append(purchases, p)
	}
	return purchases, nil
}

// FindByID retorna uma compra pelo ID
func (r *PurchaseRepository) FindByID(id int) (*models.Purchase, error) {
	query := `
		SELECT p.id, p.user_id, u.name, p.amount, p.date, p.month, p.created_at
		FROM purchases p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`
	row := r.db.QueryRow(query, id)

	var p models.Purchase
	var dateStr, createdAtStr string
	if err := row.Scan(&p.ID, &p.UserID, &p.UserName, &p.Amount, &dateStr, &p.Month, &createdAtStr); err != nil {
		return nil, err
	}
	p.Date, _ = time.Parse("2006-01-02", dateStr)
	p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	return &p, nil
}

// Delete remove uma compra
func (r *PurchaseRepository) Delete(id int) error {
	query := `DELETE FROM purchases WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// GetMonthlyTotalByUser retorna o total pago por cada usuário em um mês
func (r *PurchaseRepository) GetMonthlyTotalByUser(month string) (map[int]float64, error) {
	query := `
		SELECT user_id, SUM(amount) as total
		FROM purchases
		WHERE month = ?
		GROUP BY user_id
	`
	rows, err := r.db.Query(query, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totals := make(map[int]float64)
	for rows.Next() {
		var userID int
		var total float64
		if err := rows.Scan(&userID, &total); err != nil {
			return nil, err
		}
		totals[userID] = total
	}
	return totals, nil
}

// GetMonthlyTotal retorna o total gasto no mês
func (r *PurchaseRepository) GetMonthlyTotal(month string) (float64, error) {
	var total float64
	err := r.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM purchases WHERE month = ?`, month).Scan(&total)
	return total, err
}

// GetPurchaseCountByUser retorna a contagem de compras por usuário no mês
func (r *PurchaseRepository) GetPurchaseCountByUser(month string) (map[int]int, error) {
	query := `
		SELECT user_id, COUNT(*) as count
		FROM purchases
		WHERE month = ?
		GROUP BY user_id
	`
	rows, err := r.db.Query(query, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[int]int)
	for rows.Next() {
		var userID int
		var count int
		if err := rows.Scan(&userID, &count); err != nil {
			return nil, err
		}
		counts[userID] = count
	}
	return counts, nil
}

// GetDistinctMonths retorna lista de meses com compras
func (r *PurchaseRepository) GetDistinctMonths() ([]string, error) {
	query := `SELECT DISTINCT month FROM purchases ORDER BY month DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var months []string
	for rows.Next() {
		var month string
		if err := rows.Scan(&month); err != nil {
			return nil, err
		}
		months = append(months, month)
	}
	return months, nil
}
