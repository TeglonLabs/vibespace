package tests

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
)

func TestRepository(t *testing.T) {
	// Create a new repository
	repo := repository.NewRepository()

	// Test GetAllVibes
	vibes := repo.GetAllVibes()
	if len(vibes) != 3 {
		t.Errorf("Expected 3 vibes, got %d", len(vibes))
	}

	// Test GetVibe
	vibe, err := repo.GetVibe("focused-flow")
	if err != nil {
		t.Errorf("Error getting vibe: %v", err)
	}
	if vibe.Name != "Focused Flow" {
		t.Errorf("Expected vibe name 'Focused Flow', got '%s'", vibe.Name)
	}

	// Test GetVibe with non-existent ID
	_, err = repo.GetVibe("non-existent")
	if err != repository.ErrVibeNotFound {
		t.Errorf("Expected ErrVibeNotFound, got %v", err)
	}

	// Test AddVibe
	newVibe := models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Description: "A test vibe",
		Energy:      0.5,
		Mood:        "neutral",
		Colors:      []string{"#FF0000", "#00FF00", "#0000FF"},
	}
	err = repo.AddVibe(newVibe)
	if err != nil {
		t.Errorf("Error adding vibe: %v", err)
	}

	// Verify vibe was added
	vibe, err = repo.GetVibe("test-vibe")
	if err != nil {
		t.Errorf("Error getting vibe: %v", err)
	}
	if vibe.Name != "Test Vibe" {
		t.Errorf("Expected vibe name 'Test Vibe', got '%s'", vibe.Name)
	}

	// Test UpdateVibe
	vibe.Name = "Updated Test Vibe"
	err = repo.UpdateVibe(vibe)
	if err != nil {
		t.Errorf("Error updating vibe: %v", err)
	}

	// Verify vibe was updated
	vibe, err = repo.GetVibe("test-vibe")
	if err != nil {
		t.Errorf("Error getting vibe: %v", err)
	}
	if vibe.Name != "Updated Test Vibe" {
		t.Errorf("Expected vibe name 'Updated Test Vibe', got '%s'", vibe.Name)
	}

	// Test UpdateVibe with non-existent ID
	nonExistentVibe := models.Vibe{
		ID:   "non-existent",
		Name: "Non-existent Vibe",
	}
	err = repo.UpdateVibe(nonExistentVibe)
	if err != repository.ErrVibeNotFound {
		t.Errorf("Expected ErrVibeNotFound, got %v", err)
	}

	// Test DeleteVibe
	err = repo.DeleteVibe("test-vibe")
	if err != nil {
		t.Errorf("Error deleting vibe: %v", err)
	}

	// Verify vibe was deleted
	_, err = repo.GetVibe("test-vibe")
	if err != repository.ErrVibeNotFound {
		t.Errorf("Expected ErrVibeNotFound, got %v", err)
	}

	// Test DeleteVibe with vibe in use
	// The sample data includes a world using the "focused-flow" vibe
	err = repo.DeleteVibe("focused-flow")
	if err != repository.ErrVibeInUse {
		t.Errorf("Expected ErrVibeInUse, got %v", err)
	}

	// Test GetAllWorlds
	worlds := repo.GetAllWorlds()
	if len(worlds) != 3 {
		t.Errorf("Expected 3 worlds, got %d", len(worlds))
	}

	// Test GetWorld
	world, err := repo.GetWorld("office-space")
	if err != nil {
		t.Errorf("Error getting world: %v", err)
	}
	if world.Name != "Modern Office" {
		t.Errorf("Expected world name 'Modern Office', got '%s'", world.Name)
	}

	// Test GetWorld with non-existent ID
	_, err = repo.GetWorld("non-existent")
	if err != repository.ErrWorldNotFound {
		t.Errorf("Expected ErrWorldNotFound, got %v", err)
	}

	// Test AddWorld
	newWorld := models.World{
		ID:          "test-world",
		Name:        "Test World",
		Description: "A test world",
		Type:        models.WorldTypeVirtual,
		Location:    "https://test.world",
		CurrentVibe: "calm-clarity", // Use an existing vibe ID
		Features:    []string{"feature1", "feature2"},
	}
	err = repo.AddWorld(newWorld)
	if err != nil {
		t.Errorf("Error adding world: %v", err)
	}

	// Verify world was added
	world, err = repo.GetWorld("test-world")
	if err != nil {
		t.Errorf("Error getting world: %v", err)
	}
	if world.Name != "Test World" {
		t.Errorf("Expected world name 'Test World', got '%s'", world.Name)
	}

	// Test UpdateWorld
	world.Name = "Updated Test World"
	err = repo.UpdateWorld(world)
	if err != nil {
		t.Errorf("Error updating world: %v", err)
	}

	// Verify world was updated
	world, err = repo.GetWorld("test-world")
	if err != nil {
		t.Errorf("Error getting world: %v", err)
	}
	if world.Name != "Updated Test World" {
		t.Errorf("Expected world name 'Updated Test World', got '%s'", world.Name)
	}

	// Test SetWorldVibe
	err = repo.SetWorldVibe("test-world", "energetic-spark")
	if err != nil {
		t.Errorf("Error setting world vibe: %v", err)
	}

	// Verify vibe was set
	vibe, err = repo.GetWorldVibe("test-world")
	if err != nil {
		t.Errorf("Error getting world vibe: %v", err)
	}
	if vibe.ID != "energetic-spark" {
		t.Errorf("Expected vibe ID 'energetic-spark', got '%s'", vibe.ID)
	}

	// Test DeleteWorld
	err = repo.DeleteWorld("test-world")
	if err != nil {
		t.Errorf("Error deleting world: %v", err)
	}

	// Verify world was deleted
	_, err = repo.GetWorld("test-world")
	if err != repository.ErrWorldNotFound {
		t.Errorf("Expected ErrWorldNotFound, got %v", err)
	}

	// Test AddWorld with non-existent vibe
	badWorld := models.World{
		ID:          "bad-world",
		Name:        "Bad World",
		Description: "A world with a non-existent vibe",
		Type:        models.WorldTypeVirtual,
		CurrentVibe: "non-existent-vibe",
	}
	err = repo.AddWorld(badWorld)
	if err != repository.ErrVibeNotFound {
		t.Errorf("Expected ErrVibeNotFound, got %v", err)
	}
}

