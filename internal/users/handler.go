package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/bhumi46/users-svc/internal/httpx"
)

// Handler serves /api/users.
type Handler struct {
	store *Store
}

// NewHandler returns a Handler backed by store.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// Register adds the users routes to mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/users", h.list)
	mux.HandleFunc("GET /api/users/{id}", h.get)
	mux.HandleFunc("POST /api/users", h.create)
}

func (h *Handler) list(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, h.store.List())
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	u, err := h.store.Get(r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		httpx.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}

type createRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	if req.Name == "" || !strings.Contains(req.Email, "@") {
		httpx.WriteError(w, http.StatusBadRequest, "name and a valid email are required")
		return
	}
	u := h.store.Create(User{Name: req.Name, Email: req.Email})
	httpx.WriteJSON(w, http.StatusCreated, u)
}
