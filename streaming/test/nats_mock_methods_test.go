package test

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/stretchr/testify/assert"
)

// TestMockNATSClientErrorSetters tests the error setter methods
func TestMockNATSClientErrorSetters(t *testing.T) {
	mockClient := streaming.NewMockNATSClient()
	
	// Set connect error
	expectedError := assert.AnError
	mockClient.SetConnectError(expectedError)
	
	// Connect should return the error
	err := mockClient.Connect()
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	
	// Clear error and try again
	mockClient.SetConnectError(nil)
	err = mockClient.Connect()
	assert.NoError(t, err)
	assert.True(t, mockClient.IsConnected())
	
	// Test PublishMomentError
	mockClient.SetPublishMomentError(expectedError)
	
	// Create a moment (doesn't matter what's in it for this test)
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: 123456789,
	}
	
	// Publish should return the error
	err = mockClient.PublishWorldMoment(moment, "user1")
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	
	// Test PublishVibeError
	mockClient.SetPublishVibeError(expectedError)
	
	// Create a vibe (doesn't matter what's in it for this test)
	vibe := &models.Vibe{
		ID:   "test-vibe",
		Name: "Test Vibe",
	}
	
	// Publish should return the error
	err = mockClient.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

// TestMockNATSClientSimulateMethods tests the simulation methods
func TestMockNATSClientSimulateMethods(t *testing.T) {
	mockClient := streaming.NewMockNATSClient()
	mockClient.Connect()
	assert.True(t, mockClient.IsConnected())
	
	// Test SimulateDisconnect
	mockClient.SimulateDisconnect()
	assert.False(t, mockClient.IsConnected())
	
	// Verify disconnect count
	status := mockClient.GetConnectionStatus()
	assert.Equal(t, 1, status.DisconnectCount)
	
	// Test SimulateReconnect
	mockClient.SimulateReconnect()
	assert.True(t, mockClient.IsConnected())
	
	// Verify reconnect count
	status = mockClient.GetConnectionStatus()
	assert.Equal(t, 1, status.ReconnectCount)
	
	// Simulate multiple disconnects/reconnects
	mockClient.SimulateDisconnect()
	mockClient.SimulateDisconnect() // Should increment the counter
	mockClient.SimulateReconnect()
	mockClient.SimulateReconnect() // Should increment the counter
	
	// Verify counters
	status = mockClient.GetConnectionStatus()
	assert.Equal(t, 3, status.DisconnectCount)
	assert.Equal(t, 3, status.ReconnectCount)
}