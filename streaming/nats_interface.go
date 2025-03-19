package streaming

import (
	"github.com/bmorphism/vibespace-mcp-go/models"
)

// NATSClientInterface defines the interface for NATS client operations
// This allows for proper mocking in tests
type NATSClientInterface interface {
	// Connect establishes a connection to the NATS server
	Connect() error
	
	// Close disconnects from the NATS server
	Close()
	
	// PublishWorldMoment publishes a world moment to NATS
	PublishWorldMoment(moment *models.WorldMoment, userID string) error
	
	// PublishVibeUpdate publishes a vibe update to NATS
	PublishVibeUpdate(worldID string, vibe *models.Vibe) error
	
	// IsConnected returns the current connection status
	IsConnected() bool
	
	// GetConnectionStatus returns detailed status information about the NATS connection
	GetConnectionStatus() ConnectionStatus
}

// Ensure NATSClient implements the interface
var _ NATSClientInterface = (*NATSClient)(nil)