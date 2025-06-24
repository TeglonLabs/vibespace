package streaming

import (
	"math"
	"testing"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
)

func TestComonadicVibeContext(t *testing.T) {
	// Create test vibes
	center := &models.Vibe{
		ID:     "center",
		Energy: 0.5,
	}
	neighbors := []*models.Vibe{
		{ID: "neighbor1", Energy: 0.3},
		{ID: "neighbor2", Energy: 0.7},
	}
	temporal := []*models.Vibe{
		{ID: "past", Energy: 0.4},
	}

	ctx := &ComonadicVibeContext{
		center:    center,
		neighbors: neighbors,
		temporal:  temporal,
		gradient:  TernaryPositive,
		coherence: 0.8,
	}

	t.Run("Extract", func(t *testing.T) {
		extracted := ctx.Extract()
		assert.Equal(t, center, extracted)
		assert.Equal(t, "center", extracted.ID)
	})

	t.Run("Duplicate", func(t *testing.T) {
		duplicated := ctx.Duplicate()
		
		// Should preserve center and structure
		assert.Equal(t, center, duplicated.center)
		assert.Equal(t, ctx.gradient, duplicated.gradient)
		
		// Coherence should decay slightly
		assert.Less(t, duplicated.coherence, ctx.coherence)
		assert.Equal(t, ctx.coherence*0.9, duplicated.coherence)
		
		// Neighbors should be elevated to contexts
		assert.Len(t, duplicated.neighbors, len(neighbors))
	})

	t.Run("Extend", func(t *testing.T) {
		// Create a transformation that boosts energy
		boostTransform := func(c *ComonadicVibeContext) *models.Vibe {
			boosted := *c.center
			boosted.Energy = min(1.0, boosted.Energy+0.2)
			return &boosted
		}

		extended := ctx.Extend(boostTransform)
		
		// Center should be transformed
		assert.Greater(t, extended.center.Energy, center.Energy)
		
		// Should have same number of neighbors
		assert.Len(t, extended.neighbors, len(neighbors))
		
		// Context structure should be preserved
		assert.Equal(t, ctx.gradient, extended.gradient)
		assert.Equal(t, ctx.coherence, extended.coherence)
	})
}

func TestTernaryLogicGate(t *testing.T) {
	t.Run("Consensus Gate", func(t *testing.T) {
		gate := NewTernaryLogicGate("consensus")
		
		// Test agreement cases
		assert.Equal(t, TernaryNegative, gate.Apply(TernaryNegative, TernaryNegative))
		assert.Equal(t, TernaryNeutral, gate.Apply(TernaryNeutral, TernaryNeutral))
		assert.Equal(t, TernaryPositive, gate.Apply(TernaryPositive, TernaryPositive))
		
		// Test disagreement cases (should default to neutral)
		assert.Equal(t, TernaryNeutral, gate.Apply(TernaryNegative, TernaryPositive))
		assert.Equal(t, TernaryNeutral, gate.Apply(TernaryPositive, TernaryNegative))
	})

	t.Run("Amplify Gate", func(t *testing.T) {
		gate := NewTernaryLogicGate("amplify")
		
		// Test strengthening agreement
		assert.Equal(t, TernaryNegative, gate.Apply(TernaryNegative, TernaryNegative))
		assert.Equal(t, TernaryPositive, gate.Apply(TernaryPositive, TernaryPositive))
		
		// Test neutral preservation
		assert.Equal(t, TernaryNegative, gate.Apply(TernaryNeutral, TernaryNegative))
		assert.Equal(t, TernaryPositive, gate.Apply(TernaryNeutral, TernaryPositive))
		
		// Test disagreement weakening
		assert.Equal(t, TernaryNeutral, gate.Apply(TernaryNegative, TernaryPositive))
		assert.Equal(t, TernaryNeutral, gate.Apply(TernaryPositive, TernaryNegative))
	})

	t.Run("Inhibit Gate", func(t *testing.T) {
		gate := NewTernaryLogicGate("inhibit")
		
		// Test inhibition patterns
		assert.Equal(t, TernaryNegative, gate.Apply(TernaryNegative, TernaryPositive)) // Disinhibition
		assert.Equal(t, TernaryNeutral, gate.Apply(TernaryPositive, TernaryNegative))  // Inhibition
		assert.Equal(t, TernaryNeutral, gate.Apply(TernaryPositive, TernaryPositive))  // Self-inhibition
		
		// Test double negative
		assert.Equal(t, TernaryNegative, gate.Apply(TernaryNegative, TernaryNegative))
	})

	t.Run("Default Gate", func(t *testing.T) {
		gate := NewTernaryLogicGate("unknown")
		
		// Should default to first input wins
		assert.Equal(t, TernaryNegative, gate.Apply(TernaryNegative, TernaryPositive))
		assert.Equal(t, TernaryPositive, gate.Apply(TernaryPositive, TernaryNegative))
		assert.Equal(t, TernaryNeutral, gate.Apply(TernaryNeutral, TernaryPositive))
	})
}

