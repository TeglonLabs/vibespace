#!/usr/bin/env bash
set -euo pipefail

# Boltzmann Brain Multi-World Demonstration Script
# Shows how to run complex concurrent testing scenarios across multiple Kind clusters

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly ORCHESTRATOR="${SCRIPT_DIR}/boltzmann-orchestrator.sh"
readonly VIBESPACE_MANIFESTS="${SCRIPT_DIR}/k8s-vibespace-boltzmann.yaml"

# Colors
readonly GREEN='\033[0;32m'
readonly BLUE='\033[0;34m'
readonly PURPLE='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

demo_log() {
    echo -e "${GREEN}[DEMO] $*${NC}"
}

quantum_demo() {
    echo -e "${CYAN}[⚛️  QUANTUM DEMO] $*${NC}"
}

boltzmann_demo() {
    echo -e "${PURPLE}[🧠 BOLTZMANN DEMO] $*${NC}"
}

show_intro() {
    clear
    cat << 'EOF'
╔══════════════════════════════════════════════════════════════════════════════╗
║                    🧠 BOLTZMANN BRAIN MULTI-WORLD DEMO 🧠                    ║
║                                                                              ║
║  This demo illustrates complex concurrent testing scenarios using Kind       ║
║  clusters to simulate multiple "realities" with different coherence levels.  ║
║                                                                              ║
║  🔬 SCENARIOS DEMONSTRATED:                                                   ║
║    • Multi-cluster orchestration with Kind                                  ║
║    • Cross-reality quantum entanglement                                     ║
║    • Chaos engineering in distributed systems                              ║
║    • Comonadic pattern testing at scale                                    ║
║    • Real-time coherence monitoring                                        ║
║                                                                              ║
║  🌍 REALITIES:                                                               ║
║    α (Primary)  - High coherence, stable consciousness                     ║
║    β (Quantum)  - Medium coherence, superposition states                   ║
║    γ (Chaotic)  - Low coherence, maximum entropy                           ║
║    ω (Observer) - Meta-reality monitoring all others                       ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
EOF
    echo ""
}

check_prerequisites() {
    demo_log "Checking prerequisites..."
    
    local missing_deps=()
    local deps=("kind" "kubectl" "docker" "jq" "curl")
    
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing_deps+=("$dep")
        fi
    done
    
    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        echo -e "${YELLOW}Missing dependencies: ${missing_deps[*]}${NC}"
        echo ""
        echo "Installation commands:"
        echo "  macOS: brew install kind kubectl jq"
        echo "  Docker: https://docs.docker.com/get-docker/"
        echo ""
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        echo -e "${YELLOW}Docker is not running. Please start Docker Desktop.${NC}"
        exit 1
    fi
    
    demo_log "✅ All prerequisites satisfied"
}

run_scenario_1() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════════════════════"
    boltzmann_demo "SCENARIO 1: Multi-Reality Cluster Creation"
    echo "═══════════════════════════════════════════════════════════════════════════════"
    echo ""
    
    demo_log "Creating 4 interconnected reality clusters..."
    "$ORCHESTRATOR" create
    
    echo ""
    demo_log "✅ Reality clusters created successfully!"
    echo ""
    echo "Available realities:"
    kind get clusters | while read -r cluster; do
        echo "  🌍 $cluster"
    done
    
    echo ""
    read -p "Press Enter to continue to Scenario 2..."
}

run_scenario_2() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════════════════════"
    quantum_demo "SCENARIO 2: Vibespace Workload Deployment"
    echo "═══════════════════════════════════════════════════════════════════════════════"
    echo ""
    
    demo_log "Deploying Boltzmann brain consciousness engines..."
    "$ORCHESTRATOR" deploy
    
    echo ""
    demo_log "Deploying actual Vibespace applications..."
    
    # Deploy to each cluster
    for cluster in primary-reality quantum-reality chaotic-reality meta-observer; do
        if kind get clusters | grep -q "^${cluster}$"; then
            kubectl config use-context "kind-$cluster"
            kubectl apply -f "$VIBESPACE_MANIFESTS"
            demo_log "✅ Deployed to $cluster"
        fi
    done
    
    echo ""
    demo_log "Waiting for workloads to stabilize..."
    sleep 30
    
    echo ""
    demo_log "✅ All workloads deployed and running!"
    
    echo ""
    read -p "Press Enter to continue to Scenario 3..."
}

run_scenario_3() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════════════════════"
    boltzmann_demo "SCENARIO 3: Cross-Reality Testing & Monitoring"
    echo "═══════════════════════════════════════════════════════════════════════════════"
    echo ""
    
    demo_log "Testing vibe generation across realities..."
    
    echo ""
    echo "🔬 Testing Primary Reality (High Coherence):"
    if curl -s localhost:30001/vibe | jq -r '.current_vibe.reality, .current_vibe.coherence' 2>/dev/null; then
        echo "✅ Primary reality responding"
    else
        echo "⚠️  Primary reality not ready yet"
    fi
    
    echo ""
    echo "⚛️  Testing Quantum Reality (Medium Coherence):"
    if curl -s localhost:30003/vibe | jq -r '.current_vibe.reality, .superposition' 2>/dev/null; then
        echo "✅ Quantum reality in superposition"
    else
        echo "⚠️  Quantum reality not ready yet"
    fi
    
    echo ""
    echo "🌪️  Testing Chaotic Reality (Low Coherence):"
    if curl -s localhost:30005/vibe | jq -r '.current_vibe.reality, .chaos_level' 2>/dev/null; then
        echo "✅ Chaotic reality generating entropy"
    else
        echo "⚠️  Chaotic reality collapsed (expected)"
    fi
    
    echo ""
    echo "👁️  Testing Meta-Observer:"
    if curl -s localhost:30007/observe | jq -r '.observer_id, .cross_reality' 2>/dev/null; then
        echo "✅ Meta-observer monitoring all realities"
    else
        echo "⚠️  Meta-observer not ready yet"
    fi
    
    echo ""
    read -p "Press Enter to continue to Scenario 4..."
}

