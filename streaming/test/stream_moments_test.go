package test

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/bmorphism/vibespace-mcp-go/streaming/test/mocks"
	"github.com/bmorphism/vibespace-mcp-go/streaming/testutils"
	"github.com/stretchr/testify/assert"
)

func TestStreamMomentsProcess2(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	// Create mock client
	mockClient := streaming.NewMockNATSClient()
	mockClient.Connect()
	
	// Create service with a very short stream interval
	config := &streaming.StreamingConfig{
		StreamInterval: 10 * time.Millisecond,
		StreamID:       "test-stream",
	}
	service := testutils.CreateMockStreamingService(repo, config, mockClient)
	
	// Start streaming
	err := service.StartStreaming()
	assert.NoError(t, err)
	
	// Wait for at least one streaming cycle
	time.Sleep(20 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
	
	// Check that moments were published
	publishedMoments := mockClient.GetPublishedMoments()
	assert.NotEmpty(t, publishedMoments)
	
	// Test error conditions in streamMoments
	// Create a mock generator that returns an error
	mockGenerator := testutils.NewEnhancedMockMomentGenerator(repo)
	mockGenerator.SetGenerateError(assert.AnError)
	
	// Create a new service with the mock generator
	service = testutils.CreateMockStreamingService(repo, config, mockClient)
	service.SetMomentGenerator(mockGenerator)
	
	// Start streaming (shouldn't crash despite error)
	err = service.StartStreaming()
	assert.NoError(t, err)
	
	// Wait briefly
	time.Sleep(20 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
	
	// Now test publish error scenario
	mockClient.SetPublishMomentError(assert.AnError)
	mockGenerator.SetGenerateError(nil)
	
	// Start streaming again
	err = service.StartStreaming()
	assert.NoError(t, err)
	
	// Wait briefly
	time.Sleep(20 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
}