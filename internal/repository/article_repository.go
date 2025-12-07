package repository

import (
	"context"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"gorm.io/gorm"
)

// Baris 10: Definisi fungsi-fungsi apa saja yang bisa dilakukan ke DB.
type CreateArticleRepoFunc func(ctx context.Context, article domain.Article) (domain.Article, error)
type GetAllArticlesRepoFunc func(ctx context.Context) ([]domain.Article, error)
type GetArticleByIDRepoFunc func(ctx context.Context, id string) (domain.Article, error)

// Baris 15: Fungsi untuk simpan artikel baru ke DB.
func NewCreateArticleRepository(db *gorm.DB) CreateArticleRepoFunc {
	return func(ctx context.Context, article domain.Article) (domain.Article, error) {
		result := db.WithContext(ctx).Create(&article)
		return article, result.Error
	}
}

// Baris 23: Fungsi ambil semua artikel.
func NewGetAllArticlesRepository(db *gorm.DB) GetAllArticlesRepoFunc {
	return func(ctx context.Context) ([]domain.Article, error) {
		var articles []domain.Article
		// Find(&articles) akan mengambil semua data dari tabel articles.
		result := db.WithContext(ctx).Find(&articles)
		return articles, result.Error
	}
}

// Baris 32: Fungsi ambil satu artikel berdasarkan ID.
func NewGetArticleByIDRepository(db *gorm.DB) GetArticleByIDRepoFunc {
	return func(ctx context.Context, id string) (domain.Article, error) {
		var article domain.Article
		// First() mencari data pertama yang cocok dengan kondisi "id = ?".
		result := db.WithContext(ctx).First(&article, "id = ?", id)
		return article, result.Error
	}
}