package rpcmethods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Repository alias for the interface
type Repository = repository.VibeWorldRepository

// Define URI handlers
type vibeUriHandler struct {
	repo Repository
}

func (h *vibeUriHandler) HandleUri(uri string) (interface{}, error) {
	if uri == models.VibeListURI {
		return h.repo.GetAllVibes(), nil
	}

	if strings.HasPrefix(uri, models.VibeScheme) {
		vibeID := strings.TrimPrefix(uri, models.VibeScheme)
		vibe, err := h.repo.GetVibe(vibeID)
		if err != nil {
			return nil, err
		}
		return vibe, nil
	}

	return nil, fmt.Errorf("invalid vibe URI: %s", uri)
}

type worldUriHandler struct {
	repo Repository
}

func (h *worldUriHandler) HandleUri(uri string) (interface{}, error) {
	if uri == models.WorldListURI {
		return h.repo.GetAllWorlds(), nil
	}

	if strings.HasPrefix(uri, models.WorldScheme) {
		worldURI := strings.TrimPrefix(uri, models.WorldScheme)
		
		// Check if it's a world vibe request
		if strings.HasSuffix(worldURI, models.WorldVibeSubURI) {
			worldID := strings.TrimSuffix(worldURI, models.WorldVibeSubURI)
			vibe, err := h.repo.GetWorldVibe(worldID)
			if err != nil {
				return nil, err
			}
			return vibe, nil
		}
		
		// Regular world request
		world, err := h.repo.GetWorld(worldURI)
		if err != nil {
			return nil, err
		}
		return world, nil
	}

	return nil, fmt.Errorf("invalid world URI: %s", uri)
}

// Define tool handlers
func createVibeTools(repo Repository) map[string]interface{} {
	return map[string]interface{}{
		"create_vibe": func(req json.RawMessage) (interface{}, error) {
			var vibe models.Vibe
			if err := json.Unmarshal(req, &vibe); err != nil {
				return nil, fmt.Errorf("invalid vibe data: %v", err)
			}
			
			if err := repo.AddVibe(vibe); err != nil {
				return nil, err
			}
			
			return map[string]interface{}{
				"success": true,
				"id":      vibe.ID,
				"message": fmt.Sprintf("Vibe '%s' created successfully", vibe.Name),
			}, nil
		},
		"update_vibe": func(req json.RawMessage) (interface{}, error) {
			var vibe models.Vibe
			if err := json.Unmarshal(req, &vibe); err != nil {
				return nil, fmt.Errorf("invalid vibe data: %v", err)
			}
			
			if err := repo.UpdateVibe(vibe); err != nil {
				return nil, err
			}
			
			return map[string]interface{}{
				"success": true,
				"id":      vibe.ID,
				"message": fmt.Sprintf("Vibe '%s' updated successfully", vibe.Name),
			}, nil
		},
		"delete_vibe": func(req json.RawMessage) (interface{}, error) {
			var params struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(req, &params); err != nil {
				return nil, fmt.Errorf("invalid request: %v", err)
			}
			
			if err := repo.DeleteVibe(params.ID); err != nil {
				return nil, err
			}
			
			return map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Vibe with ID '%s' deleted successfully", params.ID),
			}, nil
		},
	}
}

func createWorldTools(repo Repository) map[string]interface{} {
	return map[string]interface{}{
		"create_world": func(req json.RawMessage) (interface{}, error) {
			var world models.World
			if err := json.Unmarshal(req, &world); err != nil {
				return nil, fmt.Errorf("invalid world data: %v", err)
			}
			
			if err := repo.AddWorld(world); err != nil {
				return nil, err
			}
			
			return map[string]interface{}{
				"success": true,
				"id":      world.ID,
				"message": fmt.Sprintf("World '%s' created successfully", world.Name),
			}, nil
		},
		"update_world": func(req json.RawMessage) (interface{}, error) {
			var world models.World
			if err := json.Unmarshal(req, &world); err != nil {
				return nil, fmt.Errorf("invalid world data: %v", err)
			}
			
			if err := repo.UpdateWorld(world); err != nil {
				return nil, err
			}
			
			return map[string]interface{}{
				"success": true,
				"id":      world.ID,
				"message": fmt.Sprintf("World '%s' updated successfully", world.Name),
			}, nil
		},
		"delete_world": func(req json.RawMessage) (interface{}, error) {
			var params struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(req, &params); err != nil {
				return nil, fmt.Errorf("invalid request: %v", err)
			}
			
			if err := repo.DeleteWorld(params.ID); err != nil {
				return nil, err
			}
			
			return map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("World with ID '%s' deleted successfully", params.ID),
			}, nil
		},
		"set_world_vibe": func(req json.RawMessage) (interface{}, error) {
			var params struct {
				WorldID string `json:"worldId"`
				VibeID  string `json:"vibeId"`
			}
			if err := json.Unmarshal(req, &params); err != nil {
				return nil, fmt.Errorf("invalid request: %v", err)
			}
			
			if err := repo.SetWorldVibe(params.WorldID, params.VibeID); err != nil {
				return nil, err
			}
			
			return map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Vibe '%s' set for world '%s'", params.VibeID, params.WorldID),
			}, nil
		},
	}
}

