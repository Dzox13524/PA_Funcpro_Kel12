package domain

import "time"

type PestReport struct{
	ID        				string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CreatedAt 				time.Time `json:"created_at"`
	UpdatedAt 				time.Time `json:"updated_at"`
	UserID    				string    `json:"user_id"`
	City      				string    `json:"city"`
	Description 			string  	`json:"description"`
	PestName  				string    `json:"pest_name"`
	Severity  				string    `json:"severity"`
	VerificationCount int				`json:"verification_count" gorm:"default:0"`
}

type PestVerification struct {
	PestReportID string    `gorm:"primaryKey;type:varchar(36)"`
	UserID       string    `gorm:"primaryKey;type:varchar(36)"` 
	CreatedAt    time.Time `json:"created_at"`
}