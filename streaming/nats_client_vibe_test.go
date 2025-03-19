package streaming

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

// TestPrepareVibeUpdateFullCoverage tests all paths in prepareVibeUpdate
func TestPrepareVibeUpdateFullCoverage(t *testing.T) {
	client := NewNATSClient("nats://test:4222")
	client.streamID = "test-stream"

	// Test case 1: Valid vibe update
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Description: "A test vibe",
		Energy:      0.5,
		Mood:        "calm",
	}

	subject, data, err := client.prepareVibeUpdate("test-world", vibe)
	assert.NoError(t, err)
	assert.Equal(t, "test-stream.world.vibe.test-world", subject)
	assert.NotNil(t, data)

	// Test case 2: Missing world ID
	subject, data, err = client.prepareVibeUpdate("", vibe)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "world ID is required")
	assert.Empty(t, subject)
	assert.Nil(t, data)

	// Test case 3: Nil vibe
	subject, data, err = client.prepareVibeUpdate("test-world", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vibe is required")
	assert.Empty(t, subject)
	assert.Nil(t, data)

	// Test additional error cases
	// Test case 3: Additional error test - bad world ID
	subject, data, err = client.prepareVibeUpdate("", vibe)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "world ID is required")
	assert.Empty(t, subject)
	assert.Nil(t, data)
	
	// Test case 4: Nil vibe test
	subject, data, err = client.prepareVibeUpdate("test-world", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vibe is required")
	assert.Empty(t, subject)
	assert.Nil(t, data)
}