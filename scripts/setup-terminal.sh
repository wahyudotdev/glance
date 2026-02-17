#!/bin/bash

# Configuration
API_URL=${MCPROXY_API:-"http://localhost:8081"}

# Colors for output
BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if we are being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo -e "${RED}Error: This script must be sourced.${NC}"
    echo "Usage: source ./scripts/setup-terminal.sh"
    exit 1
fi

# Try to fetch setup from API
SETUP_SCRIPT=$(curl -s "${API_URL}/api/client/terminal/setup")

if [ $? -ne 0 ] || [ -z "$SETUP_SCRIPT" ]; then
    echo -e "${RED}Error: Could not connect to glance API at ${API_URL}${NC}"
    echo "Make sure glance is running."
    return 1
fi

# Execute the setup (exports and aliases)
eval "$SETUP_SCRIPT"

# Extract proxy address for display
PROXY_HOST_PORT=$(echo "$SETUP_SCRIPT" | grep "export HTTP_PROXY=" | sed 's/export HTTP_PROXY=http:\/\///' | tr -d '"')

echo -e "${GREEN}✓ Terminal configured for interception${NC}"
echo -e "${BLUE}Proxy:${NC} $PROXY_HOST_PORT"
echo -e "${BLUE}CA Cert:${NC} $SSL_CERT_FILE"
echo ""
echo "Type 'unproxy' to restore original terminal settings."
