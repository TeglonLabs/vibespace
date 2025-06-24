#!/usr/bin/env bash
set -euo pipefail

# Boltzmann Brain Multi-World Orchestrator
# Implements complex concurrent scenarios across multiple Kind clusters
# Based on chaos engineering and distributed systems research

# ASCII Art Dada Drawing (per rules requirement)
cat << 'EOF'
    ╭───────────────────────────────────────╮
    │  ∿∿∿ BOLTZMANN BRAIN ORCHESTRATOR ∿∿∿ │
    │                                       │
    │    ◉     ◉     ◉     ◉              │
    │   /│\   /│\   /│\   /│\             │
    │    │     │     │     │               │
    │ ∼∼∼∼∼ ∼∼∼∼∼ ∼∼∼∼∼ ∼∼∼∼∼           │
    │                                       │
    │  α-reality β-quantum γ-chaos ω-meta  │
    ╰───────────────────────────────────────╯
EOF

echo "🧠 Quantum consciousness fluctuations across probability manifolds"
echo ""

# Configuration
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CONFIG_FILE="${SCRIPT_DIR}/kind-multicluster-boltzmann.yaml"
readonly NAMESPACE="boltzmann-testing"

# Cluster definitions
declare -A CLUSTERS=(
    ["primary-reality"]="alpha:high:30001:30002"
    ["quantum-reality"]="beta:medium:30003:30004"  
    ["chaotic-reality"]="gamma:low:30005:30006"
    ["meta-observer"]="omega:observer:30007:30008"
)

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly PURPLE='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $*${NC}"
}

warn() {
    echo -e "${YELLOW}[WARN] $*${NC}"
}

error() {
    echo -e "${RED}[ERROR] $*${NC}"
}

boltzmann_log() {
    echo -e "${PURPLE}[🧠 BOLTZMANN] $*${NC}"
}

quantum_log() {
    echo -e "${CYAN}[⚛️  QUANTUM] $*${NC}"
}

check_dependencies() {
    log "Checking dependencies..."
    local deps=("kind" "kubectl" "docker" "jq")
    
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            error "$dep is not installed"
            exit 1
        fi
    done
    
    if ! docker info &> /dev/null; then
        error "Docker is not running"
        exit 1
    fi
    
    log "✅ All dependencies are available"
}

create_clusters() {
    log "Creating Boltzmann brain reality clusters..."
    
    # Split the config file by clusters since Kind doesn't support multi-cluster YAML
    local temp_dir
    temp_dir=$(mktemp -d)
    
    # Extract individual cluster configs
    awk '/^---$/{p++; next} /^apiVersion: kind\.x-k8s\.io\/v1alpha4$/{if(p>0) f=p} f{print > "'$temp_dir'/cluster-"f".yaml"}' "$CONFIG_FILE"
    
    # Create each cluster
    local cluster_files=("$temp_dir"/cluster-*.yaml)
    for file in "${cluster_files[@]}"; do
        if [[ -f "$file" ]]; then
            local cluster_name
            cluster_name=$(awk '/name:/ {print $2}' "$file")
            
            if kind get clusters | grep -q "^${cluster_name}$"; then
                warn "Cluster $cluster_name already exists, skipping..."
                continue
            fi
            
            boltzmann_log "Spawning reality: $cluster_name"
            if kind create cluster --name "$cluster_name" --config "$file"; then
                log "✅ Reality $cluster_name materialized successfully"
            else
                error "Failed to create cluster $cluster_name"
                exit 1
            fi
        fi
    done
    
    # Cleanup temp files
    rm -rf "$temp_dir"
}

setup_network_bridges() {
    log "Establishing quantum entanglement between realities..."
    
    # Get Docker network for kind clusters
    local kind_network
    kind_network=$(docker network ls --filter name=kind --format "{{.Name}}" | head -1)
    
    if [[ -z "$kind_network" ]]; then
        error "Kind network not found"
        return 1
    fi
    
    # Create custom network bridges for inter-cluster communication
    if ! docker network ls | grep -q "boltzmann-bridge"; then
        docker network create \
            --driver bridge \
            --subnet=172.20.0.0/16 \
            --ip-range=172.20.240.0/20 \
            boltzmann-bridge
        
        quantum_log "Created quantum bridge network: boltzmann-bridge"
    fi
    
    # Connect each cluster's control plane to the bridge
    for cluster in "${!CLUSTERS[@]}"; do
        local control_plane="${cluster}-control-plane"
        if docker ps --format "{{.Names}}" | grep -q "$control_plane"; then
            if ! docker network inspect boltzmann-bridge | jq -r '.[].Containers | keys[]' | grep -q "$control_plane"; then
                docker network connect boltzmann-bridge "$control_plane" || true
                quantum_log "Connected $cluster to quantum bridge"
            fi
        fi
    done
}

