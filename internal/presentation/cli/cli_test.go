package cli

import (
	"bytes"
	"errors"
	"testing"

	"github.com/rdimidov/kvstore/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDelCmd(t *testing.T) {
	type testCase struct {
		name      string
		args      []string
		mockSetup func(m *mockdeleter)
		wantErr   string
	}

	validKey, _ := domain.NewKey("foo")

	tests := []testCase{
		{
			name: "success",
			args: []string{"foo"},
			mockSetup: func(m *mockdeleter) {
				m.On("Delete", mock.Anything, validKey).Return(nil)
			},
		},
		{
			name: "delete error",
			args: []string{"foo"},
			mockSetup: func(m *mockdeleter) {
				m.On("Delete", mock.Anything, validKey).Return(errors.New("some error"))
			},
			wantErr: "could not delete entry",
		},
		{
			name:      "invalid key format",
			args:      []string{""},
			mockSetup: nil, // Delete не вызывается
			wantErr:   "invalid key format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newMockdeleter(t)
			if tc.mockSetup != nil {
				tc.mockSetup(m)
			}

			cmd := NewDelCmd(m)
			out := &bytes.Buffer{}
			cmd.SetErr(out)
			cmd.SetArgs(tc.args)

			_ = cmd.Execute()

			if tc.wantErr != "" {
				assert.Contains(t, out.String(), tc.wantErr)
			} else {
				assert.Empty(t, out.String())
			}
		})
	}
}

func TestSetCmd(t *testing.T) {
	type testCase struct {
		name      string
		args      []string
		mockSetup func(m *mocksetter)
		wantErr   string
	}

	validKey, _ := domain.NewKey("foo")
	validValue, _ := domain.NewValue("bar")

	tests := []testCase{
		{
			name: "success",
			args: []string{"foo", "bar"},
			mockSetup: func(m *mocksetter) {
				m.On("Set", mock.Anything, validKey, validValue).Return(nil)
			},
		},
		{
			name:    "invalid key",
			args:    []string{"", "bar"},
			wantErr: "invalid key format",
		},
		{
			name:    "invalid value",
			args:    []string{"foo", ""},
			wantErr: "invalid value format",
		},
		{
			name: "set error",
			args: []string{"foo", "bar"},
			mockSetup: func(m *mocksetter) {
				m.On("Set", mock.Anything, validKey, validValue).
					Return(errors.New("save failed"))
			},
			wantErr: "could not save entry",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newMocksetter(t)
			if tc.mockSetup != nil {
				tc.mockSetup(m)
			}

			cmd := NewSetCmd(m)
			errOut := &bytes.Buffer{}
			cmd.SetErr(errOut)
			cmd.SetArgs(tc.args)

			_ = cmd.Execute()

			if tc.wantErr != "" {
				assert.Contains(t, errOut.String(), tc.wantErr)
			} else {
				assert.Empty(t, errOut.String())
			}
		})
	}
}

func TestGetCmd(t *testing.T) {
	type testCase struct {
		name       string
		args       []string
		mockSetup  func(m *mockgetter)
		wantOutput string
		wantErr    string
	}

	keyFoo, _ := domain.NewKey("foo")

	tests := []testCase{
		{
			name: "success",
			args: []string{"foo"},
			mockSetup: func(m *mockgetter) {
				m.On("Get", mock.Anything, keyFoo).
					Return(&domain.Entry{Key: keyFoo, Value: "bar"}, nil)
			},
			wantOutput: "bar",
		},
		{
			name: "key not found",
			args: []string{"foo"},
			mockSetup: func(m *mockgetter) {
				m.On("Get", mock.Anything, keyFoo).
					Return(nil, domain.ErrKeyNotFound)
			},
			wantErr: "key not found",
		},
		{
			name: "internal error",
			args: []string{"foo"},
			mockSetup: func(m *mockgetter) {
				m.On("Get", mock.Anything, keyFoo).
					Return(nil, errors.New("boom"))
			},
			wantErr: "error: boom",
		},
		{
			name:      "invalid key format",
			args:      []string{""},
			mockSetup: nil,
			wantErr:   "invalid key format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newMockgetter(t)
			if tc.mockSetup != nil {
				tc.mockSetup(m)
			}

			cmd := NewGetCmd(m)
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			cmd.SetArgs(tc.args)

			_ = cmd.Execute()

			if tc.wantOutput != "" {
				assert.Contains(t, stdout.String(), tc.wantOutput)
			}
			if tc.wantErr != "" {
				assert.Contains(t, stderr.String(), tc.wantErr)
			}
		})
	}
}
