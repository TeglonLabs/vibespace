package rpcmethods

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/bmorphism/vibespace-mcp-go/streaming"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// CategoricalTools provides MCP tools for categorical universe artifacts
type CategoricalTools struct {
	transformer *streaming.VibeContextualTransformer
}

// NewCategoricalTools creates a new categorical tools instance
func NewCategoricalTools() *CategoricalTools {
	return &CategoricalTools{
		transformer: streaming.NewVibeContextualTransformer(10),
	}
}

// CategoricalExtractRequest represents an extract operation request
type CategoricalExtractRequest struct {
	ContextID   string                           `json:"contextId"`
	VibeContext *streaming.ComonadicVibeContext  `json:"vibeContext,omitempty"`
}

// CategoricalExtractResponse represents the extracted value
type CategoricalExtractResponse struct {
	ExtractedVibe *models.Vibe `json:"extractedVibe"`
	ContextID     string       `json:"contextId"`
	Artifact      *WrapPreviewArtifact `json:"artifact"`
}

// CategoricalDuplicateRequest represents a duplicate operation request  
type CategoricalDuplicateRequest struct {
	ContextID string                           `json:"contextId"`
	Context   *streaming.ComonadicVibeContext  `json:"context"`
}

// CategoricalDuplicateResponse represents the duplicated context
type CategoricalDuplicateResponse struct {
	DuplicatedContext *streaming.ComonadicVibeContext `json:"duplicatedContext"`
	ContextTree       *ContextNavigationTree          `json:"contextTree"`
	Artifact          *WrapPreviewArtifact            `json:"artifact"`
}

// CategoricalExtendRequest represents an extend operation request
type CategoricalExtendRequest struct {
	ContextID     string                           `json:"contextId"`
	Context       *streaming.ComonadicVibeContext  `json:"context"`
	Transformation string                          `json:"transformation"`
}

// CategoricalExtendResponse represents the extended context
type CategoricalExtendResponse struct {
	ExtendedContext *streaming.ComonadicVibeContext `json:"extendedContext"`
	TransformTrace  *TransformationTrace            `json:"transformTrace"`
	Artifact        *WrapPreviewArtifact            `json:"artifact"`
}

// TernaryLogicGateRequest represents a ternary logic operation
type TernaryLogicGateRequest struct {
	GateType string                    `json:"gateType"`
	InputA   streaming.TernaryState    `json:"inputA"`
	InputB   streaming.TernaryState    `json:"inputB"`
}

// TernaryLogicGateResponse represents the gate result
type TernaryLogicGateResponse struct {
	Result     streaming.TernaryState  `json:"result"`
	TruthTable map[string]interface{}  `json:"truthTable"`
	Artifact   *WrapPreviewArtifact    `json:"artifact"`
}

// WrapPreviewArtifact represents an artifact for Warp preview
type WrapPreviewArtifact struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Content     interface{}            `json:"content"`
	Interactive bool                   `json:"interactive"`
	Metadata    map[string]interface{} `json:"metadata"`
	Wrapper     string                 `json:"wrapper"`
}

// ContextNavigationTree represents a navigable context tree
type ContextNavigationTree struct {
	Root     *ContextNode   `json:"root"`
	Current  string         `json:"current"`
	Paths    [][]string     `json:"paths"`
}

// ContextNode represents a node in the context tree
type ContextNode struct {
	ID       string         `json:"id"`
	Label    string         `json:"label"`
	Vibe     *models.Vibe   `json:"vibe"`
	Children []*ContextNode `json:"children"`
	Parent   *ContextNode   `json:"parent,omitempty"`
}

// TransformationTrace represents the trace of a transformation
type TransformationTrace struct {
	Steps       []TransformStep `json:"steps"`
	SourceCtx   string          `json:"sourceCtx"`
	TargetCtx   string          `json:"targetCtx"`
	LogicGates  []string        `json:"logicGates"`
}

