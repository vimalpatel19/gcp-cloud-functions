package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
)

const (
	linkPreviewService = "http://api.linkpreview.net"
)

//Post structure of document inserted into the Firestore database
type Post struct {
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	Image       string    `json:"image" firestore:"image"`
	URL         string    `json:"url" firestore:"url"`
	Date        time.Time `json:"date" firestore:"date"`
}

//SavePost calls Link Preview service using provided URL to start details about the link to the database
func SavePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	linkPreviewKey := os.Getenv("LINK_PREVIEW_KEY")
	projectID := os.Getenv("PROJECT_ID")
	collectionName := os.Getenv("COLLECTION_NAME")

	//Verify the 'link' parameter has been provided as a request parameter
	keys, ok := r.URL.Query()["link"]
	if !ok || len(keys[0]) < 1 {
		log.Println("URL parameter 'link' not provided")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link := keys[0]

	//Call link preview service
	result, err := callLinkPreview(link, linkPreviewKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Insert the result/post into the database
	result.Date = time.Now()
	err = insertPost(ctx, result, projectID, collectionName)
	if err != nil {
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

//insertPost helper function that connects to and stores link/post details into Firestore database
func insertPost(ctx context.Context, post Post, projectID, collectionName string) error {
	//Connect to the database
	dbClient, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	}
	defer dbClient.Close()

	//Add post to the database
	_, _, err = dbClient.Collection(collectionName).Add(ctx, post)

	if err != nil {
		log.Println("Error saving document to the database:", err)
		return err
	}

	return nil
}
