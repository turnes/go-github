package domain

import "time"

type Repository struct {
	ID          int64
	Name        string
	FullName    string
	Description string
	Private     bool
	URL         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateRepositoryInput struct {
	Name        string
	Description string
	Private     bool
}

type UpdateRepositoryInput struct {
	Name        string
	Description string
	Private     bool
}
