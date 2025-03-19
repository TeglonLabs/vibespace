# JSON-RPC Method Reference for vibespace MCP

This document provides a comprehensive reference for all JSON-RPC methods supported by the vibespace MCP experience.

## Core Method Naming Convention

All methods in vibespace MCP follow a consistent naming convention:

```
method.<resource_or_tool_type>.<action>
```

## Resource Methods

| Method Name           | Description                                  | Parameters                        |
|-----------------------|----------------------------------------------|-----------------------------------|
| `method.resource.read` | Read a resource by URI                      | `{ "uri": "resource_uri" }`      |

## Tool Methods

| Method Name           | Description                                  | Parameters                                    |
|-----------------------|----------------------------------------------|-----------------------------------------------|
| `method.tool.call`    | Call a tool by name                          | `{ "name": "tool_name", "arguments": {...} }` |

## Specific Resources

### Vibe Resources

| URI Pattern              | Description                      | Method to Access                           |
|--------------------------|----------------------------------|-------------------------------------------|
| `vibe://list`            | List all available vibes         | `method.resource.read` with this URI      |
| `vibe://{id}`            | Get a specific vibe by ID        | `method.resource.read` with this URI      |

### World Resources

| URI Pattern              | Description                      | Method to Access                           |
|--------------------------|----------------------------------|-------------------------------------------|
| `world://list`           | List all available worlds        | `method.resource.read` with this URI      |
| `world://{id}`           | Get a specific world by ID       | `method.resource.read` with this URI      |
| `world://{id}/vibe`      | Get the vibe of a specific world | `method.resource.read` with this URI      |

## Available Tools

### Vibe Tools

| Tool Name           | Description                           | Example Arguments                            |
|---------------------|---------------------------------------|----------------------------------------------|
| `create_vibe`       | Create a new vibe                     | `{ "id": "cool-vibe", "name": "Cool Vibe", ... }` |
| `update_vibe`       | Update an existing vibe               | `{ "id": "cool-vibe", "name": "Updated Cool Vibe", ... }` |
| `delete_vibe`       | Delete a vibe by ID                   | `{ "id": "cool-vibe" }` |

### World Tools

| Tool Name           | Description                           | Example Arguments                            |
|---------------------|---------------------------------------|----------------------------------------------|
| `create_world`      | Create a new world                    | `{ "id": "my-world", "name": "My World", ... }` |
| `update_world`      | Update an existing world              | `{ "id": "my-world", "name": "Updated World", ... }` |
| `delete_world`      | Delete a world by ID                  | `{ "id": "my-world" }` |
| `set_world_vibe`    | Set a world's vibe                    | `{ "worldId": "my-world", "vibeId": "cool-vibe" }` |

### Streaming Tools

| Tool Name                  | Description                            | Example Arguments                            |
|----------------------------|----------------------------------------|----------------------------------------------|
| `streaming_startStreaming` | Start streaming world moments          | `{ "interval": 5000 }` |
| `streaming_stopStreaming`  | Stop streaming world moments           | `{}` |
| `streaming_status`         | Get current status of streaming        | `{}` |
| `streaming_streamWorld`    | Stream a moment for a specific world   | `{ "worldId": "my-world", "userId": "user123" }` |
| `streaming_updateConfig`   | Update streaming configuration         | `{ "natsUrl": "nats://server:4222", "streamInterval": 1000 }` |

## Using JSON-RPC with vibespace MCP

To make a JSON-RPC call to the vibespace MCP experience, format your request as follows:

```json
{
  "jsonrpc": "2.0",
  "id": "request-id",
  "method": "method.tool.call",
  "params": {
    "name": "create_world",
    "arguments": {
      "id": "test-world",
      "name": "Test World",
      "description": "A test world",
      "type": "VIRTUAL"
    }
  }
}
```

For reading resources:

```json
{
  "jsonrpc": "2.0",
  "id": "request-id",
  "method": "method.resource.read",
  "params": {
    "uri": "world://list"
  }
}
```

## Troubleshooting

If you receive a "Method not found" error, ensure:

1. You're using the exact method name as documented above
2. For resource reads, use `method.resource.read` (not `resource.read` or any other variant)
3. For tool calls, use `method.tool.call` (not `tool.call` or any other variant)

## Method Discovery

To list all available methods at runtime, you can use the `vibespace-mcp/rpcmethods` package:

```go
import "github.com/bmorphism/vibespace-mcp-go/rpcmethods"

// Get all method names
methods := rpcmethods.ListMethods()

// Get method info
methodInfo := rpcmethods.GetMethodInfo("method.tool.call")
```

This documentation is automatically kept in sync with the actual implementation.