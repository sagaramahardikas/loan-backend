package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountStatusString(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected AccountStatus
		valid    bool
	}{
		{
			name:     "account status unspecified",
			input:    "unspecified",
			expected: AccountStatusUnspecified,
			valid:    true,
		},
		{
			name:     "account status inactive",
			input:    "inactive",
			expected: AccountStatusInactive,
			valid:    true,
		},
		{
			name:     "account status active",
			input:    "active",
			expected: AccountStatusActive,
			valid:    true,
		},
		{
			name:     "account status frozen",
			input:    "frozen",
			expected: AccountStatusFrozen,
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
			out, err := AccountStatusString(tc.input)
			assert.Equal(t, tc.expected, out)
			if tc.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestAccountStatus_String(t *testing.T) {
	testcases := []struct {
		name     string
		input    AccountStatus
		expected string
		valid    bool
	}{
		{
			name:     "account status unspecified",
			input:    AccountStatusUnspecified,
			expected: "unspecified",
			valid:    true,
		},
		{
			name:     "account status inactive",
			input:    AccountStatusInactive,
			expected: "inactive",
			valid:    true,
		},
		{
			name:     "account status active",
			input:    AccountStatusActive,
			expected: "active",
			valid:    true,
		},
		{
			name:     "account status frozen",
			input:    AccountStatusFrozen,
			expected: "frozen",
			valid:    true,
		},
		{
			name:     "invalid",
			input:    AccountStatus(99),
			expected: "AccountStatus(99)",
			valid:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.String())
			assert.Equal(t, tc.valid, tc.input.IsAAccountStatus())
		})
	}
}

func TestAccountStatusValues(t *testing.T) {
	assert.Equal(t, AccountStatusValues(), []AccountStatus{
		AccountStatusUnspecified,
		AccountStatusInactive,
		AccountStatusActive,
		AccountStatusFrozen,
	})
}

func TestAccountStatusStrings(t *testing.T) {
	// Get all available account status strings
	accountStatus := AccountStatusStrings()

	// Test cases
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "should return all account status strings",
			want: []string{
				"unspecified",
				"inactive",
				"active",
				"frozen",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if length matches
			assert.Equal(t, len(tt.want), len(accountStatus), "length of account status should match")

			// Check if all expected values are present
			for _, expectedType := range tt.want {
				assert.Contains(t, accountStatus, expectedType, "should contain %s", expectedType)
			}

			// Check if no duplicates exist
			uniqueTypes := make(map[string]bool)
			for _, status := range accountStatus {
				assert.False(t, uniqueTypes[status], "should not contain duplicates")
				uniqueTypes[status] = true
			}
		})
	}
}
