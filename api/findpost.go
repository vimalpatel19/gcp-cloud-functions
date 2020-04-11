package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

//FindPosts returns all the posts available
func FindPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	projectID := os.Getenv("PROJECT_ID")
	collectionName := os.Getenv("COLLECTION_NAME")

	//Call database to get all posts
	posts, err := getPosts(ctx, projectID, collectionName)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Proceed with returning success response if no errors encountered
	json.NewEncoder(w).Encode(posts)
	w.WriteHeader(http.StatusOK)
}

//getPosts helper function that connects to the database and makes call to return all posts
func getPosts(ctx context.Context, projectID, collectionName string) ([]Post, error) {
	var results []Post = make([]Post, 0)

	//Connect to the database
	dbClient, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Println("Error connecting to database:", err)
		return results, err
	}
	defer dbClient.Close()

	//Call database to retrieve all posts
	iter := dbClient.Collection(collectionName).Documents(ctx)

	//Iterate through the list of posts returned
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return results, err
		}

		//Decode the post and add to the list of results
		var post Post
		doc.DataTo(&post)
		results = append(results, post)
	}

	return results, nil
}
