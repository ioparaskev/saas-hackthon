package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type City struct {
	Name    string
	Country string
}

type GeoCodingResponse struct {
	GeoCodingResults []GeoCodingResult `json:"results"`
}

type GeoCodingResult struct {
	Name      string      `json:"name"`
	Latitude  json.Number `json:"latitude"`
	Longitude json.Number `json:"longitude"`
	Country   string      `json:"country"`
}

type Location struct {
	Longitude string
	Latitude  string
}

type Forecast struct {
	Temperature string
	Location    Location
}

func GetWeather(loc Location) ([]byte, error) {
	formattedURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=temperature_2m_max,temperature_2m_min",
		loc.Latitude,
		loc.Longitude,
	)
	response, err := http.Get(formattedURL)

	if err != nil {
		return []byte{}, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return responseData, nil
}

func FindCityLocation(city City) (string, string, error) {
	api_params := url.PathEscape(fmt.Sprintf("name=%s&count=10&language=en&format=json", city.Name))
	api_url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?%s", api_params)
	response, err := http.Get(api_url)

	if err != nil {
		return "", "", err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	var geocodingResponse GeoCodingResponse
	err = json.Unmarshal(responseData, &geocodingResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	for i := 0; i < len(geocodingResponse.GeoCodingResults); i++ {
		if geocodingResponse.GeoCodingResults[i].Country == city.Country {
			log.Println("Matched country location: ", geocodingResponse.GeoCodingResults[i])
			return geocodingResponse.GeoCodingResults[i].Latitude.String(), geocodingResponse.GeoCodingResults[i].Longitude.String(), nil
		}
	}

	return "", "", fmt.Errorf("Could not find a proper location match for %s of country %s", city.Name, city.Country)
}

func main() {
	// TODO: Take input
	city := flag.String("city", "", "Name of the city")
	country := flag.String("country", "", "Country of the city")
	day := flag.String("day", "", "Day for the weather forecast (e.g., '2024-10-31')")

	flag.Usage = func() {
		fmt.Println("Weather CLI Tool")
		fmt.Println("Fetches weather information for a specified city and date.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run main.go -city=\"CityName\" -country=\"CountryName\" -day=\"YYYY-MM-DD\"")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) == 1 || *city == "" || *country == "" || *day == "" {
		flag.Usage()
		os.Exit(1)
	}

	// TODO: ask weather
	lat, lon, err := FindCityLocation(City{Name: *city, Country: *country})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loc := Location{
		Latitude:  lat,
		Longitude: lon,
	}

	weather, err := GetWeather(loc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: Print weather
	fmt.Printf("Weather in %s, %s on %s: %s\n", *city, *country, *day, string(weather))
}
