package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepaymentStatusString(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected RepaymentStatus
		valid    bool
	}{
		{
			name:     "repayment status unspecified",
			input:    "unspecified",
			expected: RepaymentStatusUnspecified,
			valid:    true,
		},
		{
			name:     "repayment status created",
			input:    "created",
			expected: RepaymentStatusCreated,
			valid:    true,
		},
		{
			name:     "repayment status paid",
			input:    "paid",
			expected: RepaymentStatusPaid,
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
			out, err := RepaymentStatusString(tc.input)
			assert.Equal(t, tc.expected, out)
			if tc.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestRepaymentStatus_String(t *testing.T) {
	testcases := []struct {
		name     string
		input    RepaymentStatus
		expected string
		valid    bool
	}{
		{
			name:     "repayment status unspecified",
			input:    RepaymentStatusUnspecified,
			expected: "unspecified",
			valid:    true,
		},
		{
			name:     "repayment status created",
			input:    RepaymentStatusCreated,
			expected: "created",
			valid:    true,
		},
		{
			name:     "repayment status paid",
			input:    RepaymentStatusPaid,
			expected: "paid",
			valid:    true,
		},
		{
			name:     "invalid",
			input:    RepaymentStatus(99),
			expected: "RepaymentStatus(99)",
			valid:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.String())
			assert.Equal(t, tc.valid, tc.input.IsARepaymentStatus())
		})
	}
}

func TestRepaymentStatusValues(t *testing.T) {
	assert.Equal(t, RepaymentStatusValues(), []RepaymentStatus{
		RepaymentStatusUnspecified,
		RepaymentStatusCreated,
		RepaymentStatusPaid,
	})
}

func TestRepaymentStatusStrings(t *testing.T) {
	// Get all available repayment status strings
	repaymentStatus := RepaymentStatusStrings()

	// Test cases
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "should return all repayment status strings",
			want: []string{
				"unspecified",
				"created",
				"paid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if length matches
			assert.Equal(t, len(tt.want), len(repaymentStatus), "length of repayment status should match")

			// Check if all expected values are present
			for _, expectedType := range tt.want {
				assert.Contains(t, repaymentStatus, expectedType, "should contain %s", expectedType)
			}

			// Check if no duplicates exist
			uniqueTypes := make(map[string]bool)
			for _, status := range repaymentStatus {
				assert.False(t, uniqueTypes[status], "should not contain duplicates")
				uniqueTypes[status] = true
			}
		})
	}
}
