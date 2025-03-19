package streaming

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	
	"github.com/bmorphism/vibespace-mcp-go/models"
)

// EncodeBinaryData converts binary data to the specified encoding format
func EncodeBinaryData(data []byte, encoding models.DataEncoding) ([]byte, error) {
	switch encoding {
	case models.EncodingBinary:
		return data, nil
	case models.EncodingBase64:
		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(encoded, data)
		return encoded, nil
	case models.EncodingHex:
		encoded := make([]byte, hex.EncodedLen(len(data)))
		hex.Encode(encoded, data)
		return encoded, nil
	default:
		return nil, fmt.Errorf("unsupported encoding format: %s", encoding)
	}
}

// DecodeBinaryData converts encoded data back to binary format
func DecodeBinaryData(data []byte, encoding models.DataEncoding) ([]byte, error) {
	switch encoding {
	case models.EncodingBinary:
		return data, nil
	case models.EncodingBase64:
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		n, err := base64.StdEncoding.Decode(decoded, data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64: %w", err)
		}
		return decoded[:n], nil
	case models.EncodingHex:
		decoded := make([]byte, hex.DecodedLen(len(data)))
		n, err := hex.Decode(decoded, data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode hex: %w", err)
		}
		return decoded[:n], nil
	default:
		return nil, fmt.Errorf("unsupported encoding format: %s", encoding)
	}
}

// TernaryToBytes converts balanced ternary data to a compact byte representation
// Each byte will contain up to 4 ternary digits (trits), using 4*log2(3) ≈ 6.34 bits
func TernaryToBytes(ternary models.BalancedTernaryData) []byte {
	if len(ternary) == 0 {
		return []byte{}
	}

	// Calculate how many bytes we need
	// Each byte can store 4 trits (4 trits = 4*log2(3) ≈ 6.34 bits < 8 bits)
	numBytes := (len(ternary) + 3) / 4 // Ceiling division by 4
	
	result := make([]byte, numBytes)
	for i := 0; i < len(ternary); i++ {
		byteIdx := i / 4
		tritPos := i % 4
		
		// Map trits to bit patterns:
		// -1 (T) -> 00
		//  0     -> 01
		//  1     -> 10
		var bits byte
		switch ternary[i] {
		case -1:
			bits = 0b00
		case 0:
			bits = 0b01
		case 1:
			bits = 0b10
		default:
			bits = 0b01 // Default to 0 for invalid values
		}
		
		// Place the 2-bit pattern at the right position
		// First trit is at the lowest bits
		result[byteIdx] |= (bits << (tritPos * 2))
	}
	
	return result
}

// BytesToTernary converts a compact byte representation back to balanced ternary
func BytesToTernary(bytes []byte, numTrits int) models.BalancedTernaryData {
	if len(bytes) == 0 || numTrits == 0 {
		return models.BalancedTernaryData{}
	}

	result := make(models.BalancedTernaryData, numTrits)
	
	for i := 0; i < numTrits && i/4 < len(bytes); i++ {
		byteIdx := i / 4
		tritPos := i % 4
		
		// Extract the trit value (2 bits)
		bits := (bytes[byteIdx] >> (tritPos * 2)) & 0b11
		
		// Map bit patterns back to trits:
		// 00 -> -1 (T)
		// 01 ->  0
		// 10 ->  1
		// 11 ->  0 (invalid pattern, default to 0)
		switch bits {
		case 0b00:
			result[i] = -1
		case 0b01:
			result[i] = 0
		case 0b10:
			result[i] = 1
		default:
			result[i] = 0 // Default to 0 for invalid patterns
		}
	}
	
	return result
}

// CreateBinaryData creates a BinaryData struct with the specified encoding
func CreateBinaryData(data []byte, encoding models.DataEncoding, format string) (*models.BinaryData, error) {
	encodedData, err := EncodeBinaryData(data, encoding)
	if err != nil {
		return nil, err
	}
	
	return &models.BinaryData{
		Data:     encodedData,
		Encoding: encoding,
		Format:   format,
	}, nil
}