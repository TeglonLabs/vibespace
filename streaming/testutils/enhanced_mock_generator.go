package testutils

import (
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming/test/mocks"
)

// EnhancedMockMomentGenerator extends MockMomentGenerator with error scenarios support
type EnhancedMockMomentGenerator struct {
	repo          *mocks.MockRepository
	GenerateError error
}

// NewEnhancedMockMomentGenerator creates a new mock moment generator with error capabilities
func NewEnhancedMockMomentGenerator(repo *mocks.MockRepository) *EnhancedMockMomentGenerator {
	return &EnhancedMockMomentGenerator{
		repo: repo,
	}
}

// GenerateMoment generates a test moment or returns error
func (g *EnhancedMockMomentGenerator) GenerateMoment(worldID string) (*models.WorldMoment, error) {
	if g.GenerateError != nil {
		return nil, g.GenerateError
	}
	
	// Create a simple moment
	moment := &models.WorldMoment{
		WorldID:   worldID,
		Timestamp: time.Now().Unix(),
	}
	return moment, nil
}

// GenerateAllMoments generates test moments or returns error
func (g *EnhancedMockMomentGenerator) GenerateAllMoments() ([]*models.WorldMoment, error) {
	if g.GenerateError != nil {
		return nil, g.GenerateError
	}
	
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