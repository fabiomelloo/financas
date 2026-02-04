package controllers

import (
	"financas/internal/models"
	"financas/internal/services"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type ExpenseController struct {
	service *services.ExpenseService
}

// PageData is the structure passed to templates
type PageData struct {
	CurrentPage string
	Expenses    []models.Expense
	Expense     *models.Expense
}

func NewExpenseController(service *services.ExpenseService) *ExpenseController {
	return &ExpenseController{service: service}
}

func (c *ExpenseController) Index(w http.ResponseWriter, r *http.Request) {
	expenses, err := c.service.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/index.html",
	))

	data := PageData{
		CurrentPage: "index",
		Expenses:    expenses,
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

func (c *ExpenseController) Create(w http.ResponseWriter, r *http.Request) {
	// GET - Show the create form
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles(
			"web/templates/layout.html",
			"web/templates/create.html",
		))

		data := PageData{
			CurrentPage: "create",
		}

		tmpl.ExecuteTemplate(w, "layout", data)
		return
	}

	// POST - Process form data and save
	if r.Method == http.MethodPost {
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		date, _ := time.Parse("2006-01-02", r.FormValue("date"))

		expense := &models.Expense{
			Description: r.FormValue("description"),
			Amount:      amount,
			Type:        r.FormValue("type"),
			Category:    r.FormValue("category"),
			Date:        date,
		}
		err := c.service.Create(expense)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (c *ExpenseController) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	expense, err := c.service.FindByID(id)
	if err != nil {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/edit.html",
	))

	data := PageData{
		CurrentPage: "edit",
		Expense:     expense,
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

func (c *ExpenseController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	date, _ := time.Parse("2006-01-02", r.FormValue("date"))

	expense := &models.Expense{
		ID:          id,
		Description: r.FormValue("description"),
		Amount:      amount,
		Type:        r.FormValue("type"),
		Category:    r.FormValue("category"),
		Date:        date,
	}

	err = c.service.Update(expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *ExpenseController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
