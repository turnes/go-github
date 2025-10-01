package repository

import (
	"context"

	"github.com/google/go-github/v75/github"
	"github.com/turnes/go-github/internal/core/domain"
	"github.com/turnes/go-github/internal/core/ports"
)

type githubRepository struct {
	client *github.Client
}

func NewGithubRepository(client *github.Client) ports.RepositoryPort {
	return &githubRepository{
		client: client,
	}
}

func (g *githubRepository) List(ctx context.Context, owner string) ([]*domain.Repository, error) {
	repos, _, err := g.client.Repositories.List(ctx, owner, &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Repository, 0, len(repos))
	for _, r := range repos {
		result = append(result, g.toDomain(r))
	}
	return result, nil
}

func (g *githubRepository) Create(ctx context.Context, owner string, input domain.CreateRepositoryInput) (*domain.Repository, error) {
	repo := &github.Repository{
		Name:        github.String(input.Name),
		Description: github.String(input.Description),
		Private:     github.Bool(input.Private),
	}

	created, _, err := g.client.Repositories.Create(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return g.toDomain(created), nil
}

func (g *githubRepository) Delete(ctx context.Context, owner, name string) error {
	_, err := g.client.Repositories.Delete(ctx, owner, name)
	return err
}

func (g *githubRepository) Update(ctx context.Context, owner, name string, input domain.UpdateRepositoryInput) (*domain.Repository, error) {
	repo := &github.Repository{
		Name: github.Ptr(input.Name),
	}
	if input.Description != "" {
		repo.Description = github.Ptr(input.Description)
	}
	if input.Private {
		repo.Private = github.Ptr(input.Private)
	}

	updated, _, err := g.client.Repositories.Edit(ctx, owner, name, repo)
	if err != nil {
		return nil, err
	}

	return g.toDomain(updated), nil
}

func (g *githubRepository) Get(ctx context.Context, owner, name string) (*domain.Repository, error) {
	repo, _, err := g.client.Repositories.Get(ctx, owner, name)
	if err != nil {
		return nil, err
	}
	return g.toDomain(repo), nil
}

func (g *githubRepository) toDomain(r *github.Repository) *domain.Repository {
	return &domain.Repository{
		ID:          r.GetID(),
		Name:        r.GetName(),
		FullName:    r.GetFullName(),
		Description: r.GetDescription(),
		Private:     r.GetPrivate(),
		URL:         r.GetHTMLURL(),
		CreatedAt:   r.GetCreatedAt().Time,
		UpdatedAt:   r.GetUpdatedAt().Time,
	}
}
