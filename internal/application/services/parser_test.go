package services

import (
	"testing"

	"github.com/rdimidov/kvstore/internal/domain"
	"github.com/stretchr/testify/assert"
)

func mustKey(t *testing.T, s string) domain.Key {
	t.Helper()
	k, err := domain.NewKey(s)
	if err != nil {
		t.Fatalf("invalid key in test: %v", err)
	}
	return k
}

func mustValue(t *testing.T, s string) *domain.Value {
	t.Helper()
	v, err := domain.NewValue(s)
	if err != nil {
		t.Fatalf("invalid value in test: %v", err)
	}
	return &v
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    *Command
		wantErr error
	}{
		{
			name: "valid GET",
			raw:  "GET foo",
			want: &Command{Cmd: "GET", Key: mustKey(t, "foo")},
		},
		{
			name: "valid DEL",
			raw:  "DEL bar",
			want: &Command{Cmd: "DEL", Key: mustKey(t, "bar")},
		},
		{
			name: "valid SET",
			raw:  "SET x y",
			want: &Command{
				Cmd:   "SET",
				Key:   mustKey(t, "x"),
				Value: mustValue(t, "y"),
			},
		},
		{
			name:    "unknown command",
			raw:     "FOO bar",
			wantErr: ErrInvalidCmd,
		},
		{
			name:    "too few args",
			raw:     "GET",
			wantErr: ErrInvalidCmd,
		},
		{
			name:    "too many args for GET",
			raw:     "GET a b",
			wantErr: ErrInvalidCmd,
		},
		{
			name:    "SET missing value",
			raw:     "SET key",
			wantErr: ErrInvalidCmd,
		},
		{
			name:    "empty input",
			raw:     "",
			wantErr: ErrInvalidCmd,
		},
		{
			name:    "spaces only",
			raw:     "   ",
			wantErr: ErrInvalidCmd,
		},
	}

	p := Parser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.raw)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Cmd, got.Cmd)
				assert.Equal(t, tt.want.Key, got.Key)
				if tt.want.Value != nil {
					assert.NotNil(t, got.Value)
					assert.Equal(t, *tt.want.Value, *got.Value)
				} else {
					assert.Nil(t, got.Value)
				}
			}
		})
	}
}
