package domain

import (
	"time"

	"gorm.io/gorm"
)

// Baris 10: Struktur tabel 'articles' di database.
type Article struct {
	ID        string         `json:"id" gorm:"primaryKey"` // ID unik (UUID)
	CreatedAt time.Time      `json:"created_at"`           // Tanggal dibuat
	UpdatedAt time.Time      `json:"updated_at"`           // Tanggal update
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`       // Fitur soft delete

	Title    string `json:"title"`     // Judul Artikel
	Content  string `json:"content"`   // Isi artikel
	AuthorID string `json:"author_id"` // ID pembuat (Admin)
}

// Baris 21: Format JSON untuk membuat artikel baru.
type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}