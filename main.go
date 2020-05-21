package main

import (
	subreddit "RedditGetter/Subreddit"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// S handles the subreddit struct
var S subreddit.Subreddit

func updateReddits(w http.ResponseWriter, r *http.Request) {
	var list []string
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Please check format")
	}

	json.Unmarshal(reqBody, &list)
	S.UpdateSubbredditList(list)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(S.GetAllSubreddits())
}
func getPost(w http.ResponseWriter, r *http.Request) {
	sub := mux.Vars(r)["id"]

	url := S.GetRandomPost(sub)

	fmt.Fprint(w, url)
}

func getAllSubreddits(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(S.GetAllSubreddits())
}



func main() {
	S.Init()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/subreddit", updateReddits).Methods("POST")
	router.HandleFunc("/subreddit/{id}", getPost).Methods("GET")
	router.HandleFunc("/subreddits",getAllSubreddits).Methods("GET")

	fmt.Println("Done")

	log.Fatal(http.ListenAndServe(":8080", router))

}




