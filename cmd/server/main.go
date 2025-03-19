package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/bmorphism/vibespace-mcp-go/rpcmethods"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Pre-compiled regular expressions for better performance
var (
	vibeURIRegex       = regexp.MustCompile(models.VibeScheme + `(.+)`)
	worldURIRegex      = regexp.MustCompile(models.WorldScheme + `(.+)`)
	worldVibeURIRegex  = regexp.MustCompile(`(.*)` + models.WorldVibeSubURI + `$`)
)

func main() {
	// Create repository with sample data
	repo := repository.NewRepository()

	// Create streaming service with nonlocal.info NATS server
	streamingConfig := &streaming.StreamingConfig{
		NATSHost:       "nonlocal.info", // NATS server host
		NATSPort:       4222,            // NATS server port
		StreamID:       "ies",           // Default stream ID
		StreamInterval: 5 * time.Second, // Default interval of 5 seconds
		AutoStart:      false,           // Don't auto-start streaming
	}
	streamingService := streaming.NewStreamingService(repo, streamingConfig)
	
	// Create streaming tools
	streamingTools := streaming.NewStreamingTools(streamingService)

	// Create MCP server
	s := server.NewMCPServer("vibespace", "1.0.0")

	// Set up resource handlers
	setupVibeResources(s, repo)
	setupWorldResources(s, repo)

	// Set up tools
	setupVibeTools(s, repo)
	setupWorldTools(s, repo)
	setupStreamingTools(s, streamingTools)

	// Initialize the streaming service
	if err := streamingService.Start(); err != nil {
		log.Printf("Warning: Failed to initialize streaming service: %v\n", err)
		log.Println("Continuing without streaming functionality")
	}

	// Start the server with our improved method handling
	fmt.Println("Starting vibespace MCP experience...")
	
	// Wrap the server to provide better method handling
	wrapper := rpcmethods.WrapMCPServer(s)
	
	// Use a custom context function that wraps the message handler
	ctxFunc := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "messageHandler", wrapper)
	}
	
	// Create custom options
	options := []server.StdioOption{
		server.WithStdioContextFunc(ctxFunc),
	}
	
	// Connect the server to stdio
	if err := server.ServeStdio(s, options...); err != nil {
		// Ensure streaming service is stopped when server exits
		streamingService.Stop()
		log.Fatalf("Server error: %v\n", err)
	}
	
	// Ensure streaming service is stopped when server exits
	streamingService.Stop()
}

// Helper functions for resources and tools
func WithDescription(description string) mcp.ToolOption {
	return func(t *mcp.Tool) {
		t.Description = description
	}
}

func WithResourceDescription(resource mcp.Resource, description string) mcp.Resource {
	opts := []mcp.ResourceOption{mcp.WithResourceDescription(description)}
	for _, opt := range opts {
		opt(&resource)
	}
	return resource
}

func WithMIMEType(resource mcp.Resource, mimeType string) mcp.Resource {
	opts := []mcp.ResourceOption{mcp.WithMIMEType(mimeType)}
	for _, opt := range opts {
		opt(&resource)
	}
	return resource
}

func WithTemplateDescription(template mcp.ResourceTemplate, description string) mcp.ResourceTemplate {
	opts := []mcp.ResourceTemplateOption{mcp.WithTemplateDescription(description)}
	for _, opt := range opts {
		opt(&template)
	}
	return template
}

func WithTemplateMIMEType(template mcp.ResourceTemplate, mimeType string) mcp.ResourceTemplate {
	opts := []mcp.ResourceTemplateOption{mcp.WithTemplateMIMEType(mimeType)}
	for _, opt := range opts {
		opt(&template)
	}
	return template
}

