package streaming

import (
	"context"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/bmorphism/vibespace-mcp-go/models"
)

// Transducer represents a composable stateful transformation
type Transducer[A, B any] func(reducer Reducer[B]) Reducer[A]

// Reducer represents a reducing function with early termination
type Reducer[T any] func(acc T, input T) (result T, terminate bool)

// XForm represents a transducible collection with zero-allocation iteration
type XForm[T any] interface {
	Transduce(ctx context.Context, t Transducer[T, T], reducer Reducer[T], init T) T
}

// TernaryState represents balanced ternary logic states for quantum-inspired computation
type TernaryState int8

const (
	TernaryNegative TernaryState = -1 // T or -
	TernaryNeutral  TernaryState = 0  // 0
	TernaryPositive TernaryState = 1  // 1 or +
)

// TernaryVector uses memory-aligned packed bits for TinyGo efficiency
type TernaryVector struct {
	data   []uint64 // Pack 21 ternary states per uint64 (3^21 < 2^64)
	length int32
	_      [4]byte // Padding for alignment
}

// WorldMomentStream provides lock-free streaming with hardware memory ordering
type WorldMomentStream struct {
	buffer    []unsafe.Pointer // Circular buffer of *models.WorldMoment
	head      uint64           // Atomic head pointer
	tail      uint64           // Atomic tail pointer
	capacity  uint64           // Power of 2 for fast modulo
	mask      uint64           // capacity - 1
	_padding1 [56]byte         // Cache line padding
	closed    uint32           // Atomic closed flag
	_padding2 [60]byte         // Cache line padding
}

// NewWorldMomentStream creates a lock-free circular buffer optimized for single-producer/single-consumer
func NewWorldMomentStream(capacity uint64) *WorldMomentStream {
	// Ensure capacity is power of 2 for efficient modulo operations
	if capacity == 0 || (capacity&(capacity-1)) != 0 {
		panic("capacity must be power of 2")
	}
	
	return &WorldMomentStream{
		buffer:   make([]unsafe.Pointer, capacity),
		capacity: capacity,
		mask:     capacity - 1,
	}
}

// Push adds a moment to the stream (non-blocking for real-time systems)
func (s *WorldMomentStream) Push(moment *models.WorldMoment) bool {
	if atomic.LoadUint32(&s.closed) != 0 {
		return false
	}
	
	head := atomic.LoadUint64(&s.head)
	nextHead := (head + 1) & s.mask
	
	// Check if buffer is full (leave one slot empty to distinguish full from empty)
	if nextHead == (atomic.LoadUint64(&s.tail) & s.mask) {
		return false // Buffer full
	}
	
	// Store the moment
	atomic.StorePointer(&s.buffer[head], unsafe.Pointer(moment))
	
	// Advance head with memory barrier
	atomic.StoreUint64(&s.head, nextHead)
	return true
}

// Pop removes a moment from the stream (non-blocking)
func (s *WorldMomentStream) Pop() (*models.WorldMoment, bool) {
	tail := atomic.LoadUint64(&s.tail)
	head := atomic.LoadUint64(&s.head)
	
	// Check if buffer is empty
	if tail == head {
		return nil, false
	}
	
	// Load the moment
	ptr := atomic.LoadPointer(&s.buffer[tail])
	moment := (*models.WorldMoment)(ptr)
	
	// Clear the slot
	atomic.StorePointer(&s.buffer[tail], nil)
	
	// Advance tail
	atomic.StoreUint64(&s.tail, (tail+1)&s.mask)
	return moment, true
}

// TernaryTransducer creates composable ternary logic transformations
func TernaryTransducer[T any](logic func(T, TernaryState) (T, TernaryState)) Transducer[T, T] {
	return func(reducer Reducer[T]) Reducer[T] {
		state := TernaryNeutral
		return func(acc T, input T) (T, bool) {
			result, newState := logic(input, state)
			state = newState
			return reducer(acc, result)
		}
	}
}

