package streaming

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTernaryStates(t *testing.T) {
	// Test ternary state arithmetic
	assert.Equal(t, TernaryNegative, TernaryState(-1))
	assert.Equal(t, TernaryNeutral, TernaryState(0))
	assert.Equal(t, TernaryPositive, TernaryState(1))
	
	// Test ternary state negation
	assert.Equal(t, TernaryPositive, -TernaryNegative)
	assert.Equal(t, TernaryNeutral, -TernaryNeutral)
	assert.Equal(t, TernaryNegative, -TernaryPositive)
}

func TestWorldMomentStreamLockFree(t *testing.T) {
	capacity := uint64(1024) // Must be power of 2
	stream := NewWorldMomentStream(capacity)
	
	// Test basic push/pop
	moment1 := &models.WorldMoment{WorldID: "test1", Timestamp: 1000}
	assert.True(t, stream.Push(moment1))
	
	retrieved, ok := stream.Pop()
	assert.True(t, ok)
	assert.Equal(t, "test1", retrieved.WorldID)
	assert.Equal(t, int64(1000), retrieved.Timestamp)
	
	// Test empty stream
	_, ok = stream.Pop()
	assert.False(t, ok)
}

func TestWorldMomentStreamConcurrency(t *testing.T) {
	capacity := uint64(1024)
	stream := NewWorldMomentStream(capacity)
	
	numProducers := 4
	numConsumers := 2
	itemsPerProducer := 1000
	
	totalProduced := int64(0)
	totalConsumed := int64(0)
	
	// Start producers
	for i := 0; i < numProducers; i++ {
		go func(producerID int) {
			for j := 0; j < itemsPerProducer; j++ {
				moment := &models.WorldMoment{
					WorldID:   "producer-test",
					Timestamp: int64(producerID*1000 + j),
				}
				
				// Retry until successful (for testing, real code might handle backpressure differently)
				for !stream.Push(moment) {
					runtime.Gosched() // Yield to other goroutines
				}
				
				atomic.AddInt64(&totalProduced, 1)
			}
		}(i)
	}
	
	// Start consumers
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	for i := 0; i < numConsumers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				if moment, ok := stream.Pop(); ok {
					if moment != nil {
						assert.Equal(t, "producer-test", moment.WorldID)
						atomic.AddInt64(&totalConsumed, 1)
					}
				} else {
					runtime.Gosched()
				}
				}
			}
		}()
	}
	
	// Wait for all items to be processed
	expected := int64(numProducers * itemsPerProducer)
	
	// Wait for production to complete
	for atomic.LoadInt64(&totalProduced) < expected {
		time.Sleep(10 * time.Millisecond)
	}
	
	// Wait for consumption to complete
	deadline := time.Now().Add(2 * time.Second)
	for atomic.LoadInt64(&totalConsumed) < expected && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	
	assert.Equal(t, expected, atomic.LoadInt64(&totalProduced))
	assert.Equal(t, expected, atomic.LoadInt64(&totalConsumed))
}

func TestVibeEnergyTransducer(t *testing.T) {
	transducer := VibeEnergyTransducer()
	
	// Create a simple accumulating reducer
	reducer := func(acc *models.Vibe, input *models.Vibe) (*models.Vibe, bool) {
		return input, false // Just pass through the transformed value
	}
	
	transformedReducer := transducer(reducer)
	
	// Test energy modulation
	testCases := []struct {
		name           string
		initialEnergy  float64
		expectedRange  [2]float64 // min, max expected range
	}{
		{"low_energy", 0.2, [2]float64{0.1, 0.3}},
		{"medium_energy", 0.5, [2]float64{0.4, 0.6}},
		{"high_energy", 0.8, [2]float64{0.7, 0.9}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vibe := &models.Vibe{
				ID:     "test-vibe",
				Energy: tc.initialEnergy,
			}
			
			// Apply transformation multiple times to see state evolution
			current := vibe
			for i := 0; i < 5; i++ {
				transformed, _ := transformedReducer(nil, current)
				require.NotNil(t, transformed)
				
				// Energy should stay within reasonable bounds
				assert.GreaterOrEqual(t, transformed.Energy, 0.0)
				assert.LessOrEqual(t, transformed.Energy, 1.0)
				
				current = transformed
			}
			
			// Final energy should be within expected range (after state evolution)
			assert.GreaterOrEqual(t, current.Energy, tc.expectedRange[0]-0.1)
			assert.LessOrEqual(t, current.Energy, tc.expectedRange[1]+0.1)
		})
	}
}

