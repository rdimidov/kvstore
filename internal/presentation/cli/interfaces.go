package cli

import (
	"context"

	"github.com/rdimidov/kvstore/internal/application/services"
	"github.com/rdimidov/kvstore/internal/domain"
)

type parser interface {
	Parse(raw string) (*services.Command, error)
}

type handler interface {
	Handle(ctx context.Context, cmd *services.Command) (*domain.Entry, error)
}
