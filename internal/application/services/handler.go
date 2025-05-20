package services

import (
	"context"
	"errors"

	"github.com/rdimidov/kvstore/internal/domain"
)

var ErrInvalidCmd = errors.New("invalid command")

type application interface {
	Get(ctx context.Context, key domain.Key) (*domain.Entry, error)
	Set(ctx context.Context, key domain.Key, value domain.Value) error
	Delete(ctx context.Context, key domain.Key) error
}

// Handler acts as an adapter between commands and business logic
type Handler struct {
	application application
}

func NewHandler(application application) (*Handler, error) {
	if application == nil {
		return nil, errors.New("application is nil")
	}
	return &Handler{application: application}, nil
}

// Handle parses and routes a parsed command to the appropriate application method
func (s *Handler) Handle(ctx context.Context, cmd *Command) (*domain.Entry, error) {
	switch cmd.Cmd {
	case getCommand:
		return s.application.Get(ctx, cmd.Key)
	case setCommand:
		return nil, s.application.Set(ctx, cmd.Key, *cmd.Value)
	case delCommand:
		return nil, s.application.Delete(ctx, cmd.Key)
	default:
		return nil, ErrInvalidCmd
	}
}
