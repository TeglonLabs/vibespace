package tests

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/mark3labs/mcp-go/mcp"
)

// TestJourney1 tests the basic resource access journey
func TestJourney1(t *testing.T) {
	// JSON-RPC compatibility issue resolved with custom implementation
	
	s, _ := setupTestServer()
	ctx := context.Background()

	// Step 1: Read list of all vibes
	vibeListRes, err := readResource(ctx, s, "vibe://list")
	if err != nil {
		t.Fatalf("Journey 1 - Step 1: Error reading vibe list: %v", err)
	}

	var vibes []models.Vibe
	textContent := vibeListRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &vibes)
	if err != nil {
		t.Fatalf("Journey 1 - Step 1: Error unmarshaling vibe list: %v", err)
	}

	if len(vibes) < 3 {
		t.Errorf("Journey 1 - Step 1: Expected at least 3 vibes, got %d", len(vibes))
	}

	// Step 2: Read a specific vibe
	vibeRes, err := readResource(ctx, s, "vibe://focused-flow")
	if err != nil {
		t.Fatalf("Journey 1 - Step 2: Error reading specific vibe: %v", err)
	}

	var vibe models.Vibe
	textContent = vibeRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &vibe)
	if err != nil {
		t.Fatalf("Journey 1 - Step 2: Error unmarshaling vibe: %v", err)
	}

	if vibe.ID != "focused-flow" {
		t.Errorf("Journey 1 - Step 2: Expected vibe ID 'focused-flow', got '%s'", vibe.ID)
	}

	// Step 3: Read list of all worlds
	worldListRes, err := readResource(ctx, s, "world://list")
	if err != nil {
		t.Fatalf("Journey 1 - Step 3: Error reading world list: %v", err)
	}

	var worlds []models.World
	textContent = worldListRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &worlds)
	if err != nil {
		t.Fatalf("Journey 1 - Step 3: Error unmarshaling world list: %v", err)
	}

	if len(worlds) < 3 {
		t.Errorf("Journey 1 - Step 3: Expected at least 3 worlds, got %d", len(worlds))
	}

	// Step 4: Read a specific world
	worldRes, err := readResource(ctx, s, "world://office-space")
	if err != nil {
		t.Fatalf("Journey 1 - Step 4: Error reading specific world: %v", err)
	}

	var world models.World
	textContent = worldRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &world)
	if err != nil {
		t.Fatalf("Journey 1 - Step 4: Error unmarshaling world: %v", err)
	}

	if world.ID != "office-space" {
		t.Errorf("Journey 1 - Step 4: Expected world ID 'office-space', got '%s'", world.ID)
	}

	// Step 5: Check a world's vibe
	worldVibeRes, err := readResource(ctx, s, "world://office-space/vibe")
	if err != nil {
		t.Fatalf("Journey 1 - Step 5: Error reading world's vibe: %v", err)
	}

	var worldVibe models.Vibe
	textContent = worldVibeRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &worldVibe)
	if err != nil {
		t.Fatalf("Journey 1 - Step 5: Error unmarshaling world vibe: %v", err)
	}

	if worldVibe.ID != "focused-flow" {
		t.Errorf("Journey 1 - Step 5: Expected vibe ID 'focused-flow', got '%s'", worldVibe.ID)
	}
}

