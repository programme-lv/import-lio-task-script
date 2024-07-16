#!/bin/bash
set -ex

# Define the directory and script name
SCRIPT_NAME="import-lio-task-script"
INSTALL_DIR="/usr/local/bin"

# Ensure go.mod exists, or initialize the Go module
if [ ! -f "go.mod" ]; then
    go mod init "$SCRIPT_NAME"
fi

# Get the necessary dependency
go get github.com/pelletier/go-toml/v2

# Build the Go script
go build -o "$SCRIPT_NAME" ./cmd/main.go

# Install the built script to /usr/local/bin
sudo mv "$SCRIPT_NAME" "$INSTALL_DIR"

# Ensure the script is executable
sudo chmod +x "$INSTALL_DIR/$SCRIPT_NAME"

echo "Installation complete. You can now run the script using '$SCRIPT_NAME'."
