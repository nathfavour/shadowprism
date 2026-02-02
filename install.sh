#!/bin/bash

# ShadowPrism: Seamless Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/nathfavour/shadowprism/main/install.sh | bash

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🛡️  ShadowPrism Installation Started${NC}"

# 1. Dependency Checks
echo -e "${BLUE}🔍 Checking dependencies...${NC}"

check_cmd() {
    if ! command -v $1 &> /dev/null;
    then
        echo -e "${RED}❌ Error: $1 is not installed.${NC}"
        return 1
    fi
    return 0
}

check_cmd "go" || exit 1
check_cmd "cargo" || exit 1
check_cmd "git" || exit 1

# 2. Setup Directory Structure
echo -e "${BLUE}📂 Setting up ~/.shadowprism...${NC}"
PRISM_HOME="$HOME/.shadowprism"
mkdir -p "$PRISM_HOME/bin"
mkdir -p "$PRISM_HOME/logs"

# 3. Build Core (Rust Muscle)
echo -e "${BLUE}📦 Compiling ShadowPrism Core (Rust)...${NC}"
if [ -d "core" ]; then
    cd core
    cargo build --release
    cp target/release/shadowprism-core "$PRISM_HOME/bin/"
    cd ..
else
    echo -e "${YELLOW}⚠️  Core directory not found. Are you running this from the repo root?${NC}"
    exit 1
fi

# 4. Build CLI (Go Brain)
echo -e "${BLUE}📦 Compiling ShadowPrism CLI (Go)...${NC}"
if [ -d "cli" ]; then
    cd cli
    go build -o shadowprism .
    
    # Enforce ~/.local/bin
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
    
    echo -e "${BLUE}🚚 Installing binary to $INSTALL_DIR...${NC}"
    cp shadowprism "$INSTALL_DIR/"
    cd ..
else
    echo -e "${YELLOW}⚠️  CLI directory not found.${NC}"
    exit 1
fi

# 5. Finalize
echo -e "\n${GREEN}✅ ShadowPrism installed successfully!${NC}"
echo -e "${BLUE}--------------------------------------------------${NC}"
echo -e "🚀 ${BLUE}CLI Binary:${NC}   $INSTALL_DIR/shadowprism"
echo -e "⚙️  ${BLUE}Core Engine:${NC}  $PRISM_HOME/bin/shadowprism-core"
echo -e "📂 ${BLUE}Config Dir:${NC}   $PRISM_HOME"
echo -e "${BLUE}--------------------------------------------------${NC}"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}👉 Note: Please add $INSTALL_DIR to your PATH to run 'shadowprism' from anywhere.${NC}"
    echo -e "   export PATH=\"\$PATH:$INSTALL_DIR\""
fi

echo -e "\nRun ${GREEN}shadowprism --help${NC} to get started."
