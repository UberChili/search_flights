package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiResponse represents the top-level response from Amadeus API
type ApiResponse struct {
  Data []Flight `json:"data"`
}

// flight represents data about a flight
type Flight struct {
    Origin string `json:"origin"`
    Destination string `json:"destination"`
    DepartureDate string `json:"departureDate"`
    Price Price `json:"price"`
    Links Links `json:"links"`
}

// Links represents the two links contained inside Data (Flight).
// They have important information about each flight
type Links struct {
    FlightDates string `json:"flightDates"`
    FlightOffers string `json:"flightOffers"`
}

// Price represents the price information of a flight
type Price struct {
    Total string `json:"total"`
}

// FormattedFlight represents data about a flight to respond with
type FormattedFlight struct {
    Origin string `json:"origin"`
    Destination string `json:"destination"`
    Price string `json:"price"`
    FlightOffers string `json:flightOffers`
}

const (
    baseURL = "https://test.api.amadeus.com/v1/shopping/flight-destinations"
)

// getFlights responds with the list of all flights from a destination as JSON
func getFlights(c *gin.Context) {
    // Need AuthToken authorization
    auth_token, err := getAuthToken()
    if err != nil {
        fmt.Errorf("Error getting Auth Token: %d\n", err)
        return
    }

    // origin is obtained from query
    origin := c.Param("origin")

    // Create the request
    req, err := http. NewRequest("GET", fmt.Sprintf("%s?origin=%s&maxPrice=200", baseURL, origin), nil)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        // Should generally use c.JSON() in most cases as it is less resources heavy
        // c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Add headers 
    var full_access_token = "Bearer " + auth_token.AccessToken
    req.Header.Add("Authorization", full_access_token)

    // Make the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer resp.Body.Close()

    // Read response Body 
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Check if the response was successful
    if resp.StatusCode != http.StatusOK {
        c.IndentedJSON(resp.StatusCode, gin.H{"error": string(body)})
        return
    }

    // Parse the JSON response
    var response ApiResponse
    if err := json.Unmarshal(body, &response); err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response: " + err.Error()})
        return
    }

    // Prepare JSON to show
    var flights = []FormattedFlight{}
    var newflight FormattedFlight 

    for _, value := range response.Data {
        newflight.Origin = value.Origin
        newflight.Destination = value.Destination
        newflight.Price = value.Price.Total
        newflight.FlightOffers = value.Links.FlightOffers

        // Need the additional info 
        // get_info()

        flights = append(flights, newflight)
    }

    // c.IndentedJSON(http.StatusFound, response.Data)
    c.IndentedJSON(http.StatusFound, flights)
}

func getInfo(c *gin.Context) {
}
