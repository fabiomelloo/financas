package controllers

import (
	"financas/internal/models"
	"financas/internal/services"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type UserController struct {
	service *services.UserService
}

// UserPageData é a estrutura passada para os templates de usuários
type UserPageData struct {
	CurrentPage string
	Users       []models.User
	User        *models.User
	CSRFToken   string
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

// Index lista todos os membros da equipe
func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.FindAll()
	if err != nil {
		log.Printf("erro ao buscar usuários: %v", err)
		http.Error(w, "erro ao carregar membros", http.StatusInternalServerError)
		return
	}

	csrfToken, err := generateCSRFToken(w, r)
	if err != nil {
		log.Printf("erro ao gerar csrf token: %v", err)
		http.Error(w, "erro interno", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/users.html",
	))

	data := UserPageData{
		CurrentPage: "users",
		Users:       users,
		CSRFToken:   csrfToken,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("erro no template: %v", err)
		http.Error(w, "erro ao renderizar página", http.StatusInternalServerError)
	}
}

// Create cria um novo membro da equipe
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if !validateCSRFToken(w, r) {
		http.Error(w, "requisição inválida", http.StatusForbidden)
		return
	}

	user := &models.User{
		Name: r.FormValue("name"),
	}

	if err := c.service.Create(user); err != nil {
		log.Printf("erro ao criar usuário: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// Delete remove um membro da equipe
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := c.service.Delete(id); err != nil {
		log.Printf("erro ao remover usuário: %v", err)
		http.Error(w, "erro ao remover membro", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
