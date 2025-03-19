package streaming

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/stretchr/testify/assert"
)

// This test suite focuses on the StreamingTools component using mock NATS client

func setupToolsTest(t *testing.T) (*repository.Repository, *MockNATSClient, *StreamingTools) {
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
	
	// Create mock NATS client
	mockNats := NewMockNATSClient()
	
	// Create config
	config := &StreamingConfig{
		NATSUrl:        "nats://test:4222",
		StreamInterval: 100 * time.Millisecond,
		StreamID:       "test-stream",
		AutoStart:      false,
	}
	
	// Create service
	service := CreateStreamingService(repo, config, mockNats)
	
	// Create tools
	tools := NewStreamingTools(service)
	
	return repo, mockNats, tools
}

func TestToolsStartStreaming(t *testing.T) {
	_, mockNats, tools := setupToolsTest(t)
	
	// Call StartStreaming
	req := &StartStreamingRequest{
		Interval: 200, // 200ms
	}
	
	resp, err := tools.StartStreaming(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.Message, "Streaming started with interval")
	
	// Verify streaming is active
	assert.True(t, tools.service.IsStreaming())
	
	// Verify connected to NATS
	assert.True(t, mockNats.IsConnected())
	
	// Verify interval was updated
	assert.Equal(t, 200*time.Millisecond, tools.service.config.StreamInterval)
	
	// Call StartStreaming again
	resp, err = tools.StartStreaming(req)
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "Streaming started with interval")
}

func TestToolsStopStreaming(t *testing.T) {
	_, _, tools := setupToolsTest(t)
	
	// Set streaming active
	tools.service.streamingActive = true
	tools.service.stopChan = make(chan struct{})
	
	// Call StopStreaming
	resp, err := tools.StopStreaming()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Streaming stopped", resp.Message)
	
	// Verify streaming is stopped
	assert.False(t, tools.service.IsStreaming())
	
	// Call StopStreaming again
	resp, err = tools.StopStreaming()
	assert.NoError(t, err)
	assert.Equal(t, "Streaming stopped", resp.Message)
}

func TestToolsStatus(t *testing.T) {
	_, mockNats, tools := setupToolsTest(t)
	
	// Get status when not streaming
	status, err := tools.Status()
	assert.NoError(t, err)
	assert.Equal(t, "nats://test:4222", status.NATSUrl)
	assert.False(t, status.IsStreaming)
	assert.False(t, status.UIIndicators.StreamActive)
	assert.Equal(t, "OFFLINE", status.UIIndicators.StreamIndicator)
	assert.Equal(t, "#F44336", status.UIIndicators.StatusColor) // Red for offline
	
	// Set streaming active
	tools.service.streamingActive = true
	mockNats.Connect()
	
	// Get status when streaming
	status, err = tools.Status()
	assert.NoError(t, err)
	assert.True(t, status.IsStreaming)
	assert.True(t, status.UIIndicators.StreamActive)
	assert.Equal(t, "ACTIVE", status.UIIndicators.StreamIndicator)
	assert.Equal(t, "#4CAF50", status.UIIndicators.StatusColor) // Green for online
	
	// Test with NATS disconnected
	mockNats.Close()
	
	status, err = tools.Status()
	assert.NoError(t, err)
	assert.True(t, status.IsStreaming) // Service still thinks it's streaming
	assert.Equal(t, "ACTIVE", status.UIIndicators.StreamIndicator)
	assert.Equal(t, "#4CAF50", status.UIIndicators.StatusColor) // Green for active
}

func TestToolsStreamWorld(t *testing.T) {
	repo, mockNats, tools := setupToolsTest(t)
	
	// Add a test vibe
	vibe := models.Vibe{
		ID:     "test-vibe",
		Name:   "Test Vibe",
		Energy: 0.8,
		Mood:   models.MoodCalm,
	}
	err := repo.AddVibe(vibe)
	assert.NoError(t, err)
	
	// Associate vibe with world
	err = repo.SetWorldVibe("test-world", "test-vibe")
	assert.NoError(t, err)
	
	// Connect mock NATS
	mockNats.Connect()
	
	// Create stream request
	req := &StreamWorldRequest{
		WorldID: "test-world",
		UserID:  "user123",
	}
	
	// Call StreamWorld
	resp, err := tools.StreamWorld(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.Message, "Streamed moment for world")
	
	// Verify moment was published
	moments := mockNats.GetPublishedMoments()
	assert.Equal(t, 1, len(moments))
	assert.Equal(t, "test-world", moments[0].WorldID)
	assert.Equal(t, "test-vibe", moments[0].VibeID)
	
	// Test with non-existent world
	req.WorldID = "non-existent"
	resp, err = tools.StreamWorld(req)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
}

func TestToolsUpdateConfig(t *testing.T) {
	_, mockNats, tools := setupToolsTest(t)
	_ = mockNats // Unused but retained for consistency with other tests
	
	// Initial config values
	assert.Equal(t, "nats://test:4222", tools.service.config.NATSUrl)
	assert.Equal(t, 100*time.Millisecond, tools.service.config.StreamInterval)
	assert.False(t, tools.service.config.AutoStart)
	
	// Create config update request
	req := &UpdateConfigRequest{
		NATSUrl:        "nats://new-server:5222",
		StreamInterval: 500, // 500ms
	}
	
	// Call UpdateConfig
	resp, err := tools.UpdateConfig(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.Message, "Configuration updated successfully")
	
	// Verify config was updated
	assert.Equal(t, "nats://new-server:5222", tools.service.config.NATSUrl)
	assert.Equal(t, 500*time.Millisecond, tools.service.config.StreamInterval)
	// AutoStart is not set in this test
}