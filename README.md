# vibespace MCP Experience

An MCP (Model Context Protocol) experience implementation for managing vibes and worlds with real-time NATS streaming.

## What Makes This Special

- **🧮 Categorical Universe Artifacts**: First-class support for comonadic contexts and ternary logic
- **🎯 Interactive Previews**: WrapPreview-compatible artifacts for live exploration
- **⚡ Real-time Streaming**: NATS-powered world moment streaming
- **🔧 Production Ready**: Comprehensive testing, CI/CD, and observability
- **🌊 Transductive Chains**: Context-aware transformations with verification traces

## Overview

vibespace MCP Experience is a Go implementation of a Model Context Protocol experience that provides resources and tools for managing "vibes" (emotional atmospheres) and "worlds" (physical or virtual spaces). The server uses custom URI schemes (`vibe://` and `world://`) to access resources and supports JSON-RPC methods for communication.

## Features

- In-memory repository for storing vibes and worlds
- Custom URI schemes for resource access
- Full CRUD operations for vibes and worlds
- Support for sensor data (temperature, humidity, light, sound, movement)
- Support for world features
- Support for different world types (physical, virtual, hybrid)
- NATS integration for streaming world moments in real-time
- Stream ID namespacing for multi-tenant NATS servers
- Multiplayer support with user attribution and access control
- Connection quality monitoring and UI status indicators
- Rate limiting for stream publishing
- Binary data support with multiple encoding formats (binary, base64, hex)
- Balanced ternary numeral system support with efficient encoding/decoding

## Components

- **Models**: Data structures for Vibe, World, and WorldMoment, including SensorData
- **Repository**: In-memory storage with thread-safe CRUD operations
- **Server**: MCP protocol implementation with JSON-RPC methods
- **Streaming**: NATS integration for real-time streaming of world moments with Stream ID support

## 🚀 Quick Start

### Option 1: Download Pre-built Binary (Recommended)

```bash
# Download latest release
wget https://github.com/bmorphism/vibespace-mcp-go/releases/latest/download/vibespace-mcp-linux-amd64.tar.gz
tar -xzf vibespace-mcp-linux-amd64.tar.gz
chmod +x vibespace-mcp-linux-amd64

# Run the server
./vibespace-mcp-linux-amd64
```

### Option 2: Using Docker

```bash
docker run -p 8080:8080 ghcr.io/bmorphism/vibespace-mcp-go:latest
```

### Option 3: Build from Source

```bash
go get github.com/bmorphism/vibespace-mcp-go
go build -o vibespace-mcp ./cmd/server
./vibespace-mcp
```

## 🎯 MCP Client Integration

### Warp Terminal (Recommended)

