package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatusPending   = "pending"   
	StatusConfirmed = "confirmed" 
	StatusShipped   = "shipped"   
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

const (
	TypeReservation = "reservation"
	TypeOrder       = "order"     
)

type MarketTransaction struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	BuyerID   string  `json:"buyer_id" gorm:"type:varchar(36)"`
	Buyer     User    `json:"buyer" gorm:"foreignKey:BuyerID"` 
	
	ProductID string  `json:"product_id" gorm:"type:varchar(36)"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`

	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"` 
	Status     string  `json:"status" gorm:"default:'pending'"`
	Type       string  `json:"type"` 
	Note       string  `json:"note"`
}

type CreateOrderRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Note      string `json:"note"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"` 
}