package streaming

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/bmorphism/vibespace-mcp-go/models"
)

// NatsConnection defines the interface for NATS connection operations we use
// This allows for proper mocking in tests
type NatsConnection interface {
	Publish(subject string, data []byte) error
	IsConnected() bool
	Close()
	ConnectedServerId() string
	ConnectedUrl() string
	RTT() (time.Duration, error)
}

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	tokens      int
	maxTokens   int
	refillRate  int
	lastRefill  time.Time
	intervalMs  int
	mu          sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxTokens, refillRate, intervalMs int) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
		intervalMs: intervalMs,
	}
}

// TryAcquire attempts to acquire a token and returns true if successful
func (r *RateLimiter) TryAcquire() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Refill tokens if needed
	now := time.Now()
	elapsed := now.Sub(r.lastRefill)
	if elapsed >= time.Duration(r.intervalMs)*time.Millisecond {
		// Calculate how many intervals have passed
		intervals := int(elapsed.Milliseconds()) / r.intervalMs
		// Add tokens based on refill rate and intervals
		r.tokens += r.refillRate * intervals
		if r.tokens > r.maxTokens {
			r.tokens = r.maxTokens
		}
		r.lastRefill = now
	}
	
	// Try to acquire a token
	if r.tokens > 0 {
		r.tokens--
		return true
	}
	
	return false
}

// NATSClient handles connections to NATS server and publishing world moments
type NATSClient struct {
	conn            NatsConnection
	url             string
	streamID        string       // Stream identifier (default: "ies")
	connected       bool
	reconnectCount  int
	disconnectCount int
	lastConnectTime time.Time
	lastError       error
	rateLimiter     *RateLimiter
	mu              sync.Mutex
}

// NewNATSClient creates a new NATS client with the specified server URL
func NewNATSClient(url string) *NATSClient {
	// Create a rate limiter allowing 100 messages with 10 refills every 1000ms (1 second)
	// This effectively allows a burst of 100 messages, then 10 messages per second
	rateLimiter := NewRateLimiter(100, 10, 1000)
	
	return &NATSClient{
		url:             url,
		streamID:        "ies",         // Default stream ID
		connected:       false,
		reconnectCount:  0,
		disconnectCount: 0,
		lastConnectTime: time.Time{}, // Zero time
		rateLimiter:     rateLimiter,
	}
}

// NewNATSClientWithStreamID creates a new NATS client with custom stream ID
func NewNATSClientWithStreamID(url string, streamID string) *NATSClient {
	client := NewNATSClient(url)
	client.streamID = streamID
	return client
}

// Connect establishes a connection to the NATS server
func (c *NATSClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected && c.conn != nil && c.conn.IsConnected() {
		return nil
	}

	// Clear any previous connection
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	// Track connection metrics
	c.lastConnectTime = time.Now()

	var err error
	c.conn, err = nats.Connect(c.url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),            // Unlimited reconnect attempts
		nats.ReconnectWait(2*time.Second), // Wait 2 seconds between reconnect attempts
		nats.Timeout(5*time.Second),       // Connect timeout
		nats.PingInterval(20*time.Second), // How often to ping the server to check connection
		nats.MaxPingsOutstanding(5),       // Max number of pings in flight
		
		// Error handlers
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			c.mu.Lock()
			c.lastError = err
			c.mu.Unlock()
			fmt.Printf("NATS error: %v\n", err)
		}),
		
		// Disconnect handler
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			c.mu.Lock()
			c.connected = false
			c.disconnectCount++
			c.lastError = err
			c.mu.Unlock()
			fmt.Printf("NATS disconnected: %v\n", err)
		}),
		
		// Reconnect handler
		nats.ReconnectHandler(func(nc *nats.Conn) {
			c.mu.Lock()
			c.connected = true
			c.reconnectCount++
			c.mu.Unlock()
			fmt.Printf("NATS reconnected to %s (reconnect count: %d)\n", 
				nc.ConnectedUrl(), c.reconnectCount)
		}),
		
		// Closed handler
		nats.ClosedHandler(func(nc *nats.Conn) {
			c.mu.Lock()
			c.connected = false
			c.mu.Unlock()
			fmt.Printf("NATS connection closed\n")
		}),
	)

	if err != nil {
		c.lastError = err
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	c.connected = true
	
	// Log successful connection
	fmt.Printf("Successfully connected to NATS server at %s\n", c.url)
	
	return nil
}

