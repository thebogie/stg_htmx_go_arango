package main

import (
	"encoding/gob"
	"fmt"
	"front/pkg/api"
	"front/pkg/middle"
	"front/pkg/types"
	"front/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var secretKey = []byte(utils.GetConfig().GetString("STG_SESSIONS_KEY")) // Replace with a strong, random key
	store := sessions.NewCookieStore(secretKey)

	// Create a new Chi router
	r := chi.NewRouter()

	// Add stg-middleware (optional)
	r.Use(middleware.Logger) // Logs requests to the console
	r.Use(middle.SessionMiddleware(store))
	// Create a route along /files that will serve contents from
	// the ./static/ folder.
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/static", filesDir)

	// Parse templates
	templates := template.Must(template.ParseFiles("static/templates/layout.html", "static/templates/index.html"))

	// Define routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("session").(*sessions.Session)

		player := types.Player{Firstname: "", Email: "FISH", Password: "", AccessToken: ""}
		gob.Register(types.Player{}) // Register Player struct locally
		session.Values["currentPlayer"] = player

		err := session.Save(r, w)
		if err != nil {
			// Handle error appropriately (e.g., log the error or return an error to the client)
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
		}{
			Title: "My Front Page",
		}
		err = templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	// Mount the API handler (assuming a handler function in pkg/api)
	r.Mount("/player", api.NewAPIRouter()) // Mount under the "/api" prefix

	// Start the server on port 8080 (or any desired port)
	fmt.Println("Server listening on port 50003")
	http.ListenAndServe(":50003", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
