package user

import (
	"context"

	"github.com/Nikitapopov/Habbit/pkg/logging"
)

type Service interface {
	Create(ctx context.Context, dto CreateUserDTO) (string, error)
}

type service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository *Repository, logger *logging.Logger) Service {
	return &service{
		repository: *repository,
		logger:     logger,
	}
}

func (s *service) Create(ctx context.Context, dto CreateUserDTO) (string, error) {
	// TODO приводим к User, но значения все еще из CreateUserDTO
	return s.repository.Create(ctx, &User{
		Username:     dto.Username,
		PasswordHash: dto.Password,
		Email:        dto.Email,
	})
}
