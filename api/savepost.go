package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	linkPreviewService = "http://api.linkpreview.net"
)

//SavePost calls Link Preview service using provided URL to start details about the link to the database
func SavePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	linkPreviewKey := os.Getenv("LINK_PREVIEW_KEY")
	projectID := os.Getenv("PROJECT_ID")
	collectionName := os.Getenv("COLLECTION_NAME")

	//Verify the 'link' parameter has been provided as a request parameter
	link := GetURLParameter(r, "link")
	if link == "" {
		log.Println("Invalid value provided for 'link' parameter")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Call link preview service
	result, err := callLinkPreview(link, linkPreviewKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	//Insert the result/post into the database
	result.Date = time.Now()
	result.Likes = 0
	err = db.InsertPost(ctx, result, collectionName)

	if err != nil {
		log.Println("Error inserting post into the database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Proceed with returning response if no errors encountered
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusCreated)
}

//callLinkPreview helper function that calls and returns response from Link Preview service
func callLinkPreview(link, key string) (Post, error) {
	var result Post

	//Prepare request to retrieve the preview for the given link
	client := http.Client{}
	request, _ := http.NewRequest("GET", linkPreviewService, nil)

	//Add the request parameters
	q := request.URL.Query()
	q.Add("key", key)
	q.Add("q", link)
	request.URL.RawQuery = q.Encode()

	//Make the call
	resp, err := client.Do(request)

	if err != nil {
		log.Println("Error calling Link Preview:", err)
		return result, err
	}

	defer resp.Body.Close()

	//Decode the response body
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
