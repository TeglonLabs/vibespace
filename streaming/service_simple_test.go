package streaming

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewStreamingService(t *testing.T) {
	repo := repository.NewRepository()
	
	// Test with default values
	config := &StreamingConfig{}
	service := NewStreamingService(repo, config)
	
	assert.NotNil(t, service)
	assert.Equal(t, repo, service.repo)
	assert.NotNil(t, service.natsClient)
	assert.NotNil(t, service.momentGenerator)
	assert.Equal(t, "nats://nonlocal.info:4222", service.config.NATSUrl)
	assert.Equal(t, "ies", service.config.StreamID)
	assert.False(t, service.streamingActive)
	
	// Test with custom values
	customConfig := &StreamingConfig{
		NATSHost:       "custom.host",
		NATSPort:       5222,
		StreamID:       "custom-stream",
		StreamInterval: 10 * time.Second,
		AutoStart:      true,
	}
	
	customService := NewStreamingService(repo, customConfig)
	assert.Equal(t, "nats://custom.host:5222", customService.config.NATSUrl)
	assert.Equal(t, "custom-stream", customService.config.StreamID)
	
	// Test with direct URL
	urlConfig := &StreamingConfig{
		NATSUrl:        "nats://direct.url:4333",
		NATSHost:       "ignored.host", // Should be ignored
		StreamInterval: 10 * time.Second,
	}
	
	urlService := NewStreamingService(repo, urlConfig)
	assert.Equal(t, "nats://direct.url:4333", urlService.config.NATSUrl)
}


func TestIsStreaming(t *testing.T) {
	repo := repository.NewRepository()
	config := &StreamingConfig{}
	service := NewStreamingService(repo, config)
	
	// Initial state should be not streaming
	assert.False(t, service.IsStreaming())
	
	// Set streaming active and check
	service.streamingActive = true
	assert.True(t, service.IsStreaming())
}

func TestNewStreamingTools(t *testing.T) {
	repo := repository.NewRepository()
	config := &StreamingConfig{}
	service := NewStreamingService(repo, config)
	
	tools := NewStreamingTools(service)
	
	assert.NotNil(t, tools)
	assert.Equal(t, service, tools.service)
}

func TestNewRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(10, 2, 100)
	
	assert.NotNil(t, limiter)
	assert.Equal(t, 10, limiter.tokens)
	assert.Equal(t, 2, limiter.refillRate)
	assert.Equal(t, 100, limiter.intervalMs)
}

func TestNewNATSClient(t *testing.T) {
	url := "nats://localhost:4222"
	client := NewNATSClient(url)
	
	assert.NotNil(t, client)
	assert.Equal(t, url, client.url)
	assert.Equal(t, "ies", client.streamID)
	assert.False(t, client.connected)
}

func TestNewNATSClientWithStreamID(t *testing.T) {
	url := "nats://localhost:4222"
	streamID := "test-stream"
	client := NewNATSClientWithStreamID(url, streamID)
	
	assert.NotNil(t, client)
	assert.Equal(t, url, client.url)
	assert.Equal(t, streamID, client.streamID)
	assert.False(t, client.connected)
}