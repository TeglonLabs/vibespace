package streaming

import (
	"fmt"
	"strings"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
)

// StreamingTools provides MCP tools for controlling streaming
type StreamingTools struct {
	service *StreamingService
}

// NewStreamingTools creates a new streaming tools provider
func NewStreamingTools(service *StreamingService) *StreamingTools {
	return &StreamingTools{
		service: service,
	}
}

// StartStreamingRequest is the request for starting streaming
type StartStreamingRequest struct {
	Interval int `json:"interval"` // Stream interval in milliseconds
}

// StartStreamingResponse is the response for the start streaming request
type StartStreamingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// StartStreaming starts the streaming service
func (t *StreamingTools) StartStreaming(req *StartStreamingRequest) (*StartStreamingResponse, error) {
	// If interval is provided, update the stream interval
	if req.Interval > 0 {
		t.service.config.StreamInterval = time.Duration(req.Interval) * time.Millisecond
	}

	// Start streaming
	err := t.service.StartStreaming()
	if err != nil {
		return &StartStreamingResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to start streaming: %v", err),
		}, nil
	}

	return &StartStreamingResponse{
		Success: true,
		Message: fmt.Sprintf("Streaming started with interval %v", t.service.config.StreamInterval),
	}, nil
}

// StopStreamingResponse is the response for the stop streaming request
type StopStreamingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// StopStreaming stops the streaming service
func (t *StreamingTools) StopStreaming() (*StopStreamingResponse, error) {
	t.service.StopStreaming()

	return &StopStreamingResponse{
		Success: true,
		Message: "Streaming stopped",
	}, nil
}

// StatusResponse is the response for the status request
type StatusResponse struct {
	IsStreaming      bool              `json:"isStreaming"`
	StreamInterval   string            `json:"streamInterval"`
	NATSUrl          string            `json:"natsUrl"`
	Message          string            `json:"message"`
	UIIndicators     struct {
		StreamActive      bool   `json:"streamActive"`
		StreamIndicator   string `json:"streamIndicator"`
		StatusColor       string `json:"statusColor"`
		ConnectionQuality string `json:"connectionQuality"`
	} `json:"uiIndicators"`
	ConnectionStatus ConnectionStatus `json:"connectionStatus"`
}

// Status returns the current status of the streaming service
func (t *StreamingTools) Status() (*StatusResponse, error) {
	isStreaming := t.service.IsStreaming()
	interval := t.service.config.StreamInterval
	natsURL := t.service.config.NATSUrl
	isConnected := t.service.natsClient.IsConnected()
	
	// Get detailed connection status
	connectionStatus := t.service.natsClient.GetConnectionStatus()

	statusMsg := "Streaming inactive"
	if isStreaming {
		statusMsg = fmt.Sprintf("Streaming active with interval %v", interval)
	}

	response := &StatusResponse{
		IsStreaming:      isStreaming,
		StreamInterval:   interval.String(),
		NATSUrl:          natsURL,
		Message:          statusMsg,
		ConnectionStatus: connectionStatus,
	}

	// Set UI indicator values
	response.UIIndicators.StreamActive = isStreaming
	
	// Set stream indicator based on status
	if isStreaming {
		response.UIIndicators.StreamIndicator = "ACTIVE"
		response.UIIndicators.StatusColor = "#4CAF50" // Green
	} else if isConnected {
		response.UIIndicators.StreamIndicator = "READY"
		response.UIIndicators.StatusColor = "#2196F3" // Blue
	} else {
		response.UIIndicators.StreamIndicator = "OFFLINE"
		response.UIIndicators.StatusColor = "#F44336" // Red
	}
	
	// Set connection quality
	if isConnected {
		// Use RTT (round-trip time) to determine connection quality if available
		if connectionStatus.RTT != "" {
			rttStr := connectionStatus.RTT
			if strings.Contains(rttStr, "µs") || // microseconds
			   (strings.Contains(rttStr, "ms") && strings.HasPrefix(rttStr, "0.")) { // < 1ms
				response.UIIndicators.ConnectionQuality = "EXCELLENT"
			} else if strings.Contains(rttStr, "ms") && 
			         !strings.HasPrefix(rttStr, "0.") && 
					 !strings.HasPrefix(rttStr, "1") { // 1-9ms
				response.UIIndicators.ConnectionQuality = "GOOD"
			} else { // 10ms+ or seconds
				response.UIIndicators.ConnectionQuality = "FAIR"
			}
		} else {
			// Fall back to URL-based determination
			if strings.Contains(natsURL, "nonlocal.info") {
				response.UIIndicators.ConnectionQuality = "REMOTE"
			} else if strings.Contains(natsURL, "localhost") {
				response.UIIndicators.ConnectionQuality = "LOCAL"
			} else {
				response.UIIndicators.ConnectionQuality = "CUSTOM"
			}
		}
	} else {
		response.UIIndicators.ConnectionQuality = "DISCONNECTED"
	}

	return response, nil
}

