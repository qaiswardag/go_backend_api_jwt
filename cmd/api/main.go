package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/qaiswardag/go_backend_api_jwt/internal/pkg/support"
	"github.com/qaiswardag/go_backend_api_jwt/pkg/httpResponseMessages"
)

type Handler struct{}

// func authorize(r *http.Request) {
// 	st, err := r.Cookie("session_token")

// 	if err != nil {
// 		fmt.Println("No session token.")
// 	}
// 	fmt.Println("Session token is:", st)
// }

func login(r *http.Request, w http.ResponseWriter) {
	if r.URL.Path == "/login" && r.Method == http.MethodPost {
		sessionToken := support.GenerateToken(32)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,
		})
		// Store session_token token in database

		csrfToken := support.GenerateToken(32)
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,
		})
		// Store csrf_token token in database

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(httpResponseMessages.GetSuccessResponse())
	}
}

func getSensitiveData(r *http.Request, w http.ResponseWriter, tokenName string) {

	if r.URL.Path == "/sensitive-data" {
		// Attempt to retrieve the cookie
		cookie, err := r.Cookie(tokenName)

		if err != nil {
			// Handle the case where the cookie is not found or other errors occur
			fmt.Println("err is not nil:", err)
			http.Error(w, "Unauthorized: session token missing", http.StatusUnauthorized)
			return
		}

		// Check if the cookie value is empty
		if cookie.Value == "" {
			fmt.Printf("Cookie %s is empty: %+v\n", cookie.Name, cookie)
			http.Error(w, "Unauthorized: session token is empty", http.StatusUnauthorized)
			return
		}

		// Log the cookie name and value
		fmt.Printf("Token Name: %s, Token Value: %s\n\n", cookie.Name, cookie.Value)
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:7777")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept-Version")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Log the request method and URL path
	fmt.Printf("New:\nIncoming request: %s %s\n\n", r.Method, r.URL.Path)

	// Handle preflight request
	// GET requests don't trigger a preflight OPTIONS request, so the handler is called only once.
	// Post requests first trigger a preflight OPTIONS request, so the handler is called only twice.
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	login(r, w)

	// getSensitiveData(r, w, "session_token")
	// getSensitiveData(r, w, "csrf_token")

	if r.URL.Path != "/sensitive-data" && r.URL.Path != "/login" && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(httpResponseMessages.GetErrorResponse())
	}

}

func main() {
	handler := Handler{}

	server := http.Server{
		Addr:    ":5555",
		Handler: handler,
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("Server failed to start: ", err)
	}
}
