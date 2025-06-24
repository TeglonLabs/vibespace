package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/repository"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ProgressTracker manages progress updates for long-running operations
type ProgressTracker struct {
	mu       sync.Mutex
	progress map[string]*ProgressState
}

type ProgressState struct {
	TaskID      string    `json:"task_id"`
	Progress    float64   `json:"progress"`
	Message     string    `json:"message"`
	StartTime   time.Time `json:"start_time"`
	UpdateTime  time.Time `json:"update_time"`
	Cancelled   bool      `json:"cancelled"`
	Completed   bool      `json:"completed"`
	Error       string    `json:"error,omitempty"`
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		progress: make(map[string]*ProgressState),
	}
}

// StartTask begins tracking a new task
func (pt *ProgressTracker) StartTask(taskID string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	pt.progress[taskID] = &ProgressState{
		TaskID:     taskID,
		Progress:   0.0,
		Message:    "Starting task...",
		StartTime:  time.Now(),
		UpdateTime: time.Now(),
	}
}

// UpdateProgress updates the progress of a task
func (pt *ProgressTracker) UpdateProgress(taskID string, progress float64, message string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	if state, exists := pt.progress[taskID]; exists {
		state.Progress = progress
		state.Message = message
		state.UpdateTime = time.Now()
	}
}

// CancelTask marks a task as cancelled
func (pt *ProgressTracker) CancelTask(taskID string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	if state, exists := pt.progress[taskID]; exists {
		state.Cancelled = true
		state.UpdateTime = time.Now()
		state.Message = "Task cancelled"
	}
}

// CompleteTask marks a task as completed
func (pt *ProgressTracker) CompleteTask(taskID string, message string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	if state, exists := pt.progress[taskID]; exists {
		state.Completed = true
		state.Progress = 100.0
		state.Message = message
		state.UpdateTime = time.Now()
	}
}

// GetProgress returns the current progress state
func (pt *ProgressTracker) GetProgress(taskID string) (*ProgressState, bool) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	state, exists := pt.progress[taskID]
	if !exists {
		return nil, false
	}
	
	// Return a copy to avoid race conditions
	stateCopy := *state
	return &stateCopy, true
}

// IsCancelled checks if a task is cancelled
func (pt *ProgressTracker) IsCancelled(taskID string) bool {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	if state, exists := pt.progress[taskID]; exists {
		return state.Cancelled
	}
	return false
}

// Global progress tracker for testing
var globalProgressTracker = NewProgressTracker()