// StreamWorldRequest is the request for streaming a specific world
type StreamWorldRequest struct {
	WorldID string             `json:"worldId"`
	UserID  string             `json:"userId"`
	Sharing *SharingRequest    `json:"sharing,omitempty"`
}

// SharingRequest defines sharing settings for the stream request
type SharingRequest struct {
	IsPublic     bool     `json:"isPublic"`
	AllowedUsers []string `json:"allowedUsers,omitempty"`
	ContextLevel string   `json:"contextLevel,omitempty"`
}

// StreamWorldResponse is the response for the stream world request
type StreamWorldResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// StreamWorld streams a single world moment
func (t *StreamingTools) StreamWorld(req *StreamWorldRequest) (*StreamWorldResponse, error) {
	if req.WorldID == "" {
		return &StreamWorldResponse{
			Success: false,
			Message: "World ID is required",
		}, nil
	}

	// Require user ID for attribution and access control
	if req.UserID == "" {
		return &StreamWorldResponse{
			Success: false,
			Message: "User ID is required",
		}, nil
	}

	// Apply sharing settings if provided
	if req.Sharing != nil {
		// Get the world first to apply sharing settings
		moment, err := t.service.momentGenerator.GenerateMoment(req.WorldID)
		if err != nil {
			return &StreamWorldResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to generate world moment: %v", err),
			}, nil
		}
		
		// Apply the sharing settings
		moment.CreatorID = req.UserID
		
		// Create sharing settings from request
		moment.Sharing = models.SharingSettings{
			IsPublic: req.Sharing.IsPublic,
			AllowedUsers: req.Sharing.AllowedUsers,
		}
		
		// Apply context level if provided
		if req.Sharing.ContextLevel != "" {
			moment.Sharing.ContextLevel = models.ContextLevel(req.Sharing.ContextLevel)
		} else {
			// Default to partial if not specified
			moment.Sharing.ContextLevel = models.ContextLevelPartial
		}
		
		// Add the requesting user to viewers
		userExists := false
		for _, viewer := range moment.Viewers {
			if viewer == req.UserID {
				userExists = true
				break
			}
		}
		
		if !userExists {
			moment.Viewers = append(moment.Viewers, req.UserID)
		}
		
		// Publish directly with user's sharing preferences
		if err := t.service.natsClient.PublishWorldMoment(moment, req.UserID); err != nil {
			return &StreamWorldResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to publish world moment: %v", err),
			}, nil
		}
		
		return &StreamWorldResponse{
			Success: true,
			Message: fmt.Sprintf("Streamed moment for world %s by user %s with custom sharing settings", req.WorldID, req.UserID),
		}, nil
	}

	// If no custom sharing, use default service method
	err := t.service.StreamSingleWorld(req.WorldID, req.UserID)
	if err != nil {
		return &StreamWorldResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to stream world moment: %v", err),
		}, nil
	}

	return &StreamWorldResponse{
		Success: true,
		Message: fmt.Sprintf("Streamed moment for world %s by user %s", req.WorldID, req.UserID),
	}, nil
}

