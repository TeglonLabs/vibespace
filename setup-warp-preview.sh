#!/bin/bash

# Setup script for Warp Preview MCP integration with categorical artifacts
set -e

echo "🚀 Setting up Warp Preview with Categorical Universe Artifacts"

# Check if in correct directory
if [ ! -f "mcp.json" ]; then
    echo "❌ Error: mcp.json not found. Please run this script from the project root."
    exit 1
fi

# Check if Warp is installed (optional for WarpPreview use)
if ! command -v warp &> /dev/null; then
    echo "⚠️  Warp terminal not found locally. That's OK - you can still use WarpPreview!"
    echo "   Download Warp Preview: https://www.warp.dev/download-preview"
else
    echo "✅ Warp terminal found"
fi

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go not found. Please install Go first."
    exit 1
fi

echo "✅ Prerequisites checked"

# Build the MCP server
echo "🔨 Building MCP server..."
go build -o vibespace-mcp-server cmd/server/main.go
chmod +x vibespace-mcp-server

echo "✅ MCP server built successfully"

# Test the server
echo "🧪 Testing MCP server startup..."
timeout 3s ./vibespace-mcp-server || true
echo "✅ Server test completed"

# Setup Warp MCP configuration
WARP_CONFIG_DIR="$HOME/Library/Application Support/dev.warp.Warp-Stable/mcp"
mkdir -p "$WARP_CONFIG_DIR"

# Instructions for MCP configuration
echo "📝 Setting up Warp MCP configuration..."
echo "📋 MCP Configuration for Warp:"
echo ""
cat mcp.json
echo ""
echo "💡 You can use this configuration in Warp by:"
echo "   1. Adding servers manually in Warp UI, OR"
echo "   2. Importing this mcp.json file if your client supports it"

# Create a startup script that Warp can use
cat > "$WARP_CONFIG_DIR/start-vibespace.sh" << 'EOF'
#!/bin/bash
cd "$(dirname "$0")"
# Navigate to the project directory (adjust this path as needed)
PROJECT_DIR="${PROJECT_DIR:-/Users/$USER/infinity-topos/worlds/b/vibespace-mcp-go-ternary}"
cd "$PROJECT_DIR"
exec ./vibespace-mcp-server
EOF

chmod +x "$WARP_CONFIG_DIR/start-vibespace.sh"

echo "✅ Warp configuration created at: $WARP_CONFIG_DIR"

# Instructions for manual Warp setup
echo ""
echo "🎯 Next Steps:"
echo "1. Open Warp terminal"
echo "2. Enable Warp Preview:"
echo "   - Go to https://www.warp.dev/download-preview"
echo "   - Download and install Warp Preview"
echo ""
echo "3. Add MCP Server in Warp:"
echo "   - Open Warp Drive panel"
echo "   - Navigate to Personal > MCP Servers"
echo "   - Click '+ Add'"
echo "   - Select 'CLI Server (Command)'"
echo "   - Name: vibespace-categorical"
echo "   - Command: $WARP_CONFIG_DIR/start-vibespace.sh"
echo "   - Environment variables:"
echo "     VIBESPACE_MODE=categorical"
echo "     PREVIEW_ARTIFACTS=true"
echo "     COMONADIC_CONTEXT=enabled"
echo ""
echo "4. Start the server and test with Agent Mode:"
echo "   @categorical_extract contextId=\"test-context\""
echo "   @ternary_logic_gate gateType=\"consensus\" inputA=1 inputB=1"
echo ""
echo "🔮 Your categorical universe artifacts are ready for WrapPreview!"

# Optional: Open Warp if it's not running
if ! pgrep -f "Warp" > /dev/null; then
    echo "🚀 Attempting to open Warp..."
    open -a Warp || echo "⚠️  Please open Warp manually"
fi

echo "✨ Setup complete! Your MCP server supports:"
echo "  - Comonadic extract/duplicate/extend operations"
echo "  - Ternary logic gate computations"
echo "  - Interactive context navigation"
echo "  - Live artifact preview in Warp"
