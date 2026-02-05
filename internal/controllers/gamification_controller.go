package controllers

import (
	"financas/internal/models"
	"financas/internal/services"
	"html/template"
	"log"
	"net/http"
)

type GamificationController struct {
	gamificationService *services.GamificationService
	purchaseService     *services.PurchaseService
}

// RankingPageData é a estrutura para o template de ranking
type RankingPageData struct {
	CurrentPage string
	Ranking     []models.User
}

// AchievementsPageData é a estrutura para o template de conquistas
type AchievementsPageData struct {
	CurrentPage        string
	Achievements       []models.Achievement
	RecentAchievements []models.UserAchievement
}

func NewGamificationController(
	gamificationService *services.GamificationService,
	purchaseService *services.PurchaseService,
) *GamificationController {
	return &GamificationController{
		gamificationService: gamificationService,
		purchaseService:     purchaseService,
	}
}

// Ranking exibe o ranking de pontos
func (c *GamificationController) Ranking(w http.ResponseWriter, r *http.Request) {
	ranking, err := c.gamificationService.GetRanking()
	if err != nil {
		log.Printf("erro ao buscar ranking: %v", err)
		http.Error(w, "erro ao carregar ranking", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/ranking.html",
	))

	data := RankingPageData{
		CurrentPage: "ranking",
		Ranking:     ranking,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("erro no template: %v", err)
		http.Error(w, "erro ao renderizar página", http.StatusInternalServerError)
	}
}

// Achievements exibe as conquistas disponíveis e recentes
func (c *GamificationController) Achievements(w http.ResponseWriter, r *http.Request) {
	achievements, err := c.gamificationService.GetAllAchievements()
	if err != nil {
		log.Printf("erro ao buscar conquistas: %v", err)
		http.Error(w, "erro ao carregar conquistas", http.StatusInternalServerError)
		return
	}

	recent, err := c.gamificationService.GetRecentAchievements(15)
	if err != nil {
		log.Printf("erro ao buscar conquistas recentes: %v", err)
		http.Error(w, "erro ao carregar conquistas recentes", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/achievements.html",
	))

	data := AchievementsPageData{
		CurrentPage:        "achievements",
		Achievements:       achievements,
		RecentAchievements: recent,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("erro no template: %v", err)
		http.Error(w, "erro ao renderizar página", http.StatusInternalServerError)
	}
}