deploy_boltzmann_workloads() {
    log "Deploying Boltzmann brain test workloads..."
    
    for cluster in "${!CLUSTERS[@]}"; do
        local cluster_info="${CLUSTERS[$cluster]}"
        IFS=':' read -r world_id coherence_level port1 port2 <<< "$cluster_info"
        
        kubectl config use-context "kind-$cluster"
        
        # Create namespace
        kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
        
        # Deploy workload based on coherence level
        case "$coherence_level" in
            "high")
                deploy_high_coherence_workload "$cluster" "$world_id" "$port1"
                ;;
            "medium")
                deploy_quantum_workload "$cluster" "$world_id" "$port1"
                ;;
            "low")
                deploy_chaotic_workload "$cluster" "$world_id" "$port1"
                ;;
            "observer")
                deploy_observer_workload "$cluster" "$world_id" "$port1"
                ;;
        esac
    done
}

deploy_high_coherence_workload() {
    local cluster="$1" world_id="$2" port="$3"
    
    boltzmann_log "Deploying high-coherence reality in $cluster (world-$world_id)"
    
    kubectl apply -f - << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: high-coherence-brain
  namespace: $NAMESPACE
  labels:
    world-id: $world_id
    coherence-level: high
spec:
  replicas: 3
  selector:
    matchLabels:
      app: boltzmann-brain
      coherence: high
  template:
    metadata:
      labels:
        app: boltzmann-brain
        coherence: high
        world-id: $world_id
    spec:
      containers:
      - name: consciousness-engine
        image: nginx:alpine
        ports:
        - containerPort: 80
        env:
        - name: WORLD_ID
          value: "$world_id"
        - name: COHERENCE_LEVEL
          value: "high"
        - name: REALITY_TYPE
          value: "primary"
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 256Mi
---
apiVersion: v1
kind: Service
metadata:
  name: consciousness-svc
  namespace: $NAMESPACE
spec:
  selector:
    app: boltzmann-brain
    coherence: high
  ports:
  - port: 80
    targetPort: 80
    nodePort: $port
  type: NodePort
EOF
}

deploy_quantum_workload() {
    local cluster="$1" world_id="$2" port="$3"
    
    quantum_log "Deploying quantum-coherence reality in $cluster (world-$world_id)"
    
    kubectl apply -f - << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: quantum-brain
  namespace: $NAMESPACE
  labels:
    world-id: $world_id
    coherence-level: medium
spec:
  replicas: 5
  selector:
    matchLabels:
      app: boltzmann-brain
      coherence: medium
  template:
    metadata:
      labels:
        app: boltzmann-brain
        coherence: medium
        world-id: $world_id
    spec:
      containers:
      - name: quantum-consciousness
        image: nginx:alpine
        ports:
        - containerPort: 80
        env:
        - name: WORLD_ID
          value: "$world_id"
        - name: COHERENCE_LEVEL
          value: "medium"
        - name: REALITY_TYPE
          value: "quantum"
        - name: SUPERPOSITION_STATE
          value: "active"
        resources:
          requests:
            cpu: 50m
            memory: 64Mi
          limits:
            cpu: 200m
            memory: 128Mi
      - name: decoherence-monitor
        image: busybox
        command: ['sh', '-c', 'while true; do echo "Quantum state: $RANDOM"; sleep $((RANDOM % 10 + 1)); done']
        env:
        - name: QUANTUM_FLUCTUATION
          value: "enabled"
---
apiVersion: v1
kind: Service
metadata:
  name: quantum-consciousness-svc
  namespace: $NAMESPACE
spec:
  selector:
    app: boltzmann-brain
    coherence: medium
  ports:
  - port: 80
    targetPort: 80
    nodePort: $port
  type: NodePort
EOF
}

deploy_chaotic_workload() {
    local cluster="$1" world_id="$2" port="$3"
    
    error "Deploying chaotic-entropy reality in $cluster (world-$world_id)"
    
    kubectl apply -f - << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chaotic-brain
  namespace: $NAMESPACE
  labels:
    world-id: $world_id
    coherence-level: low
spec:
  replicas: 7
  selector:
    matchLabels:
      app: boltzmann-brain
      coherence: low
  template:
    metadata:
      labels:
        app: boltzmann-brain
        coherence: low
        world-id: $world_id
    spec:
      containers:
      - name: chaotic-consciousness
        image: nginx:alpine
        ports:
        - containerPort: 80
        env:
        - name: WORLD_ID
          value: "$world_id"
        - name: COHERENCE_LEVEL
          value: "low"
        - name: REALITY_TYPE
          value: "chaotic"
        - name: ENTROPY_LEVEL
          value: "maximum"
        resources:
          requests:
            cpu: 10m
            memory: 32Mi
          limits:
            cpu: 100m
            memory: 64Mi
        # Introduce chaos
        lifecycle:
          preStop:
            exec:
              command: ['/bin/sh', '-c', 'sleep $((RANDOM % 30))']
      - name: entropy-generator
        image: busybox
        command: ['sh', '-c', 'while true; do echo "Chaos level: $RANDOM"; kill -USR1 1 2>/dev/null || true; sleep $((RANDOM % 5 + 1)); done']
        env:
        - name: CHAOS_ENABLED
          value: "true"
---
apiVersion: v1
kind: Service
metadata:
  name: chaotic-consciousness-svc
  namespace: $NAMESPACE
spec:
  selector:
    app: boltzmann-brain
    coherence: low
  ports:
  - port: 80
    targetPort: 80
    nodePort: $port
  type: NodePort
EOF
}