// setupVibeResources adds resources for vibes
func setupVibeResources(s *server.MCPServer, repo *repository.Repository) {
	// Add resource for listing vibes
	vibeListResource := mcp.NewResource(models.VibeListURI, "List of vibes")
	vibeListResource = WithResourceDescription(vibeListResource, "List all available vibes")
	vibeListResource = WithMIMEType(vibeListResource, "application/json")
	
	s.AddResource(vibeListResource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		vibes := repo.GetAllVibes()
		vibesJSON, err := json.MarshalIndent(vibes, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error marshaling vibes list: %w", err)
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      models.VibeListURI,
				MIMEType: "application/json",
				Text:     string(vibesJSON),
			},
		}, nil
	})

	// Add resource template for specific vibes
	vibeTemplate := mcp.NewResourceTemplate(models.VibeScheme + "{id}", "Get vibe")
	vibeTemplate = WithTemplateDescription(vibeTemplate, "Get a specific vibe by ID")
	vibeTemplate = WithTemplateMIMEType(vibeTemplate, "application/json")
	
	s.AddResourceTemplate(vibeTemplate, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract vibe ID from URI using pre-compiled regex
		matches := vibeURIRegex.FindStringSubmatch(req.Params.URI)
		if len(matches) < 2 {
			return nil, fmt.Errorf("invalid vibe URI format: %s", req.Params.URI)
		}

		vibeID := matches[1]
		vibe, err := repo.GetVibe(vibeID)
		if err != nil {
			return nil, err
		}

		vibeJSON, err := json.MarshalIndent(vibe, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error marshaling vibe data: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(vibeJSON),
			},
		}, nil
	})
}

// setupWorldResources adds resources for worlds
func setupWorldResources(s *server.MCPServer, repo *repository.Repository) {
	// Add resources for worlds
	worldListResource := mcp.NewResource(models.WorldListURI, "List of worlds")
	worldListResource = WithResourceDescription(worldListResource, "List all available worlds")
	worldListResource = WithMIMEType(worldListResource, "application/json")
	
	s.AddResource(worldListResource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		worlds := repo.GetAllWorlds()
		worldsJSON, err := json.MarshalIndent(worlds, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error marshaling worlds list: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      models.WorldListURI,
				MIMEType: "application/json",
				Text:     string(worldsJSON),
			},
		}, nil
	})

	// Add resource template for specific worlds
	worldTemplate := mcp.NewResourceTemplate(models.WorldScheme + "{id}", "Get world")
	worldTemplate = WithTemplateDescription(worldTemplate, "Get a specific world by ID")
	worldTemplate = WithTemplateMIMEType(worldTemplate, "application/json")
	
	s.AddResourceTemplate(worldTemplate, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract world ID from URI using pre-compiled regex
		matches := worldURIRegex.FindStringSubmatch(req.Params.URI)
		if len(matches) < 2 {
			return nil, fmt.Errorf("invalid world URI format: %s", req.Params.URI)
		}

		worldID := matches[1]
		
		// If the URI includes "/vibe", handle it as a world's vibe request
		if worldVibe := worldVibeURIRegex.FindStringSubmatch(worldID); len(worldVibe) == 2 {
			actualWorldID := worldVibe[1]
			vibe, err := repo.GetWorldVibe(actualWorldID)
			if err != nil {
				return nil, fmt.Errorf("error retrieving world vibe: %w", err)
			}
			
			vibeJSON, err := json.MarshalIndent(vibe, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("error marshaling world vibe data: %w", err)
			}
			
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(vibeJSON),
				},
			}, nil
		}
		
		// Regular world request
		world, err := repo.GetWorld(worldID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving world: %w", err)
		}

		worldJSON, err := json.MarshalIndent(world, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error marshaling world data: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(worldJSON),
			},
		}, nil
	})
}

