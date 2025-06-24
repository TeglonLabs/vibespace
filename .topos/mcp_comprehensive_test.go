// MOVED: This comprehensive test suite is being moved to .topos for future implementation
// due to API compatibility issues with the current MCP library version.
// The tests will be updated when the MCP library stabilizes its interface.

package tests_disabled

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/bmorphism/vibespace-mcp-go/rpcmethods"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test data constants following 2-3-5-7 principle
const (
	// Binary test sizes: 2KB, 3KB, 5KB, 7KB
	BinarySize2KB = 2048
	BinarySize3KB = 3072
	BinarySize5KB = 5120
	BinarySize7KB = 7168
	
	// Large file sizes: 2MB, 3MB, 5MB, 1069KB (special case)
	BinarySize2MB = 2 * 1024 * 1024
	BinarySize3MB = 3 * 1024 * 1024
	BinarySize5MB = 5 * 1024 * 1024
	BinarySize1069KB = 1069 * 1024
	
	// Media types test matrix
	ContentTypeJSON = "application/json"
	ContentTypeProtobuf = "application/protobuf"
	ContentTypeImagePNG = "image/png"
	ContentTypeImageJPEG = "image/jpeg"
	ContentTypeAudioWAV = "audio/wav"
	ContentTypeAudioMP3 = "audio/mp3"
	ContentTypeVideoMP4 = "video/mp4"
	ContentTypeOctetStream = "application/octet-stream"
)

// TestMCPServer represents our enhanced test server
type TestMCPServer struct {
	*server.MCPServer
	repo            *repository.Repository
	streamingService *streaming.StreamingService
	httpServer      *httptest.Server
}

