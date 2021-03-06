package api

import "time"

//CountResponse response body on calls for getting counts
type CountResponse struct {
	Count int `json:"count"`
}

//Post structure of document inserted into the Firestore database
type Post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	Image       string    `json:"image" firestore:"image"`
	URL         string    `json:"url" firestore:"url"`
	Likes       int       `json:"likes" firestore:"likes"`
	Date        time.Time `json:"date" firestore:"date"`
}
