package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutationTypeString(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected MutationType
		valid    bool
	}{
		{
			name:     "mutation type unspecified",
			input:    "unspecified",
			expected: MutationTypeUnspecified,
			valid:    true,
		},
		{
			name:     "mutation type repayment",
			input:    "repayment",
			expected: MutationTypeRepayment,
			valid:    true,
		},
		{
			name:  "invalid",
			input: "invalid",
			valid: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := MutationTypeString(tc.input)
			assert.Equal(t, tc.expected, out)
			if tc.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestMutationType_String(t *testing.T) {
	testcases := []struct {
		name     string
		input    MutationType
		expected string
		valid    bool
	}{
		{
			name:     "mutation type unspecified",
			input:    MutationTypeUnspecified,
			expected: "unspecified",
			valid:    true,
		},
		{
			name:     "mutation type repayment",
			input:    MutationTypeRepayment,
			expected: "repayment",
			valid:    true,
		},
		{
			name:     "invalid",
			input:    MutationType(99),
			expected: "MutationType(99)",
			valid:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.String())
			assert.Equal(t, tc.valid, tc.input.IsAMutationType())
		})
	}
}

func TestMutationTypeValues(t *testing.T) {
	assert.Equal(t, MutationTypeValues(), []MutationType{
		MutationTypeUnspecified,
		MutationTypeRepayment,
	})
}

func TestMutationTypeStrings(t *testing.T) {
	// Get all available mutation type strings
	mutationType := MutationTypeStrings()

	// Test cases
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "should return all mutation type strings",
			want: []string{
				"unspecified",
				"repayment",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if length matches
			assert.Equal(t, len(tt.want), len(mutationType), "length of mutation type should match")

			// Check if all expected values are present
			for _, expectedType := range tt.want {
				assert.Contains(t, mutationType, expectedType, "should contain %s", expectedType)
			}

			// Check if no duplicates exist
			uniqueTypes := make(map[string]bool)
			for _, status := range mutationType {
				assert.False(t, uniqueTypes[status], "should not contain duplicates")
				uniqueTypes[status] = true
			}
		})
	}
}
