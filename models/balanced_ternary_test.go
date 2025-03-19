package models

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestNewBalancedTernaryFromString(t *testing.T) {
	testCases := []struct {
		input    string
		expected BalancedTernaryData
	}{
		{"", BalancedTernaryData{}},
		{"0", BalancedTernaryData{0}},
		{"1", BalancedTernaryData{1}},
		{"T", BalancedTernaryData{-1}},
		{"10T", BalancedTernaryData{1, 0, -1}},
		{"1T0T1", BalancedTernaryData{1, -1, 0, -1, 1}},
		{"-0-", BalancedTernaryData{-1, 0, -1}},
		{"abc", BalancedTernaryData{0, 0, 0}}, // Invalid chars default to 0
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := NewBalancedTernaryFromString(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBalancedTernaryString(t *testing.T) {
	testCases := []struct {
		input    BalancedTernaryData
		expected string
	}{
		{BalancedTernaryData{}, ""},
		{BalancedTernaryData{0}, "0"},
		{BalancedTernaryData{1}, "1"},
		{BalancedTernaryData{-1}, "T"},
		{BalancedTernaryData{1, 0, -1}, "10T"},
		{BalancedTernaryData{1, -1, 0, -1, 1}, "1T0T1"},
		{BalancedTernaryData{2}, "0"}, // Invalid values default to 0
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.input.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBalancedTernaryToDecimal(t *testing.T) {
	testCases := []struct {
		input    BalancedTernaryData
		expected int64
	}{
		{BalancedTernaryData{}, 0},
		{BalancedTernaryData{0}, 0},
		{BalancedTernaryData{1}, 1},
		{BalancedTernaryData{-1}, -1},
		{BalancedTernaryData{1, 0}, 3},    // 1×3¹ + 0×3⁰ = 3
		{BalancedTernaryData{1, 1}, 4},    // 1×3¹ + 1×3⁰ = 4
		{BalancedTernaryData{1, -1}, 2},   // 1×3¹ + (-1)×3⁰ = 2
		{BalancedTernaryData{-1, 0}, -3},  // (-1)×3¹ + 0×3⁰ = -3
		{BalancedTernaryData{-1, 1}, -2},  // (-1)×3¹ + 1×3⁰ = -2
		{BalancedTernaryData{-1, -1}, -4}, // (-1)×3¹ + (-1)×3⁰ = -4
		{BalancedTernaryData{1, 0, 1}, 10}, // 1×3² + 0×3¹ + 1×3⁰ = 10
	}

	for _, tc := range testCases {
		t.Run(tc.input.String(), func(t *testing.T) {
			result := tc.input.ToDecimal()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFromDecimal(t *testing.T) {
	testCases := []struct {
		input    int64
		expected BalancedTernaryData
	}{
		{0, BalancedTernaryData{0}},
		{1, BalancedTernaryData{1}},
		{-1, BalancedTernaryData{-1}},
		{3, BalancedTernaryData{1, 0}},
		{4, BalancedTernaryData{1, 1}},
		{2, BalancedTernaryData{1, -1}},
		{-3, BalancedTernaryData{-1, 0}},
		{-2, BalancedTernaryData{-1, 1}},
		{-4, BalancedTernaryData{-1, -1}},
		{10, BalancedTernaryData{1, 0, 1}},
		{13, BalancedTernaryData{1, 1, 1}},
		{-13, BalancedTernaryData{-1, -1, -1}},
		{40, BalancedTernaryData{1, 1, 1, 1}},
	}

	for _, tc := range testCases {
		t.Run(FromDecimal(tc.input).String(), func(t *testing.T) {
			result := FromDecimal(tc.input)
			// Check both direct comparison and converting back to decimal
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.input, result.ToDecimal())
		})
	}
}

func TestRoundTripDecimalTernary(t *testing.T) {
	// Test round trip conversion from decimal to ternary and back
	testValues := []int64{
		0, 1, -1, 2, -2, 3, -3,
		10, -10, 42, -42, 100, -100,
		1234, -1234, 9876, -9876,
		1000000, -1000000,
	}

	for _, value := range testValues {
		ternary := FromDecimal(value)
		result := ternary.ToDecimal()
		assert.Equal(t, value, result)
	}
}