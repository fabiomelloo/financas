package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"financas/internal/models"
	"financas/internal/services"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ExpenseController struct {
	service *services.ExpenseService
}

// PageData é a estrutura passada para os templates
type PageData struct {
	CurrentPage string
	Expenses    []models.Expense
	Expense     *models.Expense
	CSRFToken   string
}

func NewExpenseController(service *services.ExpenseService) *ExpenseController {
	return &ExpenseController{service: service}
}

func (c *ExpenseController) Index(w http.ResponseWriter, r *http.Request) {
	expenses, err := c.service.FindAll()
	if err != nil {
		log.Printf("error fetching expenses: %v", err)
		http.Error(w, "erro ao carregar lançamentos", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Controller Index: passing %d expenses to template\n", len(expenses))

	csrfToken, err := generateCSRFToken(w, r)
	if err != nil {
		log.Printf("error generating csrf token: %v", err)
		http.Error(w, "erro interno", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/index.html",
	))

	data := PageData{
		CurrentPage: "index",
		Expenses:    expenses,
		CSRFToken:   csrfToken,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (c *ExpenseController) Create(w http.ResponseWriter, r *http.Request) {
	// GET - Mostra o formulário de criação
	if r.Method == http.MethodGet {
		csrfToken, err := generateCSRFToken(w, r)
		if err != nil {
			log.Printf("error generating csrf token: %v", err)
			http.Error(w, "erro interno", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles(
			"web/templates/layout.html",
			"web/templates/create.html",
		))

		data := PageData{
			CurrentPage: "create",
			CSRFToken:   csrfToken,
		}

		tmpl.ExecuteTemplate(w, "layout", data)
		return
	}

	// POST - Processa os dados do formulário e salva
	if r.Method == http.MethodPost {
		if !validateCSRFToken(w, r) {
			http.Error(w, "requisição inválida", http.StatusForbidden)
			return
		}

		// Parse dos dados do formulário
		if err := r.ParseForm(); err != nil {
			log.Printf("error parsing form: %v", err)
			http.Error(w, "dados inválidos", http.StatusBadRequest)
			return
		}

		// Logs de depuração
		fmt.Printf("Form data received:\n")
		fmt.Printf("  description: %s\n", r.FormValue("description"))
		fmt.Printf("  amount: %s\n", r.FormValue("amount"))
		fmt.Printf("  type: %s\n", r.FormValue("type"))
		fmt.Printf("  category: %s\n", r.FormValue("category"))
		fmt.Printf("  category: %s\n", r.FormValue("category"))
		fmt.Printf("  payer: %s\n", r.FormValue("payer"))
		fmt.Printf("  date: %s\n", r.FormValue("date"))

		amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil {
			log.Printf("error parsing amount: %v", err)
			http.Error(w, "valor inválido", http.StatusBadRequest)
			return
		}
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			log.Printf("error parsing date: %v", err)
			http.Error(w, "data inválida", http.StatusBadRequest)
			return
		}

		expense := &models.Expense{
			Description: r.FormValue("description"),
			Amount:      amount,
			Type:        r.FormValue("type"),
			Category:    r.FormValue("category"),
			Payer:       r.FormValue("payer"),
			Date:        date,
		}

		fmt.Printf("Expense object: %+v\n", expense)

		err = c.service.Create(expense)
		if err != nil {
			log.Printf("error creating expense: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Expense created successfully!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (c *ExpenseController) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	expense, err := c.service.FindByID(id)
	if err != nil {
		http.Error(w, "lançamento não encontrado", http.StatusNotFound)
		return
	}

	csrfToken, err := generateCSRFToken(w, r)
	if err != nil {
		log.Printf("error generating csrf token: %v", err)
		http.Error(w, "erro interno", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/edit.html",
	))

	data := PageData{
		CurrentPage: "edit",
		Expense:     expense,
		CSRFToken:   csrfToken,
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

func (c *ExpenseController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
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

	expense := &models.Expense{
		ID:          id,
		Description: r.FormValue("description"),
		Amount:      amount,
		Type:        r.FormValue("type"),
		Category:    r.FormValue("category"),
		Payer:       r.FormValue("payer"),
		Date:        date,
	}

	err = c.service.Update(expense)
	if err != nil {
		log.Printf("error updating expense: %v", err)
		http.Error(w, "erro ao atualizar lançamento", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *ExpenseController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
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

	err = c.service.Delete(id)
	if err != nil {
		log.Printf("error deleting expense: %v", err)
		http.Error(w, "erro ao remover lançamento", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type InsightsPageData struct {
	CurrentPage string
	Data        *services.InsightsData
}

func (c *ExpenseController) Insights(w http.ResponseWriter, r *http.Request) {
	insights, err := c.service.GetInsights()
	if err != nil {
		log.Printf("error fetching insights: %v", err)
		http.Error(w, "erro ao carregar insights", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/insights.html",
	))

	data := InsightsPageData{
		CurrentPage: "insights",
		Data:        insights,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "erro ao renderizar página", http.StatusInternalServerError)
	}
}

type RateioPageData struct {
	CurrentPage string
	Data        *services.RateioStats
}

func (c *ExpenseController) Rateio(w http.ResponseWriter, r *http.Request) {
	stats, err := c.service.GetRateioStats()
	if err != nil {
		log.Printf("error fetching rateio stats: %v", err)
		http.Error(w, "erro ao carregar dados de rateio", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/rateio.html",
	))

	data := RateioPageData{
		CurrentPage: "rateio",
		Data:        stats,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "erro ao renderizar página", http.StatusInternalServerError)
	}
}

// Helpers CSRF (padrão double submit cookie)
func generateCSRFToken(w http.ResponseWriter, r *http.Request) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // ajustar para true se usar HTTPS
	})

	return token, nil
}

func validateCSRFToken(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return false
	}

	if err := r.ParseForm(); err != nil {
		return false
	}

	formToken := r.FormValue("csrf_token")
	if formToken == "" {
		return false
	}

	return cookie.Value == formToken
}
