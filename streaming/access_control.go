package streaming

import (
	"github.com/bmorphism/vibespace-mcp-go/models"
)

// CanAccessWorld determines if a user has permission to access a world moment
func CanAccessWorld(userID string, moment *models.WorldMoment) bool {
	// If the moment has no sharing settings defined (empty allowed users and not public), default to private
	if !moment.Sharing.IsPublic && len(moment.Sharing.AllowedUsers) == 0 {
		// Only the creator can access
		return userID == moment.CreatorID
	}
	
	// If the world is public, anyone can access
	if moment.Sharing.IsPublic {
		return true
	}
	
	// If the user is the creator, they always have access
	if userID == moment.CreatorID {
		return true
	}
	
	// Check if the user is in the allowed users list
	for _, allowedUser := range moment.Sharing.AllowedUsers {
		if userID == allowedUser {
			return true
		}
	}
	
	// Default to no access
	return false
}

// GetAccessibleContent filters the content of a WorldMoment based on the user's
// permissions and the context level specified in the sharing settings
func GetAccessibleContent(userID string, moment *models.WorldMoment) *models.WorldMoment {
	// If the user doesn't have access, return nil
	if !CanAccessWorld(userID, moment) {
		return nil
	}
	
	// If the user is the creator, they get full access
	if userID == moment.CreatorID {
		return moment
	}
	
	// Make a copy to avoid modifying the original
	result := *moment
	
	// Apply context level filtering
	switch moment.Sharing.ContextLevel {
	case models.ContextLevelNone:
		// Provide minimal information - just ID, type and public metadata
		result.CustomData = ""
		result.SensorData = models.SensorData{}
		// Remove binary and ternary data at none level
		result.BinaryData = nil
		result.BalancedTernaryData = nil
		if result.Vibe != nil {
			// Remove vibe details but keep basic info
			vibe := *result.Vibe
			vibe.SensorData = models.SensorData{}
			result.Vibe = &vibe
		}
		
	case models.ContextLevelPartial:
		// Provide moderate information but not custom/sensitive data
		result.CustomData = ""
		// May keep binary data depending on format but strip any sensitive formats
		if result.BinaryData != nil {
			// Keep only non-sensitive formats
			if result.BinaryData.Format == "application/octet-stream" || 
			   result.BinaryData.Format == "application/binary" {
				result.BinaryData = nil
			}
		}
		// Keep balanced ternary data
		
	case models.ContextLevelFull:
		// Provide all information - no changes needed
	}
	
	return &result
}