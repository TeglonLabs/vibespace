# Warp Preview Categorical Artifacts Integration

This document describes the integration between your vibespace MCP server and Warp terminal's preview system for categorical universe artifacts.

## Overview

The integration provides live interactive preview of:
- **Comonadic Contexts**: Navigate context trees with extract/duplicate/extend operations
- **Ternary Logic Gates**: Interactive truth tables and gate visualizations  
- **Transductive Chains**: Step-by-step verification traces
- **Categorical Diagrams**: Visual representation of morphisms and functors

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌───────────────────┐
│   Warp Agent    │◄──►│  MCP Protocol    │◄──►│ Categorical Tools │
│     Mode        │    │   JSON-RPC       │    │   ComonadicCtx    │
└─────────────────┘    └──────────────────┘    │   TernaryGates    │
                                                │   Transductions   │
                                                └───────────────────┘
           ▲                                              ▲
           │                                              │
    ┌──────▼──────┐                                ┌─────▼─────┐
    │ WrapPreview │                                │ Streaming │
    │  Artifacts  │                                │  Engine   │
    └─────────────┘                                └───────────┘
```

## MCP Tools Available

### Comonadic Operations

#### `categorical_extract`
Extracts the focused value from a comonadic context.

```bash
@categorical_extract contextId="ctx-001"
```

**Response:**
```json
{
  "extractedVibe": { "id": "focus", "energy": 0.6 },
  "contextId": "ctx-001", 
  "artifact": {
    "type": "comonadic_context",
    "wrapper": "ComonadicContextWrapper",
    "interactive": true
  }
}
```

#### `categorical_duplicate`
Creates a context-of-contexts via comonadic duplication.

```bash
@categorical_duplicate contextId="ctx-001"
```

**Features:**
- Interactive context tree navigation
- Visual hierarchy of nested contexts
- Hot-swappable focus points

#### `categorical_extend`
Applies context-aware transformations via comonadic extension.

```bash
@categorical_extend contextId="ctx-001" transformation="consensus"
```

**Transformation Types:**
- `consensus`: Harmonic convergence across contexts
- `amplify`: Strengthen agreement patterns
- `inhibit`: Dampen chaotic oscillations

### Ternary Logic Operations

#### `ternary_logic_gate`
Executes ternary logic operations with interactive visualization.

```bash
@ternary_logic_gate gateType="consensus" inputA=1 inputB=1
```

**Gate Types:**
- `consensus`: Agrees when inputs agree
- `amplify`: Strengthens agreement, weakens disagreement  
- `inhibit`: Complex inhibition dynamics

**Interactive Features:**
- Live truth table exploration
- State transition visualization
- Input/output tracing

## WrapPreview Artifact Types

### ComonadicContextWrapper
Renders comonadic contexts with:
- **Extract Visualization**: Focused value highlighting
- **Duplicate Tree**: Expandable context hierarchy
- **Extend Preview**: Real-time transformation preview

### TernaryLogicWrapper  
Displays ternary logic operations with:
- **Truth Table Display**: Interactive 3×3 grid
- **State Transitions**: Animated state changes
- **Interactive Gates**: Drag-and-drop gate construction

### TransductiveChainWrapper
Shows verification chains with:
- **Verification Steps**: Step-by-step proof exploration
- **Proof Tree**: Branching verification paths
- **Categorical Arrows**: Morphism visualization

## Setup Instructions

1. **Install Warp Preview**
   ```bash
   # Download from https://www.warp.dev/download-preview
   ```

2. **Run Setup Script**
   ```bash
   ./setup-warp-preview.sh
   ```

3. **Configure in Warp**
   - Open Warp Drive → Personal → MCP Servers
   - Add CLI Server: `vibespace-categorical`
   - Command: Use generated startup script
   - Environment: `PREVIEW_ARTIFACTS=true`

4. **Test Integration**
   ```bash
   @categorical_extract contextId="test"
   @ternary_logic_gate gateType="consensus" inputA=1 inputB=0
   ```

## Usage Examples

### Exploring Context Hierarchies

```bash
# Extract current focus
@categorical_extract contextId="workspace"

# Duplicate for exploration
@categorical_duplicate contextId="workspace"  

# Navigate in preview panel
# Click nodes to change focus
# Expand/collapse context branches

