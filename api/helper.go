package api

import "net/http"

//GetURLParameter returns the value from the provided URL parameter if found
func GetURLParameter(r *http.Request, param string) string {

	keys, ok := r.URL.Query()[param]
	if !ok || len(keys[0]) < 1 {
		return ""
	}

	return keys[0]
}
