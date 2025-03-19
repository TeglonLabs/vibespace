package test

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/stretchr/testify/assert"
)

// TestNATSClientMethods tests the real NATS client methods 
// with proper mocking of the actual connection
func TestNATSClientMethods(t *testing.T) {
	// Create a NATS client with a fake URL that won't connect
	client := streaming.NewNATSClient("nats://nonexistent.example:4222")
	
	// Test Close method (should not panic when not connected)
	client.Close()
	
	// Test PublishWorldMoment when not connected
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: 123456789,
	}
	err := client.PublishWorldMoment(moment, "test-user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
	
	// Test PublishVibeUpdate when not connected
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Description: "A vibe for testing",
		Energy:      0.5,
		Mood:        "calm",
	}
	err = client.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
	
	// Test GetConnectionStatus when not connected
	status := client.GetConnectionStatus()
	assert.False(t, status.IsConnected)
	assert.Empty(t, status.ServerID)
	assert.Empty(t, status.ConnectedURL)
	assert.Empty(t, status.RTT)
}

// TestNATSErrorConditions tests error handling in NATS client
func TestNATSErrorConditions(t *testing.T) {
	// Create and use a mock client directly
	mockClient := streaming.NewMockNATSClient()
	mockClient.SetConnectError(assert.AnError)
	err := mockClient.Connect()
	assert.Error(t, err)
	
	// Test with simulated publishing errors
	
	// Continue with mock client tests
	mockClient.Connect()
	
	// Set a custom error that mentions rate limiting
	rateLimitError := assert.AnError
	mockClient.SetPublishMomentError(rateLimitError)
	
	// Try to publish (should fail with our error)
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: 123456789,
	}
	err = mockClient.PublishWorldMoment(moment, "test-user")
	assert.Error(t, err)
	assert.Equal(t, rateLimitError, err)
	
	// Set error for vibe publishing
	vibeError := assert.AnError
	mockClient.SetPublishVibeError(vibeError)
	
	// Try to publish vibe
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Energy:      0.5,
	}
	err = mockClient.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Equal(t, vibeError, err)
}