func TestVibeContextualTransformer(t *testing.T) {
	transformer := NewVibeContextualTransformer(5)
	
	t.Run("Basic Initialization", func(t *testing.T) {
		assert.NotNil(t, transformer.logicGates["consensus"])
		assert.NotNil(t, transformer.logicGates["amplify"])
		assert.NotNil(t, transformer.logicGates["inhibit"])
		assert.Equal(t, 5, transformer.maxHistory)
		assert.Empty(t, transformer.history)
	})

	t.Run("Transform High Coherence", func(t *testing.T) {
		center := &models.Vibe{ID: "center", Energy: 0.5}
		neighbors := []*models.Vibe{
			{ID: "neighbor1", Energy: 0.52}, // Very similar energy = high coherence
			{ID: "neighbor2", Energy: 0.48},
		}

		result := transformer.TransformWithContext(center, neighbors)
		
		assert.NotNil(t, result)
		assert.Equal(t, "center", result.ID)
		// Energy should be influenced by consensus operation
		assert.NotEqual(t, center.Energy, result.Energy)
	})

	t.Run("Transform with History Gradient", func(t *testing.T) {
		// Build up some history first
		for i := 0; i < 3; i++ {
			vibe := &models.Vibe{
				ID:     "historical",
				Energy: 0.3 + float64(i)*0.1, // Gradually increasing energy
			}
			transformer.updateHistory(vibe)
		}

		center := &models.Vibe{ID: "center", Energy: 0.7} // Higher than history
		neighbors := []*models.Vibe{
			{ID: "neighbor", Energy: 0.9}, // Low coherence with center
		}

		result := transformer.TransformWithContext(center, neighbors)
		
		assert.NotNil(t, result)
		// Should detect positive gradient and apply amplification
		assert.NotEqual(t, center.Energy, result.Energy)
	})

	t.Run("Transform Low Coherence", func(t *testing.T) {
		center := &models.Vibe{ID: "center", Energy: 0.5}
		neighbors := []*models.Vibe{
			{ID: "neighbor1", Energy: 0.1}, // Very different energy = low coherence
			{ID: "neighbor2", Energy: 0.9},
		}

		result := transformer.TransformWithContext(center, neighbors)
		
		assert.NotNil(t, result)
		// Should apply inhibition to dampen chaos
		assert.NotEqual(t, center.Energy, result.Energy)
	})

	t.Run("History Management", func(t *testing.T) {
		transformer := NewVibeContextualTransformer(3) // Small history for testing
		
		// Add more vibes than max history
		for i := 0; i < 5; i++ {
			vibe := &models.Vibe{
				ID:     "test",
				Energy: float64(i) * 0.1,
			}
			transformer.updateHistory(vibe)
		}
		
		// Should only keep the last 3
		assert.Len(t, transformer.history, 3)
		assert.Equal(t, 0.2, transformer.history[0].Energy) // Should be the 3rd vibe (index 2)
		assert.Equal(t, 0.4, transformer.history[2].Energy) // Should be the 5th vibe (index 4)
	})
}

func TestEnergyTernaryConversion(t *testing.T) {
	testCases := []struct {
		energy   float64
		expected TernaryState
	}{
		{0.0, TernaryNegative},
		{0.2, TernaryNegative},
		{0.32, TernaryNegative},
		{0.33, TernaryNeutral},
		{0.5, TernaryNeutral},
		{0.66, TernaryNeutral},
		{0.67, TernaryPositive},
		{0.8, TernaryPositive},
		{1.0, TernaryPositive},
	}

	for _, tc := range testCases {
		t.Run("Energy to Ternary", func(t *testing.T) {
			result := energyToTernary(tc.energy)
			assert.Equal(t, tc.expected, result, "Energy %f should map to %v", tc.energy, tc.expected)
		})
	}
}

func TestTernaryEnergyConversion(t *testing.T) {
	currentEnergy := 0.5
	
	testCases := []struct {
		state    TernaryState
		expected float64
	}{
		{TernaryNegative, 0.41}, // 0.5 * 0.7 + 0.2 * 0.3
		{TernaryNeutral, 0.5},   // Should stay the same
		{TernaryPositive, 0.59}, // 0.5 * 0.7 + 0.8 * 0.3
	}

	for _, tc := range testCases {
		t.Run("Ternary to Energy", func(t *testing.T) {
			result := ternaryToEnergy(tc.state, currentEnergy)
			assert.InDelta(t, tc.expected, result, 0.01, "State %v with current %f should smooth to %f", tc.state, currentEnergy, tc.expected)
		})
	}
}

