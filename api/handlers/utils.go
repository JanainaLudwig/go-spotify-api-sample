package handlers

import (
	"encoding/json"
	"go-api-template/core"
	"log"
	"net/http"
)

func sendResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("error decoding response", err)
		//_, _ = w.Write([]byte(`{"message": "An unexpected error occurred"}`))
	}
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*core.AppError); ok {
		log.Println("error response", apiErr)
		sendResponse(w, apiErr, apiErr.Status)
		return
	}

	sendResponse(w, core.AppError{
		Message: err.Error(),
		Status:  http.StatusInternalServerError,
	}, http.StatusInternalServerError)

	return
}
