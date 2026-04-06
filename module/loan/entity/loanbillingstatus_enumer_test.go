package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoanBillingStatusString(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected LoanBillingStatus
		valid    bool
	}{
		{
			name:     "loan billing status unspecified",
			input:    "unspecified",
			expected: LoanBillingStatusUnspecified,
			valid:    true,
		},
		{
			name:     "loan billing status created",
			input:    "created",
			expected: LoanBillingStatusCreated,
			valid:    true,
		},
		{
			name:     "loan billing status paid",
			input:    "paid",
			expected: LoanBillingStatusPaid,
			valid:    true,
		},
		{
			name:     "loan billing status overdue",
			input:    "overdue",
			expected: LoanBillingStatusOverdue,
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
			out, err := LoanBillingStatusString(tc.input)
			assert.Equal(t, tc.expected, out)
			if tc.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestLoanBillingStatus_String(t *testing.T) {
	testcases := []struct {
		name     string
		input    LoanBillingStatus
		expected string
		valid    bool
	}{
		{
			name:     "loan billing status unspecified",
			input:    LoanBillingStatusUnspecified,
			expected: "unspecified",
			valid:    true,
		},
		{
			name:     "loan billing status created",
			input:    LoanBillingStatusCreated,
			expected: "created",
			valid:    true,
		},
		{
			name:     "loan billing status paid",
			input:    LoanBillingStatusPaid,
			expected: "paid",
			valid:    true,
		},
		{
			name:     "loan billing status overdue",
			input:    LoanBillingStatusOverdue,
			expected: "overdue",
			valid:    true,
		},
		{
			name:     "invalid",
			input:    LoanBillingStatus(99),
			expected: "LoanBillingStatus(99)",
			valid:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.String())
			assert.Equal(t, tc.valid, tc.input.IsALoanBillingStatus())
		})
	}
}

func TestLoanBillingStatusValues(t *testing.T) {
	assert.Equal(t, LoanBillingStatusValues(), []LoanBillingStatus{
		LoanBillingStatusUnspecified,
		LoanBillingStatusCreated,
		LoanBillingStatusPaid,
		LoanBillingStatusOverdue,
	})
}

func TestLoanBillingStatusStrings(t *testing.T) {
	// Get all available loan billing status strings
	loanBillingStatus := LoanBillingStatusStrings()

	// Test cases
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "should return all loan billing status strings",
			want: []string{
				"unspecified",
				"created",
				"paid",
				"overdue",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if length matches
			assert.Equal(t, len(tt.want), len(loanBillingStatus), "length of loan billing status should match")

			// Check if all expected values are present
			for _, expectedType := range tt.want {
				assert.Contains(t, loanBillingStatus, expectedType, "should contain %s", expectedType)
			}

			// Check if no duplicates exist
			uniqueTypes := make(map[string]bool)
			for _, status := range loanBillingStatus {
				assert.False(t, uniqueTypes[status], "should not contain duplicates")
				uniqueTypes[status] = true
			}
		})
	}
}
