package service

import (
	"context"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/google/uuid"
)

type CreateArticleServiceFunc func(ctx context.Context, authorID string, req domain.CreateArticleRequest) (domain.Article, error)
type GetAllArticlesServiceFunc func(ctx context.Context) ([]domain.Article, error)
type GetArticleByIDServiceFunc func(ctx context.Context, id string) (domain.Article, error)

func NewCreateArticleService(createRepo repository.CreateArticleRepoFunc) CreateArticleServiceFunc {
	return func(ctx context.Context, authorID string, req domain.CreateArticleRequest) (domain.Article, error) {
		newArticle := domain.Article{
			ID:       uuid.New().String(), 
			Title:    req.Title,           
			Content:  req.Content,         
			AuthorID: authorID,            
		}
		return createRepo(ctx, newArticle)
	}
}

func NewGetAllArticlesService(getAllRepo repository.GetAllArticlesRepoFunc) GetAllArticlesServiceFunc {
	return func(ctx context.Context) ([]domain.Article, error) {
		return getAllRepo(ctx)
	}
}

func NewGetArticleByIDService(getByIDRepo repository.GetArticleByIDRepoFunc) GetArticleByIDServiceFunc {
	return func(ctx context.Context, id string) (domain.Article, error) {
		return getByIDRepo(ctx, id)
	}
}