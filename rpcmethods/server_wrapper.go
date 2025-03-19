package rpcmethods

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPMethodWrapper wraps an MCP server to provide improved method handling
type MCPMethodWrapper struct {
	Server *server.MCPServer
}

// WrapMCPServer creates a new wrapper around an MCP server
func WrapMCPServer(s *server.MCPServer) *MCPMethodWrapper {
	return &MCPMethodWrapper{Server: s}
}

// HandleMessage intercepts JSON-RPC messages to provide better error messages
// for method not found errors, including suggestions for the correct method names
func (w *MCPMethodWrapper) HandleMessage(ctx context.Context, message json.RawMessage) mcp.JSONRPCMessage {
	// Parse the message to get the method name
	var req map[string]interface{}
	if err := json.Unmarshal(message, &req); err != nil {
		// Pass through to the server if we can't parse
		return w.Server.HandleMessage(ctx, message)
	}

	// Extract method name
	methodIface, ok := req["method"]
	if !ok {
		// Pass through to the server if no method in request
		return w.Server.HandleMessage(ctx, message)
	}

	method, ok := methodIface.(string)
	if !ok {
		// Pass through to the server if method is not a string
		return w.Server.HandleMessage(ctx, message)
	}

	// Let's normalize common method name variants
	normalizedMethod := normalizeMethodName(method)
	if normalizedMethod != method {
		// If we normalized to a different method, replace it in the request
		req["method"] = normalizedMethod
		newMessage, err := json.Marshal(req)
		if err != nil {
			// If we can't re-marshal, just use the original
			return w.Server.HandleMessage(ctx, message)
		}
		message = newMessage
	}

	// Forward to the server
	result := w.Server.HandleMessage(ctx, message)

	// Check if it's an error about method not found
	if jsonRPCError, isError := result.(mcp.JSONRPCError); isError {
		if jsonRPCError.Error.Code == -32601 { // Method not found error code
			// Append suggestion to the error message
			suggestion := GetMethodSuggestions(method)
			jsonRPCError.Error.Message = fmt.Sprintf("%s\n%s", jsonRPCError.Error.Message, suggestion)
			return jsonRPCError
		}
	}

	return result
}

// normalizeMethodName tries to normalize common method name variants to the standard ones
func normalizeMethodName(method string) string {
	// Replace slashes with dots
	method = strings.ReplaceAll(method, "/", ".")
	
	// Special case handling for resource reads
	if strings.Contains(method, "resource") && strings.Contains(method, "read") {
		return MethodResourceRead
	}
	
	// Special case handling for tool calls
	if strings.Contains(method, "tool") && (strings.Contains(method, "call") || 
	   strings.Contains(method, "tool.call") || strings.Contains(method, "tools.call")) {
		return MethodToolCall
	}
	
	// Handle MCP prefixes
	if strings.HasPrefix(method, "mcp.") {
		method = method[4:] // Remove mcp. prefix
	}
	
	// Check exact matches
	if IsValidMethod(method) {
		return method
	}
	
	// Check other common patterns
	switch method {
	case "resource.read", "resources.read", "read_resource", "get_resource":
		return MethodResourceRead
	case "tool.call", "tools.call", "call_tool", "execute_tool":
		return MethodToolCall
	case "Resources.Read":
		return MethodResourceRead
	case "Tools.Call":
		return MethodToolCall
	}
	
	// If no matches, return original
	return method
}