package utils

import "net/http"

// HandleError helps standardize error responses
func HandleError(w http.ResponseWriter, msg string, code int) {
	http.Error(w, msg, code)
}
