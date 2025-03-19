package rpcmethods

import (
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestWrapMCPServer(t *testing.T) {
	// Create a nil server for testing (we're just checking that the wrapper is created correctly)
	var s *server.MCPServer
	
	wrapper := WrapMCPServer(s)
	
	if wrapper == nil {
		t.Error("WrapMCPServer returned nil")
	}
	
	if wrapper.Server != s {
		t.Errorf("WrapMCPServer did not correctly set the Server field")
	}
}