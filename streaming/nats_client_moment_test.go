package streaming

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

// TestCreateMomentSubjectsFullCoverage tests all paths in createMomentSubjects
func TestCreateMomentSubjectsFullCoverage(t *testing.T) {
	client := NewNATSClient("nats://test:4222")
	client.streamID = "test-stream"

	// Test case 1: Basic moment with no allowed users
	moment1 := &models.WorldMoment{
		WorldID:   "test-world",
		CreatorID: "user1",
		Sharing: models.SharingSettings{
			IsPublic:     true,
			AllowedUsers: []string{},
		},
	}

	subjects1, err := client.createMomentSubjects(moment1)
	assert.NoError(t, err)
	assert.Len(t, subjects1, 2, "Should have 2 subjects: world and creator")
	
	// Verify world subject
	worldSubject := "test-stream.world.moment.test-world"
	assert.Contains(t, subjects1, worldSubject)
	
	// Verify creator subject
	creatorSubject := "test-stream.world.moment.test-world.user.user1"
	assert.Contains(t, subjects1, creatorSubject)
	
	// Test case 2: Moment with allowed users (including creator)
	moment2 := &models.WorldMoment{
		WorldID:   "test-world",
		CreatorID: "user1",
		Sharing: models.SharingSettings{
			IsPublic:     false,
			AllowedUsers: []string{"user1", "user2", "user3"},
		},
	}
	
	subjects2, err := client.createMomentSubjects(moment2)
	assert.NoError(t, err)
	assert.Len(t, subjects2, 3, "Should have 3 subjects: creator and 2 other users")
	
	// Verify user subjects (user1 is creator so already has a subject)
	user2Subject := "test-stream.world.moment.test-world.user.user2"
	user3Subject := "test-stream.world.moment.test-world.user.user3"
	assert.Contains(t, subjects2, creatorSubject)
	assert.Contains(t, subjects2, user2Subject)
	assert.Contains(t, subjects2, user3Subject)
	
	// Test case 3: World ID validation error
	badMoment := &models.WorldMoment{
		WorldID:   "", // Empty world ID should trigger an error
		CreatorID: "user1",
		Sharing: models.SharingSettings{
			IsPublic: true,
		},
	}
	
	_, err = client.createMomentSubjects(badMoment)
	assert.Error(t, err, "Should error on empty world ID")
	assert.Contains(t, err.Error(), "world ID is required")
}