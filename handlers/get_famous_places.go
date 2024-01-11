package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RSO-project-Prepih/AI-service/prometheus"
)

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Place struct {
	Properties struct {
		Rate  int    `json:"rate"`
		Name  string `json:"name"`
		Kinds string `json:"kinds"`
	} `json:"properties"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
}

type ResponseData struct {
	Features []Place `json:"features"`
}

func FetchFamousPlaces() ([]Place, error) {
	startTime := time.Now()

	apiKey := os.Getenv("OPENTRIPMAP_API_KEY")
	if apiKey == "" {
		log.Println("OPENTRIPMAP_API_KEY environment variable not set")
		return nil, fmt.Errorf("OPENTRIPMAP_API_KEY environment variable not set")
	}

	lat := "45.7769"
	lon := "14.2166"
	radius := "7000"

	url := fmt.Sprintf("https://opentripmap-places-v1.p.rapidapi.com/en/places/radius?radius=%s&lon=%s&lat=%s", radius, lon, lat)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", "opentripmap-places-v1.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return nil, err
	}

	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("Error parsing JSON response: %s", err)
		return nil, err
	}

	const MIN_POPULARITY_RATING = 6
	var popularPlaces []Place
	for _, place := range responseData.Features {
		if place.Properties.Rate >= MIN_POPULARITY_RATING {
			popularPlaces = append(popularPlaces, place)
		}
	}

	duration := time.Since(startTime)
	prometheus.HTTPRequestDuration.WithLabelValues("opentripmap").Observe(duration.Seconds())

	return popularPlaces, nil
}
