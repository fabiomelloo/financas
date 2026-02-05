package models

import "time"

// User representa um membro da equipe no sistema de rateio
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
