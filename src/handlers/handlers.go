package handlers

import (
	"fmt"
	"net/http"
)

// IndexHandler handle the requests to /
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World !")
}
