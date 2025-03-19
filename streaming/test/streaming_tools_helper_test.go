package test

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/bmorphism/vibespace-mcp-go/streaming/test/mocks"
	"github.com/stretchr/testify/assert"
)

// TestTestHelperFunctions tests the streaming test helper functions
func TestTestHelperFunctions(t *testing.T) {
	// Create a NATS client
	client := streaming.NewNATSClient("test-url")
	
	// Create a rate limiter
	rateLimiter := streaming.NewRateLimiter(5, 1, 100)
	
	// Test SetRateLimiter method
	client.SetRateLimiter(rateLimiter)
	
	// Create a streaming service
	service := &streaming.StreamingService{}
	
	// Create a config
	config := &streaming.StreamingConfig{
		NATSUrl: "test-url",
		StreamID: "test-stream",
	}
	
	// Test SetConfig method
	service.SetConfig(config)
	assert.Equal(t, config, service.GetConfig())
	
	// Test SetClient method
	service.SetClient(client)
	assert.Equal(t, client, service.GetNATSClient())
	
	// Test SetStreamingActive method
	service.SetStreamingActive(true)
	assert.True(t, service.IsStreaming())
	
	// Test SetMomentGenerator and SetRepository with mock objects
	mockRepo := mocks.NewMockRepository()
	momentGenerator := streaming.NewMomentGenerator(nil)
	
	// These will use reflection
	service.SetMomentGenerator(momentGenerator)
	service.SetRepository(mockRepo)
}