package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoanStatusString(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected LoanStatus
		valid    bool
	}{
		{
			name:     "loan status unspecified",
			input:    "unspecified",
			expected: LoanStatusUnspecified,
			valid:    true,
		},
		{
			name:     "loan status proposed",
			input:    "proposed",
			expected: LoanStatusProposed,
			valid:    true,
		},
		{
			name:     "loan status approved",
			input:    "approved",
			expected: LoanStatusApproved,
			valid:    true,
		},
		{
			name:     "loan status invested",
			input:    "invested",
			expected: LoanStatusInvested,
			valid:    true,
		},
		{
			name:     "loan status disbursed",
			input:    "disbursed",
			expected: LoanStatusDisbursed,
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
			out, err := LoanStatusString(tc.input)
			assert.Equal(t, tc.expected, out)
			if tc.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestLoanStatus_String(t *testing.T) {
	testcases := []struct {
		name     string
		input    LoanStatus
		expected string
		valid    bool
	}{
		{
			name:     "loan status unspecified",
			input:    LoanStatusUnspecified,
			expected: "unspecified",
			valid:    true,
		},
		{
			name:     "loan status proposed",
			input:    LoanStatusProposed,
			expected: "proposed",
			valid:    true,
		},
		{
			name:     "loan status approved",
			input:    LoanStatusApproved,
			expected: "approved",
			valid:    true,
		},
		{
			name:     "loan status invested",
			input:    LoanStatusInvested,
			expected: "invested",
			valid:    true,
		},
		{
			name:     "loan status disbursed",
			input:    LoanStatusDisbursed,
			expected: "disbursed",
			valid:    true,
		},
		{
			name:     "invalid",
			input:    LoanStatus(99),
			expected: "LoanStatus(99)",
			valid:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.String())
			assert.Equal(t, tc.valid, tc.input.IsALoanStatus())
		})
	}
}

func TestLoanStatusValues(t *testing.T) {
	assert.Equal(t, LoanStatusValues(), []LoanStatus{
		LoanStatusUnspecified,
		LoanStatusProposed,
		LoanStatusApproved,
		LoanStatusInvested,
		LoanStatusDisbursed,
	})
}

func TestLoanStatusStrings(t *testing.T) {
	// Get all available loan status strings
	loanStatus := LoanStatusStrings()

	// Test cases
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "should return all loan status strings",
			want: []string{
				"unspecified",
				"proposed",
				"approved",
				"invested",
				"disbursed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if length matches
			assert.Equal(t, len(tt.want), len(loanStatus), "length of loan status should match")

			// Check if all expected values are present
			for _, expectedType := range tt.want {
				assert.Contains(t, loanStatus, expectedType, "should contain %s", expectedType)
			}

			// Check if no duplicates exist
			uniqueTypes := make(map[string]bool)
			for _, status := range loanStatus {
				assert.False(t, uniqueTypes[status], "should not contain duplicates")
				uniqueTypes[status] = true
			}
		})
	}
}