// CreateMCPRequestHandler creates an HTTP handler for MCP requests
func CreateMCPRequestHandler(repo Repository, streamingTools *streaming.StreamingTools) http.Handler {
	// Create MCP server with name and version
	mcpServer := server.NewMCPServer("vibespace-mcp-go", "1.0.0")
	
	// Add resource handlers (equivalent to URI handlers)
	mcpServer.AddResource(mcp.Resource{
		URI:         "vibe://",
		Name:        "vibe",
		Description: "Vibe resource handler",
		MIMEType:    "application/json",
	}, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Convert to vibeUriHandler call
		handler := &vibeUriHandler{repo: repo}
		return handler.HandleRead(ctx, request)
	})
	
	mcpServer.AddResource(mcp.Resource{
		URI:         "world://",
		Name:        "world",
		Description: "World resource handler",
		MIMEType:    "application/json",
	}, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Convert to worldUriHandler call
		handler := &worldUriHandler{repo: repo}
		return handler.HandleRead(ctx, request)
	})
	
	// Add vibe tools
	for name, toolFunc := range createVibeTools(repo) {
		mcpServer.AddTool(mcp.Tool{
			Name:        name,
			Description: fmt.Sprintf("Vibe tool: %s", name),
		}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Convert arguments to JSON first
			var args json.RawMessage
			if request.Params.Arguments != nil {
				if argBytes, err := json.Marshal(request.Params.Arguments); err == nil {
					args = argBytes
				} else {
					return nil, fmt.Errorf("failed to marshal arguments: %v", err)
				}
			}
			
			// Convert the old tool function to new format
			if fn, ok := toolFunc.(func(json.RawMessage) (interface{}, error)); ok {
				result, err := fn(args)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(fmt.Sprintf("%v", result)), nil
			} else {
				return nil, fmt.Errorf("invalid tool function type: %T", toolFunc)
			}
		})
	}
	
	// Add world tools
	for name, toolFunc := range createWorldTools(repo) {
		mcpServer.AddTool(mcp.Tool{
			Name:        name,
			Description: fmt.Sprintf("World tool: %s", name),
		}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Convert arguments to JSON first
			var args json.RawMessage
			if request.Params.Arguments != nil {
				if argBytes, err := json.Marshal(request.Params.Arguments); err == nil {
					args = argBytes
				} else {
					return nil, fmt.Errorf("failed to marshal arguments: %v", err)
				}
			}
			
			// Convert the old tool function to new format
			if fn, ok := toolFunc.(func(json.RawMessage) (interface{}, error)); ok {
				result, err := fn(args)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(fmt.Sprintf("%v", result)), nil
			} else {
				return nil, fmt.Errorf("invalid tool function type: %T", toolFunc)
			}
		})
	}
	
	// Add streaming tools
	for name := range streaming.GetStreamingToolMethods() {
		switch name {
		case "streaming_startStreaming":
			mcpServer.AddTool(mcp.Tool{
				Name:        name,
				Description: "Start streaming",
			}, createStreamingToolHandler(streamingTools.StartStreaming))
		case "streaming_stopStreaming":
			mcpServer.AddTool(mcp.Tool{
				Name:        name,
				Description: "Stop streaming",
			}, createStreamingToolHandler(streamingTools.StopStreaming))
		case "streaming_status":
			mcpServer.AddTool(mcp.Tool{
				Name:        name,
				Description: "Get streaming status",
			}, createStreamingToolHandler(streamingTools.Status))
		case "streaming_streamWorld":
			mcpServer.AddTool(mcp.Tool{
				Name:        name,
				Description: "Stream world",
			}, createStreamingToolHandler(streamingTools.StreamWorld))
		case "streaming_updateConfig":
			mcpServer.AddTool(mcp.Tool{
				Name:        name,
				Description: "Update streaming config",
			}, createStreamingToolHandler(streamingTools.UpdateConfig))
		}
	}
	
	// Create and return the wrapped handler to improve error messages
	wrapper := WrapMCPServer(mcpServer)
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// Read request body
		var body json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}
		
		// Process the request
		response := wrapper.HandleMessage(r.Context(), body)
		
		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

// Helper function to convert URI handlers to resource handlers
func (h *vibeUriHandler) HandleRead(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	result, err := h.HandleUri(request.Params.URI)
	if err != nil {
		return nil, err
	}
	
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

func (h *worldUriHandler) HandleRead(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	result, err := h.HandleUri(request.Params.URI)
	if err != nil {
		return nil, err
	}
	
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

// Helper function to create streaming tool handlers
func createStreamingToolHandler(toolFunc interface{}) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Convert arguments to the format expected by streaming tools
		var result interface{}
		var err error
		
		// Call the tool function based on its type
		switch f := toolFunc.(type) {
		case func() (interface{}, error):
			result, err = f()
		case func(map[string]interface{}) (interface{}, error):
			args := make(map[string]interface{})
			if request.Params.Arguments != nil {
				// Convert arguments to JSON bytes first
				argBytes, err := json.Marshal(request.Params.Arguments)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal arguments: %v", err)
				}
				if err := json.Unmarshal(argBytes, &args); err != nil {
					return nil, fmt.Errorf("failed to unmarshal arguments: %v", err)
				}
			}
			result, err = f(args)
		default:
			return nil, fmt.Errorf("unsupported tool function type: %T", toolFunc)
		}
		
		if err != nil {
			return nil, err
		}
		
		return mcp.NewToolResultText(fmt.Sprintf("%v", result)), nil
	}
}