// setupProgressTrackingServer creates a test server with progress tracking capabilities
func setupProgressTrackingServer(t *testing.T) *server.MCPServer {
	mcpServer := server.NewMCPServer("vibespace-progress-test", "1.0.0")
	
	// Long-running analysis tool with progress tracking
	analysisTool := mcp.NewTool("analyze_large_dataset", func(t *mcp.Tool) {
		t.Description = "Analyze large dataset with progress tracking and cancellation support"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"dataset_size": {Type: "integer", Description: "Number of records to analyze"},
				"analysis_type": {Type: "string", Enum: []string{"full", "partial", "summary"}, Description: "Type of analysis"},
				"batch_size": {Type: "integer", Description: "Records per batch (default: 100)"},
			},
			Required: []string{"dataset_size", "analysis_type"},
		}
	})
	
	mcpServer.AddTool(analysisTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		datasetSize, _ := args["dataset_size"].(float64)
		analysisType, _ := args["analysis_type"].(string)
		batchSize, _ := args["batch_size"].(float64)
		
		if batchSize == 0 {
			batchSize = 100
		}
		
		// Generate unique task ID
		taskID := fmt.Sprintf("analysis_%d_%s", time.Now().Unix(), analysisType)
		
		// Start progress tracking
		globalProgressTracker.StartTask(taskID)
		
		// Simulate long-running analysis with progress updates
		go func() {
			totalBatches := int(datasetSize / batchSize)
			if totalBatches == 0 {
				totalBatches = 1
			}
			
			for i := 0; i < totalBatches; i++ {
				// Check for cancellation
				if globalProgressTracker.IsCancelled(taskID) {
					return
				}
				
				// Simulate batch processing time
				time.Sleep(50 * time.Millisecond)
				
				progress := float64(i+1) / float64(totalBatches) * 100
				message := fmt.Sprintf("Processed batch %d/%d", i+1, totalBatches)
				globalProgressTracker.UpdateProgress(taskID, progress, message)
			}
			
			// Complete the task
			if !globalProgressTracker.IsCancelled(taskID) {
				globalProgressTracker.CompleteTask(taskID, "Analysis completed successfully")
			}
		}()
		
		result := map[string]interface{}{
			"task_id": taskID,
			"status": "started",
			"dataset_size": datasetSize,
			"analysis_type": analysisType,
			"batch_size": batchSize,
			"estimated_duration": fmt.Sprintf("%.1f seconds", datasetSize/batchSize*0.05),
		}
		
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
	
	// Progress query tool
	progressTool := mcp.NewTool("get_task_progress", func(t *mcp.Tool) {
		t.Description = "Get current progress of a running task"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"task_id": {Type: "string", Description: "Task identifier"},
			},
			Required: []string{"task_id"},
		}
	})
	
	mcpServer.AddTool(progressTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		taskID, _ := args["task_id"].(string)
		
		state, exists := globalProgressTracker.GetProgress(taskID)
		if !exists {
			return mcp.NewToolResultError("Task not found"), nil
		}
		
		stateJSON, _ := json.Marshal(state)
		return mcp.NewToolResultText(string(stateJSON)), nil
	})
	
	// Cancellation tool
	cancelTool := mcp.NewTool("cancel_task", func(t *mcp.Tool) {
		t.Description = "Cancel a running task"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"task_id": {Type: "string", Description: "Task identifier to cancel"},
			},
			Required: []string{"task_id"},
		}
	})
	
	mcpServer.AddTool(cancelTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		taskID, _ := args["task_id"].(string)
		
		globalProgressTracker.CancelTask(taskID)
		
		result := map[string]interface{}{
			"task_id": taskID,
			"status": "cancelled",
			"cancelled_at": time.Now().Format(time.RFC3339),
		}
		
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
	
	// Batch processing tool with detailed progress
	batchTool := mcp.NewTool("process_files_batch", func(t *mcp.Tool) {
		t.Description = "Process files in batches with detailed progress tracking"
		t.InputSchema = &mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]mcp.ToolInputSchema{
				"file_count": {Type: "integer", Description: "Number of files to process"},
				"processing_type": {Type: "string", Enum: []string{"image", "video", "audio", "document"}, Description: "Type of files"},
				"parallel_workers": {Type: "integer", Description: "Number of parallel workers (default: 2)"},
			},
			Required: []string{"file_count", "processing_type"},
		}
	})
	
	mcpServer.AddTool(batchTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}
		
		fileCount, _ := args["file_count"].(float64)
		processingType, _ := args["processing_type"].(string)
		workers, _ := args["parallel_workers"].(float64)
		
		if workers == 0 {
			workers = 2
		}
		
		taskID := fmt.Sprintf("batch_%s_%d", processingType, time.Now().Unix())
		globalProgressTracker.StartTask(taskID)
		
		// Simulate batch processing
		go func() {
			filesPerWorker := int(fileCount / workers)
			var wg sync.WaitGroup
			processedFiles := 0
			var mu sync.Mutex
			
			for w := 0; w < int(workers); w++ {
				wg.Add(1)
				go func(workerID int) {
					defer wg.Done()
					
					start := workerID * filesPerWorker
					end := start + filesPerWorker
					if workerID == int(workers)-1 {
						end = int(fileCount) // Last worker handles remaining files
					}
					
					for i := start; i < end; i++ {
						if globalProgressTracker.IsCancelled(taskID) {
							return
						}
						
						// Simulate file processing time based on type
						processingTime := map[string]time.Duration{
							"image":    30 * time.Millisecond,
							"video":    100 * time.Millisecond,
							"audio":    50 * time.Millisecond,
							"document": 20 * time.Millisecond,
						}[processingType]
						
						time.Sleep(processingTime)
						
						mu.Lock()
						processedFiles++
						progress := float64(processedFiles) / fileCount * 100
						message := fmt.Sprintf("Worker %d: Processed %d/%d files (%s)", 
							workerID+1, processedFiles, int(fileCount), processingType)
						globalProgressTracker.UpdateProgress(taskID, progress, message)
						mu.Unlock()
					}
				}(w)
			}
			
			wg.Wait()
			
			if !globalProgressTracker.IsCancelled(taskID) {
				globalProgressTracker.CompleteTask(taskID, 
					fmt.Sprintf("Batch processing completed: %d %s files processed", 
						int(fileCount), processingType))
			}
		}()
		
		result := map[string]interface{}{
			"task_id": taskID,
			"status": "started",
			"file_count": fileCount,
			"processing_type": processingType,
			"parallel_workers": workers,
			"estimated_duration": fmt.Sprintf("%.1f seconds", fileCount*0.05/workers),
		}
		
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
	
	return mcpServer
}

