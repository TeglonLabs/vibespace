package streaming

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

// TestStreamSingleWorld tests the StreamSingleWorld method
func TestStreamSingleWorld(t *testing.T) {
	// Create a mock client
	mockClient := NewMockNATSClient()
	mockClient.Connect()
	
	// Create a mock repository with test data
	mockRepo := &MockRepository{
		Worlds: []models.World{
			{
				ID:   "test-world",
				Name: "Test World",
			},
		},
		Vibes: map[string]*models.Vibe{
			"test-world": {
				ID:     "test-vibe",
				Name:   "Test Vibe",
				Energy: 0.5,
			},
		},
	}
	
	// Create a test config
	config := &StreamingConfig{
		StreamInterval: 100 * time.Millisecond,
	}
	
	// Create service
	service := &StreamingService{}
	service.SetConfig(config)
	service.SetClient(mockClient)
	service.SetRepository(mockRepo)
	service.SetMomentGenerator(NewMomentGenerator(mockRepo))
	
	// Test streaming a single world
	err := service.StreamSingleWorld("test-world", "test-user")
	assert.NoError(t, err)
	
	// Verify the moment was published
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
	
	// Create a generator that produces errors
	mockGenerator := &MockMomentGenerator{
		GenerateError: assert.AnError,
	}
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

// MockRepository for service tests
type MockRepository struct {
	Worlds []models.World
	Vibes  map[string]*models.Vibe
}

func (r *MockRepository) GetWorld(id string) (models.World, error) {
	for _, world := range r.Worlds {
		if world.ID == id {
			return world, nil
		}
	}
	return models.World{}, assert.AnError
}

func (r *MockRepository) GetAllWorlds() []models.World {
	return r.Worlds
}

func (r *MockRepository) GetWorldVibe(worldID string) (models.Vibe, error) {
	if vibe, ok := r.Vibes[worldID]; ok {
		return *vibe, nil
	}
	return models.Vibe{}, assert.AnError
}

// MockMomentGenerator for service tests
type MockMomentGenerator struct {
	GenerateError    error
	GenerateAllError error
}

func (g *MockMomentGenerator) GenerateMoment(worldID string) (*models.WorldMoment, error) {
	if g.GenerateError != nil {
		return nil, g.GenerateError
	}
	
	moment := &models.WorldMoment{
		WorldID:   worldID,
		Timestamp: time.Now().Unix(),
	}
	return moment, nil
}

func (g *MockMomentGenerator) GenerateAllMoments() ([]*models.WorldMoment, error) {
	if g.GenerateAllError != nil {
		return nil, g.GenerateAllError
	}
	
	moments := []*models.WorldMoment{
		{
			WorldID:   "test-world-1",
			Timestamp: time.Now().Unix(),
		},
		{
			WorldID:   "test-world-2",
			Timestamp: time.Now().Unix(),
		},
	}
	return moments, nil
}