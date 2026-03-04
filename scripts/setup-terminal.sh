#!/bin/bash

# Configuration
API_URL=${GLANCE_API:-"http://localhost:15501"}

# Colors for output
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if we are being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo -e "${RED}Error: This script must be sourced.${NC}"
    echo "Usage: source ./scripts/setup-terminal.sh"
    exit 1
fi

# Prepare query parameters
QUERY=""
if [[ -v NO_PROXY ]]; then
    # NO_PROXY is defined (could be empty)
    if [ -z "$NO_PROXY" ]; then
        QUERY="?no_proxy=__EMPTY__"
    else
        # Simple URL encoding for NO_PROXY
        ENCODED_NO_PROXY=$(echo "$NO_PROXY" | sed 's/ /%20/g; s/,/%2C/g')
        QUERY="?no_proxy=$ENCODED_NO_PROXY"
    fi
fi

# Try to fetch setup from API
SETUP_SCRIPT=$(curl -s "${API_URL}/api/client/terminal/setup${QUERY}")

if [ $? -ne 0 ] || [ -z "$SETUP_SCRIPT" ]; then
    echo -e "${RED}Error: Could not connect to Glance API at ${API_URL}${NC}"
    echo "Make sure Glance is running and API_URL is correct."
    echo "You can set GLANCE_API environment variable if you are using a non-default port."
    return 1
fi

# Execute the setup (exports, aliases, and logging)
eval "$SETUP_SCRIPT"