// TestSensorData tests the functionality for managing vibe sensor data
func TestSensorData(t *testing.T) {
	// Create a new repository
	repo := repository.NewRepository()
	
	// Test 1: Create a vibe with sensor data
	temperature := 22.5
	humidity := 40.0
	light := 650.0
	sound := 35.0
	movement := 0.3
	
	newVibe := models.Vibe{
		ID:          "sensor-test-vibe",
		Name:        "Sensor Test Vibe",
		Description: "A vibe for testing sensor data",
		Energy:      0.6,
		Mood:        "productive",
		Colors:      []string{"#222222", "#444444", "#666666"},
		SensorData: models.SensorData{
			Temperature: &temperature,
			Humidity:    &humidity,
			Light:       &light,
			Sound:       &sound,
			Movement:    &movement,
		},
	}
	
	err := repo.AddVibe(newVibe)
	if err != nil {
		t.Errorf("Error adding vibe with sensor data: %v", err)
	}
	
	// Test 2: Retrieve and verify sensor data
	retrievedVibe, err := repo.GetVibe("sensor-test-vibe")
	if err != nil {
		t.Errorf("Error getting vibe with sensor data: %v", err)
	}
	
	if retrievedVibe.SensorData.Temperature == nil || *retrievedVibe.SensorData.Temperature != temperature {
		t.Errorf("Expected temperature %f, got %v", temperature, retrievedVibe.SensorData.Temperature)
	}
	
	if retrievedVibe.SensorData.Humidity == nil || *retrievedVibe.SensorData.Humidity != humidity {
		t.Errorf("Expected humidity %f, got %v", humidity, retrievedVibe.SensorData.Humidity)
	}
	
	if retrievedVibe.SensorData.Light == nil || *retrievedVibe.SensorData.Light != light {
		t.Errorf("Expected light %f, got %v", light, retrievedVibe.SensorData.Light)
	}
	
	// Test 3: Update only the sensor data
	newTemperature := 24.0
	newHumidity := 45.0
	
	retrievedVibe.SensorData.Temperature = &newTemperature
	retrievedVibe.SensorData.Humidity = &newHumidity
	
	err = repo.UpdateVibe(retrievedVibe)
	if err != nil {
		t.Errorf("Error updating vibe sensor data: %v", err)
	}
	
	// Test 4: Verify the updated sensor data
	updatedVibe, err := repo.GetVibe("sensor-test-vibe")
	if err != nil {
		t.Errorf("Error getting updated vibe: %v", err)
	}
	
	if updatedVibe.SensorData.Temperature == nil || *updatedVibe.SensorData.Temperature != newTemperature {
		t.Errorf("Expected updated temperature %f, got %v", newTemperature, updatedVibe.SensorData.Temperature)
	}
	
	if updatedVibe.SensorData.Humidity == nil || *updatedVibe.SensorData.Humidity != newHumidity {
		t.Errorf("Expected updated humidity %f, got %v", newHumidity, updatedVibe.SensorData.Humidity)
	}
	
	// Test 5: Create a vibe without sensor data, then add it later
	basicVibe := models.Vibe{
		ID:          "basic-vibe",
		Name:        "Basic Vibe",
		Description: "A vibe without initial sensor data",
		Energy:      0.5,
		Mood:        "neutral",
		Colors:      []string{"#AAAAAA", "#BBBBBB", "#CCCCCC"},
	}
	
	err = repo.AddVibe(basicVibe)
	if err != nil {
		t.Errorf("Error adding basic vibe: %v", err)
	}
	
	// Add sensor data later
	retrievedBasicVibe, err := repo.GetVibe("basic-vibe")
	if err != nil {
		t.Errorf("Error getting basic vibe: %v", err)
	}
	
	addedTemperature := 21.0
	addedHumidity := 38.0
	addedLight := 500.0
	
	retrievedBasicVibe.SensorData = models.SensorData{
		Temperature: &addedTemperature,
		Humidity:    &addedHumidity,
		Light:       &addedLight,
	}
	
	err = repo.UpdateVibe(retrievedBasicVibe)
	if err != nil {
		t.Errorf("Error adding sensor data to existing vibe: %v", err)
	}
	
	// Verify added sensor data
	vibeWithAddedData, err := repo.GetVibe("basic-vibe")
	if err != nil {
		t.Errorf("Error getting vibe with added sensor data: %v", err)
	}
	
	if vibeWithAddedData.SensorData.Temperature == nil || *vibeWithAddedData.SensorData.Temperature != addedTemperature {
		t.Errorf("Expected added temperature %f, got %v", addedTemperature, vibeWithAddedData.SensorData.Temperature)
	}
	
	// Test 6: Check partial sensor data (some fields nil)
	partialVibe := models.Vibe{
		ID:          "partial-sensor-vibe",
		Name:        "Partial Sensor Vibe",
		Description: "A vibe with partial sensor data",
		Energy:      0.7,
		Mood:        "relaxed",
		Colors:      []string{"#112233", "#445566", "#778899"},
		SensorData: models.SensorData{
			Temperature: &temperature, // Only temperature is set
		},
	}
	
	err = repo.AddVibe(partialVibe)
	if err != nil {
		t.Errorf("Error adding vibe with partial sensor data: %v", err)
	}
	
	// Verify partial sensor data
	retrievedPartialVibe, err := repo.GetVibe("partial-sensor-vibe")
	if err != nil {
		t.Errorf("Error getting vibe with partial sensor data: %v", err)
	}
	
	if retrievedPartialVibe.SensorData.Temperature == nil {
		t.Errorf("Expected temperature to be set, got nil")
	}
	
	if retrievedPartialVibe.SensorData.Humidity != nil {
		t.Errorf("Expected humidity to be nil, got %v", *retrievedPartialVibe.SensorData.Humidity)
	}
}