// VibeEnergyTransducer applies quantum-inspired energy level transformations
func VibeEnergyTransducer() Transducer[*models.Vibe, *models.Vibe] {
	return TernaryTransducer(func(vibe *models.Vibe, state TernaryState) (*models.Vibe, TernaryState) {
		if vibe == nil {
			return vibe, state
		}
		
		// Create a copy to avoid mutation
		newVibe := *vibe
		
		// Apply ternary energy modulation
		switch state {
		case TernaryNegative:
			newVibe.Energy = max(0.0, vibe.Energy-0.1)
		case TernaryNeutral:
			// Maintain energy level
		case TernaryPositive:
			newVibe.Energy = min(1.0, vibe.Energy+0.1)
		}
		
		// Determine next state based on energy level
		var nextState TernaryState
		if newVibe.Energy < 0.33 {
			nextState = TernaryNegative
		} else if newVibe.Energy > 0.66 {
			nextState = TernaryPositive
		} else {
			nextState = TernaryNeutral
		}
		
		return &newVibe, nextState
	})
}

// StreamTransformer provides composable stream transformations with memory efficiency
type StreamTransformer[T any] struct {
	source      chan T
	transforms  []Transducer[T, T]
	bufferPool  sync.Pool
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewStreamTransformer creates a composable stream transformer
func NewStreamTransformer[T any](capacity int) *StreamTransformer[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return &StreamTransformer[T]{
		source: make(chan T, capacity),
		ctx:    ctx,
		cancel: cancel,
		bufferPool: sync.Pool{
			New: func() interface{} {
				return make([]T, 0, 256) // Pre-allocate slice capacity
			},
		},
	}
}

// Chain adds a transducer to the transformation pipeline
func (st *StreamTransformer[T]) Chain(t Transducer[T, T]) *StreamTransformer[T] {
	st.transforms = append(st.transforms, t)
	return st
}

// Process starts the transformation pipeline with optimal memory usage
func (st *StreamTransformer[T]) Process(output chan<- T) {
	defer close(output)
	
	// Get buffer from pool
	buffer := st.bufferPool.Get().([]T)
	defer st.bufferPool.Put(buffer[:0]) // Reset slice but keep capacity
	
	for {
		select {
		case input, ok := <-st.source:
			if !ok {
				return
			}
			
			// Apply all transformations in sequence
			current := input
			for _, transform := range st.transforms {
				// Create a simple accumulating reducer
				reducer := func(acc T, val T) (T, bool) {
					return val, false // Never terminate early
				}
				
				// Apply transformation
				transformed := transform(reducer)
				var zero T
				current, _ = transformed(zero, current)
			}
			
			select {
			case output <- current:
			case <-st.ctx.Done():
				return
			}
			
		case <-st.ctx.Done():
			return
		}
	}
}

// HardwareOptimizedVibe provides SIMD-friendly vibe processing for TinyGo
type HardwareOptimizedVibe struct {
	// Pack data for cache efficiency (64-byte cache line)
	ID          [16]byte  // Fixed-size ID for SIMD operations
	Energy      float32   // 32-bit for SIMD vector operations
	Temperature float32
	Humidity    float32
	Light       float32
	Checksum    uint32    // Hardware CRC32 if available
	_padding    [12]byte  // Align to cache line
}

// ConvertToOptimized converts a standard Vibe to hardware-optimized format
func ConvertToOptimized(vibe *models.Vibe) *HardwareOptimizedVibe {
	opt := &HardwareOptimizedVibe{
		Energy: float32(vibe.Energy),
	}
	
	// Copy ID with bounds checking
	copy(opt.ID[:], []byte(vibe.ID))
	
	// Extract sensor data if available
	if vibe.SensorData.Temperature != nil {
		opt.Temperature = float32(*vibe.SensorData.Temperature)
	}
	if vibe.SensorData.Humidity != nil {
		opt.Humidity = float32(*vibe.SensorData.Humidity)
	}
	if vibe.SensorData.Light != nil {
		opt.Light = float32(*vibe.SensorData.Light)
	}
	
	// Compute checksum for integrity (can use hardware CRC32 on supported platforms)
	opt.computeChecksum()
	
	return opt
}

// computeChecksum calculates a safe checksum without unsafe pointer operations
func (hv *HardwareOptimizedVibe) computeChecksum() {
	// Use a safe approach that doesn't rely on unsafe pointer operations
	// to avoid memory alignment issues and potential crashes
	var sum uint32
	
	// Hash the individual fields safely
	for _, b := range hv.ID {
		sum ^= uint32(b)
	}
	sum ^= math.Float32bits(hv.Energy)
	sum ^= math.Float32bits(hv.Temperature)
	sum ^= math.Float32bits(hv.Humidity)
	sum ^= math.Float32bits(hv.Light)
	
	hv.Checksum = sum
}

// Validate checks data integrity using hardware acceleration where possible
func (hv *HardwareOptimizedVibe) Validate() bool {
	oldChecksum := hv.Checksum
	hv.computeChecksum()
	valid := hv.Checksum == oldChecksum
	hv.Checksum = oldChecksum // Restore original
	return valid
}

// QuantumInspiredProcessor provides quantum-inspired parallel processing
type QuantumInspiredProcessor struct {
	superposition []TernaryState // Quantum-like state superposition
	entanglement  map[int]int    // Entangled state pairs
	measurement   chan TernaryState
	coherence     int64 // Coherence time in nanoseconds
}

// NewQuantumInspiredProcessor creates a quantum-inspired processor for ternary logic
func NewQuantumInspiredProcessor(qubits int) *QuantumInspiredProcessor {
	return &QuantumInspiredProcessor{
		superposition: make([]TernaryState, qubits),
		entanglement:  make(map[int]int),
		measurement:   make(chan TernaryState, qubits),
		coherence:     1000000, // 1ms default coherence time
	}
}

// Entangle creates quantum-like correlations between ternary states
func (qip *QuantumInspiredProcessor) Entangle(qubit1, qubit2 int) {
	qip.entanglement[qubit1] = qubit2
	qip.entanglement[qubit2] = qubit1
}

// Evolve applies unitary evolution to the ternary system
func (qip *QuantumInspiredProcessor) Evolve(vibes []*models.Vibe) []*models.Vibe {
	if len(vibes) != len(qip.superposition) {
		return vibes // Dimension mismatch
	}
	
	// Prepare superposition based on vibe energies
	for i, vibe := range vibes {
		if vibe == nil {
			qip.superposition[i] = TernaryNeutral
			continue
		}
		
		if vibe.Energy < 0.33 {
			qip.superposition[i] = TernaryNegative
		} else if vibe.Energy > 0.66 {
			qip.superposition[i] = TernaryPositive
		} else {
			qip.superposition[i] = TernaryNeutral
		}
	}
	
	// Apply entanglement correlations
	for qubit1, qubit2 := range qip.entanglement {
		if qubit1 < len(qip.superposition) && qubit2 < len(qip.superposition) {
			// Simple entanglement: if one changes, the other mirrors
			if qip.superposition[qubit1] != TernaryNeutral {
				qip.superposition[qubit2] = -qip.superposition[qubit1]
			}
		}
	}
	
	// Apply evolved states back to vibes
	result := make([]*models.Vibe, len(vibes))
	for i, vibe := range vibes {
		if vibe == nil {
			result[i] = nil
			continue
		}
		
		newVibe := *vibe // Copy
		state := qip.superposition[i]
		
		// Translate ternary state back to energy
		switch state {
		case TernaryNegative:
			newVibe.Energy = max(0.0, vibe.Energy-0.2)
		case TernaryNeutral:
			// Maintain energy
		case TernaryPositive:
			newVibe.Energy = min(1.0, vibe.Energy+0.2)
		}
		
		result[i] = &newVibe
	}
	
	return result
}

// RTStreamProcessor provides real-time stream processing with deterministic latency
type RTStreamProcessor struct {
	inputBuffer  *WorldMomentStream
	outputBuffer *WorldMomentStream
	processor    func(*models.WorldMoment) *models.WorldMoment
	deadline     time.Duration
	stats        struct {
		processed uint64
		dropped   uint64
		maxLatency time.Duration
	}
}

// NewRTStreamProcessor creates a real-time processor with deadline guarantees
func NewRTStreamProcessor(capacity uint64, deadline time.Duration) *RTStreamProcessor {
	return &RTStreamProcessor{
		inputBuffer:  NewWorldMomentStream(capacity),
		outputBuffer: NewWorldMomentStream(capacity),
		deadline:     deadline,
	}
}

// ProcessRealTime runs the real-time processing loop with deadline enforcement
func (rtp *RTStreamProcessor) ProcessRealTime(ctx context.Context) {
	ticker := time.NewTicker(rtp.deadline / 10) // Process at 10x deadline frequency
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			start := time.Now()
			
			// Process one item with deadline enforcement
			if moment, ok := rtp.inputBuffer.Pop(); ok {
				processed := rtp.processor(moment)
				
				elapsed := time.Since(start)
				if elapsed <= rtp.deadline {
					// Within deadline, try to output
					if !rtp.outputBuffer.Push(processed) {
						atomic.AddUint64(&rtp.stats.dropped, 1)
					} else {
						atomic.AddUint64(&rtp.stats.processed, 1)
						
						// Update max latency stats
						if elapsed > rtp.stats.maxLatency {
							rtp.stats.maxLatency = elapsed
						}
					}
				} else {
					// Deadline missed, drop the item
					atomic.AddUint64(&rtp.stats.dropped, 1)
				}
			}
			
		case <-ctx.Done():
			return
		}
	}
}

