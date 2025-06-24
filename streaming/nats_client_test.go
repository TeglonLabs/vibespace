package streaming

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

// mockNatsConn is a minimal mock for NatsConnection interface
type mockNatsConn struct {
	connected bool
	closed    bool
}

func (m *mockNatsConn) Publish(subject string, data []byte) error {
	return nil
}

func (m *mockNatsConn) IsConnected() bool {
	return m.connected
}

func (m *mockNatsConn) Close() {
	m.closed = true
	m.connected = false
}

func (m *mockNatsConn) ConnectedServerId() string {
	return "test-server"
}

func (m *mockNatsConn) ConnectedUrl() string {
	return "nats://test:4222"
}

func (m *mockNatsConn) RTT() (time.Duration, error) {
	return time.Duration(1 * time.Millisecond), nil
}

func TestNATSClientBasics(t *testing.T) {
	// Test client creation
	client := NewNATSClient("nats://localhost:4222")
	assert.NotNil(t, client)
	assert.False(t, client.IsConnected()) // Initially not connected
	
	// Test URL setting
	assert.Equal(t, "nats://localhost:4222", client.url)
	
	// Test that client can be created with different URLs
	client2 := NewNATSClient("nats://other:5222")
	assert.Equal(t, "nats://other:5222", client2.url)
}

func TestNATSClientClose(t *testing.T) {
	tests := []struct {
		name     string
		hasConn  bool
		expected bool
	}{
		{
			name:     "close with valid connection",
			hasConn:  true,
			expected: true,
		},
		{
			name:     "close with no connection",
			hasConn:  false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewNATSClient("nats://localhost:4222")
			
			if tt.hasConn {
				mockConn := &mockNatsConn{connected: true}
				client.conn = mockConn
				client.connected = true
				
				client.Close()
				assert.True(t, mockConn.closed, "Connection's Close method should be called")
				assert.False(t, client.IsConnected(), "Client should report as not connected")
			} else {
				client.Close() // Should not panic
				assert.False(t, client.IsConnected())
			}
		})
	}
}

func TestNATSClientPublishWorldMoment(t *testing.T) {
	tests := []struct {
		name        string
		connected   bool
		moment      *models.WorldMoment
		userID      string
		expectError bool
	}{
		{
			name:      "successful publish",
			connected: true,
			moment: &models.WorldMoment{
				WorldID:   "test-world",
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
				Sharing: models.SharingSettings{IsPublic: true},
			},
			userID:      "user123",
			expectError: false,
		},
		{
			name:      "publish when not connected",
			connected: false,
			moment: &models.WorldMoment{
				WorldID:   "test-world",
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
				Sharing: models.SharingSettings{IsPublic: true},
			},
			userID:      "user123",
			expectError: true,
		},
		{
			name:      "empty world ID",
			connected: true,
			moment: &models.WorldMoment{
				WorldID:   "", // Invalid
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
				Sharing: models.SharingSettings{IsPublic: true},
			},
			userID:      "user123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewNATSClient("nats://localhost:4222")
			
			if tt.connected {
				mockConn := &mockNatsConn{connected: true}
				client.conn = mockConn
				client.connected = true
			}

			err := client.PublishWorldMoment(tt.moment, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNATSClientIsConnected(t *testing.T) {
	client := NewNATSClient("nats://localhost:4222")
	
	// Initially not connected
	assert.False(t, client.IsConnected())
	
	// Simulate connection
	mockConn := &mockNatsConn{connected: true}
	client.conn = mockConn
	client.connected = true
	
	assert.True(t, client.IsConnected())
	
	// Simulate disconnection
	client.connected = false
	assert.False(t, client.IsConnected())
}

func TestNATSClientConcurrency(t *testing.T) {
	client := NewNATSClient("nats://localhost:4222")
	mockConn := &mockNatsConn{connected: true}
	client.conn = mockConn
	client.connected = true

	// Test concurrent moment publishes don't cause race conditions
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			moment := &models.WorldMoment{
				WorldID:   "test-world",
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
				Sharing: models.SharingSettings{IsPublic: true},
			}
			// Just test that concurrent access doesn't panic
			// Don't assert NoError since we might hit rate limits
			_ = client.PublishWorldMoment(moment, "user123")
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Keep the higher-level behavioral tests for mock client
func TestNATSClientHelpers(t *testing.T) {
	// Test mock client behavior for WorldMoment
	client := NewMockNATSClient()
	client.url = "nats://localhost:4222"
	err := client.Connect()
	assert.NoError(t, err)
	
	// Create a test moment
	tempVal := 22.5
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		SensorData: models.SensorData{
			Temperature: &tempVal,
		},
		CreatorID: "user123",
		Sharing: models.SharingSettings{
			IsPublic: true,
		},
	}
	
	// Publish the moment
	err = client.PublishWorldMoment(moment, "user123")
	assert.NoError(t, err)
	
	// Verify it was published
	publishedMoments := client.GetPublishedMoments()
	assert.Equal(t, 1, len(publishedMoments))
	assert.Equal(t, "test-world", publishedMoments[0].WorldID)
	
	// Test vibe update
	vibe := &models.Vibe{
		ID:      "vibe1",
		Name:    "Test Vibe",
		Energy:  0.8,
		Mood:    models.MoodEnergetic,
	}
	
	err = client.PublishVibeUpdate("test-world", vibe)
	assert.NoError(t, err)
	
	publishedVibes := client.GetPublishedVibes()
	assert.Equal(t, 1, len(publishedVibes))
	assert.Equal(t, "vibe1", publishedVibes["test-world"].ID)
}
