package models

import (
	"bytes"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestAttachBinaryData(t *testing.T) {
	moment := &WorldMoment{
		WorldID:   "test-world",
		Timestamp: 1234567890,
	}
	
	testData := []byte{0x01, 0x02, 0x03}
	format := "application/octet-stream"

	// Test binary encoding
	err := moment.AttachBinaryData(testData, EncodingBinary, format)
	assert.NoError(t, err)
	assert.NotNil(t, moment.BinaryData)
	assert.Equal(t, EncodingBinary, moment.BinaryData.Encoding)
	assert.Equal(t, format, moment.BinaryData.Format)
	assert.Equal(t, testData, moment.BinaryData.Data)
	
	// Test base64 encoding
	err = moment.AttachBinaryData(testData, EncodingBase64, format)
	assert.NoError(t, err)
	assert.NotNil(t, moment.BinaryData)
	assert.Equal(t, EncodingBase64, moment.BinaryData.Encoding)
	
	// Test hex encoding
	err = moment.AttachBinaryData(testData, EncodingHex, format)
	assert.NoError(t, err)
	assert.NotNil(t, moment.BinaryData)
	assert.Equal(t, EncodingHex, moment.BinaryData.Encoding)
	
	// Test unsupported encoding
	err = moment.AttachBinaryData(testData, "unsupported", format)
	assert.Error(t, err)
}

func TestGetBinaryData(t *testing.T) {
	moment := &WorldMoment{
		WorldID:   "test-world",
		Timestamp: 1234567890,
	}
	
	testData := []byte{0x01, 0x02, 0x03}
	format := "application/octet-stream"

	// Test with no binary data
	data, err := moment.GetBinaryData()
	assert.Error(t, err)
	assert.Nil(t, data)
	
	// Test with binary encoding
	moment.AttachBinaryData(testData, EncodingBinary, format)
	data, err = moment.GetBinaryData()
	assert.NoError(t, err)
	assert.Equal(t, testData, data)
	
	// Test with base64 encoding
	moment.AttachBinaryData(testData, EncodingBase64, format)
	data, err = moment.GetBinaryData()
	assert.NoError(t, err)
	assert.Equal(t, testData, data)
	
	// Test with hex encoding
	moment.AttachBinaryData(testData, EncodingHex, format)
	data, err = moment.GetBinaryData()
	assert.NoError(t, err)
	assert.Equal(t, testData, data)
	
	// Test with invalid base64 data
	moment.BinaryData = &BinaryData{
		Data:     []byte("invalid!!"),
		Encoding: EncodingBase64,
		Format:   format,
	}
	data, err = moment.GetBinaryData()
	assert.Error(t, err)
	
	// Test with invalid hex data
	moment.BinaryData = &BinaryData{
		Data:     []byte("invalid!!"),
		Encoding: EncodingHex,
		Format:   format,
	}
	data, err = moment.GetBinaryData()
	assert.Error(t, err)
	
	// Test with unsupported encoding
	moment.BinaryData = &BinaryData{
		Data:     testData,
		Encoding: "unsupported",
		Format:   format,
	}
	data, err = moment.GetBinaryData()
	assert.Error(t, err)
}

func TestAttachBalancedTernaryData(t *testing.T) {
	moment := &WorldMoment{
		WorldID:   "test-world",
		Timestamp: 1234567890,
	}
	
	// Test attaching ternary data
	ternary := BalancedTernaryData{1, 0, -1}
	moment.AttachBalancedTernaryData(ternary)
	assert.NotNil(t, moment.BalancedTernaryData)
	assert.Equal(t, ternary, *moment.BalancedTernaryData)
	
	// Verify that a copy was made (not reference)
	ternary[0] = -1
	assert.NotEqual(t, ternary, *moment.BalancedTernaryData)
}

func TestAttachBalancedTernaryFromString(t *testing.T) {
	moment := &WorldMoment{
		WorldID:   "test-world",
		Timestamp: 1234567890,
	}
	
	// Test attaching ternary data from string
	moment.AttachBalancedTernaryFromString("10T")
	assert.NotNil(t, moment.BalancedTernaryData)
	
	expected := BalancedTernaryData{1, 0, -1}
	assert.Equal(t, expected, *moment.BalancedTernaryData)
}

func TestAttachBalancedTernaryFromDecimal(t *testing.T) {
	moment := &WorldMoment{
		WorldID:   "test-world",
		Timestamp: 1234567890,
	}
	
	// Test attaching ternary data from decimal
	moment.AttachBalancedTernaryFromDecimal(4) // Decimal 4 is 1T1 in balanced ternary
	assert.NotNil(t, moment.BalancedTernaryData)
	
	// Verify by converting back to decimal
	decimal := (*moment.BalancedTernaryData).ToDecimal()
	assert.Equal(t, int64(4), decimal)
}

func TestRoundTripBinaryData(t *testing.T) {
	moment := &WorldMoment{
		WorldID:   "test-world",
		Timestamp: 1234567890,
	}
	
	testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	format := "application/octet-stream"
	
	// Test round trip with all encodings
	for _, encoding := range []DataEncoding{EncodingBinary, EncodingBase64, EncodingHex} {
		t.Run(string(encoding), func(t *testing.T) {
			err := moment.AttachBinaryData(testData, encoding, format)
			assert.NoError(t, err)
			
			result, err := moment.GetBinaryData()
			assert.NoError(t, err)
			assert.True(t, bytes.Equal(testData, result))
		})
	}
}