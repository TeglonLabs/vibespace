package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/bmorphism/vibespace-mcp-go/rpcmethods"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverPort         = 8080
	startupMessageVibe = "Experience running at http://localhost:8080 - Ready to vibe!"
)

func main() {
	// Create a repository
	repo := repository.NewRepository()

	// Add some initial vibes
	addInitialVibes(repo)

	// Add some initial worlds
	addInitialWorlds(repo)

	// Set up NATS streaming configuration
	streamingConfig := &streaming.StreamingConfig{
		NATSHost:       "nonlocal.info",
		NATSPort:       4222,
		NATSUrl:        "nats://nonlocal.info:4222",
		StreamID:       "preworm",
		StreamInterval: 5 * time.Second,
		AutoStart:      false,
	}

	// Start the streaming service
	streamingService := streaming.NewStreamingService(repo, streamingConfig)

	// Set up streaming tools
	streamingTools := streaming.NewStreamingTools(streamingService)

	// Create mcp server
	mcpServer := server.NewMCPServer("vibespace-mcp", "1.0.0")

	// Get the streaming tool methods and register them
	fmt.Println("Registering streaming tools:")
	
	// Register streaming tools with the MCPServer
	startStreamingTool := mcp.NewTool("streaming_startStreaming", func(t *mcp.Tool) {
		t.Description = "Start streaming world moments"
	})
	mcpServer.AddTool(startStreamingTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Convert JSON-RPC request to our internal format
		interval := 0
		if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
			if intervalVal, ok := args["interval"]; ok {
				if intervalFloat, ok := intervalVal.(float64); ok {
					interval = int(intervalFloat)
				}
			}
		}
		
		// Call the streaming tool
		response, err := streamingTools.StartStreaming(&streaming.StartStreamingRequest{
			Interval: interval,
		})
		
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		// Convert response to JSON-RPC format
		resultText := response.Message
		if !response.Success {
			return mcp.NewToolResultError(resultText), nil
		}
		
		return mcp.NewToolResultText(resultText), nil
	})
	fmt.Println("  - streaming_startStreaming: Start streaming world moments")
	
	// Register stop streaming tool
	stopStreamingTool := mcp.NewTool("streaming_stopStreaming", func(t *mcp.Tool) {
		t.Description = "Stop streaming world moments"
	})
	mcpServer.AddTool(stopStreamingTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		response, err := streamingTools.StopStreaming()
		
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		resultText := response.Message
		if !response.Success {
			return mcp.NewToolResultError(resultText), nil
		}
		
		return mcp.NewToolResultText(resultText), nil
	})
	fmt.Println("  - streaming_stopStreaming: Stop streaming world moments")
	
	// Register status tool
	statusTool := mcp.NewTool("streaming_status", func(t *mcp.Tool) {
		t.Description = "Get current streaming status"
	})
	mcpServer.AddTool(statusTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		status, err := streamingTools.Status()
		
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		// Convert to JSON
		statusJSON, err := json.Marshal(status)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error marshaling status: %v", err)), nil
		}
		
		return mcp.NewToolResultText(string(statusJSON)), nil
	})
	fmt.Println("  - streaming_status: Get current streaming status")
	
	// Register streamWorld tool
	streamWorldTool := mcp.NewTool("streaming_streamWorld", func(t *mcp.Tool) {
		t.Description = "Stream a single world moment"
	})
	mcpServer.AddTool(streamWorldTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		var worldID, userID string
		if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
			worldID, _ = args["worldId"].(string)
			userID, _ = args["userId"].(string)
		}
		
		// Create request
		streamReq := &streaming.StreamWorldRequest{
			WorldID: worldID,
			UserID:  userID,
		}
		
		// Check for sharing settings
		if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
			if sharingMap, ok := args["sharing"].(map[string]interface{}); ok {
				sharing := &streaming.SharingRequest{}
			
			// Extract sharing fields
			if isPublic, ok := sharingMap["isPublic"].(bool); ok {
				sharing.IsPublic = isPublic
			}
			
			if contextLevel, ok := sharingMap["contextLevel"].(string); ok {
				sharing.ContextLevel = contextLevel
			}
			
			if allowedUsersInterface, ok := sharingMap["allowedUsers"].([]interface{}); ok {
				allowedUsers := make([]string, len(allowedUsersInterface))
				for i, u := range allowedUsersInterface {
					if userStr, ok := u.(string); ok {
						allowedUsers[i] = userStr
					}
				}
				sharing.AllowedUsers = allowedUsers
			}
			
			streamReq.Sharing = sharing
			}
		}
		// Call the streaming tool
		response, err := streamingTools.StreamWorld(streamReq)
		
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		resultText := response.Message
		if !response.Success {
			return mcp.NewToolResultError(resultText), nil
		}
		
		return mcp.NewToolResultText(resultText), nil
	})
	fmt.Println("  - streaming_streamWorld: Stream a single world moment")
	
	// Register updateConfig tool
	updateConfigTool := mcp.NewTool("streaming_updateConfig", func(t *mcp.Tool) {
		t.Description = "Update streaming configuration"
	})
	mcpServer.AddTool(updateConfigTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		config := &streaming.UpdateConfigRequest{}
		
		if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
			if natsHost, ok := args["natsHost"].(string); ok {
				config.NATSHost = natsHost
			}
			
			if natsPort, ok := args["natsPort"].(float64); ok {
				config.NATSPort = int(natsPort)
			}
			
			if natsURL, ok := args["natsUrl"].(string); ok {
				config.NATSUrl = natsURL
			}
			
			if streamID, ok := args["streamId"].(string); ok {
				config.StreamID = streamID
			}
			
			if streamInterval, ok := args["streamInterval"].(float64); ok {
				config.StreamInterval = int(streamInterval)
			}
		}
		
		// Call the streaming tool
		response, err := streamingTools.UpdateConfig(config)
		
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		resultText := response.Message
		if !response.Success {
			return mcp.NewToolResultError(resultText), nil
		}
		
		return mcp.NewToolResultText(resultText), nil
	})
	fmt.Println("  - streaming_updateConfig: Update streaming configuration")

	// Register resource handlers from server wrapper
	handler := rpcmethods.WrapMCPServer(mcpServer)

	// Configure HTTP server
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", serverPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handle MCP RPC requests
			if r.Method == http.MethodPost {
				// Read the request body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
					return
				}
				defer r.Body.Close()
				
				// Process the request
				response := handler.HandleMessage(r.Context(), body)
				
				// Marshal the response
				responseJSON, err := json.Marshal(response)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error marshaling response: %v", err), http.StatusInternalServerError)
					return
				}
				
				// Set content type and write response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(responseJSON)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method not allowed"))
			}
		}),
	}

	// Start the server
	fmt.Println(startupMessageVibe)
	log.Fatal(httpServer.ListenAndServe())
}

