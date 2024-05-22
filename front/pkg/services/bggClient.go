package services

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"front/pkg/types"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"strings"
)

type Items struct {
	XMLName    xml.Name `xml:"items"`
	Total      int      `xml:"total,attr"`
	TermsOfUse string   `xml:"termsofuse,attr"`
	Items      []Item   `xml:"item"`
}

type Item struct {
	Type          string `xml:"type,attr"`
	ID            int    `xml:"id,attr"`
	Name          Name   `xml:"name"`
	YearPublished Year   `xml:"yearpublished"`
}

type Name struct {
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

type Year struct {
	Value int `xml:"value,attr"`
}

func GetGameListFromBGG(query string) (types.Games, error) {
	// Construct the search URL
	url := os.Getenv("BGG_API_URL") + strings.ReplaceAll(query, " ", "+")
	var gamesFound types.Games

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return gamesFound, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gamesFound, fmt.Errorf("failed to read response body: %w", err)
	}

	//convert from XML
	var items Items
	err = xml.Unmarshal([]byte(body), &items)
	if err != nil {
		return gamesFound, fmt.Errorf("error unmarshaling XML: %v", err)
	}

	// Convert to Json
	jsonData, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return gamesFound, fmt.Errorf("error marshaling JSON: %v", err)
	}

	gameData := string(jsonData)
	total := gjson.Get(gameData, "Total")
	//fill up gamesFound, if any
	if total.Int() > 0 {

		items := gjson.Get(gameData, "Items")
		var game types.Game
		items.ForEach(func(_, item gjson.Result) bool {

			game.BGGId = int(item.Get("ID").Int())
			game.Name = item.Get("Name.Value").String()
			game.YearPublished = int(item.Get("YearPublished.Value").Int())

			gamesFound.List = append(gamesFound.List, game)
			return true
		})

	}

	return gamesFound, nil
}
