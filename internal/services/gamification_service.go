package services

import (
	"financas/internal/models"
	"financas/internal/repositories"
	"math"
	"sort"
)

// GamificationService gerencia o sistema de pontos e conquistas
type GamificationService struct {
	userRepo        *repositories.UserRepository
	purchaseRepo    *repositories.PurchaseRepository
	achievementRepo *repositories.AchievementRepository
}

func NewGamificationService(
	userRepo *repositories.UserRepository,
	purchaseRepo *repositories.PurchaseRepository,
	achievementRepo *repositories.AchievementRepository,
) *GamificationService {
	return &GamificationService{
		userRepo:        userRepo,
		purchaseRepo:    purchaseRepo,
		achievementRepo: achievementRepo,
	}
}

// Constantes de pontos
const (
	PointsPaidSnack       = 10  // Pagou o lanche do dia
	PointsAboveAverage    = 5   // Pagou acima da média do mês
	PointsTopCreditor     = 20  // Maior crédito do mês (ajudou o time)
	PointsTopDebtor       = -10 // Maior débito do mês
	PointsNoParticipation = -15 // Não participou de nenhuma compra no mês
)

// AwardPointsForPurchase atribui pontos quando alguém paga um lanche
func (s *GamificationService) AwardPointsForPurchase(userID int) error {
	return s.userRepo.UpdatePoints(userID, PointsPaidSnack)
}

// ProcessMonthlyGamification processa pontos e conquistas do mês
// Deve ser chamado no fechamento do mês
func (s *GamificationService) ProcessMonthlyGamification(month string) error {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	// Obter totais pagos por usuário
	totals, err := s.purchaseRepo.GetMonthlyTotalByUser(month)
	if err != nil {
		return err
	}

	// Obter contagem de compras por usuário
	counts, err := s.purchaseRepo.GetPurchaseCountByUser(month)
	if err != nil {
		return err
	}

	// Calcular total e média
	totalSpent, _ := s.purchaseRepo.GetMonthlyTotal(month)
	share := totalSpent / float64(len(users))

	// Calcular balanços
	type userBalance struct {
		UserID  int
		Paid    float64
		Count   int
		Balance float64
	}

	var balances []userBalance
	for _, user := range users {
		paid := totals[user.ID]
		count := counts[user.ID]
		balance := paid - share
		balances = append(balances, userBalance{
			UserID:  user.ID,
			Paid:    paid,
			Count:   count,
			Balance: balance,
		})
	}

	// Ordenar por balanço (maior crédito primeiro)
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].Balance > balances[j].Balance
	})

	// Processar pontos e conquistas para cada usuário
	for i, b := range balances {
		// Quem pagou acima da média ganha pontos extras
		if b.Paid > share && b.Paid > 0 {
			s.userRepo.UpdatePoints(b.UserID, PointsAboveAverage)
		}

		// Quem não participou perde pontos
		if b.Count == 0 {
			s.userRepo.UpdatePoints(b.UserID, PointsNoParticipation)
		}

		// Maior crédito do mês (primeiro da lista ordenada)
		if i == 0 && b.Balance > 0 {
			s.userRepo.UpdatePoints(b.UserID, PointsTopCreditor)
			s.awardAchievement(b.UserID, "Mecenas", month)
		}

		// Maior débito do mês (último da lista com balanço negativo)
		if i == len(balances)-1 && b.Balance < 0 {
			s.userRepo.UpdatePoints(b.UserID, PointsTopDebtor)
			s.awardAchievement(b.UserID, "Caloteiro Simpático", month)
		}

		// Saldo equilibrado (próximo de zero, margem de 5%)
		if share > 0 && math.Abs(b.Balance) <= share*0.05 {
			s.awardAchievement(b.UserID, "Equilibrado", month)
		}
	}

	// Quem mais comprou no mês (mais vezes pagou)
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].Count > balances[j].Count
	})
	if len(balances) > 0 && balances[0].Count > 0 {
		s.awardAchievement(balances[0].UserID, "Contador", month)
	}

	// Maior gasto individual (ordenar por valor pago)
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].Paid > balances[j].Paid
	})
	if len(balances) > 0 && balances[0].Paid > 0 {
		s.awardAchievement(balances[0].UserID, "Mão Aberta", month)
	}

	return nil
}

// awardAchievement atribui uma conquista a um usuário
func (s *GamificationService) awardAchievement(userID int, achievementName, month string) error {
	achievement, err := s.achievementRepo.GetByName(achievementName)
	if err != nil {
		return err
	}
	return s.achievementRepo.AwardToUser(userID, achievement.ID, month)
}

// GetRanking retorna o ranking geral de pontos
func (s *GamificationService) GetRanking() ([]models.User, error) {
	return s.userRepo.GetRanking()
}

// GetAllAchievements retorna todas as conquistas disponíveis
func (s *GamificationService) GetAllAchievements() ([]models.Achievement, error) {
	return s.achievementRepo.GetAll()
}

// GetUserAchievements retorna as conquistas de um usuário
func (s *GamificationService) GetUserAchievements(userID int) ([]models.UserAchievement, error) {
	return s.achievementRepo.GetUserAchievements(userID)
}

// GetMonthlyAchievements retorna as conquistas de um mês
func (s *GamificationService) GetMonthlyAchievements(month string) ([]models.UserAchievement, error) {
	return s.achievementRepo.GetMonthlyAchievements(month)
}

// GetRecentAchievements retorna as últimas conquistas
func (s *GamificationService) GetRecentAchievements(limit int) ([]models.UserAchievement, error) {
	return s.achievementRepo.GetRecentAchievements(limit)
}

// DashboardData contém dados para o dashboard de gamificação
type DashboardData struct {
	Ranking            []models.User
	RecentAchievements []models.UserAchievement
	AllAchievements    []models.Achievement
	CurrentMonth       string
}

// GetDashboardData retorna dados consolidados para o dashboard
func (s *GamificationService) GetDashboardData(currentMonth string) (*DashboardData, error) {
	ranking, err := s.GetRanking()
	if err != nil {
		return nil, err
	}

	recent, err := s.GetRecentAchievements(10)
	if err != nil {
		return nil, err
	}

	all, err := s.GetAllAchievements()
	if err != nil {
		return nil, err
	}

	return &DashboardData{
		Ranking:            ranking,
		RecentAchievements: recent,
		AllAchievements:    all,
		CurrentMonth:       currentMonth,
	}, nil
}
