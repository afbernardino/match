package response

import (
	"log"
	"net/http"
)

const (
	ErrNotFound            string = `{"error":"not_found"}`
	ErrBadRequest          string = `{"error":"bad_request"}`
	ErrInternalServerError string = `{"error":"internal_server_error"}`
)

// Write writes byte array to http.ResponseWriter.
func Write(w http.ResponseWriter, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		log.Printf("error writing byte array: %v\n", err)
	}
}

// WriteInternalServerError writes internal server error response to http.ResponseWriter.
func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	Write(w, []byte(ErrInternalServerError))
}
