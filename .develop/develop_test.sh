#!/bin/sh

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if required commands exist
REQUIRED_COMMANDS=("cd" "go" "golangci-lint")

for cmd in "${REQUIRED_COMMANDS[@]}"; do
    if ! command_exists "$cmd"; then
        echo "Error: $cmd is not installed." >&2
        exit 1
    fi
done

# Navigate to parent directory
cd ..

# Run go mod tidy
go mod tidy

# Run golangci-lint
golangci-lint run -v --timeout=10m --fix

# Run go tests with race detector and coverage
go test -race -v ./... -coverprofile ./coverage.txt

# Generate HTML coverage report
go tool cover -html=./coverage.txt