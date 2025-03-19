package test

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/stretchr/testify/assert"
)

func TestGetStreamingToolMethods(t *testing.T) {
	// Get the tool methods
	methods := streaming.GetStreamingToolMethods()
	
	// Check that all expected methods exist
	expectedMethods := []string{
		"streaming_startStreaming",
		"streaming_stopStreaming",
		"streaming_status",
		"streaming_streamWorld",
		"streaming_updateConfig",
	}
	
	for _, method := range expectedMethods {
		assert.Contains(t, methods, method)
		assert.NotNil(t, methods[method])
	}
	
	// Verify the methods exist (don't need to check specific implementation types)
	assert.NotNil(t, methods["streaming_startStreaming"])
	assert.NotNil(t, methods["streaming_stopStreaming"])
	assert.NotNil(t, methods["streaming_status"])
	assert.NotNil(t, methods["streaming_streamWorld"])
	assert.NotNil(t, methods["streaming_updateConfig"])
}