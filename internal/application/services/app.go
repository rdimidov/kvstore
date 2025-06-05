package services

import (
	"context"
	"errors"

	"github.com/rdimidov/kvstore/internal/domain"

	"go.uber.org/zap"
)

type repository interface {
	Set(context.Context, domain.Key, domain.Value) error
	Get(context.Context, domain.Key) (*domain.Entry, error)
	Delete(context.Context, domain.Key) error
}

type WALogger interface {
	WriteSet(domain.Key, domain.Value) error
	WriteDel(domain.Key) error
	Recover(ctx context.Context) error
}

// Application defines application-level operations and coordinates between
// the domain logic and the data persistence layer.
type Application struct {
	repo   repository
	wal    WALogger
	logger *zap.SugaredLogger
}

func NewApplication(ctx context.Context, repo repository, logger *zap.SugaredLogger, wal WALogger) (*Application, error) {
	if wal != nil {
		if err := wal.Recover(ctx); err != nil {
			return nil, err
		}
	}
	return &Application{
		repo:   repo,
		wal:    wal,
		logger: logger,
	}, nil
}

func (c *Application) Set(ctx context.Context, key domain.Key, value domain.Value) error {
	c.logger.Debugw("setting", "key", key, "value", value)

	if c.wal != nil {
		if err := c.wal.WriteSet(key, value); err != nil {
			return err
		}
	}

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

	if c.wal != nil {
		if err := c.wal.WriteDel(key); err != nil {
			return err
		}
	}

	err := c.repo.Delete(ctx, key)
	if err != nil && !errors.Is(err, domain.ErrKeyNotFound) {
		c.logger.Errorf("failed to delete key: %s, err: %v", key, err)
	}
	return err
}
