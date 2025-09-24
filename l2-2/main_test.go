package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpack(t *testing.T) {
	var tests = []struct {
		in          string
		expected    string
		expectedErr error
	}{
		{
			in:          "a4bc2d5e",
			expected:    "aaaabccddddde",
			expectedErr: nil,
		},
		{
			in:          "abcd",
			expected:    "abcd",
			expectedErr: nil,
		},
		{
			in:          "45",
			expected:    "",
			expectedErr: fmt.Errorf("invalid string"),
		},
		{
			in:          "",
			expected:    "",
			expectedErr: nil,
		},
		{
			in:          `qwe\4\5`,
			expected:    "qwe45",
			expectedErr: nil,
		},
		{
			in:          `qwe\45`,
			expected:    "qwe44444",
			expectedErr: nil,
		},
		{
			in:          `qwe\\5`,
			expected:    `qwe\\\\\`,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := unpack(tt.in)
			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err, tt.expectedErr)
			}

			assert.Equal(t, tt.expected, out, "failed hahahahahahha")
		})
	}
}
