package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStatusString(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected UserStatus
		valid    bool
	}{
		{
			name:     "user status unspecified",
			input:    "unspecified",
			expected: UserStatusUnspecified,
			valid:    true,
		},
		{
			name:     "user status inactive",
			input:    "inactive",
			expected: UserStatusInactive,
			valid:    true,
		},
		{
			name:     "user status active",
			input:    "active",
			expected: UserStatusActive,
			valid:    true,
		},
		{
			name:     "user status delinquent",
			input:    "delinquent",
			expected: UserStatusDelinquent,
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
			out, err := UserStatusString(tc.input)
			assert.Equal(t, tc.expected, out)
			if tc.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestUserStatus_String(t *testing.T) {
	testcases := []struct {
		name     string
		input    UserStatus
		expected string
		valid    bool
	}{
		{
			name:     "user status unspecified",
			input:    UserStatusUnspecified,
			expected: "unspecified",
			valid:    true,
		},
		{
			name:     "user status inactive",
			input:    UserStatusInactive,
			expected: "inactive",
			valid:    true,
		},
		{
			name:     "user status active",
			input:    UserStatusActive,
			expected: "active",
			valid:    true,
		},
		{
			name:     "user status delinquent",
			input:    UserStatusDelinquent,
			expected: "delinquent",
			valid:    true,
		},
		{
			name:     "invalid",
			input:    UserStatus(99),
			expected: "UserStatus(99)",
			valid:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.String())
			assert.Equal(t, tc.valid, tc.input.IsAUserStatus())
		})
	}
}

func TestUserStatusValues(t *testing.T) {
	assert.Equal(t, UserStatusValues(), []UserStatus{
		UserStatusUnspecified,
		UserStatusInactive,
		UserStatusActive,
		UserStatusDelinquent,
	})
}

func TestUserStatusStrings(t *testing.T) {
	// Get all available user status strings
	userStatus := UserStatusStrings()

	// Test cases
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "should return all user status strings",
			want: []string{
				"unspecified",
				"inactive",
				"active",
				"delinquent",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if length matches
			assert.Equal(t, len(tt.want), len(userStatus), "length of user status should match")

			// Check if all expected values are present
			for _, expectedType := range tt.want {
				assert.Contains(t, userStatus, expectedType, "should contain %s", expectedType)
			}

			// Check if no duplicates exist
			uniqueTypes := make(map[string]bool)
			for _, status := range userStatus {
				assert.False(t, uniqueTypes[status], "should not contain duplicates")
				uniqueTypes[status] = true
			}
		})
	}
}
