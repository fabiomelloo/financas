package models

import "time"

// Achievement representa uma conquista/badge do sistema de gamificação
type Achievement struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"` // Emoji ou ícone
}

// UserAchievement registra quando um usuário ganhou uma conquista
type UserAchievement struct {
	UserID          int       `json:"user_id"`
	UserName        string    `json:"user_name"` // Para exibição
	AchievementID   int       `json:"achievement_id"`
	AchievementName string    `json:"achievement_name"` // Para exibição
	AchievementIcon string    `json:"achievement_icon"` // Para exibição
	Month           string    `json:"month"`            // Mês de referência
	AwardedAt       time.Time `json:"awarded_at"`
}