// TestWorldFeatures tests the functionality for managing world features
func TestWorldFeatures(t *testing.T) {
	// Create a new repository
	repo := repository.NewRepository()

	// Test 1: Create a world with multiple features
	multiFeatureWorld := models.World{
		ID:          "multi-feature-world",
		Name:        "Multi-Feature World",
		Description: "A world with multiple features for testing",
		Type:        models.WorldTypeVirtual,
		Location:    "https://features.world",
		CurrentVibe: "calm-clarity", // Use an existing vibe ID
		Features:    []string{"feature1", "feature2", "feature3", "feature4"},
	}

	err := repo.AddWorld(multiFeatureWorld)
	if err != nil {
		t.Errorf("Error adding world with multiple features: %v", err)
	}

	// Test 2: Retrieve and verify features
	retrievedWorld, err := repo.GetWorld("multi-feature-world")
	if err != nil {
		t.Errorf("Error getting world with features: %v", err)
	}

	if len(retrievedWorld.Features) != 4 {
		t.Errorf("Expected 4 features, got %d", len(retrievedWorld.Features))
	}

	// Check if all features are present
	expectedFeatures := map[string]bool{
		"feature1": false,
		"feature2": false,
		"feature3": false,
		"feature4": false,
	}

	for _, feature := range retrievedWorld.Features {
		expectedFeatures[feature] = true
	}

	for feature, found := range expectedFeatures {
		if !found {
			t.Errorf("Expected feature '%s' not found in world", feature)
		}
	}

	// Test 3: Update world with modified features (add and remove)
	retrievedWorld.Features = []string{"feature1", "feature3", "feature5", "feature6"} // Remove 2 & 4, add 5 & 6

	err = repo.UpdateWorld(retrievedWorld)
	if err != nil {
		t.Errorf("Error updating world features: %v", err)
	}

	// Test 4: Verify feature updates
	updatedWorld, err := repo.GetWorld("multi-feature-world")
	if err != nil {
		t.Errorf("Error getting updated world: %v", err)
	}

	if len(updatedWorld.Features) != 4 {
		t.Errorf("Expected 4 features after update, got %d", len(updatedWorld.Features))
	}

	// Check if updated features are present
	updatedExpectedFeatures := map[string]bool{
		"feature1": false,
		"feature3": false,
		"feature5": false,
		"feature6": false,
	}

	for _, feature := range updatedWorld.Features {
		updatedExpectedFeatures[feature] = true
	}

	for feature, found := range updatedExpectedFeatures {
		if !found {
			t.Errorf("Expected updated feature '%s' not found in world", feature)
		}
	}

	// Test 5: Create a world without features, then add them later
	basicWorld := models.World{
		ID:          "basic-world",
		Name:        "Basic World",
		Description: "A world without initial features",
		Type:        models.WorldTypePhysical,
		Location:    "Floor 2, Building B",
		CurrentVibe: "energetic-spark", // Use an existing vibe ID
	}

	err = repo.AddWorld(basicWorld)
	if err != nil {
		t.Errorf("Error adding basic world: %v", err)
	}

	// Verify no features
	retrievedBasicWorld, err := repo.GetWorld("basic-world")
	if err != nil {
		t.Errorf("Error getting basic world: %v", err)
	}

	if len(retrievedBasicWorld.Features) != 0 {
		t.Errorf("Expected no features initially, got %d", len(retrievedBasicWorld.Features))
	}

	// Add features
	retrievedBasicWorld.Features = []string{"added-feature1", "added-feature2"}

	err = repo.UpdateWorld(retrievedBasicWorld)
	if err != nil {
		t.Errorf("Error adding features to existing world: %v", err)
	}

	// Verify added features
	worldWithAddedFeatures, err := repo.GetWorld("basic-world")
	if err != nil {
		t.Errorf("Error getting world with added features: %v", err)
	}

	if len(worldWithAddedFeatures.Features) != 2 {
		t.Errorf("Expected 2 added features, got %d", len(worldWithAddedFeatures.Features))
	}

	// Test 6: Empty an existing world's features
	worldToEmpty, err := repo.GetWorld("multi-feature-world")
	if err != nil {
		t.Errorf("Error getting world to empty features: %v", err)
	}

	// Ensure it has features before emptying
	if len(worldToEmpty.Features) == 0 {
		t.Errorf("Expected world to have features before emptying")
	}

	// Empty features
	worldToEmpty.Features = []string{}
	err = repo.UpdateWorld(worldToEmpty)
	if err != nil {
		t.Errorf("Error emptying world features: %v", err)
	}

	// Verify features were emptied
	emptyFeaturesWorld, err := repo.GetWorld("multi-feature-world")
	if err != nil {
		t.Errorf("Error getting world with emptied features: %v", err)
	}

	if len(emptyFeaturesWorld.Features) != 0 {
		t.Errorf("Expected 0 features after emptying, got %d", len(emptyFeaturesWorld.Features))
	}
}

