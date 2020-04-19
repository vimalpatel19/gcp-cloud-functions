package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

//FindPosts returns all the posts found
func FindPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	projectID := os.Getenv("PROJECT_ID")
	collectionName := os.Getenv("COLLECTION_NAME")
	size, skip := 0, 0
	var val string
	var err error

	//Get URL query parameters
	val = GetURLParameter(r, "size")
	if val != "" {
		size, err = strconv.Atoi(val)

		if err != nil {
			log.Println("Invalid value provided for 'size' parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	val = GetURLParameter(r, "skip")
	if val != "" {
		skip, err = strconv.Atoi(val)

		if err != nil {
			log.Println("Invalid value provided for 'skip' parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//Connect to the database
	db, err := ConnectToDatabase(ctx, projectID)

	if err != nil {
		log.Println("Error connecting to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Get all posts from the database
	posts, err := db.GetPosts(ctx, collectionName, size, skip)

	if err != nil {
		log.Println("Error retrieving posts from the database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Proceed with returning success response if no errors encountered
	json.NewEncoder(w).Encode(posts)
	w.WriteHeader(http.StatusOK)
}