func TestHardwareOptimizedVibe(t *testing.T) {
	// Create test vibe
	originalVibe := &models.Vibe{
		ID:          "test-vibe-12345",
		Name:        "Test Vibe",
		Energy:      0.75,
		SensorData: models.SensorData{
			Temperature: &[]float64{23.5}[0],
			Humidity:    &[]float64{65.0}[0],
			Light:       &[]float64{800.0}[0],
		},
	}
	
	// Convert to optimized format
	optimized := ConvertToOptimized(originalVibe)
	
	// Verify data integrity
	assert.True(t, optimized.Validate())
	assert.Equal(t, float32(0.75), optimized.Energy)
	assert.Equal(t, float32(23.5), optimized.Temperature)
	assert.Equal(t, float32(65.0), optimized.Humidity)
	assert.Equal(t, float32(800.0), optimized.Light)
	
	// Test ID packing (first 16 bytes)
	idBytes := []byte("test-vibe-12345")
	for i, b := range idBytes {
		if i < 16 {
			assert.Equal(t, b, optimized.ID[i])
		}
	}
}

func TestQuantumInspiredProcessor(t *testing.T) {
	qubits := 4
	processor := NewQuantumInspiredProcessor(qubits)
	
	// Create entangled pairs
	processor.Entangle(0, 1)
	processor.Entangle(2, 3)
	
	// Create test vibes with different energy levels
	vibes := []*models.Vibe{
		{ID: "vibe0", Energy: 0.2}, // Low -> TernaryNegative
		{ID: "vibe1", Energy: 0.5}, // Medium -> TernaryNeutral  
		{ID: "vibe2", Energy: 0.8}, // High -> TernaryPositive
		{ID: "vibe3", Energy: 0.4}, // Medium -> TernaryNeutral
	}
	
	// Apply quantum evolution
	evolved := processor.Evolve(vibes)
	require.Equal(t, len(vibes), len(evolved))
	
	// Check that entanglement affected the states
	// vibe0 (negative) should have affected vibe1 (its entangled partner)
	// vibe2 (positive) should have affected vibe3 (its entangled partner)
	
	for i, vibe := range evolved {
		assert.NotNil(t, vibe)
		assert.GreaterOrEqual(t, vibe.Energy, 0.0)
		assert.LessOrEqual(t, vibe.Energy, 1.0)
		t.Logf("Vibe %d: %s - Energy: %.3f -> %.3f", i, vibe.ID, vibes[i].Energy, vibe.Energy)
	}
}

func TestRTStreamProcessor(t *testing.T) {
	capacity := uint64(256)
	deadline := 10 * time.Millisecond
	
	processor := NewRTStreamProcessor(capacity, deadline)
	
	// Set up a simple processing function
	processor.processor = func(moment *models.WorldMoment) *models.WorldMoment {
		// Simulate some processing time
		time.Sleep(1 * time.Millisecond)
		
		// Transform the moment
		processed := *moment
		processed.Activity = moment.Activity * 1.1 // Increase activity by 10%
		return &processed
	}
	
	// Add test moments to input buffer
	for i := 0; i < 10; i++ {
		moment := &models.WorldMoment{
			WorldID:   "test-world",
			Timestamp: int64(i * 1000),
			Activity:  float64(i) * 0.1,
		}
		assert.True(t, processor.inputBuffer.Push(moment))
	}
	
	// Process for a short time
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	go processor.ProcessRealTime(ctx)
	
	// Wait for processing
	time.Sleep(50 * time.Millisecond)
	
	// Check that some items were processed
	processed := atomic.LoadUint64(&processor.stats.processed)
	assert.Greater(t, processed, uint64(0))
	
	// Verify output
	for processed > 0 {
		if moment, ok := processor.outputBuffer.Pop(); ok {
			assert.Equal(t, "test-world", moment.WorldID)
			processed--
		} else {
			break
		}
	}
}

func TestMemoryPools(t *testing.T) {
	// Test vibe pool
	vibe1 := GetVibe()
	assert.NotNil(t, vibe1)
	
	vibe1.ID = "test"
	vibe1.Energy = 0.5
	PutVibe(vibe1)
	
	vibe2 := GetVibe()
	assert.NotNil(t, vibe2)
	// Should be cleared
	assert.Equal(t, "", vibe2.ID)
	assert.Equal(t, 0.0, vibe2.Energy)
	
	// Test moment pool
	moment1 := GetMoment()
	assert.NotNil(t, moment1)
	
	moment1.WorldID = "test"
	moment1.Viewers = []string{"user1", "user2"}
	PutMoment(moment1)
	
	moment2 := GetMoment()
	assert.NotNil(t, moment2)
	// Should be cleared but slice capacity preserved
	assert.Equal(t, "", moment2.WorldID)
	assert.Equal(t, 0, len(moment2.Viewers))
	assert.GreaterOrEqual(t, cap(moment2.Viewers), 2) // Capacity should be preserved
}

