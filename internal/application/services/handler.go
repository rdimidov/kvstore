package services

import (
	"context"
	"errors"

	"github.com/rdimidov/kvstore/internal/domain"
)

var ErrInvalidCmd = errors.New("invalid command")

// StringHandler acts as an adapter between raw string commands and business logic
type StringHandler struct {
	compute compute
}

func NewStringHandler(compute compute) (*StringHandler, error) {
	if compute == nil {
		return nil, errors.New("compute is nil")
	}
	return &StringHandler{compute: compute}, nil
}

// Handle parses and routes a parsed command to the appropriate compute method
func (s *StringHandler) Handle(ctx context.Context, cmd *Command) (*domain.Entry, error) {
	switch cmd.Cmd {
	case getCommand:
		return s.compute.Get(ctx, cmd.Key)
	case setCommand:
		return nil, s.compute.Set(ctx, cmd.Key, *cmd.Value)
	case delCommand:
		return nil, s.compute.Delete(ctx, cmd.Key)
	default:
		return nil, ErrInvalidCmd
	}
}