// TestJourney2 tests creating and managing vibes
func TestJourney2(t *testing.T) {
	// JSON-RPC compatibility issue resolved with custom implementation
	
	s, _ := setupTestServer()
	ctx := context.Background()

	// Step 1: Create a new vibe
	createVibeArgs := map[string]interface{}{
		"id":          "journey-test-vibe",
		"name":        "Journey Test Vibe",
		"description": "A vibe for journey testing",
		"energy":      0.6,
		"mood":        "adventurous",
		"colors":      []interface{}{"#123456", "#789ABC", "#DEF012"},
	}

	createVibeRes, err := callTool(ctx, s, "create_vibe", createVibeArgs)
	if err != nil {
		t.Fatalf("Journey 2 - Step 1: Error creating vibe: %v", err)
	}

	if createVibeRes.IsError {
		t.Errorf("Journey 2 - Step 1: Unexpected error creating vibe: %s", createVibeRes.Text)
	}

	// Step 2: Read the vibe to verify it was created
	vibeRes, err := readResource(ctx, s, "vibe://journey-test-vibe")
	if err != nil {
		t.Fatalf("Journey 2 - Step 2: Error reading created vibe: %v", err)
	}

	var vibe models.Vibe
	textContent := vibeRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &vibe)
	if err != nil {
		t.Fatalf("Journey 2 - Step 2: Error unmarshaling vibe: %v", err)
	}

	if vibe.ID != "journey-test-vibe" || vibe.Name != "Journey Test Vibe" {
		t.Errorf("Journey 2 - Step 2: Vibe data mismatch. Expected ID 'journey-test-vibe' and name 'Journey Test Vibe', got ID '%s' and name '%s'", vibe.ID, vibe.Name)
	}

	// Step 3: Update the vibe
	updateVibeArgs := map[string]interface{}{
		"id":     "journey-test-vibe",
		"name":   "Updated Journey Vibe",
		"energy": 0.8,
		"mood":   "excited",
	}

	updateVibeRes, err := callTool(ctx, s, "update_vibe", updateVibeArgs)
	if err != nil {
		t.Fatalf("Journey 2 - Step 3: Error updating vibe: %v", err)
	}

	if updateVibeRes.IsError {
		t.Errorf("Journey 2 - Step 3: Unexpected error updating vibe: %s", updateVibeRes.Text)
	}

	// Step 4: Verify the changes
	vibeRes, err = readResource(ctx, s, "vibe://journey-test-vibe")
	if err != nil {
		t.Fatalf("Journey 2 - Step 4: Error reading updated vibe: %v", err)
	}

	textContent = vibeRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &vibe)
	if err != nil {
		t.Fatalf("Journey 2 - Step 4: Error unmarshaling vibe: %v", err)
	}

	if vibe.Name != "Updated Journey Vibe" || vibe.Energy != 0.8 || vibe.Mood != "excited" {
		t.Errorf("Journey 2 - Step 4: Vibe update verification failed. Expected name 'Updated Journey Vibe', energy 0.8, mood 'excited', got name '%s', energy %f, mood '%s'", vibe.Name, vibe.Energy, vibe.Mood)
	}

	// Step 5: Delete the vibe
	deleteVibeArgs := map[string]interface{}{
		"id": "journey-test-vibe",
	}

	deleteVibeRes, err := callTool(ctx, s, "delete_vibe", deleteVibeArgs)
	if err != nil {
		t.Fatalf("Journey 2 - Step 5: Error deleting vibe: %v", err)
	}

	if deleteVibeRes.IsError {
		t.Errorf("Journey 2 - Step 5: Unexpected error deleting vibe: %s", deleteVibeRes.Text)
	}

	// Verify the vibe was deleted
	_, err = readResource(ctx, s, "vibe://journey-test-vibe")
	if err == nil {
		t.Errorf("Journey 2 - Step 5: Vibe still exists after deletion")
	}
}

