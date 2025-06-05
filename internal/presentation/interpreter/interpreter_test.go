package interpreter

import (
	"context"
	"testing"

	"github.com/rdimidov/kvstore/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInterpreter_Execute(t *testing.T) {
	// sample key and value
	key, _ := domain.NewKey("foo")
	val, _ := domain.NewValue("bar")
	entry := &domain.Entry{Key: key, Value: val}

	tests := []struct {
		name      string
		input     string
		setup     func(app *mockhandler)
		wantEntry *domain.Entry
		wantErr   error
	}{
		{
			name:  "GET success",
			input: "GET foo",
			setup: func(app *mockhandler) {
				app.On("Get", mock.Anything, key).Return(entry, nil)
			},
			wantEntry: entry,
			wantErr:   nil,
		},
		{
			name:    "GET invalid args",
			input:   "GET",
			setup:   func(app *mockhandler) {},
			wantErr: ErrInvalidCmd,
		},
		{
			name:    "GET more invalid args",
			input:   "GET foo foo",
			setup:   func(app *mockhandler) {},
			wantErr: ErrInvalidCmd,
		},
		{
			name:  "DEL success",
			input: "DEL foo",
			setup: func(app *mockhandler) {
				app.On("Delete", mock.Anything, key).Return(nil)
			},
			wantEntry: nil,
			wantErr:   nil,
		},
		{
			name:    "DEL invalid args",
			input:   "DEL bar foo",
			setup:   func(app *mockhandler) {},
			wantErr: ErrInvalidCmd,
		},
		{
			name:  "SET success",
			input: "SET foo bar",
			setup: func(app *mockhandler) {
				app.On("Set", mock.Anything, key, val).Return(nil)
			},
			wantEntry: nil,
			wantErr:   nil,
		},
		{
			name:    "SET invalid args",
			input:   "SET foo",
			setup:   func(app *mockhandler) {},
			wantErr: ErrInvalidCmd,
		},
		{
			name:    "Unknown command",
			input:   "FOO foo",
			setup:   func(app *mockhandler) {},
			wantErr: ErrInvalidCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create mock
			appMock := newMockhandler(t)
			tt.setup(appMock)

			interp, err := New(appMock)
			assert.NoError(t, err)

			gotEntry, err := interp.Execute(context.Background(), tt.input)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantEntry, gotEntry)
			}
			appMock.AssertExpectations(t)
		})
	}
}
