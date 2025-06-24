package streaming

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNATSClientUtilities tests the utility methods for NATSClient
func TestNATSClientUtilities(t *testing.T) {
	// Create a client for testing
	client := NewNATSClient("nats://test:4222")
	
	// Test SetRateLimiter
	newLimiter := NewRateLimiter(200, 20, 1000)
	client.SetRateLimiter(newLimiter)
	
	// No direct way to test if rate limiter was set correctly
	
	// Test SetConnectedState
	client.SetConnectedState(true)
	
	// Avoid any potential race conditions by direct field access
	client.mu.Lock()
	isConnected := client.connected
	client.mu.Unlock()
	
	assert.True(t, isConnected)
	client.SetConnectedState(false)
	assert.False(t, client.IsConnected())
}

// Custom mock for NATS connection
type customNatsConnMock struct {
	mock.Mock
}

func (m *customNatsConnMock) Publish(subject string, data []byte) error {
	args := m.Called(subject, data)
	return args.Error(0)
}

func (m *customNatsConnMock) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *customNatsConnMock) Close() {
	m.Called()
}

func (m *customNatsConnMock) ConnectedServerId() string {
	args := m.Called()
	return args.String(0)
}

func (m *customNatsConnMock) ConnectedUrl() string {
	args := m.Called()
	return args.String(0)
}

func (m *customNatsConnMock) RTT() (time.Duration, error) {
	args := m.Called()
	return args.Get(0).(time.Duration), args.Error(1)
}

// TestInjectConnection tests the InjectConnection helper method
func TestInjectConnection(t *testing.T) {
	// Create a client for testing
	client := NewNATSClient("nats://test:4222")
	
	// Create a mock NATS connection
	mockConn := new(customNatsConnMock)
	mockConn.On("IsConnected").Return(true)
	mockConn.On("ConnectedServerId").Return("mock-server")
	mockConn.On("ConnectedUrl").Return("nats://mock:4222")
	mockConn.On("RTT").Return(time.Duration(5*time.Millisecond), nil)
	
	// Inject the mock connection
	client.InjectConnection(mockConn)
	
	// Verify the connection was injected
	assert.Equal(t, mockConn, client.conn)
	
	// Test using the injected connection
	status := client.GetConnectionStatus()
	assert.Equal(t, "mock-server", status.ServerID)
	assert.Equal(t, "nats://mock:4222", status.ConnectedURL)
	
	// Verify mock was called
	mockConn.AssertCalled(t, "ConnectedServerId")
	mockConn.AssertCalled(t, "ConnectedUrl")
}

// TestStreamingServiceHelpers tests the helper methods for StreamingService
func TestStreamingServiceHelpers(t *testing.T) {
	// Create a service for testing
	service := &StreamingService{}
	
	// Test GetNATSClient (initially nil)
	assert.Nil(t, service.GetNATSClient())
	
	// Set and test client
	mockClient := NewMockNATSClient()
	service.SetClient(mockClient)
	assert.Equal(t, mockClient, service.GetNATSClient())
	
	// Test SetStreamingActive
	assert.False(t, service.IsStreaming())
	service.SetStreamingActive(true)
	assert.True(t, service.IsStreaming())
	service.SetStreamingActive(false)
	assert.False(t, service.IsStreaming())
}