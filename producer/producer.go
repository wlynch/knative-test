package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")

	client := github.NewClient(nil)
	pr, resp, err := client.PullRequests.Get(r.Context(), "tektoncd", "triggers", 1)
	if err != nil {
		http.Error(w, err.Error(), resp.StatusCode)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(pr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	log.Print("Hello world sample started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
