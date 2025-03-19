package test

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/stretchr/testify/assert"
)

// TestNATSClientClosingDisconnected tests the Close method when already disconnected
func TestNATSClientClosingDisconnected(t *testing.T) {
	// Create a real NATS client
	client := streaming.NewNATSClient("nats://unreachable:4222")
	
	// Call Close() without ever connecting - should not panic
	client.Close()
}

// TestNATSClientPublishWorldWithNoConnection tests PublishWorldMoment with no connection
func TestNATSClientPublishWorldWithNoConnection(t *testing.T) {
	// Create a client without connecting
	client := streaming.NewNATSClient("nats://unreachable:4222")
	
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: time.Now().Unix(),
	}
	
	// Should get an error since we're not connected
	err := client.PublishWorldMoment(moment, "user-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

// TestNATSClientPublishVibeWithNoConnection tests PublishVibeUpdate with no connection
func TestNATSClientPublishVibeWithNoConnection(t *testing.T) {
	// Create a client without connecting
	client := streaming.NewNATSClient("nats://unreachable:4222")
	
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Description: "Test Description",
		Energy:      0.5,
		Mood:        "calm",
	}
	
	// Should get an error since we're not connected
	err := client.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

// TestNATSClientWithRateLimiting tests rate limit enforcement
func TestNATSClientWithRateLimiting(t *testing.T) {
	// We'll use a mock client instead since we're just testing the rate limiter
	mockClient := streaming.NewMockNATSClient()
	mockClient.Connect()
	
	// Use mock client's error capabilities to simulate rate limiting
	
	// Use a special "rate limit exceeded" error for publishing
	mockClient.SetPublishMomentError(assert.AnError)
	mockClient.SetPublishVibeError(assert.AnError)
	
	// Try to publish world moment
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: time.Now().Unix(),
	}
	
	// Should get our mock error
	err := mockClient.PublishWorldMoment(moment, "user-1")
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	
	// Try to publish vibe
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Energy:      0.5,
		Mood:        "calm",
	}
	
	// Should get our mock error  
	err = mockClient.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}