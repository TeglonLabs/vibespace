package repository

import (
	"errors"
	"sync"

	"github.com/bmorphism/vibespace-mcp-go/models"
)

var (
	ErrVibeNotFound  = errors.New("vibe not found")
	ErrWorldNotFound = errors.New("world not found")
	ErrVibeInUse     = errors.New("vibe is currently used by one or more worlds")
)

// Define interfaces for better testability and separation of concerns

// VibeRepository defines the interface for vibe operations
type VibeRepository interface {
	GetVibe(id string) (models.Vibe, error)
	GetAllVibes() []models.Vibe
	AddVibe(vibe models.Vibe) error
	UpdateVibe(vibe models.Vibe) error
	DeleteVibe(id string) error
}

// WorldRepository defines the interface for world operations
type WorldRepository interface {
	GetWorld(id string) (models.World, error)
	GetAllWorlds() []models.World
	AddWorld(world models.World) error
	UpdateWorld(world models.World) error
	DeleteWorld(id string) error
}

// VibeWorldRepository combines both interfaces and adds relation operations
type VibeWorldRepository interface {
	VibeRepository
	WorldRepository
	SetWorldVibe(worldID, vibeID string) error
	GetWorldVibe(worldID string) (models.Vibe, error)
}

// Repository handles the storage and retrieval of vibes and worlds
type Repository struct {
	vibes  map[string]models.Vibe
	worlds map[string]models.World
	mu     sync.RWMutex
}

// Ensure Repository implements VibeWorldRepository interface
var _ VibeWorldRepository = (*Repository)(nil)

// NewRepository creates a new repository with sample data
func NewRepository() *Repository {
	return NewRepositoryWithSampleData(true)
}

// NewRepositoryWithSampleData creates a new repository with optional sample data
func NewRepositoryWithSampleData(includeSampleData bool) *Repository {
	r := &Repository{
		vibes:  make(map[string]models.Vibe),
		worlds: make(map[string]models.World),
	}

	if !includeSampleData {
		return r
	}

	// Add sample vibes
	temperature := 21.5
	humidity := 45.0
	light := 500.0
	sound := 30.0
	movement := 0.1

	focusedVibe := models.Vibe{
		ID:          "focused-flow",
		Name:        "Focused Flow",
		Description: "A concentration-enhancing atmosphere for deep work",
		Energy:      0.7,
		Mood:        "focused",
		Colors:      []string{"#1A2B3C", "#2C4B6C", "#3C5B7C"},
		SensorData: models.SensorData{
			Temperature: &temperature,
			Humidity:    &humidity,
			Light:       &light,
			Sound:       &sound,
			Movement:    &movement,
		},
	}

	temperature = 23.0
	humidity = 50.0
	light = 300.0
	sound = 20.0
	movement = 0.05

	calmVibe := models.Vibe{
		ID:          "calm-clarity",
		Name:        "Calm Clarity",
		Description: "A peaceful atmosphere for meditation and mindfulness",
		Energy:      0.3,
		Mood:        "calm",
		Colors:      []string{"#8ECAE6", "#219EBC", "#023047"},
		SensorData: models.SensorData{
			Temperature: &temperature,
			Humidity:    &humidity,
			Light:       &light,
			Sound:       &sound,
			Movement:    &movement,
		},
	}

	temperature = 22.0
	humidity = 40.0
	light = 800.0
	sound = 60.0
	movement = 0.8

	energeticVibe := models.Vibe{
		ID:          "energetic-spark",
		Name:        "Energetic Spark",
		Description: "A high-energy atmosphere for creativity and brainstorming",
		Energy:      0.9,
		Mood:        "energetic",
		Colors:      []string{"#F94144", "#F8961E", "#F9C74F"},
		SensorData: models.SensorData{
			Temperature: &temperature,
			Humidity:    &humidity,
			Light:       &light,
			Sound:       &sound,
			Movement:    &movement,
		},
	}

	r.vibes[focusedVibe.ID] = focusedVibe
	r.vibes[calmVibe.ID] = calmVibe
	r.vibes[energeticVibe.ID] = energeticVibe

	// Add sample worlds
	officeWorld := models.World{
		ID:          "office-space",
		Name:        "Modern Office",
		Description: "An open-concept workspace designed for collaboration",
		Type:        models.WorldTypePhysical,
		Location:    "Floor 3, Building A",
		CurrentVibe: focusedVibe.ID,
		Size:        "Medium (500 sqm)",
		Features:    []string{"standing desks", "natural light", "acoustic panels"},
	}

	virtualWorld := models.World{
		ID:          "virtual-garden",
		Name:        "Zen Garden",
		Description: "A virtual peaceful garden for mental relaxation",
		Type:        models.WorldTypeVirtual,
		Location:    "https://garden.vibespace.io",
		CurrentVibe: calmVibe.ID,
		Features:    []string{"water sounds", "interactive plants", "meditation spots"},
	}

	hybridWorld := models.World{
		ID:          "hybrid-studio",
		Name:        "Creative Studio",
		Description: "A hybrid space for both physical and virtual creative collaboration",
		Type:        models.WorldTypeHybrid,
		Location:    "Floor 5, Innovation Center + VR instance",
		CurrentVibe: energeticVibe.ID,
		Size:        "Large (1000 sqm physical + unlimited virtual)",
		Features:    []string{"AR overlays", "digital whiteboard", "spatial audio"},
	}

	r.worlds[officeWorld.ID] = officeWorld
	r.worlds[virtualWorld.ID] = virtualWorld
	r.worlds[hybridWorld.ID] = hybridWorld

	return r
}

