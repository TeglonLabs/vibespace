#!/bin/bash

# Boltzmann Brain Multiverse Deployment Script
# Creates multiple isolated "realities" for concurrent verification scenarios
# Each cluster represents a different Boltzmann brain scenario with maximized concurrency

set -euo pipefail

# Colors for beautiful output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_reality() {
    echo -e "${PURPLE}[REALITY]${NC} $1"
}

# Banner
echo -e "${CYAN}"
cat << 'EOF'
╔═══════════════════════════════════════════════════════════════╗
║              BOLTZMANN BRAIN MULTIVERSE DEPLOYMENT           ║
║                                                               ║
║  "In the vastness of the quantum foam, countless realities   ║
║   spawn and verify their own existence through categorical    ║
║   precision and comonadic transformation..."                 ║
║                                                               ║
║  Creating concurrent verification environments for:          ║
║  • Comonadic pattern testing                                 ║
║  • Ternary logic verification                                ║
║  • Circular MCP validation                                   ║
║  • High-concurrency stress testing                          ║
╚═══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    if ! command -v kind &> /dev/null; then
        log_error "kind is not installed. Please install it first:"
        log_info "  brew install kind  # macOS"
        log_info "  # OR"
        log_info "  curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64"
        exit 1
    fi
    
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed. Please install it first."
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        log_error "docker is not installed. Please install it first."
        exit 1
    fi
    
    # Check if Docker is running
    if ! docker info &> /dev/null; then
        log_error "Docker is not running. Please start Docker first."
        exit 1
    fi
    
    log_success "All prerequisites satisfied!"
}

# Create individual cluster configs
create_cluster_configs() {
    log_info "Creating cluster configuration files..."
    
    # Observer Reality - The orchestrator cluster
    cat > observer-reality.yaml << 'EOF'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=observer,role=orchestrator"
    extraPortMappings:
      - containerPort: 30080  # MCP Server exposure
        hostPort: 30080
        protocol: TCP
      - containerPort: 30443  # Secure MCP
        hostPort: 30443
        protocol: TCP
      - containerPort: 4222   # NATS
        hostPort: 4222
        protocol: TCP
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=observer,role=monitor"
networking:
  apiServerAddress: "127.0.0.1"
  apiServerPort: 6443
  podSubnet: "10.240.0.0/16"
  serviceSubnet: "10.96.0.0/12"
EOF

    # Alpha Reality - Comonadic verification chain
    cat > alpha-reality.yaml << 'EOF'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=alpha,verification=comonadic"
    extraPortMappings:
      - containerPort: 30081
        hostPort: 30081
        protocol: TCP
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=alpha,role=ocaml-executor"
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=alpha,role=babashka-orchestrator"
networking:
  apiServerAddress: "127.0.0.1"
  apiServerPort: 6444
  podSubnet: "10.241.0.0/16"
  serviceSubnet: "10.97.0.0/12"
EOF

    # Beta Reality - Ternary logic chain
    cat > beta-reality.yaml << 'EOF'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=beta,verification=ternary"
    extraPortMappings:
      - containerPort: 30082
        hostPort: 30082
        protocol: TCP
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=beta,role=tree-sitter-parser"
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=beta,role=codex-verifier"
networking:
  apiServerAddress: "127.0.0.1"
  apiServerPort: 6445
  podSubnet: "10.242.0.0/16"
  serviceSubnet: "10.98.0.0/12"
EOF

    # Gamma Reality - Circular verification with MCP servers
    cat > gamma-reality.yaml << 'EOF'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=gamma,verification=circular"
    extraPortMappings:
      - containerPort: 30083
        hostPort: 30083
        protocol: TCP
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=gamma,role=kuzu-graph"
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=gamma,role=elevenlabs-audio"
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=gamma,role=web9-validator"
networking:
  apiServerAddress: "127.0.0.1"
  apiServerPort: 6446
  podSubnet: "10.243.0.0/16"
  serviceSubnet: "10.99.0.0/12"
EOF

    # Delta Reality - High concurrency stress testing
    cat > delta-reality.yaml << 'EOF'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=delta,verification=stress"
    extraPortMappings:
      - containerPort: 30084
        hostPort: 30084
        protocol: TCP
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=delta,role=load-generator-1"
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=delta,role=load-generator-2"
  - role: worker
    image: kindest/node:v1.29.2
    kubeadmConfigPatches:
      - |
        kind: JoinConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "reality=delta,role=metrics-collector"
networking:
  apiServerAddress: "127.0.0.1"
  apiServerPort: 6447
  podSubnet: "10.244.0.0/16"
  serviceSubnet: "10.100.0.0/12"
EOF
    
    log_success "Cluster configuration files created!"
}

