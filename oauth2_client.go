// File: oauth2_client.go

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
)

var (
	clientID      = flag.String("client_id", "test-client", "OAuth2 client ID")
	clientSecret  = flag.String("client_secret", "test-secret", "OAuth2 client secret")
	authServerURL = flag.String("auth_server", "http://localhost:9999", "OAuth2 authorization server URL")
)

func main() {
	flag.Parse()

	config := &oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  *authServerURL + "/authorize",
			TokenURL: *authServerURL + "/token",
		},
		RedirectURL: "http://localhost:8081/callback", // This URL won't be used, but is required
		Scopes:      []string{"openid", "profile", "email"},
	}

	// Generate authorization URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Please visit this URL to authorize the application: %v\n", authURL)

	// Prompt user to enter the authorization code
	fmt.Print("Enter the authorization code: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	code := scanner.Text()

	// Exchange authorization code for token
	token, err := exchangeCodeForToken(config, code)
	if err != nil {
		log.Fatalf("Unable to exchange code for token: %v", err)
	}

	fmt.Printf("Access Token: %s\n", token.AccessToken)

	// Get user info
	userInfo, err := getUserInfo(*authServerURL, token.AccessToken)
	if err != nil {
		log.Fatalf("Unable to get user info: %v", err)
	}

	fmt.Printf("User Info: %+v\n", userInfo)
}

func exchangeCodeForToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	values := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"redirect_uri":  {config.RedirectURL},
	}

	resp, err := http.PostForm(config.Endpoint.TokenURL, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %s - %s", resp.Status, string(body))
	}

	var token oauth2.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func getUserInfo(serverURL, accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", serverURL+"/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s - %s", resp.Status, string(body))
	}

	var userInfo map[string]interface{}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
