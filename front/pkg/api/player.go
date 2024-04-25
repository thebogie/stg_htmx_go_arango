package api

import (
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
		"static/templates/login/login.html",
		"static/templates/login/auth.html"))

	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {

		session := r.Context().Value("session").(*sessions.Session)
		currentPlayerI := session.Values["currentPlayer"]

		player, ok := currentPlayerI.(types.Player)
		if !ok {
			fmt.Println("Error: Unexpected value for currentPlayer")
			return
		}

		fmt.Println("Player Email:", player.Email)
		fmt.Println("Player Firstname:", player.Firstname)

		data := struct {
			Title string
		}{
			Title: "Login",
		}
		err := templates.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	router.Post("/auth", func(w http.ResponseWriter, r *http.Request) {
		//session := r.Context().Value("session").(*sessions.Session)
		//currentPlayerI := session.Values["currentPlayer"]

		query, err := os.ReadFile("pkg/graphql/login.graphql")
		if err != nil {
			log.Fatal(err)
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("POST DATA: " + r.FormValue("email"))

		//get the motts
		services.Graphql_(string(query))

		data := struct {
			Title string
		}{
			Title: "Auth",
		}
		err = templates.ExecuteTemplate(w, "auth.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	return router
}