// GetVibe retrieves a vibe by ID
func (r *Repository) GetVibe(id string) (models.Vibe, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	vibe, ok := r.vibes[id]
	if !ok {
		return models.Vibe{}, ErrVibeNotFound
	}
	return vibe, nil
}

// GetAllVibes returns all vibes
func (r *Repository) GetAllVibes() []models.Vibe {
	r.mu.RLock()
	defer r.mu.RUnlock()

	vibes := make([]models.Vibe, 0, len(r.vibes))
	for _, vibe := range r.vibes {
		vibes = append(vibes, vibe)
	}
	return vibes
}

// AddVibe adds a new vibe
func (r *Repository) AddVibe(vibe models.Vibe) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.vibes[vibe.ID] = vibe
	return nil
}

// UpdateVibe updates an existing vibe
func (r *Repository) UpdateVibe(vibe models.Vibe) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.vibes[vibe.ID]
	if !ok {
		return ErrVibeNotFound
	}
	r.vibes[vibe.ID] = vibe
	return nil
}

// DeleteVibe removes a vibe
func (r *Repository) DeleteVibe(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if vibe exists
	_, ok := r.vibes[id]
	if !ok {
		return ErrVibeNotFound
	}

	// Check if vibe is in use by any world
	for _, world := range r.worlds {
		if world.CurrentVibe == id {
			return ErrVibeInUse
		}
	}

	delete(r.vibes, id)
	return nil
}

// GetWorld retrieves a world by ID
func (r *Repository) GetWorld(id string) (models.World, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	world, ok := r.worlds[id]
	if !ok {
		return models.World{}, ErrWorldNotFound
	}
	return world, nil
}

// GetAllWorlds returns all worlds
func (r *Repository) GetAllWorlds() []models.World {
	r.mu.RLock()
	defer r.mu.RUnlock()

	worlds := make([]models.World, 0, len(r.worlds))
	for _, world := range r.worlds {
		worlds = append(worlds, world)
	}
	return worlds
}

// AddWorld adds a new world
func (r *Repository) AddWorld(world models.World) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// If world has a vibe assigned, check if it exists
	if world.CurrentVibe != "" {
		_, ok := r.vibes[world.CurrentVibe]
		if !ok {
			return ErrVibeNotFound
		}
	}

	r.worlds[world.ID] = world
	return nil
}

// UpdateWorld updates an existing world
func (r *Repository) UpdateWorld(world models.World) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.worlds[world.ID]
	if !ok {
		return ErrWorldNotFound
	}

	// If world has a vibe assigned, check if it exists
	if world.CurrentVibe != "" {
		_, ok := r.vibes[world.CurrentVibe]
		if !ok {
			return ErrVibeNotFound
		}
	}

	r.worlds[world.ID] = world
	return nil
}

// DeleteWorld removes a world
func (r *Repository) DeleteWorld(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.worlds[id]
	if !ok {
		return ErrWorldNotFound
	}

	delete(r.worlds, id)
	return nil
}

// SetWorldVibe sets a world's vibe
func (r *Repository) SetWorldVibe(worldID, vibeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	world, ok := r.worlds[worldID]
	if !ok {
		return ErrWorldNotFound
	}

	_, ok = r.vibes[vibeID]
	if !ok {
		return ErrVibeNotFound
	}

	world.CurrentVibe = vibeID
	r.worlds[worldID] = world
	return nil
}

// GetWorldVibe gets a world's vibe
func (r *Repository) GetWorldVibe(worldID string) (models.Vibe, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	world, ok := r.worlds[worldID]
	if !ok {
		return models.Vibe{}, ErrWorldNotFound
	}

	if world.CurrentVibe == "" {
		return models.Vibe{}, ErrVibeNotFound
	}

	vibe, ok := r.vibes[world.CurrentVibe]
	if !ok {
		return models.Vibe{}, ErrVibeNotFound
	}

	return vibe, nil
}