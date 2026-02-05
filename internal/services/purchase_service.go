package services

import (
	"errors"
	"financas/internal/models"
	"financas/internal/repositories"
	"time"
)

type PurchaseService struct {
	purchaseRepo *repositories.PurchaseRepository
	userRepo     *repositories.UserRepository
}

func NewPurchaseService(purchaseRepo *repositories.PurchaseRepository, userRepo *repositories.UserRepository) *PurchaseService {
	return &PurchaseService{
		purchaseRepo: purchaseRepo,
		userRepo:     userRepo,
	}
}

// Create registra uma nova compra de lanche
func (s *PurchaseService) Create(purchase *models.Purchase) error {
	if purchase.Amount <= 0 {
		return errors.New("o valor deve ser maior que 0")
	}
	if purchase.UserID <= 0 {
		return errors.New("usuário inválido")
	}
	if purchase.Date.IsZero() {
		return errors.New("a data não pode ser vazia")
	}

	// Verificar se o usuário existe
	_, err := s.userRepo.FindByID(purchase.UserID)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	return s.purchaseRepo.Create(purchase)
}

// FindAll retorna todas as compras
func (s *PurchaseService) FindAll() ([]models.Purchase, error) {
	return s.purchaseRepo.FindAll()
}

// FindByMonth retorna compras de um mês específico
func (s *PurchaseService) FindByMonth(month string) ([]models.Purchase, error) {
	return s.purchaseRepo.FindByMonth(month)
}

// FindByID busca uma compra pelo ID
func (s *PurchaseService) FindByID(id int) (*models.Purchase, error) {
	return s.purchaseRepo.FindByID(id)
}

// Delete remove uma compra
func (s *PurchaseService) Delete(id int) error {
	return s.purchaseRepo.Delete(id)
}

// GetCurrentMonth retorna o mês atual no formato "2006-01"
func (s *PurchaseService) GetCurrentMonth() string {
	return time.Now().Format("2006-01")
}

// GetDistinctMonths retorna lista de meses com compras
func (s *PurchaseService) GetDistinctMonths() ([]string, error) {
	return s.purchaseRepo.GetDistinctMonths()
}

// RateioData contém os dados de rateio para um mês
type RateioData struct {
	Month          string
	TotalSpent     float64
	SharePerPerson float64
	MemberCount    int
	MemberStats    []MemberRateioStat
}

type MemberRateioStat struct {
	UserID   int
	UserName string
	Paid     float64
	Share    float64
	Balance  float64 // Positivo = crédito, Negativo = débito
}

// CalculateRateio calcula o rateio para um mês
func (s *PurchaseService) CalculateRateio(month string) (*RateioData, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	memberCount := len(users)
	if memberCount == 0 {
		return &RateioData{
			Month:          month,
			TotalSpent:     0,
			SharePerPerson: 0,
			MemberCount:    0,
			MemberStats:    []MemberRateioStat{},
		}, nil
	}

	// Obter totais por usuário
	totals, err := s.purchaseRepo.GetMonthlyTotalByUser(month)
	if err != nil {
		return nil, err
	}

	// Calcular total gasto no mês
	total, err := s.purchaseRepo.GetMonthlyTotal(month)
	if err != nil {
		return nil, err
	}

	// Calcular cota por pessoa
	share := total / float64(memberCount)

	// Montar estatísticas de cada membro
	var stats []MemberRateioStat
	for _, user := range users {
		paid := totals[user.ID] // 0 se não pagou nada
		balance := paid - share

		stats = append(stats, MemberRateioStat{
			UserID:   user.ID,
			UserName: user.Name,
			Paid:     paid,
			Share:    share,
			Balance:  balance,
		})
	}

	return &RateioData{
		Month:          month,
		TotalSpent:     total,
		SharePerPerson: share,
		MemberCount:    memberCount,
		MemberStats:    stats,
	}, nil
}
