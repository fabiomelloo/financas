package repositories

import (
	"database/sql"
	"financas/internal/models"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create insere um novo usuário/membro da equipe
func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (name, points) VALUES (?, ?)`
	result, err := r.db.Exec(query, user.Name, user.Points)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

// FindAll retorna todos os usuários
func (r *UserRepository) FindAll() ([]models.User, error) {
	query := `SELECT id, name, points, created_at, updated_at FROM users ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createdAt, updatedAt string
		if err := rows.Scan(&user.ID, &user.Name, &user.Points, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		users = append(users, user)
	}
	return users, nil
}

// FindByID busca um usuário pelo ID
func (r *UserRepository) FindByID(id int) (*models.User, error) {
	query := `SELECT id, name, points, created_at, updated_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var user models.User
	var createdAt, updatedAt string
	if err := row.Scan(&user.ID, &user.Name, &user.Points, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
	return &user, nil
}

// UpdatePoints atualiza os pontos de um usuário
func (r *UserRepository) UpdatePoints(userID int, points int) error {
	query := `UPDATE users SET points = points + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, points, userID)
	return err
}

// GetRanking retorna os usuários ordenados por pontos (ranking)
func (r *UserRepository) GetRanking() ([]models.User, error) {
	query := `SELECT id, name, points, created_at, updated_at FROM users ORDER BY points DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createdAt, updatedAt string
		if err := rows.Scan(&user.ID, &user.Name, &user.Points, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		users = append(users, user)
	}
	return users, nil
}

// Delete remove um usuário (hard delete)
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// Count retorna o número total de usuários
func (r *UserRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}
