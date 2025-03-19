package mocks

import (
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
)

// MockMomentGenerator implements a moment generator for testing
type MockMomentGenerator struct {
	Repo               *MockRepository
	GenerateError      error
	GenerateAllError   error
	GeneratedMoments   []*models.WorldMoment
	GeneratedAllCalls  int
}

// NewMockMomentGenerator creates a new mock moment generator
func NewMockMomentGenerator(repo *MockRepository) *MockMomentGenerator {
	return &MockMomentGenerator{
		Repo:             repo,
		GeneratedMoments: []*models.WorldMoment{},
	}
}

// GenerateMoment generates a test moment or returns configured error
func (g *MockMomentGenerator) GenerateMoment(worldID string) (*models.WorldMoment, error) {
	if g.GenerateError != nil {
		return nil, g.GenerateError
	}
	
	// Create a simple moment
	moment := &models.WorldMoment{
		WorldID:   worldID,
		Timestamp: time.Now().Unix(),
		CreatorID: "test-user",
		Sharing: models.SharingSettings{
			ContextLevel: models.ContextLevelPartial,
		},
	}
	g.GeneratedMoments = append(g.GeneratedMoments, moment)
	return moment, nil
}

// GenerateAllMoments generates test moments for all worlds or returns error
func (g *MockMomentGenerator) GenerateAllMoments() ([]*models.WorldMoment, error) {
	g.GeneratedAllCalls++
	
	if g.GenerateAllError != nil {
		return nil, g.GenerateAllError
	}
	
	worlds := g.Repo.GetAllWorlds()
	moments := make([]*models.WorldMoment, 0, len(worlds))
	
	for _, world := range worlds {
		moment := &models.WorldMoment{
			WorldID:   world.ID,
			Timestamp: time.Now().Unix(),
			CreatorID: world.CreatorID,
			Sharing: models.SharingSettings{
				ContextLevel: models.ContextLevelPartial,
			},
		}
		moments = append(moments, moment)
	}
	
	return moments, nil
}

// Verify that MockMomentGenerator implements the MomentGeneratorInterface
var _ streaming.MomentGeneratorInterface = (*MockMomentGenerator)(nil)