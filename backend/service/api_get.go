package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"net/url"
	"github.com/joho/godotenv"
	"encoding/json"
)

type ApiResponse struct {
	Items    []map[string]interface{} `json:"items"`
	NextPage string                   `json:"next_page,omitempty"`
}

// ApiGet fetches data from the API and returns it in a JSON format
func ApiGet() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api_url := os.Getenv("API_URL")
	bearer_token := "Bearer " + os.Getenv("BEARER_TOKEN")

	var allItems []map[string]interface{}
	nextPage := ""

	for {
		request_url := api_url

		if nextPage != "" {
			parsedUrl, err := url.Parse(api_url)
			if err != nil {
				log.Fatal("Error parsing URL")
			}
			query := parsedUrl.Query()
			query.Set("next_page", nextPage)
			parsedUrl.RawQuery = query.Encode()
			request_url = parsedUrl.String()
		}

		req, err := http.NewRequest("GET", request_url, nil)
		if err != nil {
			log.Fatal("Error creating request")
		}
		req.Header.Add("Authorization", bearer_token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error sending request")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body")
		}

		var apiResponse ApiResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			log.Println("Error unmarshalling response body")
		}

		allItems = append(allItems, apiResponse.Items...)

		if apiResponse.NextPage == "" {
			break
		}
		nextPage = apiResponse.NextPage

		fmt.Printf("Retrieved %d items, fetching next page: %s\n", len(apiResponse.Items), nextPage)
	}
	finalResponse := map[string]interface{}{
		"total_items": len(allItems),
		"items": allItems,
	}

	finalJSON, err := json.MarshalIndent(finalResponse, "", "  ")
	if err != nil {
		log.Println("Error marshalling final response")
	}
	fmt.Println(string(finalJSON))

}
