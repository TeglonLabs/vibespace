package test

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/bmorphism/vibespace-mcp-go/streaming/test/mocks"
	"github.com/stretchr/testify/assert"
)

// MockMomentGenerator for service testing
type ExtendedMockGenerator struct {
	GenerateError       error
	GenerateAllError    error
	GeneratedMoments    []*models.WorldMoment
	GeneratedAllMoments [][]*models.WorldMoment
}

func NewExtendedMockGenerator() *ExtendedMockGenerator {
	return &ExtendedMockGenerator{
		GeneratedMoments:    []*models.WorldMoment{},
		GeneratedAllMoments: [][]*models.WorldMoment{},
	}
}

func (g *ExtendedMockGenerator) GenerateMoment(worldID string) (*models.WorldMoment, error) {
	if g.GenerateError != nil {
		return nil, g.GenerateError
	}
	
	moment := &models.WorldMoment{
		WorldID:   worldID,
		Timestamp: time.Now().Unix(),
	}
	g.GeneratedMoments = append(g.GeneratedMoments, moment)
	return moment, nil
}

func (g *ExtendedMockGenerator) GenerateAllMoments() ([]*models.WorldMoment, error) {
	if g.GenerateAllError != nil {
		return nil, g.GenerateAllError
	}
	
	moments := []*models.WorldMoment{
		{
			WorldID:   "world-1",
			Timestamp: time.Now().Unix(),
		},
		{
			WorldID:   "world-2",
			Timestamp: time.Now().Unix(),
		},
	}
	g.GeneratedAllMoments = append(g.GeneratedAllMoments, moments)
	return moments, nil
}

// TestStreamingServiceStartWithError tests the service Start method with error
func TestStreamingServiceStartWithError(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	config := &streaming.StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
		NATSUrl:        "nats://nonexistent:4222",
		StreamID:       "test-stream",
	}
	
	// Create a mock client that returns an error on Connect
	mockClient := streaming.NewMockNATSClient()
	mockClient.SetConnectError(assert.AnError)
	
	// Create service
	service := LocalCreateMockStreamingService(repo, config, mockClient)
	
	// Start should return the connect error
	err := service.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
}

// LocalCreateMockStreamingService creates a streaming service for testing
func LocalCreateMockStreamingService(repo *mocks.MockRepository, config *streaming.StreamingConfig, natsClient streaming.NATSClientInterface) *streaming.StreamingService {
	service := &streaming.StreamingService{}
	
	// Use the setter methods from the test_helper.go file
	service.SetConfig(config)
	service.SetClient(natsClient)
	service.SetMomentGenerator(NewExtendedMockGenerator())
	service.SetRepository(repo)
	
	return service
}

// TestStreamingServiceStartWithAutoStart tests auto-starting
func TestStreamingServiceStartWithAutoStart(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	config := &streaming.StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
		NATSUrl:        "nats://nonexistent:4222",
		StreamID:       "test-stream",
		AutoStart:      true,
	}
	
	// Create a mock client
	mockClient := streaming.NewMockNATSClient()
	
	// Create service
	service := LocalCreateMockStreamingService(repo, config, mockClient)
	
	// Start should also start streaming
	err := service.Start()
	assert.NoError(t, err)
	assert.True(t, service.IsStreaming())
	
	// Verify that the client was connected
	assert.True(t, mockClient.IsConnected())
	
	// Stop the service
	service.Stop()
	assert.False(t, service.IsStreaming())
	assert.False(t, mockClient.IsConnected())
}

