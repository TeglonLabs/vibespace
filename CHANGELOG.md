# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - Package Updates & Compatibility Fixes

### Added
- Comprehensive package updates to latest stable versions
- Enhanced nil pointer safety in streaming service
- Better error handling and logging in streaming components
- Type safety improvements for MCP library compatibility

### Changed
- **BREAKING**: Updated Go version from 1.23.0 to 1.23.5
- **BREAKING**: Updated MCP library to v0.32.0 with new RequestId handling
- Updated key dependencies:
  - `github.com/klauspost/compress` v1.17.2 → v1.18.0
  - `github.com/nats-io/nats.go` v1.33.1 → v1.43.0
  - `github.com/nats-io/nkeys` v0.4.7 → v0.4.11
  - `github.com/spf13/cast` v1.7.1 → v1.9.2
  - `golang.org/x/crypto` v0.35.0 → v0.39.0
  - `golang.org/x/text` v0.22.0 → v0.26.0

### Fixed
- **CRITICAL**: Fixed repository interface compatibility issues in main.go
- **CRITICAL**: Fixed RequestId type handling for MCP library v0.32.0
- **CRITICAL**: Fixed argument type assertions from `map[string]interface{}` to `any`
- Fixed nil pointer dereference in streaming service `streamMoments` function
- Fixed all test helpers to work with new MCP library type system
- Improved type safety across all MCP-related code

### Technical Details

#### MCP Library Compatibility (v0.32.0)
- `RequestId` is now a struct requiring `mcp.NewRequestId(value)` constructor
- Tool call arguments are now of type `any` instead of `map[string]interface{}`
- Resource content handling updated for new type system

#### Repository Interface Updates
- Added proper `VibeWorldRepository` interface implementation
- Fixed repository type passing in main server initialization
- Enhanced interface compliance checking

#### Test Infrastructure
- All test helpers updated for new MCP library compatibility
- Enhanced error handling in test mock implementations
- Improved type assertion safety throughout test suite

### Performance & Reliability
- Enhanced memory safety with nil pointer checks
- Improved error propagation in streaming components
- Better resource cleanup in test infrastructure
- Reduced potential race conditions in concurrent operations

### Documentation
- Updated README.md with current dependency versions
- Enhanced RPC_METHODS.md with latest method signatures
- Improved TESTING.md with current test status
- Added comprehensive CHANGELOG.md

### Build & Test Status
- ✅ **Build**: All packages compile successfully
- ✅ **Core Tests**: All primary functionality tests pass
- ✅ **Models Tests**: Balanced ternary and binary data tests pass
- ✅ **RPC Methods Tests**: Method discovery and formatting tests pass
- ✅ **Repository Tests**: CRUD operations and concurrency tests pass
- ⚠️ **Streaming Tests**: Most tests pass; some concurrency tests have known issues

### Upgrade Path
1. Update Go to version 1.23.5 or later
2. Run `go mod tidy` to clean up dependencies
3. Update any custom MCP implementations to use `mcp.NewRequestId()`
4. Update argument handling from `map[string]interface{}` to type assertions with `any`
5. Test all MCP-related functionality thoroughly

### Known Issues
- Some streaming concurrency tests may experience race conditions (non-breaking)
- Legacy test helpers in streaming module need refactoring (tracked)

### Compatibility
- **Go**: Requires Go 1.23.5+
- **MCP**: Compatible with MCP library v0.32.0
- **NATS**: Compatible with NATS v1.43.0+
- **Dependencies**: All dependencies updated to latest stable versions

---

## Previous Versions

### [v1.0.0] - Initial Release
- Basic MCP server implementation
- Vibe and World management
- NATS streaming integration
- Balanced ternary support
- Binary data encoding
- Comprehensive testing suite
