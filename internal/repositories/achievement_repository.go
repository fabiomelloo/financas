package repositories

import (
	"database/sql"
	"financas/internal/models"
)

type AchievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

// GetAll retorna todas as conquistas disponíveis
func (r *AchievementRepository) GetAll() ([]models.Achievement, error) {
	query := `SELECT id, name, description, icon FROM achievements ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []models.Achievement
	for rows.Next() {
		var a models.Achievement
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Icon); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

// GetByID retorna uma conquista pelo ID
func (r *AchievementRepository) GetByID(id int) (*models.Achievement, error) {
	query := `SELECT id, name, description, icon FROM achievements WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var a models.Achievement
	if err := row.Scan(&a.ID, &a.Name, &a.Description, &a.Icon); err != nil {
		return nil, err
	}
	return &a, nil
}

// GetByName retorna uma conquista pelo nome
func (r *AchievementRepository) GetByName(name string) (*models.Achievement, error) {
	query := `SELECT id, name, description, icon FROM achievements WHERE name = ?`
	row := r.db.QueryRow(query, name)

	var a models.Achievement
	if err := row.Scan(&a.ID, &a.Name, &a.Description, &a.Icon); err != nil {
		return nil, err
	}
	return &a, nil
}

// AwardToUser atribui uma conquista a um usuário para um mês específico
func (r *AchievementRepository) AwardToUser(userID, achievementID int, month string) error {
	query := `INSERT OR IGNORE INTO user_achievements (user_id, achievement_id, month) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, userID, achievementID, month)
	return err
}

// GetUserAchievements retorna as conquistas de um usuário
func (r *AchievementRepository) GetUserAchievements(userID int) ([]models.UserAchievement, error) {
	query := `
		SELECT ua.user_id, u.name, ua.achievement_id, a.name, a.icon, ua.month, ua.awarded_at
		FROM user_achievements ua
		JOIN users u ON ua.user_id = u.id
		JOIN achievements a ON ua.achievement_id = a.id
		WHERE ua.user_id = ?
		ORDER BY ua.awarded_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAchievements []models.UserAchievement
	for rows.Next() {
		var ua models.UserAchievement
		var awardedAtStr string
		if err := rows.Scan(&ua.UserID, &ua.UserName, &ua.AchievementID, &ua.AchievementName, &ua.AchievementIcon, &ua.Month, &awardedAtStr); err != nil {
			return nil, err
		}
		userAchievements = append(userAchievements, ua)
	}
	return userAchievements, nil
}

// GetMonthlyAchievements retorna todas as conquistas de um mês
func (r *AchievementRepository) GetMonthlyAchievements(month string) ([]models.UserAchievement, error) {
	query := `
		SELECT ua.user_id, u.name, ua.achievement_id, a.name, a.icon, ua.month, ua.awarded_at
		FROM user_achievements ua
		JOIN users u ON ua.user_id = u.id
		JOIN achievements a ON ua.achievement_id = a.id
		WHERE ua.month = ?
		ORDER BY ua.achievement_id
	`
	rows, err := r.db.Query(query, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAchievements []models.UserAchievement
	for rows.Next() {
		var ua models.UserAchievement
		var awardedAtStr string
		if err := rows.Scan(&ua.UserID, &ua.UserName, &ua.AchievementID, &ua.AchievementName, &ua.AchievementIcon, &ua.Month, &awardedAtStr); err != nil {
			return nil, err
		}
		userAchievements = append(userAchievements, ua)
	}
	return userAchievements, nil
}

// GetRecentAchievements retorna as últimas conquistas atribuídas
func (r *AchievementRepository) GetRecentAchievements(limit int) ([]models.UserAchievement, error) {
	query := `
		SELECT ua.user_id, u.name, ua.achievement_id, a.name, a.icon, ua.month, ua.awarded_at
		FROM user_achievements ua
		JOIN users u ON ua.user_id = u.id
		JOIN achievements a ON ua.achievement_id = a.id
		ORDER BY ua.awarded_at DESC
		LIMIT ?
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAchievements []models.UserAchievement
	for rows.Next() {
		var ua models.UserAchievement
		var awardedAtStr string
		if err := rows.Scan(&ua.UserID, &ua.UserName, &ua.AchievementID, &ua.AchievementName, &ua.AchievementIcon, &ua.Month, &awardedAtStr); err != nil {
			return nil, err
		}
		userAchievements = append(userAchievements, ua)
	}
	return userAchievements, nil
}
