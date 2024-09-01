#!/bin/sh
set -e

# Start the OAuth mock server in the background
/oauth-mock-server &

# Store the PID of the server
SERVER_PID=$!

# Function to check if the server is up
check_server() {
    curl -s -o /dev/null -w "%{http_code}" http://localhost:${PORT:-8080}/authorize
}

# Wait for the server to be ready
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

# Keep the action running
tail -f /dev/null &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?