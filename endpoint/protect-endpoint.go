package endpoint

import (
	"github.com/gorilla/mux"
	"net/http"
)

// AuthMiddleware is a middleware function to check the token.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the request header.
		token := r.Header.Get("Authorization")

		// Perform token validation logic here.
		// You can check the token against a database or some other source.
		// For simplicity, let's just check if the token is "YourSecretToken".
		if token != "YourSecretToken" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		// Token is valid; proceed to the next handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()

	// Use the AuthMiddleware for the "/protected" route.
	r.Handle("/protected", AuthMiddleware(http.HandlerFunc(ProtectedHandler))).Methods("GET")

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}

// ProtectedHandler is a handler for the protected route.
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Access to protected route granted"))
}
