package streaming

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

// TestStartStopStreaming tests the start and stop streaming functionality
func TestStartStopStreaming(t *testing.T) {
	// Create a mock client
	mockClient := NewMockNATSClient()
	mockClient.Connect()
	
	// Create a test config with a very short interval
	config := &StreamingConfig{
		StreamInterval: 10 * time.Millisecond,
	}
	
	// Create mock repository
	mockRepo := &MockRepository{
		Worlds: []models.World{
			{ID: "world-1", Name: "World 1"},
			{ID: "world-2", Name: "World 2"},
		},
	}
	
	// Create service
	service := &StreamingService{}
	service.SetConfig(config)
	service.SetClient(mockClient)
	service.SetRepository(mockRepo)
	service.SetMomentGenerator(NewMomentGenerator(mockRepo))
	
	// Test starting streaming
	err := service.StartStreaming()
	assert.NoError(t, err)
	assert.True(t, service.IsStreaming())
	
	// Wait for at least one streaming cycle
	time.Sleep(15 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
	assert.False(t, service.IsStreaming())
	
	// Verify that moments were published
	publishedMoments := mockClient.GetPublishedMoments()
	assert.NotEmpty(t, publishedMoments)
	
	// Test that starting again works
	err = service.StartStreaming()
	assert.NoError(t, err)
	assert.True(t, service.IsStreaming())
	
	// Wait for another streaming cycle
	time.Sleep(15 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
	assert.False(t, service.IsStreaming())
	
	// Test starting when already streaming (should be a no-op)
	service.StartStreaming()
	wasStreaming := service.IsStreaming()
	service.StartStreaming() // Call again
	stillStreaming := service.IsStreaming()
	assert.Equal(t, wasStreaming, stillStreaming)
	
	// Test stopping when not streaming (should be a no-op)
	service.StopStreaming()
	wasNotStreaming := !service.IsStreaming()
	service.StopStreaming() // Call again
	stillNotStreaming := !service.IsStreaming()
	assert.Equal(t, wasNotStreaming, stillNotStreaming)
}

// TestStreamMomentsWithErrors tests the streamMoments method with various error conditions
func TestStreamMomentsWithErrors(t *testing.T) {
	// Test 1: Generator error
	mockClient := NewMockNATSClient()
	mockClient.Connect()
	
	config := &StreamingConfig{
		StreamInterval: 10 * time.Millisecond,
	}
	
	// Create service with mock generator that returns errors
	service := &StreamingService{}
	service.SetConfig(config)
	service.SetClient(mockClient)
	mockGenerator := &MockMomentGenerator{
		GenerateAllError: assert.AnError,
	}
	service.SetMomentGenerator(mockGenerator)
	
	// Start streaming
	err := service.StartStreaming()
	assert.NoError(t, err)
	
	// Wait for at least one streaming cycle
	time.Sleep(15 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
	
	// No moments should be published due to the error
	publishedMoments := mockClient.GetPublishedMoments()
	assert.Empty(t, publishedMoments)
	
	// Test 2: Publish error
	mockClient = NewMockNATSClient()
	mockClient.Connect()
	mockClient.SetPublishMomentError(assert.AnError)
	
	// Create a new service
	service = &StreamingService{}
	service.SetConfig(config)
	service.SetClient(mockClient)
	
	// Use a working generator
	mockGenerator = &MockMomentGenerator{}
	service.SetMomentGenerator(mockGenerator)
	
	// Start streaming
	err = service.StartStreaming()
	assert.NoError(t, err)
	
	// Wait for at least one streaming cycle
	time.Sleep(15 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
	
	// No moments should be published due to the error
	publishedMoments = mockClient.GetPublishedMoments()
	assert.Empty(t, publishedMoments)
	
	// Test 3: Test stop channel works correctly
	mockClient = NewMockNATSClient()
	mockClient.Connect()
	
	// Create a new service with a deliberately long interval
	longConfig := &StreamingConfig{
		StreamInterval: 1 * time.Second, // Long interval
	}
	service = &StreamingService{}
	service.SetConfig(longConfig)
	service.SetClient(mockClient)
	service.SetMomentGenerator(mockGenerator)
	
	// Start streaming
	err = service.StartStreaming()
	assert.NoError(t, err)
	
	// Immediately stop streaming (before the ticker fires)
	service.StopStreaming()
	
	// Wait a moment to ensure any in-flight processing completes
	time.Sleep(10 * time.Millisecond)
	
	// No moments should be published since we stopped before the interval
	publishedMoments = mockClient.GetPublishedMoments()
	assert.Empty(t, publishedMoments)
}

// TestStartWithConnectionError tests starting the service with a connection error
func TestStartWithConnectionError(t *testing.T) {
	// Create a mock client that returns an error on connect
	mockClient := NewMockNATSClient()
	mockClient.SetConnectError(assert.AnError)
	
	// Create service
	config := &StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
	}
	service := &StreamingService{}
	service.SetConfig(config)
	service.SetClient(mockClient)
	
	// Start should return the connect error
	err := service.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
	
	// Service should not be streaming
	assert.False(t, service.IsStreaming())
}

// TestAutoStart tests the auto-start functionality
func TestAutoStart(t *testing.T) {
	// Create a mock client
	mockClient := NewMockNATSClient()
	
	// Create service with auto-start enabled
	config := &StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
		AutoStart:      true,
	}
	service := &StreamingService{}
	service.SetConfig(config)
	service.SetClient(mockClient)
	
	// Start should also start streaming
	err := service.Start()
	assert.NoError(t, err)
	assert.True(t, service.IsStreaming())
	
	// Clean up
	service.Stop()
	assert.False(t, service.IsStreaming())
}

// TestStartStreamingWithConnectionError tests starting streaming with a connection error
func TestStartStreamingWithConnectionError(t *testing.T) {
	// Create a mock client that's not connected
	mockClient := NewMockNATSClient()
	mockClient.SetConnectError(assert.AnError)
	
	// Create service
	config := &StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
	}
	service := &StreamingService{}
	service.SetConfig(config)
	service.SetClient(mockClient)
	
	// StartStreaming should try to connect and return the error
	err := service.StartStreaming()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
}

// TestDefaultConfig tests the default configuration values
func TestDefaultConfig(t *testing.T) {
	// Test with empty config
	emptyConfig := &StreamingConfig{}
	service := NewStreamingService(&MockRepository{}, emptyConfig)
	
	// Check default values
	config := service.GetConfig()
	assert.Equal(t, 4222, config.NATSPort)
	assert.Equal(t, "ies", config.StreamID)
	assert.Contains(t, config.NATSUrl, "nats://")
}