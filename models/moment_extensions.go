package models

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// AttachBinaryData adds binary data to a WorldMoment with the specified encoding
func (m *WorldMoment) AttachBinaryData(data []byte, encoding DataEncoding, format string) error {
	switch encoding {
	case EncodingBinary:
		m.BinaryData = &BinaryData{
			Data:     data,
			Encoding: encoding,
			Format:   format,
		}
	case EncodingBase64:
		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(encoded, data)
		m.BinaryData = &BinaryData{
			Data:     encoded,
			Encoding: encoding,
			Format:   format,
		}
	case EncodingHex:
		encoded := make([]byte, hex.EncodedLen(len(data)))
		hex.Encode(encoded, data)
		m.BinaryData = &BinaryData{
			Data:     encoded,
			Encoding: encoding,
			Format:   format,
		}
	default:
		return fmt.Errorf("unsupported encoding format: %s", encoding)
	}
	return nil
}

// GetBinaryData retrieves the decoded binary data from a WorldMoment
func (m *WorldMoment) GetBinaryData() ([]byte, error) {
	if m.BinaryData == nil {
		return nil, fmt.Errorf("no binary data attached")
	}

	switch m.BinaryData.Encoding {
	case EncodingBinary:
		return m.BinaryData.Data, nil
	case EncodingBase64:
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(m.BinaryData.Data)))
		n, err := base64.StdEncoding.Decode(decoded, m.BinaryData.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64: %w", err)
		}
		return decoded[:n], nil
	case EncodingHex:
		decoded := make([]byte, hex.DecodedLen(len(m.BinaryData.Data)))
		n, err := hex.Decode(decoded, m.BinaryData.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode hex: %w", err)
		}
		return decoded[:n], nil
	default:
		return nil, fmt.Errorf("unsupported encoding format: %s", m.BinaryData.Encoding)
	}
}

// AttachBalancedTernaryData adds balanced ternary data to a WorldMoment
func (m *WorldMoment) AttachBalancedTernaryData(data BalancedTernaryData) {
	ternary := make(BalancedTernaryData, len(data))
	copy(ternary, data)
	m.BalancedTernaryData = &ternary
}

// AttachBalancedTernaryFromString adds balanced ternary data from a string representation
func (m *WorldMoment) AttachBalancedTernaryFromString(ternaryStr string) {
	ternary := NewBalancedTernaryFromString(ternaryStr)
	m.BalancedTernaryData = &ternary
}

// AttachBalancedTernaryFromDecimal adds balanced ternary data from a decimal integer
func (m *WorldMoment) AttachBalancedTernaryFromDecimal(decimal int64) {
	ternary := FromDecimal(decimal)
	m.BalancedTernaryData = &ternary
}