// Test progress tracking for long-running operations
func TestProgressTracking(t *testing.T) {
	mcpServer := setupProgressTrackingServer(t)
	
	t.Run("Start Analysis Task with Progress Tracking", func(t *testing.T) {
		// Start analysis task
		request := mcp.CallToolRequest{
			Method: "tools/call",
			Params: mcp.CallToolParams{
				Name: "analyze_large_dataset",
				Arguments: map[string]interface{}{
					"dataset_size": 500.0,
					"analysis_type": "full",
					"batch_size": 50.0,
				},
			},
		}
		
		result, err := mcpServer.HandleToolCall(context.Background(), request)
		require.NoError(t, err)
		
		// Parse result to get task ID
		content := result.Content[0]
		var taskResult map[string]interface{}
		err = json.Unmarshal([]byte(content.Text), &taskResult)
		require.NoError(t, err)
		
		taskID := taskResult["task_id"].(string)
		assert.NotEmpty(t, taskID)
		assert.Equal(t, "started", taskResult["status"])
		assert.Equal(t, 500.0, taskResult["dataset_size"])
		
		// Wait a bit for processing to start
		time.Sleep(100 * time.Millisecond)
		
		// Check progress
		progressRequest := mcp.CallToolRequest{
			Method: "tools/call",
			Params: mcp.CallToolParams{
				Name: "get_task_progress",
				Arguments: map[string]interface{}{
					"task_id": taskID,
				},
			},
		}
		
		progressResult, err := mcpServer.HandleToolCall(context.Background(), progressRequest)
		require.NoError(t, err)
		
		var progressState ProgressState
		err = json.Unmarshal([]byte(progressResult.Content[0].Text), &progressState)
		require.NoError(t, err)
		
		assert.Equal(t, taskID, progressState.TaskID)
		assert.Greater(t, progressState.Progress, 0.0)
		assert.False(t, progressState.Cancelled)
		assert.Contains(t, progressState.Message, "batch")
		
		// Wait for completion
		for i := 0; i < 20; i++ {
			time.Sleep(100 * time.Millisecond)
			
			progressResult, err = mcpServer.HandleToolCall(context.Background(), progressRequest)
			require.NoError(t, err)
			
			err = json.Unmarshal([]byte(progressResult.Content[0].Text), &progressState)
			require.NoError(t, err)
			
			if progressState.Completed {
				break
			}
		}
		
		assert.True(t, progressState.Completed)
		assert.Equal(t, 100.0, progressState.Progress)
		assert.Contains(t, progressState.Message, "completed")
	})
	
	t.Run("Cancel Running Task", func(t *testing.T) {
		// Start a longer task
		request := mcp.CallToolRequest{
			Method: "tools/call",
			Params: mcp.CallToolParams{
				Name: "analyze_large_dataset",
				Arguments: map[string]interface{}{
					"dataset_size": 1000.0,
					"analysis_type": "full",
					"batch_size": 25.0,
				},
			},
		}
		
		result, err := mcpServer.HandleToolCall(context.Background(), request)
		require.NoError(t, err)
		
		var taskResult map[string]interface{}
		err = json.Unmarshal([]byte(result.Content[0].Text), &taskResult)
		require.NoError(t, err)
		
		taskID := taskResult["task_id"].(string)
		
		// Wait for task to start processing
		time.Sleep(150 * time.Millisecond)
		
		// Cancel the task
		cancelRequest := mcp.CallToolRequest{
			Method: "tools/call",
			Params: mcp.CallToolParams{
				Name: "cancel_task",
				Arguments: map[string]interface{}{
					"task_id": taskID,
				},
			},
		}
		
		cancelResult, err := mcpServer.HandleToolCall(context.Background(), cancelRequest)
		require.NoError(t, err)
		
		var cancelResponse map[string]interface{}
		err = json.Unmarshal([]byte(cancelResult.Content[0].Text), &cancelResponse)
		require.NoError(t, err)
		
		assert.Equal(t, taskID, cancelResponse["task_id"])
		assert.Equal(t, "cancelled", cancelResponse["status"])
		
		// Verify task is cancelled
		progressRequest := mcp.CallToolRequest{
			Method: "tools/call",
			Params: mcp.CallToolParams{
				Name: "get_task_progress",
				Arguments: map[string]interface{}{
					"task_id": taskID,
				},
			},
		}
		
		progressResult, err := mcpServer.HandleToolCall(context.Background(), progressRequest)
		require.NoError(t, err)
		
		var progressState ProgressState
		err = json.Unmarshal([]byte(progressResult.Content[0].Text), &progressState)
		require.NoError(t, err)
		
		assert.True(t, progressState.Cancelled)
		assert.Equal(t, "Task cancelled", progressState.Message)
		assert.False(t, progressState.Completed)
	})
	
	t.Run("Parallel Batch Processing with Progress", func(t *testing.T) {
		// Test parallel file processing following 2-3-5-7 principle
		testCases := []struct {
			fileCount int
			fileType  string
			workers   int
		}{
			{200, "image", 2},
			{300, "video", 3},
			{500, "audio", 5},
			{700, "document", 7},
		}
		
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s_%d_files_%d_workers", tc.fileType, tc.fileCount, tc.workers), func(t *testing.T) {
				request := mcp.CallToolRequest{
					Method: "tools/call",
					Params: mcp.CallToolParams{
						Name: "process_files_batch",
						Arguments: map[string]interface{}{
							"file_count": float64(tc.fileCount),
							"processing_type": tc.fileType,
							"parallel_workers": float64(tc.workers),
						},
					},
				}
				
				result, err := mcpServer.HandleToolCall(context.Background(), request)
				require.NoError(t, err)
				
				var taskResult map[string]interface{}
				err = json.Unmarshal([]byte(result.Content[0].Text), &taskResult)
				require.NoError(t, err)
				
				taskID := taskResult["task_id"].(string)
				assert.Equal(t, float64(tc.fileCount), taskResult["file_count"])
				assert.Equal(t, tc.fileType, taskResult["processing_type"])
				assert.Equal(t, float64(tc.workers), taskResult["parallel_workers"])
				
				// Monitor progress until completion
				progressRequest := mcp.CallToolRequest{
					Method: "tools/call",
					Params: mcp.CallToolParams{
						Name: "get_task_progress",
						Arguments: map[string]interface{}{
							"task_id": taskID,
						},
					},
				}
				
				var finalProgress ProgressState
				for i := 0; i < 100; i++ {
					time.Sleep(50 * time.Millisecond)
					
					progressResult, err := mcpServer.HandleToolCall(context.Background(), progressRequest)
					require.NoError(t, err)
					
					err = json.Unmarshal([]byte(progressResult.Content[0].Text), &finalProgress)
					require.NoError(t, err)
					
					if finalProgress.Completed {
						break
					}
					
					// Validate progress is increasing
					assert.GreaterOrEqual(t, finalProgress.Progress, 0.0)
					assert.LessOrEqual(t, finalProgress.Progress, 100.0)
					assert.Contains(t, finalProgress.Message, "Worker")
				}
				
				assert.True(t, finalProgress.Completed)
				assert.Equal(t, 100.0, finalProgress.Progress)
				assert.Contains(t, finalProgress.Message, "Batch processing completed")
				assert.Contains(t, finalProgress.Message, tc.fileType)
			})
		}
	})
	
	t.Run("Progress Tracking Edge Cases", func(t *testing.T) {
		// Test with invalid task ID
		progressRequest := mcp.CallToolRequest{
			Method: "tools/call",
			Params: mcp.CallToolParams{
				Name: "get_task_progress",
				Arguments: map[string]interface{}{
					"task_id": "nonexistent-task-id",
				},
			},
		}
		
		result, err := mcpServer.HandleToolCall(context.Background(), progressRequest)
		require.NoError(t, err)
		
		assert.Equal(t, "error", result.Content[0].Type)
		assert.Contains(t, result.Content[0].Text, "Task not found")
		
		// Test cancelling non-existent task
		cancelRequest := mcp.CallToolRequest{
			Method: "tools/call",
			Params: mcp.CallToolParams{
				Name: "cancel_task",
				Arguments: map[string]interface{}{
					"task_id": "nonexistent-task-id",
				},
			},
		}
		
		cancelResult, err := mcpServer.HandleToolCall(context.Background(), cancelRequest)
		require.NoError(t, err)
		
		// Should succeed even for non-existent task (idempotent operation)
		var cancelResponse map[string]interface{}
		err = json.Unmarshal([]byte(cancelResult.Content[0].Text), &cancelResponse)
		require.NoError(t, err)
		assert.Equal(t, "cancelled", cancelResponse["status"])
	})
}