func TestCPUOptimizedBatch(t *testing.T) {
	// Create test vibes
	numVibes := 100
	vibes := make([]*models.Vibe, numVibes)
	for i := 0; i < numVibes; i++ {
		vibes[i] = &models.Vibe{
			ID:     fmt.Sprintf("vibe-%d", i),
			Energy: float64(i) / float64(numVibes),
		}
	}
	
	// Track processing
	processedCount := int64(0)
	processor := func(batch []*models.Vibe) {
		for _, vibe := range batch {
			assert.NotNil(t, vibe)
			atomic.AddInt64(&processedCount, 1)
		}
	}
	
	// Process with optimal batch size
	batchSize := 10
	CPUOptimizedBatch(vibes, batchSize, processor)
	
	// Verify all vibes were processed
	assert.Equal(t, int64(numVibes), atomic.LoadInt64(&processedCount))
}

func TestStreamTransformerComposition(t *testing.T) {
	transformer := NewStreamTransformer[*models.Vibe](10)
	defer transformer.cancel()
	
	// Chain multiple transformations
	energyBoost := TernaryTransducer(func(vibe *models.Vibe, state TernaryState) (*models.Vibe, TernaryState) {
		if vibe == nil {
			return vibe, state
		}
		boosted := *vibe
		boosted.Energy = min(1.0, vibe.Energy+0.1)
		return &boosted, TernaryPositive
	})
	
	transformer.Chain(energyBoost)
	
	// Test with sample vibe
	input := &models.Vibe{
		ID:     "test",
		Energy: 0.5,
	}
	
	output := make(chan *models.Vibe, 1)
	go transformer.Process(output)
	
	// Send input
	transformer.source <- input
	close(transformer.source)
	
	// Receive output
	result := <-output
	assert.NotNil(t, result)
	assert.Equal(t, "test", result.ID)
	assert.Greater(t, result.Energy, input.Energy)
}

// Benchmark tests for performance validation

func BenchmarkWorldMomentStreamPush(b *testing.B) {
	stream := NewWorldMomentStream(1024)
	moment := &models.WorldMoment{WorldID: "bench", Timestamp: 1000}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !stream.Push(moment) {
			// Buffer full, pop one to make space
			stream.Pop()
			stream.Push(moment)
		}
	}
}

func BenchmarkWorldMomentStreamPop(b *testing.B) {
	stream := NewWorldMomentStream(1024)
	moment := &models.WorldMoment{WorldID: "bench", Timestamp: 1000}
	
	// Fill buffer
	for i := 0; i < 512; i++ {
		stream.Push(moment)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := stream.Pop(); !ok {
			// Buffer empty, add one
			stream.Push(moment)
			stream.Pop()
		}
	}
}

func BenchmarkHardwareOptimizedVibeConversion(b *testing.B) {
	vibe := &models.Vibe{
		ID:     "benchmark-vibe",
		Energy: 0.75,
		SensorData: models.SensorData{
			Temperature: &[]float64{23.5}[0],
			Humidity:    &[]float64{65.0}[0],
			Light:       &[]float64{800.0}[0],
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimized := ConvertToOptimized(vibe)
		_ = optimized.Validate()
	}
}

func BenchmarkQuantumInspiredProcessor(b *testing.B) {
	processor := NewQuantumInspiredProcessor(8)
	processor.Entangle(0, 1)
	processor.Entangle(2, 3)
	processor.Entangle(4, 5)
	processor.Entangle(6, 7)
	
	vibes := make([]*models.Vibe, 8)
	for i := 0; i < 8; i++ {
		vibes[i] = &models.Vibe{
			ID:     fmt.Sprintf("bench-vibe-%d", i),
			Energy: float64(i) / 8.0,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processor.Evolve(vibes)
	}
}

func BenchmarkMemoryPools(b *testing.B) {
	b.Run("VibePool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			vibe := GetVibe()
			vibe.ID = "bench"
			vibe.Energy = 0.5
			PutVibe(vibe)
		}
	})
	
	b.Run("MomentPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			moment := GetMoment()
			moment.WorldID = "bench"
			moment.Viewers = append(moment.Viewers, "user1")
			PutMoment(moment)
		}
	})
}
