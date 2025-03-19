package streaming

import (
	"errors"
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/stretchr/testify/assert"
)

func TestStreamingServiceWithMock(t *testing.T) {
	// Create test repository
	repo := repository.NewRepository()
	
	// Add a test world
	world := models.World{
		ID:          "test-world",
		Name:        "Test World",
		Description: "A world for testing",
		Type:        models.WorldTypeVirtual,
		CreatorID:   "user123",
		Sharing: models.SharingSettings{
			IsPublic:     true,
			ContextLevel: models.ContextLevelFull,
		},
	}
	err := repo.AddWorld(world)
	assert.NoError(t, err)
	
	// Add a test vibe
	vibe := models.Vibe{
		ID:     "test-vibe",
		Name:   "Test Vibe",
		Energy: 0.8,
		Mood:   models.MoodCalm,
	}
	err = repo.AddVibe(vibe)
	assert.NoError(t, err)
	
	// Associate vibe with world
	err = repo.SetWorldVibe("test-world", "test-vibe")
	assert.NoError(t, err)
	
	// Create config
	config := &StreamingConfig{
		NATSUrl:        "nats://test:4222",
		StreamInterval: 100 * time.Millisecond,
		StreamID:       "test-stream",
		AutoStart:      false,
	}
	
	t.Run("Start connects to NATS", func(t *testing.T) {
		// Create mock NATS client
		mockNats := NewMockNATSClient()
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Call Start
		err := service.Start()
		assert.NoError(t, err)
		
		// Verify Connect was called
		assert.True(t, mockNats.IsConnected())
		
		// AutoStart is false, so streaming should not be active
		assert.False(t, service.IsStreaming())
	})
	
	t.Run("Start with AutoStart", func(t *testing.T) {
		// Create mock NATS client
		mockNats := NewMockNATSClient()
		
		// Create config with AutoStart true
		autoStartConfig := &StreamingConfig{
			NATSUrl:        "nats://test:4222",
			StreamInterval: 100 * time.Millisecond,
			StreamID:       "test-stream",
			AutoStart:      true,
		}
		
		// Create service with mock
		service := CreateStreamingService(repo, autoStartConfig, mockNats)
		
		// Call Start
		err := service.Start()
		assert.NoError(t, err)
		
		// Verify Connect was called
		assert.True(t, mockNats.IsConnected())
		
		// AutoStart is true, so streaming should be active
		assert.True(t, service.IsStreaming())
		
		// Cleanup
		service.Stop()
	})
	
	t.Run("Start with connect error", func(t *testing.T) {
		// Create mock NATS client with error
		mockNats := NewMockNATSClient()
		mockNats.SetConnectError(errors.New("connect error"))
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Call Start
		err := service.Start()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connect error")
		
		// Streaming should not be active
		assert.False(t, service.IsStreaming())
	})
	
	t.Run("Stop ends streaming", func(t *testing.T) {
		// Create mock NATS client
		mockNats := NewMockNATSClient()
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Set streaming active
		service.streamingActive = true
		service.stopChan = make(chan struct{})
		
		// Call Stop
		service.Stop()
		
		// Verify streaming is stopped
		assert.False(t, service.IsStreaming())
		
		// Verify NATS connection is closed
		assert.False(t, mockNats.IsConnected())
	})
	
	t.Run("StartStreaming connects and starts streaming", func(t *testing.T) {
		// Create mock NATS client
		mockNats := NewMockNATSClient()
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Call StartStreaming
		err := service.StartStreaming()
		assert.NoError(t, err)
		
		// Verify streaming is active
		assert.True(t, service.IsStreaming())
		
		// Verify connected to NATS
		assert.True(t, mockNats.IsConnected())
		
		// Calling again should be a no-op
		err = service.StartStreaming()
		assert.NoError(t, err)
		
		// Cleanup
		service.Stop()
	})
	
	t.Run("StartStreaming with connect error", func(t *testing.T) {
		// Create mock NATS client with error
		mockNats := NewMockNATSClient()
		mockNats.SetConnectError(errors.New("connect error"))
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Call StartStreaming
		err := service.StartStreaming()
		assert.Error(t, err)
		
		// Streaming should not be active
		assert.False(t, service.IsStreaming())
	})
	
	t.Run("StopStreaming stops streaming", func(t *testing.T) {
		// Create mock NATS client
		mockNats := NewMockNATSClient()
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Set streaming active
		service.streamingActive = true
		service.stopChan = make(chan struct{})
		
		// Call StopStreaming
		service.StopStreaming()
		
		// Verify streaming is stopped
		assert.False(t, service.IsStreaming())
	})
	
	t.Run("StreamSingleWorld publishes a moment", func(t *testing.T) {
		// Create mock NATS client
		mockNats := NewMockNATSClient()
		mockNats.Connect()
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Call StreamSingleWorld
		err := service.StreamSingleWorld("test-world", "user123")
		assert.NoError(t, err)
		
		// Verify moment was published
		moments := mockNats.GetPublishedMoments()
		assert.Equal(t, 1, len(moments))
		assert.Equal(t, "test-world", moments[0].WorldID)
		assert.Equal(t, "test-vibe", moments[0].VibeID)
	})
	
	t.Run("StreamSingleWorld with non-existent world", func(t *testing.T) {
		// Create mock NATS client
		mockNats := NewMockNATSClient()
		mockNats.Connect()
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Call StreamSingleWorld with non-existent world
		err := service.StreamSingleWorld("non-existent", "user123")
		assert.Error(t, err)
		
		// Verify no moments were published
		moments := mockNats.GetPublishedMoments()
		assert.Equal(t, 0, len(moments))
	})
	
	t.Run("StreamSingleWorld with publish error", func(t *testing.T) {
		// Create mock NATS client with error
		mockNats := NewMockNATSClient()
		mockNats.Connect()
		mockNats.SetPublishMomentError(errors.New("publish error"))
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Call StreamSingleWorld
		err := service.StreamSingleWorld("test-world", "user123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "publish error")
	})
	
	t.Run("PublishVibeUpdate publishes a vibe update", func(t *testing.T) {
		// Create mock NATS client
		mockNats := NewMockNATSClient()
		mockNats.Connect()
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Get the vibe
		worldVibe, err := repo.GetWorldVibe("test-world")
		assert.NoError(t, err)
		
		// Call PublishVibeUpdate
		err = service.PublishVibeUpdate("test-world", &worldVibe)
		assert.NoError(t, err)
		
		// Verify vibe was published
		vibes := mockNats.GetPublishedVibes()
		assert.Equal(t, 1, len(vibes))
		assert.Equal(t, "test-vibe", vibes["test-world"].ID)
	})
	
	t.Run("PublishVibeUpdate with publish error", func(t *testing.T) {
		// Create mock NATS client with error
		mockNats := NewMockNATSClient()
		mockNats.Connect()
		mockNats.SetPublishVibeError(errors.New("publish error"))
		
		// Create service with mock
		service := CreateStreamingService(repo, config, mockNats)
		
		// Get the vibe
		worldVibe, err := repo.GetWorldVibe("test-world")
		assert.NoError(t, err)
		
		// Call PublishVibeUpdate
		err = service.PublishVibeUpdate("test-world", &worldVibe)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "publish error")
	})
}