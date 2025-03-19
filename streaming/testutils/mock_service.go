package testutils

import (
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/bmorphism/vibespace-mcp-go/streaming/test/mocks"
)

// MockMomentGenerator is a simplified moment generator for testing
type MockMomentGenerator struct {
	repo *mocks.MockRepository
}

// NewMockMomentGenerator creates a new mock moment generator
func NewMockMomentGenerator(repo *mocks.MockRepository) *MockMomentGenerator {
	return &MockMomentGenerator{
		repo: repo,
	}
}

// GenerateMoment generates a test moment
func (g *MockMomentGenerator) GenerateMoment(worldID string) (*models.WorldMoment, error) {
	// Create a simple moment
	moment := &models.WorldMoment{
		WorldID:   worldID,
		Timestamp: time.Now().Unix(),
	}
	return moment, nil
}

// GenerateAllMoments generates test moments for all worlds
func (g *MockMomentGenerator) GenerateAllMoments() ([]*models.WorldMoment, error) {
	worlds := g.repo.GetAllWorlds()
	moments := make([]*models.WorldMoment, 0, len(worlds))
	
	for _, world := range worlds {
		moment := &models.WorldMoment{
			WorldID:   world.ID,
			Timestamp: time.Now().Unix(),
		}
		moments = append(moments, moment)
	}
	
	return moments, nil
}

// CreateMockStreamingService creates a streaming service for testing
func CreateMockStreamingService(repo *mocks.MockRepository, config *streaming.StreamingConfig, natsClient streaming.NATSClientInterface) *streaming.StreamingService {
	service := &streaming.StreamingService{}
	
	// Create a mock moment generator
	mockGenerator := NewMockMomentGenerator(repo)
	
	// Set test helpers
	service.SetConfig(config)
	service.SetClient(natsClient)
	service.SetMomentGenerator(mockGenerator)
	service.SetRepository(repo)
	
	return service
}