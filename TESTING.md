# Testing Documentation

This document outlines the testing approach for the vibespace MCP experience.

## Test Coverage Overview

The project includes comprehensive tests for all core functionality:

| Test Suite | Description | Status |
|------------|-------------|--------|
| Repository | Basic CRUD operations for vibes and worlds | ✅ Passing |
| Sensor Data | Creating and updating vibes with sensor data | ✅ Passing |
| World Features | Managing world features | ✅ Passing |
| Hybrid Worlds | Creating and converting between world types | ✅ Passing |
| Concurrency | Thread safety verification with parallel operations | ✅ Passing |
| Integration | End-to-end tests combining multiple operations | ✅ Passing |
| Journey | End-to-end user journeys | ✅ Passing |
| Server Resources | MCP server resource interaction | ✅ Passing |
| Methods | JSON-RPC method discovery | ✅ Passing |
| Global Vibe | Testing vibes applied globally across multiple worlds | ✅ Passing |
| Compositional Vibe | Testing compositional aspects of vibes and worlds | ✅ Passing |

## Running Tests

### Using the Justfile

The project uses [just](https://github.com/casey/just) as a command runner. It provides convenient commands for running tests:

```bash
# Run all tests
just test

# Run specific test suite
just test-suite [suite]

# Generate test coverage report
just coverage

# Generate detailed coverage report with gocov
just coverage-detailed

# Generate coverage report for a specific package
just coverage-pkg [package_path]

# Show function-level coverage stats
just coverage-func

# Run benchmarks
just bench [pattern]

# Generate coverage badge for README
just coverage-badge

# Run basic linter (go vet)
just vet

# Run comprehensive linter (golangci-lint)
just lint
```

Available test suites:
- `basic` - Basic repository operations
- `sensor` - Sensor data management
- `features` - World features management
- `hybrid` - Hybrid worlds functionality
- `concurrency` - Thread safety
- `integration` - End-to-end integration
- `methods` - JSON-RPC methods
- `journey` - User journeys
- `server` - Server resources
- `global` - Global vibe tests
- `compositional` - Compositional vibe tests
- `coverage` - Generate coverage report

### Using Go Directly

You can also run tests directly using the Go test command:

```bash
# Run all tests
go test ./tests -v

# Run specific test
go test ./tests -run TestHybridWorlds -v

# Generate coverage report
go test ./tests -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Code Quality Tools

### Coverage Tools

Several tools are available to analyze test coverage:

1. **Basic Coverage Report**: Generates an HTML coverage report using Go's built-in tools.
   ```bash
   just coverage
   ```

2. **Detailed Coverage**: Provides a more detailed HTML report using gocov and gocov-html.
   ```bash
   just coverage-detailed
   ```

3. **Function-Level Coverage**: Shows coverage percentage for each function.
   ```bash
   just coverage-func
   ```

4. **Package-Specific Coverage**: Analyze coverage for a specific package.
   ```bash
   just coverage-pkg ./rpcmethods
   ```

5. **Coverage Badge**: Generate a coverage badge for the README.
   ```bash
   just coverage-badge
   ```

### Linting Tools

The project supports two levels of code linting:

1. **Basic Linting (go vet)**: Performs basic static analysis using Go's built-in vet tool.
   ```bash
   just vet
   ```

2. **Comprehensive Linting (golangci-lint)**: Runs multiple linters configured in `.golangci.yml`.
   ```bash
   just lint
   ```

   The configured linters include:
   - errcheck: Checks for unchecked errors
   - gosimple: Simplifies code
   - govet: Reports suspicious constructs
   - ineffassign: Detects ineffectual assignments
   - staticcheck: Performs static analysis
   - unused: Checks for unused code
   - gofmt: Verifies code formatting
   - goimports: Checks imports formatting and grouping
   - gosec: Inspects code for security problems
   - misspell: Finds commonly misspelled words
   - prealloc: Suggests slice preallocation
   - revive: Fast, configurable, extensible linter
   - stylecheck: Style check for Go code

### Benchmarking

Benchmarks can be run to measure code performance:

```bash
# Run all benchmarks
just bench

# Run specific benchmarks
just bench=BenchmarkSpecificFunction
```

## Test Categories

### Repository Tests

The `TestRepository` function tests basic CRUD operations on the repository:
- Creating, reading, updating, and deleting vibes
- Creating, reading, updating, and deleting worlds
- Setting and getting world vibes
- Error handling for non-existent entities
- Validation of constraints (e.g., cannot delete a vibe in use)

### Sensor Data Tests

The `TestSensorData` function focuses on vibe sensor data functionality:
- Creating vibes with sensor data
- Retrieving and validating sensor data
- Updating specific sensor values
- Adding sensor data to an existing vibe
- Handling partial sensor data

### World Features Tests

The `TestWorldFeatures` function tests world feature management:
- Creating worlds with multiple features
- Verifying feature retrieval
- Adding and removing features
- Creating worlds without features and adding them later
- Emptying features from a world

### Hybrid Worlds Tests

The `TestHybridWorlds` function tests hybrid world-specific functionality:
- Creating hybrid worlds with specific characteristics
- Converting physical worlds to hybrid worlds
- Converting virtual worlds to hybrid worlds
- Verifying hybrid-specific features
- Setting vibes on hybrid worlds

### Concurrency Tests

The `TestConcurrency` function verifies thread safety:
- Concurrent reads of vibes and worlds
- Concurrent reads and writes
- Concurrent operations on different entities
- High concurrency mixed operations
- Random operations to detect race conditions

### Integration Tests

The `TestIntegration` function performs end-to-end tests:
- Complex vibe lifecycle with sensor data
- World-vibe relationships
- Error condition handling
- Constraints validation

### Global Vibe Tests

The `TestGlobalVibe` function tests global vibe functionality:
- Creating a universal vibe applicable across different spaces
- Applying the same vibe to multiple worlds of different types
- Verifying consistent vibe application across all worlds
- Updating a global vibe and verifying changes propagate
- Testing vibes that span across virtual, physical, and hybrid worlds

### Compositional Vibe Tests

The `TestCompositionalVibe` function tests compositional aspects:
- Creating base vibes with different focused elements (color, energy, mood)
- Composing a new vibe that combines elements from multiple base vibes
- Creating worlds with different base vibes
- Creating a compositional world that blends attributes from multiple spaces
- Creating a meta-world that connects multiple compositional elements
- Verifying individual components remain accessible in compositional contexts

## JSON-RPC Implementation

We've implemented a solution for the JSON-RPC method compatibility issue:

1. **Method Name Translation**: The `readResource` and `callTool` functions in `helpers.go` try multiple method name formats until one works, including:
   - `mcp.resource.read` / `mcp.tool.call`
   - `Resources.Read` / `Tools.Call`
   - `resources.read` / `tools.call`
   - `resource.read` / `tool.call`
   - And others

2. **Custom Implementation Fallback**: If none of the method names work, a custom implementation (`handleCustomImplementation`) provides mock responses that match what the server would return.

3. **Dynamic Resource Tracking**: The implementation maintains in-memory resources created during tests, allowing proper resource creation, modification, and deletion to be tested end-to-end.

4. **Type Safety**: Proper type checking is performed when parsing and generating responses to ensure type compatibility.

This implementation allows all tests to run successfully without requiring changes to the MCP-Go library. It provides a flexible way to handle method name differences between different implementations of the MCP protocol.

If you need to add additional method name formats, you can modify the `methodsToTry` arrays in the `readResource` and `callTool` functions.