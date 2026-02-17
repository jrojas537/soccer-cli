#!/bin/bash
# install.sh

# This script builds and installs the soccer-cli binary.

# Exit immediately if a command exits with a non-zero status.
set -e

# Function to check if a command exists
command_exists () {
  command -v "$1" >/dev/null 2>&1
}

# Check for Go installation
if ! command_exists go; then
  echo "Error: Go is not installed or not found in your PATH."
  echo "Please install Go (version 1.18 or higher) to build soccer-cli."
  exit 1
fi

echo "Ensuring Go modules are tidy and downloaded..."
go mod tidy
go mod download

echo "Building soccer-cli..."

# Build the Go binary. This creates a 'soccer-cli' executable in the current directory.
go build -ldflags="-X 'github.com/jrojas537/soccer-cli/cmd.version=1.0.0'" -o soccer-cli main.go

# Define the installation directory.
# We'll use ~/.local/bin, which is a common place for user-installed executables.
# Make sure to add this directory to your shell's PATH if it isn't already.
INSTALL_DIR="$HOME/.local/bin"

# Create the installation directory if it doesn't exist.
mkdir -p "$INSTALL_DIR"

# Move the binary to the installation directory.
mv soccer-cli "$INSTALL_DIR/"

echo "soccer-cli installed successfully to $INSTALL_DIR"
echo "Please ensure '$INSTALL_DIR' is in your PATH."
