package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ApiResponse struct {
	Items    []map[string]interface{} `json:"items"`
	NextPage string                   `json:"next_page,omitempty"`
}

// ApiGet fetches data from the API and saves it to the database
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
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error: Status code %d", resp.StatusCode)
			break
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body")
			break
		}

		var apiResponse ApiResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			log.Println("Error unmarshalling response body")
			break
		}

		allItems = append(allItems, apiResponse.Items...)

		if apiResponse.NextPage == "" {
			break
		}
		nextPage = apiResponse.NextPage

		fmt.Printf("Retrieved %d items, fetching next page: %s\n", len(apiResponse.Items), nextPage)
	}
	fmt.Printf("🔄 Starting data processing with %d total items...\n", len(allItems))

	recommendations := make([]RecommendationData, 0, len(allItems))

	fmt.Println("🔄 Converting API data to RecommendationData structs...")
	//Transform API data to RecommendationData structs
	for i, item := range allItems {
		if i%100 == 0 { // Progress every 100 items
			fmt.Printf("�� Processing item %d/%d\n", i, len(allItems))
		}

		rec, err := convertRecommendationData(item)
		if err != nil {
			log.Printf("Error converting recommendation data %d: %v", i, err)
			continue
		}
		recommendations = append(recommendations, rec)
	}

	fmt.Printf("✅ Conversion complete! %d recommendations ready for database\n", len(recommendations))

	fmt.Println("💾 Starting database insertion...")
	//Store to database
	err = SaveRecommendations(recommendations)
	if err != nil {
		log.Printf("❌ Error saving recommendations: %v", err)
	} else {
		fmt.Printf("✅ Database insertion complete! %d recommendations saved\n", len(recommendations))
	}

	fmt.Println("📄 Generating final JSON response...")
	//Log JSON response
	finalReponse := map[string]interface{}{
		"total_items": len(allItems),
		"items":       allItems,
	}
	finalJSON, err := json.MarshalIndent(finalReponse, "", "  ")
	if err != nil {
		log.Println("❌ Error marshalling final response")
	} else {
		fmt.Println("✅ JSON response generated successfully")
	}
	fmt.Println(string(finalJSON))
}

func convertRecommendationData(item map[string]interface{}) (RecommendationData, error) {
	var rec RecommendationData

	// Helper function to get string from interface{}
	getString := func(key string) string {
		if val, ok := item[key]; ok && val != nil {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	rec.Ticker = getString("ticker")
	rec.TargetFrom = getString("target_from")
	rec.TargetTo = getString("target_to")
	rec.Company = getString("company")
	rec.Action = getString("action")
	rec.Brokerage = getString("brokerage")
	rec.RatingFrom = getString("rating_from")
	rec.RatingTo = getString("rating_to")

	//Helper to parse time
	timeStr := getString("time")
	if timeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return rec, fmt.Errorf("failed to parse time: %v", err)
		}
		rec.Time = parsedTime
	} else {
		return rec, fmt.Errorf("time is required")
	}
	return rec, nil
}
