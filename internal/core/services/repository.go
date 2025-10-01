package services

import (
	"context"
	"errors"

	"github.com/turnes/go-github/internal/core/domain"
	"github.com/turnes/go-github/internal/core/ports"
)

type repositoryService struct {
	repo ports.RepositoryPort
}

func NewRepositoryService(repo ports.RepositoryPort) ports.RepositoryService {
	return &repositoryService{
		repo: repo,
	}
}

func (s *repositoryService) ListRepositories(ctx context.Context, owner string) ([]*domain.Repository, error) {
	if owner == "" {
		return nil, errors.New("owner cannot be empty")
	}
	return s.repo.List(ctx, owner)
}

func (s *repositoryService) CreateRepository(ctx context.Context, owner string, input domain.CreateRepositoryInput) (*domain.Repository, error) {
	if owner == "" {
		return nil, errors.New("owner cannot be empty")
	}
	if input.Name == "" {
		return nil, errors.New("repository name cannot be empty")
	}
	return s.repo.Create(ctx, owner, input)
}

func (s *repositoryService) DeleteRepository(ctx context.Context, owner, name string) error {
	if owner == "" || name == "" {
		return errors.New("owner and name cannot be empty")
	}
	return s.repo.Delete(ctx, owner, name)
}

func (s *repositoryService) UpdateRepository(ctx context.Context, owner, name string, input domain.UpdateRepositoryInput) (*domain.Repository, error) {
	if owner == "" || name == "" {
		return nil, errors.New("owner and name cannot be empty")
	}
	if input.Name == "" {
		return nil, errors.New("repository name cannot be empty")
	}
	return s.repo.Update(ctx, owner, name, input)
}

func (s *repositoryService) GetRepository(ctx context.Context, owner, name string) (*domain.Repository, error) {
	if owner == "" || name == "" {
		return nil, errors.New("owner and name cannot be empty")
	}
	return s.repo.Get(ctx, owner, name)
}
