package routes

import (
	"encoding/json"
	"fmt"
	"front/pkg/middle"
	"front/pkg/services"
	"front/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
)

func PlayerRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/login", PlayerLogin)
	r.Post("/auth", PlayerAuth)
	//r.Route("/{id}", func(r chi.Router) {
	//	r.Get("/", getUserByID)
	//	r.Put("/", updateUser)
	//	r.Delete("/", deleteUser)
	//})
	return r
}

func PlayerAuth(w http.ResponseWriter, r *http.Request) {
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
	// Redirect to target URL
	http.Redirect(w, r, "/", http.StatusFound)

}

func PlayerLogin(w http.ResponseWriter, r *http.Request) {
	// Parse templates
	templates := template.Must(template.ParseFiles(
		"static/templates/layout.html",
		"static/templates/index.html",
		"static/templates/login/loginForm.html"))

	currentPlayer := middle.GetCurrentPlayer(r)
	fmt.Println("session" + currentPlayer.Firstname)

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

}