1. **Install Warp Preview**: [Download here](https://www.warp.dev/download-preview)
2. **Add MCP Server**: 
   - Open Warp Drive → Personal → MCP Servers
   - Click "+ Add" → CLI Server
   - Name: `vibespace-categorical`
   - Command: `go run cmd/server/main.go`
   - Environment:
     ```
     VIBESPACE_MODE=categorical
     PREVIEW_ARTIFACTS=true
     COMONADIC_CONTEXT=enabled
     ```
3. **Test Integration**:
   ```bash
   @categorical_extract contextId="test-workspace"
   @ternary_logic_gate gateType="consensus" inputA=1 inputB=1
   ```

### Claude Desktop

Add to your `mcp.json`:

```json
{
  "mcpServers": {
    "vibespace-categorical": {
      "command": "./vibespace-mcp-server",
      "env": {
        "VIBESPACE_MODE": "categorical",
        "PREVIEW_ARTIFACTS": "true"
      }
    }
  }
}
```

### VS Code with Continue

```yaml
# config.yaml
mcp:
  vibespace:
    command: ["./vibespace-mcp-server"]
    env:
      VIBESPACE_MODE: categorical
      PREVIEW_ARTIFACTS: "true"
```

### Using MCP Inspector (Development)

```bash
npx @modelcontextprotocol/inspector ./vibespace-mcp-server
# Opens browser at http://localhost:5173
```

## Running the Server

To run the server:

```bash
just run
```

Or build and run manually:

```bash
just build
./bin/vibespace-mcp
```

### NATS Subscriber Example

The repository includes an example NATS subscriber to listen for world moments:

```bash
cd examples
./run_nats_subscriber.sh [nats-url] [stream-id] [user-id]
```

Parameters:
- `nats-url`: NATS server URL (defaults to `nats://nonlocal.info:4222`)
- `stream-id`: Stream ID for subject namespacing (defaults to `ies`)
- `user-id`: Optional user ID for user-specific streams

## Testing

The server includes comprehensive tests for all functionality:

```bash
# Run all tests
just test

# Run a specific test suite
just test-suite hybrid

# Generate test coverage report
just coverage
```

See [TESTING.md](./TESTING.md) for detailed information about the testing approach and available test suites.

### Test Coverage

- **Repository Tests**: Basic CRUD operations for vibes and worlds
- **Sensor Data Tests**: Creating and updating vibes with sensor data
- **World Features Tests**: Managing world features (adding, removing, etc.)
- **Hybrid Worlds Tests**: Creating and converting between world types
- **Concurrency Tests**: Thread safety verification with parallel operations
- **Integration Tests**: End-to-end tests combining multiple operations
- **Method Tests**: JSON-RPC method discovery and compatibility

## 🧮 Categorical Tools (New!)

The server now includes specialized tools for categorical universe exploration:

### Comonadic Operations
- **`categorical_extract`**: Extract focused values from comonadic contexts
- **`categorical_duplicate`**: Create navigable context-of-contexts
- **`categorical_extend`**: Apply context-aware transformations

### Ternary Logic Gates
- **`ternary_logic_gate`**: Execute consensus/amplify/inhibit operations
- Interactive truth tables and state visualizations
- Real-time WrapPreview artifact generation

### Example Usage

```bash
# Extract a focused vibe from context
@categorical_extract contextId="workspace-flow"

# Create navigable context tree
@categorical_duplicate contextId="workspace-flow"

# Apply consensus transformation
@categorical_extend contextId="workspace-flow" transformation="consensus"

# Execute ternary logic gate
@ternary_logic_gate gateType="consensus" inputA=1 inputB=1
```

## MCP Protocol

The server implements the Model Context Protocol providing:

- **Resources**: `vibe://list`, `vibe://{id}`, `world://list`, `world://{id}`, `world://{id}/vibe`
- **Tools**: 
  - **Vibe Tools**: `create_vibe`, `update_vibe`, `delete_vibe`
  - **World Tools**: `create_world`, `update_world`, `delete_world`, `set_world_vibe`
  - **Streaming Tools**: `streaming_startStreaming`, `streaming_stopStreaming`, `streaming_status`, `streaming_streamWorld`, `streaming_updateConfig`
  - **Categorical Tools**: `categorical_extract`, `categorical_duplicate`, `categorical_extend`, `ternary_logic_gate`

For more details on the streaming capabilities, see [STREAMING.md](./STREAMING.md).

## JSON-RPC Method Documentation

The vibespace MCP experience provides clear documentation and helper utilities for JSON-RPC method names, preventing common issues with method name mismatches. See [RPC_METHODS.md](./RPC_METHODS.md) for complete method reference.

### Method Name Constants

Use the standard method name constants from the `rpcmethods` package in your code:

```go
import "github.com/bmorphism/vibespace-mcp-go/rpcmethods"

// Standard method names
const ResourceReadMethod = rpcmethods.MethodResourceRead  // "method.resource.read"
const ToolCallMethod = rpcmethods.MethodToolCall          // "method.tool.call"
```

### Helper Functions

The `rpcmethods` package provides utilities to simplify JSON-RPC interactions:

```go
// For formatting proper JSON-RPC requests
resourceRequest := rpcmethods.FormatResourceRequest("world://list", "request-id")
toolRequest := rpcmethods.FormatToolRequest("create_world", worldArgs, "request-id")

// For discovering method names
methods := rpcmethods.ListMethods()
info := rpcmethods.GetMethodInfo("method.resource.read")
suggestions := rpcmethods.GetMethodSuggestions("resource.read") // Gives helpful suggestions
```

### Improved Error Messages

If you use an incorrect method name, the server will suggest the correct method:

```
Error: Method 'resource.read' not found. Did you mean: method.resource.read?
```

## JSON-RPC Compatibility Note

Some tests are currently skipped due to JSON-RPC method compatibility issues with the current version of the MCP-Go library. This affects the journey and server integration tests, but our method name handling improvements should resolve most issues for clients.

## User Example: Team Topos Shared Experience

Here's how the Topos team can leverage the vibespace MCP experience to create a cohesive team environment while working in different tools:

### Scenario: Team Topos Collaborative Development

Members of **Team Topos** are distributed across tools: **Alice** (using Claude Desktop), **Bob** (using Goose MCP client), and **Carol** (using VS Code with Claude). They want to maintain a sense of "vibing together" while developing categorical network models across different interfaces.

**Setup:**

1. A shared vibespace MCP experience runs on a dedicated Topos server
2. A NATS server on `nats://nonlocal.info:4222` handles real-time streaming
3. Each team member connects with a unique user ID but shares the team-topos stream ID

**Step 1: Creating the Team Topos World**

Alice initializes the Team Topos world through Claude Desktop:

```javascript
// Alice creates the Team Topos world
const response = await mcp.tools.create_world({
  name: "Team Topos Workspace",
  description: "A shared categorical space for infinity-topos development",
  type: "hybrid",
  features: ["categorical-thinking", "research-sharing", "presence-indication", "distributed-computation"],
  sharingSettings: {
    isPublic: false,
    allowedUsers: ["alice@topos", "bob@topos", "carol@topos"],
    contextLevel: "full"
  }
});

// Alice sets the initial team vibe
await mcp.tools.set_world_vibe({
  worldId: response.id,
  vibeId: "categorical-flow-state"
});

console.log(`Join our Team Topos workspace with ID: ${response.id}`);
```

**Step 2: Joining the Team Topos Workspace**

Bob connects through the Goose MCP client:

```python
# Bob connects to Team Topos with his user ID
world_id = "topos_world_123"  # ID shared by Alice
client.streaming_startStreaming(
    url="nats://nonlocal.info:4222",
    stream_id="team-topos",
    world_id=world_id,
    user_id="bob@topos"
)

# Bob indicates he's working on sheaf theory implementation
client.update_world_feature(
    world_id=world_id,
    feature_name="research-focus",
    feature_value="Sheaf cohomology for network models"
)

# Bob shares his computational resources
client.update_world_feature(
    world_id=world_id,
    feature_name="shared-computation",
    feature_value={"available_cores": 16, "gpu_memory": "24GB"}
)
```

Carol joins through VS Code with Claude:

```typescript
// Carol connects to Team Topos from VS Code
const connection = await mcpClient.connect({
  worldId: "topos_world_123",
  userId: "carol@topos",
  streamId: "team-topos",
  statusCallback: (status) => {
    vscode.window.setStatusBarMessage(`⊤ Team Topos: ${status.connectedUsers.length} members active`);
  }
});

// Carol shares her current research direction
await mcpClient.tools.update_world_feature({
  worldId: "topos_world_123",
  featureName: "active-models",
  featureValue: ["dynamic-markov-blankets", "hypergraph-topologies"]
});

// Carol updates the shared mathematical context
await mcpClient.tools.update_vibe({
  worldId: "topos_world_123",
  attribute: "mathematical-coherence",
  value: 92  // Scale 0-100 indicating conceptual clarity
});
```

**Step 3: Experiencing Shared Topos Development**

As the team works:

1. **Alice** sees real-time visualizations in Claude Desktop showing which categorical constructs Bob and Carol are exploring, with directed graphs of collaboration paths
   
2. **Bob** receives ambient notifications in Goose when Alice or Carol make mathematical breakthroughs, with the interface subtly reflecting the team's collective understanding through topological visualizations
   
3. **Carol** sees VS Code extensions automatically activating based on the team's current focus areas, with sidebar indicators showing which mathematical domains are being actively explored
   
4. All three experience synchronized development environments with:
   - Shared LaTeX symbol suggestions based on team usage patterns
   - Collective computation resources allocated based on current needs
   - Ambient music that subtly shifts to match the team's cognitive state
   - Visual themes that align with the categorical structures being explored

5. The vibespace MCP experience enables:
   - Automatic documentation generation from the team's shared mathematical context
   - Synchronized whiteboards that appear when conceptual misalignments are detected
   - Distributed computation of complex categorical models across team members' machines

**Step 4: Leveraging Binary and Balanced Ternary Data**

As the team works on advanced mathematical models, they use the binary and balanced ternary features:

```typescript
// Carol shares balanced ternary data representing a categorical network structure
const ternaryEncoding = "10T01T001T10T01"; // Balanced ternary representing network topology
await mcpClient.tools.streaming_streamWorld({
  worldId: "topos_world_123",
  userId: "carol@topos",
  customData: {
    action: "share_mathematical_model",
    modelName: "HypergraphCohomology"
  },
  // Add balanced ternary data representing the mathematical structure
  balancedTernaryData: ternaryEncoding
});
```

Bob processes the mathematical structure using Python:

```python
# Bob subscribes to the moment and extracts balanced ternary data
def on_world_moment(message):
    moment = json.loads(message.data.decode())
    if moment.get("binaryData") and moment["binaryData"]["format"] == "application/balanced-ternary":
        # Extract balanced ternary representation
        binary_data = base64.b64decode(moment["binaryData"]["data"])
        
        # Calculate how many trits were encoded (each 4 trits use 1 byte)
        num_trits = len(binary_data) * 4
        
        # Convert back to ternary representation
        ternary_data = bytes_to_ternary(binary_data, num_trits)
        
        # Process mathematical structure
        if "HypergraphCohomology" in moment.get("customData", {}).get("modelName", ""):
            process_categorical_network(ternary_data)
    
# Convert binary to balanced ternary (Python implementation)
def bytes_to_ternary(bytes_data, num_trits):
    result = []
    for i in range(min(num_trits, len(bytes_data) * 4)):
        byte_idx = i // 4
        trit_pos = i % 4
        
        # Extract 2 bits representing this trit
        bits = (bytes_data[byte_idx] >> (trit_pos * 2)) & 0b11
        
        # Map bits to trits: 00->-1, 01->0, 10->1
        if bits == 0:
            result.append(-1)  # 'T'
        elif bits == 1:
            result.append(0)   # '0'
        elif bits == 2:
            result.append(1)   # '1'
        else:
            result.append(0)   # Invalid pattern, default to 0
    
    return result
```

Alice uses the data to visualize the categorical network:

```javascript
// Alice subscribes to world moments and extracts mathematical models
natsClient.subscribe(`${streamId}.world.moment.${worldId}`, (msg) => {
  const moment = JSON.parse(msg.data);
  
  // Check for balanced ternary data
  if (moment.binaryData && moment.binaryData.format === "application/balanced-ternary") {
    // Decode the binary data
    const binaryData = atob(moment.binaryData.data); // Base64 decode
    
    // Create a Uint8Array from the binary string
    const bytes = new Uint8Array(binaryData.length);
    for (let i = 0; i < binaryData.length; i++) {
      bytes[i] = binaryData.charCodeAt(i);
    }
    
    // Convert to balanced ternary
    const ternaryData = bytesToTernary(bytes, bytes.length * 4);
    
    // Visualize the mathematical structure
    visualizeCategoricalNetwork(ternaryData);
  }
});

// Convert binary to balanced ternary (JavaScript implementation)
function bytesToTernary(bytes, numTrits) {
  const result = [];
  
  for (let i = 0; i < Math.min(numTrits, bytes.length * 4); i++) {
    const byteIdx = Math.floor(i / 4);
    const tritPos = i % 4;
    
    // Extract the trit value (2 bits)
    const bits = (bytes[byteIdx] >> (tritPos * 2)) & 0b11;
    
    // Map bit patterns to trits
    if (bits === 0) {
      result.push(-1); // 'T'
    } else if (bits === 1) {
      result.push(0);  // '0'
    } else if (bits === 2) {
      result.push(1);  // '1'
    } else {
      result.push(0);  // Invalid pattern, default to 0
    }
  }
  
  return result;
}
```

**Benefits for Team Topos:**

- **Unified mathematical context**: The team maintains conceptual alignment across different tools
- **Ambient categorical awareness**: Non-intrusive indicators of theoretical explorations
- **Collective intelligence**: Shared vibes that enhance group mathematical intuition
- **Tool interoperability**: Different interfaces participate in the same categorical workspace
- **Resource optimization**: Computational tasks distributed based on available resources
- **Cognitive synchronization**: Team members' thought patterns become more aligned over time
- **Efficient mathematical representation**: Balanced ternary provides compact encoding of mathematical structures
- **Cross-language compatibility**: Binary data can be processed consistently across different programming languages

This example demonstrates how vibespace MCP experience enables Team Topos to collaboratively explore mathematical structures across different tools and interfaces, creating a shared cognitive space that transcends individual environments.