// setupVibeTools sets up tools for manipulating vibes
func setupVibeTools(s *server.MCPServer, repo *repository.Repository) {
	// Create a new vibe
	createVibe := mcp.NewTool("create_vibe", WithDescription("Create a new vibe with the specified properties"))

	s.AddTool(createVibe, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Safely extract parameters with type checking
		idVal, ok := req.Params.Arguments["id"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: id"), nil
		}
		id, ok := idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'id': expected string"), nil
		}

		nameVal, ok := req.Params.Arguments["name"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: name"), nil
		}
		name, ok := nameVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'name': expected string"), nil
		}

		descriptionVal, ok := req.Params.Arguments["description"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: description"), nil
		}
		description, ok := descriptionVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'description': expected string"), nil
		}

		energyVal, ok := req.Params.Arguments["energy"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: energy"), nil
		}
		energy, ok := energyVal.(float64)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'energy': expected float64"), nil
		}

		moodVal, ok := req.Params.Arguments["mood"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: mood"), nil
		}
		mood, ok := moodVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'mood': expected string"), nil
		}

		// Convert colors from interface{} array to string array with safety checks
		colorsVal, ok := req.Params.Arguments["colors"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: colors"), nil
		}
		colorsInterface, ok := colorsVal.([]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'colors': expected array"), nil
		}
		
		colors := make([]string, len(colorsInterface))
		for i, c := range colorsInterface {
			colorStr, ok := c.(string)
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid type for color at index %d: expected string", i)), nil
			}
			colors[i] = colorStr
		}

		// Create vibe object
		vibe := models.Vibe{
			ID:          id,
			Name:        name,
			Description: description,
			Energy:      energy,
			Mood:        mood,
			Colors:      colors,
		}

		// Handle optional sensor data if provided
		if sensorDataRaw, ok := req.Params.Arguments["sensorData"]; ok && sensorDataRaw != nil {
			sensorData := models.SensorData{}
			sensorDataMap, ok := sensorDataRaw.(map[string]interface{})
			if !ok {
				return mcp.NewToolResultError("Invalid type for parameter 'sensorData': expected object"), nil
			}

			if temp, ok := sensorDataMap["temperature"]; ok && temp != nil {
				tempFloat, ok := temp.(float64)
				if !ok {
					return mcp.NewToolResultError("Invalid type for 'temperature': expected float64"), nil
				}
				sensorData.Temperature = &tempFloat
			}
			if humidity, ok := sensorDataMap["humidity"]; ok && humidity != nil {
				humidityFloat, ok := humidity.(float64)
				if !ok {
					return mcp.NewToolResultError("Invalid type for 'humidity': expected float64"), nil
				}
				sensorData.Humidity = &humidityFloat
			}
			if light, ok := sensorDataMap["light"]; ok && light != nil {
				lightFloat, ok := light.(float64)
				if !ok {
					return mcp.NewToolResultError("Invalid type for 'light': expected float64"), nil
				}
				sensorData.Light = &lightFloat
			}
			if sound, ok := sensorDataMap["sound"]; ok && sound != nil {
				soundFloat, ok := sound.(float64)
				if !ok {
					return mcp.NewToolResultError("Invalid type for 'sound': expected float64"), nil
				}
				sensorData.Sound = &soundFloat
			}
			if movement, ok := sensorDataMap["movement"]; ok && movement != nil {
				movementFloat, ok := movement.(float64)
				if !ok {
					return mcp.NewToolResultError("Invalid type for 'movement': expected float64"), nil
				}
				sensorData.Movement = &movementFloat
			}

			vibe.SensorData = sensorData
		}

		// Check if vibe already exists
		_, err := repo.GetVibe(id)
		if err == nil {
			return mcp.NewToolResultError(fmt.Sprintf("Vibe with ID '%s' already exists", id)), nil
		} else if err != repository.ErrVibeNotFound {
			return mcp.NewToolResultError(fmt.Sprintf("Error checking vibe existence: %v", err)), nil
		}

		// Add vibe to repository
		if err := repo.AddVibe(vibe); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create vibe: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully created vibe '%s'", name)), nil
	})

	// Update an existing vibe
	updateVibe := mcp.NewTool("update_vibe", WithDescription("Update an existing vibe's properties"))

	s.AddTool(updateVibe, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract vibe ID
		id := req.Params.Arguments["id"].(string)

		// Get existing vibe
		existingVibe, err := repo.GetVibe(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Vibe with ID '%s' not found", id)), nil
		}

		// Update fields if provided
		if name, ok := req.Params.Arguments["name"].(string); ok {
			existingVibe.Name = name
		}
		if description, ok := req.Params.Arguments["description"].(string); ok {
			existingVibe.Description = description
		}
		if energy, ok := req.Params.Arguments["energy"].(float64); ok {
			existingVibe.Energy = energy
		}
		if mood, ok := req.Params.Arguments["mood"].(string); ok {
			existingVibe.Mood = mood
		}

		// Update colors if provided
		if colorsInterface, ok := req.Params.Arguments["colors"].([]interface{}); ok {
			colors := make([]string, len(colorsInterface))
			for i, c := range colorsInterface {
				colors[i] = c.(string)
			}
			existingVibe.Colors = colors
		}

		// Update sensor data if provided
		if sensorDataRaw, ok := req.Params.Arguments["sensorData"]; ok && sensorDataRaw != nil {
			sensorData := models.SensorData{}
			sensorDataMap := sensorDataRaw.(map[string]interface{})

			// Keep existing values if new ones are not provided
			if existingVibe.SensorData.Temperature != nil {
				sensorData.Temperature = existingVibe.SensorData.Temperature
			}
			if existingVibe.SensorData.Humidity != nil {
				sensorData.Humidity = existingVibe.SensorData.Humidity
			}
			if existingVibe.SensorData.Light != nil {
				sensorData.Light = existingVibe.SensorData.Light
			}
			if existingVibe.SensorData.Sound != nil {
				sensorData.Sound = existingVibe.SensorData.Sound
			}
			if existingVibe.SensorData.Movement != nil {
				sensorData.Movement = existingVibe.SensorData.Movement
			}

			// Update with new values
			if temp, ok := sensorDataMap["temperature"]; ok && temp != nil {
				tempFloat := temp.(float64)
				sensorData.Temperature = &tempFloat
			}
			if humidity, ok := sensorDataMap["humidity"]; ok && humidity != nil {
				humidityFloat := humidity.(float64)
				sensorData.Humidity = &humidityFloat
			}
			if light, ok := sensorDataMap["light"]; ok && light != nil {
				lightFloat := light.(float64)
				sensorData.Light = &lightFloat
			}
			if sound, ok := sensorDataMap["sound"]; ok && sound != nil {
				soundFloat := sound.(float64)
				sensorData.Sound = &soundFloat
			}
			if movement, ok := sensorDataMap["movement"]; ok && movement != nil {
				movementFloat := movement.(float64)
				sensorData.Movement = &movementFloat
			}

			existingVibe.SensorData = sensorData
		}

		// Update vibe in repository
		if err := repo.UpdateVibe(existingVibe); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update vibe: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully updated vibe '%s'", existingVibe.Name)), nil
	})

	// Delete a vibe
	deleteVibe := mcp.NewTool("delete_vibe", WithDescription("Delete a vibe by ID"))

	s.AddTool(deleteVibe, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Safely extract ID parameter
		idVal, ok := req.Params.Arguments["id"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: id"), nil
		}
		id, ok := idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'id': expected string"), nil
		}

		// Try to get the vibe to ensure it exists
		vibe, err := repo.GetVibe(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Vibe with ID '%s' not found", id)), nil
		}

		// Delete the vibe
		if err := repo.DeleteVibe(id); err != nil {
			if err == repository.ErrVibeInUse {
				// Find out which worlds are using this vibe
				worlds := repo.GetAllWorlds()
				worldsUsingVibe := []string{}
				
				for _, world := range worlds {
					if world.CurrentVibe == id {
						worldsUsingVibe = append(worldsUsingVibe, world.ID)
					}
				}
				
				return mcp.NewToolResultError(fmt.Sprintf("Cannot delete vibe '%s'. It's currently used by worlds: %v", 
					vibe.Name, worldsUsingVibe)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete vibe: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted vibe '%s'", vibe.Name)), nil
	})
}

