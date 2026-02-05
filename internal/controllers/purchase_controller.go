package controllers

import (
	"financas/internal/models"
	"financas/internal/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PurchaseController struct {
	purchaseService     *services.PurchaseService
	userService         *services.UserService
	gamificationService *services.GamificationService
}

// PurchasePageData é a estrutura passada para os templates de compras
type PurchasePageData struct {
	CurrentPage  string
	Purchases    []models.Purchase
	Users        []models.User
	RateioData   *services.RateioData
	Months       []string
	CurrentMonth string
	CSRFToken    string
}

func NewPurchaseController(
	purchaseService *services.PurchaseService,
	userService *services.UserService,
	gamificationService *services.GamificationService,
) *PurchaseController {
	return &PurchaseController{
		purchaseService:     purchaseService,
		userService:         userService,
		gamificationService: gamificationService,
	}
}

// Index lista todas as compras e mostra o rateio
func (c *PurchaseController) Index(w http.ResponseWriter, r *http.Request) {
	currentMonth := c.purchaseService.GetCurrentMonth()

	// Verificar se foi passado um mês específico
	month := r.URL.Query().Get("month")
	if month == "" {
		month = currentMonth
	}

	purchases, err := c.purchaseService.FindByMonth(month)
	if err != nil {
		log.Printf("erro ao buscar compras: %v", err)
		http.Error(w, "erro ao carregar compras", http.StatusInternalServerError)
		return
	}

	users, err := c.userService.FindAll()
	if err != nil {
		log.Printf("erro ao buscar usuários: %v", err)
		http.Error(w, "erro ao carregar usuários", http.StatusInternalServerError)
		return
	}

	rateio, err := c.purchaseService.CalculateRateio(month)
	if err != nil {
		log.Printf("erro ao calcular rateio: %v", err)
		http.Error(w, "erro ao calcular rateio", http.StatusInternalServerError)
		return
	}

	months, _ := c.purchaseService.GetDistinctMonths()

	csrfToken, err := generateCSRFToken(w, r)
	if err != nil {
		log.Printf("erro ao gerar csrf token: %v", err)
		http.Error(w, "erro interno", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/purchases.html",
	))

	data := PurchasePageData{
		CurrentPage:  "purchases",
		Purchases:    purchases,
		Users:        users,
		RateioData:   rateio,
		Months:       months,
		CurrentMonth: month,
		CSRFToken:    csrfToken,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("erro no template: %v", err)
		http.Error(w, "erro ao renderizar página", http.StatusInternalServerError)
	}
}

// Create registra uma nova compra de lanche
func (c *PurchaseController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if !validateCSRFToken(w, r) {
		http.Error(w, "requisição inválida", http.StatusForbidden)
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, "usuário inválido", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "valor inválido", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "data inválida", http.StatusBadRequest)
		return
	}

	purchase := &models.Purchase{
		UserID: userID,
		Amount: amount,
		Date:   date,
	}

	if err := c.purchaseService.Create(purchase); err != nil {
		log.Printf("erro ao criar compra: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Atribuir pontos pela compra (+10)
	c.gamificationService.AwardPointsForPurchase(userID)

	http.Redirect(w, r, "/purchases", http.StatusSeeOther)
}

// Delete remove uma compra
func (c *PurchaseController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if !validateCSRFToken(w, r) {
		http.Error(w, "requisição inválida", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := c.purchaseService.Delete(id); err != nil {
		log.Printf("erro ao remover compra: %v", err)
		http.Error(w, "erro ao remover compra", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/purchases", http.StatusSeeOther)
}

// ProcessMonth processa pontos e conquistas do mês (fechamento)
func (c *PurchaseController) ProcessMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if !validateCSRFToken(w, r) {
		http.Error(w, "requisição inválida", http.StatusForbidden)
		return
	}

	month := r.FormValue("month")
	if month == "" {
		month = c.purchaseService.GetCurrentMonth()
	}

	if err := c.gamificationService.ProcessMonthlyGamification(month); err != nil {
		log.Printf("erro ao processar gamificação: %v", err)
		http.Error(w, "erro ao processar mês", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/ranking", http.StatusSeeOther)
}
