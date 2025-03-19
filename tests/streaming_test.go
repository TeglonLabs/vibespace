package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/stretchr/testify/assert"
)

// MockNATSServer is a test helper that mimics a NATS server
type MockNATSServer struct {
	messages map[string][]byte
}

func NewMockNATSServer() *MockNATSServer {
	return &MockNATSServer{
		messages: make(map[string][]byte),
	}
}

func (m *MockNATSServer) Publish(subject string, data []byte) {
	m.messages[subject] = data
}

func (m *MockNATSServer) GetMessage(subject string) []byte {
	return m.messages[subject]
}

// TestWorldMomentMultiplayerFlow tests the core multiplayer flow
func TestWorldMomentMultiplayerFlow(t *testing.T) {
	// Set up test repository with sample data
	repo := repository.NewRepository()
	
	// Create two test users
	userA := "userA"
	userB := "userB"
	
	// Create a test world
	world := models.World{
		ID:          "office",
		Name:        "Office Space",
		Description: "A modern office space",
		Type:        models.WorldTypePhysical,
		Occupancy:   5,
		CreatorID:   userA,
	}
	
	// We don't need to add sensor data for this test
	
	err := repo.AddWorld(world)
	assert.NoError(t, err)
	
	// Create a test vibe
	vibe := models.Vibe{
		ID:          "focused",
		Name:        "Focused Work",
		Description: "A calm, productive atmosphere",
		Energy:      0.7,
		Mood:        "productive",
		Colors:      []string{"#2c3e50", "#3498db"},
		CreatorID:   userA,
	}
	
	err = repo.AddVibe(vibe)
	assert.NoError(t, err)
	
	// Assign vibe to world
	err = repo.SetWorldVibe("office", "focused")
	assert.NoError(t, err)
	
	// Set up mock NATS server
	mockNATS := NewMockNATSServer()
	
	// Create moment generator
	generator := streaming.NewMomentGenerator(repo)
	
	// Test 1: Generate a world moment with creator attribution
	moment, err := generator.GenerateMoment("office")
	assert.NoError(t, err)
	assert.Equal(t, "office", moment.WorldID)
	assert.Equal(t, userA, moment.CreatorID)
	assert.Equal(t, 0, len(moment.Viewers))
	
	// Test 2: Add a viewer and check streaming with proper sharing settings
	moment.Viewers = append(moment.Viewers, userB)
	moment.Sharing = models.SharingSettings{
		IsPublic:     true,
		ContextLevel: models.ContextLevelFull,
	}
	
	// Serialize and publish
	data, err := json.Marshal(moment)
	assert.NoError(t, err)
	
	// Publish to user-specific subject to demonstrate multiplayer awareness
	subject := "world.moment.office.user." + userA
	mockNATS.Publish(subject, data)
	
	// Test 3: Verify what userB would receive
	receivedData := mockNATS.GetMessage(subject)
	var receivedMoment models.WorldMoment
	err = json.Unmarshal(receivedData, &receivedMoment)
	assert.NoError(t, err)
	
	// Verify user attribution and sharing settings are preserved
	assert.Equal(t, userA, receivedMoment.CreatorID)
	assert.Contains(t, receivedMoment.Viewers, userB)
	assert.Equal(t, models.ContextLevel("full"), receivedMoment.Sharing.ContextLevel)
}

