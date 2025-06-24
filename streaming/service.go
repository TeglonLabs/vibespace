package streaming

import (
	"fmt"
	"sync"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
)

// StreamingConfig holds configuration for the streaming service
type StreamingConfig struct {
	NATSHost       string        // NATS host (e.g., "nonlocal.info")
	NATSPort       int           // NATS port (default: 4222)
	NATSUrl        string        // Complete NATS URL (overrides NATSHost/NATSPort if set)
	StreamID       string        // Stream identifier (default: "ies")
	StreamInterval time.Duration // Interval between streaming moments
	AutoStart      bool          // Whether to start streaming automatically
}

// StreamingService manages NATS streaming for world moments
type StreamingService struct {
	natsClient      NATSClientInterface
	momentGenerator MomentGeneratorInterface
	config          *StreamingConfig
	repo            RepositoryInterface
	streamingActive bool
	stopChan        chan struct{}
	mu              sync.RWMutex // Use RWMutex for better read concurrency
	once            sync.Once    // Ensure single initialization
}

// NewStreamingService creates a new streaming service
func NewStreamingService(repo RepositoryInterface, config *StreamingConfig) *StreamingService {
	// Set default values if not provided
	if config.NATSPort == 0 {
		config.NATSPort = 4222 // Default NATS port
	}
	
	if config.StreamID == "" {
		config.StreamID = "ies" // Default stream ID
	}
	
	// If NATSUrl is not provided, construct it from host and port
	if config.NATSUrl == "" && config.NATSHost != "" {
		config.NATSUrl = fmt.Sprintf("nats://%s:%d", config.NATSHost, config.NATSPort)
	} else if config.NATSUrl == "" {
		// Default to nonlocal.info if nothing is provided
		config.NATSUrl = fmt.Sprintf("nats://nonlocal.info:%d", config.NATSPort)
	}
	
	// Create NATS client with the configured stream ID
	natsClient := NewNATSClientWithStreamID(config.NATSUrl, config.StreamID)
	
	return CreateStreamingService(repo, config, natsClient)
}

// CreateStreamingService creates a new streaming service with a custom NATS client
// This allows dependency injection for testing
func CreateStreamingService(repo RepositoryInterface, config *StreamingConfig, natsClient NATSClientInterface) *StreamingService {
	return &StreamingService{
		natsClient:      natsClient,
		momentGenerator: NewMomentGenerator(repo),
		config:          config,
		repo:            repo,
		streamingActive: false,
		stopChan:        make(chan struct{}),
	}
}

// Start initializes the streaming service and begins streaming if autoStart is true
func (s *StreamingService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Connect to NATS
	if err := s.natsClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Start streaming if autoStart is enabled
	if s.config.AutoStart {
		return s.startStreaming()
	}

	return nil
}

// Stop terminates the streaming service
func (s *StreamingService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Stop streaming if active
	if s.streamingActive {
		s.stopStreaming()
	}

	// Close NATS connection
	s.natsClient.Close()
}

// StartStreaming begins the streaming of world moments
func (s *StreamingService) StartStreaming() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.startStreaming()
}

// startStreaming is the internal method to start streaming (not thread-safe)
func (s *StreamingService) startStreaming() error {
	if s.streamingActive {
		return nil // Already streaming
	}

	// Make sure we're connected to NATS
	if !s.natsClient.IsConnected() {
		if err := s.natsClient.Connect(); err != nil {
			return fmt.Errorf("failed to connect to NATS: %w", err)
		}
	}

	// Reset the stop channel
	s.stopChan = make(chan struct{})
	s.streamingActive = true

	// Start the streaming goroutine
	go s.streamMoments()

	return nil
}

// StopStreaming stops the streaming of world moments
func (s *StreamingService) StopStreaming() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopStreaming()
}

