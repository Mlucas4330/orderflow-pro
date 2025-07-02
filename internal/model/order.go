package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPaid      Status = "paid"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
	StatusRefunded  Status = "refunded"
)

type PaymentMethod string

const (
	PaymentCreditCard   PaymentMethod = "credit_card"
	PaymentPix          PaymentMethod = "pix"
	PaymentBoleto       PaymentMethod = "boleto"
	PaymentPaypal       PaymentMethod = "paypal"
	PaymentCash         PaymentMethod = "cash"
	PaymentBankTransfer PaymentMethod = "bank_transfer"
)

type Order struct {
	ID              uuid.UUID       `db:"id"`
	CustomerID      *uuid.UUID      `db:"customer_id"`
	Status          Status          `db:"status"`
	TotalAmount     decimal.Decimal `db:"total_amount"`
	ShippingCost    decimal.Decimal `db:"shipping_cost"`
	DiscountAmount  decimal.Decimal `db:"discount_amount"`
	Currency        string          `db:"currency"`
	PaymentMethod   PaymentMethod   `db:"payment_method"`
	ShippingAddress *uuid.UUID      `db:"shipping_address_id"`
	BillingAddress  *uuid.UUID      `db:"billing_address_id"`
	ItemsCount      int             `db:"items_count"`
	Notes           *string         `db:"notes"`
	IsTest          bool            `db:"is_test"`
	Metadata        json.RawMessage `db:"metadata"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}