# Apply consensus transformation
@categorical_extend contextId="workspace" transformation="consensus"
```

### Ternary Logic Exploration

```bash
# Test all consensus combinations
@ternary_logic_gate gateType="consensus" inputA=-1 inputB=-1  # → -1
@ternary_logic_gate gateType="consensus" inputA=0 inputB=0    # → 0
@ternary_logic_gate gateType="consensus" inputA=1 inputB=1    # → 1
@ternary_logic_gate gateType="consensus" inputA=1 inputB=-1   # → 0

# Interactive truth table appears in preview
# Click cells to see detailed computations
# Drag inputs to test combinations
```

### Verification Chains

```bash
# Create transductive verification chain
@transductive_chain_verify 
  sourceContext="workspace" 
  targetContext="flow-state"
  steps=["focus", "eliminate-distractions", "enter-flow"]

# Preview shows:
# - Verification steps with proofs
# - Failed verification highlights  
# - Alternative path suggestions
```

## Configuration

### warp-mcp-config.json

```json
{
  "mcpServers": {
    "vibespace-categorical": {
      "command": "./run-mcp.sh",
      "startOnLaunch": true,
      "env": {
        "VIBESPACE_MODE": "categorical",
        "PREVIEW_ARTIFACTS": "true", 
        "COMONADIC_CONTEXT": "enabled"
      }
    }
  },
  "warp_preview": {
    "artifact_rendering": {
      "comonadic_contexts": {
        "extract_visualization": true,
        "duplicate_tree": true,
        "extend_preview": true
      },
      "ternary_logic": {
        "truth_table_display": true,
        "state_transitions": true,
        "interactive_gates": true
      }
    },
    "live_preview": {
      "auto_refresh": true,
      "debounce_ms": 300,
      "hot_reload": true
    }
  }
}
```

### Environment Variables

- `VIBESPACE_MODE=categorical`: Enable categorical features
- `PREVIEW_ARTIFACTS=true`: Generate WrapPreview artifacts
- `COMONADIC_CONTEXT=enabled`: Enable comonadic operations
- `DEBUG_ARTIFACTS=true`: Additional debug information

## Troubleshooting

### Common Issues

1. **Server Won't Start**
   ```bash
   # Check Go installation
   go version
   
   # Rebuild server
   go build -o vibespace-mcp-server cmd/server/main.go
   
   # Test directly
   ./vibespace-mcp-server
   ```

2. **Warp Can't Find Server**
   ```bash
   # Check startup script path
   ls -la "$HOME/Library/Application Support/dev.warp.Warp-Stable/mcp/"
   
   # Verify permissions
   chmod +x start-vibespace.sh
   ```

3. **Artifacts Not Rendering**
   - Ensure `PREVIEW_ARTIFACTS=true` in environment
   - Check Warp Preview version is latest
   - Verify MCP server logs in Warp

### Debug Commands

```bash
# Test MCP connection
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","id":1}'

# Check artifact generation
@categorical_extract contextId="debug" --verbose

# Inspect server logs
tail -f "$HOME/Library/Application Support/dev.warp.Warp-Stable/mcp/logs/vibespace-categorical.log"
```

## Advanced Features

### Custom Artifact Wrappers

Create custom wrappers for domain-specific visualizations:

```go
type CustomWrapper struct {
    Type        string      `json:"type"`
    Component   string      `json:"component"`
    Props       interface{} `json:"props"`
    Interactive bool        `json:"interactive"`
}
```

### Live Context Streaming

Enable real-time context updates:

```bash
@streaming_startStreaming interval=1000
# Artifacts automatically update with streaming data
```

### Collaborative Exploration

Share context explorations across team members:

```bash
@categorical_extend contextId="shared-workspace" transformation="consensus"
# All team members see live updates in their Warp Preview
```

## Integration with Existing Tools

### NATS Streaming
Categorical artifacts integrate with your existing NATS streaming infrastructure for real-time updates.

### Comonadic Transformations  
Leverages your sophisticated `ComonadicVibeContext` system for authentic category theory operations.

### Ternary Logic Engine
Uses your existing ternary logic gates for computation and visualization.

## Future Enhancements

- [ ] Visual diagram editor for categorical morphisms
- [ ] Collaborative real-time context editing
- [ ] Export artifacts to computational notebooks
- [ ] Integration with proof assistants
- [ ] Custom visualization plugins
- [ ] Advanced animation for transformations

---

Your categorical universe artifacts are now fully integrated with Warp Preview for unprecedented interactive exploration! 🔮
