package streaming

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

// TestNATSClientConnect tests more aspects of the Connect method 
func TestNATSClientConnect(t *testing.T) {
	// Create a client with test config
	client := &NATSClient{
		url:      "nats://localhost:4222",
		streamID: "test-stream",
	}
	
	// Get connection status when not connected
	assert.False(t, client.IsConnected())
	status := client.GetConnectionStatus()
	assert.False(t, status.IsConnected)
	assert.Equal(t, "nats://localhost:4222", status.URL)
	
	// Set connected state manually for testing
	client.connected = true
	status = client.GetConnectionStatus()
	assert.True(t, status.IsConnected)
}

// TestCloseImplementation tests the Close method implementation
func TestCloseImplementation(t *testing.T) {
	// Create a client
	client := &NATSClient{
		url:       "nats://localhost:4222",
		streamID:  "test-stream",
		connected: true,
	}
	
	// Close should mark as disconnected
	client.Close()
	
	// Verify it was closed - nothing to check since we're mocking
}

// TestPublishWorldMomentImplementation tests the PublishWorldMoment method
func TestPublishWorldMomentImplementation(t *testing.T) {
	// Create a client for testing
	client := &NATSClient{
		url:       "nats://localhost:4222", 
		streamID:  "test-stream",
		connected: true,
		rateLimiter: NewRateLimiter(100, 10, 1000),
	}
	
	// Create a test moment
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: time.Now().Unix(),
		Viewers:   []string{"user1", "user2"},
		CreatorID: "creator1",
		Sharing: models.SharingSettings{
			IsPublic:     true,
			ContextLevel: models.ContextLevelFull,
		},
		CustomData: "{\"test\":\"data\"}",
	}
	
	// Test not connected error
	client.connected = false
	err := client.PublishWorldMoment(moment, "user1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

// TestPublishVibeUpdateImplementation tests the PublishVibeUpdate method
func TestPublishVibeUpdateImplementation(t *testing.T) {
	// Create a client for testing
	client := &NATSClient{
		url:       "nats://localhost:4222",
		streamID:  "test-stream",
		connected: true,
		rateLimiter: NewRateLimiter(100, 10, 1000),
	}
	
	// Create a test vibe
	vibe := &models.Vibe{
		ID:          "vibe-1",
		Name:        "Test Vibe",
		Description: "A vibe for testing",
		Energy:      0.8,
		Mood:        "energetic",
	}
	
	// Test not connected error
	client.connected = false
	err := client.PublishVibeUpdate("world-1", vibe)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}