package rpcmethods

import (
	"testing"
)

func TestDefaultMethodFormatter(t *testing.T) {
	testCases := []struct {
		category string
		action   string
		expected string
	}{
		{
			category: "resource",
			action:   "read",
			expected: "method.resource.read",
		},
		{
			category: "tool",
			action:   "call",
			expected: "method.tool.call",
		},
		{
			category: "test",
			action:   "action",
			expected: "method.test.action",
		},
		{
			category: "",
			action:   "",
			expected: "method..",
		},
	}

	for _, tc := range testCases {
		result := DefaultMethodFormatter(tc.category, tc.action)
		if result != tc.expected {
			t.Errorf("DefaultMethodFormatter(%s, %s) = %s, want %s",
				tc.category, tc.action, result, tc.expected)
		}
	}
}