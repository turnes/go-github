package ports

import (
	"context"

	"github.com/turnes/go-github/internal/core/domain"
)

type RepositoryService interface {
	ListRepositories(ctx context.Context, owner string) ([]*domain.Repository, error)
	CreateRepository(ctx context.Context, owner string, input domain.CreateRepositoryInput) (*domain.Repository, error)
	DeleteRepository(ctx context.Context, owner, name string) error
	UpdateRepository(ctx context.Context, owner, name string, input domain.UpdateRepositoryInput) (*domain.Repository, error)
	GetRepository(ctx context.Context, owner, name string) (*domain.Repository, error)
}

type RepositoryPort interface {
	List(ctx context.Context, owner string) ([]*domain.Repository, error)
	Create(ctx context.Context, owner string, input domain.CreateRepositoryInput) (*domain.Repository, error)
	Delete(ctx context.Context, owner, name string) error
	Update(ctx context.Context, owner, name string, input domain.UpdateRepositoryInput) (*domain.Repository, error)
	Get(ctx context.Context, owner, name string) (*domain.Repository, error)
}
