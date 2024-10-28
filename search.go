package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Main, "parent" response
type ApiResponse struct {
    Data []Flight `json:"data"`
    Dictionaries DictionaryInfo `json:"dictionaries"`
}

// Flight contains information about a flight, it depends on other structs
type Flight struct {
    Id string `json:"id"`
    Itineraries []struct {
        Segments []struct {
            Departure struct {
                IataCode string `json:"iataCode"`
            } `json:"departure"`
            Arrival struct {
                IataCode string `json:"iataCode"`
            } `json:"arrival"`
            CarrierCode string `json:"carrierCode"`
            Number      string `json:"number"`
        } `json:"segments"`
    } `json:"itineraries"`
    Price Price `json:"price"`
}

// Price contains the actual field we want (total)
type Price struct {
    Currency string `json:"currency"`
    Total string `json:"total"`
}

type DictionaryInfo struct {
    Locations map[string]Location `json:"locations"`
    Carriers map[string]string  `json:"carriers"`
}

type Location struct {
    CityCode string `json:"cityCode"`
    CountryCode string `json:"countryCode"`
}

// Auxiliary struct to nicely format our output
type SimplifiedFlight struct {
    Id string `json:"id"`
    Origin struct {
        Code    string `json:"code"`
        Country string `json:"country"`
    } `json:"origin"`
    Destination struct {
        Code    string `json:"code"`
        Country string `json:"country"`
    } `json:"destination"`
    Airline      string `json:"airline"`
    FlightNumber string `json:"flightNumber"`
    Price        Price  `json:"price"`
}

const (
    // baseURL = "https://test.api.amadeus.com/v1/shopping/flight-destinations"
    // Was wrongly using the above url for hours, making everything more difficult. The (more) correct URL is the following:
    baseURL = "https://test.api.amadeus.com/v2/shopping/flight-offers"
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
    destination := c.Param("destination")

    // Get tomorrow's date in the required format (YYYY-MM-DD)
    tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

    resp, err := makeAmadeusRequest(authToken, origin, destination, tomorrow)
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

    // This responds with an intended JSON, which formats some characters weirdly, like tha ampersand
    // It also is a little more resource heavy. Ideally we should use just c.JSON()
    c.IndentedJSON(http.StatusOK, simplifyFlights(response))
}

// Neatly formats our response to get a new slice of simplified flights
func simplifyFlights(flights ApiResponse) []SimplifiedFlight {
    simplifiedFlights := make([]SimplifiedFlight, 0)

    for _, flight := range flights.Data {
        segment := flight.Itineraries[0].Segments[0]
        simplified := SimplifiedFlight {
            Price: flight.Price,
            Id: flight.Id,
        }

        // Set origin
        simplified.Origin.Code = segment.Departure.IataCode
        simplified.Origin.Country = flights.Dictionaries.Locations[segment.Departure.IataCode].CountryCode
        // Set destination
        simplified.Destination.Code = segment.Arrival.IataCode
        simplified.Destination.Country = flights.Dictionaries.Locations[segment.Arrival.IataCode].CountryCode
        // Set airline and flight number
        simplified.Airline = flights.Dictionaries.Carriers[segment.CarrierCode]
        simplified.FlightNumber = segment.CarrierCode + segment.Number

        simplifiedFlights = append(simplifiedFlights, simplified)
    }
    return simplifiedFlights
}

// Helper function that performs the actual calls to the API
func makeAmadeusRequest(authToken AuthToken, origin string, destination string, tomorrow string) (*http.Response, error) {
    req, err := http.NewRequest("GET", fmt.Sprintf(
        "%s?originLocationCode=%s&destinationLocationCode=%s&departureDate=%s&adults=1&currencyCode=%s",
        baseURL, origin, destination, tomorrow, "USD"), nil)
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
