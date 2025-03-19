package mocks

import (
	"errors"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
)

// MockRepository implements a simple in-memory repository for testing
type MockRepository struct {
	vibes  map[string]models.Vibe
	worlds map[string]models.World
}

// NewMockRepository creates a new mock repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		vibes:  make(map[string]models.Vibe),
		worlds: make(map[string]models.World),
	}
}

// GetVibe retrieves a vibe by ID
func (r *MockRepository) GetVibe(id string) (models.Vibe, error) {
	vibe, ok := r.vibes[id]
	if !ok {
		return models.Vibe{}, errors.New("vibe not found")
	}
	return vibe, nil
}

// GetAllVibes returns all vibes
func (r *MockRepository) GetAllVibes() []models.Vibe {
	vibes := make([]models.Vibe, 0, len(r.vibes))
	for _, vibe := range r.vibes {
		vibes = append(vibes, vibe)
	}
	return vibes
}

// SaveVibe saves a vibe
func (r *MockRepository) SaveVibe(vibe *models.Vibe) error {
	r.vibes[vibe.ID] = *vibe
	return nil
}

// GetWorld retrieves a world by ID
func (r *MockRepository) GetWorld(id string) (models.World, error) {
	world, ok := r.worlds[id]
	if !ok {
		return models.World{}, errors.New("world not found")
	}
	return world, nil
}

// GetAllWorlds returns all worlds
func (r *MockRepository) GetAllWorlds() []models.World {
	worlds := make([]models.World, 0, len(r.worlds))
	for _, world := range r.worlds {
		worlds = append(worlds, world)
	}
	return worlds
}

// SaveWorld saves a world
func (r *MockRepository) SaveWorld(world *models.World) error {
	r.worlds[world.ID] = *world
	return nil
}

// GetWorldVibe gets a world's vibe
func (r *MockRepository) GetWorldVibe(worldID string) (models.Vibe, error) {
	world, err := r.GetWorld(worldID)
	if err != nil {
		return models.Vibe{}, err
	}
	
	if world.CurrentVibe == "" {
		return models.Vibe{}, errors.New("world has no vibe")
	}
	
	return r.GetVibe(world.CurrentVibe)
}

// AddTestData adds test data to the repository
func (r *MockRepository) AddTestData() {
	// Add a test world
	world := &models.World{
		ID:          "test-world",
		Name:        "Test World",
		Description: "A world for testing",
		Type:        models.WorldTypePhysical,
		CreatorID:   "test-user",
	}
	r.SaveWorld(world)
	
	// Add a test vibe
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Description: "A vibe for testing",
		Energy:      0.5,
		Mood:        "neutral",
	}
	r.SaveVibe(vibe)
	
	// Link vibe to world
	world.CurrentVibe = vibe.ID
	r.SaveWorld(world)
}

// Verify that MockRepository implements the RepositoryInterface
var _ streaming.RepositoryInterface = (*MockRepository)(nil)