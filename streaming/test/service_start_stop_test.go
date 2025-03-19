package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/bmorphism/vibespace-mcp-go/streaming/test/mocks"
	"github.com/bmorphism/vibespace-mcp-go/streaming/testutils"
	"github.com/stretchr/testify/assert"
)

func TestStreamingServiceStartStop(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	config := &streaming.StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
		NATSUrl:        "nats://nonexistent:4222",
		StreamID:       "test-stream",
		AutoStart:      false,
	}
	
	// Mock client is needed to avoid actual connection attempts
	mockClient := streaming.NewMockNATSClient()
	
	// Create service with mock
	service := testutils.CreateMockStreamingService(repo, config, mockClient)
	
	// Test Start method
	err := service.Start()
	assert.NoError(t, err)
	assert.True(t, mockClient.IsConnected())
	
	// Test Stop method
	service.Stop()
	assert.False(t, mockClient.IsConnected())
	
	// Test with error during connect
	mockClient.SetConnectError(assert.AnError)
	err = service.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
	
	// Test auto-start behavior
	mockClient.SetConnectError(nil)
	config.AutoStart = true
	service = testutils.CreateMockStreamingService(repo, config, mockClient)
	
	// Start should also start streaming
	err = service.Start()
	assert.NoError(t, err)
	assert.True(t, service.IsStreaming())
	
	// Clean up
	service.Stop()
}

func TestServicePublishVibeUpdateStartStop(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	// Create mock client
	mockClient := streaming.NewMockNATSClient()
	
	// Create service
	config := &streaming.StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
	}
	service := testutils.CreateMockStreamingService(repo, config, mockClient)
	
	// Connect mock client
	mockClient.Connect()
	
	// Create a test vibe
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Description: "A vibe for testing",
		Energy:      0.5,
		Mood:        "calm",
	}
	
	// Test publishing when connected
	err := service.PublishVibeUpdate("test-world", vibe)
	assert.NoError(t, err)
	
	// Verify the vibe was published
	publishedVibes := mockClient.GetPublishedVibes()
	assert.Contains(t, publishedVibes, "test-world")
	assert.Equal(t, vibe, publishedVibes["test-world"])
	
	// Test when not connected
	mockClient.SimulateDisconnect()
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err, "Should return error when not connected")
	
	// Test with connect error
	mockClient.SetConnectError(fmt.Errorf("simulated connection error"))
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err, "Should return error with connection error")
	
	// Test with publish error
	mockClient.SetConnectError(nil)
	mockClient.Connect()
	mockClient.SetPublishVibeError(assert.AnError)
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
}