// TestStreamMomentsProcess tests the streamMoments method
func TestStreamMomentsProcess(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	// Create mock client
	mockClient := streaming.NewMockNATSClient()
	mockClient.Connect()
	
	// Create a custom generator that we can track
	mockGenerator := NewExtendedMockGenerator()
	
	// Create service with a very short stream interval
	config := &streaming.StreamingConfig{
		StreamInterval: 10 * time.Millisecond,
		StreamID:       "test-stream",
	}
	service := LocalCreateMockStreamingService(repo, config, mockClient)
	service.SetMomentGenerator(mockGenerator)
	
	// Start streaming
	err := service.StartStreaming()
	assert.NoError(t, err)
	
	// Wait for at least two streaming cycles
	time.Sleep(25 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
	
	// Verify that moments were generated
	assert.GreaterOrEqual(t, len(mockGenerator.GeneratedAllMoments), 1, "Should have generated moments")
	
	// Verify that moments were published
	publishedMoments := mockClient.GetPublishedMoments()
	assert.NotEmpty(t, publishedMoments, "Should have published moments")
}

// TestStreamMomentsWithGenerateError tests handling of generator errors
func TestStreamMomentsWithGenerateError(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	// Create mock client
	mockClient := streaming.NewMockNATSClient()
	mockClient.Connect()
	
	// Create a custom generator that returns an error
	mockGenerator := NewExtendedMockGenerator()
	mockGenerator.GenerateAllError = assert.AnError
	
	// Create service with a very short stream interval
	config := &streaming.StreamingConfig{
		StreamInterval: 10 * time.Millisecond,
		StreamID:       "test-stream",
	}
	service := LocalCreateMockStreamingService(repo, config, mockClient)
	service.SetMomentGenerator(mockGenerator)
	
	// Start streaming (should not panic despite error)
	err := service.StartStreaming()
	assert.NoError(t, err)
	
	// Wait briefly
	time.Sleep(25 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
	
	// No moments should be published due to the error
	publishedMoments := mockClient.GetPublishedMoments()
	assert.Empty(t, publishedMoments, "Should not have published moments due to error")
}

// TestStreamMomentsWithPublishError tests handling of publish errors
func TestStreamMomentsWithPublishError(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	// Create mock client with publish error
	mockClient := streaming.NewMockNATSClient()
	mockClient.Connect()
	mockClient.SetPublishMomentError(assert.AnError)
	
	// Create a custom generator
	mockGenerator := NewExtendedMockGenerator()
	
	// Create service with a very short stream interval
	config := &streaming.StreamingConfig{
		StreamInterval: 10 * time.Millisecond,
		StreamID:       "test-stream",
	}
	service := LocalCreateMockStreamingService(repo, config, mockClient)
	service.SetMomentGenerator(mockGenerator)
	
	// Start streaming (should not panic despite error)
	err := service.StartStreaming()
	assert.NoError(t, err)
	
	// Wait briefly
	time.Sleep(25 * time.Millisecond)
	
	// Stop streaming
	service.StopStreaming()
	
	// Moments should be generated but not published
	assert.GreaterOrEqual(t, len(mockGenerator.GeneratedAllMoments), 1, "Should have generated moments")
	publishedMoments := mockClient.GetPublishedMoments()
	assert.Empty(t, publishedMoments, "Should not have published moments due to error")
}

// TestServicePublishVibeUpdate tests the PublishVibeUpdate method
func TestServicePublishVibeUpdate(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	// Create a mock client
	mockClient := streaming.NewMockNATSClient()
	mockClient.Connect()
	
	// Create service
	config := &streaming.StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
	}
	service := LocalCreateMockStreamingService(repo, config, mockClient)
	
	// Create a test vibe
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Description: "For testing",
		Energy:      0.8,
		Mood:        "energetic",
	}
	
	// Publish the vibe
	err := service.PublishVibeUpdate("test-world", vibe)
	assert.NoError(t, err)
	
	// Verify it was published
	publishedVibes := mockClient.GetPublishedVibes()
	assert.Contains(t, publishedVibes, "test-world")
	assert.Equal(t, vibe, publishedVibes["test-world"])
	
	// Test with connection error
	mockClient.Close()
	mockClient.SetConnectError(assert.AnError)
	
	// Should get an error
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
	
	// Reset client and test with publish error
	mockClient.SetConnectError(nil)
	mockClient.Connect()
	mockClient.SetPublishVibeError(assert.AnError)
	
	// Should get an error
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
}

// TestStreamSingleWorld tests the StreamSingleWorld method
func TestStreamSingleWorld(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	repo.AddTestData()
	
	// Create a mock client
	mockClient := streaming.NewMockNATSClient()
	mockClient.Connect()
	
	// Create service
	config := &streaming.StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
	}
	service := LocalCreateMockStreamingService(repo, config, mockClient)
	
	// Test streaming a single world
	err := service.StreamSingleWorld("test-world", "test-user")
	assert.NoError(t, err)
	
	// Verify it was published
	publishedMoments := mockClient.GetPublishedMoments()
	assert.NotEmpty(t, publishedMoments)
	assert.Equal(t, "test-world", publishedMoments[0].WorldID)
	assert.Equal(t, "test-user", publishedMoments[0].CreatorID)
	assert.Contains(t, publishedMoments[0].Viewers, "test-user")
	
	// Test error handling when not connected
	mockClient.Close()
	mockClient.SetConnectError(assert.AnError)
	
	err = service.StreamSingleWorld("test-world", "test-user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
	
	// Test error handling for moment generation
	mockClient.SetConnectError(nil)
	mockClient.Connect()
	
	mockGenerator := NewExtendedMockGenerator()
	mockGenerator.GenerateError = assert.AnError
	service.SetMomentGenerator(mockGenerator)
	
	err = service.StreamSingleWorld("test-world", "test-user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate moment")
	
	// Test error handling for publishing moment
	mockGenerator.GenerateError = nil
	mockClient.SetPublishMomentError(assert.AnError)
	
	err = service.StreamSingleWorld("test-world", "test-user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish moment")
}