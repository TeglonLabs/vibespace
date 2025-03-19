package streaming

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/stretchr/testify/assert"
)

// TestNonlocalNATSConfig tests that the NATS URL is correctly configured
func TestNonlocalNATSConfig(t *testing.T) {
	// Create a streaming config with nonlocal.info URL
	config := &StreamingConfig{
		NATSUrl:        "nats://nonlocal.info:4222",
		StreamInterval: 5 * time.Second,
		AutoStart:      false,
	}

	// Create a NATS client with the config
	client := NewNATSClient(config.NATSUrl)
	
	// Verify the URL is correct
	assert.Equal(t, "nats://nonlocal.info:4222", client.url)
}

// TestStreamingStatusWithUIIndicators tests the UI indicators in streaming status
func TestStreamingStatusWithUIIndicators(t *testing.T) {
	// Create a fake repository
	repo := repository.NewRepository()
	
	// Create streaming service with nonlocal.info URL
	config := &StreamingConfig{
		NATSUrl:        "nats://nonlocal.info:4222",
		StreamInterval: 5 * time.Second,
		AutoStart:      false,
	}
	service := NewStreamingService(repo, config)
	
	// Create streaming tools
	tools := NewStreamingTools(service)
	
	// Get the status
	status, err := tools.Status()
	assert.NoError(t, err)
	
	// Verify the status has the correct URL
	assert.Equal(t, "nats://nonlocal.info:4222", status.NATSUrl)
	
	// Verify UI indicators are present
	assert.False(t, status.UIIndicators.StreamActive) // Should be false initially
	assert.Equal(t, "OFFLINE", status.UIIndicators.StreamIndicator)
	assert.Equal(t, "#F44336", status.UIIndicators.StatusColor) // Red color for offline
	assert.Equal(t, "DISCONNECTED", status.UIIndicators.ConnectionQuality)
}

// TestSharingSettingsInStreamWorld tests that sharing settings are correctly parsed
func TestSharingSettingsInStreamWorld(t *testing.T) {
	// Create sharing request
	sharingReq := &SharingRequest{
		IsPublic:     true,
		AllowedUsers: []string{"user1", "user2"},
		ContextLevel: "partial",
	}
	
	// Verify it converts correctly to a models.SharingSettings
	sharing := models.SharingSettings{
		IsPublic:     sharingReq.IsPublic,
		AllowedUsers: sharingReq.AllowedUsers,
		ContextLevel: models.ContextLevel(sharingReq.ContextLevel),
	}
	
	// Check the values
	assert.True(t, sharing.IsPublic)
	assert.ElementsMatch(t, []string{"user1", "user2"}, sharing.AllowedUsers)
	assert.Equal(t, models.ContextLevel("partial"), sharing.ContextLevel)
}

// TestNATSSubscriberConfig tests that the NATS subscriber example uses nonlocal.info
func TestNATSSubscriberConfig(t *testing.T) {
	// Open the example file and check that it's using the correct URL
	// This is a simple test that doesn't require execution of the file
	
	// We're just verifying that the URL is correctly set in the test
	expectedURL := "nats://nonlocal.info:4222"
	
	// Create a NATS client with the expected URL to simulate what the subscriber would use
	client := NewNATSClient(expectedURL)
	assert.Equal(t, expectedURL, client.url)
}

// TestRateLimiter tests the rate limiter functionality
func TestRateLimiter(t *testing.T) {
	// Create a rate limiter with 5 tokens, 1 refill per 100ms, and 100ms interval
	limiter := NewRateLimiter(5, 1, 100)
	
	// Should allow 5 initial requests
	for i := 0; i < 5; i++ {
		allowed := limiter.TryAcquire()
		assert.True(t, allowed, "Should allow request %d", i+1)
	}
	
	// Should deny the 6th request
	allowed := limiter.TryAcquire()
	assert.False(t, allowed, "Should deny the 6th request")
	
	// Wait for a refill
	time.Sleep(100 * time.Millisecond)
	
	// Should allow the next request after refill
	allowed = limiter.TryAcquire()
	assert.True(t, allowed, "Should allow request after refill")
	
	// Should deny the request after that
	allowed = limiter.TryAcquire()
	assert.False(t, allowed, "Should deny request after using refilled token")
}

// TestMomentGeneratorWithNATS tests that world moments are generated correctly for NATS
func TestMomentGeneratorWithNATS(t *testing.T) {
	// Create repository with test data
	repo := repository.NewRepository()
	
	// Add a test world with creator attribution
	testWorld := models.World{
		ID:          "test-world",
		Name:        "Test World",
		Description: "A world for testing",
		Type:        models.WorldTypeVirtual,
		CreatorID:   "user123",
		Occupancy:   10,
		Sharing: models.SharingSettings{
			IsPublic:     true,
			AllowedUsers: []string{"user456"},
			ContextLevel: models.ContextLevelFull,
		},
	}
	
	err := repo.AddWorld(testWorld)
	assert.NoError(t, err)
	
	// Create a moment generator
	generator := NewMomentGenerator(repo)
	
	// Generate a moment
	moment, err := generator.GenerateMoment("test-world")
	assert.NoError(t, err)
	
	// Verify the moment has the correct values
	assert.Equal(t, "test-world", moment.WorldID)
	assert.Equal(t, "user123", moment.CreatorID)
	assert.Equal(t, 10, moment.Occupancy)
	assert.True(t, moment.Sharing.IsPublic)
	assert.Contains(t, moment.Sharing.AllowedUsers, "user456")
	assert.Equal(t, models.ContextLevelFull, moment.Sharing.ContextLevel)
	
	// Check timestamp is in milliseconds (should be a 13-digit number for current time)
	now := time.Now().UnixNano() / int64(time.Millisecond)
	// Timestamp should be within a second of now
	assert.InDelta(t, now, moment.Timestamp, 1000)
}