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
