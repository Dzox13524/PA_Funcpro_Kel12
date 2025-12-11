package domain

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        string         `json:"id" gorm:"primaryKey"` 
	CreatedAt time.Time      `json:"created_at"`           
	UpdatedAt time.Time      `json:"updated_at"`         
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`      

	Title    string `json:"title"`  
	Content  string `json:"content"`   
	AuthorID string `json:"author_id"` 
}

type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}