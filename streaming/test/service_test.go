package test

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/bmorphism/vibespace-mcp-go/streaming/test/mocks"
	"github.com/bmorphism/vibespace-mcp-go/streaming/testutils"
	"github.com/stretchr/testify/assert"
)

func setupTestService() (*streaming.StreamingService, *streaming.MockNATSClient, *mocks.MockRepository) {
	// Create repository
	repo := mocks.NewMockRepository()
	
	// Create a test world
	world := &models.World{
		ID:          "test-world",
		Name:        "Test World",
		Description: "A world for testing",
		CreatorID:   "test-user",
	}
	repo.SaveWorld(world)
	
	// Create config
	config := &streaming.StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
		StreamID:       "test-stream",
		AutoStart:      false,
	}
	
	// Create mock NATS client
	mockClient := streaming.NewMockNATSClient()
	
	// Create service with mock client
	service := testutils.CreateMockStreamingService(repo, config, mockClient)
	
	return service, mockClient, repo
}

func TestStreamingServiceStart(t *testing.T) {
	service, mockClient, _ := setupTestService()
	
	// Test start without auto-start
	err := service.Start()
	assert.NoError(t, err)
	assert.True(t, mockClient.IsConnected(), "Client should be connected after Start()")
	assert.False(t, service.IsStreaming(), "Service should not be streaming with autoStart=false")
	
	// Test start with connection error
	mockClient.Close()
	mockClient.SetConnectError(assert.AnError)
	err = service.Start()
	assert.Error(t, err)
	
	// Test with auto-start
	mockClient.SetConnectError(nil)
	service, mockClient, _ = setupTestService()
	// Set auto start (modify the config directly)
	service.SetConfig(&streaming.StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
		StreamID:       "test-stream",
		AutoStart:      true,
	})
	
	err = service.Start()
	assert.NoError(t, err)
	assert.True(t, mockClient.IsConnected(), "Client should be connected after Start()")
	assert.True(t, service.IsStreaming(), "Service should be streaming with autoStart=true")
	
	// Clean up
	service.Stop()
}

func TestStreamingServiceStartStopStreaming(t *testing.T) {
	service, mockClient, _ := setupTestService()
	
	// Start service first to connect
	err := service.Start()
	assert.NoError(t, err)
	
	// Test start streaming
	err = service.StartStreaming()
	assert.NoError(t, err)
	assert.True(t, service.IsStreaming())
	
	// Starting again should be a no-op
	err = service.StartStreaming()
	assert.NoError(t, err)
	
	// Test stop streaming
	service.StopStreaming()
	assert.False(t, service.IsStreaming())
	
	// Stopping again should be a no-op
	service.StopStreaming()
	assert.False(t, service.IsStreaming())
	
	// Test start with connection error
	mockClient.Close()
	mockClient.SetConnectError(assert.AnError)
	err = service.StartStreaming()
	assert.Error(t, err)
	
	// Clean up
	service.Stop()
}

func TestStreamingServiceStop(t *testing.T) {
	service, mockClient, _ := setupTestService()
	
	// Start service and streaming
	service.Start()
	service.StartStreaming()
	assert.True(t, service.IsStreaming())
	assert.True(t, mockClient.IsConnected())
	
	// Stop service
	service.Stop()
	assert.False(t, service.IsStreaming())
	assert.False(t, mockClient.IsConnected())
}

func TestStreamSingleWorldBasic(t *testing.T) {
	service, mockClient, _ := setupTestService()
	
	// Start service
	service.Start()
	
	// Stream a single world
	err := service.StreamSingleWorld("test-world", "test-user")
	assert.NoError(t, err)
	
	// Check if a moment was published
	publishedMoments := mockClient.GetPublishedMoments()
	assert.Len(t, publishedMoments, 1)
	assert.Equal(t, "test-world", publishedMoments[0].WorldID)
	assert.Equal(t, "test-user", publishedMoments[0].CreatorID)
	
	// Test with connection error
	mockClient.SimulateDisconnect()
	mockClient.SetConnectError(assert.AnError)
	err = service.StreamSingleWorld("test-world", "test-user")
	assert.Error(t, err)
	
	// Test with publish error
	mockClient.SetConnectError(nil)
	mockClient.Connect()
	mockClient.SetPublishMomentError(assert.AnError)
	err = service.StreamSingleWorld("test-world", "test-user")
	assert.Error(t, err)
	
	// Test with invalid world ID
	mockClient.SetPublishMomentError(nil)
	err = service.StreamSingleWorld("non-existent-world", "test-user")
	assert.Error(t, err)
	
	// Clean up
	service.Stop()
}

func TestPublishVibeUpdate(t *testing.T) {
	service, mockClient, _ := setupTestService()
	
	// Start service
	service.Start()
	
	// Create a vibe
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "happiness",
		Description: "Test vibe",
		Energy:      0.8,
		Mood:        "happy",
	}
	
	// Publish vibe update
	err := service.PublishVibeUpdate("test-world", vibe)
	assert.NoError(t, err)
	
	// Check if vibe was published
	publishedVibes := mockClient.GetPublishedVibes()
	assert.Contains(t, publishedVibes, "test-world")
	assert.Equal(t, vibe, publishedVibes["test-world"])
	
	// Test with connection error
	mockClient.Close()
	mockClient.SetConnectError(assert.AnError)
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	
	// Test with publish error
	mockClient.SetConnectError(nil)
	mockClient.Connect()
	mockClient.SetPublishVibeError(assert.AnError)
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	
	// Clean up
	service.Stop()
}