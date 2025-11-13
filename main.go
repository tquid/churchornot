package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/joho/godotenv"

	"net/http"
	"os"
)

type AddressQuery struct {
	TextQuery string `json:"textQuery"`
}

// Load API key from .env file

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getPlacesInfo(apiKey string, address string) string {
	p := AddressQuery{
		TextQuery: address,
	}

	jsonData, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Error marshaling request body: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://places.googleapis.com/v1/places:searchText", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating http request: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Goog-Api-Key", apiKey)
	req.Header.Add("X-Goog-FieldMask", "places.id,places.types,places.primaryType")

	fmt.Printf("Request body: %s\n", string(jsonData))
	fmt.Printf("API Key (first 10 chars): %s...\n", apiKey[:10])

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		errorBody, _ := io.ReadAll(resp.Body)
		log.Fatalf("Server returned error: %d %s - %s", resp.StatusCode, resp.Status, string(errorBody))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}

	return string(body)
}

func main() {
	// Grab our CSV file with address info

	apiKey := os.Getenv("APIKEY")

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Error trying to open '%s': %s", os.Args[1], err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error trying to read CSV file: %s", err)
	}
	for i := 1; i <= 10; i++ {
		fmt.Printf("%s\n", records[i])
		addressInfo := getPlacesInfo(apiKey, strings.Join(records[i], " "))
		fmt.Printf("Reply: %s", addressInfo)
	}
}
