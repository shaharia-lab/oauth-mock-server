#!/bin/sh
set -e

# Start the OAuth mock server in the background
/oauth-mock-server &

# Wait for the server to be ready
for i in $(seq 1 30); do
    if curl -s http://localhost:${PORT:-8080}/authorize > /dev/null; then
        echo "OAuth mock server is up and running"
        exit 0
    fi
    echo "Waiting for OAuth mock server to be ready..."
    sleep 1
done

echo "Failed to start OAuth mock server"
exit 1