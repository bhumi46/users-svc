package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer() *http.ServeMux {
	mux := http.NewServeMux()
	NewHandler(NewStore(DefaultSeed()...)).Register(mux)
	return mux
}

func do(t *testing.T, mux http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestList(t *testing.T) {
	rec := do(t, newTestServer(), http.MethodGet, "/api/users", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var got []User
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].ID != "1" || got[2].ID != "3" {
		t.Fatalf("unexpected users: %+v", got)
	}
}

func TestGet(t *testing.T) {
	mux := newTestServer()
	tests := []struct {
		id   string
		want int
	}{
		{"1", http.StatusOK},
		{"999", http.StatusNotFound},
	}
	for _, tt := range tests {
		if rec := do(t, mux, http.MethodGet, "/api/users/"+tt.id, ""); rec.Code != tt.want {
			t.Errorf("GET /api/users/%s = %d, want %d", tt.id, rec.Code, tt.want)
		}
	}
}

func TestCreate(t *testing.T) {
	mux := newTestServer()
	tests := []struct {
		name string
		body string
		want int
	}{
		{"ok", `{"name":"Linus","email":"linus@example.com"}`, http.StatusCreated},
		{"bad json", `{`, http.StatusBadRequest},
		{"missing name", `{"email":"x@example.com"}`, http.StatusBadRequest},
		{"bad email", `{"name":"x","email":"nope"}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(t, mux, http.MethodPost, "/api/users", tt.body)
			if rec.Code != tt.want {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tt.want, rec.Body)
			}
		})
	}

	// The created user gets the next ID and can be fetched.
	if rec := do(t, mux, http.MethodGet, "/api/users/4", ""); rec.Code != http.StatusOK {
		t.Fatalf("created user not found: %d", rec.Code)
	}
}
