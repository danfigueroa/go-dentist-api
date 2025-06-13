package models

import (
	"fmt"
	"time"
)

// PaymentMethod representa os métodos de pagamento
type PaymentMethod string

const (
	PaymentMethodCash       PaymentMethod = "cash"
	PaymentMethodCard       PaymentMethod = "card"
	PaymentMethodPix        PaymentMethod = "pix"
	PaymentMethodBankSlip   PaymentMethod = "bank_slip"
	PaymentMethodInsurance  PaymentMethod = "insurance"
)

// PaymentStatus representa o status do pagamento
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// Revenue representa uma receita da clínica
type Revenue struct {
	ID            string        `json:"id"`
	Description   string        `json:"description"`
	Amount        float64       `json:"amount"`
	PatientID     string        `json:"patient_id"`
	ProcedureID   string        `json:"procedure_id,omitempty"`
	AppointmentID string        `json:"appointment_id,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	DueDate       time.Time     `json:"due_date"`
	PaidDate      *time.Time    `json:"paid_date,omitempty"`
	InvoiceID     string        `json:"invoice_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// IsValid verifica se os campos obrigatórios da receita estão preenchidos
func (r *Revenue) IsValid() error {
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if r.PatientID == "" {
		return fmt.Errorf("patient ID is required")
	}
	if r.PaymentMethod == "" {
		return fmt.Errorf("payment method is required")
	}
	if r.PaymentStatus == "" {
		return fmt.Errorf("payment status is required")
	}
	if r.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}

	return nil
}