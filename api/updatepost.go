package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

//UpdateLikes updates the number of likes for the given post
func UpdateLikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	projectID := os.Getenv("PROJECT_ID")
	collectionName := os.Getenv("COLLECTION_NAME")
	updateLikes := 0
	var err error

	//Get URL query parameters
	id := GetURLParameter(r, "id")
	if id == "" {
		log.Println("Invalid value provided for 'id' parameter")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val := GetURLParameter(r, "count")
	if val != "" {
		updateLikes, err = strconv.Atoi(val)

		if err != nil {
			log.Println("Invalid value provided for 'count' parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		log.Println("No value provided for 'count' parameter")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Connect to the database
	db, err := ConnectToDatabase(ctx, projectID)
	defer db.Client.Close()

	if err != nil {
		log.Println("Error connecting to database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Update the post
	err = db.UpdateLikes(ctx, updateLikes, id, collectionName)

	if err != nil {
		log.Println("Error updating like count in database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Get the updated like count
	updated, err := db.GetPostByID(ctx, collectionName, id)

	if err != nil {
		log.Println("Error validating updated count in database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Proceed with returning success response if no errors encountered
	json.NewEncoder(w).Encode(CountResponse{Count: updated.Likes})
	w.WriteHeader(http.StatusOK)
}
