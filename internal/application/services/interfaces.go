package services

import (
	"context"

	"github.com/rdimidov/kvstore/internal/domain"
)

type compute interface {
	Get(ctx context.Context, key domain.Key) (*domain.Entry, error)
	Set(ctx context.Context, key domain.Key, value domain.Value) error
	Delete(ctx context.Context, key domain.Key) error
}

type repository interface {
	Set(ctx context.Context, key domain.Key, value domain.Value) error
	Get(ctx context.Context, key domain.Key) (*domain.Entry, error)
	Delete(ctx context.Context, key domain.Key) error
}
