package middle

import (
	"context"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"front/pkg/services"
	"front/pkg/types"
	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Store is the session store
var Store *redistore.RediStore

// Init initializes the session store
func InitSession() {
	gob.Register(&types.Player{})

	var err error

	for i := 0; i < 10; i++ {
		Store, err = redistore.NewRediStore(10, "tcp", os.Getenv("STG_REDIS_LOCATION"), "redispassword", []byte(os.Getenv("STG_REDIS_KEY")))
		if err == nil {
			// RediStore created successfully, break out of the loop
			break
		}

		// Wait for 1 second before trying again
		time.Sleep(1 * time.Second)
	}

	// Check if the RediStore was created successfully
	if err != nil {
		panic(err)
	} else {
		fmt.Println("RediStore created successfully!")
	}

}

// Middleware is the session middleware

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := Store.Get(r, "session-id")
		if err != nil {
			// Handle the error
			return
		}

		// Check if the session is new
		if session.IsNew || session.Values["currentPlayer"] == nil {
			// Generate a new session ID
			//session.Values["user_id"] = generateUniqueID()
			session.Values["currentPlayer"] = types.Player{
				Firstname:   "empty",
				Email:       "empty",
				Password:    "",
				AccessToken: "",
			}
			err := session.Save(r, w)
			if err != nil {
				panic(err)
				return
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "session", session)
		ctx = services.WithGraphQLClient(ctx, services.InitGraphQL(os.Getenv("STG_GRAPHQL_API")))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetSession retrieves the session from the request context
func GetSession(r *http.Request) *sessions.Session {
	return r.Context().Value("session").(*sessions.Session)
}

// GetUserID retrieves the user ID from the request context
func GetCurrentPlayer(r *http.Request) *types.Player {
	retVal := &types.Player{}
	session := r.Context().Value("session").(*sessions.Session)
	currentPlayer := session.Values["currentPlayer"]

	_, ok := currentPlayer.(*types.Player)
	if currentPlayer != nil && !ok {
		fmt.Errorf("Unexpected value for currentPlayer: %T", currentPlayer)

	}

	retVal = currentPlayer.(*types.Player)
	return retVal
}

// getUserID generates a unique user ID for the request
func generateUniqueID(r *http.Request) string {

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic("Rand didnt work in generateID???")
	}

	// Encode the random sequence as a URL-safe base64 string
	return base64.URLEncoding.EncodeToString(b)
}

const (
	sessionName       = "stg_session"
	sessionContextKey = "stg_session_store"
)
