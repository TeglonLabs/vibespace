# vibespace MCP Experience Justfile

# Default recipe
default:
    @just --list

# Run the server
run:
    go run cmd/server/main.go

# Build the server
build:
    go build -o bin/vibespace-mcp cmd/server/main.go

# Run all tests
test:
    @echo -e "\033[0;34m====================================\033[0m"
    @echo -e "\033[0;34m   vibespace MCP Experience Tests  \033[0m"
    @echo -e "\033[0;34m====================================\033[0m"
    @echo -e "\n\033[1;33mRunning all tests...\033[0m"
    go test -v ./tests
    @echo -e "\n\033[0;34m====================================\033[0m"
    @echo -e "\033[0;32mTesting complete!\033[0m"

# Run specific test suite
test-suite suite="":
    #!/usr/bin/env bash
    set -e
    
    # Colors for output
    GREEN='\033[0;32m'
    BLUE='\033[0;34m'
    YELLOW='\033[1;33m'
    RED='\033[0;31m'
    NC='\033[0m' # No Color
    
    echo -e "${BLUE}====================================${NC}"
    echo -e "${BLUE}   vibespace MCP Experience Tests  ${NC}"
    echo -e "${BLUE}====================================${NC}"
    
    # Function to run tests with a specific pattern
    run_test_pattern() {
        pattern=$1
        description=$2
        options=$3
        
        echo -e "\n${YELLOW}Running $description tests...${NC}"
        go test ./tests -run "$pattern" $options
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ $description tests passed${NC}"
        else
            echo -e "${RED}✗ $description tests failed${NC}"
            exit 1
        fi
    }
    
    # Run specific test suite based on argument
    case "{{suite}}" in
        "basic")
            run_test_pattern "TestRepository" "basic repository" -v
            ;;
        "sensor")
            run_test_pattern "TestSensorData" "sensor data" -v
            ;;
        "features")
            run_test_pattern "TestWorldFeatures" "world features" -v
            ;;
        "hybrid")
            run_test_pattern "TestHybridWorlds" "hybrid worlds" -v
            ;;
        "concurrency")
            run_test_pattern "TestConcurrency" "concurrency" -v
            ;;
        "integration")
            run_test_pattern "TestIntegration" "integration" -v
            ;;
        "methods")
            run_test_pattern "TestMethods" "JSON-RPC methods" -v
            ;;
        "journey")
            run_test_pattern "TestJourney" "journey" -v
            ;;
        "server")
            run_test_pattern "TestResource" "server resource" -v
            ;;
        "global")
            run_test_pattern "TestGlobalVibe" "global vibe" -v
            ;;
        "compositional")
            run_test_pattern "TestCompositionalVibe" "compositional vibe" -v
            ;;
        "coverage")
            echo -e "\n${YELLOW}Running all tests with coverage...${NC}"
            go test ./rpcmethods/... ./streaming/... ./tests/... -coverprofile=coverage.out
            go tool cover -html=coverage.out -o coverage.html
            echo -e "${GREEN}✓ Coverage report generated in coverage.html${NC}"
            # Try to open the coverage report if possible
            if command -v open >/dev/null 2>&1; then
                open coverage.html
            else
                echo "Open coverage.html in your browser to view the report"
            fi
            ;;
        "")
            # No argument provided, show available test suites
            echo -e "\n${YELLOW}Please specify a test suite:${NC}"
            echo "Available: basic, sensor, features, hybrid, concurrency, integration, methods, journey, server, global, compositional, coverage"
            echo "Example: just test-suite hybrid"
            exit 1
            ;;
        *)
            echo -e "${RED}Unknown test suite: {{suite}}${NC}"
            echo "Available: basic, sensor, features, hybrid, concurrency, integration, methods, journey, server, global, compositional, coverage"
            exit 1
            ;;
    esac
    
    echo -e "\n${BLUE}====================================${NC}"
    echo -e "${GREEN}Testing complete!${NC}"

# Generate test coverage report
coverage:
    go test ./rpcmethods/... ./streaming/... ./tests/... -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"
    @if command -v open >/dev/null 2>&1; then open coverage.html; fi

# Generate detailed coverage report with gocov
coverage-detailed:
    #!/usr/bin/env bash
    set -e
    echo "Installing coverage tools if needed..."
    go install github.com/axw/gocov/gocov@latest
    go install github.com/matm/gocov-html/cmd/gocov-html@latest
    
    echo "Generating detailed coverage report..."
    $(go env GOPATH)/bin/gocov test ./rpcmethods/... ./streaming/... ./tests/... | $(go env GOPATH)/bin/gocov-html > coverage-detailed.html
    echo "Detailed coverage report generated: coverage-detailed.html"
    if command -v open >/dev/null 2>&1; then open coverage-detailed.html; fi

# Run coverage for a specific package
coverage-pkg pkg="":
    #!/usr/bin/env bash
    set -e
    if [ -z "{{pkg}}" ]; then
        echo "Please specify a package path (e.g., ./rpcmethods)"
        exit 1
    fi
    go test -coverprofile=coverage.out {{pkg}}
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated for {{pkg}}: coverage.html"
    if command -v open >/dev/null 2>&1; then open coverage.html; fi

# Generate coverage report with function-level details
coverage-func:
    go test ./rpcmethods/... ./streaming/... ./tests/... -coverprofile=coverage.out
    go tool cover -func=coverage.out

# Check if code builds
check:
    go build ./...

# Tidy dependencies
tidy:
    go mod tidy

# Clean build artifacts
clean:
    rm -rf bin/
    rm -f coverage.out coverage.html

# Format code
fmt:
    go fmt ./...

# Run basic linter (go vet)
vet:
    go vet ./rpcmethods/... ./streaming/... ./tests/...

# Run comprehensive linter (golangci-lint)
lint:
    #!/usr/bin/env bash
    set -e
    echo "Installing golangci-lint if needed..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    echo "Running golangci-lint..."
    golangci-lint run --timeout=5m

# Fix test imports
fix-tests:
    go mod tidy
    go fmt ./tests

# Run benchmarks
bench pattern="":
    #!/usr/bin/env bash
    set -e
    if [ -z "{{pattern}}" ]; then
        go test -bench=. -benchmem ./rpcmethods/... ./streaming/... ./tests/...
    else
        go test -bench={{pattern}} -benchmem ./rpcmethods/... ./streaming/... ./tests/...
    fi

# Generate test coverage badge for README
coverage-badge:
    #!/usr/bin/env bash
    set -e
    echo "Installing coverage tools if needed..."
    go install github.com/axw/gocov/gocov@latest
    go install github.com/matm/gocov-html/cmd/gocov-html@latest
    
    echo "Generating coverage statistics..."
    go test ./rpcmethods/... ./streaming/... ./tests/... -coverprofile=coverage.out
    
    # Extract coverage percentage
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo "Code coverage: $COVERAGE"
    
    echo "To add this badge to your README.md, use:"
    echo "![Coverage](https://img.shields.io/badge/coverage-$COVERAGE-brightgreen)"