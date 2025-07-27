package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody PreviewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println("1111 Failed to decode request body:", err)
		http.Error(w, "failed to parse JSON body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// validation
	if reqBody.UserID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	jsonBody, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonBody)

	resp, err := http.Post("http://trip-service:8083/preview", "Application/json", reader)
	if err != nil {
		log.Println("Failed to call trip service:", err)
		return
	}
	defer resp.Body.Close()

	var respBody any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		log.Println("Failed to decode response body:", err)
		return
	}

	response := contracts.APIResponse{Data: respBody}
	writeJson(w, http.StatusCreated, response)
}
