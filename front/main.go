package main

import (
	"fmt"
	"front/pkg/middle"
	"front/pkg/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the session store
	middle.InitSession()

	// Create a new Chi router
	r := chi.NewRouter()
	//ctx := context.Background()
	//ctx = services.WithRedisClient(ctx, redisClient)
	//ctx = services.WithGraphQLClient(ctx, graphqlClient)
	// Add stg-middleware (optional)
	r.Use(middleware.Logger) // Logs requests to the console
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middle.SessionMiddleware)

	// Create a route along /files that will serve contents from
	// the ./static/ folder.
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/static", filesDir)

	// Define routes
	r.Mount("/", routes.IndexRoutes())
	r.Mount("/profile", routes.ProfileRoutes())
	r.Mount("/player", routes.PlayerRoutes())
	r.Mount("/contest", routes.ContestRoutes())
	r.Mount("/game", routes.GameRoutes())

	// Mount the API handler (assuming a handler function in pkg/api)
	//r.Mount("/player", api.NewAPIRouter()) // Mount under the "/api" prefix

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