// TestJourney3 tests creating and managing worlds
func TestJourney3(t *testing.T) {
	// JSON-RPC compatibility issue resolved with custom implementation
	
	s, _ := setupTestServer()
	ctx := context.Background()

	// Step 1: Create a new world
	createWorldArgs := map[string]interface{}{
		"id":          "journey-test-world",
		"name":        "Journey Test World",
		"description": "A world for journey testing",
		"type":        "VIRTUAL",
	}

	createWorldRes, err := callTool(ctx, s, "create_world", createWorldArgs)
	if err != nil {
		t.Fatalf("Journey 3 - Step 1: Error creating world: %v", err)
	}

	if createWorldRes.IsError {
		t.Errorf("Journey 3 - Step 1: Unexpected error creating world: %s", createWorldRes.Text)
	}

	// Step 2: Read the world to verify it was created
	worldRes, err := readResource(ctx, s, "world://journey-test-world")
	if err != nil {
		t.Fatalf("Journey 3 - Step 2: Error reading created world: %v", err)
	}

	var world models.World
	textContent := worldRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &world)
	if err != nil {
		t.Fatalf("Journey 3 - Step 2: Error unmarshaling world: %v", err)
	}

	if world.ID != "journey-test-world" || world.Name != "Journey Test World" {
		t.Errorf("Journey 3 - Step 2: World data mismatch. Expected ID 'journey-test-world' and name 'Journey Test World', got ID '%s' and name '%s'", world.ID, world.Name)
	}

	// Step 3: Set the world's vibe
	setVibeArgs := map[string]interface{}{
		"worldId": "journey-test-world",
		"vibeId":  "calm-clarity",
	}

	setVibeRes, err := callTool(ctx, s, "set_world_vibe", setVibeArgs)
	if err != nil {
		t.Fatalf("Journey 3 - Step 3: Error setting world vibe: %v", err)
	}

	if setVibeRes.IsError {
		t.Errorf("Journey 3 - Step 3: Unexpected error setting world vibe: %s", setVibeRes.Text)
	}

	// Step 4: Verify the vibe was set
	worldVibeRes, err := readResource(ctx, s, "world://journey-test-world/vibe")
	if err != nil {
		t.Fatalf("Journey 3 - Step 4: Error reading world vibe: %v", err)
	}

	var vibe models.Vibe
	textContent = worldVibeRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &vibe)
	if err != nil {
		t.Fatalf("Journey 3 - Step 4: Error unmarshaling world vibe: %v", err)
	}

	if vibe.ID != "calm-clarity" {
		t.Errorf("Journey 3 - Step 4: World vibe mismatch. Expected 'calm-clarity', got '%s'", vibe.ID)
	}

	// Step 5: Update the world
	updateWorldArgs := map[string]interface{}{
		"id":          "journey-test-world",
		"name":        "Updated Journey World",
		"description": "An updated world for journey testing",
		"location":    "https://updated.journey.world",
	}

	updateWorldRes, err := callTool(ctx, s, "update_world", updateWorldArgs)
	if err != nil {
		t.Fatalf("Journey 3 - Step 5: Error updating world: %v", err)
	}

	if updateWorldRes.IsError {
		t.Errorf("Journey 3 - Step 5: Unexpected error updating world: %s", updateWorldRes.Text)
	}

	// Verify the world was updated
	worldRes, err = readResource(ctx, s, "world://journey-test-world")
	if err != nil {
		t.Fatalf("Journey 3 - Step 5: Error reading updated world: %v", err)
	}

	textContent = worldRes[0].(mcp.TextResourceContents)
	err = json.Unmarshal([]byte(textContent.Text), &world)
	if err != nil {
		t.Fatalf("Journey 3 - Step 5: Error unmarshaling world: %v", err)
	}

	if world.Name != "Updated Journey World" || world.Location != "https://updated.journey.world" {
		t.Errorf("Journey 3 - Step 5: World update verification failed. Expected name 'Updated Journey World' and location 'https://updated.journey.world', got name '%s' and location '%s'", world.Name, world.Location)
	}

	// Step 6: Delete the world
	deleteWorldArgs := map[string]interface{}{
		"id": "journey-test-world",
	}

	deleteWorldRes, err := callTool(ctx, s, "delete_world", deleteWorldArgs)
	if err != nil {
		t.Fatalf("Journey 3 - Step 6: Error deleting world: %v", err)
	}

	if deleteWorldRes.IsError {
		t.Errorf("Journey 3 - Step 6: Unexpected error deleting world: %s", deleteWorldRes.Text)
	}

	// Verify the world was deleted
	_, err = readResource(ctx, s, "world://journey-test-world")
	if err == nil {
		t.Errorf("Journey 3 - Step 6: World still exists after deletion")
	}
}

// TestJourney4 tests error handling
func TestJourney4(t *testing.T) {
	// JSON-RPC compatibility issue resolved with custom implementation
	
	s, _ := setupTestServer()
	ctx := context.Background()

	// Step 1: Try to get a non-existent vibe
	_, err := readResource(ctx, s, "vibe://non-existent-vibe")
	if err == nil {
		t.Errorf("Journey 4 - Step 1: Expected error for non-existent vibe, got none")
	}

	// Step 2: Try to update a non-existent world
	updateWorldArgs := map[string]interface{}{
		"id":   "non-existent-world",
		"name": "This world doesn't exist",
	}

	updateWorldRes, err := callTool(ctx, s, "update_world", updateWorldArgs)
	if err != nil {
		t.Fatalf("Journey 4 - Step 2: Unexpected error calling update_world: %v", err)
	}

	if !updateWorldRes.IsError {
		t.Errorf("Journey 4 - Step 2: Expected error for updating non-existent world, got none")
	}

	// Step 3: Try to create a world with a reference to a non-existent vibe
	createWorldArgs := map[string]interface{}{
		"id":          "bad-world",
		"name":        "Bad World",
		"description": "A world with a non-existent vibe",
		"type":        "VIRTUAL",
		"currentVibe": "non-existent-vibe",
	}

	createWorldRes, err := callTool(ctx, s, "create_world", createWorldArgs)
	if err != nil {
		t.Fatalf("Journey 4 - Step 3: Unexpected error calling create_world: %v", err)
	}

	if !createWorldRes.IsError {
		t.Errorf("Journey 4 - Step 3: Expected error for creating world with non-existent vibe, got none")
	}

	// Step 4: Try to delete a vibe that's currently assigned to a world
	deleteVibeArgs := map[string]interface{}{
		"id": "focused-flow", // This is used by a sample world
	}

	deleteVibeRes, err := callTool(ctx, s, "delete_vibe", deleteVibeArgs)
	if err != nil {
		t.Fatalf("Journey 4 - Step 4: Unexpected error calling delete_vibe: %v", err)
	}

	if !deleteVibeRes.IsError {
		t.Errorf("Journey 4 - Step 4: Expected error for deleting vibe in use, got none")
	}
}