// TestHybridWorlds tests the functionality specific to hybrid worlds
func TestHybridWorlds(t *testing.T) {
	// Create a new repository
	repo := repository.NewRepository()

	// Test 1: Create a new hybrid world with specific hybrid characteristics
	hybridWorld := models.World{
		ID:          "test-hybrid-world",
		Name:        "Test Hybrid World",
		Description: "A hybrid world for specific testing",
		Type:        models.WorldTypeHybrid,
		Location:    "Floor 1, Innovation Hub + https://metaverse.example.com/world-123",
		CurrentVibe: "energetic-spark", // Use an existing vibe ID
		Size:        "Medium physical (200sqm) + Unlimited virtual",
		Features:    []string{"AR overlays", "digital twin", "physical anchors", "spatial audio"},
	}

	err := repo.AddWorld(hybridWorld)
	if err != nil {
		t.Errorf("Error adding hybrid world: %v", err)
	}

	// Test 2: Retrieve and verify hybrid world properties
	retrievedHybrid, err := repo.GetWorld("test-hybrid-world")
	if err != nil {
		t.Errorf("Error getting hybrid world: %v", err)
	}

	if retrievedHybrid.Type != models.WorldTypeHybrid {
		t.Errorf("Expected world type '%s', got '%s'", models.WorldTypeHybrid, retrievedHybrid.Type)
	}

	// Test 3: Convert a physical world to a hybrid world
	physicalWorld := models.World{
		ID:          "physical-to-hybrid",
		Name:        "Physical World",
		Description: "A physical world that will be converted to hybrid",
		Type:        models.WorldTypePhysical,
		Location:    "Floor 3, Building C",
		Size:        "Small (50sqm)",
		Features:    []string{"meeting spaces", "whiteboards"},
	}

	err = repo.AddWorld(physicalWorld)
	if err != nil {
		t.Errorf("Error adding physical world: %v", err)
	}

	// Convert to hybrid
	retrievedPhysical, err := repo.GetWorld("physical-to-hybrid")
	if err != nil {
		t.Errorf("Error getting physical world: %v", err)
	}

	// Update to hybrid
	retrievedPhysical.Type = models.WorldTypeHybrid
	retrievedPhysical.Location = retrievedPhysical.Location + " + https://hybrid.example.com/extension"
	retrievedPhysical.Features = append(retrievedPhysical.Features, "virtual extensions", "digital overlays")
	retrievedPhysical.Size = retrievedPhysical.Size + " + Unlimited virtual"

	err = repo.UpdateWorld(retrievedPhysical)
	if err != nil {
		t.Errorf("Error updating world to hybrid: %v", err)
	}

	// Verify conversion
	convertedWorld, err := repo.GetWorld("physical-to-hybrid")
	if err != nil {
		t.Errorf("Error getting converted world: %v", err)
	}

	if convertedWorld.Type != models.WorldTypeHybrid {
		t.Errorf("Expected converted world type '%s', got '%s'", models.WorldTypeHybrid, convertedWorld.Type)
	}

	// Test 4: Convert a virtual world to a hybrid world
	virtualWorld := models.World{
		ID:          "virtual-to-hybrid",
		Name:        "Virtual World",
		Description: "A virtual world that will be converted to hybrid",
		Type:        models.WorldTypeVirtual,
		Location:    "https://virtual.example.com/space",
		Features:    []string{"3D spaces", "avatars"},
	}

	err = repo.AddWorld(virtualWorld)
	if err != nil {
		t.Errorf("Error adding virtual world: %v", err)
	}

	// Convert to hybrid
	retrievedVirtual, err := repo.GetWorld("virtual-to-hybrid")
	if err != nil {
		t.Errorf("Error getting virtual world: %v", err)
	}

	// Update to hybrid
	retrievedVirtual.Type = models.WorldTypeHybrid
	retrievedVirtual.Location = "Innovation Center Pod 5 + " + retrievedVirtual.Location
	retrievedVirtual.Features = append(retrievedVirtual.Features, "physical installation", "haptic feedback")
	retrievedVirtual.Size = "Medium (100sqm) physical + Virtual"

	err = repo.UpdateWorld(retrievedVirtual)
	if err != nil {
		t.Errorf("Error updating virtual world to hybrid: %v", err)
	}

	// Verify conversion
	convertedVirtualWorld, err := repo.GetWorld("virtual-to-hybrid")
	if err != nil {
		t.Errorf("Error getting converted virtual world: %v", err)
	}

	if convertedVirtualWorld.Type != models.WorldTypeHybrid {
		t.Errorf("Expected converted world type '%s', got '%s'", models.WorldTypeHybrid, convertedVirtualWorld.Type)
	}

	// Test 5: Validate special hybrid features exist in a sample hybrid world
	sampleHybrid, err := repo.GetWorld("hybrid-studio") // This is from the sample data
	if err != nil {
		t.Errorf("Error getting sample hybrid world: %v", err)
	}

	if sampleHybrid.Type != models.WorldTypeHybrid {
		t.Errorf("Expected sample hybrid world type '%s', got '%s'", models.WorldTypeHybrid, sampleHybrid.Type)
	}

	// Check for hybrid-appropriate features
	hybridFeatures := map[string]bool{
		"AR overlays":        false,
		"digital whiteboard": false,
		"spatial audio":      false,
	}

	for _, feature := range sampleHybrid.Features {
		if _, exists := hybridFeatures[feature]; exists {
			hybridFeatures[feature] = true
		}
	}

	for feature, found := range hybridFeatures {
		if !found {
			t.Errorf("Expected hybrid feature '%s' not found in sample hybrid world", feature)
		}
	}

	// Test 6: Ensure hybrid world can have a vibe assigned
	err = repo.SetWorldVibe("test-hybrid-world", "calm-clarity")
	if err != nil {
		t.Errorf("Error setting vibe for hybrid world: %v", err)
	}

	hybridVibe, err := repo.GetWorldVibe("test-hybrid-world")
	if err != nil {
		t.Errorf("Error getting hybrid world vibe: %v", err)
	}

	if hybridVibe.ID != "calm-clarity" {
		t.Errorf("Expected hybrid world vibe ID 'calm-clarity', got '%s'", hybridVibe.ID)
	}
}

