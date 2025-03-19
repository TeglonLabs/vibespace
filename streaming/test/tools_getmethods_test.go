package test

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/stretchr/testify/assert"
)

// TestGetStreamingToolMethods tests full coverage of GetStreamingToolMethods
func TestGetStreamingToolMethodsFull(t *testing.T) {
	// Get the tool methods
	methods := streaming.GetStreamingToolMethods()
	
	// Check all methods are present
	assert.Len(t, methods, 5, "Should have 5 methods")
	
	// Check method names match expected
	expectedNames := []string{
		"streaming_startStreaming",
		"streaming_stopStreaming",
		"streaming_status",
		"streaming_streamWorld", 
		"streaming_updateConfig",
	}
	
	for _, name := range expectedNames {
		assert.Contains(t, methods, name, "Should contain method "+name)
		assert.NotNil(t, methods[name], "Method should not be nil")
	}
	
	// Test a specific method to ensure it's a proper function
	startStreamingMethod := methods["streaming_startStreaming"]
	assert.NotNil(t, startStreamingMethod)
}