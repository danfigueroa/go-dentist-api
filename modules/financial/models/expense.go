package models

import (
	"fmt"
	"time"
)

// ExpenseCategory representa as categorias de gastos
type ExpenseCategory string

const (
	ExpenseCategoryMaterials  ExpenseCategory = "materials"
	ExpenseCategoryRent       ExpenseCategory = "rent"
	ExpenseCategoryUtilities  ExpenseCategory = "utilities"
	ExpenseCategoryStaff      ExpenseCategory = "staff"
	ExpenseCategoryEquipment  ExpenseCategory = "equipment"
	ExpenseCategoryOther      ExpenseCategory = "other"
)

// Expense representa um gasto da clínica
type Expense struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Amount      float64         `json:"amount"`
	Category    ExpenseCategory `json:"category"`
	Date        time.Time       `json:"date"`
	Supplier    string          `json:"supplier,omitempty"`
	InvoiceID   string          `json:"invoice_id,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// IsValid verifica se os campos obrigatórios do gasto estão preenchidos
func (e *Expense) IsValid() error {
	if e.Description == "" {
		return fmt.Errorf("description is required")
	}
	if e.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if e.Category == "" {
		return fmt.Errorf("category is required")
	}
	if e.Date.IsZero() {
		return fmt.Errorf("date is required")
	}

	return nil
}