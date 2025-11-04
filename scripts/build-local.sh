#!/bin/bash

# Script to build and install Terraform provider locally for testing
# Usage: ./scripts/build-local.sh [version] [architecture]

set -e

# Default values
VERSION=${1:-100.0.0}
ARCH=${2:-darwin_arm64}

echo "Building local Terraform provider version ${VERSION} for ${ARCH}..."

# Step 1: Ensure we're in the project directory
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod not found. Please run this script from the project root."
    exit 1
fi

# Step 2: Create the local version folder
# Using logzio.io as the hostname for local development
PLUGIN_PATH=~/.terraform.d/plugins/logzio.io/logzio/logzio/${VERSION}/${ARCH}
echo "Creating plugin directory: ${PLUGIN_PATH}"
mkdir -p ${PLUGIN_PATH}

# Step 3: Build the provider
echo "Building provider..."
GO111MODULE=on go build -o ./build/terraform-provider-logzio

if [ $? -ne 0 ]; then
    echo "Error: Build failed"
    exit 1
fi

# Step 4: Copy provider to local plugins
echo "Copying provider to ${PLUGIN_PATH}..."
cp ./build/terraform-provider-logzio ${PLUGIN_PATH}/

# Make it executable
chmod +x ${PLUGIN_PATH}/terraform-provider-logzio

echo ""
echo "✅ Local provider installed successfully!"
echo ""
echo "To use it in your Terraform configuration, add:"
echo ""
echo "terraform {"
echo "  required_providers {"
echo "    logzio = {"
echo "      source  = \"logzio.io/logzio/logzio\""
echo "      version = \"${VERSION}\""
echo "    }"
echo "  }"
echo "}"
echo ""
echo "Then run: terraform init"

