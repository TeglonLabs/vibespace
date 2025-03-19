package models

import (
	"strings"
)

// Sensor data constants
const (
	// Energy level constants
	EnergyMin float64 = 0.0
	EnergyMax float64 = 1.0
	
	// Movement level constants
	MovementMin float64 = 0.0
	MovementMax float64 = 1.0
	
	// Common mood types
	MoodCalm         string = "calm"
	MoodFocused      string = "focused"
	MoodRelaxed      string = "relaxed"
	MoodEnergetic    string = "energetic"
	MoodCreative     string = "creative"
	MoodContemplative string = "contemplative"
	MoodProductive   string = "productive"
	MoodNeutral      string = "neutral"
)

// SensorData represents environmental sensor data that contributes to a vibe
type SensorData struct {
	Temperature *float64 `json:"temperature,omitempty"` // in Celsius
	Humidity    *float64 `json:"humidity,omitempty"`    // percentage
	Light       *float64 `json:"light,omitempty"`       // in lux
	Sound       *float64 `json:"sound,omitempty"`       // in dB
	Movement    *float64 `json:"movement,omitempty"`    // relative activity level 0-1
}

// ContextLevel defines how much context is shared with viewers
type ContextLevel string

// Context level options
const (
	ContextLevelNone    ContextLevel = "none"    // Share minimal information
	ContextLevelPartial ContextLevel = "partial" // Share some details but not all
	ContextLevelFull    ContextLevel = "full"    // Share complete context
)

// SharingSettings controls how a world or vibe is shared with others
type SharingSettings struct {
	IsPublic     bool         `json:"isPublic"`                // Whether this is visible to all users
	AllowedUsers []string     `json:"allowedUsers,omitempty"`  // Specific users who can access this
	ContextLevel ContextLevel `json:"contextLevel"`            // Amount of context to share
}

// Vibe represents the emotional atmosphere of a space
type Vibe struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Energy      float64        `json:"energy"`  // from 0.0 (low) to 1.0 (high)
	Mood        string         `json:"mood"`    // e.g., "relaxed", "energetic", "contemplative"
	Colors      []string       `json:"colors"`  // hex colors
	SensorData  SensorData     `json:"sensorData,omitempty"`
	CreatorID   string         `json:"creatorId,omitempty"`    // User who created this vibe
	Sharing     SharingSettings `json:"sharing,omitempty"`     // How this vibe is shared
}

// WorldType represents the type of a world
type WorldType string

const (
	WorldTypePhysical WorldType = "PHYSICAL" // A real-world physical space
	WorldTypeVirtual  WorldType = "VIRTUAL"  // A digital/virtual space
	WorldTypeHybrid   WorldType = "HYBRID"   // A combination of physical and virtual
)

// Resource URI schemes
const (
	VibeScheme      string = "vibe://"
	WorldScheme     string = "world://"
	VibeListURI     string = "vibe://list"
	WorldListURI    string = "world://list"
	WorldVibeSubURI string = "/vibe"
)

// World represents a physical or virtual world space
type World struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        WorldType      `json:"type"`
	Location    string         `json:"location,omitempty"`     // physical location or virtual address
	CurrentVibe string         `json:"currentVibe,omitempty"`  // ID of current vibe
	Size        string         `json:"size,omitempty"`         // description of size/scale
	Features    []string       `json:"features,omitempty"`     // special characteristics
	CreatorID   string         `json:"creatorId,omitempty"`    // User who created this world
	Sharing     SharingSettings `json:"sharing,omitempty"`     // How this world is shared
	Occupancy   int            `json:"occupancy,omitempty"`    // Current number of people
}

// DataEncoding represents the encoding format of binary or ternary data
type DataEncoding string

const (
	EncodingBinary          DataEncoding = "binary"           // Raw binary data
	EncodingBase64          DataEncoding = "base64"           // Base64 encoded binary data
	EncodingBalancedTernary DataEncoding = "balanced-ternary" // Balanced ternary encoding
	EncodingHex             DataEncoding = "hex"              // Hexadecimal encoded data
	EncodingJSON            DataEncoding = "json"             // JSON string (default for CustomData)
)

