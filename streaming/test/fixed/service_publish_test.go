package fixed

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/bmorphism/vibespace-mcp-go/streaming/test/mocks"
	"github.com/bmorphism/vibespace-mcp-go/streaming/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestServicePublishVibeUpdated tests the PublishVibeUpdate method with proper error mocking
// This test demonstrates a better approach to mocking that properly handles errors
func TestServicePublishVibeUpdated(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	
	// Create mock client that will properly register error expectations
	mockConn := new(MockNatsConnWithForceError)
	mockConn.On("IsConnected").Return(true)
	mockConn.On("Close").Return()

	// Create a mock NATS client
	mockClient := new(MockNATSClient)
	
	// Set up initial state (connected)
	mockClient.On("IsConnected").Return(true).Maybe()
	
	// Set expected behavior for successful publish
	mockClient.On("PublishVibeUpdate", "test-world", mock.Anything).Return(nil).Once()
	
	// Create service with minimal configuration
	config := &streaming.StreamingConfig{
		StreamID: "test-stream",
	}
	service := testutils.CreateMockStreamingService(repo, config, mockClient)
	
	// Create a test vibe
	vibe := &models.Vibe{
		ID:          "test-vibe",
		Name:        "Test Vibe",
		Description: "A vibe for testing",
		Energy:      0.5,
	}
	
	// Test 1: Successful publish
	err := service.PublishVibeUpdate("test-world", vibe)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
	
	// Test 2: Error during publish
	mockClient.On("PublishVibeUpdate", "test-world", mock.Anything).Return(assert.AnError).Once()
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	// Service will add its own error message wrapper
	assert.Contains(t, err.Error(), "assert.AnError")
	mockClient.AssertExpectations(t)
	
	// Test 3: Disconnected state
	// Create a new mock client for the disconnected test
	mockDisconnectedClient := new(MockNATSClient)
	mockDisconnectedClient.On("IsConnected").Return(false)
	mockDisconnectedClient.On("Connect").Return(assert.AnError).Once()
	
	// Replace the client in the service
	service.SetClient(mockDisconnectedClient)
	
	err = service.PublishVibeUpdate("test-world", vibe)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "assert.AnError")
	mockDisconnectedClient.AssertExpectations(t)
}

// MockNatsConnWithForceError is a mock for the NATS connection
type MockNatsConnWithForceError struct {
	mock.Mock
}

// Publish implements the Publish method
func (m *MockNatsConnWithForceError) Publish(subject string, data []byte) error {
	args := m.Called(subject, data)
	return args.Error(0)
}

// IsConnected implements the IsConnected method
func (m *MockNatsConnWithForceError) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

// Close implements the Close method
func (m *MockNatsConnWithForceError) Close() {
	m.Called()
}

// ConnectedServerId implements the ConnectedServerId method
func (m *MockNatsConnWithForceError) ConnectedServerId() string {
	return "test-server"
}

// ConnectedUrl implements the ConnectedUrl method
func (m *MockNatsConnWithForceError) ConnectedUrl() string {
	return "nats://localhost:4222"
}

// RTT implements the RTT method for the mock
func (m *MockNatsConnWithForceError) RTT() (time.Duration, error) {
	return time.Duration(5*time.Millisecond), nil
}

// TestServicePublishWorldMoment tests the StreamSingleWorld method with proper mocking
// This test demonstrates a better approach to testing that correctly handles errors
func TestServicePublishWorldMoment(t *testing.T) {
	// Create test objects
	repo := mocks.NewMockRepository()
	
	// Create a mock NATS client for the successful test case
	mockClient := new(MockNATSClient)
	mockClient.On("IsConnected").Return(true).Maybe()
	mockClient.On("PublishWorldMoment", mock.Anything, "user-123").Return(nil).Once()
	
	// Create service with minimal configuration
	config := &streaming.StreamingConfig{
		StreamID: "test-stream",
	}
	service := testutils.CreateMockStreamingService(repo, config, mockClient)
	
	// Test 1: Successful publish
	err := service.StreamSingleWorld("test-world", "user-123")
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
	
	// Test 2: Error during publish
	mockClient.On("PublishWorldMoment", mock.Anything, "user-123").Return(assert.AnError).Once()
	err = service.StreamSingleWorld("test-world", "user-123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "assert.AnError")
	mockClient.AssertExpectations(t)
	
	// Test 3: Disconnected state
	mockDisconnectedClient := new(MockNATSClient)
	mockDisconnectedClient.On("IsConnected").Return(false)
	mockDisconnectedClient.On("Connect").Return(assert.AnError).Once()
	
	// Replace the client in the service
	service.SetClient(mockDisconnectedClient)
	
	err = service.StreamSingleWorld("test-world", "user-123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "assert.AnError")
	mockDisconnectedClient.AssertExpectations(t)
}

// MockNATSClient implements the NATSClientInterface for testing
type MockNATSClient struct {
	mock.Mock
}

// Connect implements the Connect method
func (m *MockNATSClient) Connect() error {
	args := m.Called()
	return args.Error(0)
}

// Close implements the Close method
func (m *MockNATSClient) Close() {
	m.Called()
}

// PublishWorldMoment implements the PublishWorldMoment method
func (m *MockNATSClient) PublishWorldMoment(moment *models.WorldMoment, userID string) error {
	args := m.Called(moment, userID)
	return args.Error(0)
}

// PublishVibeUpdate implements the PublishVibeUpdate method
func (m *MockNATSClient) PublishVibeUpdate(worldID string, vibe *models.Vibe) error {
	args := m.Called(worldID, vibe)
	return args.Error(0)
}

// IsConnected implements the IsConnected method
func (m *MockNATSClient) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

// GetConnectionStatus implements the GetConnectionStatus method
func (m *MockNATSClient) GetConnectionStatus() streaming.ConnectionStatus {
	args := m.Called()
	if status, ok := args.Get(0).(streaming.ConnectionStatus); ok {
		return status
	}
	return streaming.ConnectionStatus{
		IsConnected: false, 
		URL: "test-url",
	}
}