// Test concurrent progress tracking
func TestConcurrentProgressTracking(t *testing.T) {
	mcpServer := setupProgressTrackingServer(t)
	
	t.Run("Multiple Concurrent Tasks", func(t *testing.T) {
		// Start multiple tasks concurrently following 2-3-5 principle
		taskCounts := []int{200, 300, 500}
		var taskIDs []string
		var wg sync.WaitGroup
		
		// Start all tasks
		for i, count := range taskCounts {
			wg.Add(1)
			go func(idx, fileCount int) {
				defer wg.Done()
				
				request := mcp.CallToolRequest{
					Method: "tools/call",
					Params: mcp.CallToolParams{
						Name: "process_files_batch",
						Arguments: map[string]interface{}{
							"file_count": float64(fileCount),
							"processing_type": "image",
							"parallel_workers": 2.0,
						},
					},
				}
				
				result, err := mcpServer.HandleToolCall(context.Background(), request)
				require.NoError(t, err)
				
				var taskResult map[string]interface{}
				err = json.Unmarshal([]byte(result.Content[0].Text), &taskResult)
				require.NoError(t, err)
				
				taskIDs = append(taskIDs, taskResult["task_id"].(string))
			}(i, count)
		}
		
		wg.Wait()
		
		// Monitor all tasks until completion
		allCompleted := false
		for i := 0; i < 200 && !allCompleted; i++ {
			time.Sleep(50 * time.Millisecond)
			
			completedCount := 0
			for _, taskID := range taskIDs {
				progressRequest := mcp.CallToolRequest{
					Method: "tools/call",
					Params: mcp.CallToolParams{
						Name: "get_task_progress",
						Arguments: map[string]interface{}{
							"task_id": taskID,
						},
					},
				}
				
				progressResult, err := mcpServer.HandleToolCall(context.Background(), progressRequest)
				require.NoError(t, err)
				
				var progressState ProgressState
				err = json.Unmarshal([]byte(progressResult.Content[0].Text), &progressState)
				require.NoError(t, err)
				
				if progressState.Completed {
					completedCount++
				}
				
				// Validate concurrent execution doesn't interfere
				assert.GreaterOrEqual(t, progressState.Progress, 0.0)
				assert.LessOrEqual(t, progressState.Progress, 100.0)
				assert.False(t, progressState.Cancelled)
			}
			
			if completedCount == len(taskIDs) {
				allCompleted = true
			}
		}
		
		assert.True(t, allCompleted, "All concurrent tasks should complete")
		assert.Equal(t, 3, len(taskIDs), "Should have started 3 concurrent tasks")
	})
}

// Benchmark progress tracking performance
func BenchmarkProgressTracking(b *testing.B) {
	tracker := NewProgressTracker()
	
	b.Run("UpdateProgress", func(b *testing.B) {
		taskID := "benchmark-task"
		tracker.StartTask(taskID)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			progress := float64(i%100) + 1
			tracker.UpdateProgress(taskID, progress, fmt.Sprintf("Progress %d", i))
		}
	})
	
	b.Run("GetProgress", func(b *testing.B) {
		taskID := "benchmark-task-get"
		tracker.StartTask(taskID)
		tracker.UpdateProgress(taskID, 50.0, "Halfway done")
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, exists := tracker.GetProgress(taskID)
			if !exists {
				b.Fatal("Task should exist")
			}
		}
	})
	
	b.Run("ConcurrentAccess", func(b *testing.B) {
		taskID := "concurrent-benchmark"
		tracker.StartTask(taskID)
		
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				if i%2 == 0 {
					tracker.UpdateProgress(taskID, float64(i%100), fmt.Sprintf("Update %d", i))
				} else {
					tracker.GetProgress(taskID)
				}
				i++
			}
		})
	})
}
