package tests

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Helpers for server options
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

// Register JSON-RPC handlers for resources and tools
// This function is not currently used due to compatibility issues
func registerJSONRPCHandlers(s interface{}) {
	// This is a placeholder function that will be updated in the future
	// when we have a way to register custom JSON-RPC handlers
}

// setupTestServer creates a test server with the repository
func setupTestServer() (*server.MCPServer, *repository.Repository) {
	// Use empty repository without sample data for tests
	repo := repository.NewRepositoryWithSampleData(false)
	s := server.NewMCPServer("vibespace-test", "1.0.0")
	
	// Set up resources
	setupTestVibeResources(s, repo)
	setupTestWorldResources(s, repo)
	
	// Set up tools
	setupTestVibeTools(s, repo)
	setupTestWorldTools(s, repo)
	
	return s, repo
}

// Copy of setupVibeResources from main.go for testing
func setupTestVibeResources(s *server.MCPServer, repo *repository.Repository) {
	// Add resource for listing vibes
	vibeListResource := mcp.NewResource("vibe://list", "List of vibes")
	vibeListResource = WithResourceDescription(vibeListResource, "List all available vibes")
	vibeListResource = WithMIMEType(vibeListResource, "application/json")
	
	s.AddResource(vibeListResource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		vibes := repo.GetAllVibes()
		vibesJSON, _ := json.MarshalIndent(vibes, "", "  ")
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "vibe://list",
				MIMEType: "application/json",
				Text:     string(vibesJSON),
			},
		}, nil
	})

	// Add resource template for specific vibes
	vibeTemplate := mcp.NewResourceTemplate("vibe://{id}", "Get vibe")
	vibeTemplate = WithTemplateDescription(vibeTemplate, "Get a specific vibe by ID")
	vibeTemplate = WithTemplateMIMEType(vibeTemplate, "application/json")
	
	s.AddResourceTemplate(vibeTemplate, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Simple regex implementation for testing
		re := regexp.MustCompile(`vibe://(.+)`)
		matches := re.FindStringSubmatch(req.Params.URI)
		if len(matches) < 2 {
			return nil, json.Unmarshal([]byte(`{"error":"Invalid vibe URI"}`), nil)
		}
	
		vibeID := matches[1]
		vibe, err := repo.GetVibe(vibeID)
		if err != nil {
			return nil, err
		}
		vibeJSON, _ := json.MarshalIndent(vibe, "", "  ")
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(vibeJSON),
			},
		}, nil
	})
}

// Copy of setupWorldResources from main.go for testing
func setupTestWorldResources(s *server.MCPServer, repo *repository.Repository) {
	// Add resources for worlds
	worldListResource := mcp.NewResource("world://list", "List of worlds")
	worldListResource = WithResourceDescription(worldListResource, "List all available worlds")
	worldListResource = WithMIMEType(worldListResource, "application/json")
	
	s.AddResource(worldListResource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		worlds := repo.GetAllWorlds()
		worldsJSON, _ := json.MarshalIndent(worlds, "", "  ")
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "world://list",
				MIMEType: "application/json",
				Text:     string(worldsJSON),
			},
		}, nil
	})

	// Add resource template for specific worlds
	worldTemplate := mcp.NewResourceTemplate("world://{id}", "Get world")
	worldTemplate = WithTemplateDescription(worldTemplate, "Get a specific world by ID")
	worldTemplate = WithTemplateMIMEType(worldTemplate, "application/json")
	
	s.AddResourceTemplate(worldTemplate, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Simple implementation for testing
		re := regexp.MustCompile(`world://(.+)`)
		matches := re.FindStringSubmatch(req.Params.URI)
		if len(matches) < 2 {
			return nil, json.Unmarshal([]byte(`{"error":"Invalid world URI"}`), nil)
		}
		
		worldID := matches[1]
		// If the URI includes "/vibe", handle it as a world's vibe request
		if re := regexp.MustCompile(`(.*)/vibe$`); re.MatchString(worldID) {
			matches := re.FindStringSubmatch(worldID)
			if len(matches) < 2 {
				return nil, json.Unmarshal([]byte(`{"error":"Invalid world/vibe URI"}`), nil)
			}
			
			actualWorldID := matches[1]
			vibe, err := repo.GetWorldVibe(actualWorldID)
			if err != nil {
				return nil, err
			}
			vibeJSON, _ := json.MarshalIndent(vibe, "", "  ")
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(vibeJSON),
				},
			}, nil
		}
		
		world, err := repo.GetWorld(worldID)
		if err != nil {
			return nil, err
		}
		worldJSON, _ := json.MarshalIndent(world, "", "  ")
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(worldJSON),
			},
		}, nil
	})
}

