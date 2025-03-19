#!/bin/bash

# Determine NATS server URL and Stream ID
NATS_URL="nats://nonlocal.info:4222"
STREAM_ID="ies"
USER_ID=""

# Process command line arguments
if [ "$#" -ge 1 ]; then
    NATS_URL="$1"
fi

if [ "$#" -ge 2 ]; then
    STREAM_ID="$2"
fi

if [ "$#" -ge 3 ]; then
    USER_ID="$3"
fi

# Ensure we're in the right directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Compile the subscriber
echo "Compiling NATS subscriber..."
go build -o nats_subscriber nats_subscriber.go

if [ $? -ne 0 ]; then
    echo "Error compiling subscriber. Exiting."
    exit 1
fi

# Run the subscriber
echo "Starting NATS subscriber..."
echo "- Server URL: $NATS_URL"
echo "- Stream ID:  $STREAM_ID"
if [ -n "$USER_ID" ]; then
    echo "- User ID:    $USER_ID"
    ./nats_subscriber "$NATS_URL" "$STREAM_ID" "$USER_ID"
else
    echo "- User ID:    (none - subscribing to public streams only)"
    ./nats_subscriber "$NATS_URL" "$STREAM_ID"
fi