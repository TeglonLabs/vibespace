package test

import (
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/stretchr/testify/assert"
)

// TestRateLimiter tests the token bucket rate limiter
func TestRateLimiter(t *testing.T) {
	// Create a rate limiter with 5 max tokens, refilling 2 tokens every 100ms
	limiter := streaming.NewRateLimiter(5, 2, 100)
	
	// Initially we should have 5 tokens
	for i := 0; i < 5; i++ {
		assert.True(t, limiter.TryAcquire(), "Should acquire token %d", i+1)
	}
	
	// No more tokens available
	assert.False(t, limiter.TryAcquire(), "Should not acquire token after max is used")
	
	// Wait for refill (slightly more than 100ms to be safe)
	time.Sleep(110 * time.Millisecond)
	
	// Should have 2 more tokens now
	assert.True(t, limiter.TryAcquire(), "Should acquire token after refill")
	assert.True(t, limiter.TryAcquire(), "Should acquire second token after refill")
	assert.False(t, limiter.TryAcquire(), "Should not acquire token after refill tokens exhausted")
	
	// Wait for multiple intervals to test multiple refills
	time.Sleep(210 * time.Millisecond) // Just over 2 intervals
	
	// Should have 4 more tokens now (2 per interval × 2 intervals)
	for i := 0; i < 4; i++ {
		assert.True(t, limiter.TryAcquire(), "Should acquire token %d after multiple refills", i+1)
	}
	assert.False(t, limiter.TryAcquire(), "Should not acquire token after multiple refill tokens exhausted")
	
	// Wait long enough to fully refill
	time.Sleep(300 * time.Millisecond) // 3 intervals
	
	// Should be back at max (5)
	for i := 0; i < 5; i++ {
		assert.True(t, limiter.TryAcquire(), "Should acquire token %d after full refill", i+1)
	}
	assert.False(t, limiter.TryAcquire(), "Should not acquire token after full refill tokens exhausted")
}