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

// flight represents data about a flight gotten from the API
// This worked
// type Flight struct {
//     Origin string `json:"origin"`
//     Destination string `json:"destination"`
//     DepartureDate string `json:"departureDate"`
//     Price Price `json:"price"`
//     Links Links `json:"links"`
// }

type Flight struct {
    Origin string `json:"origin"`
    Destination string `json:"destination"`
    DepartureDate string `json:"departureDate"`
    Price Price `json:"price"`
    Links Links `json:"links"`
    // FlightDates string `json:"flightDates"`
    // FlightOffers string `json:"flightOffers"`
}

// Helper that represents data about a flight to respond with
// type FormattedFlight struct {
//     Origin string `json:"origin"`
//     Destination string `json:"destination"`
//     Price string `json:"price"`
//     FlightDates string `json:"flightDates"`
//     FlightOffers string `json:"flightOffers"`
// }

// Price contains the actual field we want (total)
type Price struct {
    Total string `json:"total"`
}

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

    // formattedFlights := formatFlights(response.Data)
    // c.IndentedJSON(http.StatusOK, formattedFlights)

    // trying to deserialize directly
    actualFlights := formatFlights(response.Data)
    c.IndentedJSON(http.StatusOK, actualFlights)
}

func formatFlights(flights []Flight) []Flight {
    actual_flights := make([]Flight, len(flights))
    for i, flight := range flights {
        actual_flights[i] = Flight{
            Origin: flight.Origin,
            Destination: flight.Destination,
            Price: flight.Price,
            Links: flight.Links,
        }
    }
    return actual_flights
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

// Formats the flights structs
// func formatFlights(flights []Flight) []FormattedFlight {
//     formatted := make([]FormattedFlight, len(flights))
//     for i, flight := range flights {
//         formatted[i] = FormattedFlight{
//             Origin: flight.Origin,
//             Price: flight.Price.Total,
//             FlightDates: flight.Links.FlightDates,
//             FlightOffers: flight.Links.FlightOffers,
//         }
//     }
//     return formatted
// }

// Helps avoid repetitive error-handling code
func respondWithError(c *gin.Context, status int, message string) {
    c.IndentedJSON(status, gin.H{"error": message})
}