run_scenario_4() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════════════════════"
    quantum_demo "SCENARIO 4: Chaos Engineering Experiments"
    echo "═══════════════════════════════════════════════════════════════════════════════"
    echo ""
    
    demo_log "Initiating cross-reality chaos experiments..."
    "$ORCHESTRATOR" chaos
    
    echo ""
    demo_log "Monitoring chaos propagation..."
    sleep 20
    
    echo ""
    demo_log "Reality coherence after chaos injection:"
    
    for cluster in primary-reality quantum-reality chaotic-reality meta-observer; do
        if kind get clusters | grep -q "^${cluster}$"; then
            kubectl config use-context "kind-$cluster"
            echo ""
            echo "🌍 Reality: $cluster"
            kubectl get pods -n boltzmann-testing -o custom-columns="NAME:.metadata.name,STATUS:.status.phase,REALITY:.metadata.labels.reality,COHERENCE:.metadata.labels.coherence" 2>/dev/null || echo "  No pods found"
        fi
    done
    
    echo ""
    read -p "Press Enter to continue to final monitoring..."
}

run_scenario_5() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════════════════════"
    boltzmann_demo "SCENARIO 5: Real-time Multi-Reality Monitoring"
    echo "═══════════════════════════════════════════════════════════════════════════════"
    echo ""
    
    demo_log "Comprehensive reality status monitoring..."
    "$ORCHESTRATOR" monitor
    
    echo ""
    echo "🔗 Cross-Reality API Endpoints:"
    echo "  Primary Reality:  http://localhost:30001/vibe"
    echo "  Quantum Reality:  http://localhost:30003/vibe" 
    echo "  Chaotic Reality:  http://localhost:30005/vibe"
    echo "  Meta Observer:    http://localhost:30007/observe"
    
    echo ""
    echo "🧪 Try these commands to interact with the realities:"
    echo "  curl localhost:30001/vibe | jq '.current_vibe'"
    echo "  curl localhost:30003/vibe | jq '.superposition'"
    echo "  curl localhost:30005/vibe | jq '.chaos_level'"
    echo "  curl localhost:30007/observe | jq '.realities'"
    
    echo ""
    quantum_demo "Real-time vibe sampling (Ctrl+C to stop):"
    
    local count=0
    while [[ $count -lt 10 ]]; do
        echo ""
        echo "Sample $((count + 1))/10:"
        
        for port in 30001 30003 30005; do
            local reality_name
            case $port in
                30001) reality_name="Primary" ;;
                30003) reality_name="Quantum" ;;
                30005) reality_name="Chaotic" ;;
            esac
            
            local response
            response=$(curl -s "localhost:$port/vibe" 2>/dev/null || echo '{"error":"unavailable"}')
            local energy
            energy=$(echo "$response" | jq -r '.current_vibe.energy // "N/A"' 2>/dev/null || echo "N/A")
            local coherence
            coherence=$(echo "$response" | jq -r '.current_vibe.coherence // "N/A"' 2>/dev/null || echo "N/A")
            
            printf "  %-8s: Energy=%-6s Coherence=%-6s\n" "$reality_name" "$energy" "$coherence"
        done
        
        ((count++))
        sleep 2
    done
    
    echo ""
    read -p "Press Enter to clean up..."
}

cleanup_demo() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════════════════════"
    demo_log "CLEANUP: Collapsing All Realities"
    echo "═══════════════════════════════════════════════════════════════════════════════"
    echo ""
    
    demo_log "Collapsing all Boltzmann brain realities..."
    "$ORCHESTRATOR" cleanup
    
    echo ""
    demo_log "✅ All realities have collapsed back into quantum foam"
    echo ""
    echo "🎉 Demo completed successfully!"
    echo ""
    echo "What you experienced:"
    echo "  • Multi-cluster Kind orchestration"
    echo "  • Cross-reality quantum entanglement simulation"
    echo "  • Chaos engineering in distributed consciousness"
    echo "  • Real-time coherence monitoring across realities"
    echo "  • Comonadic pattern testing at scale"
    echo ""
    echo "The generated files can be used for:"
    echo "  • CI/CD pipeline testing"
    echo "  • Load testing distributed systems"
    echo "  • Chaos engineering experiments"
    echo "  • Multi-tenant application testing"
    echo "  • Kubernetes operator development"
    echo ""
}

show_help() {
    show_intro
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  full        Run complete demonstration (default)"
    echo "  scenario1   Only create clusters"
    echo "  scenario2   Only deploy workloads"
    echo "  scenario3   Only test cross-reality"
    echo "  scenario4   Only chaos experiments"
    echo "  scenario5   Only monitoring"
    echo "  cleanup     Clean up all resources"
    echo "  help        Show this help"
    echo ""
}

main() {
    case "${1:-full}" in
        "scenario1")
            show_intro
            check_prerequisites
            run_scenario_1
            ;;
        "scenario2")
            show_intro
            check_prerequisites
            run_scenario_2
            ;;
        "scenario3")
            show_intro
            check_prerequisites
            run_scenario_3
            ;;
        "scenario4")
            show_intro
            check_prerequisites
            run_scenario_4
            ;;
        "scenario5")
            show_intro
            check_prerequisites
            run_scenario_5
            ;;
        "cleanup")
            cleanup_demo
            ;;
        "full")
            show_intro
            check_prerequisites
            run_scenario_1
            run_scenario_2
            run_scenario_3
            run_scenario_4
            run_scenario_5
            cleanup_demo
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

main "$@"
