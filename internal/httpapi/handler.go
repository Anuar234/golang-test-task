package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
)

type NumbersStore interface {
	AddAndList(ctx context.Context, value int64) ([]int64, error)
}

type Handler struct {
	store NumbersStore
}

func NewRouter(store NumbersStore) http.Handler {
	h := &Handler{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("/numbers", h.handleNumbers)
	mux.HandleFunc("/healthz", h.handleHealth)
	return mux
}

func (h *Handler) handleNumbers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Number *int64 `json:"number"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Number == nil {
		http.Error(w, "missing number", http.StatusBadRequest)
		return
	}

	numbers, err := h.store.AddAndList(r.Context(), *req.Number)
	if err != nil {
		http.Error(w, "failed to store number", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string][]int64{"numbers": numbers})
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