// TestConcurrency tests thread safety of the repository implementation
func TestConcurrency(t *testing.T) {
	// Create a new repository
	repo := repository.NewRepository()
	
	// Test 1: Concurrent reads of vibes
	t.Run("ConcurrentVibeReads", func(t *testing.T) {
		var wg sync.WaitGroup
		errorCh := make(chan error, 100)
		
		// Launch 10 goroutines to read vibes concurrently
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				// Each goroutine will read all vibes 5 times
				for j := 0; j < 5; j++ {
					vibes := repo.GetAllVibes()
					if len(vibes) < 3 {
						errorCh <- fmt.Errorf("goroutine %d, iteration %d: expected at least 3 vibes, got %d", id, j, len(vibes))
					}
					
					// Also try to get a specific vibe
					_, err := repo.GetVibe("focused-flow")
					if err != nil {
						errorCh <- fmt.Errorf("goroutine %d, iteration %d: error getting vibe: %v", id, j, err)
					}
				}
			}(i)
		}
		
		wg.Wait()
		close(errorCh)
		
		// Check for errors
		for err := range errorCh {
			t.Error(err)
		}
	})
	
	// Test 2: Concurrent reads and writes to vibes
	t.Run("ConcurrentVibeReadsAndWrites", func(t *testing.T) {
		var wg sync.WaitGroup
		errorCh := make(chan error, 100)
		
		// Launch 5 goroutines to read vibes concurrently
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				// Each goroutine will read all vibes 10 times
				for j := 0; j < 10; j++ {
					vibes := repo.GetAllVibes()
					if len(vibes) < 3 {
						errorCh <- fmt.Errorf("reader %d, iteration %d: expected at least 3 vibes, got %d", id, j, len(vibes))
					}
					
					// Small sleep to increase chance of interleaving
					time.Sleep(time.Millisecond * 5)
				}
			}(i)
		}
		
		// Launch 5 goroutines to add and delete vibes concurrently
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				// Each goroutine will add a unique vibe and then delete it
				vibeID := fmt.Sprintf("concurrent-vibe-%d", id)
				
				// Create a vibe
				vibe := models.Vibe{
					ID:          vibeID,
					Name:        fmt.Sprintf("Concurrent Vibe %d", id),
					Description: fmt.Sprintf("A vibe created by goroutine %d", id),
					Energy:      0.5,
					Mood:        "neutral",
					Colors:      []string{"#FFFFFF", "#000000"},
				}
				
				// Add it
				err := repo.AddVibe(vibe)
				if err != nil {
					errorCh <- fmt.Errorf("writer %d: error adding vibe: %v", id, err)
					return
				}
				
				// Small sleep to increase chance of interleaving
				time.Sleep(time.Millisecond * 10)
				
				// Verify it was added
				_, err = repo.GetVibe(vibeID)
				if err != nil {
					errorCh <- fmt.Errorf("writer %d: error getting added vibe: %v", id, err)
					return
				}
				
				// Small sleep to increase chance of interleaving
				time.Sleep(time.Millisecond * 10)
				
				// Delete it
				err = repo.DeleteVibe(vibeID)
				if err != nil {
					errorCh <- fmt.Errorf("writer %d: error deleting vibe: %v", id, err)
				}
			}(i)
		}
		
		wg.Wait()
		close(errorCh)
		
		// Check for errors
		for err := range errorCh {
			t.Error(err)
		}
	})
	
	// Test 3: Concurrent reads and writes to worlds
	t.Run("ConcurrentWorldOperations", func(t *testing.T) {
		var wg sync.WaitGroup
		errorCh := make(chan error, 100)
		
		// Launch 5 goroutines to read worlds concurrently
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				// Each goroutine will read all worlds and specific worlds
				for j := 0; j < 5; j++ {
					worlds := repo.GetAllWorlds()
					if len(worlds) < 3 {
						errorCh <- fmt.Errorf("world reader %d, iteration %d: expected at least 3 worlds, got %d", id, j, len(worlds))
					}
					
					// Get a specific world
					_, err := repo.GetWorld("office-space")
					if err != nil {
						errorCh <- fmt.Errorf("world reader %d, iteration %d: error getting world: %v", id, j, err)
					}
					
					// Small sleep to increase chance of interleaving
					time.Sleep(time.Millisecond * 3)
				}
			}(i)
		}
		
		// Launch 5 goroutines to create and modify worlds concurrently
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				worldID := fmt.Sprintf("concurrent-world-%d", id)
				
				// Create a world
				world := models.World{
					ID:          worldID,
					Name:        fmt.Sprintf("Concurrent World %d", id),
					Description: fmt.Sprintf("A world created by goroutine %d", id),
					Type:        models.WorldTypeVirtual,
					Location:    fmt.Sprintf("https://example.com/world-%d", id),
					CurrentVibe: "calm-clarity", // Use an existing vibe ID
				}
				
				// Add the world
				err := repo.AddWorld(world)
				if err != nil {
					errorCh <- fmt.Errorf("world writer %d: error adding world: %v", id, err)
					return
				}
				
				// Small sleep to increase chance of interleaving
				time.Sleep(time.Millisecond * 5)
				
				// Update the world
				retrievedWorld, err := repo.GetWorld(worldID)
				if err != nil {
					errorCh <- fmt.Errorf("world writer %d: error getting world for update: %v", id, err)
					return
				}
				
				retrievedWorld.Features = []string{fmt.Sprintf("feature-%d-1", id), fmt.Sprintf("feature-%d-2", id)}
				
				err = repo.UpdateWorld(retrievedWorld)
				if err != nil {
					errorCh <- fmt.Errorf("world writer %d: error updating world: %v", id, err)
					return
				}
				
				// Small sleep to increase chance of interleaving
				time.Sleep(time.Millisecond * 5)
				
				// Change the world's vibe
				err = repo.SetWorldVibe(worldID, "energetic-spark")
				if err != nil {
					errorCh <- fmt.Errorf("world writer %d: error changing world vibe: %v", id, err)
					return
				}
				
				// Small sleep to increase chance of interleaving
				time.Sleep(time.Millisecond * 5)
				
				// Delete the world
				err = repo.DeleteWorld(worldID)
				if err != nil {
					errorCh <- fmt.Errorf("world writer %d: error deleting world: %v", id, err)
				}
			}(i)
		}
		
		wg.Wait()
		close(errorCh)
		
		// Check for errors
		for err := range errorCh {
			t.Error(err)
		}
	})
	
	// Test 4: Highly concurrent mixed operations
	t.Run("HighConcurrencyMixedOperations", func(t *testing.T) {
		var wg sync.WaitGroup
		errorCh := make(chan error, 200)
		
		// Create 20 goroutines with random operations
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				// Seed random generator to get different patterns on each run
				rand.Seed(time.Now().UnixNano() + int64(id))
				
				// Each goroutine will perform 10 random operations
				for j := 0; j < 10; j++ {
					// Choose a random operation
					op := rand.Intn(8)
					
					switch op {
					case 0:
						// Get all vibes
						vibes := repo.GetAllVibes()
						if len(vibes) < 3 {
							errorCh <- fmt.Errorf("high concurrency %d: expected at least 3 vibes, got %d", id, len(vibes))
						}
					case 1:
						// Get a specific vibe
						_, err := repo.GetVibe("focused-flow")
						if err != nil {
							errorCh <- fmt.Errorf("high concurrency %d: error getting vibe: %v", id, err)
						}
					case 2:
						// Add a temporary vibe
						tempVibeID := fmt.Sprintf("temp-vibe-%d-%d", id, j)
						tempVibe := models.Vibe{
							ID:          tempVibeID,
							Name:        fmt.Sprintf("Temp Vibe %d-%d", id, j),
							Description: "Temporary vibe for concurrency testing",
							Energy:      0.5,
							Mood:        "neutral",
							Colors:      []string{"#CCCCCC"},
						}
						err := repo.AddVibe(tempVibe)
						if err != nil {
							errorCh <- fmt.Errorf("high concurrency %d: error adding temp vibe: %v", id, err)
						}
						// Immediately delete it to avoid resource buildup
						repo.DeleteVibe(tempVibeID)
					case 3:
						// Get all worlds
						worlds := repo.GetAllWorlds()
						if len(worlds) < 3 {
							errorCh <- fmt.Errorf("high concurrency %d: expected at least 3 worlds, got %d", id, len(worlds))
						}
					case 4:
						// Get a specific world
						_, err := repo.GetWorld("office-space")
						if err != nil {
							errorCh <- fmt.Errorf("high concurrency %d: error getting world: %v", id, err)
						}
					case 5:
						// Try to get a world's vibe
						_, err := repo.GetWorldVibe("office-space")
						if err != nil {
							errorCh <- fmt.Errorf("high concurrency %d: error getting world vibe: %v", id, err)
						}
					case 6:
						// Add and immediately delete a temporary world
						tempWorldID := fmt.Sprintf("temp-world-%d-%d", id, j)
						tempWorld := models.World{
							ID:          tempWorldID,
							Name:        fmt.Sprintf("Temp World %d-%d", id, j),
							Description: "Temporary world for concurrency testing",
							Type:        models.WorldTypeVirtual,
						}
						err := repo.AddWorld(tempWorld)
						if err != nil {
							errorCh <- fmt.Errorf("high concurrency %d: error adding temp world: %v", id, err)
						} else {
							// Immediately delete it to avoid resource buildup
							repo.DeleteWorld(tempWorldID)
						}
					case 7:
						// Try updating an existing world (non-destructively)
						worldIDs := []string{"office-space", "virtual-garden", "hybrid-studio"}
						worldIndex := rand.Intn(len(worldIDs))
						worldID := worldIDs[worldIndex]
						
						world, err := repo.GetWorld(worldID)
						if err != nil {
							errorCh <- fmt.Errorf("high concurrency %d: error getting world for update: %v", id, err)
						} else {
							// Just update the world with the same data (non-destructive operation)
							err = repo.UpdateWorld(world)
							if err != nil {
								errorCh <- fmt.Errorf("high concurrency %d: error updating world: %v", id, err)
							}
						}
					}
					
					// Small sleep to increase chance of interleaving
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)))
				}
			}(i)
		}
		
		wg.Wait()
		close(errorCh)
		
		// Check for errors
		for err := range errorCh {
			t.Error(err)
		}
	})
}

