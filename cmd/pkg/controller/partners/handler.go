package partners

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"match/cmd/pkg/controller/response"
	"match/cmd/pkg/models"
	"match/cmd/pkg/repository"

	"github.com/gorilla/mux"
)

// Database can communicate with the persistent storage for our partners.
type Database interface {
	// GetMatches returns the best match for the customer, i.e. returns the partners that are experienced with the given materials
	// ordered by the highest score and closest location.
	GetMatches(ctx context.Context, materials []uint, lat, long float32) ([]models.Partner, error)

	// GetPartnerById returns a partner by id.
	GetPartnerById(ctx context.Context, id uint) (models.Partner, error)
}

// Handler handles '/partners' requests.
type Handler struct {
	db Database
}

// NewHandler creates a new Handler.
func NewHandler(db Database) Handler {
	return Handler{db: db}
}

// GetMatches returns the best match for the customer, i.e. returns the partners that are experienced with the given materials
// ordered by the highest score and closest location.
func (h *Handler) GetMatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	var reqBody models.MatchRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Printf("error decoding request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		response.Write(w, []byte(response.ErrBadRequest))
		return
	}

	var a models.Address
	if reqBody.Address == a || len(reqBody.Materials) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response.Write(w, []byte(response.ErrBadRequest))
		return
	}

	var matches []models.Partner
	matches, err = h.db.GetMatches(ctx, reqBody.Materials, reqBody.Address.Lat, reqBody.Address.Long)
	if err != nil {
		log.Printf("error retrieving matches from the database: %v\n", err)
		response.WriteInternalServerError(w)
		return
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(matches)
	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
		response.WriteInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.Write(w, jsonBytes)
}

// GetPartnerById returns a partner by id.
func (h *Handler) GetPartnerById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Write(w, []byte(response.ErrBadRequest))
		return
	}

	p, err := h.db.GetPartnerById(ctx, uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			response.Write(w, []byte(response.ErrNotFound))
			return
		}
		log.Printf("error retrieving the partner from the database: %v\n", err)
		response.WriteInternalServerError(w)
		return
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(p)
	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
		response.WriteInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.Write(w, jsonBytes)
}
