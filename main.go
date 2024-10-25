package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	// "github.com/gin-gonic/gin"
)

const (
    access_token_req_url = "https://test.api.amadeus.com/v1/security/oauth2/token"
)

// Auth Token struct
type AuthToken struct {
    Type string `json:"type"`
    Username string `json:"username"`
    Client_id string `json:"client_id"`
    TokenType string `json:"token_type"`
    AccessToken string `json:"access_token"`
}


func main() {
    my_auth_token, err := getAuthToken()
    if err != nil {
        fmt.Println("Failed to get acces token")
        return
    }

    fmt.Println(my_auth_token.Type)
    fmt.Println(my_auth_token.Username)
    fmt.Println(my_auth_token.Client_id)
    fmt.Println(my_auth_token.TokenType)
    fmt.Println(my_auth_token.AccessToken)

}

// getAuthToken responds with filling an AuthToken struct needed to get the AccessToken
func getAuthToken() (AuthToken, error) {
    var API_KEY = os.Getenv("AMAD_API_KEY")
    var SECRET = os.Getenv("AMAD_API_SECRET")

    fmt.Println("API_KEY: ", API_KEY)
    fmt.Println("API_SECRET: ", SECRET)

    var authToken AuthToken

    data :="grant_type=client_credentials&client_id=" + API_KEY + "&client_secret=" + SECRET
    req, err := http.NewRequest("POST", access_token_req_url, strings.NewReader(data))
    if err != nil {
        return authToken, err
    }
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return authToken, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return authToken, fmt.Errorf("failed to retrieve access token: status code %d\n", resp.StatusCode)
    }

    if err := json.NewDecoder(resp.Body).Decode(&authToken); err != nil {
        return authToken, err
    }

    // c.IndentedJSON(http.StatusOK, authToken)
    return authToken, nil
}