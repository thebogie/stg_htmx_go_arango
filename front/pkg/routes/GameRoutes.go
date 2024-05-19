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

func GameRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/search", search)
	//r.Route("/{id}", func(r chi.Router) {
	//	r.Get("/", getUserByID)
	//	r.Put("/", updateUser)
	//	r.Delete("/", deleteUser)
	//})
	return r
}

func search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//session := r.Context().Value("session").(*sessions.Session)
	templates := template.Must(template.ParseFiles(
		"static/templates/game/gameReturnSearch.html"))

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gql, ok := services.GraphqlClientFromContext(ctx)
	if !ok {
		fmt.Println("Error: the Unexpected graphql client not found")
		return
	}

	query, err := os.ReadFile("pkg/graphql/gamesearch.graphql")
	if err != nil {
		log.Fatal(err)
	}
	variables := map[string]interface{}{
		"name": r.FormValue("game-search"),
	}
	//req.Var("input", variables["input"])
	var result []byte
	result = []byte{}
	err = gql.Query(ctx, string(query), variables, &result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into the LoginUser Graphql struct
	var gamesFound types.FindGameAPI
	err = json.Unmarshal(result, &gamesFound)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	gamesFound.SortByName()

	// Parse the template file
	err = templates.ExecuteTemplate(w, "gameReturnSearch.html", gamesFound)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
