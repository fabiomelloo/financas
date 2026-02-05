package services

import (
	"errors"
	"financas/internal/models"
	"financas/internal/repositories"
	"strings"
)

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{repository: repository}
}

// Create cria um novo membro da equipe
func (s *UserService) Create(user *models.User) error {
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("o nome não pode ser vazio")
	}
	user.Name = strings.TrimSpace(user.Name)
	user.Points = 0 // Inicia com 0 pontos
	return s.repository.Create(user)
}

// FindAll retorna todos os membros
func (s *UserService) FindAll() ([]models.User, error) {
	return s.repository.FindAll()
}

// FindByID busca um membro pelo ID
func (s *UserService) FindByID(id int) (*models.User, error) {
	return s.repository.FindByID(id)
}

// GetRanking retorna o ranking de pontos
func (s *UserService) GetRanking() ([]models.User, error) {
	return s.repository.GetRanking()
}

// UpdatePoints atualiza os pontos de um usuário
func (s *UserService) UpdatePoints(userID int, points int) error {
	return s.repository.UpdatePoints(userID, points)
}

// Delete remove um usuário
func (s *UserService) Delete(id int) error {
	return s.repository.Delete(id)
}

// Count retorna o número de membros
func (s *UserService) Count() (int, error) {
	return s.repository.Count()
}
