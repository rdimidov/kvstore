package services

import (
	"context"
	"errors"
	"testing"

	"github.com/rdimidov/kvstore/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCompute_Set(t *testing.T) {
	t.Parallel()

	type args struct {
		key   domain.Key
		value domain.Value
	}

	tests := []struct {
		name        string
		args        args
		mockSetup   func(r *mockrepository)
		expectError bool
	}{
		{
			name: "successfully sets key",
			args: args{key: "foo", value: "bar"},
			mockSetup: func(r *mockrepository) {
				r.On("Set", mock.Anything, domain.Key("foo"), domain.Value("bar")).Return(nil)
			},
		},
		{
			name: "repo failure",
			args: args{key: "foo", value: "bar"},
			mockSetup: func(r *mockrepository) {
				r.On("Set", mock.Anything, domain.Key("foo"), domain.Value("bar")).Return(errors.New("oops.."))
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := newMockrepository(t)
			tt.mockSetup(mockRepo)

			app := NewApplication(mockRepo, zap.NewNop().Sugar())
			err := app.Set(context.Background(), tt.args.key, tt.args.value)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCompute_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		key domain.Key
	}

	tests := []struct {
		name        string
		args        args
		mockSetup   func(r *mockrepository)
		expectError bool
	}{
		{
			name: "successfully gets key",
			args: args{key: "foo"},
			mockSetup: func(r *mockrepository) {
				r.On("Get", mock.Anything, domain.Key("foo")).Return(&domain.Entry{
					Key: domain.Key("foo"), Value: domain.Value("bar"),
				}, nil)
			},
		},
		{
			name: "repo failure",
			args: args{key: "foo"},
			mockSetup: func(r *mockrepository) {
				r.On("Get", mock.Anything, domain.Key("foo")).Return(nil, errors.New("oops.."))
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := newMockrepository(t)
			tt.mockSetup(mockRepo)

			app := NewApplication(mockRepo, zap.NewNop().Sugar())
			entry, err := app.Get(context.Background(), tt.args.key)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, entry.Key, tt.args.key)
			}
		})
	}
}

func TestCompute_Delete(t *testing.T) {
	t.Parallel()

	type args struct {
		key domain.Key
	}

	tests := []struct {
		name        string
		args        args
		mockSetup   func(r *mockrepository)
		expectError bool
	}{
		{
			name: "successfully deletes key",
			args: args{key: "foo"},
			mockSetup: func(r *mockrepository) {
				r.On("Delete", mock.Anything, domain.Key("foo")).Return(nil)
			},
		},
		{
			name: "repo failure",
			args: args{key: "foo"},
			mockSetup: func(r *mockrepository) {
				r.On("Delete", mock.Anything, domain.Key("foo")).Return(errors.New("oops.."))
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := newMockrepository(t)
			tt.mockSetup(mockRepo)

			app := NewApplication(mockRepo, zap.NewNop().Sugar())

			err := app.Delete(context.Background(), tt.args.key)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