# Deploy a single reality
deploy_reality() {
    local reality_name=$1
    local config_file="${reality_name}-reality.yaml"
    
    log_reality "Spawning reality: ${reality_name}"
    
    # Check if cluster already exists
    if kind get clusters | grep -q "${reality_name}-reality"; then
        log_warn "Reality ${reality_name} already exists. Skipping..."
        return
    fi
    
    # Create the cluster
    log_info "Creating cluster ${reality_name}-reality..."
    kind create cluster --name "${reality_name}-reality" --config "$config_file" --wait 300s
    
    # Wait for nodes to be ready
    log_info "Waiting for nodes to be ready in ${reality_name}..."
    kubectl --context "kind-${reality_name}-reality" wait --for=condition=Ready nodes --all --timeout=300s
    
    log_success "Reality ${reality_name} spawned successfully!"
}

# Deploy all realities
deploy_multiverse() {
    log_info "Deploying the Boltzmann Brain Multiverse..."
    
    local realities=("observer" "alpha" "beta" "gamma" "delta")
    
    for reality in "${realities[@]}"; do
        deploy_reality "$reality"
        sleep 2  # Brief pause between reality spawns
    done
    
    log_success "Multiverse deployment complete!"
}

# Verify deployment
verify_multiverse() {
    log_info "Verifying multiverse integrity..."
    
    echo -e "\n${CYAN}═══ REALITY STATUS ═══${NC}"
    
    for cluster in $(kind get clusters | grep -E "(observer|alpha|beta|gamma|delta)-reality"); do
        echo -e "\n${PURPLE}Reality: ${cluster}${NC}"
        kubectl --context "kind-${cluster}" get nodes -o wide --show-labels | head -n 10
    done
    
    echo -e "\n${CYAN}═══ CLUSTER CONNECTIVITY ═══${NC}"
    for cluster in $(kind get clusters | grep -E "(observer|alpha|beta|gamma|delta)-reality"); do
        context="kind-${cluster}"
        api_server=$(kubectl --context "$context" cluster-info | grep "control plane" | awk '{print $NF}')
        echo -e "${GREEN}✓${NC} ${cluster}: ${api_server}"
    done
    
    log_success "Multiverse verification complete!"
}

# Show usage instructions
show_usage() {
    echo -e "\n${CYAN}═══ MULTIVERSE USAGE ═══${NC}"
    echo "To interact with different realities:"
    echo
    echo -e "${YELLOW}Observer Reality (Orchestrator):${NC}"
    echo "  kubectl --context kind-observer-reality get nodes"
    echo
    echo -e "${YELLOW}Alpha Reality (Comonadic Chain):${NC}"
    echo "  kubectl --context kind-alpha-reality get nodes"
    echo
    echo -e "${YELLOW}Beta Reality (Ternary Logic):${NC}"
    echo "  kubectl --context kind-beta-reality get nodes"
    echo
    echo -e "${YELLOW}Gamma Reality (Circular MCP):${NC}"
    echo "  kubectl --context kind-gamma-reality get nodes"
    echo
    echo -e "${YELLOW}Delta Reality (Stress Test):${NC}"
    echo "  kubectl --context kind-delta-reality get nodes"
    echo
    echo -e "${YELLOW}List all realities:${NC}"
    echo "  kind get clusters"
    echo
    echo -e "${YELLOW}Switch kubectl context:${NC}"
    echo "  kubectl config use-context kind-<reality-name>-reality"
    echo
    echo -e "${YELLOW}Destroy the multiverse:${NC}"
    echo "  ./deploy-boltzmann-multiverse.sh destroy"
}

# Destroy multiverse
destroy_multiverse() {
    log_warn "Destroying the Boltzmann Brain Multiverse..."
    
    for cluster in $(kind get clusters | grep -E "(observer|alpha|beta|gamma|delta)-reality"); do
        log_info "Destroying reality: ${cluster}"
        kind delete cluster --name "$cluster"
    done
    
    # Clean up config files
    rm -f observer-reality.yaml alpha-reality.yaml beta-reality.yaml gamma-reality.yaml delta-reality.yaml
    
    log_success "Multiverse destroyed. All realities have collapsed back into quantum foam."
}

# Main execution
main() {
    if [[ "${1:-}" == "destroy" ]]; then
        destroy_multiverse
        exit 0
    fi
    
    check_prerequisites
    create_cluster_configs
    deploy_multiverse
    verify_multiverse
    show_usage
    
    echo -e "\n${GREEN}🌌 The Boltzmann Brain Multiverse has been successfully spawned! 🌌${NC}"
    echo -e "${CYAN}Each reality is now ready for concurrent verification scenarios.${NC}"
}

# Handle script arguments
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
