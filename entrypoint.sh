#!/bin/sh
set -e

# Function to check if the server is up
check_server() {
    curl -s -o /dev/null -w "%{http_code}" http://localhost:${PORT:-8080}/authorize
}

# Wait for the server to be ready
echo "Starting OAuth mock server..."
/oauth-mock-server &
SERVER_PID=$!

echo "Waiting for OAuth mock server to be ready..."
for i in $(seq 1 30); do
    if [ "$(check_server)" = "200" ]; then
        echo "OAuth mock server is up and running"
        break
    fi
    sleep 1
done

if [ "$i" = "30" ]; then
    echo "Failed to start OAuth mock server"
    exit 1
fi

# Keep the action running by tailing the server process
tail --pid=$SERVER_PID -f /dev/null