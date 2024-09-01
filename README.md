# OAuth Mock Server

This is a simple OAuth 2.0 mock server implemented in Go. It's designed to simulate an OAuth 2.0 provider for testing purposes, especially in CI/CD environments.

## Features

- Simulates OAuth 2.0 authorization flow
- Provides JWT access tokens
- Returns static user information
- Easy to configure via environment variables
- Lightweight and easy to deploy

## Prerequisites

- Go 1.16 or higher
- Docker (optional, for containerized deployment)

## Configuration

The server can be configured using the following environment variables:

- `PORT`: The port on which the server will run (default: 8080)
- `CLIENT_ID`: The client ID to use (default: "test-client")
- `CLIENT_SECRET`: The client secret to use (default: "test-secret")

## Running the Server

### Locally

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/oauth-mock-server.git
   cd oauth-mock-server
   ```

2. Run the server:
   ```
   go run main.go
   ```

### Using Docker

1. Build the Docker image:
   ```
   docker build -t oauth-mock-server .
   ```

2. Run the container:
   ```
   docker run -p 8080:8080 -e CLIENT_ID=my-client -e CLIENT_SECRET=my-secret oauth-mock-server
   ```

## Usage

1. Initiate the OAuth flow by sending a GET request to `/authorize` with the following parameters:
   - `client_id`: Your client ID
   - `redirect_uri`: Your callback URL
   - `state`: A random state value for security

2. The server will display a consent screen and automatically approve after 2 seconds, redirecting to your `redirect_uri` with an authorization code.

3. Exchange the authorization code for an access token by sending a POST request to `/token` with the following parameters:
   - `grant_type`: "authorization_code"
   - `code`: The authorization code received in step 2
   - `client_id`: Your client ID
   - `client_secret`: Your client secret

4. Use the received access token to make requests to the `/userinfo` endpoint by including it in the Authorization header:
   ```
   Authorization: Bearer <your_access_token>
   ```

## License

This project is licensed under the MIT License.