// Memory pool for zero-allocation processing
var (
	vibePool = sync.Pool{
		New: func() interface{} {
			return &models.Vibe{}
		},
	}
	
	momentPool = sync.Pool{
		New: func() interface{} {
			return &models.WorldMoment{}
		},
	}
)

// GetVibe retrieves a vibe from the pool for zero-allocation processing
func GetVibe() *models.Vibe {
	return vibePool.Get().(*models.Vibe)
}

// PutVibe returns a vibe to the pool after clearing sensitive data
func PutVibe(vibe *models.Vibe) {
	// Clear the vibe but keep allocated slices
	*vibe = models.Vibe{}
	vibePool.Put(vibe)
}

// GetMoment retrieves a world moment from the pool
func GetMoment() *models.WorldMoment {
	return momentPool.Get().(*models.WorldMoment)
}

// PutMoment returns a world moment to the pool
func PutMoment(moment *models.WorldMoment) {
	// Clear but preserve slice capacity
	moment.Viewers = moment.Viewers[:0]
	*moment = models.WorldMoment{Viewers: moment.Viewers}
	momentPool.Put(moment)
}

// CPUOptimizedBatch processes vibes in CPU-cache-friendly batches
func CPUOptimizedBatch(vibes []*models.Vibe, batchSize int, processor func([]*models.Vibe)) {
	// Use runtime.GOMAXPROCS to determine optimal parallelism
	numCPU := runtime.GOMAXPROCS(0)
	
	// Process in cache-friendly batches
	for i := 0; i < len(vibes); i += batchSize {
		end := min(i+batchSize, len(vibes))
		batch := vibes[i:end]
		
		// Use worker pool based on CPU count
		if len(batch) > numCPU {
			// Parallel processing for large batches
			var wg sync.WaitGroup
			chunkSize := len(batch) / numCPU
			
			for j := 0; j < numCPU; j++ {
				start := j * chunkSize
				chunkEnd := start + chunkSize
				if j == numCPU-1 {
					chunkEnd = len(batch) // Handle remainder
				}
				
				if start < len(batch) {
					wg.Add(1)
					go func(chunk []*models.Vibe) {
						defer wg.Done()
						processor(chunk)
					}(batch[start:chunkEnd])
				}
			}
			
			wg.Wait()
		} else {
			// Sequential processing for small batches
			processor(batch)
		}
	}
}

