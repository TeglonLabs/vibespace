package streaming

import (
	"github.com/bmorphism/vibespace-mcp-go/models"
)

// MomentGeneratorInterface defines the interface for generating world moments
// This enables proper mocking for testing
type MomentGeneratorInterface interface {
	// GenerateMoment generates a moment for a specific world
	GenerateMoment(worldID string) (*models.WorldMoment, error)
	
	// GenerateAllMoments generates moments for all worlds
	GenerateAllMoments() ([]*models.WorldMoment, error)
}

// Ensure MomentGenerator implements the interface
var _ MomentGeneratorInterface = (*MomentGenerator)(nil)