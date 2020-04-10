package main

import (
	"net/http"

	"github.com/vimalpatel19/gcp-cloud-functions/api"
)

func main() {
	http.HandleFunc("/posts-save", api.SavePost)
	http.ListenAndServe(":8080", nil)
}