func TestCoherenceCalculation(t *testing.T) {
	transformer := NewVibeContextualTransformer(5)
	
	t.Run("Perfect Coherence in Isolation", func(t *testing.T) {
		center := &models.Vibe{Energy: 0.5}
		neighbors := []*models.Vibe{}
		
		coherence := transformer.calculateCoherence(center, neighbors)
		assert.Equal(t, 1.0, coherence)
	})

	t.Run("High Coherence with Similar Neighbors", func(t *testing.T) {
		center := &models.Vibe{Energy: 0.5}
		neighbors := []*models.Vibe{
			{Energy: 0.51},
			{Energy: 0.49},
		}
		
		coherence := transformer.calculateCoherence(center, neighbors)
		assert.Greater(t, coherence, 0.9) // Should be very high
	})

	t.Run("Low Coherence with Dissimilar Neighbors", func(t *testing.T) {
		center := &models.Vibe{Energy: 0.2} // Center different from neighbor average
		neighbors := []*models.Vibe{
			{Energy: 0.8}, // Average will be 0.8, big difference from center 0.2
		}
		
		coherence := transformer.calculateCoherence(center, neighbors)
		// With energy difference of 0.6, exp(-0.6) ≈ 0.55
		assert.Less(t, coherence, 0.6) // Should be low
	})

	t.Run("Handles Nil Neighbors", func(t *testing.T) {
		center := &models.Vibe{Energy: 0.5}
		neighbors := []*models.Vibe{
			{Energy: 0.6},
			nil, // Should be handled gracefully
			{Energy: 0.4},
		}
		
		coherence := transformer.calculateCoherence(center, neighbors)
		assert.False(t, math.IsNaN(coherence), "Coherence should not be NaN")
		assert.GreaterOrEqual(t, coherence, 0.0)
		assert.LessOrEqual(t, coherence, 1.0)
	})
}

func TestGradientCalculation(t *testing.T) {
	transformer := NewVibeContextualTransformer(5)
	
	t.Run("Neutral with No History", func(t *testing.T) {
		vibe := &models.Vibe{Energy: 0.5}
		gradient := transformer.calculateGradient(vibe)
		assert.Equal(t, TernaryNeutral, gradient)
	})

	t.Run("Neutral with Insufficient History", func(t *testing.T) {
		transformer.updateHistory(&models.Vibe{Energy: 0.4})
		vibe := &models.Vibe{Energy: 0.5}
		gradient := transformer.calculateGradient(vibe)
		assert.Equal(t, TernaryNeutral, gradient)
	})

	t.Run("Positive Gradient", func(t *testing.T) {
		transformer.updateHistory(&models.Vibe{Energy: 0.3})
		transformer.updateHistory(&models.Vibe{Energy: 0.4})
		vibe := &models.Vibe{Energy: 0.6} // Significantly higher
		gradient := transformer.calculateGradient(vibe)
		assert.Equal(t, TernaryPositive, gradient)
	})

	t.Run("Negative Gradient", func(t *testing.T) {
		transformer.updateHistory(&models.Vibe{Energy: 0.7})
		transformer.updateHistory(&models.Vibe{Energy: 0.6})
		vibe := &models.Vibe{Energy: 0.4} // Significantly lower
		gradient := transformer.calculateGradient(vibe)
		assert.Equal(t, TernaryNegative, gradient)
	})

	t.Run("Stable Gradient", func(t *testing.T) {
		transformer.updateHistory(&models.Vibe{Energy: 0.5})
		transformer.updateHistory(&models.Vibe{Energy: 0.52})
		vibe := &models.Vibe{Energy: 0.51} // Small change
		gradient := transformer.calculateGradient(vibe)
		assert.Equal(t, TernaryNeutral, gradient)
	})
}

// Benchmark tests for performance validation
func BenchmarkComonadicTransform(b *testing.B) {
	transformer := NewVibeContextualTransformer(10)
	center := &models.Vibe{ID: "center", Energy: 0.5}
	neighbors := []*models.Vibe{
		{ID: "n1", Energy: 0.4},
		{ID: "n2", Energy: 0.6},
		{ID: "n3", Energy: 0.55},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = transformer.TransformWithContext(center, neighbors)
	}
}

func BenchmarkTernaryLogicGate(b *testing.B) {
	gate := NewTernaryLogicGate("consensus")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gate.Apply(TernaryPositive, TernaryNegative)
	}
}

func BenchmarkComonadicExtension(b *testing.B) {
	ctx := &ComonadicVibeContext{
		center: &models.Vibe{ID: "center", Energy: 0.5},
		neighbors: []*models.Vibe{
			{ID: "n1", Energy: 0.4},
			{ID: "n2", Energy: 0.6},
		},
		gradient:  TernaryNeutral,
		coherence: 0.8,
	}

	transform := func(c *ComonadicVibeContext) *models.Vibe {
		boosted := *c.center
		boosted.Energy = min(1.0, boosted.Energy+0.1)
		return &boosted
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Extend(transform)
	}
}
