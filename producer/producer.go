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
	log.Print("Producer received a request.")

	payload := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("in: %+v", payload)

	client := github.NewClient(nil)
	commit, resp, err := client.Repositories.GetCommit(r.Context(), "tektoncd", "triggers", "master")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}
	payload["commit"] = commit

	repo, resp, err := client.Repositories.Get(r.Context(), "tektoncd", "triggers")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}
	payload["repo"] = repo

	payload["build"] = "true"

	//log.Printf("out: %+v", payload)

	enc := json.NewEncoder(w)
	if err := enc.Encode(payload); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	log.Print("Producer started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
