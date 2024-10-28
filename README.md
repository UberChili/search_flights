# search_flights

Simple API that retrieves flight information from the Amadeus API

## Idea
Designed in the Go Programming Language using the Gin framework. This API consults the Self-Servide APIs from Amadeus, retrieves a JSON file and parses it to get:
- Origin: IATA and country codes of origin airport
- Destination: IATA and country codes of destination airport
- Airline
- Flight number
- Price in USD

as output in JSON format.

## Prerequisites
- Go 1.16 or later
- Amadeus API credentials (API Key and API Secret)

### Usage
1. Clone the repository:
```sh
git clone git@github.com:UberChili/search_flights.git
```
2. Set the Amadeus API credentials as environment variables:
```sh
export AMAD_API_KEY=[amadeus-api-key]
export AMAD_API_SECRET=[amadeus-api-secret]
```
3. Build and run the application:
```sh
go build -o search_flights
./search_flights
# or:
go run .
```
The API will be available at `http://localhost:8080`

### API Endpoints
- `GET /flights/:origin/:destination`: Retrieves a list of flights between the specified origin and destination.

Example request: `http://localhost:8080/flights/MEX/MTY`
