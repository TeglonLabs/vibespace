package streaming

import (
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
)

// MomentGenerator creates WorldMoment objects from repository data
type MomentGenerator struct {
	repo RepositoryInterface
}

// NewMomentGenerator creates a new moment generator with the given repository
func NewMomentGenerator(repo RepositoryInterface) *MomentGenerator {
	return &MomentGenerator{
		repo: repo,
	}
}

// GenerateMoment creates a WorldMoment for the specified world
func (g *MomentGenerator) GenerateMoment(worldID string) (*models.WorldMoment, error) {
	// Get the world from the repository
	world, err := g.repo.GetWorld(worldID)
	if err != nil {
		return nil, err
	}

	// Get the world's vibe
	var vibePtr *models.Vibe
	vibe, err := g.repo.GetWorldVibe(worldID)
	if err == nil {
		// Make a copy of the vibe
		vibeCopy := vibe
		vibePtr = &vibeCopy
	}

	// Create the moment - convert timestamp to milliseconds for consistent handling
	now := time.Now().UTC()
	timestamp := now.UnixNano() / int64(time.Millisecond)
	
	// Calculate activity from world properties
	activity := float64(world.Occupancy) / 100.0
	if activity > 1.0 {
		activity = 1.0
	}
	
	// Create sharing settings if not present on world
	sharing := models.SharingSettings{
		IsPublic: false,
		AllowedUsers: []string{},
		ContextLevel: models.ContextLevelPartial,
	}
	if len(world.Sharing.AllowedUsers) > 0 || world.Sharing.IsPublic {
		// If world has sharing settings, inherit them
		sharing = models.SharingSettings{
			IsPublic: world.Sharing.IsPublic,
			AllowedUsers: append([]string{}, world.Sharing.AllowedUsers...),
			ContextLevel: world.Sharing.ContextLevel,
		}
	}
	
	moment := &models.WorldMoment{
		WorldID:     worldID,
		Timestamp:   timestamp,
		VibeID:      world.CurrentVibe,
		Vibe:        vibePtr,
		Occupancy:   world.Occupancy,     // Use current world occupancy
		Activity:    activity,            // Use calculated activity level
		SensorData:  models.SensorData{}, // Initialize empty sensor data
		CreatorID:   world.CreatorID,     // Inherit creator from world
		Viewers:     []string{},          // Initialize empty viewers list
		Sharing:     sharing,             // Use the sharing settings
		CustomData:  "",                  // Initialize empty custom data
	}

	return moment, nil
}

// GenerateAllMoments creates WorldMoment objects for all worlds in the repository
func (g *MomentGenerator) GenerateAllMoments() ([]*models.WorldMoment, error) {
	// Get all worlds
	worlds := g.repo.GetAllWorlds()
	
	// Check if worlds is nil (this actually checks for the nil slice condition)
	if worlds == nil {
		return nil, &GeneratorError{message: "failed to get worlds: nil slice returned"}
	}
	
	moments := make([]*models.WorldMoment, 0, len(worlds))
	
	// Generate a moment for each world
	for _, world := range worlds {
		moment, err := g.GenerateMoment(world.ID)
		if err != nil {
			continue // Skip worlds with errors
		}
		moments = append(moments, moment)
	}
	
	return moments, nil
}

// GeneratorError represents an error when generating moments
type GeneratorError struct {
	message string
}

// Error implements the error interface
func (e *GeneratorError) Error() string {
	return e.message
}

// calculateActivity determines the activity level of a world based on its properties
func calculateActivity(world models.World) float64 {
	// This is a placeholder implementation
	// In a real system, this would consider recent changes, occupancy trends, sensor data, etc.
	
	// For now, just use occupancy as a simple activity metric
	activity := float64(world.Occupancy) / 100.0
	if activity > 1.0 {
		activity = 1.0
	}
	
	return activity
}