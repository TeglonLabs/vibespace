package streaming

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	// Create a mock client
	client := NewMockNATSClient()
	client.url = "nats://localhost:4222"
	
	// Test connection
	err := client.Connect()
	assert.NoError(t, err)
	assert.True(t, client.IsConnected())
	
	// Connecting when already connected should work
	err = client.Connect()
	assert.NoError(t, err)
}

func TestClose(t *testing.T) {
	// Create and connect a mock client
	client := NewMockNATSClient()
	client.url = "nats://localhost:4222"
	err := client.Connect()
	assert.NoError(t, err)
	
	// Close the connection
	client.Close()
	assert.False(t, client.IsConnected())
	
	// Closing an already closed connection should be a no-op
	client.Close()
	assert.False(t, client.IsConnected())
}

func TestPublishWorldMoment(t *testing.T) {
	// Create and connect a mock client
	client := NewMockNATSClient()
	client.url = "nats://localhost:4222"
	err := client.Connect()
	assert.NoError(t, err)
	
	// Create a test moment
	tempVal := 22.5
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		SensorData: models.SensorData{
			Temperature: &tempVal,
		},
		CreatorID: "user123",
		Sharing: models.SharingSettings{
			IsPublic: true,
		},
	}
	
	// Publish the moment
	err = client.PublishWorldMoment(moment, "user123")
	assert.NoError(t, err)
	
	// Verify it was published
	publishedMoments := client.GetPublishedMoments()
	assert.Equal(t, 1, len(publishedMoments))
	assert.Equal(t, "test-world", publishedMoments[0].WorldID)
	
	// Test publishing when disconnected
	client.Close()
	err = client.PublishWorldMoment(moment, "user123")
	assert.Error(t, err)
}

func TestPublishVibeUpdate(t *testing.T) {
	// Create and connect a mock client
	client := NewMockNATSClient()
	client.url = "nats://localhost:4222"
	err := client.Connect()
	assert.NoError(t, err)
	
	// Create a test vibe
	vibe := &models.Vibe{
		ID:      "vibe1",
		Name:    "Test Vibe",
		Energy:  0.8,
		Mood:    models.MoodEnergetic,
	}
	
	// Publish the vibe update
	err = client.PublishVibeUpdate("test-world", vibe)
	assert.NoError(t, err)
	
	// Verify it was published
	publishedVibes := client.GetPublishedVibes()
	assert.Equal(t, 1, len(publishedVibes))
	assert.Equal(t, "vibe1", publishedVibes["test-world"].ID)
	
	// Test publishing when disconnected
	client.Close()
	err = client.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
}