package models

import (
	"fmt"
	"time"
)

// InvoiceType representa o tipo de nota fiscal
type InvoiceType string

const (
	InvoiceTypeService InvoiceType = "service"
	InvoiceTypeProduct InvoiceType = "product"
	InvoiceTypeMixed   InvoiceType = "mixed"
)

// InvoiceStatus representa o status da nota fiscal
type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusIssued    InvoiceStatus = "issued"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

// InvoiceItem representa um item da nota fiscal
type InvoiceItem struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

// Invoice representa uma nota fiscal
type Invoice struct {
	ID           string          `json:"id"`
	Number       string          `json:"number"`
	Type         InvoiceType     `json:"type"`
	Status       InvoiceStatus   `json:"status"`
	PatientID    string          `json:"patient_id"`
	PatientName  string          `json:"patient_name"`
	PatientEmail string          `json:"patient_email"`
	Items        []InvoiceItem   `json:"items"`
	Subtotal     float64         `json:"subtotal"`
	TaxAmount    float64         `json:"tax_amount"`
	TotalAmount  float64         `json:"total_amount"`
	IssueDate    time.Time       `json:"issue_date"`
	DueDate      time.Time       `json:"due_date"`
	Notes        string          `json:"notes,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// IsValid verifica se os campos obrigatórios da nota fiscal estão preenchidos
func (i *Invoice) IsValid() error {
	if i.Number == "" {
		return fmt.Errorf("invoice number is required")
	}
	if i.Type == "" {
		return fmt.Errorf("invoice type is required")
	}
	if i.PatientID == "" {
		return fmt.Errorf("patient ID is required")
	}
	if i.PatientName == "" {
		return fmt.Errorf("patient name is required")
	}
	if len(i.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	if i.TotalAmount <= 0 {
		return fmt.Errorf("total amount must be greater than zero")
	}
	if i.IssueDate.IsZero() {
		return fmt.Errorf("issue date is required")
	}
	if i.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}

	return nil
}

// CalculateTotals calcula os totais da nota fiscal
func (i *Invoice) CalculateTotals() {
	i.Subtotal = 0
	for idx := range i.Items {
		i.Items[idx].TotalPrice = float64(i.Items[idx].Quantity) * i.Items[idx].UnitPrice
		i.Subtotal += i.Items[idx].TotalPrice
	}
	i.TotalAmount = i.Subtotal + i.TaxAmount
}