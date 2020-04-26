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
		post.ID = doc.Ref.ID

		results = append(results, post)
	}

	return results, nil
}

//GetPostByID returns all posts matching the given id
func (db *Database) GetPostByID(ctx context.Context, collectionName, id string) (Post, error) {
	var post Post

	doc, err := db.Client.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		return post, err
	}

	doc.DataTo(&post)
	post.ID = doc.Ref.ID

	return post, nil
}

//InsertPost inserts the provided post into the database
func (db *Database) InsertPost(ctx context.Context, post Post, collectionName string) error {
	//Add post to the database
	_, _, err := db.Client.Collection(collectionName).Add(ctx, post)

	if err != nil {
		return err
	}

	return nil
}

//UpdateLikes updates the number of likes for the given post
func (db *Database) UpdateLikes(ctx context.Context, likes int, id, collectionName string) error {
	post := db.Client.Collection(collectionName).Doc(id)

	_, err := post.Update(ctx, []firestore.Update{
		{Path: "likes", Value: firestore.Increment(likes)},
	})

	if err != nil {
		return err
	}

	return nil
}
