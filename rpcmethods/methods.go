// Package rpcmethods provides utilities for working with JSON-RPC methods in the vibespace MCP experience.
// It offers method discovery, documentation, and helper functions to simplify JSON-RPC interaction.
package rpcmethods

import (
	"fmt"
	"sort"
	"strings"
)

// Constants for standard method names used in vibespace MCP
const (
	// Core methods
	MethodResourceRead = "method.resource.read"
	MethodToolCall     = "method.tool.call"

	// For potential future extensions
	MethodPromptCreate       = "method.prompt.create"
	MethodNotificationListen = "method.notification.listen"
)

// MethodInfo contains documentation and usage information about a JSON-RPC method
type MethodInfo struct {
	Name        string            // The method name to use in JSON-RPC requests
	Description string            // Human-readable description
	Category    string            // Category (e.g., "resource", "tool", etc.)
	Parameters  map[string]string // Parameter names and descriptions
	Example     string            // Example JSON payload
}

// Global registry of all methods and their documentation
var methodRegistry = map[string]MethodInfo{
	MethodResourceRead: {
		Name:        MethodResourceRead,
		Description: "Read a resource by URI",
		Category:    "resource",
		Parameters: map[string]string{
			"uri": "The URI of the resource to read",
		},
		Example: `{
  "jsonrpc": "2.0",
  "id": "request-id",
  "method": "method.resource.read",
  "params": {
    "uri": "world://list"
  }
}`,
	},
	MethodToolCall: {
		Name:        MethodToolCall,
		Description: "Call a tool by name with arguments",
		Category:    "tool",
		Parameters: map[string]string{
			"name":      "The name of the tool to call",
			"arguments": "Object containing tool-specific arguments",
		},
		Example: `{
  "jsonrpc": "2.0",
  "id": "request-id",
  "method": "method.tool.call",
  "params": {
    "name": "create_world",
    "arguments": {
      "id": "test-world",
      "name": "Test World",
      "description": "A test world",
      "type": "VIRTUAL"
    }
  }
}`,
	},
}

// MethodFormatter is a function type that formats method names
// This allows for customization of method name formats
type MethodFormatter func(category, action string) string

// DefaultMethodFormatter formats methods in the standard vibespace MCP pattern: method.<category>.<action>
func DefaultMethodFormatter(category, action string) string {
	return fmt.Sprintf("method.%s.%s", category, action)
}

// ListMethods returns a sorted list of all registered method names
func ListMethods() []string {
	methods := make([]string, 0, len(methodRegistry))
	for method := range methodRegistry {
		methods = append(methods, method)
	}
	sort.Strings(methods)
	return methods
}

// GetMethodInfo returns detailed information about a specific method
// Returns nil if the method is not registered
func GetMethodInfo(methodName string) *MethodInfo {
	info, exists := methodRegistry[methodName]
	if !exists {
		return nil
	}
	return &info
}

// FindMethod searches for a method with fuzzy matching, helping users who might
// be using incorrect method name formats
func FindMethod(partialName string) []string {
	matches := []string{}
	lowercaseName := strings.ToLower(partialName)
	
	for method := range methodRegistry {
		if strings.Contains(strings.ToLower(method), lowercaseName) {
			matches = append(matches, method)
		}
	}
	
	sort.Strings(matches)
	return matches
}

// FormatResourceRequest helps create a properly formatted JSON-RPC resource read request
func FormatResourceRequest(uri string, requestID string) map[string]interface{} {
	if requestID == "" {
		requestID = "request-id"
	}
	
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      requestID,
		"method":  MethodResourceRead,
		"params": map[string]interface{}{
			"uri": uri,
		},
	}
}

// FormatToolRequest helps create a properly formatted JSON-RPC tool call request
func FormatToolRequest(toolName string, arguments map[string]interface{}, requestID string) map[string]interface{} {
	if requestID == "" {
		requestID = "request-id"
	}
	
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      requestID,
		"method":  MethodToolCall,
		"params": map[string]interface{}{
			"name":      toolName,
			"arguments": arguments,
		},
	}
}

// IsValidMethod checks if a method name is registered
func IsValidMethod(methodName string) bool {
	_, exists := methodRegistry[methodName]
	return exists
}

// GetMethodSuggestions provides suggestions when an invalid method is used
// This is helpful for error messages to guide users to the correct method
func GetMethodSuggestions(invalidMethod string) string {
	matches := FindMethod(invalidMethod)
	
	if len(matches) == 0 {
		return fmt.Sprintf("Method '%s' not found. Try using one of the standard methods: %s, %s", 
			invalidMethod, MethodResourceRead, MethodToolCall)
	}
	
	return fmt.Sprintf("Method '%s' not found. Did you mean: %s?", 
		invalidMethod, strings.Join(matches, ", "))
}