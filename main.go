package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Server struct {
	clients      map[string]string // clientID -> clientSecret
	authCodes    map[string]string // code -> clientID
	accessTokens map[string]string // token -> userID
	privateKey   *rsa.PrivateKey
	publicKey    *rsa.PublicKey
	mu           sync.Mutex
}

func NewServer() (*Server, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &Server{
		clients:      make(map[string]string),
		authCodes:    make(map[string]string),
		accessTokens: make(map[string]string),
		privateKey:   privateKey,
		publicKey:    &privateKey.PublicKey,
	}, nil
}

func (s *Server) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")

	if _, ok := s.clients[clientID]; !ok {
		http.Error(w, "Invalid client", http.StatusBadRequest)
		return
	}

	html := fmt.Sprintf(`
		<html>
			<body>
				<h1>Authorize Application</h1>
				<p>The application %s is requesting access to your account.</p>
				<p>This screen will automatically approve in 2 seconds...</p>
				<script>
					setTimeout(function() {
						window.location.href = "/authorize/approve?client_id=%s&redirect_uri=%s&state=%s";
					}, 2000);
				</script>
			</body>
		</html>
	`, clientID, clientID, redirectURI, state)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (s *Server) ApproveHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")

	code := uuid.New().String()
	s.mu.Lock()
	s.authCodes[code] = clientID
	s.mu.Unlock()

	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *Server) TokenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	grantType := r.Form.Get("grant_type")
	if grantType != "authorization_code" {
		http.Error(w, "Unsupported grant type", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")

	s.mu.Lock()
	storedClientID, codeValid := s.authCodes[code]
	delete(s.authCodes, code)
	s.mu.Unlock()

	if !codeValid || storedClientID != clientID || s.clients[clientID] != clientSecret {
		http.Error(w, "Invalid client credentials or authorization code", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	s.accessTokens[tokenString] = "user123"
	s.mu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokenString,
		"token_type":   "Bearer",
		"expires_in":   "3600",
	})
}

func (s *Server) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	accessToken := tokenParts[1]

	s.mu.Lock()
	_, valid := s.accessTokens[accessToken]
	s.mu.Unlock()

	if !valid {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"sub":   "user123",
		"name":  "John Doe",
		"email": "john.doe@example.com",
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	clientID := os.Getenv("CLIENT_ID")
	if clientID == "" {
		clientID = "test-client"
	}

	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = "test-secret"
	}

	server, err := NewServer()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	server.clients[clientID] = clientSecret

	http.HandleFunc("/authorize", server.AuthorizeHandler)
	http.HandleFunc("/authorize/approve", server.ApproveHandler)
	http.HandleFunc("/token", server.TokenHandler)
	http.HandleFunc("/userinfo", server.UserInfoHandler)

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