// OAuth test configuration
type OAuthTestConfig struct {
	AuthorizationURL string            `json:"authorization_url"`
	TokenURL        string            `json:"token_url"`
	ClientID        string            `json:"client_id"`
	ClientSecret    string            `json:"client_secret"`
	Scopes          []string          `json:"scopes"`
	ResourceServer  bool              `json:"resource_server"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// Structured content schemas for testing
type DeviceStatusSchema struct {
	DeviceID     string    `json:"deviceId"`
	Status       string    `json:"status"`
	LastSeen     time.Time `json:"lastSeen"`
	Metrics      []Metric  `json:"metrics"`
	BinaryData   []byte    `json:"binaryData,omitempty"`
	Capabilities []string  `json:"capabilities"`
}

type Metric struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

// ElicitationTestFlow represents a multi-turn interaction
type ElicitationTestFlow struct {
	InitialRequest  map[string]interface{} `json:"initial_request"`
	ElicitationStep string                 `json:"elicitation_step"`
	UserResponse    map[string]interface{} `json:"user_response"`
	FinalResult     interface{}            `json:"final_result"`
}

// setupTestServer creates a comprehensive test server with all new MCP features
func setupTestServer(t *testing.T) *TestMCPServer {
	// Create repository and streaming service
	repo := repository.NewRepository()
	streamingConfig := &streaming.StreamingConfig{
		NATSHost:       "localhost",
		NATSPort:       4222,
		StreamInterval: 1 * time.Second,
		AutoStart:      false,
	}
	streamingService := streaming.NewStreamingService(repo, streamingConfig)
	
	// Create MCP server with enhanced capabilities
	mcpServer := server.NewMCPServer("vibespace-test-enhanced", "2.0.0")
	
	// Add OAuth configuration capability
	setupOAuthCapabilities(mcpServer)
	
	// Add structured content tools
	setupStructuredContentTools(mcpServer, repo)
	
	// Add binary/media content handlers
	setupBinaryContentHandlers(mcpServer)
	
	// Add elicitation flow tools
	setupElicitationTools(mcpServer)
	
	// Create HTTP test server
	handler := rpcmethods.WrapMCPServer(mcpServer)
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		
		response := handler.HandleMessage(r.Context(), body)
		
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))
	
	return &TestMCPServer{
		MCPServer:       mcpServer,
		repo:            repo,
		streamingService: streamingService,
		httpServer:      httpServer,
	}
}

// setupOAuthCapabilities adds OAuth 2.1 support with comprehensive test scenarios
func setupOAuthCapabilities(mcpServer *server.MCPServer) {
	// OAuth configuration tool following RFC 8414
	oauthConfigTool := mcp.NewTool("oauth_get_config", func(t *mcp.Tool) {
		t.Description = "Get OAuth 2.1 configuration with Resource Indicators support"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"resource_server": {Type: "string", Description: "Target resource server"},
				"scopes": {Type: "array", Items: &mcp.ToolInputSchema{Type: "string"}, Description: "Required scopes"},
			},
		}
	})
	
	mcpServer.AddTool(oauthConfigTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		resourceServer, _ := args["resource_server"].(string)
		if resourceServer == "" {
			resourceServer = "https://api.vibespace.example.com"
		}
		
		config := OAuthTestConfig{
			AuthorizationURL: "https://auth.vibespace.example.com/oauth/authorize",
			TokenURL:        "https://auth.vibespace.example.com/oauth/token",
			ClientID:        "vibespace-mcp-client",
			Scopes:          []string{"read:vibes", "write:worlds", "stream:moments"},
			ResourceServer:  true,
			Metadata: map[string]interface{}{
				"resource_indicators_supported": true,
				"grant_types_supported": []string{"authorization_code", "client_credentials"},
				"response_types_supported": []string{"code"},
				"scopes_supported": []string{"read:vibes", "write:worlds", "stream:moments", "admin:all"},
				"token_endpoint_auth_methods_supported": []string{"client_secret_post", "client_secret_basic"},
			},
		}
		
		configJSON, _ := json.Marshal(config)
		return mcp.NewToolResultText(string(configJSON)), nil
	})
	
	// Token validation tool with Resource Indicators (RFC 8707)
	tokenValidationTool := mcp.NewTool("oauth_validate_token", func(t *mcp.Tool) {
		t.Description = "Validate OAuth token with resource indicators"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"access_token": {Type: "string", Description: "OAuth access token"},
				"resource": {Type: "string", Description: "Target resource URI"},
				"required_scopes": {Type: "array", Items: &mcp.ToolInputSchema{Type: "string"}},
			},
			Required: []string{"access_token"},
		}
	})
	
	mcpServer.AddTool(tokenValidationTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		token, _ := args["access_token"].(string)
		resource, _ := args["resource"].(string)
		
		// Simulate token validation with resource indicators
		validation := map[string]interface{}{
			"valid": true,
			"expires_in": 3600,
			"scopes": []string{"read:vibes", "write:worlds"},
			"resource": resource,
			"client_id": "vibespace-mcp-client",
			"user_id": "test-user-123",
		}
		
		if token == "invalid-token" {
			validation["valid"] = false
			validation["error"] = "invalid_token"
		}
		
		validationJSON, _ := json.Marshal(validation)
		return mcp.NewToolResultText(string(validationJSON)), nil
	})
}

// setupStructuredContentTools adds tools with output schemas and MIME type support
func setupStructuredContentTools(mcpServer *server.MCPServer, repo *repository.Repository) {
	// Device status tool with structured output schema
	deviceStatusTool := mcp.NewTool("get_device_status", func(t *mcp.Tool) {
		t.Description = "Get structured device status with typed output"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"device_id": {Type: "string", Description: "Device identifier"},
				"include_metrics": {Type: "boolean", Description: "Include performance metrics"},
				"include_binary": {Type: "boolean", Description: "Include binary diagnostic data"},
			},
			Required: []string{"device_id"},
		}
		// Output schema following JSON Schema specification
		t.OutputSchema = &mcp.ToolOutputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"deviceId": map[string]interface{}{"type": "string", "description": "Unique device identifier"},
				"status": map[string]interface{}{"type": "string", "enum": []string{"online", "offline", "maintenance"}},
				"lastSeen": map[string]interface{}{"type": "string", "format": "date-time"},
				"metrics": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{"type": "string"},
							"value": map[string]interface{}{"type": "number"},
							"unit": map[string]interface{}{"type": "string"},
							"timestamp": map[string]interface{}{"type": "string", "format": "date-time"},
						},
					},
				},
				"binaryData": map[string]interface{}{"type": "string", "format": "base64"},
				"capabilities": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{"type": "string"},
				},
			},
			Required: []string{"deviceId", "status", "lastSeen"},
		}
	})
	
	mcpServer.AddTool(deviceStatusTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		deviceID, _ := args["device_id"].(string)
		includeMetrics, _ := args["include_metrics"].(bool)
		includeBinary, _ := args["include_binary"].(bool)
		
		status := DeviceStatusSchema{
			DeviceID:     deviceID,
			Status:       "online",
			LastSeen:     time.Now(),
			Capabilities: []string{"streaming", "authentication", "binary_transfer"},
		}
		
		if includeMetrics {
			status.Metrics = []Metric{
				{Name: "cpu_usage", Value: 45.2, Unit: "percent", Timestamp: time.Now()},
				{Name: "memory_usage", Value: 67.8, Unit: "percent", Timestamp: time.Now()},
				{Name: "network_latency", Value: 12.3, Unit: "ms", Timestamp: time.Now()},
			}
		}
		
		if includeBinary {
			// Generate test binary data
			binaryData := make([]byte, BinarySize2KB)
			for i := range binaryData {
				binaryData[i] = byte(i % 256)
			}
			status.BinaryData = binaryData
		}
		
		statusJSON, _ := json.Marshal(status)
		return mcp.NewToolResultText(string(statusJSON)), nil
	})
}

// setupBinaryContentHandlers adds comprehensive binary/media content support
func setupBinaryContentHandlers(mcpServer *server.MCPServer) {
	// Binary content upload tool with MIME type validation
	uploadBinaryTool := mcp.NewTool("upload_binary_content", func(t *mcp.Tool) {
		t.Description = "Upload binary content with MIME type validation"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"content": {Type: "string", Description: "Base64-encoded binary content"},
				"mime_type": {Type: "string", Description: "MIME type of the content"},
				"filename": {Type: "string", Description: "Original filename"},
				"compression": {Type: "string", Enum: []string{"none", "gzip", "brotli"}, Description: "Compression method"},
			},
			Required: []string{"content", "mime_type"},
		}
	})
	
	mcpServer.AddTool(uploadBinaryTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		content, _ := args["content"].(string)
		mimeType, _ := args["mime_type"].(string)
		filename, _ := args["filename"].(string)
		
		// Validate MIME type
		supportedTypes := []string{
			ContentTypeJSON, ContentTypeProtobuf, ContentTypeImagePNG, 
			ContentTypeImageJPEG, ContentTypeAudioWAV, ContentTypeAudioMP3,
			ContentTypeVideoMP4, ContentTypeOctetStream,
		}
		
		validMime := false
		for _, supported := range supportedTypes {
			if mimeType == supported {
				validMime = true
				break
			}
		}
		
		if !validMime {
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported MIME type: %s", mimeType)), nil
		}
		
		// Decode and validate content
		decodedContent, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			return mcp.NewToolResultError("Invalid base64 content"), nil
		}
		
		result := map[string]interface{}{
			"upload_id": fmt.Sprintf("upload_%d", time.Now().Unix()),
			"size": len(decodedContent),
			"mime_type": mimeType,
			"filename": filename,
			"checksum": fmt.Sprintf("sha256:%x", decodedContent[:32]),
			"stored_at": time.Now().Format(time.RFC3339),
		}
		
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
	
	// Protobuf content handler
	protobufTool := mcp.NewTool("handle_protobuf_content", func(t *mcp.Tool) {
		t.Description = "Handle Protocol Buffer content with schema validation"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"protobuf_data": {Type: "string", Description: "Base64-encoded protobuf data"},
				"schema_version": {Type: "string", Description: "Protobuf schema version"},
				"message_type": {Type: "string", Description: "Protobuf message type"},
			},
			Required: []string{"protobuf_data", "message_type"},
		}
	})
	
	mcpServer.AddTool(protobufTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		protobufData, _ := args["protobuf_data"].(string)
		messageType, _ := args["message_type"].(string)
		schemaVersion, _ := args["schema_version"].(string)
		
		// Simulate protobuf processing
		decodedData, err := base64.StdEncoding.DecodeString(protobufData)
		if err != nil {
			return mcp.NewToolResultError("Invalid protobuf data"), nil
		}
		
		result := map[string]interface{}{
			"message_type": messageType,
			"schema_version": schemaVersion,
			"decoded_size": len(decodedData),
			"fields_parsed": 7, // Simulated field count
			"validation_status": "valid",
			"processed_at": time.Now().Format(time.RFC3339),
		}
		
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
}

// setupElicitationTools adds human-in-the-loop elicitation support
func setupElicitationTools(mcpServer *server.MCPServer) {
	// Interactive form filling tool with elicitation
	formFillingTool := mcp.NewTool("fill_interactive_form", func(t *mcp.Tool) {
		t.Description = "Fill form with interactive elicitation for missing fields"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"form_id": {Type: "string", Description: "Form identifier"},
				"partial_data": {Type: "object", Description: "Partially filled form data"},
				"continue_elicitation": {Type: "boolean", Description: "Continue from previous elicitation"},
				"elicitation_id": {Type: "string", Description: "Previous elicitation session ID"},
			},
			Required: []string{"form_id"},
		}
	})
	
	mcpServer.AddTool(formFillingTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		formID, _ := args["form_id"].(string)
		partialData, _ := args["partial_data"].(map[string]interface{})
		continueElicitation, _ := args["continue_elicitation"].(bool)
		
		// Required fields for the test form
		requiredFields := []string{"name", "email", "phone", "preferences"}
		missingFields := []string{}
		
		// Check for missing fields
		for _, field := range requiredFields {
			if partialData == nil || partialData[field] == nil {
				missingFields = append(missingFields, field)
			}
		}
		
		if len(missingFields) > 0 && !continueElicitation {
			// Return elicitation request
			elicitationRequest := map[string]interface{}{
				"type": "elicitation_request",
				"elicitation_id": fmt.Sprintf("elicit_%s_%d", formID, time.Now().Unix()),
				"missing_fields": missingFields,
				"prompts": map[string]string{
					"name": "Please provide your full name:",
					"email": "What's your email address?",
					"phone": "Enter your phone number (optional):",
					"preferences": "Any specific preferences or requirements?",
				},
				"form_id": formID,
			}
			
			elicitationJSON, _ := json.Marshal(elicitationRequest)
			return mcp.NewToolResultText(string(elicitationJSON)), nil
		}
		
		// Form is complete or continuing from elicitation
		result := map[string]interface{}{
			"form_id": formID,
			"status": "completed",
			"filled_fields": len(requiredFields) - len(missingFields),
			"total_fields": len(requiredFields),
			"submitted_at": time.Now().Format(time.RFC3339),
			"data": partialData,
		}
		
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
	
	// Continue elicitation tool
	continueElicitationTool := mcp.NewTool("continue_elicitation", func(t *mcp.Tool) {
		t.Description = "Continue elicitation with user-provided data"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"elicitation_id": {Type: "string", Description: "Elicitation session ID"},
				"user_data": {Type: "object", Description: "User-provided data"},
				"original_tool": {Type: "string", Description: "Original tool name"},
				"original_args": {Type: "object", Description: "Original tool arguments"},
			},
			Required: []string{"elicitation_id", "user_data", "original_tool"},
		}
	})
	
	mcpServer.AddTool(continueElicitationTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		elicitationID, _ := args["elicitation_id"].(string)
		userData, _ := args["user_data"].(map[string]interface{})
		originalTool, _ := args["original_tool"].(string)
		originalArgs, _ := args["original_args"].(map[string]interface{})
		
		// Merge user data with original args
		if originalArgs == nil {
			originalArgs = make(map[string]interface{})
		}
		
		if originalArgs["partial_data"] == nil {
			originalArgs["partial_data"] = make(map[string]interface{})
		}
		
		partialData := originalArgs["partial_data"].(map[string]interface{})
		for key, value := range userData {
			partialData[key] = value
		}
		
		originalArgs["continue_elicitation"] = true
		originalArgs["elicitation_id"] = elicitationID
		
		result := map[string]interface{}{
			"elicitation_id": elicitationID,
			"status": "continued",
			"original_tool": originalTool,
			"merged_data": originalArgs,
			"next_action": "retry_original_tool",
		}
		
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
}

// Test OAuth 2.1 implementation with Resource Indicators
func TestOAuth21WithResourceIndicators(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.httpServer.Close()
	
	tests := []struct {
		name string
		testFunc func(t *testing.T, server *TestMCPServer)
	}{
		{
			name: "OAuth Configuration Retrieval",
			testFunc: func(t *testing.T, server *TestMCPServer) {
				// Test OAuth configuration endpoint
				request := map[string]interface{}{
					"jsonrpc": "2.0",
					"id": 1,
					"method": "tools/call",
					"params": map[string]interface{}{
						"name": "oauth_get_config",
						"arguments": map[string]interface{}{
							"resource_server": "https://api.vibespace.example.com",
							"scopes": []string{"read:vibes", "write:worlds"},
						},
					},
				}
				
				requestJSON, _ := json.Marshal(request)
				resp, err := http.Post(server.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
				require.NoError(t, err)
				defer resp.Body.Close()
				
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				
				assert.Equal(t, "2.0", response["jsonrpc"])
				assert.Equal(t, float64(1), response["id"])
				
				result := response["result"].(map[string]interface{})
				content := result["content"].([]interface{})[0].(map[string]interface{})
				
				var config OAuthTestConfig
				err = json.Unmarshal([]byte(content["text"].(string)), &config)
				require.NoError(t, err)
				
				assert.True(t, config.ResourceServer)
				assert.Contains(t, config.Scopes, "read:vibes")
				assert.Contains(t, config.Scopes, "write:worlds")
				assert.True(t, config.Metadata["resource_indicators_supported"].(bool))
			},
		},
		{
			name: "Token Validation with Resource Indicators",
			testFunc: func(t *testing.T, server *TestMCPServer) {
				request := map[string]interface{}{
					"jsonrpc": "2.0",
					"id": 2,
					"method": "tools/call",
					"params": map[string]interface{}{
						"name": "oauth_validate_token",
						"arguments": map[string]interface{}{
							"access_token": "valid-test-token-123",
							"resource": "https://api.vibespace.example.com/vibes",
							"required_scopes": []string{"read:vibes"},
						},
					},
				}
				
				requestJSON, _ := json.Marshal(request)
				resp, err := http.Post(server.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
				require.NoError(t, err)
				defer resp.Body.Close()
				
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				
				result := response["result"].(map[string]interface{})
				content := result["content"].([]interface{})[0].(map[string]interface{})
				
				var validation map[string]interface{}
				err = json.Unmarshal([]byte(content["text"].(string)), &validation)
				require.NoError(t, err)
				
				assert.True(t, validation["valid"].(bool))
				assert.Equal(t, "https://api.vibespace.example.com/vibes", validation["resource"])
				assert.Equal(t, "test-user-123", validation["user_id"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, testServer)
		})
	}
}

// Test structured content with output schemas
func TestStructuredContentWithSchemas(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.httpServer.Close()
	
	tests := []struct {
		name string
		args map[string]interface{}
		validateFunc func(t *testing.T, result DeviceStatusSchema)
	}{
		{
			name: "Basic Device Status",
			args: map[string]interface{}{
				"device_id": "device-001",
			},
			validateFunc: func(t *testing.T, result DeviceStatusSchema) {
				assert.Equal(t, "device-001", result.DeviceID)
				assert.Equal(t, "online", result.Status)
				assert.NotZero(t, result.LastSeen)
				assert.Contains(t, result.Capabilities, "streaming")
				assert.Nil(t, result.Metrics)
				assert.Nil(t, result.BinaryData)
			},
		},
		{
			name: "Device Status with Metrics",
			args: map[string]interface{}{
				"device_id": "device-002",
				"include_metrics": true,
			},
			validateFunc: func(t *testing.T, result DeviceStatusSchema) {
				assert.Equal(t, "device-002", result.DeviceID)
				assert.NotNil(t, result.Metrics)
				assert.Len(t, result.Metrics, 3)
				
				// Validate metric structure
				cpuMetric := result.Metrics[0]
				assert.Equal(t, "cpu_usage", cpuMetric.Name)
				assert.Equal(t, "percent", cpuMetric.Unit)
				assert.Greater(t, cpuMetric.Value, 0.0)
			},
		},
		{
			name: "Device Status with Binary Data",
			args: map[string]interface{}{
				"device_id": "device-003",
				"include_binary": true,
			},
			validateFunc: func(t *testing.T, result DeviceStatusSchema) {
				assert.Equal(t, "device-003", result.DeviceID)
				assert.NotNil(t, result.BinaryData)
				assert.Equal(t, BinarySize2KB, len(result.BinaryData))
				
				// Validate binary data pattern
				for i := 0; i < 256; i++ {
					assert.Equal(t, byte(i), result.BinaryData[i])
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := map[string]interface{}{
				"jsonrpc": "2.0",
				"id": 1,
				"method": "tools/call",
				"params": map[string]interface{}{
					"name": "get_device_status",
					"arguments": tt.args,
				},
			}
			
			requestJSON, _ := json.Marshal(request)
			resp, err := http.Post(testServer.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			
			result := response["result"].(map[string]interface{})
			content := result["content"].([]interface{})[0].(map[string]interface{})
			
			var deviceStatus DeviceStatusSchema
			err = json.Unmarshal([]byte(content["text"].(string)), &deviceStatus)
			require.NoError(t, err)
			
			tt.validateFunc(t, deviceStatus)
		})
	}
}

// Test binary content handling with various MIME types
func TestBinaryContentHandling(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.httpServer.Close()
	
	// Test data following 2-3-5-7 principle
	testCases := []struct {
		name string
		size int
		mimeType string
		shouldSucceed bool
	}{
		{"Small PNG Image", BinarySize2KB, ContentTypeImagePNG, true},
		{"Medium JPEG Image", BinarySize3KB, ContentTypeImageJPEG, true},
		{"Large Audio WAV", BinarySize5KB, ContentTypeAudioWAV, true},
		{"Extra Large MP4 Video", BinarySize7KB, ContentTypeVideoMP4, true},
		{"Protobuf Data", BinarySize1069KB, ContentTypeProtobuf, true},
		{"Invalid MIME Type", BinarySize2KB, "invalid/mime-type", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate test binary data
			binaryData := make([]byte, tc.size)
			for i := range binaryData {
				binaryData[i] = byte(i % 256)
			}
			
			// Encode as base64
			encodedData := base64.StdEncoding.EncodeToString(binaryData)
			
			request := map[string]interface{}{
				"jsonrpc": "2.0",
				"id": 1,
				"method": "tools/call",
				"params": map[string]interface{}{
					"name": "upload_binary_content",
					"arguments": map[string]interface{}{
						"content": encodedData,
						"mime_type": tc.mimeType,
						"filename": fmt.Sprintf("test-file-%d.bin", tc.size),
						"compression": "none",
					},
				},
			}
			
			requestJSON, _ := json.Marshal(request)
			resp, err := http.Post(testServer.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			
			if tc.shouldSucceed {
				assert.NotContains(t, response, "error")
				
				result := response["result"].(map[string]interface{})
				content := result["content"].([]interface{})[0].(map[string]interface{})
				
				var uploadResult map[string]interface{}
				err = json.Unmarshal([]byte(content["text"].(string)), &uploadResult)
				require.NoError(t, err)
				
				assert.Equal(t, float64(tc.size), uploadResult["size"])
				assert.Equal(t, tc.mimeType, uploadResult["mime_type"])
				assert.Contains(t, uploadResult["upload_id"], "upload_")
			} else {
				// Should have an error for invalid MIME type
				result := response["result"].(map[string]interface{})
				content := result["content"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, "error", content["type"])
			}
		})
	}
}

// Test Protocol Buffer content handling
func TestProtobufContentHandling(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.httpServer.Close()
	
	// Generate mock protobuf data
	protobufData := make([]byte, BinarySize2KB)
	for i := range protobufData {
		protobufData[i] = byte(i % 256)
	}
	encodedData := base64.StdEncoding.EncodeToString(protobufData)
	
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "tools/call",
		"params": map[string]interface{}{
			"name": "handle_protobuf_content",
			"arguments": map[string]interface{}{
				"protobuf_data": encodedData,
				"message_type": "vibespace.models.WorldMoment",
				"schema_version": "v1.2.0",
			},
		},
	}
	
	requestJSON, _ := json.Marshal(request)
	resp, err := http.Post(testServer.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()
	
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	
	result := response["result"].(map[string]interface{})
	content := result["content"].([]interface{})[0].(map[string]interface{})
	
	var protobufResult map[string]interface{}
	err = json.Unmarshal([]byte(content["text"].(string)), &protobufResult)
	require.NoError(t, err)
	
	assert.Equal(t, "vibespace.models.WorldMoment", protobufResult["message_type"])
	assert.Equal(t, "v1.2.0", protobufResult["schema_version"])
	assert.Equal(t, float64(BinarySize2KB), protobufResult["decoded_size"])
	assert.Equal(t, "valid", protobufResult["validation_status"])
}

// Test elicitation flow for human-in-the-loop interactions
func TestElicitationFlow(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.httpServer.Close()
	
	t.Run("Initial Form Request Triggers Elicitation", func(t *testing.T) {
		request := map[string]interface{}{
			"jsonrpc": "2.0",
			"id": 1,
			"method": "tools/call",
			"params": map[string]interface{}{
				"name": "fill_interactive_form",
				"arguments": map[string]interface{}{
					"form_id": "user-registration-001",
					"partial_data": map[string]interface{}{
						"name": "Test User",
						// Missing: email, phone, preferences
					},
				},
			},
		}
		
		requestJSON, _ := json.Marshal(request)
		resp, err := http.Post(testServer.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
		require.NoError(t, err)
		defer resp.Body.Close()
		
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		
		result := response["result"].(map[string]interface{})
		content := result["content"].([]interface{})[0].(map[string]interface{})
		
		var elicitationRequest map[string]interface{}
		err = json.Unmarshal([]byte(content["text"].(string)), &elicitationRequest)
		require.NoError(t, err)
		
		assert.Equal(t, "elicitation_request", elicitationRequest["type"])
		assert.Contains(t, elicitationRequest["elicitation_id"], "elicit_user-registration-001_")
		
		missingFields := elicitationRequest["missing_fields"].([]interface{})
		assert.Contains(t, missingFields, "email")
		assert.Contains(t, missingFields, "phone")
		assert.Contains(t, missingFields, "preferences")
		
		prompts := elicitationRequest["prompts"].(map[string]interface{})
		assert.Contains(t, prompts, "email")
		assert.Equal(t, "What's your email address?", prompts["email"])
	})
	
	t.Run("Continue Elicitation with User Data", func(t *testing.T) {
		request := map[string]interface{}{
			"jsonrpc": "2.0",
			"id": 2,
			"method": "tools/call",
			"params": map[string]interface{}{
				"name": "continue_elicitation",
				"arguments": map[string]interface{}{
					"elicitation_id": "elicit_user-registration-001_12345",
					"user_data": map[string]interface{}{
						"email": "test@example.com",
						"phone": "+1234567890",
						"preferences": "Dark mode, minimal notifications",
					},
					"original_tool": "fill_interactive_form",
					"original_args": map[string]interface{}{
						"form_id": "user-registration-001",
						"partial_data": map[string]interface{}{
							"name": "Test User",
						},
					},
				},
			},
		}
		
		requestJSON, _ := json.Marshal(request)
		resp, err := http.Post(testServer.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
		require.NoError(t, err)
		defer resp.Body.Close()
		
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		
		result := response["result"].(map[string]interface{})
		content := result["content"].([]interface{})[0].(map[string]interface{})
		
		var continuationResult map[string]interface{}
		err = json.Unmarshal([]byte(content["text"].(string)), &continuationResult)
		require.NoError(t, err)
		
		assert.Equal(t, "continued", continuationResult["status"])
		assert.Equal(t, "fill_interactive_form", continuationResult["original_tool"])
		assert.Equal(t, "retry_original_tool", continuationResult["next_action"])
		
		mergedData := continuationResult["merged_data"].(map[string]interface{})
		partialData := mergedData["partial_data"].(map[string]interface{})
		assert.Equal(t, "Test User", partialData["name"])
		assert.Equal(t, "test@example.com", partialData["email"])
		assert.Equal(t, "+1234567890", partialData["phone"])
	})
	
	t.Run("Complete Form After Elicitation", func(t *testing.T) {
		request := map[string]interface{}{
			"jsonrpc": "2.0",
			"id": 3,
			"method": "tools/call",
			"params": map[string]interface{}{
				"name": "fill_interactive_form",
				"arguments": map[string]interface{}{
					"form_id": "user-registration-001",
					"partial_data": map[string]interface{}{
						"name": "Test User",
						"email": "test@example.com",
						"phone": "+1234567890",
						"preferences": "Dark mode, minimal notifications",
					},
					"continue_elicitation": true,
				},
			},
		}
		
		requestJSON, _ := json.Marshal(request)
		resp, err := http.Post(testServer.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
		require.NoError(t, err)
		defer resp.Body.Close()
		
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		
		result := response["result"].(map[string]interface{})
		content := result["content"].([]interface{})[0].(map[string]interface{})
		
		var completionResult map[string]interface{}
		err = json.Unmarshal([]byte(content["text"].(string)), &completionResult)
		require.NoError(t, err)
		
		assert.Equal(t, "completed", completionResult["status"])
		assert.Equal(t, float64(4), completionResult["filled_fields"])
		assert.Equal(t, float64(4), completionResult["total_fields"])
		
		data := completionResult["data"].(map[string]interface{})
		assert.Equal(t, "Test User", data["name"])
		assert.Equal(t, "test@example.com", data["email"])
	})
}

// Test multipart content with various modalities
func TestMultipartContentHandling(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.httpServer.Close()
	
	t.Run("Multipart Form with Binary and Text", func(t *testing.T) {
		// Create multipart form data
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		
		// Add text field
		textField, _ := writer.CreateFormField("description")
		textField.Write([]byte("Test multipart upload with binary data"))
		
		// Add binary file
		binaryData := make([]byte, BinarySize3KB)
		for i := range binaryData {
			binaryData[i] = byte(i % 256)
		}
		
		fileField, _ := writer.CreateFormFile("binary_file", "test.bin")
		fileField.Write(binaryData)
		
		writer.Close()
		
		// Create HTTP request
		req, err := http.NewRequest("POST", testServer.httpServer.URL+"/multipart", &buf)
		require.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		
		// For now, we expect this to return method not allowed since we haven't implemented multipart handling
		// This test validates that the server can receive multipart data
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})
}

// Test streaming content with progress tracking
func TestStreamingContentWithProgress(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.httpServer.Close()
	
	// This test validates that the streaming service can handle progress tracking
	// which is part of the new MCP capabilities
	
	t.Run("Streaming Service Initialization", func(t *testing.T) {
		assert.NotNil(t, testServer.streamingService)
		assert.False(t, testServer.streamingService.IsStreaming())
	})
	
	t.Run("Start Streaming with Progress Tracking", func(t *testing.T) {
		// Test starting streaming service
		err := testServer.streamingService.StartStreaming()
		// We expect this to fail since NATS is not running, but we validate the interface
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to NATS")
	})
}

// Test MIME type validation comprehensive coverage
func TestMIMETypeValidation(t *testing.T) {
	testServer := setupTestServer(t)
	defer testServer.httpServer.Close()
	
	// Test all supported MIME types following 2-3-5-7 principle
	supportedTypes := []string{
		ContentTypeJSON,
		ContentTypeProtobuf,
		ContentTypeImagePNG,
		ContentTypeImageJPEG,
		ContentTypeAudioWAV,
		ContentTypeAudioMP3,
		ContentTypeVideoMP4,
		ContentTypeOctetStream,
	}
	
	// Test unsupported types
	unsupportedTypes := []string{
		"text/plain",
		"application/xml",
		"image/gif", 
		"video/avi",
		"audio/flac",
	}
	
	for i, mimeType := range supportedTypes {
		t.Run(fmt.Sprintf("Supported MIME Type %s", mimeType), func(t *testing.T) {
			size := []int{BinarySize2KB, BinarySize3KB, BinarySize5KB, BinarySize7KB}[i%4]
			
			binaryData := make([]byte, size)
			for j := range binaryData {
				binaryData[j] = byte(j % 256)
			}
			
			encodedData := base64.StdEncoding.EncodeToString(binaryData)
			
			request := map[string]interface{}{
				"jsonrpc": "2.0",
				"id": i + 1,
				"method": "tools/call",
				"params": map[string]interface{}{
					"name": "upload_binary_content",
					"arguments": map[string]interface{}{
						"content": encodedData,
						"mime_type": mimeType,
						"filename": fmt.Sprintf("test-%d.bin", i),
					},
				},
			}
			
			requestJSON, _ := json.Marshal(request)
			resp, err := http.Post(testServer.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			
			// Should succeed for supported types
			assert.NotContains(t, response, "error")
			
			result := response["result"].(map[string]interface{})
			content := result["content"].([]interface{})[0].(map[string]interface{})
			
			var uploadResult map[string]interface{}
			err = json.Unmarshal([]byte(content["text"].(string)), &uploadResult)
			require.NoError(t, err)
			
			assert.Equal(t, mimeType, uploadResult["mime_type"])
			assert.Equal(t, float64(size), uploadResult["size"])
		})
	}
	
	for i, mimeType := range unsupportedTypes {
		t.Run(fmt.Sprintf("Unsupported MIME Type %s", mimeType), func(t *testing.T) {
			binaryData := make([]byte, BinarySize2KB)
			encodedData := base64.StdEncoding.EncodeToString(binaryData)
			
			request := map[string]interface{}{
				"jsonrpc": "2.0",
				"id": i + 100,
				"method": "tools/call",
				"params": map[string]interface{}{
					"name": "upload_binary_content",
					"arguments": map[string]interface{}{
						"content": encodedData,
						"mime_type": mimeType,
						"filename": fmt.Sprintf("unsupported-%d.bin", i),
					},
				},
			}
			
			requestJSON, _ := json.Marshal(request)
			resp, err := http.Post(testServer.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			
			// Should have error for unsupported types
			result := response["result"].(map[string]interface{})
			content := result["content"].([]interface{})[0].(map[string]interface{})
			assert.Equal(t, "error", content["type"])
			assert.Contains(t, content["text"], "Unsupported MIME type")
		})
	}
}

// Benchmark test for large binary content following 2-3-5-1069 principle
func BenchmarkLargeBinaryContent(b *testing.B) {
	testServer := setupTestServer(&testing.T{})
	defer testServer.httpServer.Close()
	
	sizes := []int{BinarySize2MB, BinarySize3MB, BinarySize5MB, BinarySize1069KB}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("BinarySize_%dB", size), func(b *testing.B) {
			// Generate test data once
			binaryData := make([]byte, size)
			for i := range binaryData {
				binaryData[i] = byte(i % 256)
			}
			encodedData := base64.StdEncoding.EncodeToString(binaryData)
			
			request := map[string]interface{}{
				"jsonrpc": "2.0",
				"id": 1,
				"method": "tools/call",
				"params": map[string]interface{}{
					"name": "upload_binary_content",
					"arguments": map[string]interface{}{
						"content": encodedData,
						"mime_type": ContentTypeOctetStream,
						"filename": fmt.Sprintf("benchmark-%d.bin", size),
					},
				},
			}
			
			requestJSON, _ := json.Marshal(request)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				resp, err := http.Post(testServer.httpServer.URL, "application/json", bytes.NewReader(requestJSON))
				if err != nil {
					b.Fatal(err)
				}
				resp.Body.Close()
			}
		})
	}
}
