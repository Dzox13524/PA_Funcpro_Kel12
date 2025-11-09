package domain

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        string 
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
	Name string
	Email string
	Password string
	Role string
}