// BinaryData represents a binary or ternary data payload
type BinaryData struct {
	Data     []byte       `json:"data"`               // Binary data as byte array
	Encoding DataEncoding `json:"encoding,omitempty"` // Encoding format
	Format   string       `json:"format,omitempty"`   // Format description or MIME type
}

// BalancedTernaryData represents data in balanced ternary format
// Values are represented as -1, 0, 1 (T, 0, 1 in traditional notation)
type BalancedTernaryData []int8

// NewBalancedTernaryFromString creates balanced ternary data from a string
// where 'T' or '-' represents -1, '0' represents 0, and '1' represents 1
func NewBalancedTernaryFromString(s string) BalancedTernaryData {
	result := make(BalancedTernaryData, len(s))
	for i, c := range s {
		switch c {
		case 'T', '-', 't':
			result[i] = -1
		case '0':
			result[i] = 0
		case '1':
			result[i] = 1
		default:
			// Invalid character, default to 0
			result[i] = 0
		}
	}
	return result
}

// String returns the balanced ternary data as a string
// using 'T' for -1, '0' for 0, and '1' for 1
func (bt BalancedTernaryData) String() string {
	var b strings.Builder
	for _, v := range bt {
		switch v {
		case -1:
			b.WriteRune('T')
		case 0:
			b.WriteRune('0')
		case 1:
			b.WriteRune('1')
		default:
			// Invalid value, default to 0
			b.WriteRune('0')
		}
	}
	return b.String()
}

// ToDecimal converts balanced ternary data to a decimal integer
func (bt BalancedTernaryData) ToDecimal() int64 {
	var result int64
	power := int64(1)
	
	for i := len(bt) - 1; i >= 0; i-- {
		result += int64(bt[i]) * power
		power *= 3
	}
	
	return result
}

// FromDecimal converts a decimal integer to balanced ternary
func FromDecimal(n int64) BalancedTernaryData {
	if n == 0 {
		return BalancedTernaryData{0}
	}

	var result BalancedTernaryData
	
	// Calculate the absolute value for processing
	abs := n
	if abs < 0 {
		abs = -abs
	}
	
	// This algorithm works by iteratively finding the remainder when
	// dividing by 3, then adjusting if the remainder is 2
	for abs > 0 {
		rem := abs % 3
		
		if rem == 0 {
			result = append(result, 0)
		} else if rem == 1 {
			result = append(result, 1)
		} else { // rem == 2
			result = append(result, -1)
			abs += 3 // Carry 1 to the next position
		}
		
		abs /= 3
	}
	
	// If the original number was negative, negate all digits
	if n < 0 {
		for i := range result {
			result[i] = -result[i]
		}
	}
	
	// Reverse the result (since we built it backwards)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// WorldMoment represents a moment in time for a world with its current state
type WorldMoment struct {
	WorldID     string          `json:"worldId"`              // ID of the world
	Timestamp   int64           `json:"timestamp"`            // Unix timestamp in milliseconds
	VibeID      string          `json:"vibeId,omitempty"`     // Current vibe ID
	Vibe        *Vibe           `json:"vibe,omitempty"`       // Full vibe object (optional)
	SensorData  SensorData      `json:"sensorData,omitempty"` // Current sensor data
	Occupancy   int             `json:"occupancy,omitempty"`  // Number of people/entities in the world
	Activity    float64         `json:"activity,omitempty"`   // Activity level 0.0-1.0
	CustomData  string          `json:"customData,omitempty"` // Custom JSON data for extension
	
	// Binary and balanced ternary data
	BinaryData          *BinaryData          `json:"binaryData,omitempty"`          // Binary payload
	BalancedTernaryData *BalancedTernaryData `json:"balancedTernaryData,omitempty"` // Balanced ternary data
	
	// Multiplayer additions
	CreatorID   string          `json:"creatorId,omitempty"`  // User who created this moment
	Viewers     []string        `json:"viewers,omitempty"`    // Users currently viewing this world
	Sharing     SharingSettings `json:"sharing,omitempty"`    // How this moment should be shared
}