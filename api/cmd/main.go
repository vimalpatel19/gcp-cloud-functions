package main

import (
	"net/http"

	"github.com/vimalpatel19/gcp-cloud-functions/api"
)

func main() {
	http.HandleFunc("/posts-save", api.SavePost)
	http.HandleFunc("/posts-find", api.FindPosts)
	http.HandleFunc("/posts-count", api.CountPosts)
	http.ListenAndServe(":8080", nil)
}
