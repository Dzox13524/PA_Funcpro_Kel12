package repository

import (
	"context"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"gorm.io/gorm"
)

type CreateArticleRepoFunc func(ctx context.Context, article domain.Article) (domain.Article, error)
type GetAllArticlesRepoFunc func(ctx context.Context) ([]domain.Article, error)
type GetArticleByIDRepoFunc func(ctx context.Context, id string) (domain.Article, error)

func NewCreateArticleRepository(db *gorm.DB) CreateArticleRepoFunc {
	return func(ctx context.Context, article domain.Article) (domain.Article, error) {
		result := db.WithContext(ctx).Create(&article)
		return article, result.Error
	}
}

func NewGetAllArticlesRepository(db *gorm.DB) GetAllArticlesRepoFunc {
	return func(ctx context.Context) ([]domain.Article, error) {
		var articles []domain.Article
		result := db.WithContext(ctx).Find(&articles)
		return articles, result.Error
	}
}

func NewGetArticleByIDRepository(db *gorm.DB) GetArticleByIDRepoFunc {
	return func(ctx context.Context, id string) (domain.Article, error) {
		var article domain.Article
		result := db.WithContext(ctx).First(&article, "id = ?", id)
		return article, result.Error
	}
}