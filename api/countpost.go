package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

//CountPosts returns the number of posts found
func CountPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	projectID := os.Getenv("PROJECT_ID")
	collectionName := os.Getenv("COLLECTION_NAME")

	//Connect to the database
	db, err := ConnectToDatabase(ctx, projectID)

	if err != nil {
		log.Println("Error connecting to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Return count of posts found
	count, err := db.CountPosts(ctx, collectionName)

	if err != nil {
		log.Println("Error retrieving count from the database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Proceed with returning success response if no errors encountered
	json.NewEncoder(w).Encode(CountResponse{Count: count})
	w.WriteHeader(http.StatusOK)
}
