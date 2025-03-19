package fixed

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/bmorphism/vibespace-mcp-go/streaming/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Import the MockNATSClient from the other file
// We don't need to redeclare it

// MockMomentGenerator is a mock for the MomentGeneratorInterface
type MockMomentGenerator struct {
	mock.Mock
}

// GenerateMoment is a mock for the GenerateMoment method
func (m *MockMomentGenerator) GenerateMoment(worldID string) (*models.WorldMoment, error) {
	args := m.Called(worldID)
	if moment, ok := args.Get(0).(*models.WorldMoment); ok {
		return moment, args.Error(1)
	}
	return nil, args.Error(1)
}

// GenerateAllMoments is a mock for the GenerateAllMoments method
func (m *MockMomentGenerator) GenerateAllMoments() ([]*models.WorldMoment, error) {
	args := m.Called()
	if moments, ok := args.Get(0).([]*models.WorldMoment); ok {
		return moments, args.Error(1)
	}
	return nil, args.Error(1)
}

// TestStreamServiceSafeShutdown verifies that the service can be safely shut down after starting
func TestStreamServiceSafeShutdown(t *testing.T) {
	// Create NATS client
	mockClient := new(MockNATSClient)
	mockClient.On("IsConnected").Return(true).Maybe()
	mockClient.On("Connect").Return(nil).Maybe()
	mockClient.On("Close").Return().Maybe()
	mockClient.On("PublishWorldMoment", mock.Anything, mock.Anything).Return(nil).Maybe()
	
	// Create moment generator that returns test data
	mockGenerator := new(MockMomentGenerator)
	testMoments := []*models.WorldMoment{
		{
			WorldID: "test-world-1",
			Timestamp: time.Now().Unix(),
			CreatorID: "system",
			Sharing: models.SharingSettings{
				IsPublic: true,
			},
		},
		{
			WorldID: "test-world-2",
			Timestamp: time.Now().Unix(),
			CreatorID: "system",
			Sharing: models.SharingSettings{
				IsPublic: true,
			},
		},
	}
	mockGenerator.On("GenerateAllMoments").Return(testMoments, nil).Maybe()
	
	// Create service with 50ms stream interval
	config := &streaming.StreamingConfig{
		StreamID: "test-stream",
		StreamInterval: 50 * time.Millisecond,
		AutoStart: false,
	}
	
	// Create repository
	repo := mocks.NewMockRepository()
	
	// Create service
	service := streaming.CreateStreamingService(repo, config, mockClient)
	
	// Replace with our mock generator
	service.SetMomentGenerator(mockGenerator)
	
	// Start streaming and verify it's running
	err := service.StartStreaming()
	assert.NoError(t, err)
	assert.True(t, service.IsStreaming())
	
	// Give it time to do some streaming (at least one iteration)
	time.Sleep(100 * time.Millisecond)
	
	// Stop streaming and verify it's stopped
	service.StopStreaming()
	assert.False(t, service.IsStreaming())
	
	// Verify the mock was called
	mockClient.AssertCalled(t, "IsConnected")
	mockGenerator.AssertCalled(t, "GenerateAllMoments")
}