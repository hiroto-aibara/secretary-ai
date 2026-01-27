#!/bin/bash
set -e

# Trust mise config in this directory
mise trust

# Install tools (go, node, golangci-lint, goimports)
mise install

# Install root dependencies (husky, lint-staged)
npm install

# Install frontend dependencies
cd web && npm install && cd ..

# Download Go modules
mise exec -- go mod download

echo "Setup complete"