// setupStreamingTools sets up tools for controlling the streaming service
func setupStreamingTools(s *server.MCPServer, tools *streaming.StreamingTools) {
	// Register all streaming tools
	
	// Start streaming
	startStreaming := mcp.NewTool("streaming_startStreaming", WithDescription("Start streaming world moments to NATS"))
	s.AddTool(startStreaming, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract interval parameter if provided
		intervalVal, ok := req.Params.Arguments["interval"]
		var interval int
		if ok {
			intervalFloat, ok := intervalVal.(float64)
			if !ok {
				return mcp.NewToolResultError("Invalid type for parameter 'interval': expected number"), nil
			}
			interval = int(intervalFloat)
		}
		
		// Call the streaming tool method
		startReq := &streaming.StartStreamingRequest{Interval: interval}
		result, err := tools.StartStreaming(startReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error starting streaming: %v", err)), nil
		}
		
		if !result.Success {
			return mcp.NewToolResultError(result.Message), nil
		}
		
		return mcp.NewToolResultText(result.Message), nil
	})
	
	// Stop streaming
	stopStreaming := mcp.NewTool("streaming_stopStreaming", WithDescription("Stop streaming world moments"))
	s.AddTool(stopStreaming, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := tools.StopStreaming()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error stopping streaming: %v", err)), nil
		}
		
		if !result.Success {
			return mcp.NewToolResultError(result.Message), nil
		}
		
		return mcp.NewToolResultText(result.Message), nil
	})
	
	// Get streaming status
	streamingStatus := mcp.NewTool("streaming_status", WithDescription("Get current status of streaming service"))
	s.AddTool(streamingStatus, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		status, err := tools.Status()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error getting streaming status: %v", err)), nil
		}
		
		// Format response
		statusJSON, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error formatting status: %v", err)), nil
		}
		
		return mcp.NewToolResultText(string(statusJSON)), nil
	})
	
	// Stream a specific world
	streamWorld := mcp.NewTool("streaming_streamWorld", WithDescription("Stream a moment for a specific world"))
	s.AddTool(streamWorld, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract world ID parameter
		worldIDVal, ok := req.Params.Arguments["worldId"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: worldId"), nil
		}
		worldID, ok := worldIDVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'worldId': expected string"), nil
		}
		
		// Extract user ID parameter
		userIDVal, ok := req.Params.Arguments["userId"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: userId"), nil
		}
		userID, ok := userIDVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'userId': expected string"), nil
		}
		
		// Extract sharing settings if provided
		streamReq := &streaming.StreamWorldRequest{
			WorldID: worldID,
			UserID:  userID,
		}
		
		// Handle sharing settings if provided
		if sharingVal, ok := req.Params.Arguments["sharing"]; ok && sharingVal != nil {
			sharingMap, ok := sharingVal.(map[string]interface{})
			if !ok {
				return mcp.NewToolResultError("Invalid type for parameter 'sharing': expected object"), nil
			}
			
			sharing := &streaming.SharingRequest{}
			
			// Parse IsPublic
			if isPublicVal, ok := sharingMap["isPublic"]; ok {
				isPublic, ok := isPublicVal.(bool)
				if !ok {
					return mcp.NewToolResultError("Invalid type for 'sharing.isPublic': expected boolean"), nil
				}
				sharing.IsPublic = isPublic
			}
			
			// Parse AllowedUsers
			if allowedUsersVal, ok := sharingMap["allowedUsers"]; ok && allowedUsersVal != nil {
				allowedUsersInterface, ok := allowedUsersVal.([]interface{})
				if !ok {
					return mcp.NewToolResultError("Invalid type for 'sharing.allowedUsers': expected array"), nil
				}
				
				allowedUsers := make([]string, len(allowedUsersInterface))
				for i, u := range allowedUsersInterface {
					userStr, ok := u.(string)
					if !ok {
						return mcp.NewToolResultError(fmt.Sprintf("Invalid type for user at index %d: expected string", i)), nil
					}
					allowedUsers[i] = userStr
				}
				sharing.AllowedUsers = allowedUsers
			}
			
			// Parse ContextLevel
			if contextLevelVal, ok := sharingMap["contextLevel"]; ok {
				contextLevel, ok := contextLevelVal.(string)
				if !ok {
					return mcp.NewToolResultError("Invalid type for 'sharing.contextLevel': expected string"), nil
				}
				sharing.ContextLevel = contextLevel
			}
			
			streamReq.Sharing = sharing
		}
		
		// Call the streaming tool method
		result, err := tools.StreamWorld(streamReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error streaming world: %v", err)), nil
		}
		
		if !result.Success {
			return mcp.NewToolResultError(result.Message), nil
		}
		
		return mcp.NewToolResultText(result.Message), nil
	})
	
	// Update streaming configuration
	updateConfig := mcp.NewTool("streaming_updateConfig", WithDescription("Update streaming configuration"))
	s.AddTool(updateConfig, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters
		configReq := &streaming.UpdateConfigRequest{}
		
		if natsURLVal, ok := req.Params.Arguments["natsUrl"]; ok && natsURLVal != nil {
			natsURL, ok := natsURLVal.(string)
			if !ok {
				return mcp.NewToolResultError("Invalid type for parameter 'natsUrl': expected string"), nil
			}
			configReq.NATSUrl = natsURL
		}
		
		if intervalVal, ok := req.Params.Arguments["streamInterval"]; ok && intervalVal != nil {
			intervalFloat, ok := intervalVal.(float64)
			if !ok {
				return mcp.NewToolResultError("Invalid type for parameter 'streamInterval': expected number"), nil
			}
			configReq.StreamInterval = int(intervalFloat)
		}
		
		// Call the streaming tool method
		result, err := tools.UpdateConfig(configReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error updating configuration: %v", err)), nil
		}
		
		if !result.Success {
			return mcp.NewToolResultError(result.Message), nil
		}
		
		return mcp.NewToolResultText(result.Message), nil
	})
}

