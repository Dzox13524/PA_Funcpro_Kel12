package service

import (
	"context"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/google/uuid"
)

// Baris 11: Definisi Service
type CreateArticleServiceFunc func(ctx context.Context, authorID string, req domain.CreateArticleRequest) (domain.Article, error)
type GetAllArticlesServiceFunc func(ctx context.Context) ([]domain.Article, error)
type GetArticleByIDServiceFunc func(ctx context.Context, id string) (domain.Article, error)

// Baris 16: Logika membuat artikel (POST).
func NewCreateArticleService(createRepo repository.CreateArticleRepoFunc) CreateArticleServiceFunc {
	return func(ctx context.Context, authorID string, req domain.CreateArticleRequest) (domain.Article, error) {
		// Kita buat objek Article baru.
		newArticle := domain.Article{
			ID:       uuid.New().String(), // Generate ID unik
			Title:    req.Title,           // Ambil judul dari input
			Content:  req.Content,         // Ambil isi dari input
			AuthorID: authorID,            // Simpan siapa pembuatnya (dari Token)
		}
		return createRepo(ctx, newArticle)
	}
}

// Baris 29: Logika ambil semua (GET List).
func NewGetAllArticlesService(getAllRepo repository.GetAllArticlesRepoFunc) GetAllArticlesServiceFunc {
	return func(ctx context.Context) ([]domain.Article, error) {
		return getAllRepo(ctx)
	}
}

// Baris 36: Logika ambil detail (GET Detail).
func NewGetArticleByIDService(getByIDRepo repository.GetArticleByIDRepoFunc) GetArticleByIDServiceFunc {
	return func(ctx context.Context, id string) (domain.Article, error) {
		return getByIDRepo(ctx, id)
	}
}