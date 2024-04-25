package middle

import (
	"context"
	"github.com/gorilla/sessions"
	"net/http"
)

func SessionMiddleware(store sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "stg-session") // Replace with your session name
			if err != nil {
				// Handle error (e.g. create new session)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer session.Save(r, w) // Save session data after request

			// Add session to request context for access in handlers
			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
