package streaming

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/stretchr/testify/assert"
)

// TestGenerateAllMomentsError tests error handling in GenerateAllMoments
func TestGenerateAllMomentsError(t *testing.T) {
	// Create a mock repository that returns errors
	repo := &errorRepo{}
	
	// Create a moment generator with the error repository
	generator := NewMomentGenerator(repo)
	
	// Test error path - repository fails to get worlds
	moments, err := generator.GenerateAllMoments()
	assert.Error(t, err)
	assert.Nil(t, moments)
	assert.Contains(t, err.Error(), "failed to get worlds")
}

// errorRepo is a mock repository that returns errors
type errorRepo struct{}

func (r *errorRepo) GetVibe(id string) (models.Vibe, error) {
	return models.Vibe{}, assert.AnError
}

func (r *errorRepo) GetAllVibes() []models.Vibe {
	return []models.Vibe{}
}

func (r *errorRepo) AddVibe(vibe models.Vibe) error {
	return assert.AnError
}

func (r *errorRepo) UpdateVibe(vibe models.Vibe) error {
	return assert.AnError
}

func (r *errorRepo) DeleteVibe(id string) error {
	return assert.AnError
}

func (r *errorRepo) GetWorld(id string) (models.World, error) {
	return models.World{}, assert.AnError
}

func (r *errorRepo) GetAllWorlds() []models.World {
	// Return a nil slice to simulate the error condition
	// Go's nil slices are actually represented as []Type(nil)
	// So we need to return that specifically instead of just nil
	var nilSlice []models.World
	return nilSlice
}

func (r *errorRepo) AddWorld(world models.World) error {
	return assert.AnError
}

func (r *errorRepo) UpdateWorld(world models.World) error {
	return assert.AnError
}

func (r *errorRepo) DeleteWorld(id string) error {
	return assert.AnError
}

func (r *errorRepo) SetWorldVibe(worldID, vibeID string) error {
	return assert.AnError
}

func (r *errorRepo) GetWorldVibe(worldID string) (models.Vibe, error) {
	return models.Vibe{}, assert.AnError
}

