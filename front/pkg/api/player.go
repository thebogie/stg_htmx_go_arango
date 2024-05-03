package api

import (
	"encoding/json"
	"fmt"
	"front/pkg/services"
	"front/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
)

// NewAPIRouter creates a sub-router for API endpoints
func NewAPIRouter() chi.Router {
	// Create a sub-router
	router := chi.NewRouter()

	// Add middleware (optional)
	router.Use(middleware.Logger) // Logs API requests

	// Parse templates
	templates := template.Must(template.ParseFiles(
		"static/templates/login/loginForm.html",
		"static/templates/login/auth.html"))

	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {

		session := r.Context().Value("session").(*sessions.Session)
		currentPlayer := session.Values["currentPlayer"]

		_, ok := currentPlayer.(*types.Player)
		if currentPlayer != nil && !ok {
			fmt.Println("Error: Unexpected value for currentPlayer")
			return
		}

		data := struct {
			Title string
		}{
			Title: "Login",
		}
		err := templates.ExecuteTemplate(w, "loginForm.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	router.Post("/auth", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session := r.Context().Value("session").(*sessions.Session)

		// "input": {    "username": "mitch@gmail.com",    "password": "letmein"  }
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		gql, ok := services.GraphqlClientFromContext(ctx)
		if !ok {
			fmt.Println("Error: Unexpected graphql client not found")
			return
		}

		query, err := os.ReadFile("pkg/graphql/login.graphql")
		if err != nil {
			log.Fatal(err)
		}
		variables := map[string]interface{}{
			"input": map[string]interface{}{
				"username": r.FormValue("email"),
				"password": r.FormValue("password"),
			},
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
		var loginUser types.LoginUserAPI
		err = json.Unmarshal(result, &loginUser)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		session.Values["currentPlayer"] = types.Player{
			Firstname:   loginUser.LoginUser.Userdata.Firstname,
			Email:       loginUser.LoginUser.Userdata.Email,
			Password:    "",
			AccessToken: loginUser.LoginUser.Token,
		}
		err = session.Save(r, w)
		if err != nil {
			panic(err)
			return
		}

		// Parse templates
		templates := template.Must(template.ParseFiles("static/templates/layout.html", "static/templates/index.html"))
		data := types.HeaderInfo{
			Title: "Smacktalk Gaming",
			User: types.HeaderUser{
				Email:     loginUser.LoginUser.Userdata.Email,
				Firstname: "",
			},
		}
		err = templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	return router
}
