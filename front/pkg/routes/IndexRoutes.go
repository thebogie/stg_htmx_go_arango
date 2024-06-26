package routes

import (
	"front/pkg/middle"
	"front/pkg/types"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

func IndexRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", mainPage)
	//r.Post("/", createUser)
	//r.Route("/{id}", func(r chi.Router) {
	//	r.Get("/", getUserByID)
	//	r.Put("/", updateUser)
	//	r.Delete("/", deleteUser)
	//})
	return r
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	//session := middle.GetSession(r)
	currentPlayer := middle.GetCurrentPlayer(r)

	// Parse templates
	templates := template.Must(template.ParseFiles(
		"static/templates/layout.html",
		"static/templates/nav.html",
		"static/templates/index.html"))

	data := types.HeaderInfo{
		Title: "Smacktalk Gaming",
		User: types.HeaderUser{
			Email:     currentPlayer.Email,
			Firstname: currentPlayer.Firstname,
		},
	}
	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