// Simplified versions of the tool handlers for testing
func setupTestVibeTools(s *server.MCPServer, repo *repository.Repository) {
	// Create a new vibe
	createVibe := mcp.NewTool("create_vibe", WithDescription("Create a new vibe with the specified properties"))

	s.AddTool(createVibe, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := req.Params.Arguments["id"].(string)
		name := req.Params.Arguments["name"].(string)
		description := req.Params.Arguments["description"].(string)
		energy := req.Params.Arguments["energy"].(float64)
		mood := req.Params.Arguments["mood"].(string)
		
		colorsInterface := req.Params.Arguments["colors"].([]interface{})
		colors := make([]string, len(colorsInterface))
		for i, c := range colorsInterface {
			colors[i] = c.(string)
		}
		
		vibe := models.Vibe{
			ID:          id,
			Name:        name,
			Description: description,
			Energy:      energy,
			Mood:        mood,
			Colors:      colors,
		}
		
		err := repo.AddVibe(vibe)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		return mcp.NewToolResultText("Vibe created successfully"), nil
	})
	
	// Delete a vibe
	deleteVibe := mcp.NewTool("delete_vibe", WithDescription("Delete a vibe by ID"))
	
	s.AddTool(deleteVibe, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := req.Params.Arguments["id"].(string)
		err := repo.DeleteVibe(id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		return mcp.NewToolResultText("Vibe deleted successfully"), nil
	})
	
	// Update a vibe
	updateVibe := mcp.NewTool("update_vibe", WithDescription("Update an existing vibe's properties"))
	
	s.AddTool(updateVibe, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := req.Params.Arguments["id"].(string)
		
		// Get existing vibe
		vibe, err := repo.GetVibe(id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		// Update fields if provided
		if name, ok := req.Params.Arguments["name"].(string); ok {
			vibe.Name = name
		}
		if description, ok := req.Params.Arguments["description"].(string); ok {
			vibe.Description = description
		}
		if energy, ok := req.Params.Arguments["energy"].(float64); ok {
			vibe.Energy = energy
		}
		if mood, ok := req.Params.Arguments["mood"].(string); ok {
			vibe.Mood = mood
		}
		
		err = repo.UpdateVibe(vibe)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		return mcp.NewToolResultText("Vibe updated successfully"), nil
	})
}

func setupTestWorldTools(s *server.MCPServer, repo *repository.Repository) {
	// Create a new world
	createWorld := mcp.NewTool("create_world", WithDescription("Create a new world with the specified properties"))
	
	s.AddTool(createWorld, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := req.Params.Arguments["id"].(string)
		name := req.Params.Arguments["name"].(string)
		description := req.Params.Arguments["description"].(string)
		typeStr := req.Params.Arguments["type"].(string)
		
		world := models.World{
			ID:          id,
			Name:        name,
			Description: description,
			Type:        models.WorldType(typeStr),
		}
		
		if vibeID, ok := req.Params.Arguments["currentVibe"].(string); ok {
			world.CurrentVibe = vibeID
		}
		
		err := repo.AddWorld(world)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		return mcp.NewToolResultText("World created successfully"), nil
	})
	
	// Set a world's vibe
	setWorldVibe := mcp.NewTool("set_world_vibe", WithDescription("Set a world's vibe"))
	
	s.AddTool(setWorldVibe, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		worldID := req.Params.Arguments["worldId"].(string)
		vibeID := req.Params.Arguments["vibeId"].(string)
		
		err := repo.SetWorldVibe(worldID, vibeID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		return mcp.NewToolResultText("World vibe set successfully"), nil
	})
	
	// Update a world
	updateWorld := mcp.NewTool("update_world", WithDescription("Update an existing world's properties"))
	
	s.AddTool(updateWorld, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := req.Params.Arguments["id"].(string)
		
		world, err := repo.GetWorld(id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		if name, ok := req.Params.Arguments["name"].(string); ok {
			world.Name = name
		}
		if description, ok := req.Params.Arguments["description"].(string); ok {
			world.Description = description
		}
		if location, ok := req.Params.Arguments["location"].(string); ok {
			world.Location = location
		}
		
		err = repo.UpdateWorld(world)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		return mcp.NewToolResultText("World updated successfully"), nil
	})
	
	// Delete a world
	deleteWorld := mcp.NewTool("delete_world", WithDescription("Delete a world by ID"))
	
	s.AddTool(deleteWorld, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := req.Params.Arguments["id"].(string)
		
		err := repo.DeleteWorld(id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		
		return mcp.NewToolResultText("World deleted successfully"), nil
	})
}

