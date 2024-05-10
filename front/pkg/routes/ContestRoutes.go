package routes

import (
	"fmt"
	"front/pkg/middle"
	"front/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

func ContestRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", contestPage)
	//r.Post("/", createUser)
	//r.Route("/{id}", func(r chi.Router) {
	//	r.Get("/", getUserByID)
	//	r.Put("/", updateUser)
	//	r.Delete("/", deleteUser)
	//})
	return r
}

func contestPage(w http.ResponseWriter, r *http.Request) {
	//session := middle.GetSession(r)
	currentPlayer := middle.GetCurrentPlayer(r)

	if !middle.CheckAuth(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		//w.Header().Set("HX-Redirect", "/")
		//w.WriteHeader(http.StatusOK)
		return
	}

	// Parse templates
	templates := template.Must(template.ParseFiles(
		"static/templates/layout.html",
		"static/templates/nav.html",
		"static/templates/contest/contest.html"))

	test := r.Context().Value("session").(*sessions.Session)
	fmt.Println("TEST" + test.ID)

	data := types.HeaderInfo{
		Title: "Smacktalk Gaming",
		User: types.HeaderUser{
			Email:     currentPlayer.Email,
			Firstname: currentPlayer.Firstname,
		},
	}
	err := templates.ExecuteTemplate(w, "contest.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
