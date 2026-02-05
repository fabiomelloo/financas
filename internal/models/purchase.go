package models

import "time"

// Purchase representa uma compra de lanche feita por um membro da equipe
type Purchase struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"user_name"` // Para exibição (preenchido via JOIN)
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	Month     string    `json:"month"` // Formato "2026-02" para agrupamento mensal
	CreatedAt time.Time `json:"created_at"`
}