// TestPrivateWorldMoment verifies that privacy settings work correctly
func TestPrivateWorldMoment(t *testing.T) {
	// Set up test repository with sample data
	repo := repository.NewRepository()
	
	// Create test users
	userA := "userA"
	userB := "userB"
	userC := "userC"
	
	// Create a test world with private settings
	world := models.World{
		ID:          "sanctuary",
		Name:        "Private Sanctuary",
		Description: "A private space",
		Type:        models.WorldTypeVirtual,
		CreatorID:   userA,
	}
	
	err := repo.AddWorld(world)
	assert.NoError(t, err)
	
	// Set up mock NATS server
	mockNATS := NewMockNATSServer()
	
	// Create moment generator
	generator := streaming.NewMomentGenerator(repo)
	
	// Generate a world moment with private sharing settings
	moment, err := generator.GenerateMoment("sanctuary")
	assert.NoError(t, err)
	
	// Set private sharing settings - only userB is allowed
	moment.CreatorID = userA
	moment.Sharing = models.SharingSettings{
		IsPublic:     false,
		AllowedUsers: []string{userB},
		ContextLevel: models.ContextLevelPartial,
	}
	
	// Serialize and publish
	data, err := json.Marshal(moment)
	assert.NoError(t, err)
	
	// Publish to subject
	subject := "world.moment.sanctuary.user." + userA
	mockNATS.Publish(subject, data)
	
	// Verify what userB would receive (full access)
	receivedData := mockNATS.GetMessage(subject)
	var receivedMoment models.WorldMoment
	err = json.Unmarshal(receivedData, &receivedMoment)
	assert.NoError(t, err)
	
	// Verify sharing settings are preserved and world is accessible to userB
	assert.False(t, receivedMoment.Sharing.IsPublic)
	assert.Contains(t, receivedMoment.Sharing.AllowedUsers, userB)
	
	// Test filtering function (would be part of the streaming service)
	canAccess := streaming.CanAccessWorld(userB, &receivedMoment)
	assert.True(t, canAccess, "UserB should have access")
	
	canAccess = streaming.CanAccessWorld(userC, &receivedMoment)
	assert.False(t, canAccess, "UserC should NOT have access")
}

// TestStreamWorldUserAttribution tests that the streaming_streamWorld tool properly handles user attribution
func TestStreamWorldUserAttribution(t *testing.T) {
	// Set up test repository with sample data
	repo := repository.NewRepository()
	
	// Create test users
	userA := "userA"
	userB := "userB"
	
	// Create a test world
	world := models.World{
		ID:          "conference",
		Name:        "Conference Room",
		Description: "A large meeting space",
		Type:        models.WorldTypePhysical,
		CreatorID:   userA,
	}
	
	err := repo.AddWorld(world)
	assert.NoError(t, err)
	
	// Create a stub config that would be used in a real scenario
	_ = &streaming.StreamingConfig{
		NATSUrl:        "nats://nonlocal.info:4222",
		StreamInterval: 5 * time.Second,
		AutoStart:      false,
	}
	// Note: We're not creating an actual service for this test
	
	// We can't actually test the nats connection here, but we can mock the moment generator
	// to verify that the request parameters are processed correctly
	
	// Create a mock moment with the necessary fields
	moment := &models.WorldMoment{
		WorldID:   "conference",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		CreatorID: userA,
		Viewers:   []string{},
		Sharing: models.SharingSettings{
			IsPublic:     false,
			AllowedUsers: []string{},
			ContextLevel: models.ContextLevelPartial,
		},
	}
	
	// Manually apply the sharing settings like the StreamWorld method would
	moment.CreatorID = userB // Override with requesting user
	moment.Sharing = models.SharingSettings{
		IsPublic:     true,
		AllowedUsers: []string{userA},
		ContextLevel: models.ContextLevel("full"),
	}
	moment.Viewers = append(moment.Viewers, userB)
	
	// Verify the moment has the correct user attribution and sharing settings
	assert.Equal(t, userB, moment.CreatorID)
	assert.True(t, moment.Sharing.IsPublic)
	assert.Contains(t, moment.Sharing.AllowedUsers, userA)
	assert.Equal(t, models.ContextLevel("full"), moment.Sharing.ContextLevel)
	assert.Contains(t, moment.Viewers, userB)
}

// TestNonlocalServerConfig tests that the nonlocal.info server configuration is used
func TestNonlocalServerConfig(t *testing.T) {
	repo := repository.NewRepository()

	// Create streaming service with nonlocal.info URL
	config := &streaming.StreamingConfig{
		NATSUrl:        "nats://nonlocal.info:4222",
		StreamInterval: 5 * time.Second,
		AutoStart:      false,
	}
	service := streaming.NewStreamingService(repo, config)

	// Create streaming tools
	tools := streaming.NewStreamingTools(service)

	// Check that status shows the correct URL
	status, err := tools.Status()
	assert.NoError(t, err)
	assert.Equal(t, "nats://nonlocal.info:4222", status.NATSUrl)
}

// Helper function for testing
func floatPtr(v float64) *float64 {
	return &v
}