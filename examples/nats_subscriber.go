package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	url := "nats://nonlocal.info:4222"
	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	// Set default stream ID
	streamID := "ies"
	if len(os.Args) > 2 {
		streamID = os.Args[2]
	}

	// Optional user ID
	userID := ""
	if len(os.Args) > 3 {
		userID = os.Args[3]
	}

	fmt.Printf("Connecting to NATS server at %s with Stream ID '%s'...\n", url, streamID)
	nc, err := nats.Connect(url,
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("Disconnected from NATS: %v\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("Reconnected to NATS server %s\n", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	fmt.Println("Connected to NATS server successfully!")
	
	// Subscribe to all world moments (namespaced with stream ID)
	worldSubject := fmt.Sprintf("%s.world.moment.*", streamID)
	fmt.Printf("Subscribing to %s and %s.world.vibe.* subjects...\n", worldSubject, streamID)

	// Subscribe to world moments
	_, err = nc.Subscribe(worldSubject, func(msg *nats.Msg) {
		var moment models.WorldMoment
		if err := json.Unmarshal(msg.Data, &moment); err != nil {
			fmt.Printf("Error unmarshaling world moment: %v\n", err)
			return
		}

		// Extract components from subject ({streamID}.world.moment.{worldID}[.user.{userID}])
		parts := strings.Split(msg.Subject, ".")

		// Check if this is a user-specific message
		isUserSpecific := len(parts) > 5 && parts[4] == "user"

		// Convert timestamp from milliseconds to time.Time for display
		timestamp := time.Unix(0, moment.Timestamp * int64(time.Millisecond))
		
		fmt.Printf("Received world moment for %s at %v:\n", 
			moment.WorldID, 
			timestamp.Format(time.RFC3339),
		)
		fmt.Printf("  - Subject: %s\n", msg.Subject)
		if isUserSpecific {
			fmt.Printf("  - User-specific message for: %s\n", parts[5])
		} else {
			fmt.Printf("  - Public message\n")
		}
		fmt.Printf("  - Creator ID: %s\n", moment.CreatorID)
		fmt.Printf("  - Occupancy: %d\n", moment.Occupancy)
		fmt.Printf("  - Activity: %.2f\n", moment.Activity)
		
		// Display sharing information
		fmt.Printf("  - Sharing: Public=%v, Context=%s\n", 
			moment.Sharing.IsPublic,
			moment.Sharing.ContextLevel,
		)
		
		// Display viewers
		fmt.Printf("  - Viewers: %d\n", len(moment.Viewers))
		if len(moment.Viewers) > 0 {
			fmt.Printf("    - Viewers list: %v\n", moment.Viewers)
		}

		if moment.Vibe != nil {
			fmt.Printf("  - Current Vibe: %s (Energy: %.2f, Mood: %s)\n", 
				moment.Vibe.Name, 
				moment.Vibe.Energy,
				moment.Vibe.Mood,
			)
		} else {
			fmt.Println("  - No current vibe")
		}

		if moment.SensorData.Temperature != nil {
			fmt.Printf("  - Temperature: %.2f\n", *moment.SensorData.Temperature)
		}
		if moment.SensorData.Humidity != nil {
			fmt.Printf("  - Humidity: %.2f\n", *moment.SensorData.Humidity)
		}
		if moment.SensorData.Light != nil {
			fmt.Printf("  - Light: %.2f\n", *moment.SensorData.Light)
		}
		if moment.SensorData.Sound != nil {
			fmt.Printf("  - Sound: %.2f\n", *moment.SensorData.Sound)
		}
		if moment.SensorData.Movement != nil {
			fmt.Printf("  - Movement: %.2f\n", *moment.SensorData.Movement)
		}

		fmt.Println("-------------------------------------------")
	})
	if err != nil {
		log.Fatalf("Error subscribing to world moments: %v", err)
	}

	// Subscribe to vibe updates
	vibeSubject := fmt.Sprintf("%s.world.vibe.*", streamID)
	_, err = nc.Subscribe(vibeSubject, func(msg *nats.Msg) {
		var vibe models.Vibe
		if err := json.Unmarshal(msg.Data, &vibe); err != nil {
			fmt.Printf("Error unmarshaling vibe: %v\n", err)
			return
		}

		// Extract world ID from subject ({streamID}.world.vibe.{worldID})
		parts := strings.Split(msg.Subject, ".")
		worldID := "unknown"
		if len(parts) > 3 {
			worldID = parts[3]
		}

		fmt.Printf("Vibe update for world %s:\n", worldID)
		fmt.Printf("  - Subject: %s\n", msg.Subject)
		fmt.Printf("  - Vibe: %s\n", vibe.Name)
		fmt.Printf("  - Description: %s\n", vibe.Description)
		fmt.Printf("  - Energy: %.2f\n", vibe.Energy)
		fmt.Printf("  - Mood: %s\n", vibe.Mood)
		fmt.Printf("  - Colors: %v\n", vibe.Colors)

		if vibe.SensorData.Temperature != nil {
			fmt.Printf("  - Target Temperature: %.2f\n", *vibe.SensorData.Temperature)
		}
		if vibe.SensorData.Humidity != nil {
			fmt.Printf("  - Target Humidity: %.2f\n", *vibe.SensorData.Humidity)
		}
		if vibe.SensorData.Light != nil {
			fmt.Printf("  - Target Light: %.2f\n", *vibe.SensorData.Light)
		}
		if vibe.SensorData.Sound != nil {
			fmt.Printf("  - Target Sound: %.2f\n", *vibe.SensorData.Sound)
		}
		if vibe.SensorData.Movement != nil {
			fmt.Printf("  - Target Movement: %.2f\n", *vibe.SensorData.Movement)
		}

		fmt.Println("-------------------------------------------")
	})
	if err != nil {
		log.Fatalf("Error subscribing to vibe updates: %v", err)
	}

	// Subscribe to user-specific streams if userID is provided
	if userID != "" {
		userSubject := fmt.Sprintf("%s.world.moment.*.user.%s", streamID, userID)
		fmt.Printf("Subscribing to user-specific messages: %s\n", userSubject)
		
		_, err = nc.Subscribe(userSubject, func(msg *nats.Msg) {
			var moment models.WorldMoment
			if err := json.Unmarshal(msg.Data, &moment); err != nil {
				fmt.Printf("Error unmarshaling world moment: %v\n", err)
				return
			}
			
			// Convert timestamp from milliseconds to time.Time for display
			timestamp := time.Unix(0, moment.Timestamp * int64(time.Millisecond))
			
			fmt.Printf("📬 USER-SPECIFIC MESSAGE FOR %s:\n", userID)
			fmt.Printf("Received world moment for %s at %v:\n", 
				moment.WorldID, 
				timestamp.Format(time.RFC3339),
			)
			fmt.Printf("  - Subject: %s\n", msg.Subject)
			fmt.Printf("  - Creator ID: %s\n", moment.CreatorID)
			
			// Check if this is your own world
			if moment.CreatorID == userID {
				fmt.Println("  - YOUR WORLD (you are the creator)")
			}
			
			// Check if you're viewing the world
			isViewer := false
			for _, viewer := range moment.Viewers {
				if viewer == userID {
					isViewer = true
					break
				}
			}
			
			if isViewer {
				fmt.Println("  - You are actively viewing this world")
			}
			
			// Display sharing context level
			fmt.Printf("  - Sharing context level: %s\n", moment.Sharing.ContextLevel)
			
			// Display more details
			fmt.Printf("  - Occupancy: %d\n", moment.Occupancy)
			fmt.Printf("  - Activity: %.2f\n", moment.Activity)
			
			if moment.Vibe != nil {
				fmt.Printf("  - Current Vibe: %s (Energy: %.2f, Mood: %s)\n", 
					moment.Vibe.Name, 
					moment.Vibe.Energy,
					moment.Vibe.Mood,
				)
			}
			
			fmt.Println("-------------------------------------------")
		})
		if err != nil {
			log.Fatalf("Error subscribing to user-specific streams: %v", err)
		}
	}

	fmt.Println("Subscriber is running. Press Ctrl+C to exit.")

	// Wait for interruption signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\nExiting...")
}