// Close disconnects from the NATS server
func (c *NATSClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
	}
	c.connected = false
}

// prepareWorldMoment prepares a world moment for publishing
// This function is extracted to make it testable without an actual NATS connection
func (c *NATSClient) prepareWorldMoment(moment *models.WorldMoment, userID string) (*models.WorldMoment, error) {
	// If creator is not set, use the provided userID
	if moment.CreatorID == "" {
		moment.CreatorID = userID
	}
	
	// Validate moment data
	if moment.WorldID == "" {
		return nil, fmt.Errorf("world ID is required")
	}
	
	return moment, nil
}

// createMomentSubject creates subject strings for a moment
// This function is extracted to make it testable without an actual NATS connection
func (c *NATSClient) createMomentSubjects(moment *models.WorldMoment) (map[string][]byte, error) {
	subjects := make(map[string][]byte)
	
	// Basic validation of required fields
	if moment.WorldID == "" {
		return nil, fmt.Errorf("world ID is required")
	}
	
	// Create subjects for different access patterns (including stream ID)
	// 1. World-specific subject (for public moments)
	worldSubject := fmt.Sprintf("%s.world.moment.%s", c.streamID, moment.WorldID)
	
	// 2. Creator-specific subject for their worlds
	creatorSubject := fmt.Sprintf("%s.world.moment.%s.user.%s", c.streamID, moment.WorldID, moment.CreatorID)
	
	// Serialize the moment to JSON
	data, err := json.Marshal(moment)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal world moment: %w", err)
	}
	
	// Add main subjects with their data
	if moment.Sharing.IsPublic {
		subjects[worldSubject] = data
	}
	subjects[creatorSubject] = data
	
	// Prepare subject data for each allowed user with filtered content
	for _, allowedUserID := range moment.Sharing.AllowedUsers {
		// Skip if it's the creator (already published)
		if allowedUserID == moment.CreatorID {
			continue
		}
		
		// Get filtered content for this user
		filteredMoment := GetAccessibleContent(allowedUserID, moment)
		if filteredMoment == nil {
			continue
		}
		
		// Create user-specific subject
		userSubject := fmt.Sprintf("%s.world.moment.%s.user.%s", c.streamID, moment.WorldID, allowedUserID)
		
		// Serialize the filtered moment
		filteredData, err := json.Marshal(filteredMoment)
		if err != nil {
			fmt.Printf("Warning: Failed to marshal filtered moment for user %s: %v\n", allowedUserID, err)
			continue
		}
		
		subjects[userSubject] = filteredData
	}
	
	return subjects, nil
}

