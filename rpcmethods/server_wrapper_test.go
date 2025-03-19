package rpcmethods

import (
	"testing"
)

func TestIsValidMethod(t *testing.T) {
	// Test with valid methods
	if !IsValidMethod(MethodResourceRead) {
		t.Errorf("IsValidMethod(%s) = false, want true", MethodResourceRead)
	}
	
	if !IsValidMethod(MethodToolCall) {
		t.Errorf("IsValidMethod(%s) = false, want true", MethodToolCall)
	}
	
	// Test with invalid methods
	if IsValidMethod("invalid.method") {
		t.Errorf("IsValidMethod(invalid.method) = true, want false")
	}
	
	if IsValidMethod("") {
		t.Errorf("IsValidMethod(empty string) = true, want false")
	}
}

func TestNormalizeMethodName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Already standard resource read",
			input:    MethodResourceRead,
			expected: MethodResourceRead,
		},
		{
			name:     "Already standard tool call",
			input:    MethodToolCall,
			expected: MethodToolCall,
		},
		{
			name:     "Replace slashes with dots",
			input:    "resource/read",
			expected: MethodResourceRead,
		},
		{
			name:     "Resource read variant",
			input:    "resource.read",
			expected: MethodResourceRead,
		},
		{
			name:     "Tool call variant",
			input:    "tool.call",
			expected: MethodToolCall,
		},
		{
			name:     "PascalCase variant",
			input:    "Resources.Read",
			expected: MethodResourceRead,
		},
		{
			name:     "MCP prefixed",
			input:    "mcp.resource.read",
			expected: MethodResourceRead,
		},
		{
			name:     "Snake case variant",
			input:    "read_resource",
			expected: MethodResourceRead,
		},
		{
			name:     "Different prefix",
			input:    "get_resource",
			expected: MethodResourceRead,
		},
		{
			name:     "Unknown method",
			input:    "custom.method",
			expected: "custom.method",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := normalizeMethodName(tc.input)
			if result != tc.expected {
				t.Errorf("normalizeMethodName(%s) = %s, want %s", tc.input, result, tc.expected)
			}
		})
	}
}