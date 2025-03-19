package streaming

import (
	"testing"
	"sync"
	
	"github.com/stretchr/testify/assert"
)

// TestCloseComplete provides 100% coverage for the Close method
func TestCloseComplete(t *testing.T) {
	// Case 1: With a valid connection
	client := NewNATSClient("nats://localhost:4222")
	mockConn := &mockNatsConn{connected: true}
	client.conn = mockConn
	client.connected = true
	
	// Close should call Close() on the connection
	client.Close()
	assert.True(t, mockConn.closed, "Connection's Close method should be called")
	assert.False(t, client.connected, "Client should be marked as disconnected")
	
	// Case 2: When conn is nil
	client = NewNATSClient("nats://localhost:4222")
	client.conn = nil
	client.connected = true
	
	// This should not panic
	client.Close()
	assert.False(t, client.connected, "Client should be marked as disconnected")
	
	// Case 3: Test concurrency - multiple Close() calls
	client = NewNATSClient("nats://localhost:4222")
	mockConn = &mockNatsConn{connected: true}
	client.conn = mockConn
	client.connected = true
	
	// Call Close concurrently from multiple goroutines
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.Close()
		}()
	}
	wg.Wait()
	
	// Final check
	assert.False(t, client.connected, "Client should be marked as disconnected")
	assert.True(t, mockConn.closed, "Connection's Close method should be called")
}