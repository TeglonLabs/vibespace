package rpcmethods

import (
	"reflect"
	"strings"
	"testing"
)

func TestListMethods(t *testing.T) {
	methods := ListMethods()
	
	// Should return at least our core methods
	if len(methods) < 2 {
		t.Errorf("Expected at least 2 methods, got %d", len(methods))
	}
	
	// Check for specific methods
	containsMethod := func(methodName string) bool {
		for _, m := range methods {
			if m == methodName {
				return true
			}
		}
		return false
	}
	
	if !containsMethod(MethodResourceRead) {
		t.Errorf("Expected methods to contain '%s'", MethodResourceRead)
	}
	
	if !containsMethod(MethodToolCall) {
		t.Errorf("Expected methods to contain '%s'", MethodToolCall)
	}
}

func TestGetMethodInfo(t *testing.T) {
	// Test with valid method
	info := GetMethodInfo(MethodResourceRead)
	if info == nil {
		t.Fatalf("Expected method info for '%s', got nil", MethodResourceRead)
	}
	
	if info.Name != MethodResourceRead {
		t.Errorf("Expected name '%s', got '%s'", MethodResourceRead, info.Name)
	}
	
	if info.Category != "resource" {
		t.Errorf("Expected category 'resource', got '%s'", info.Category)
	}
	
	// Test with invalid method
	info = GetMethodInfo("nonexistent.method")
	if info != nil {
		t.Errorf("Expected nil for nonexistent method, got info")
	}
}

func TestFindMethod(t *testing.T) {
	// Test exact match
	matches := FindMethod(MethodResourceRead)
	if len(matches) != 1 || matches[0] != MethodResourceRead {
		t.Errorf("Expected 1 exact match '%s', got %v", MethodResourceRead, matches)
	}
	
	// Test partial match
	matches = FindMethod("resource")
	if len(matches) < 1 || !contains(matches, MethodResourceRead) {
		t.Errorf("Expected partial matches including '%s', got %v", MethodResourceRead, matches)
	}
	
	// Test case insensitive
	matches = FindMethod("RESOURCE")
	if len(matches) < 1 || !contains(matches, MethodResourceRead) {
		t.Errorf("Expected case-insensitive matches including '%s', got %v", MethodResourceRead, matches)
	}
	
	// Test no matches
	matches = FindMethod("xyz123")
	if len(matches) != 0 {
		t.Errorf("Expected no matches for nonsense query, got %v", matches)
	}
}

func TestFormatRequests(t *testing.T) {
	// Test resource request formatting
	resourceReq := FormatResourceRequest("world://list", "test-id")
	expectedResourceReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "test-id",
		"method":  MethodResourceRead,
		"params": map[string]interface{}{
			"uri": "world://list",
		},
	}
	
	if !reflect.DeepEqual(resourceReq, expectedResourceReq) {
		t.Errorf("Resource request formatting does not match expected.\nGot: %v\nExpected: %v", resourceReq, expectedResourceReq)
	}
	
	// Test tool request formatting
	toolArgs := map[string]interface{}{
		"id":   "test-world",
		"name": "Test World",
	}
	
	toolReq := FormatToolRequest("create_world", toolArgs, "test-id")
	expectedToolReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "test-id",
		"method":  MethodToolCall,
		"params": map[string]interface{}{
			"name": "create_world",
			"arguments": map[string]interface{}{
				"id":   "test-world",
				"name": "Test World",
			},
		},
	}
	
	if !reflect.DeepEqual(toolReq, expectedToolReq) {
		t.Errorf("Tool request formatting does not match expected.\nGot: %v\nExpected: %v", toolReq, expectedToolReq)
	}
}

func TestGetMethodSuggestions(t *testing.T) {
	// Test with a close match
	suggestion := GetMethodSuggestions("resource.read")
	if !contains([]string{suggestion}, "method.resource.read") {
		t.Errorf("Expected suggestion to contain method.resource.read, got: %s", suggestion)
	}
	
	// Test with no match
	suggestion = GetMethodSuggestions("completely.wrong")
	if !contains([]string{suggestion}, MethodResourceRead) || !contains([]string{suggestion}, MethodToolCall) {
		t.Errorf("Expected suggestion to include standard methods, got: %s", suggestion)
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item || strings.Contains(s, item) {
			return true
		}
	}
	return false
}