package streaming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMockSimulateFunctions tests the simulation methods in MockNATSClient
func TestMockSimulateFunctions(t *testing.T) {
	// Create a mock client
	mockClient := NewMockNATSClient()
	
	// Test initial state
	assert.False(t, mockClient.IsConnected())
	status := mockClient.GetConnectionStatus()
	assert.Equal(t, 0, status.ReconnectCount)
	assert.Equal(t, 0, status.DisconnectCount)
	
	// Test connect
	mockClient.Connect()
	assert.True(t, mockClient.IsConnected())
	
	// Test SimulateDisconnect
	mockClient.SimulateDisconnect()
	assert.False(t, mockClient.IsConnected())
	
	// Check disconnect count was incremented
	status = mockClient.GetConnectionStatus()
	assert.Equal(t, 0, status.ReconnectCount)
	assert.Equal(t, 1, status.DisconnectCount)
	
	// Test SimulateReconnect
	mockClient.SimulateReconnect()
	assert.True(t, mockClient.IsConnected())
	
	// Check reconnect count was incremented
	status = mockClient.GetConnectionStatus()
	assert.Equal(t, 1, status.ReconnectCount)
	assert.Equal(t, 1, status.DisconnectCount)
	
	// Test multiple disconnects/reconnects
	mockClient.SimulateDisconnect()
	mockClient.SimulateDisconnect() // Should still increment even if already disconnected
	mockClient.SimulateReconnect()
	mockClient.SimulateReconnect() // Should still increment even if already connected
	
	// Check final counts
	status = mockClient.GetConnectionStatus()
	assert.Equal(t, 3, status.ReconnectCount)
	assert.Equal(t, 3, status.DisconnectCount)
}