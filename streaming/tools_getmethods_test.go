package streaming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetStreamingToolMethods tests the GetStreamingToolMethods function
func TestGetStreamingToolMethods(t *testing.T) {
	// Get the streaming tool methods
	methods := GetStreamingToolMethods()
	
	// Verify the correct methods are included
	assert.Len(t, methods, 5, "Should return 5 methods")
	
	// Check each expected method is present by key name
	expectedPrefixes := []string{
		"streaming_startStreaming",
		"streaming_stopStreaming",
		"streaming_status",
		"streaming_streamWorld",
		"streaming_updateConfig",
	}
	
	for _, prefix := range expectedPrefixes {
		_, exists := methods[prefix]
		assert.True(t, exists, "Method %s should be included", prefix)
	}
	
	// Check all methods are functions
	for key, method := range methods {
		assert.NotNil(t, method, "Method should not be nil for "+key)
	}
}