// PublishWorldMoment publishes a world moment to NATS
func (c *NATSClient) PublishWorldMoment(moment *models.WorldMoment, userID string) error {
	// Check connection first without holding the main lock
	if !c.IsConnected() {
		return fmt.Errorf("not connected to NATS server")
	}
	
	// Apply rate limiting
	if !c.rateLimiter.TryAcquire() {
		return fmt.Errorf("rate limit exceeded, too many messages being published")
	}

	// Prepare the moment (set creator ID, validate)
	preparedMoment, err := c.prepareWorldMoment(moment, userID)
	if err != nil {
		return fmt.Errorf("failed to prepare world moment: %w", err)
	}
	
	// Handle binary data optimization if present
	if preparedMoment.BinaryData != nil && preparedMoment.BinaryData.Encoding == models.EncodingBinary {
		// Make sure the binary data is correctly encoded according to its declared format
		// For raw binary data, this is a no-op, for other encodings we need to ensure it's correct
		if _, err := preparedMoment.GetBinaryData(); err != nil {
			return fmt.Errorf("invalid binary data: %w", err)
		}
	}
	
	// Handle balanced ternary data optimization
	if preparedMoment.BalancedTernaryData != nil {
		// Encode balanced ternary data to a compact binary representation
		ternaryData := TernaryToBytes(*preparedMoment.BalancedTernaryData)
		
		// If the moment doesn't have binary data yet, create it for the ternary data
		if preparedMoment.BinaryData == nil {
			preparedMoment.BinaryData = &models.BinaryData{
				Data:     ternaryData,
				Encoding: models.EncodingBinary,
				Format:   "application/balanced-ternary",
			}
		}
	}
	
	// Create all subject mappings
	subjectData, err := c.createMomentSubjects(preparedMoment)
	if err != nil {
		return fmt.Errorf("failed to create subject mappings: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check connection after acquiring lock
	if !c.connected || c.conn == nil {
		return fmt.Errorf("not connected to NATS server")
	}

	// Publish to all subjects
	publishCount := 0
	for subject, data := range subjectData {
		if err := c.conn.Publish(subject, data); err != nil {
			return fmt.Errorf("failed to publish to subject %s: %w", subject, err)
		}
		publishCount++
	}

	// Log successful publishing
	fmt.Printf("Published world moment for %s to %d subjects\n", moment.WorldID, publishCount)
	
	return nil
}

// prepareVibeUpdate prepares a vibe update for publishing
// This function is extracted to make it testable without an actual NATS connection
func (c *NATSClient) prepareVibeUpdate(worldID string, vibe *models.Vibe) (string, []byte, error) {
	// Validate inputs
	if worldID == "" {
		return "", nil, fmt.Errorf("world ID is required")
	}
	
	if vibe == nil {
		return "", nil, fmt.Errorf("vibe is required")
	}
	
	// Create the subject for vibe updates with stream ID
	subject := fmt.Sprintf("%s.world.vibe.%s", c.streamID, worldID)
	
	// Serialize the vibe to JSON
	data, err := json.Marshal(vibe)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal vibe: %w", err)
	}
	
	return subject, data, nil
}

// PublishVibeUpdate publishes a vibe update to NATS
func (c *NATSClient) PublishVibeUpdate(worldID string, vibe *models.Vibe) error {
	// Check connection first without holding the main lock
	if !c.IsConnected() {
		return fmt.Errorf("not connected to NATS server")
	}
	
	// Apply rate limiting
	if !c.rateLimiter.TryAcquire() {
		return fmt.Errorf("rate limit exceeded, too many messages being published")
	}

	// Prepare the vibe update (validate and serialize)
	subject, data, err := c.prepareVibeUpdate(worldID, vibe)
	if err != nil {
		return fmt.Errorf("failed to prepare vibe update: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check connection after acquiring lock
	if !c.connected || c.conn == nil {
		return fmt.Errorf("not connected to NATS server")
	}
	
	// Publish the data
	err = c.conn.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish vibe update: %w", err)
	}
	
	// Log successful publishing
	fmt.Printf("Published vibe update for world %s\n", worldID)
	
	return nil
}

// IsConnected returns the current connection status
func (c *NATSClient) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected && c.conn != nil && c.conn.IsConnected()
}

// ConnectionStatus returns detailed status information about the NATS connection
type ConnectionStatus struct {
	IsConnected      bool      `json:"isConnected"`
	URL              string    `json:"url"`
	ReconnectCount   int       `json:"reconnectCount"`
	DisconnectCount  int       `json:"disconnectCount"`
	LastConnectTime  time.Time `json:"lastConnectTime"`
	LastErrorMessage string    `json:"lastErrorMessage,omitempty"`
	ServerID         string    `json:"serverId,omitempty"`
	ConnectedURL     string    `json:"connectedUrl,omitempty"`
	RTT              string    `json:"rtt,omitempty"` // Round-trip time
}

// GetConnectionStatus returns detailed status information about the NATS connection
func (c *NATSClient) GetConnectionStatus() ConnectionStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	status := ConnectionStatus{
		IsConnected:     c.connected,
		URL:             c.url,
		ReconnectCount:  c.reconnectCount,
		DisconnectCount: c.disconnectCount,
		LastConnectTime: c.lastConnectTime,
	}
	
	if c.lastError != nil {
		status.LastErrorMessage = c.lastError.Error()
	}
	
	// Add server-specific information if connected
	if c.connected && c.conn != nil {
		status.ServerID = c.conn.ConnectedServerId()
		status.ConnectedURL = c.conn.ConnectedUrl()
		
		// Get RTT (round-trip time) if available
		if rtt, err := c.conn.RTT(); err == nil {
			status.RTT = rtt.String()
		}
	}
	
	return status
}