// Helper functions for min/max that work with TinyGo
func min[T ~int | ~float32 | ~float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T ~int | ~float32 | ~float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// ComonadicVibeContext represents a vibe within its spatial-temporal context
type ComonadicVibeContext struct {
	center     *models.Vibe
	neighbors  []*models.Vibe
	temporal   []*models.Vibe // Past and future states
	gradient   TernaryState   // Directional tendency
	coherence  float64        // Context coherence measure
}

// Extract implements the comonadic extract operation - get the focused value
func (cvc *ComonadicVibeContext) Extract() *models.Vibe {
	return cvc.center
}

// Duplicate implements comonadic duplication - create context-of-contexts
func (cvc *ComonadicVibeContext) Duplicate() *ComonadicVibeContext {
	// Create a meta-context where each neighbor becomes a center
	neighborContexts := make([]*models.Vibe, len(cvc.neighbors))
	for i, neighbor := range cvc.neighbors {
		// Each neighbor gets elevated to be a potential center
		neighborContexts[i] = neighbor
	}
	
	return &ComonadicVibeContext{
		center:    cvc.center,
		neighbors: neighborContexts,
		temporal:  cvc.temporal,
		gradient:  cvc.gradient,
		coherence: cvc.coherence * 0.9, // Slight coherence decay in duplication
	}
}

// Extend implements comonadic extension - apply context-aware transformation
func (cvc *ComonadicVibeContext) Extend(f func(*ComonadicVibeContext) *models.Vibe) *ComonadicVibeContext {
	newCenter := f(cvc)
	
	// Apply the transformation to each neighbor context
	newNeighbors := make([]*models.Vibe, len(cvc.neighbors))
	for i, neighbor := range cvc.neighbors {
		// Create mini-context for each neighbor
		neighborCtx := &ComonadicVibeContext{
			center:    neighbor,
			neighbors: []*models.Vibe{cvc.center}, // Neighbor sees center as its neighbor
			temporal:  cvc.temporal,
			gradient:  -cvc.gradient, // Invert gradient for neighbor perspective
			coherence: cvc.coherence,
		}
		newNeighbors[i] = f(neighborCtx)
	}
	
	return &ComonadicVibeContext{
		center:    newCenter,
		neighbors: newNeighbors,
		temporal:  cvc.temporal,
		gradient:  cvc.gradient,
		coherence: cvc.coherence,
	}
}

// TernaryLogicGate represents sophisticated ternary operations
type TernaryLogicGate struct {
	operator string
	truthTable map[[2]TernaryState]TernaryState
}

// NewTernaryLogicGate creates a gate with specified truth table
func NewTernaryLogicGate(operator string) *TernaryLogicGate {
	gate := &TernaryLogicGate{
		operator: operator,
		truthTable: make(map[[2]TernaryState]TernaryState),
	}
	
	// Define truth tables for different operators
	switch operator {
	case "consensus":
		// Ternary consensus: agrees when inputs agree, neutral otherwise
		gate.defineTruthTable(map[[2]TernaryState]TernaryState{
			{TernaryNegative, TernaryNegative}: TernaryNegative,
			{TernaryNeutral, TernaryNeutral}:   TernaryNeutral,
			{TernaryPositive, TernaryPositive}: TernaryPositive,
			// All other combinations default to neutral
		})
	case "amplify":
		// Amplification: strengthens agreement, weakens disagreement
		gate.defineTruthTable(map[[2]TernaryState]TernaryState{
			{TernaryNegative, TernaryNegative}: TernaryNegative,
			{TernaryPositive, TernaryPositive}: TernaryPositive,
			{TernaryNegative, TernaryPositive}: TernaryNeutral,
			{TernaryPositive, TernaryNegative}: TernaryNeutral,
			// Neutral input preserves the other
			{TernaryNeutral, TernaryNegative}:  TernaryNegative,
			{TernaryNeutral, TernaryPositive}:  TernaryPositive,
			{TernaryNegative, TernaryNeutral}:  TernaryNegative,
			{TernaryPositive, TernaryNeutral}:  TernaryPositive,
			{TernaryNeutral, TernaryNeutral}:   TernaryNeutral,
		})
	case "inhibit":
		// Inhibition: positive inhibits, negative disinhibits
		gate.defineTruthTable(map[[2]TernaryState]TernaryState{
			{TernaryNegative, TernaryPositive}: TernaryNegative, // Negative disinhibits positive
			{TernaryPositive, TernaryNegative}: TernaryNeutral,  // Positive inhibits negative
			{TernaryPositive, TernaryPositive}: TernaryNeutral,  // Positive inhibits positive
			{TernaryNegative, TernaryNegative}: TernaryNegative, // Double negative
			{TernaryNeutral, TernaryNeutral}:   TernaryNeutral,
			{TernaryNeutral, TernaryPositive}:  TernaryNeutral,
			{TernaryNeutral, TernaryNegative}:  TernaryNegative,
			{TernaryPositive, TernaryNeutral}:  TernaryNeutral,
			{TernaryNegative, TernaryNeutral}:  TernaryNegative,
		})
	default:
		// Default to identity operation
		for _, a := range []TernaryState{TernaryNegative, TernaryNeutral, TernaryPositive} {
			for _, b := range []TernaryState{TernaryNegative, TernaryNeutral, TernaryPositive} {
				gate.truthTable[[2]TernaryState{a, b}] = a // First input wins
			}
		}
	}
	
	return gate
}

// defineTruthTable helper to set up truth table with defaults
func (tlg *TernaryLogicGate) defineTruthTable(table map[[2]TernaryState]TernaryState) {
	// Initialize all combinations to neutral by default
	for _, a := range []TernaryState{TernaryNegative, TernaryNeutral, TernaryPositive} {
		for _, b := range []TernaryState{TernaryNegative, TernaryNeutral, TernaryPositive} {
			tlg.truthTable[[2]TernaryState{a, b}] = TernaryNeutral
		}
	}
	
	// Override with specific table
	for key, value := range table {
		tlg.truthTable[key] = value
	}
}

// Apply executes the ternary logic operation
func (tlg *TernaryLogicGate) Apply(a, b TernaryState) TernaryState {
	if result, exists := tlg.truthTable[[2]TernaryState{a, b}]; exists {
		return result
	}
	return TernaryNeutral // Safe default
}

// VibeContextualTransformer applies context-aware transformations using comonads
type VibeContextualTransformer struct {
	logicGates map[string]*TernaryLogicGate
	history    []*models.Vibe // Temporal context
	maxHistory int
}

// NewVibeContextualTransformer creates a sophisticated vibe processor
func NewVibeContextualTransformer(maxHistory int) *VibeContextualTransformer {
	return &VibeContextualTransformer{
		logicGates: map[string]*TernaryLogicGate{
			"consensus": NewTernaryLogicGate("consensus"),
			"amplify":   NewTernaryLogicGate("amplify"),
			"inhibit":   NewTernaryLogicGate("inhibit"),
		},
		history:    make([]*models.Vibe, 0, maxHistory),
		maxHistory: maxHistory,
	}
}

// TransformWithContext applies comonadic transformations to vibes
func (vct *VibeContextualTransformer) TransformWithContext(center *models.Vibe, neighbors []*models.Vibe) *models.Vibe {
	// Create comonadic context
	ctx := &ComonadicVibeContext{
		center:    center,
		neighbors: neighbors,
		temporal:  append([]*models.Vibe{}, vct.history...), // Copy history
		gradient:  vct.calculateGradient(center),
		coherence: vct.calculateCoherence(center, neighbors),
	}
	
	// Apply comonadic extension with context-aware transformation
	transformed := ctx.Extend(func(c *ComonadicVibeContext) *models.Vibe {
		return vct.contextAwareTransform(c)
	})
	
	// Update history
	vct.updateHistory(center)
	
	return transformed.Extract()
}

// calculateGradient determines the directional tendency of a vibe
func (vct *VibeContextualTransformer) calculateGradient(vibe *models.Vibe) TernaryState {
	if len(vct.history) < 2 {
		return TernaryNeutral
	}
	
	// Compare with recent history
	recentEnergy := vct.history[len(vct.history)-1].Energy
	if vibe.Energy > recentEnergy + 0.1 {
		return TernaryPositive // Rising
	} else if vibe.Energy < recentEnergy - 0.1 {
		return TernaryNegative // Falling
	}
	return TernaryNeutral // Stable
}

// calculateCoherence measures how well the vibe fits with its context
func (vct *VibeContextualTransformer) calculateCoherence(center *models.Vibe, neighbors []*models.Vibe) float64 {
	if len(neighbors) == 0 {
		return 1.0 // Perfect coherence in isolation
	}
	
	// Calculate average neighbor energy
	var totalEnergy float64
	for _, neighbor := range neighbors {
		if neighbor != nil {
			totalEnergy += neighbor.Energy
		}
	}
	avgEnergy := totalEnergy / float64(len(neighbors))
	
	// Coherence is inverse of energy difference
	energy_diff := math.Abs(center.Energy - avgEnergy)
	return math.Exp(-energy_diff) // Exponential decay with difference
}

// contextAwareTransform performs the actual transformation using context
func (vct *VibeContextualTransformer) contextAwareTransform(ctx *ComonadicVibeContext) *models.Vibe {
	vibe := *ctx.center // Copy
	
	// Apply different gates based on context
	if ctx.coherence > 0.8 {
		// High coherence: use consensus
		for _, neighbor := range ctx.neighbors {
			if neighbor != nil {
				centerState := energyToTernary(vibe.Energy)
				neighborState := energyToTernary(neighbor.Energy)
				result := vct.logicGates["consensus"].Apply(centerState, neighborState)
				vibe.Energy = ternaryToEnergy(result, vibe.Energy)
				break // Only apply to first neighbor for simplicity
			}
		}
	} else if ctx.gradient != TernaryNeutral {
		// Trending: use amplification
		centerState := energyToTernary(vibe.Energy)
		result := vct.logicGates["amplify"].Apply(centerState, ctx.gradient)
		vibe.Energy = ternaryToEnergy(result, vibe.Energy)
	} else {
		// Low coherence: use inhibition to dampen chaos
		for _, neighbor := range ctx.neighbors {
			if neighbor != nil {
				centerState := energyToTernary(vibe.Energy)
				neighborState := energyToTernary(neighbor.Energy)
				result := vct.logicGates["inhibit"].Apply(centerState, neighborState)
				vibe.Energy = ternaryToEnergy(result, vibe.Energy)
				break
			}
		}
	}
	
	return &vibe
}

// energyToTernary converts energy level to ternary state
func energyToTernary(energy float64) TernaryState {
	if energy < 0.33 {
		return TernaryNegative
	} else if energy > 0.66 {
		return TernaryPositive
	}
	return TernaryNeutral
}

// ternaryToEnergy converts ternary state back to energy with smoothing
func ternaryToEnergy(state TernaryState, currentEnergy float64) float64 {
	var targetEnergy float64
	
	switch state {
	case TernaryNegative:
		targetEnergy = 0.2
	case TernaryNeutral:
		targetEnergy = 0.5
	case TernaryPositive:
		targetEnergy = 0.8
	}
	
	// Smooth transition - don't jump directly to target
	smoothing := 0.3
	return currentEnergy*(1-smoothing) + targetEnergy*smoothing
}

// updateHistory maintains temporal context
func (vct *VibeContextualTransformer) updateHistory(vibe *models.Vibe) {
	vct.history = append(vct.history, vibe)
	if len(vct.history) > vct.maxHistory {
		// Remove oldest entry
		vct.history = vct.history[1:]
	}
}
