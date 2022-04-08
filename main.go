package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Login string
}

func checkhealth() error {
	return nil // FIXME: return nil
}

func userFromToken(token string) *User {
	// FIXME: JWT, Oauth2 ...you can user JWT Or Oauth2
	if token == "harit123Pass" { // dummy token "harit123Pass"
		return &User{"Ha-rit Kumsan"} // user what are you want maybe another token or something??
	}
	return nil
}


// Check Serve request
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := checkhealth(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK... Health Check Failed!!\n")
}


// check if user is nil?? when we check http heads is working
func requireAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := authToken(r) // Token are we used maybe not the same this maybe JWT or Oauth2 ...
		fmt.Println("The Token thus received in the request is : ", token)

		user := userFromToken(token)
		fmt.Println("The User Object thus received from Token : ", user)

		// this check
		if user == nil {
			http.Error(w, "bad authentication", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
	
}
// Main
func main() {
	http.HandleFunc("/health", healthHandler)
	h := requireAuth(http.HandlerFunc(messageHandler))
	http.Handle("/messages", h)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// endpoint messageHandler
func messageHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*User)
	if !ok {
		http.Error(w, "no user", http.StatusInternalServerError)
		return
	}
	log.Printf("User thus receive from Context is : %s", user)

	// FIXME: message to frontend when status 200
	fmt.Fprint(w, "[passed!!]/n")
}

// FIXME: Header of request what are http headers, In this is Authorization is a headers of http headers
func authToken(r *http.Request) string {
	hdr := r.Header.Get("Authorization")
	return strings.TrimPrefix(hdr, "Bearer")
}