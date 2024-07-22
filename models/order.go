package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderID     int64      `json:"order_id"`
	CustomerID  uuid.UUID  `json:"customer_id"`
	LineItem    []LineItem `json:"line_item"`
	OrderStatus string     `json:"order_status"`
	CreatedAt   *time.Time `json:"created_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type LineItem struct {
	ItemId   uuid.UUID `json:"item_id"`
	Quantity uint      `json:"quantity"`
	Price    uint      `json:"price"`
}
