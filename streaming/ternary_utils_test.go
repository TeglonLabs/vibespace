package streaming

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"testing"
	
	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

func TestEncodeBinaryData(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		encoding models.DataEncoding
		expected []byte
		hasError bool
	}{
		{
			name:     "Binary encoding",
			data:     []byte{0x01, 0x02, 0x03},
			encoding: models.EncodingBinary,
			expected: []byte{0x01, 0x02, 0x03},
			hasError: false,
		},
		{
			name:     "Base64 encoding",
			data:     []byte{0x01, 0x02, 0x03},
			encoding: models.EncodingBase64,
			expected: []byte(base64.StdEncoding.EncodeToString([]byte{0x01, 0x02, 0x03})),
			hasError: false,
		},
		{
			name:     "Hex encoding",
			data:     []byte{0x01, 0x02, 0x03},
			encoding: models.EncodingHex,
			expected: []byte(hex.EncodeToString([]byte{0x01, 0x02, 0x03})),
			hasError: false,
		},
		{
			name:     "Unsupported encoding",
			data:     []byte{0x01, 0x02, 0x03},
			encoding: "unsupported",
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := EncodeBinaryData(tc.data, tc.encoding)
			
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestDecodeBinaryData(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		encoding models.DataEncoding
		expected []byte
		hasError bool
	}{
		{
			name:     "Binary decoding",
			data:     []byte{0x01, 0x02, 0x03},
			encoding: models.EncodingBinary,
			expected: []byte{0x01, 0x02, 0x03},
			hasError: false,
		},
		{
			name:     "Base64 decoding",
			data:     []byte(base64.StdEncoding.EncodeToString([]byte{0x01, 0x02, 0x03})),
			encoding: models.EncodingBase64,
			expected: []byte{0x01, 0x02, 0x03},
			hasError: false,
		},
		{
			name:     "Hex decoding",
			data:     []byte(hex.EncodeToString([]byte{0x01, 0x02, 0x03})),
			encoding: models.EncodingHex,
			expected: []byte{0x01, 0x02, 0x03},
			hasError: false,
		},
		{
			name:     "Invalid base64",
			data:     []byte("invalid!!"),
			encoding: models.EncodingBase64,
			expected: nil,
			hasError: true,
		},
		{
			name:     "Invalid hex",
			data:     []byte("invalid!!"),
			encoding: models.EncodingHex,
			expected: nil,
			hasError: true,
		},
		{
			name:     "Unsupported encoding",
			data:     []byte{0x01, 0x02, 0x03},
			encoding: "unsupported",
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := DecodeBinaryData(tc.data, tc.encoding)
			
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestTernaryToBytes(t *testing.T) {
	tests := []struct {
		name     string
		ternary  models.BalancedTernaryData
		expected []byte
	}{
		{
			name:     "Empty ternary",
			ternary:  models.BalancedTernaryData{},
			expected: []byte{},
		},
		{
			name:     "Single trit",
			ternary:  models.BalancedTernaryData{1},
			expected: []byte{0x02}, // 1 -> 10 in bits = 2 decimal
		},
		{
			name:     "Multiple trits in one byte",
			ternary:  models.BalancedTernaryData{1, 0, -1},
			expected: []byte{0x06}, // [1,0,-1] -> [10,01,00] -> 00000110
		},
		{
			name:     "Trits spanning multiple bytes",
			ternary:  models.BalancedTernaryData{1, 0, -1, 1, 0, -1},
			expected: []byte{0x86, 0x01}, // First 4 trits in first byte, next 2 in second byte
		},
		{
			name:     "All possible trit values",
			ternary:  models.BalancedTernaryData{-1, 0, 1},
			expected: []byte{0x24}, // [-1,0,1] -> [00,01,10] -> 00100100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TernaryToBytes(tt.ternary)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytesToTernary(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		numTrits int
		expected models.BalancedTernaryData
	}{
		{
			name:     "Empty bytes",
			bytes:    []byte{},
			numTrits: 0,
			expected: models.BalancedTernaryData{},
		},
		{
			name:     "Single trit",
			bytes:    []byte{0x02}, // 00000010 -> 10 -> 1
			numTrits: 1,
			expected: models.BalancedTernaryData{1},
		},
		{
			name:     "Multiple trits in one byte",
			bytes:    []byte{0x06}, // 00000110 -> [10,01,00] -> [1,0,-1]
			numTrits: 3,
			expected: models.BalancedTernaryData{1, 0, -1},
		},
		{
			name:     "Trits spanning multiple bytes",
			bytes:    []byte{0x86, 0x01}, // First byte: [1,0,-1,1], Second byte: [0,-1]
			numTrits: 6,
			expected: models.BalancedTernaryData{1, 0, -1, 1, 0, -1},
		},
		{
			name:     "All possible trit values",
			bytes:    []byte{0x24}, // 00100100 -> [00,01,10] -> [-1,0,1]
			numTrits: 3,
			expected: models.BalancedTernaryData{-1, 0, 1},
		},
		{
			name:     "Partial read (more trits than available)",
			bytes:    []byte{0x06}, // 00000110 -> [10,01,00] -> [1,0,-1]
			numTrits: 5,
			expected: models.BalancedTernaryData{1, 0, -1, -1, 0}, // Last 2 trits based on remaining bits
		},
		{
			name:     "Invalid bit pattern",
			bytes:    []byte{0xFF}, // All bits set - contains invalid patterns
			numTrits: 4,
			expected: models.BalancedTernaryData{0, 0, 0, 0}, // All 11 patterns map to 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BytesToTernary(tt.bytes, tt.numTrits)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoundTripTernary(t *testing.T) {
	// Test round trip conversion from ternary to bytes and back
	testCases := []models.BalancedTernaryData{
		{},
		{0},
		{1},
		{-1},
		{1, 0, -1},
		{-1, -1, -1, -1, -1},
		{1, 1, 1, 1, 1},
		{0, 0, 0, 0, 0},
		{1, 0, -1, 1, 0, -1, 1, 0, -1, 1, 0},
	}

	for i, tc := range testCases {
		t.Run(string(byte('A'+i)), func(t *testing.T) {
			bytes := TernaryToBytes(tc)
			result := BytesToTernary(bytes, len(tc))
			assert.Equal(t, tc, result)
		})
	}
}

func TestCreateBinaryData(t *testing.T) {
	testData := []byte{0x01, 0x02, 0x03}
	format := "application/octet-stream"

	testCases := []struct {
		name     string
		data     []byte
		encoding models.DataEncoding
		format   string
		hasError bool
	}{
		{
			name:     "Binary encoding",
			data:     testData,
			encoding: models.EncodingBinary,
			format:   format,
			hasError: false,
		},
		{
			name:     "Base64 encoding",
			data:     testData,
			encoding: models.EncodingBase64,
			format:   format,
			hasError: false,
		},
		{
			name:     "Hex encoding",
			data:     testData,
			encoding: models.EncodingHex,
			format:   format,
			hasError: false,
		},
		{
			name:     "Unsupported encoding",
			data:     testData,
			encoding: "unsupported",
			format:   format,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateBinaryData(tc.data, tc.encoding, tc.format)
			
			if tc.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.encoding, result.Encoding)
				assert.Equal(t, tc.format, result.Format)
				
				// Check the encoded data
				expected, _ := EncodeBinaryData(tc.data, tc.encoding)
				assert.True(t, bytes.Equal(expected, result.Data))
			}
		})
	}
}