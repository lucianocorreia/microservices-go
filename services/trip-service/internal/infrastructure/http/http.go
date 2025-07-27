package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
)

type HttpHandler struct {
	Service domain.TripService
}

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (h *HttpHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	fare := &domain.RideFareModel{
		UserID: "42",
	}

	ctx := r.Context()

	fmt.Println("Received request to create trip for user:", fare.UserID)
	t, err := h.Service.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
	}

	// log the details of the trip
	log.Printf("Trip created: ID=%s, UserID=%s, Status=%s, Fare=%v\n",
		t.ID.Hex(), t.UserID, t.Status, t.RideFare)

	writeJson(w, http.StatusOK, t)
}

// TODO: move to shared package
func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
