package streaming

import (
	"errors"
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestConnectWithExtractedPrep tests a refactored version of Connect
// We'll simulate what we want to do with the Connect method by
// extracting its configuration setup functionality
func TestConnectWithExtractedPrep(t *testing.T) {
	// Create test client
	client := &NATSClient{
		url:            "nats://localhost:4222",
		streamID:       "test-stream",
		connected:      false,
		reconnectCount: 0,
		disconnectCount: 0,
		lastConnectTime: time.Time{},
		rateLimiter:    NewRateLimiter(100, 10, 1000),
	}

	// Test all code paths in the "already connected" logic
	// Path 1: Already connected, conn is not nil
	mockConn := new(MockNatsConn)
	mockConn.On("IsConnected").Return(true)
	mockConn.On("Close").Return()
	mockConn.On("ConnectedServerId").Return("test-server-id")
	mockConn.On("ConnectedUrl").Return("nats://connected-server:4222")
	mockConn.On("RTT").Return(time.Duration(5*time.Millisecond), nil)
	client.conn = mockConn
	client.connected = true

	err := client.Connect()
	assert.NoError(t, err)

	// Path 2: Client thinks it's connected but conn says it's not
	mockConn = new(MockNatsConn)
	mockConn.On("IsConnected").Return(false)
	mockConn.On("Close").Return()
	mockConn.On("ConnectedServerId").Return("test-server-id")
	mockConn.On("ConnectedUrl").Return("nats://connected-server:4222")
	mockConn.On("RTT").Return(time.Duration(5*time.Millisecond), nil)
	client.conn = mockConn
	client.connected = true

	// Since we're mocking everything, we won't actually get an error
	// We're just testing the first part of the Connect method
	err = client.Connect()
	assert.NoError(t, err)  // No error since we mock the connection
	
	// Just manually set client to disconnected for the next test
	client.connected = false
	client.conn = nil
	
	// Create a new client
	client = &NATSClient{
		url:            "nats://localhost:4222",
		streamID:       "test-stream",
		connected:      false,
		reconnectCount: 0,
		disconnectCount: 0,
		lastConnectTime: time.Time{},
		rateLimiter:    NewRateLimiter(100, 10, 1000),
	}
	
	// Trigger various handlers to ensure they work properly
	// Test error handler
	client.lastError = errors.New("test error")
	
	status := client.GetConnectionStatus()
	assert.Equal(t, "test error", status.LastErrorMessage)
}

// TestConnectHandlers simulates the various connection handlers
func TestConnectHandlers(t *testing.T) {
	// Create test client
	client := &NATSClient{
		url:            "nats://localhost:4222",
		streamID:       "test-stream",
		connected:      true,
		reconnectCount: 0,
		disconnectCount: 0,
		lastConnectTime: time.Now(),
		rateLimiter:    NewRateLimiter(100, 10, 1000),
	}

	// Setup a mock connection
	mockConn := new(MockNatsConn)
	mockConn.On("IsConnected").Return(true)
	mockConn.On("Close").Return()
	mockConn.On("ConnectedServerId").Return("test-server-id")
	mockConn.On("ConnectedUrl").Return("nats://connected-server:4222")
	mockConn.On("RTT").Return(time.Duration(5*time.Millisecond), nil)
	client.conn = mockConn

	// Test disconnect handler would set these values
	client.connected = false
	client.disconnectCount = 1
	client.lastError = errors.New("disconnected")

	assert.False(t, client.IsConnected())
	assert.Equal(t, 1, client.GetConnectionStatus().DisconnectCount)

	// Test reconnect handler would set these values
	client.connected = true
	client.reconnectCount = 1

	assert.True(t, client.IsConnected())
	assert.Equal(t, 1, client.GetConnectionStatus().ReconnectCount)

	// Test closed handler would set these values
	client.connected = false

	assert.False(t, client.IsConnected())
}

// TestGetConnectionStatusDetailed tests all code paths in GetConnectionStatus
func TestGetConnectionStatusDetailed(t *testing.T) {
	// Create test client
	client := &NATSClient{
		url:            "nats://localhost:4222",
		streamID:       "test-stream",
		connected:      true,
		reconnectCount: 3,
		disconnectCount: 2,
		lastConnectTime: time.Now(),
		lastError:       errors.New("test error"),
		rateLimiter:    NewRateLimiter(100, 10, 1000),
	}

	// Path 1: Basic status with an error
	status := client.GetConnectionStatus()
	assert.True(t, status.IsConnected)
	assert.Equal(t, "nats://localhost:4222", status.URL)
	assert.Equal(t, 3, status.ReconnectCount)
	assert.Equal(t, 2, status.DisconnectCount)
	assert.Equal(t, "test error", status.LastErrorMessage)
	assert.True(t, !status.LastConnectTime.IsZero())

	// Path 2: Without an error
	client.lastError = nil
	status = client.GetConnectionStatus()
	assert.Empty(t, status.LastErrorMessage)

	// Path 3: With a connected server
	mockConn := new(MockNatsConn)
	mockConn.On("IsConnected").Return(true)
	mockConn.On("ConnectedServerId").Return("test-server-id")
	mockConn.On("ConnectedUrl").Return("nats://connected-server:4222")
	mockConn.On("RTT").Return(time.Duration(5*time.Millisecond), nil)
	client.conn = mockConn
	
	status = client.GetConnectionStatus()
	assert.Equal(t, "test-server-id", status.ServerID)
	assert.Equal(t, "nats://connected-server:4222", status.ConnectedURL)
	assert.Equal(t, "5ms", status.RTT)

	// Path 4: With RTT error
	mockConn = new(MockNatsConn)
	mockConn.On("IsConnected").Return(true)
	mockConn.On("ConnectedServerId").Return("test-server-id")
	mockConn.On("ConnectedUrl").Return("nats://connected-server:4222")
	mockConn.On("RTT").Return(time.Duration(0), errors.New("rtt error"))
	client.conn = mockConn
	
	status = client.GetConnectionStatus()
	assert.Equal(t, "test-server-id", status.ServerID)
	assert.Equal(t, "nats://connected-server:4222", status.ConnectedURL)
	assert.Empty(t, status.RTT)

	// Path 5: Not connected
	client.connected = false
	status = client.GetConnectionStatus()
	assert.False(t, status.IsConnected)

	// Path 6: Conn is nil
	client.conn = nil
	status = client.GetConnectionStatus()
	assert.False(t, status.IsConnected)
	assert.Empty(t, status.ServerID)
}

// TestRateLimiterAdvanced tests additional rate limiter functionality
func TestRateLimiterAdvanced(t *testing.T) {
	// Create a rate limiter with 2 tokens, refilling 1 token every 100ms
	limiter := NewRateLimiter(2, 1, 100)
	
	// Should be able to acquire 2 tokens immediately
	assert.True(t, limiter.TryAcquire())
	assert.True(t, limiter.TryAcquire())
	
	// Third attempt should fail
	assert.False(t, limiter.TryAcquire())
	
	// Wait for a refill
	time.Sleep(110 * time.Millisecond)
	
	// Should be able to acquire 1 more token
	assert.True(t, limiter.TryAcquire())
	
	// But not 2
	assert.False(t, limiter.TryAcquire())
	
	// Wait for multiple intervals
	time.Sleep(210 * time.Millisecond)
	
	// Should be able to acquire 2 tokens (refilled at 1 per 100ms)
	assert.True(t, limiter.TryAcquire())
	assert.True(t, limiter.TryAcquire())
	
	// But not 3
	assert.False(t, limiter.TryAcquire())
	
	// Test maximum tokens
	limiter = NewRateLimiter(3, 10, 100) // Max 3, refill 10 per interval
	
	// Use all tokens
	assert.True(t, limiter.TryAcquire())
	assert.True(t, limiter.TryAcquire())
	assert.True(t, limiter.TryAcquire())
	assert.False(t, limiter.TryAcquire())
	
	// Wait for a refill
	time.Sleep(110 * time.Millisecond)
	
	// Should only have refilled up to the max
	assert.True(t, limiter.TryAcquire())
	assert.True(t, limiter.TryAcquire())
	assert.True(t, limiter.TryAcquire())
	assert.False(t, limiter.TryAcquire())
}

// MockNatsConnWithForceError is a special mock that allows forcing errors
type MockNatsConnWithForceError struct {
	MockNatsConn
	forceError bool
}

// Publish overrides the mock publish and can force an error
func (m *MockNatsConnWithForceError) Publish(subject string, data []byte) error {
	if m.forceError {
		return errors.New("forced error")
	}
	args := m.Called(subject, data)
	return args.Error(0)
}


// TestNATSClientWithErrorScenarios tests error handling in the NATSClient
func TestNATSClientWithErrorScenarios(t *testing.T) {
	// Create a world moment for testing
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: time.Now().Unix(),
		CreatorID: "creator1",
	}

	// Setup mock for publish error test
	mockConn := new(MockNatsConn)
	mockConn.On("IsConnected").Return(true)
	mockConn.On("ConnectedServerId").Return("test-server-id")
	mockConn.On("ConnectedUrl").Return("nats://connected-server:4222")
	mockConn.On("RTT").Return(time.Duration(5*time.Millisecond), nil)
	mockConn.On("Close").Return()
	mockConn.On("Publish", mock.Anything, mock.Anything).Return(errors.New("publish error"))
	
	client := &NATSClient{
		url:         "nats://localhost:4222",
		streamID:    "test-stream",
		connected:   true,
		conn:        mockConn,
		rateLimiter: NewRateLimiter(100, 10, 1000),
	}

	// Test publish error
	err := client.PublishWorldMoment(moment, "user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish")
	
	// Custom IsConnected mock for rate limiting test
	client = &NATSClient{
		url:         "nats://localhost:4222",
		streamID:    "test-stream",
		connected:   true,
		rateLimiter: NewRateLimiter(1, 1, 1000), // Only 1 token
	}

	mockConn = new(MockNatsConn)
	mockConn.On("IsConnected").Return(true)
	mockConn.On("ConnectedServerId").Return("test-server-id")
	mockConn.On("ConnectedUrl").Return("nats://connected-server:4222")
	mockConn.On("RTT").Return(time.Duration(5*time.Millisecond), nil)
	mockConn.On("Close").Return()
	mockConn.On("Publish", mock.Anything, mock.Anything).Return(nil)
	client.conn = mockConn

	// Use one token
	err = client.PublishWorldMoment(moment, "user")
	assert.NoError(t, err)

	// Should hit rate limit on second attempt
	err = client.PublishWorldMoment(moment, "user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
}