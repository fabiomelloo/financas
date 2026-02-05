package models

// MonthlyBalance armazena o balanço mensal de cada membro
// Balance positivo = Crédito (pagou mais que a cota)
// Balance negativo = Débito (pagou menos que a cota)
type MonthlyBalance struct {
	UserID     int     `json:"user_id"`
	UserName   string  `json:"user_name"` // Para exibição
	Month      string  `json:"month"`     // Formato "2026-02"
	TotalPaid  float64 `json:"total_paid"`
	ShareValue float64 `json:"share_value"` // Cota = TotalMês / 6
	Balance    float64 `json:"balance"`     // TotalPaid - ShareValue
}
