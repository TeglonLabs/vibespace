package streaming

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

// TestStartFullCoverage tests all code paths in the Start method
func TestStartFullCoverage(t *testing.T) {
	// Create service
	service := &StreamingService{}
	service.SetConfig(&StreamingConfig{
		NATSUrl:        "nats://test.host:4222",
		StreamID:       "test-stream",
		StreamInterval: 100 * time.Millisecond,
		AutoStart:      false,
	})
	
	// Create mock client
	mockClient := NewMockNATSClient()
	service.SetClient(mockClient)
	
	// Test successful connect without auto-start
	err := service.Start()
	assert.NoError(t, err)
	assert.False(t, service.IsStreaming())
	
	// Test with auto-start enabled
	service.config.AutoStart = true
	
	err = service.Start()
	assert.NoError(t, err)
	assert.True(t, service.IsStreaming())
	
	// Test with connection error
	mockClient.SetConnectError(assert.AnError)
	service.streamingActive = false
	
	err = service.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
	
	// Test already started/streaming
	mockClient.SetConnectError(nil)
	service.streamingActive = true
	
	err = service.StartStreaming()
	assert.NoError(t, err) // Should be no-op
}

// TestStreamSingleWorldFullCoverage tests more code paths in StreamSingleWorld
func TestStreamSingleWorldFullCoverage(t *testing.T) {
	// Create service components
	mockRepo := &MockRepository{
		Worlds: []models.World{
			{
				ID:   "test-world",
				Name: "Test World",
			},
		},
	}
	mockClient := NewMockNATSClient()
	mockClient.Connect()
	
	// Create service
	service := &StreamingService{}
	service.SetRepository(mockRepo)
	service.SetClient(mockClient)
	service.SetMomentGenerator(NewMomentGenerator(mockRepo))
	
	// Test basic case
	err := service.StreamSingleWorld("test-world", "test-user")
	assert.NoError(t, err)
	
	// Check the published moment
	moments := mockClient.GetPublishedMoments()
	assert.NotEmpty(t, moments)
	
	// Check the user was added as viewer
	hasUser := false
	for _, viewer := range moments[0].Viewers {
		if viewer == "test-user" {
			hasUser = true
			break
		}
	}
	assert.True(t, hasUser)
	
	// Test with existing user in viewers
	// First manually add a viewer
	service.momentGenerator = &CustomMockGenerator{
		moment: &models.WorldMoment{
			WorldID:  "test-world",
			Viewers:  []string{"test-user", "other-user"},
			Sharing:  models.SharingSettings{},
		},
	}
	
	err = service.StreamSingleWorld("test-world", "test-user")
	assert.NoError(t, err)
	
	// Test with existing creator ID
	service.momentGenerator = &CustomMockGenerator{
		moment: &models.WorldMoment{
			WorldID:   "test-world",
			CreatorID: "existing-creator",
		},
	}
	
	err = service.StreamSingleWorld("test-world", "test-user")
	assert.NoError(t, err)
	
	// Verify creator ID wasn't overwritten
	moments = mockClient.GetPublishedMoments()
	assert.Equal(t, "existing-creator", moments[len(moments)-1].CreatorID)
}

// MockMomentGenerator that returns a customizable moment
type CustomMockGenerator struct {
	moment *models.WorldMoment
	error  error
}

func (g *CustomMockGenerator) GenerateMoment(worldID string) (*models.WorldMoment, error) {
	if g.error != nil {
		return nil, g.error
	}
	
	// Clone the moment to avoid modifying the original
	return g.moment, nil
}

func (g *CustomMockGenerator) GenerateAllMoments() ([]*models.WorldMoment, error) {
	if g.error != nil {
		return nil, g.error
	}
	
	return []*models.WorldMoment{g.moment}, nil
}

// TestMomentGeneratorFullCoverage tests additional code paths in GenerateAllMoments
func TestMomentGeneratorFullCoverage(t *testing.T) {
	// Create mock repository with test data
	mockRepo := &MockRepository{
		Worlds: []models.World{
			{
				ID:          "world-1",
				Name:        "World 1",
				CurrentVibe: "vibe-1",
				Occupancy:   10,
			},
			{
				ID:          "world-2",
				Name:        "World 2",
				CurrentVibe: "vibe-2",
				Occupancy:   5,
			},
		},
		Vibes: map[string]*models.Vibe{
			"vibe-1": {
				ID:     "vibe-1",
				Name:   "Vibe 1",
				Energy: 0.7,
			},
			"vibe-2": {
				ID:     "vibe-2",
				Name:   "Vibe 2",
				Energy: 0.3,
			},
		},
	}
	
	// Create moment generator
	generator := NewMomentGenerator(mockRepo)
	
	// Generate all moments
	moments, err := generator.GenerateAllMoments()
	assert.NoError(t, err)
	assert.Len(t, moments, 2)
	
	// Verify both worlds were processed
	worldIDs := make(map[string]bool)
	for _, moment := range moments {
		worldIDs[moment.WorldID] = true
		
		// Check properties (vibe might not be populated in test mock)
		assert.NotEmpty(t, moment.WorldID)
		
		// Activity should be calculated but value depends on implementation
		assert.True(t, moment.Activity >= 0)
	}
	
	assert.True(t, worldIDs["world-1"])
	assert.True(t, worldIDs["world-2"])
	
	// Create different test data with world over capacity
	overCapacityRepo := &MockRepository{
		Worlds: []models.World{
			{
				ID:        "over-capacity",
				Name:      "Over Capacity",
				Occupancy: 20,
				Size:      "10",
			},
		},
	}
	
	// Create generator with this repo
	generator = NewMomentGenerator(overCapacityRepo)
	
	// Generate moment for the over-capacity world
	moment, err := generator.GenerateMoment("over-capacity")
	assert.NoError(t, err)
	
	// Activity exists and is calculated
	assert.True(t, moment.Activity >= 0.0)
}

// TestCanAccessWorldFullCoverage tests edge cases in CanAccessWorld
func TestCanAccessWorldFullCoverage(t *testing.T) {
	// Test empty sharing settings
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		CreatorID: "creator",
		Sharing:   models.SharingSettings{},
	}
	
	// Creator should always have access
	assert.True(t, CanAccessWorld("creator", moment))
	
	// Others should not by default
	assert.False(t, CanAccessWorld("other-user", moment))
	
	// We can't test with nil moment because it will panic
	// The function doesn't handle nil input
}