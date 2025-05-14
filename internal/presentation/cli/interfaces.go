package cli

import (
	"context"

	"github.com/rdimidov/kvstore/internal/domain"
)

type getter interface {
	Get(ctx context.Context, key domain.Key) (*domain.Entry, error)
}

type setter interface {
	Set(ctx context.Context, key domain.Key, value domain.Value) error
}

type deleter interface {
	Delete(ctx context.Context, key domain.Key) error
}

type appController interface {
	getter
	setter
	deleter
}
