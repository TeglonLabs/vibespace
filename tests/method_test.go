package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/rpcmethods"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestMethods(t *testing.T) {
	// Print out key constants and types that might help us identify the correct method names
	t.Logf("JSONRPC_VERSION: %s", mcp.JSONRPC_VERSION)
	
	// Create a test server with a resource and tool
	s := server.NewMCPServer("Test", "1.0")
	
	// Add a resource
	resource := mcp.NewResource("test://resource", "Test Resource")
	s.AddResource(resource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "test://resource",
				MIMEType: "text/plain",
				Text:     "Hello, world!",
			},
		}, nil
	})
	
	// Add a tool
	tool := mcp.NewTool("test_tool", func(t *mcp.Tool) {
		t.Description = "Test tool"
	})
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("Hello from tool!"), nil
	})
	
	// Use reflection to examine the server object and find registered methods
	serverValue := reflect.ValueOf(s).Elem()
	t.Logf("Server type fields:")
	for i := 0; i < serverValue.NumField(); i++ {
		field := serverValue.Type().Field(i)
		t.Logf("  %s: %s", field.Name, field.Type)
		
		// Try to find handlers or methods
		if field.Name == "handlers" || field.Name == "methods" || field.Name == "jsonrpcHandlers" {
			fieldValue := serverValue.Field(i)
			t.Logf("  %s contents:", field.Name)
			
			// If it's a map, print the keys
			if fieldValue.Kind() == reflect.Map {
				for _, key := range fieldValue.MapKeys() {
					t.Logf("    %v", key)
				}
			}
		}
	}
	
	// Print out the available methods from our package
	t.Logf("Available methods from rpcmethods package:")
	for _, method := range rpcmethods.ListMethods() {
		t.Logf("  - %s", method)
	}

	// Try various possible method names
	methodsToTry := []struct {
		method string
		params string
	}{
		// Use our standard method names from the rpcmethods package
		{rpcmethods.MethodResourceRead, `{"uri": "test://resource"}`},
		{rpcmethods.MethodToolCall, `{"name": "test_tool", "arguments": {}}`},
		
		// Also try the common variations to show they don't work
		{"mcp.resource.read", `{"uri": "test://resource"}`},
		{"resource.read", `{"uri": "test://resource"}`},
		{"mcp/resource/read", `{"uri": "test://resource"}`},
		{"get", `{"uri": "test://resource"}`},
		{"read", `{"uri": "test://resource"}`},
		{"call", `{"name": "test_tool", "arguments": {}}`},
		{"mcp/v1/resource.read", `{"uri": "test://resource"}`},
		{"mcp.resources.read", `{"uri": "test://resource"}`},
		{"mcp/v2/resource.read", `{"uri": "test://resource"}`},
		{"mcpresourceread", `{"uri": "test://resource"}`},
		{"Resources.Read", `{"uri": "test://resource"}`},
	}
	
	t.Log("Testing various method names:")
	for _, m := range methodsToTry {
		requestJSON := []byte(fmt.Sprintf(`{
			"jsonrpc": "2.0", 
			"id": "test-id",
			"method": "%s",
			"params": %s
		}`, m.method, m.params))
		
		response := s.HandleMessage(context.Background(), requestJSON)
		jsonRPCError, isError := response.(mcp.JSONRPCError)
		
		if isError {
			t.Logf("Method '%s' response: %s", m.method, jsonRPCError.Error.Message)
			
			// For methods that don't work, show suggestions
			if !rpcmethods.IsValidMethod(m.method) {
				suggestion := rpcmethods.GetMethodSuggestions(m.method)
				t.Logf("   Suggestion: %s", suggestion)
			}
		} else {
			t.Logf("Method '%s' seems to work!", m.method)
			jsonBytes, _ := json.MarshalIndent(response, "", "  ")
			t.Logf("Response: %s", jsonBytes)
		}
	}
}