package streaming

import (
	"errors"
	"sync"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
)

// MockNATSClient implements a mock version of NATSClientInterface for testing
type MockNATSClient struct {
	url               string
	streamID          string
	connected         bool
	reconnectCount    int
	disconnectCount   int
	lastConnectTime   time.Time
	lastError         error
	publishedMoments  []*models.WorldMoment
	publishedVibes    map[string]*models.Vibe
	connectError      error
	publishMomentError error
	publishVibeError   error
	mu                sync.Mutex
}

// NewMockNATSClient creates a new mock NATS client for testing
func NewMockNATSClient() *MockNATSClient {
	return &MockNATSClient{
		url:              "nats://mock.server:4222",
		streamID:         "test-stream",
		connected:        false,
		publishedMoments: []*models.WorldMoment{},
		publishedVibes:   make(map[string]*models.Vibe),
	}
}

// SetConnectError sets an error to be returned by Connect
func (m *MockNATSClient) SetConnectError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connectError = err
}

// SetPublishMomentError sets an error to be returned by PublishWorldMoment
func (m *MockNATSClient) SetPublishMomentError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishMomentError = err
}

// SetPublishVibeError sets an error to be returned by PublishVibeUpdate
func (m *MockNATSClient) SetPublishVibeError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishVibeError = err
}

// Connect implements the NATSClientInterface.Connect method
func (m *MockNATSClient) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.connectError != nil {
		return m.connectError
	}
	
	m.connected = true
	m.lastConnectTime = time.Now()
	return nil
}

// Close implements the NATSClientInterface.Close method
func (m *MockNATSClient) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = false
}

// PublishWorldMoment implements the NATSClientInterface.PublishWorldMoment method
func (m *MockNATSClient) PublishWorldMoment(moment *models.WorldMoment, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.publishMomentError != nil {
		return m.publishMomentError
	}
	
	if !m.connected {
		return errors.New("not connected to NATS server")
	}
	
	m.publishedMoments = append(m.publishedMoments, moment)
	return nil
}

// PublishVibeUpdate implements the NATSClientInterface.PublishVibeUpdate method
func (m *MockNATSClient) PublishVibeUpdate(worldID string, vibe *models.Vibe) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.publishVibeError != nil {
		return m.publishVibeError
	}
	
	if !m.connected {
		return errors.New("not connected to NATS server")
	}
	
	m.publishedVibes[worldID] = vibe
	return nil
}

// IsConnected implements the NATSClientInterface.IsConnected method
func (m *MockNATSClient) IsConnected() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.connected
}

// GetConnectionStatus implements the NATSClientInterface.GetConnectionStatus method
func (m *MockNATSClient) GetConnectionStatus() ConnectionStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var lastErrorMsg string
	if m.lastError != nil {
		lastErrorMsg = m.lastError.Error()
	}
	
	return ConnectionStatus{
		IsConnected:      m.connected,
		URL:              m.url,
		ReconnectCount:   m.reconnectCount,
		DisconnectCount:  m.disconnectCount,
		LastConnectTime:  m.lastConnectTime,
		LastErrorMessage: lastErrorMsg,
	}
}

// GetPublishedMoments returns all published moments for testing verification
func (m *MockNATSClient) GetPublishedMoments() []*models.WorldMoment {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.publishedMoments
}

// GetPublishedVibes returns all published vibes for testing verification
func (m *MockNATSClient) GetPublishedVibes() map[string]*models.Vibe {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.publishedVibes
}

// SimulateDisconnect simulates a disconnection from the NATS server
func (m *MockNATSClient) SimulateDisconnect() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = false
	m.disconnectCount++
}

// SimulateReconnect simulates a reconnection to the NATS server
func (m *MockNATSClient) SimulateReconnect() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = true
	m.reconnectCount++
}