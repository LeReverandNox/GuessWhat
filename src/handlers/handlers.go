package handlers

import "net/http"

// IndexHandler handle the requests to /
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}
