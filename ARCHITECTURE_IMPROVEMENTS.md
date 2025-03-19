# Architecture Improvements for JSON-RPC Method Handling

This document outlines the architectural improvements made to address the JSON-RPC method name handling issues in the vibespace MCP experience.

## Problem Statement

The original architecture had several issues:

1. **Hidden Method Name Mapping**: No documentation or API to expose method name mapping between resources/tools and JSON-RPC methods
2. **No Method Name Constants**: No exported constants for standard method names
3. **No Method Discovery**: No way to list or validate method names
4. **Unhelpful Error Messages**: Method not found errors didn't provide guidance on what method names to use
5. **Inconsistent Method Naming**: Different clients might use different conventions (e.g., `resource.read` vs `method.resource.read`)

These issues resulted in frustrating developer experiences, with tests showing many failed attempts to use different method name formats.

## Solution Architecture

We implemented a comprehensive solution with several components:

### 1. Documentation & Constants (`RPC_METHODS.md` and `rpcmethods` package)

- Created clear documentation listing all supported method names
- Defined constants for standard method names (`MethodResourceRead`, `MethodToolCall`)
- Provided examples showing the correct JSON-RPC format

### 2. Method Discovery & Utilities (`rpcmethods` package)

- `ListMethods()`: Returns all registered method names
- `GetMethodInfo()`: Gets detailed information about a specific method
- `IsValidMethod()`: Validates if a method name is registered
- `FindMethod()`: Provides fuzzy matching for method names
- `GetMethodSuggestions()`: Gives helpful suggestions for invalid method names

### 3. Request Formatting Helpers

- `FormatResourceRequest()`: Creates properly formatted resource read requests
- `FormatToolRequest()`: Creates properly formatted tool call requests

### 4. Improved Method Handling (`MCPMethodWrapper`)

- Created a wrapper around the MCP server
- Intercepts JSON-RPC messages to normalize method names
- Provides helpful error messages when invalid methods are used
- Suggests the correct method name based on the attempted method

### 5. Method Name Normalization

- Handles common variants (`resource.read` → `method.resource.read`)
- Converts slashes to dots (`mcp/resource/read` → `method.resource.read`)
- Special case handling for resource reads and tool calls

## Benefits of New Architecture

1. **Clear Documentation**: Developers know exactly which method names to use
2. **Method Name Constants**: No guesswork when writing code
3. **Helpful Error Messages**: When a wrong method is used, users get suggestions
4. **Automatic Normalization**: Common method name variants work automatically
5. **Request Helpers**: Simplified request construction

## Architectural Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         vibespace MCP Experience                         │
│                                                                         │
│  ┌───────────────────┐     ┌──────────────────┐     ┌───────────────┐   │
│  │                   │     │                  │     │               │   │
│  │  RPC_METHODS.md   │     │  MCPMethodWrapper│     │  rpcmethods   │   │
│  │  Documentation    │◄────┤  Server Wrapper  │◄────┤  Package      │   │
│  │                   │     │                  │     │               │   │
│  └───────────────────┘     └──────────────────┘     └───────────────┘   │
│          ▲                         ▲                       ▲             │
│          │                         │                       │             │
│          ▼                         ▼                       ▼             │
│  ┌─────────────────┐       ┌──────────────────┐      ┌────────────────┐ │
│  │                 │       │                  │      │                │ │
│  │ JSON-RPC APIs   │       │  Method          │      │ Helper         │ │
│  │ & Format        │       │  Normalization   │      │ Functions      │ │
│  │                 │       │                  │      │                │ │
│  └─────────────────┘       └──────────────────┘      └────────────────┘ │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Usage Example

Client code using our improved method handling:

```go
import "github.com/bmorphism/vibespace-mcp-go/rpcmethods"

// Get list of supported methods
methods := rpcmethods.ListMethods()

// Format proper JSON-RPC request to read a world list
request := rpcmethods.FormatResourceRequest("world://list", "request-1")

// Make the request to the server...
// If an error occurs, the server will suggest the correct method name
```

The architectural improvements make it much easier for developers to integrate with the vibespace MCP experience and avoid common JSON-RPC method naming issues.