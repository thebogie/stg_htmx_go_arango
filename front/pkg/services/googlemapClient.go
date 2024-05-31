package services

import (
	"encoding/json"
	"fmt"
	"front/pkg/types"
	"io"
	"net/http"
	"os"
)

type Prediction struct {
	Description string `json:"description"`
	Place_ID    string `json:"place_id"`
}

type GoogleMapResponse struct {
	Predictions []Prediction `json:"predictions"`
}

func GetAddressListFromGoogleMaps(query string) (types.Venues, error) {
	// Construct the search URL
	url := os.Getenv("GOOGLEMAP_API_URL") + "?input=" + query + "&types=geocode&key=" + os.Getenv("GOOGLE_LOCATION_API")

	var venuesFound types.Venues

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return venuesFound, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return venuesFound, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response
	var response GoogleMapResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return venuesFound, fmt.Errorf("failed to read response body: %w", err)
	}

	for _, prediction := range response.Predictions {
		var venue types.Venue

		venue.PlaceID = prediction.Place_ID
		venue.Address = prediction.Description

		venuesFound.List = append(venuesFound.List, venue)
	}

	//foundData := string(jsonData)
	//total := gjson.Get(foundData, "Total")
	//fill up gamesFound, if any
	/*if total.Int() > 0 {

		items := gjson.Get(foundData, "Items")
		var venue types.Venue
		items.ForEach(func(_, item gjson.Result) bool {

			//venue.BGGId = int(item.Get("ID").Int())
			//venue.Name = item.Get("Name.Value").String()
			//venue.YearPublished = int(item.Get("YearPublished.Value").Int())

			venuesFound.List = append(venuesFound.List, venue)
			return true
		})

	} */

	return venuesFound, nil
}
