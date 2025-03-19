package streaming

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

// TestServicePublishVibeUpdate tests the PublishVibeUpdate method on the streaming service
func TestServicePublishVibeUpdate(t *testing.T) {
	// Create a mock client
	mockClient := NewMockNATSClient()
	mockClient.Connect()
	
	// Create a simple service
	service := &StreamingService{}
	service.SetClient(mockClient)
	
	// Create a test vibe
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Description: "For testing",
		Energy:      0.8,
		Mood:        "energetic",
	}
	
	// Test successful publish
	err := service.PublishVibeUpdate("test-world", vibe)
	assert.NoError(t, err)
	
	// Verify it was published
	publishedVibes := mockClient.GetPublishedVibes()
	assert.Contains(t, publishedVibes, "test-world")
	assert.Equal(t, vibe, publishedVibes["test-world"])
	
	// Test when disconnected
	mockClient.Close()
	
	// Should get an error
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to NATS")
	
	// Reset client and test with publish error
	mockClient.SetConnectError(nil)
	mockClient.Connect()
	mockClient.SetPublishVibeError(assert.AnError)
	
	// Should get an error
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
}