// stopStreaming is the internal method to stop streaming (not thread-safe)
func (s *StreamingService) stopStreaming() {
	if !s.streamingActive {
		return // Not streaming
	}

	// Signal the streaming goroutine to stop
	// Use select to avoid panic if channel is already closed
	select {
	case <-s.stopChan:
		// Already closed
	default:
		close(s.stopChan)
	}
	s.streamingActive = false
}

// streamMoments is the main streaming loop that publishes world moments at regular intervals
func (s *StreamingService) streamMoments() {
	// Get snapshot of configuration and stop channel to avoid races
	s.mu.RLock()
	interval := s.config.StreamInterval
	stopChan := s.stopChan
	momGen := s.momentGenerator
	client := s.natsClient
	s.mu.RUnlock()

	if momGen == nil {
		fmt.Printf("Error: momentGenerator is nil in streamMoments\n")
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Generate and publish moments for all worlds
			moments, err := momGen.GenerateAllMoments()
			if err != nil {
				fmt.Printf("Error generating moments: %v\n", err)
				continue
			}

			// Publish each moment
			for _, moment := range moments {
				// For automatic streaming, we use the "system" as the creator ID
				// if it's not already set
				creatorID := moment.CreatorID
				if creatorID == "" {
					creatorID = "system"
				}
				
				// Set default sharing settings for automated moments if needed
				if !moment.Sharing.IsPublic && len(moment.Sharing.AllowedUsers) == 0 && moment.Sharing.ContextLevel == "" {
					// By default, system-generated moments are public with partial context
					moment.Sharing = models.SharingSettings{
						IsPublic:     true,
						AllowedUsers: []string{},
						ContextLevel: models.ContextLevelPartial,
					}
				}
				
				if err := client.PublishWorldMoment(moment, creatorID); err != nil {
					fmt.Printf("Error publishing moment for world %s: %v\n", moment.WorldID, err)
				}
			}

		case <-stopChan:
			// Streaming has been stopped
			return
		}
	}
}

// StreamSingleWorld generates and streams a moment for a single world
func (s *StreamingService) StreamSingleWorld(worldID string, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Make sure we're connected to NATS
	if !s.natsClient.IsConnected() {
		if err := s.natsClient.Connect(); err != nil {
			return fmt.Errorf("failed to connect to NATS: %w", err)
		}
	}

	// Generate a moment for the world
	moment, err := s.momentGenerator.GenerateMoment(worldID)
	if err != nil {
		return fmt.Errorf("failed to generate moment: %w", err)
	}
	
	// Set creator ID if not already set
	if moment.CreatorID == "" {
		moment.CreatorID = userID
	}
	
	// If world already has viewers, add this user if not already there
	userExists := false
	for _, viewer := range moment.Viewers {
		if viewer == userID {
			userExists = true
			break
		}
	}
	
	if !userExists {
		moment.Viewers = append(moment.Viewers, userID)
	}
	
	// Set default sharing settings if needed (completely empty or only default values)
	if !moment.Sharing.IsPublic && len(moment.Sharing.AllowedUsers) == 0 && moment.Sharing.ContextLevel == "" {
		moment.Sharing = models.SharingSettings{
			IsPublic:     false,
			AllowedUsers: []string{},
			ContextLevel: models.ContextLevelPartial,
		}
	}

	// Publish the moment with user information
	if err := s.natsClient.PublishWorldMoment(moment, userID); err != nil {
		return fmt.Errorf("failed to publish moment: %w", err)
	}

	return nil
}

// IsStreaming returns whether the service is currently streaming
func (s *StreamingService) IsStreaming() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.streamingActive
}

// PublishVibeUpdate publishes a vibe update for a specific world
func (s *StreamingService) PublishVibeUpdate(worldID string, vibe *models.Vibe) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we're connected to NATS first
	if !s.natsClient.IsConnected() {
		return fmt.Errorf("not connected to NATS")
	}

	// Publish the vibe update
	if err := s.natsClient.PublishVibeUpdate(worldID, vibe); err != nil {
		return fmt.Errorf("failed to publish vibe update: %w", err)
	}
	
	return nil
}