// Main tests for normal operation
func TestGenerateAllMoments(t *testing.T) {
	// Create a new empty repository for this test
	repo := repository.NewRepository()
	
	// Clear any existing worlds if using a shared repository
	for _, world := range repo.GetAllWorlds() {
		repo.DeleteWorld(world.ID)
	}
	
	// Create test worlds
	worlds := []models.World{
		{
			ID:        "world1",
			Name:      "Test World 1",
			Type:      models.WorldTypeVirtual,
			Occupancy: 10,
			CreatorID: "user1",
			Sharing: models.SharingSettings{
				IsPublic:     true,
				ContextLevel: models.ContextLevelFull,
			},
		},
		{
			ID:        "world2",
			Name:      "Test World 2",
			Type:      models.WorldTypePhysical,
			Occupancy: 50,
			CreatorID: "user2",
			Sharing: models.SharingSettings{
				IsPublic:     false,
				AllowedUsers: []string{"user3", "user4"},
				ContextLevel: models.ContextLevelPartial,
			},
		},
		{
			ID:        "world3",
			Name:      "Test World 3",
			Type:      models.WorldTypeHybrid,
			Occupancy: 120, // Will be capped at 1.0 activity
			CreatorID: "user5",
		},
	}
	
	// Add worlds to repository
	for _, world := range worlds {
		err := repo.AddWorld(world)
		assert.NoError(t, err)
	}
	
	// Create vibes and associate them with worlds
	vibes := []models.Vibe{
		{
			ID:   "vibe1",
			Name: "Happy Vibe",
			Mood: models.MoodEnergetic,
			Energy: 0.9,
		},
		{
			ID:   "vibe2",
			Name: "Calm Vibe",
			Mood: models.MoodCalm,
			Energy: 0.3,
		},
	}
	
	for _, vibe := range vibes {
		err := repo.AddVibe(vibe)
		assert.NoError(t, err)
	}
	
	// Associate vibes with worlds
	err := repo.SetWorldVibe("world1", "vibe1")
	assert.NoError(t, err)
	
	err = repo.SetWorldVibe("world2", "vibe2")
	assert.NoError(t, err)
	
	// Create a moment generator
	generator := NewMomentGenerator(repo)
	
	// Generate all moments
	moments, err := generator.GenerateAllMoments()
	assert.NoError(t, err)
	
	// Should have 3 moments (one for each world)
	assert.Equal(t, 3, len(moments))
	
	// Check each moment
	worldMap := make(map[string]*models.WorldMoment)
	for _, moment := range moments {
		worldMap[moment.WorldID] = moment
	}
	
	// Verify world1
	moment1 := worldMap["world1"]
	assert.NotNil(t, moment1)
	assert.Equal(t, "world1", moment1.WorldID)
	assert.Equal(t, "vibe1", moment1.VibeID)
	assert.NotNil(t, moment1.Vibe)
	assert.Equal(t, "Happy Vibe", moment1.Vibe.Name)
	assert.Equal(t, 10, moment1.Occupancy)
	assert.Equal(t, 0.1, moment1.Activity)
	assert.Equal(t, "user1", moment1.CreatorID)
	assert.True(t, moment1.Sharing.IsPublic)
	assert.Equal(t, models.ContextLevelFull, moment1.Sharing.ContextLevel)
	
	// Verify world2
	moment2 := worldMap["world2"]
	assert.NotNil(t, moment2)
	assert.Equal(t, "world2", moment2.WorldID)
	assert.Equal(t, "vibe2", moment2.VibeID)
	assert.NotNil(t, moment2.Vibe)
	assert.Equal(t, "Calm Vibe", moment2.Vibe.Name)
	assert.Equal(t, 50, moment2.Occupancy)
	assert.Equal(t, 0.5, moment2.Activity)
	assert.Equal(t, "user2", moment2.CreatorID)
	assert.False(t, moment2.Sharing.IsPublic)
	assert.Equal(t, 2, len(moment2.Sharing.AllowedUsers))
	assert.Contains(t, moment2.Sharing.AllowedUsers, "user3")
	assert.Contains(t, moment2.Sharing.AllowedUsers, "user4")
	assert.Equal(t, models.ContextLevelPartial, moment2.Sharing.ContextLevel)
	
	// Verify world3
	moment3 := worldMap["world3"]
	assert.NotNil(t, moment3)
	assert.Equal(t, "world3", moment3.WorldID)
	assert.Equal(t, "", moment3.VibeID) // No vibe set
	assert.Nil(t, moment3.Vibe)
	assert.Equal(t, 120, moment3.Occupancy)
	assert.Equal(t, 1.0, moment3.Activity) // Capped at 1.0
	assert.Equal(t, "user5", moment3.CreatorID)
	// Default sharing settings
	assert.False(t, moment3.Sharing.IsPublic)
	assert.Equal(t, 0, len(moment3.Sharing.AllowedUsers))
	assert.Equal(t, models.ContextLevelPartial, moment3.Sharing.ContextLevel)
}

func TestCalculateActivity(t *testing.T) {
	testCases := []struct {
		name         string
		world        models.World
		expected     float64
	}{
		{
			name: "Zero occupancy",
			world: models.World{
				ID:        "world1",
				Occupancy: 0,
			},
			expected: 0.0,
		},
		{
			name: "Normal occupancy",
			world: models.World{
				ID:        "world2",
				Occupancy: 75,
			},
			expected: 0.75,
		},
		{
			name: "Full occupancy",
			world: models.World{
				ID:        "world3",
				Occupancy: 100,
			},
			expected: 1.0,
		},
		{
			name: "Over capacity",
			world: models.World{
				ID:        "world4",
				Occupancy: 150,
			},
			expected: 1.0, // Capped at 1.0
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			activity := calculateActivity(tc.world)
			assert.Equal(t, tc.expected, activity)
		})
	}
}