#!/bin/bash
# install-bv.sh - Install bv and dependencies
set -e

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

echo "=== Installing bv dependencies ==="

# Install Go if needed
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    GO_VERSION="1.22.5"
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" | sudo tar -C /usr/local -xzf -
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
else
    echo "Go already installed: $(go version)"
fi

# Install dolt if needed
if ! command -v dolt &> /dev/null; then
    echo "Installing dolt..."
    sudo curl -L https://github.com/dolthub/dolt/releases/latest/download/install.sh | sudo bash
else
    echo "Dolt already installed: $(dolt version)"
fi

echo "=== Building bv ==="

# Clone and build
TMPDIR=$(mktemp -d)
git clone https://github.com/Jaggerxtrm/beads_viewer.git "$TMPDIR/bv"
cd "$TMPDIR/bv"
go build -o "$INSTALL_DIR/bv" ./cmd/bv/
rm -rf "$TMPDIR"

echo "=== Done ==="
echo "bv installed to $INSTALL_DIR/bv"
echo ""
echo "Usage: bv --robot-triage"
echo "Set BEADS_DIR=/path/to/project/.beads if not in project root"