func addInitialVibes(repo repository.VibeRepository) {
	// Add some pre-configured vibes
	vibes := []models.Vibe{
		{
			ID:          "calm",
			Name:        "Calm",
			Description: "A peaceful and serene vibe",
			Energy:      0.3,
			Mood:        models.MoodCalm,
			Colors:      []string{"#6A98DC", "#B5CAE8", "#DAF1F9"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelFull},
		},
		{
			ID:          "focus",
			Name:        "Focused",
			Description: "A concentrated and productive vibe",
			Energy:      0.6,
			Mood:        models.MoodFocused,
			Colors:      []string{"#2D3E50", "#34495E", "#5D6D7E"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelFull},
		},
		{
			ID:          "energetic",
			Name:        "Energetic",
			Description: "A high-energy, vibrant atmosphere",
			Energy:      0.9,
			Mood:        models.MoodEnergetic,
			Colors:      []string{"#F39C12", "#E74C3C", "#9B59B6"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelFull},
		},
		{
			ID:          "creative",
			Name:        "Creative",
			Description: "An inspiring and imaginative vibe",
			Energy:      0.7,
			Mood:        models.MoodCreative,
			Colors:      []string{"#1ABC9C", "#3498DB", "#F1C40F"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelFull},
		},
		{
			ID:          "contemplative",
			Name:        "Contemplative",
			Description: "A thoughtful, reflective atmosphere",
			Energy:      0.4,
			Mood:        models.MoodContemplative,
			Colors:      []string{"#8E44AD", "#2C3E50", "#34495E"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelFull},
		},
	}

	for _, vibe := range vibes {
		err := repo.AddVibe(vibe)
		if err != nil {
			fmt.Printf("Error creating vibe %s: %v\n", vibe.Name, err)
		}
	}
}

func addInitialWorlds(repo repository.WorldRepository) {
	// Add some pre-configured worlds
	worlds := []models.World{
		{
			ID:          "office",
			Name:        "Office Space",
			Description: "A modern collaborative workspace",
			Type:        models.WorldTypePhysical,
			Location:    "Building A, Floor 2",
			CurrentVibe: "focus",
			Features:    []string{"standing desks", "natural light", "sound dampening"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelPartial},
		},
		{
			ID:          "home-office",
			Name:        "Home Office",
			Description: "A comfortable work-from-home setup",
			Type:        models.WorldTypePhysical,
			Location:    "Home",
			CurrentVibe: "calm",
			Features:    []string{"ergonomic chair", "plant", "coffee machine"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelPartial},
		},
		{
			ID:          "virtual-cafe",
			Name:        "Virtual Café",
			Description: "A digital space with café ambiance",
			Type:        models.WorldTypeVirtual,
			CurrentVibe: "creative",
			Features:    []string{"ambient sounds", "customizable decor", "shared whiteboard"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelFull},
		},
		{
			ID:          "conference-room",
			Name:        "Conference Room",
			Description: "A hybrid meeting space for team collaboration",
			Type:        models.WorldTypeHybrid,
			Location:    "Building A, Conference Room 3",
			CurrentVibe: "focus",
			Features:    []string{"video conferencing", "digital whiteboard", "sound system"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelPartial},
		},
		{
			ID:          "study-lounge",
			Name:        "Study Lounge",
			Description: "A quiet space for focused learning",
			Type:        models.WorldTypePhysical,
			Location:    "Library, 3rd Floor",
			CurrentVibe: "contemplative",
			Features:    []string{"bookshelf", "individual desks", "natural light"},
			CreatorID:   "system",
			Sharing:     models.SharingSettings{IsPublic: true, ContextLevel: models.ContextLevelPartial},
		},
	}

	for _, world := range worlds {
		err := repo.AddWorld(world)
		if err != nil {
			fmt.Printf("Error creating world %s: %v\n", world.Name, err)
		}
	}
}