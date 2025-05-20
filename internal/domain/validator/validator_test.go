package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
	}{
		{
			name:    "single letter",
			input:   "a",
			isValid: true,
		},
		{
			name:    "single digit",
			input:   "5",
			isValid: true,
		},
		{
			name:    "single punctuation star",
			input:   "*",
			isValid: true,
		},
		{
			name:    "single punctuation slash",
			input:   "/",
			isValid: true,
		},
		{
			name:    "underscore",
			input:   "_",
			isValid: true,
		},
		{
			name:    "multiple mixed valid",
			input:   "a1_*Z",
			isValid: true,
		},
		{
			name:    "empty string",
			input:   "",
			isValid: false,
		},
		{
			name:    "space in string",
			input:   "a b",
			isValid: false,
		},
		{
			name:    "invalid dash",
			input:   "a-b",
			isValid: false,
		},
		{
			name:    "invalid at sign",
			input:   "@home",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidString(tt.input)
			require.Equal(t, tt.isValid, got)
		})
	}
}
