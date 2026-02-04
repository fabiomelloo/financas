package services

import (
	"errors"
	"financas/internal/models"
	"financas/internal/repositories"
)

type ExpenseService struct {
	repository *repositories.ExpenseRepository
}

func NewExpenseService(repository *repositories.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repository: repository}
}

func (s *ExpenseService) Create(expense *models.Expense) error {
	if expense.Amount <= 0 {
		return errors.New("o valor deve ser maior que 0")
	}
	if expense.Description == "" {
		return errors.New("a descrição não pode ser vazia")
	}
	if expense.Type == "" {
		return errors.New("o tipo não pode ser vazio")
	}
	if expense.Category == "" {
		return errors.New("a categoria não pode ser vazia")
	}
	if expense.Date.IsZero() {
		return errors.New("a data não pode ser vazia")
	}
	return s.repository.Create(expense)
}

func (s *ExpenseService) FindAll() ([]models.Expense, error) {
	return s.repository.FindAll()
}

func (s *ExpenseService) FindByID(id int) (*models.Expense, error) {
	return s.repository.FindByID(id)
}

func (s *ExpenseService) Update(expense *models.Expense) error {
	if expense.Amount <= 0 {
		return errors.New("o valor deve ser maior que 0")
	}
	if expense.Description == "" {
		return errors.New("a descrição não pode ser vazia")
	}
	if expense.Type == "" {
		return errors.New("o tipo não pode ser vazio")
	}
	if expense.Category == "" {
		return errors.New("a categoria não pode ser vazia")
	}
	if expense.Date.IsZero() {
		return errors.New("a data não pode ser vazia")
	}
	return s.repository.Update(expense)
}

func (s *ExpenseService) Delete(id int) error {
	return s.repository.Delete(id)
}

// InsightsData agrega todos os dados de relatório
type InsightsData struct {
	TotalIncome       float64                       `json:"total_income"`
	TotalExpense      float64                       `json:"total_expense"`
	Balance           float64                       `json:"balance"`
	CategoryStats     []repositories.CategoryMetric `json:"category_stats"`
	MonthlyStats      []repositories.MonthlyMetric  `json:"monthly_stats"`
	TypeStats         []repositories.TypeMetric     `json:"type_stats"`
	TopExpenses       []models.Expense              `json:"top_expenses"`
	TotalTransactions int                           `json:"total_transactions"`
}

func (s *ExpenseService) GetInsights() (*InsightsData, error) {
	income, expense, balance, err := s.repository.GetSummary()
	if err != nil {
		return nil, err
	}

	cats, err := s.repository.GetCategoryBreakdown()
	if err != nil {
		return nil, err
	}

	monthly, err := s.repository.GetMonthlyBreakdown()
	if err != nil {
		return nil, err
	}

	typeStats, err := s.repository.GetTypeBreakdown()
	if err != nil {
		return nil, err
	}

	topExpenses, err := s.repository.GetTopExpenses(5)
	if err != nil {
		return nil, err
	}

	totalTransactions := 0
	for _, ts := range typeStats {
		totalTransactions += ts.Count
	}

	return &InsightsData{
		TotalIncome:       income,
		TotalExpense:      expense,
		Balance:           balance,
		CategoryStats:     cats,
		MonthlyStats:      monthly,
		TypeStats:         typeStats,
		TopExpenses:       topExpenses,
		TotalTransactions: totalTransactions,
	}, nil
}

// RateioStats armazena dados para a cota de lanche da equipe
type RateioStats struct {
	TotalSpent     float64      `json:"total_spent"`
	SharePerPerson float64      `json:"share_per_person"`
	MemberStats    []MemberStat `json:"member_stats"`
}

type MemberStat struct {
	Name    string  `json:"name"`
	Paid    float64 `json:"paid"`
	Balance float64 `json:"balance"` // Positivo = Crédito, Negativo = Débito
}

func (s *ExpenseService) GetRateioStats() (*RateioStats, error) {
	expenses, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	total := 0.0
	payerTotals := make(map[string]float64)

	for _, e := range expenses {
		// Considerar apenas Despesas (Type = despesa) para o rateio?
		// O usuário disse "cota para o lanche".
		// Vamos filtrar por Type="despesa" por segurança.
		if e.Type == "despesa" {
			total += e.Amount

			name := e.Payer
			if name == "" {
				name = "Outros"
			}
			payerTotals[name] += e.Amount
		}
	}

	// 6 pessoas
	share := total / 6.0

	var stats []MemberStat
	for name, paid := range payerTotals {
		balance := paid - share
		stats = append(stats, MemberStat{
			Name:    name,
			Paid:    paid,
			Balance: balance,
		})
	}

	return &RateioStats{
		TotalSpent:     total,
		SharePerPerson: share,
		MemberStats:    stats,
	}, nil
}
