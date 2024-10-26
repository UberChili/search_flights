package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
  Data []Flight `json:"data"`
}

// Flight contains main information about a Flight
type Flight struct {
    Origin string `json:"origin"`
    Destination string `json:"destination"`
    Price Price `json:"price"`
    Links Links `json:"links"`
}

// Price contains the actual field we want (total)
type Price struct {
    Total string `json:"total"`
}

// Links we need to call in order to get additional information 
type Links struct {
    FlightDates string `json:"flightDates"`
    FlightOffers string `json:"flightOffers"`
}

const (
    baseURL = "https://test.api.amadeus.com/v1/shopping/flight-destinations"
)

// getFlights responds with the list of all flights from a destination as JSON
func getFlights(c *gin.Context) {
    authToken, err := getAuthToken()
    if err != nil {
        respondWithError(c, http.StatusUnauthorized, "Error getting Auth Token")
        return
    }

    // Use origin from the query and make request
    origin := c.Param("origin")
    resp, err := makeAmadeusRequest(authToken, origin)
    if err != nil {
        respondWithError(c, http.StatusInternalServerError, "Error making request to Amadeus API")
        return
    }
    defer resp.Body.Close()

    var response ApiResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        respondWithError(c, http.StatusInternalServerError, "Failed to parse API response")
        return
    }

    c.IndentedJSON(http.StatusOK, response)
}

// Actually calls the API
func makeAmadeusRequest(authToken AuthToken, origin string) (*http.Response, error) {
    req, err := http.NewRequest("GET", fmt.Sprintf("%s?origin=%s&maxPrice=200", baseURL, origin), nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer " + authToken.AccessToken)
    client := &http.Client{}
    return client.Do(req)
}

// Helps avoid repetitive error-handling code
func respondWithError(c *gin.Context, status int, message string) {
    c.IndentedJSON(status, gin.H{"error": message})
}