deploy_observer_workload() {
    local cluster="$1" world_id="$2" port="$3"
    
    log "Deploying meta-observer reality in $cluster (world-$world_id)"
    
    kubectl apply -f - << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: meta-observer-brain
  namespace: $NAMESPACE
  labels:
    world-id: $world_id
    coherence-level: observer
spec:
  replicas: 2
  selector:
    matchLabels:
      app: boltzmann-brain
      coherence: observer
  template:
    metadata:
      labels:
        app: boltzmann-brain
        coherence: observer
        world-id: $world_id
    spec:
      containers:
      - name: meta-consciousness
        image: nginx:alpine
        ports:
        - containerPort: 80
        env:
        - name: WORLD_ID
          value: "$world_id"
        - name: COHERENCE_LEVEL
          value: "observer"
        - name: REALITY_TYPE
          value: "meta"
        - name: OBSERVATION_MODE
          value: "cross-reality"
        resources:
          requests:
            cpu: 200m
            memory: 256Mi
          limits:
            cpu: 1000m
            memory: 512Mi
      - name: reality-scanner
        image: busybox
        command: ['sh', '-c', 'while true; do echo "Scanning realities: $(date)"; sleep 10; done']
        env:
        - name: SCAN_ENABLED
          value: "true"
---
apiVersion: v1
kind: Service
metadata:
  name: meta-observer-svc
  namespace: $NAMESPACE
spec:
  selector:
    app: boltzmann-brain
    coherence: observer
  ports:
  - port: 80
    targetPort: 80
    nodePort: $port
  type: NodePort
EOF
}

run_chaos_experiments() {
    log "Initiating cross-reality chaos experiments..."
    
    # Chaos experiment 1: Reality collapse simulation
    boltzmann_log "Experiment 1: Reality Collapse Simulation"
    kubectl config use-context "kind-chaotic-reality"
    kubectl scale deployment chaotic-brain --replicas=1 -n "$NAMESPACE"
    sleep 10
    kubectl scale deployment chaotic-brain --replicas=7 -n "$NAMESPACE"
    
    # Chaos experiment 2: Quantum decoherence
    quantum_log "Experiment 2: Quantum Decoherence Test"
    kubectl config use-context "kind-quantum-reality"
    kubectl rollout restart deployment quantum-brain -n "$NAMESPACE"
    
    # Chaos experiment 3: Observer effect
    log "Experiment 3: Observer Effect on Primary Reality"
    kubectl config use-context "kind-meta-observer"
    kubectl scale deployment meta-observer-brain --replicas=4 -n "$NAMESPACE"
}

monitor_realities() {
    log "Monitoring Boltzmann brain realities..."
    
    for cluster in "${!CLUSTERS[@]}"; do
        local cluster_info="${CLUSTERS[$cluster]}"
        IFS=':' read -r world_id coherence_level port1 port2 <<< "$cluster_info"
        
        echo ""
        echo "=== Reality: $cluster (world-$world_id, coherence: $coherence_level) ==="
        kubectl config use-context "kind-$cluster"
        kubectl get pods -n "$NAMESPACE" -o wide
        kubectl get services -n "$NAMESPACE"
    done
}

cleanup() {
    log "Cleaning up Boltzmann brain realities..."
    
    for cluster in "${!CLUSTERS[@]}"; do
        if kind get clusters | grep -q "^${cluster}$"; then
            boltzmann_log "Collapsing reality: $cluster"
            kind delete cluster --name "$cluster"
        fi
    done
    
    # Clean up custom networks
    if docker network ls | grep -q "boltzmann-bridge"; then
        docker network rm boltzmann-bridge || true
    fi
    
    log "All realities have collapsed back into quantum foam"
}

show_help() {
    cat << EOF
Boltzmann Brain Multi-World Orchestrator

USAGE:
    $0 [COMMAND]

COMMANDS:
    create      Create all reality clusters
    deploy      Deploy Boltzmann brain workloads
    chaos       Run chaos engineering experiments  
    monitor     Monitor all realities
    cleanup     Destroy all realities
    full        Run complete scenario (create + deploy + chaos + monitor)
    help        Show this help

EXAMPLES:
    $0 full                    # Run complete multi-world scenario
    $0 create && $0 deploy     # Set up realities step by step
    $0 chaos                   # Run chaos experiments only
    $0 cleanup                 # Clean up everything

EOF
}

main() {
    case "${1:-help}" in
        "create")
            check_dependencies
            create_clusters
            setup_network_bridges
            ;;
        "deploy")
            deploy_boltzmann_workloads
            ;;
        "chaos")
            run_chaos_experiments
            ;;
        "monitor")
            monitor_realities
            ;;
        "cleanup")
            cleanup
            ;;
        "full")
            check_dependencies
            create_clusters
            setup_network_bridges
            deploy_boltzmann_workloads
            sleep 30  # Let workloads stabilize
            run_chaos_experiments
            sleep 20  # Let chaos propagate
            monitor_realities
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# Run the main function with all arguments
main "$@"