// Test resource listing
func TestResourceList(t *testing.T) {
	s, _ := setupTestServer()
	ctx := context.Background()
	
	// JSON-RPC compatibility issue resolved with custom implementation
	
	// Test vibe://list resource
	vibeListRes, err := readResource(ctx, s, "vibe://list")
	if err != nil {
		t.Fatalf("Error reading vibe list: %v", err)
	}
	
	if len(vibeListRes) == 0 {
		t.Error("Expected non-empty vibe list response")
	}
	
	// Unmarshal the vibes to check content
	var vibes []models.Vibe
	textContent := vibeListRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &vibes)
	if err != nil {
		t.Fatalf("Error unmarshaling vibe list: %v", err)
	}
	
	if len(vibes) < 3 {
		t.Errorf("Expected at least 3 vibes, got %d", len(vibes))
	}
	
	// Test world://list resource
	worldListRes, err := readResource(ctx, s, "world://list")
	if err != nil {
		t.Fatalf("Error reading world list: %v", err)
	}
	
	if len(worldListRes) == 0 {
		t.Error("Expected non-empty world list response")
	}
	
	// Unmarshal the worlds to check content
	var worlds []models.World
	textContent = worldListRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &worlds)
	if err != nil {
		t.Fatalf("Error unmarshaling world list: %v", err)
	}
	
	if len(worlds) < 3 {
		t.Errorf("Expected at least 3 worlds, got %d", len(worlds))
	}
}

// Test getting resources by ID
func TestResourceById(t *testing.T) {
	s, _ := setupTestServer()
	ctx := context.Background()
	
	// JSON-RPC compatibility issue resolved with custom implementation
	
	// Test vibe://{id}
	vibeRes, err := readResource(ctx, s, "vibe://focused-flow")
	if err != nil {
		t.Fatalf("Error reading vibe: %v", err)
	}
	
	if len(vibeRes) == 0 {
		t.Error("Expected non-empty vibe response")
	}
	
	// Unmarshal the vibe to check content
	var vibe models.Vibe
	textContent := vibeRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &vibe)
	if err != nil {
		t.Fatalf("Error unmarshaling vibe: %v", err)
	}
	
	if vibe.ID != "focused-flow" {
		t.Errorf("Expected vibe ID 'focused-flow', got '%s'", vibe.ID)
	}
	
	// Test invalid vibe ID
	_, err = readResource(ctx, s, "vibe://non-existent")
	if err == nil {
		t.Error("Expected error for non-existent vibe, got nil")
	}
	
	// Test world://{id}
	worldRes, err := readResource(ctx, s, "world://office-space")
	if err != nil {
		t.Fatalf("Error reading world: %v", err)
	}
	
	if len(worldRes) == 0 {
		t.Error("Expected non-empty world response")
	}
	
	// Unmarshal the world to check content
	var world models.World
	textContent = worldRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &world)
	if err != nil {
		t.Fatalf("Error unmarshaling world: %v", err)
	}
	
	if world.ID != "office-space" {
		t.Errorf("Expected world ID 'office-space', got '%s'", world.ID)
	}
	
	// Test world's vibe resource
	worldVibeRes, err := readResource(ctx, s, "world://office-space/vibe")
	if err != nil {
		t.Fatalf("Error reading world vibe: %v", err)
	}
	
	if len(worldVibeRes) == 0 {
		t.Error("Expected non-empty world vibe response")
	}
	
	// Unmarshal the vibe to check content
	textContent = worldVibeRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &vibe)
	if err != nil {
		t.Fatalf("Error unmarshaling world vibe: %v", err)
	}
	
	// The office-space world should have the focused-flow vibe
	if vibe.ID != "focused-flow" {
		t.Errorf("Expected vibe ID 'focused-flow', got '%s'", vibe.ID)
	}
}