// TransformStep represents a single transformation step
type TransformStep struct {
	Operation   string                 `json:"operation"`
	Input       interface{}            `json:"input"`
	Output      interface{}            `json:"output"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// RegisterCategoricalTools registers all categorical MCP tools
func RegisterCategoricalTools(mcpServer *server.MCPServer, tools *CategoricalTools) {
	// Register categorical extract tool
	extractTool := mcp.NewTool("categorical_extract", func(t *mcp.Tool) {
		t.Description = "Extract the focused value from a comonadic context"
	})
	
	mcpServer.AddTool(extractTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.CategoricalExtract(ctx, req)
	})

	// Register categorical duplicate tool
	duplicateTool := mcp.NewTool("categorical_duplicate", func(t *mcp.Tool) {
		t.Description = "Create a context-of-contexts via comonadic duplication"
	})
	
	mcpServer.AddTool(duplicateTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.CategoricalDuplicate(ctx, req)
	})

	// Register categorical extend tool
	extendTool := mcp.NewTool("categorical_extend", func(t *mcp.Tool) {
		t.Description = "Apply context-aware transformation via comonadic extension"
	})
	
	mcpServer.AddTool(extendTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.CategoricalExtend(ctx, req)
	})

	// Register ternary logic gate tool  
	ternaryTool := mcp.NewTool("ternary_logic_gate", func(t *mcp.Tool) {
		t.Description = "Execute ternary logic gate operations"
	})
	
	mcpServer.AddTool(ternaryTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.TernaryLogicGate(ctx, req)
	})
}

// CategoricalExtract implements the comonadic extract operation
func (ct *CategoricalTools) CategoricalExtract(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}
	
	contextID, ok := args["contextId"].(string)
	if !ok {
		return mcp.NewToolResultError("contextId is required"), nil
	}

	// Create a sample comonadic context (in real implementation, retrieve from storage)
	vibeCtx := &streaming.ComonadicVibeContext{
		// This would be populated from actual context storage
	}
	
	// Extract the focused vibe
	extractedVibe := vibeCtx.Extract()
	
	// Create WrapPreview artifact
	artifact := &WrapPreviewArtifact{
		Type:        "comonadic_context",
		Title:       fmt.Sprintf("Extracted Vibe from Context %s", contextID),
		Content:     extractedVibe,
		Interactive: true,
		Metadata: map[string]interface{}{
			"operation": "extract",
			"contextId": contextID,
		},
		Wrapper: "ComonadicContextWrapper",
	}
	
	response := &CategoricalExtractResponse{
		ExtractedVibe: extractedVibe,
		ContextID:     contextID,
		Artifact:      artifact,
	}
	
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error marshaling response: %v", err)), nil
	}
	
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// CategoricalDuplicate implements the comonadic duplicate operation
func (ct *CategoricalTools) CategoricalDuplicate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}
	
	contextID, ok := args["contextId"].(string)
	if !ok {
		return mcp.NewToolResultError("contextId is required"), nil
	}

	// Create sample context (retrieve from storage in real implementation)
	vibeCtx := &streaming.ComonadicVibeContext{
		// Populated from context storage
	}
	
	// Perform duplication
	duplicatedCtx := vibeCtx.Duplicate()
	
	// Create navigation tree
	contextTree := &ContextNavigationTree{
		Root: &ContextNode{
			ID:    contextID,
			Label: "Original Context",
		},
		Current: contextID,
		Paths:   [][]string{{contextID}},
	}
	
	// Create WrapPreview artifact
	artifact := &WrapPreviewArtifact{
		Type:        "comonadic_context",
		Title:       fmt.Sprintf("Duplicated Context Tree from %s", contextID),
		Content:     duplicatedCtx,
		Interactive: true,
		Metadata: map[string]interface{}{
			"operation":   "duplicate",
			"sourceCtx":   contextID,
			"treeDepth":   1,
			"navigation":  true,
		},
		Wrapper: "ComonadicContextWrapper",
	}
	
	response := &CategoricalDuplicateResponse{
		DuplicatedContext: duplicatedCtx,
		ContextTree:       contextTree,
		Artifact:          artifact,
	}
	
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error marshaling response: %v", err)), nil
	}
	
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// CategoricalExtend implements the comonadic extend operation
func (ct *CategoricalTools) CategoricalExtend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}
	
	contextID, ok := args["contextId"].(string)
	if !ok {
		return mcp.NewToolResultError("contextId is required"), nil
	}
	
	transformation, ok := args["transformation"].(string)
	if !ok {
		return mcp.NewToolResultError("transformation is required"), nil
	}

	// Create sample context
	vibeCtx := &streaming.ComonadicVibeContext{
		// Populated from context storage
	}
	
	// Define transformation function based on type
	var transformFunc func(*streaming.ComonadicVibeContext) *models.Vibe
	switch transformation {
	case "consensus":
		transformFunc = func(c *streaming.ComonadicVibeContext) *models.Vibe {
			return ct.transformer.TransformWithContext(c.Extract(), c.Neighbors())
		}
	case "amplify":
		transformFunc = func(c *streaming.ComonadicVibeContext) *models.Vibe {
			// Amplification logic
			return c.Extract()
		}
	case "inhibit":
		transformFunc = func(c *streaming.ComonadicVibeContext) *models.Vibe {
			// Inhibition logic
			return c.Extract()
		}
	default:
		return mcp.NewToolResultError("Unknown transformation type"), nil
	}
	
	// Perform extension
	extendedCtx := vibeCtx.Extend(transformFunc)
	
	// Create transformation trace
	trace := &TransformationTrace{
		Steps: []TransformStep{
			{
				Operation: transformation,
				Input:     vibeCtx,
				Output:    extendedCtx,
				Metadata: map[string]interface{}{
					"gateType": transformation,
				},
			},
		},
		SourceCtx:  contextID,
		TargetCtx:  fmt.Sprintf("%s_extended", contextID),
		LogicGates: []string{transformation},
	}
	
	// Create WrapPreview artifact
	artifact := &WrapPreviewArtifact{
		Type:        "comonadic_context",
		Title:       fmt.Sprintf("Extended Context via %s", transformation),
		Content:     extendedCtx,
		Interactive: true,
		Metadata: map[string]interface{}{
			"operation":      "extend",
			"transformation": transformation,
			"sourceCtx":      contextID,
			"trace":          trace,
		},
		Wrapper: "ComonadicContextWrapper",
	}
	
	response := &CategoricalExtendResponse{
		ExtendedContext: extendedCtx,
		TransformTrace:  trace,
		Artifact:        artifact,
	}
	
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error marshaling response: %v", err)), nil
	}
	
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// TernaryLogicGate executes ternary logic operations
func (ct *CategoricalTools) TernaryLogicGate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}
	
	gateType, ok := args["gateType"].(string)
	if !ok {
		return mcp.NewToolResultError("gateType is required"), nil
	}
	
	inputA_float, ok := args["inputA"].(float64)
	if !ok {
		return mcp.NewToolResultError("inputA is required"), nil
	}
	inputA := streaming.TernaryState(int(inputA_float))
	
	inputB_float, ok := args["inputB"].(float64)
	if !ok {
		return mcp.NewToolResultError("inputB is required"), nil
	}
	inputB := streaming.TernaryState(int(inputB_float))
	
	// Create and execute the logic gate
	gate := streaming.NewTernaryLogicGate(gateType)
	result := gate.Apply(inputA, inputB)
	
	// Create truth table for visualization
	truthTable := make(map[string]interface{})
	for _, a := range []streaming.TernaryState{streaming.TernaryNegative, streaming.TernaryNeutral, streaming.TernaryPositive} {
		for _, b := range []streaming.TernaryState{streaming.TernaryNegative, streaming.TernaryNeutral, streaming.TernaryPositive} {
			key := fmt.Sprintf("%d,%d", int(a), int(b))
			truthTable[key] = int(gate.Apply(a, b))
		}
	}
	
	// Create WrapPreview artifact
	artifact := &WrapPreviewArtifact{
		Type:        "ternary_logic_result",
		Title:       fmt.Sprintf("Ternary %s Gate: %d ⊕ %d = %d", gateType, int(inputA), int(inputB), int(result)),
		Content: map[string]interface{}{
			"gate":       gateType,
			"inputs":     []int{int(inputA), int(inputB)},
			"result":     int(result),
			"truthTable": truthTable,
		},
		Interactive: true,
		Metadata: map[string]interface{}{
			"gateType":   gateType,
			"operation":  "ternary_logic",
			"truthTable": truthTable,
		},
		Wrapper: "TernaryLogicWrapper",
	}
	
	response := &TernaryLogicGateResponse{
		Result:     result,
		TruthTable: truthTable,
		Artifact:   artifact,
	}
	
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error marshaling response: %v", err)), nil
	}
	
	return mcp.NewToolResultText(string(responseJSON)), nil
}
