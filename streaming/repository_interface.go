package streaming

import (
	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
)

// RepositoryInterface defines the required repository operations
// This enables proper mocking for testing
type RepositoryInterface interface {
	// GetWorld retrieves a world by ID
	GetWorld(id string) (models.World, error)
	
	// GetAllWorlds returns all worlds
	GetAllWorlds() []models.World
	
	// GetWorldVibe gets a world's vibe
	GetWorldVibe(worldID string) (models.Vibe, error)
}

// Ensure Repository implements the interface
var _ RepositoryInterface = (*repository.Repository)(nil)