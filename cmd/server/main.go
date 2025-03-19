package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/bmorphism/vibespace-mcp-go/rpcmethods"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
)

const (
	serverPort         = 8080
	startupMessageVibe = "Experience running at http://localhost:8080 - Ready to vibe!"
)

func main() {
	// Create a repository
	repo := repository.NewInMemoryRepository()

	// Add some initial vibes
	addInitialVibes(repo)

	// Add some initial worlds
	addInitialWorlds(repo)

	// Set up NATS streaming configuration
	streamingConfig := &streaming.StreamingConfig{
		NATSHost:       "nonlocal.info",
		NATSPort:       4222,
		StreamID:       "preworm",
		StreamInterval: 5 * time.Second,
		AutoStart:      false,
	}

	// Start the streaming service
	streamingService := streaming.NewStreamingService(repo, streamingConfig)

	// Set up streaming tools
	streamingTools := streaming.NewStreamingTools(streamingService, streamingConfig)

	// Configure MCP server
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", serverPort),
	}

	// Create the handler
	handler := rpcmethods.CreateMCPRequestHandler(repo, streamingTools)

	// Set the handler
	server.Handler = handler

	// Start the server
	fmt.Println(startupMessageVibe)
	log.Fatal(server.ListenAndServe())
}

func addInitialVibes(repo repository.Repository) {
	// Add some pre-configured vibes
	vibes := []*models.Vibe{
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
		_, err := repo.CreateVibe(vibe)
		if err != nil {
			fmt.Printf("Error creating vibe %s: %v\n", vibe.Name, err)
		}
	}
}

func addInitialWorlds(repo repository.Repository) {
	// Add some pre-configured worlds
	worlds := []*models.World{
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
		_, err := repo.CreateWorld(world)
		if err != nil {
			fmt.Printf("Error creating world %s: %v\n", world.Name, err)
		}
	}
}