package test

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/stretchr/testify/assert"
)

func TestMockNATSClientConnect(t *testing.T) {
	client := streaming.NewMockNATSClient()
	err := client.Connect()
	assert.NoError(t, err)
	assert.True(t, client.IsConnected())
	
	// Test with error
	expectedErr := assert.AnError
	client.SetConnectError(expectedErr)
	err = client.Connect()
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestMockNATSClientClose(t *testing.T) {
	client := streaming.NewMockNATSClient()
	client.Connect()
	assert.True(t, client.IsConnected())
	
	client.Close()
	assert.False(t, client.IsConnected())
}

func TestMockNATSClientPublishWorldMoment(t *testing.T) {
	client := streaming.NewMockNATSClient()
	
	// Test error when not connected
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: time.Now().Unix(),
	}
	err := client.PublishWorldMoment(moment, "user1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
	
	// Test successful publish
	client.Connect()
	err = client.PublishWorldMoment(moment, "user1")
	assert.NoError(t, err)
	
	// Verify moment was published
	publishedMoments := client.GetPublishedMoments()
	assert.Len(t, publishedMoments, 1)
	assert.Equal(t, moment, publishedMoments[0])
	
	// Test with error
	expectedErr := assert.AnError
	client.SetPublishMomentError(expectedErr)
	err = client.PublishWorldMoment(moment, "user1")
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestMockNATSClientPublishVibeUpdate(t *testing.T) {
	client := streaming.NewMockNATSClient()
	
	// Test error when not connected
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "happiness",
		Description: "Test vibe",
		Energy:      0.8,
		Mood:        "happy",
	}
	err := client.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
	
	// Test successful publish
	client.Connect()
	err = client.PublishVibeUpdate("test-world", vibe)
	assert.NoError(t, err)
	
	// Verify vibe was published
	publishedVibes := client.GetPublishedVibes()
	assert.Contains(t, publishedVibes, "test-world")
	assert.Equal(t, vibe, publishedVibes["test-world"])
	
	// Test with error
	expectedErr := assert.AnError
	client.SetPublishVibeError(expectedErr)
	err = client.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestMockNATSClientGetConnectionStatus(t *testing.T) {
	client := streaming.NewMockNATSClient()
	
	// Not connected
	status := client.GetConnectionStatus()
	assert.False(t, status.IsConnected)
	
	// Connected
	client.Connect()
	status = client.GetConnectionStatus()
	assert.True(t, status.IsConnected)
	
	// After disconnect
	client.SimulateDisconnect()
	status = client.GetConnectionStatus()
	assert.False(t, status.IsConnected)
	assert.Equal(t, 1, status.DisconnectCount)
	
	// After reconnect
	client.SimulateReconnect()
	status = client.GetConnectionStatus()
	assert.True(t, status.IsConnected)
	assert.Equal(t, 1, status.ReconnectCount)
}

func TestMockNATSClientSimulateDisconnectReconnect(t *testing.T) {
	client := streaming.NewMockNATSClient()
	client.Connect()
	assert.True(t, client.IsConnected())
	
	client.SimulateDisconnect()
	assert.False(t, client.IsConnected())
	
	client.SimulateReconnect()
	assert.True(t, client.IsConnected())
	
	status := client.GetConnectionStatus()
	assert.Equal(t, 1, status.DisconnectCount)
	assert.Equal(t, 1, status.ReconnectCount)
}