// UpdateConfigRequest is the request for updating streaming configuration
type UpdateConfigRequest struct {
	NATSHost       string `json:"natsHost"`       // NATS host (e.g., "nonlocal.info")
	NATSPort       int    `json:"natsPort"`       // NATS port (default: 4222)
	NATSUrl        string `json:"natsUrl"`        // Complete NATS URL (overrides NATSHost/NATSPort if set)
	StreamID       string `json:"streamId"`       // Stream identifier (default: "ies")
	StreamInterval int    `json:"streamInterval"` // in milliseconds
}

// UpdateConfigResponse is the response for the update config request
type UpdateConfigResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateConfig updates the streaming configuration
func (t *StreamingTools) UpdateConfig(req *UpdateConfigRequest) (*UpdateConfigResponse, error) {
	// Track if we need to create a new client
	needNewClient := false
	
	// Acquire write lock to prevent races with streamMoments
	t.service.mu.Lock()
	wasStreaming := t.service.streamingActive
	
	// Update NATS URL if provided directly
	if req.NATSUrl != "" {
		t.service.config.NATSUrl = req.NATSUrl
		needNewClient = true
	} else {
		// Update host and port independently if provided
		if req.NATSHost != "" {
			t.service.config.NATSHost = req.NATSHost
			needNewClient = true
		}
		
		if req.NATSPort > 0 {
			t.service.config.NATSPort = req.NATSPort
			needNewClient = true
		}
		
		// Reconstruct the URL if host or port changed
		if needNewClient {
			t.service.config.NATSUrl = fmt.Sprintf("nats://%s:%d", 
				t.service.config.NATSHost, 
				t.service.config.NATSPort)
		}
	}
	
	// Update stream ID if provided
	streamIDChanged := false
	if req.StreamID != "" {
		t.service.config.StreamID = req.StreamID
		streamIDChanged = true
	}
	
	// If we need to create a new client due to URL or stream ID changes
	if needNewClient || streamIDChanged {
		// Stop streaming if active (must be called without holding the lock)
		if wasStreaming {
			t.service.stopStreaming()
		}
		
		// Close the old connection
		t.service.natsClient.Close()
		
		// Create a new client with the updated configuration
		t.service.natsClient = NewNATSClientWithStreamID(
			t.service.config.NATSUrl, 
			t.service.config.StreamID)
		
		// Reconnect and resume streaming if needed
		err := t.service.natsClient.Connect()
		if err != nil {
			return &UpdateConfigResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to connect to new NATS server: %v", err),
			}, nil
		}
		
		if wasStreaming {
			if err := t.service.startStreaming(); err != nil {
				t.service.mu.Unlock()
				return &UpdateConfigResponse{
					Success: false,
					Message: fmt.Sprintf("Failed to restart streaming: %v", err),
				}, nil
			}
		}
	}

	// Update stream interval if provided
	if req.StreamInterval > 0 {
		t.service.config.StreamInterval = time.Duration(req.StreamInterval) * time.Millisecond
	}
	
	// Release the lock before returning
	t.service.mu.Unlock()

	return &UpdateConfigResponse{
		Success: true,
		Message: fmt.Sprintf("Configuration updated successfully (NATS: %s, Stream ID: %s)",
			t.service.config.NATSUrl,
			t.service.config.StreamID),
	}, nil
}

// GetStreamingToolMethods returns the available streaming tool methods
func GetStreamingToolMethods() map[string]interface{} {
	return map[string]interface{}{
		"streaming_startStreaming": (*StreamingTools).StartStreaming,
		"streaming_stopStreaming":  (*StreamingTools).StopStreaming,
		"streaming_status":         (*StreamingTools).Status,
		"streaming_streamWorld":    (*StreamingTools).StreamWorld,
		"streaming_updateConfig":   (*StreamingTools).UpdateConfig,
	}
}