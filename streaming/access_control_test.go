package streaming

import (
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

func TestCanAccessWorld(t *testing.T) {
	creatorID := "creator123"
	allowedUserID := "user456"
	unrelatedUserID := "stranger789"

	// Test cases
	testCases := []struct {
		name     string
		userID   string
		moment   *models.WorldMoment
		expected bool
	}{
		{
			name:   "Creator can access private world",
			userID: creatorID,
			moment: &models.WorldMoment{
				WorldID:   "world1",
				CreatorID: creatorID,
				Sharing: models.SharingSettings{
					IsPublic:     false,
					AllowedUsers: []string{},
				},
			},
			expected: true,
		},
		{
			name:   "Allowed user can access shared world",
			userID: allowedUserID,
			moment: &models.WorldMoment{
				WorldID:   "world2",
				CreatorID: creatorID,
				Sharing: models.SharingSettings{
					IsPublic:     false,
					AllowedUsers: []string{allowedUserID, "otheruser"},
				},
			},
			expected: true,
		},
		{
			name:   "Anyone can access public world",
			userID: unrelatedUserID,
			moment: &models.WorldMoment{
				WorldID:   "world3",
				CreatorID: creatorID,
				Sharing: models.SharingSettings{
					IsPublic:     true,
					AllowedUsers: []string{},
				},
			},
			expected: true,
		},
		{
			name:   "Unrelated user cannot access private world",
			userID: unrelatedUserID,
			moment: &models.WorldMoment{
				WorldID:   "world4",
				CreatorID: creatorID,
				Sharing: models.SharingSettings{
					IsPublic:     false,
					AllowedUsers: []string{},
				},
			},
			expected: false,
		},
		{
			name:   "Unrelated user cannot access shared world if not in allowed list",
			userID: unrelatedUserID,
			moment: &models.WorldMoment{
				WorldID:   "world5",
				CreatorID: creatorID,
				Sharing: models.SharingSettings{
					IsPublic:     false,
					AllowedUsers: []string{allowedUserID},
				},
			},
			expected: false,
		},
		{
			name:   "Default private when sharing settings are empty",
			userID: unrelatedUserID,
			moment: &models.WorldMoment{
				WorldID:   "world6",
				CreatorID: creatorID,
				Sharing:   models.SharingSettings{},
			},
			expected: false,
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CanAccessWorld(tc.userID, tc.moment)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetAccessibleContent(t *testing.T) {
	creatorID := "creator123"
	allowedUserID := "user456"
	unrelatedUserID := "stranger789"

	tempVal := 22.5
	lightVal := 0.8
	movementVal := 0.9
	vibeSensorData := models.SensorData{
		Temperature: &tempVal,
		Light:       &lightVal,
		Movement:    &movementVal,
	}

	testVibe := &models.Vibe{
		ID:         "vibe1",
		Name:       "Test Vibe",
		SensorData: vibeSensorData,
	}

	tempVal2 := 25.5
	humidityVal := 0.7
	worldSensorData := models.SensorData{
		Temperature: &tempVal2,
		Humidity:    &humidityVal,
	}

	customData := `{"secretKey": "12345", "privateInfo": "Confidential data"}`

	// Test cases
	testCases := []struct {
		name           string
		userID         string
		moment         *models.WorldMoment
		expected       *models.WorldMoment
		expectNil      bool
		checkVibe      bool
		checkSensor    bool
		checkCustom    bool
		expectCustom   bool
		expectVibeData bool
	}{
		{
			name:   "Creator gets full access",
			userID: creatorID,
			moment: &models.WorldMoment{
				WorldID:    "world1",
				CreatorID:  creatorID,
				CustomData: customData,
				SensorData: worldSensorData,
				Vibe:       testVibe,
				Sharing: models.SharingSettings{
					IsPublic:     false,
					ContextLevel: models.ContextLevelFull,
				},
			},
			expectNil:      false,
			checkVibe:      true,
			checkSensor:    true,
			checkCustom:    true,
			expectCustom:   true,
			expectVibeData: true,
		},
		{
			name:   "No access returns nil",
			userID: unrelatedUserID,
			moment: &models.WorldMoment{
				WorldID:   "world2",
				CreatorID: creatorID,
				Sharing: models.SharingSettings{
					IsPublic:     false,
					ContextLevel: models.ContextLevelFull,
				},
			},
			expectNil: true,
		},
		{
			name:   "Full context level gives all data",
			userID: allowedUserID,
			moment: &models.WorldMoment{
				WorldID:    "world3",
				CreatorID:  creatorID,
				CustomData: customData,
				SensorData: worldSensorData,
				Vibe:       testVibe,
				Sharing: models.SharingSettings{
					IsPublic:     false,
					AllowedUsers: []string{allowedUserID},
					ContextLevel: models.ContextLevelFull,
				},
			},
			expectNil:      false,
			checkVibe:      true,
			checkSensor:    true,
			checkCustom:    true,
			expectCustom:   true,
			expectVibeData: true,
		},
		{
			name:   "Partial context level removes custom data",
			userID: allowedUserID,
			moment: &models.WorldMoment{
				WorldID:    "world4",
				CreatorID:  creatorID,
				CustomData: customData,
				SensorData: worldSensorData,
				Vibe:       testVibe,
				Sharing: models.SharingSettings{
					IsPublic:     false,
					AllowedUsers: []string{allowedUserID},
					ContextLevel: models.ContextLevelPartial,
				},
			},
			expectNil:      false,
			checkVibe:      true,
			checkSensor:    true,
			checkCustom:    true,
			expectCustom:   false,
			expectVibeData: true,
		},
		{
			name:   "None context level removes most data",
			userID: allowedUserID,
			moment: &models.WorldMoment{
				WorldID:    "world5",
				CreatorID:  creatorID,
				CustomData: customData,
				SensorData: worldSensorData,
				Vibe:       testVibe,
				Sharing: models.SharingSettings{
					IsPublic:     false,
					AllowedUsers: []string{allowedUserID},
					ContextLevel: models.ContextLevelNone,
				},
			},
			expectNil:      false,
			checkVibe:      true,
			checkSensor:    true,
			checkCustom:    true,
			expectCustom:   false,
			expectVibeData: false,
		},
		{
			name:   "Public world with no context level gives minimal data",
			userID: unrelatedUserID,
			moment: &models.WorldMoment{
				WorldID:    "world6",
				CreatorID:  creatorID,
				CustomData: customData,
				SensorData: worldSensorData,
				Vibe:       testVibe,
				Sharing: models.SharingSettings{
					IsPublic:     true,
					ContextLevel: models.ContextLevelNone,
				},
			},
			expectNil:      false,
			checkVibe:      true,
			checkSensor:    true,
			checkCustom:    true,
			expectCustom:   false,
			expectVibeData: false,
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetAccessibleContent(tc.userID, tc.moment)
			
			if tc.expectNil {
				assert.Nil(t, result)
				return
			}
			
			assert.NotNil(t, result)
			assert.Equal(t, tc.moment.WorldID, result.WorldID)
			assert.Equal(t, tc.moment.CreatorID, result.CreatorID)
			
			if tc.checkCustom {
				if tc.expectCustom {
					assert.Equal(t, tc.moment.CustomData, result.CustomData)
				} else {
					assert.Empty(t, result.CustomData)
				}
			}
			
			if tc.checkSensor {
				if tc.moment.Sharing.ContextLevel == models.ContextLevelNone {
					assert.Empty(t, result.SensorData)
				} else {
					assert.Equal(t, tc.moment.SensorData, result.SensorData)
				}
			}
			
			if tc.checkVibe {
				assert.NotNil(t, result.Vibe)
				if tc.expectVibeData {
					assert.Equal(t, tc.moment.Vibe.SensorData, result.Vibe.SensorData)
				} else {
					assert.Empty(t, result.Vibe.SensorData)
				}
			}
		})
	}
}