package services

import (
	"context"
	"errors"

	"github.com/rdimidov/kvstore/internal/domain"
	"go.uber.org/zap"
)

type repository interface {
	Set(ctx context.Context, key domain.Key, value domain.Value) error
	Get(ctx context.Context, key domain.Key) (*domain.Entry, error)
	Delete(ctx context.Context, key domain.Key) error
}

// Application defines application-level operations and coordinates between
// the domain logic and the data persistence layer.
type Application struct {
	repo   repository
	logger *zap.SugaredLogger
}

func NewApplication(repo repository, logger *zap.SugaredLogger) *Application {
	return &Application{
		repo:   repo,
		logger: logger,
	}
}
func (c *Application) Set(ctx context.Context, key domain.Key, value domain.Value) error {
	c.logger.Debugw("setting", "key", key, "value", value)
	err := c.repo.Set(ctx, key, value)
	if err != nil {
		c.logger.Errorf("failed to set key: %s, err: %v", key, err)
	}
	return err
}

func (c *Application) Get(ctx context.Context, key domain.Key) (*domain.Entry, error) {
	c.logger.Debugw("getting", "key", key)
	entry, err := c.repo.Get(ctx, key)
	if err != nil && !errors.Is(err, domain.ErrKeyNotFound) {
		c.logger.Errorf("failed to get key: %s, err: %v", key, err)
	}
	return entry, err
}

func (c *Application) Delete(ctx context.Context, key domain.Key) error {
	c.logger.Debugw("deleting", "key", key)
	err := c.repo.Delete(ctx, key)
	if err != nil && !errors.Is(err, domain.ErrKeyNotFound) {
		c.logger.Errorf("failed to delete key: %s, err: %v", key, err)
	}
	return err
}
