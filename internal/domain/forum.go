package domain

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	UserID    string    `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	Category  string    `json:"category"`
	
	Answers   []Answer  `json:"answers,omitempty" gorm:"foreignKey:QuestionID"`
	
	LikesCount    int64 `json:"likes_count" gorm:"-"`
	IsLiked       bool  `json:"is_liked" gorm:"-"`
	IsFavorited   bool  `json:"is_favorited" gorm:"-"`
}

type Answer struct {
	ID         string         `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	
	QuestionID string    `json:"question_id"`
	UserID     string    `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	
	Content    string    `json:"content"`
	IsSolution bool      `json:"is_solution"`

	LikesCount int64 `json:"likes_count" gorm:"-"`
	IsLiked    bool  `json:"is_liked" gorm:"-"`
}

type QuestionLike struct {
	UserID     string `gorm:"primaryKey"`
	QuestionID string `gorm:"primaryKey"`
	CreatedAt  time.Time
}

type AnswerLike struct {
	UserID    string `gorm:"primaryKey"`
	AnswerID  string `gorm:"primaryKey"`
	CreatedAt time.Time
}

type Favorite struct {
	UserID     string `gorm:"primaryKey"`
	QuestionID string `gorm:"primaryKey"`
	CreatedAt  time.Time
}

type CreateQuestionReq struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type CreateAnswerReq struct {
	Content string `json:"content"`
}

