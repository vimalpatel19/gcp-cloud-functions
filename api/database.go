package api

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

//Database struct for database operationss
type Database struct {
	Client *firestore.Client
}

//ConnectToDatabase attempts to the database returns the database client
func ConnectToDatabase(ctx context.Context, projectID string) (*Database, error) {

	dbClient, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		return nil, err
	}

	return &Database{Client: dbClient}, nil
}

//CountPosts returns the count of posts found in the database
func (db *Database) CountPosts(ctx context.Context, collectionName string) (int, error) {
	defer db.Client.Close()

	//Call database to retrieve all posts
	iter := db.Client.Collection(collectionName).OrderBy("date", firestore.Desc).Documents(ctx)

	docs, err := iter.GetAll()
	if err != nil {
		return 0, err
	}

	return len(docs), nil
}

//GetPosts returns all posts from the database
func (db *Database) GetPosts(ctx context.Context, collectionName string, size, skip int) ([]Post, error) {
	defer db.Client.Close()

	var results []Post = make([]Post, 0)

	//Build query
	query := db.Client.Collection(collectionName).OrderBy("date", firestore.Desc)

	//Update
	if size != 0 && skip != 0 {
		query = db.Client.Collection(collectionName).OrderBy("date", firestore.Desc).Limit(size).Offset(skip)
	} else if size != 0 {
		query = db.Client.Collection(collectionName).OrderBy("date", firestore.Desc).Limit(size)
	} else if skip != 0 {
		query = db.Client.Collection(collectionName).OrderBy("date", firestore.Desc).Offset(skip)
	}

	//Call database to retrieve all posts
	iter := query.Documents(ctx)

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
