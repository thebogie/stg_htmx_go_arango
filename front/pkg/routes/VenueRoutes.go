package routes

import (
	"encoding/json"
	"fmt"
	"front/pkg/services"
	"front/pkg/types"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"os"
)

func VenueRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/search", venueSearch)
	//r.Route("/{id}", func(r chi.Router) {
	//	r.Get("/", getUserByID)
	//	r.Put("/", updateUser)
	//	r.Delete("/", deleteUser)
	//})
	return r
}

func venueSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var venuesFound types.Venues
	//session := r.Context().Value("session").(*sessions.Session)
	templates := template.Must(template.ParseFiles(
		"static/templates/Venue/venueReturnSearch.html"))

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	searchQuery := r.FormValue("venue_search")

	gql, ok := services.GraphqlClientFromContext(ctx)
	if !ok {
		fmt.Println("Error: the Unexpected graphql client not found")
		return
	}

	query, err := os.ReadFile("pkg/graphql/venuesearch.graphql")
	if err != nil {
		log.Fatal(err)
	}
	variables := map[string]interface{}{
		"name": searchQuery,
	}
	//req.Var("input", variables["input"])
	var result []byte
	result = []byte{}
	err = gql.Query(ctx, string(query), variables, &result)
	if err != nil {
		gql.CheckLoginRefresh(w, r, err)
		return
	}

	var venuesFoundAPI types.FindVenueAPI
	var venuesLocalFound types.Venues
	var venuesGoogleFound types.Venues

	err = json.Unmarshal(result, &venuesFoundAPI)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	var markedVenues []types.Venue // Create an empty slice to store modified games
	for _, venue := range venuesFoundAPI.FindVenue {
		markedVenues = append(markedVenues, types.Venue{
			InDatabase: true,
			Id:         venue.Id,
			Address:    venue.Address,
			PlaceID:    venue.PlaceID,
			Lat:        venue.Lat,
			Lng:        venue.Lng,
		})
	}
	venuesLocalFound.List = markedVenues

	venuesGoogleFound, err = services.GetAddressListFromGoogleMaps(searchQuery)
	if err != nil {
		fmt.Println("Error:googlemap api issue")
		return
	}

	venuesLocalFound.SortByAddress()
	venuesGoogleFound.SortByAddress()
	venuesFound.List = append(venuesLocalFound.List, venuesGoogleFound.List...)

	// Parse the template file
	err = templates.ExecuteTemplate(w, "venueReturnSearch.html", venuesFound)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