func setupWorldTools(s *server.MCPServer, repo *repository.Repository) {
	// Create a new world
	createWorld := mcp.NewTool("create_world", WithDescription("Create a new world with the specified properties"))

	s.AddTool(createWorld, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters with type safety
		idVal, ok := req.Params.Arguments["id"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: id"), nil
		}
		id, ok := idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'id': expected string"), nil
		}

		nameVal, ok := req.Params.Arguments["name"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: name"), nil
		}
		name, ok := nameVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'name': expected string"), nil
		}

		descriptionVal, ok := req.Params.Arguments["description"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: description"), nil
		}
		description, ok := descriptionVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'description': expected string"), nil
		}

		typeVal, ok := req.Params.Arguments["type"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: type"), nil
		}
		typeStr, ok := typeVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid type for parameter 'type': expected string"), nil
		}
		
		// Validate world type
		worldType := models.WorldType(typeStr)
		if worldType != models.WorldTypePhysical && 
		   worldType != models.WorldTypeVirtual && 
		   worldType != models.WorldTypeHybrid {
			return mcp.NewToolResultError("Invalid world type: must be 'PHYSICAL', 'VIRTUAL', or 'HYBRID'"), nil
		}

		// Create world object
		world := models.World{
			ID:          id,
			Name:        name,
			Description: description,
			Type:        worldType,
		}

		// Handle optional fields with type safety
		if locationVal, ok := req.Params.Arguments["location"]; ok && locationVal != nil {
			location, ok := locationVal.(string)
			if !ok {
				return mcp.NewToolResultError("Invalid type for parameter 'location': expected string"), nil
			}
			world.Location = location
		}
		
		if currentVibeVal, ok := req.Params.Arguments["currentVibe"]; ok && currentVibeVal != nil {
			currentVibe, ok := currentVibeVal.(string)
			if !ok {
				return mcp.NewToolResultError("Invalid type for parameter 'currentVibe': expected string"), nil
			}
			world.CurrentVibe = currentVibe
		}
		
		if sizeVal, ok := req.Params.Arguments["size"]; ok && sizeVal != nil {
			size, ok := sizeVal.(string)
			if !ok {
				return mcp.NewToolResultError("Invalid type for parameter 'size': expected string"), nil
			}
			world.Size = size
		}

		// Convert features if provided with type safety
		if featuresVal, ok := req.Params.Arguments["features"]; ok && featuresVal != nil {
			featuresInterface, ok := featuresVal.([]interface{})
			if !ok {
				return mcp.NewToolResultError("Invalid type for parameter 'features': expected array"), nil
			}
			
			features := make([]string, len(featuresInterface))
			for i, f := range featuresInterface {
				featureStr, ok := f.(string)
				if !ok {
					return mcp.NewToolResultError(fmt.Sprintf("Invalid type for feature at index %d: expected string", i)), nil
				}
				features[i] = featureStr
			}
			world.Features = features
		}

		// Check if world already exists
		_, err := repo.GetWorld(id)
		if err == nil {
			return mcp.NewToolResultError(fmt.Sprintf("World with ID '%s' already exists", id)), nil
		} else if err != repository.ErrWorldNotFound {
			return mcp.NewToolResultError(fmt.Sprintf("Error checking world existence: %v", err)), nil
		}

		// Add world to repository
		if err := repo.AddWorld(world); err != nil {
			if err == repository.ErrVibeNotFound {
				return mcp.NewToolResultError(fmt.Sprintf("Referenced vibe '%s' doesn't exist", world.CurrentVibe)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create world: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully created world '%s'", name)), nil
	})

	// Update an existing world
	updateWorld := mcp.NewTool("update_world", WithDescription("Update an existing world's properties"))

	s.AddTool(updateWorld, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract world ID
		id := req.Params.Arguments["id"].(string)

		// Get existing world
		existingWorld, err := repo.GetWorld(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("World with ID '%s' not found", id)), nil
		}

		// Update fields if provided
		if name, ok := req.Params.Arguments["name"].(string); ok {
			existingWorld.Name = name
		}
		if description, ok := req.Params.Arguments["description"].(string); ok {
			existingWorld.Description = description
		}
		if typeStr, ok := req.Params.Arguments["type"].(string); ok {
			existingWorld.Type = models.WorldType(typeStr)
		}
		if location, ok := req.Params.Arguments["location"].(string); ok {
			existingWorld.Location = location
		}
		if currentVibe, ok := req.Params.Arguments["currentVibe"].(string); ok {
			existingWorld.CurrentVibe = currentVibe
		}
		if size, ok := req.Params.Arguments["size"].(string); ok {
			existingWorld.Size = size
		}

		// Update features if provided
		if featuresInterface, ok := req.Params.Arguments["features"].([]interface{}); ok {
			features := make([]string, len(featuresInterface))
			for i, f := range featuresInterface {
				features[i] = f.(string)
			}
			existingWorld.Features = features
		}

		// Update world in repository
		if err := repo.UpdateWorld(existingWorld); err != nil {
			if err == repository.ErrVibeNotFound {
				return mcp.NewToolResultError(fmt.Sprintf("Referenced vibe '%s' doesn't exist", existingWorld.CurrentVibe)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update world: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully updated world '%s'", existingWorld.Name)), nil
	})

	// Delete a world
	deleteWorld := mcp.NewTool("delete_world", WithDescription("Delete a world by ID"))

	s.AddTool(deleteWorld, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := req.Params.Arguments["id"].(string)

		// Try to get the world to ensure it exists
		world, err := repo.GetWorld(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("World with ID '%s' not found", id)), nil
		}

		// Delete the world
		if err := repo.DeleteWorld(id); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete world: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted world '%s'", world.Name)), nil
	})

	// Set a world's vibe
	setWorldVibe := mcp.NewTool("set_world_vibe", WithDescription("Set a world's vibe"))

	s.AddTool(setWorldVibe, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		worldID := req.Params.Arguments["worldId"].(string)
		vibeID := req.Params.Arguments["vibeId"].(string)

		// Get the world and vibe to ensure they exist
		world, err := repo.GetWorld(worldID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("World with ID '%s' not found", worldID)), nil
		}

		vibe, err := repo.GetVibe(vibeID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Vibe with ID '%s' not found", vibeID)), nil
		}

		// Set the world's vibe
		if err := repo.SetWorldVibe(worldID, vibeID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set world vibe: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully set world '%s' to vibe '%s'", world.Name, vibe.Name)), nil
	})
}