// TestIntegration performs end-to-end tests that combine multiple operations
func TestIntegration(t *testing.T) {
	// Create a new repository
	repo := repository.NewRepository()
	
	// Test 1: Complex vibe lifecycle
	t.Run("ComplexVibeLifecycle", func(t *testing.T) {
		// Step 1: Create a new vibe with sensor data
		temperature := 24.5
		humidity := 55.0
		light := 750.0
		
		complexVibe := models.Vibe{
			ID:          "complex-vibe",
			Name:        "Complex Test Vibe",
			Description: "A vibe for complex lifecycle testing",
			Energy:      0.65,
			Mood:        "creative",
			Colors:      []string{"#FF5733", "#33FF57", "#3357FF"},
			SensorData: models.SensorData{
				Temperature: &temperature,
				Humidity:    &humidity,
				Light:       &light,
			},
		}
		
		err := repo.AddVibe(complexVibe)
		if err != nil {
			t.Fatalf("Error adding complex vibe: %v", err)
		}
		
		// Step 2: Create a new world with this vibe
		complexWorld := models.World{
			ID:          "complex-world",
			Name:        "Complex Test World",
			Description: "A world for complex lifecycle testing",
			Type:        models.WorldTypeHybrid,
			Location:    "Building X + https://complex.example.com",
			CurrentVibe: "complex-vibe",
			Features:    []string{"feature-a", "feature-b"},
		}
		
		err = repo.AddWorld(complexWorld)
		if err != nil {
			t.Fatalf("Error adding complex world: %v", err)
		}
		
		// Step 3: Verify the world has the correct vibe
		worldVibe, err := repo.GetWorldVibe("complex-world")
		if err != nil {
			t.Fatalf("Error getting world vibe: %v", err)
		}
		
		if worldVibe.ID != "complex-vibe" {
			t.Errorf("Expected world vibe ID 'complex-vibe', got '%s'", worldVibe.ID)
		}
		
		// Step 4: Update the vibe's sensor data
		retrievedVibe, err := repo.GetVibe("complex-vibe")
		if err != nil {
			t.Fatalf("Error getting vibe for update: %v", err)
		}
		
		newTemperature := 22.0
		retrievedVibe.SensorData.Temperature = &newTemperature
		
		sound := 45.0
		retrievedVibe.SensorData.Sound = &sound
		
		err = repo.UpdateVibe(retrievedVibe)
		if err != nil {
			t.Fatalf("Error updating vibe: %v", err)
		}
		
		// Step 5: Verify the world still sees the updated vibe
		updatedWorldVibe, err := repo.GetWorldVibe("complex-world")
		if err != nil {
			t.Fatalf("Error getting updated world vibe: %v", err)
		}
		
		if updatedWorldVibe.SensorData.Temperature == nil || *updatedWorldVibe.SensorData.Temperature != newTemperature {
			t.Errorf("Expected updated temperature %f, got %v", newTemperature, updatedWorldVibe.SensorData.Temperature)
		}
		
		if updatedWorldVibe.SensorData.Sound == nil || *updatedWorldVibe.SensorData.Sound != sound {
			t.Errorf("Expected sound %f, got %v", sound, updatedWorldVibe.SensorData.Sound)
		}
		
		// Step 6: Try to delete the vibe (should fail because it's in use)
		err = repo.DeleteVibe("complex-vibe")
		if err != repository.ErrVibeInUse {
			t.Errorf("Expected ErrVibeInUse when deleting vibe in use, got %v", err)
		}
		
		// Step 7: Change the world to use a different vibe
		err = repo.SetWorldVibe("complex-world", "calm-clarity")
		if err != nil {
			t.Fatalf("Error changing world vibe: %v", err)
		}
		
		// Step 8: Now we should be able to delete the original vibe
		err = repo.DeleteVibe("complex-vibe")
		if err != nil {
			t.Errorf("Error deleting vibe after removing from world: %v", err)
		}
		
		// Step 9: Verify the world has the new vibe
		finalWorldVibe, err := repo.GetWorldVibe("complex-world")
		if err != nil {
			t.Fatalf("Error getting final world vibe: %v", err)
		}
		
		if finalWorldVibe.ID != "calm-clarity" {
			t.Errorf("Expected final world vibe ID 'calm-clarity', got '%s'", finalWorldVibe.ID)
		}
		
		// Step 10: Delete the test world
		err = repo.DeleteWorld("complex-world")
		if err != nil {
			t.Errorf("Error deleting test world: %v", err)
		}
	})
	
	// Test 2: Integration test with error conditions
	t.Run("IntegrationWithErrors", func(t *testing.T) {
		// Try creating a world with a non-existent vibe
		invalidWorld := models.World{
			ID:          "invalid-world",
			Name:        "Invalid World",
			Description: "A world with an invalid vibe reference",
			Type:        models.WorldTypePhysical,
			CurrentVibe: "non-existent-vibe",
		}
		
		err := repo.AddWorld(invalidWorld)
		if err != repository.ErrVibeNotFound {
			t.Errorf("Expected ErrVibeNotFound when adding world with non-existent vibe, got %v", err)
		}
		
		// Create a valid vibe and world
		testVibe := models.Vibe{
			ID:          "error-test-vibe",
			Name:        "Error Test Vibe",
			Description: "A vibe for error testing",
			Energy:      0.5,
			Mood:        "neutral",
			Colors:      []string{"#AAAAAA"},
		}
		
		err = repo.AddVibe(testVibe)
		if err != nil {
			t.Fatalf("Error adding test vibe: %v", err)
		}
		
		validWorld := models.World{
			ID:          "error-test-world",
			Name:        "Error Test World",
			Description: "A world for error testing",
			Type:        models.WorldTypePhysical,
			CurrentVibe: "error-test-vibe",
		}
		
		err = repo.AddWorld(validWorld)
		if err != nil {
			t.Fatalf("Error adding valid world: %v", err)
		}
		
		// Try updating the world to use a non-existent vibe
		retrievedWorld, err := repo.GetWorld("error-test-world")
		if err != nil {
			t.Fatalf("Error getting world for update: %v", err)
		}
		
		retrievedWorld.CurrentVibe = "another-non-existent-vibe"
		
		err = repo.UpdateWorld(retrievedWorld)
		if err != repository.ErrVibeNotFound {
			t.Errorf("Expected ErrVibeNotFound when updating world with non-existent vibe, got %v", err)
		}
		
		// Try to set the world's vibe to a non-existent one
		err = repo.SetWorldVibe("error-test-world", "yet-another-non-existent-vibe")
		if err != repository.ErrVibeNotFound {
			t.Errorf("Expected ErrVibeNotFound when setting world vibe to non-existent vibe, got %v", err)
		}
		
		// Try to get a non-existent world's vibe
		_, err = repo.GetWorldVibe("non-existent-world")
		if err != repository.ErrWorldNotFound {
			t.Errorf("Expected ErrWorldNotFound when getting non-existent world's vibe, got %v", err)
		}
		
		// Clean up
		err = repo.DeleteWorld("error-test-world")
		if err != nil {
			t.Errorf("Error deleting test world: %v", err)
		}
		
		err = repo.DeleteVibe("error-test-vibe")
		if err != nil {
			t.Errorf("Error deleting test vibe: %v", err)
		}
	})
	
	// Skip server integration test in repository test - this is covered in separate server tests
	t.Run("ServerIntegration", func(t *testing.T) {
		t.Skip("Server integration is covered in dedicated server tests")
	})
}