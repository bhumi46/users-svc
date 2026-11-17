// Package users implements the users API with an in-memory store.
package users

import (
	"errors"
	"sort"
	"strconv"
	"sync"
)

// User is a user of the lab.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrNotFound is returned when a user does not exist.
var ErrNotFound = errors.New("user not found")

// Store is a concurrency-safe in-memory user store.
type Store struct {
	mu     sync.RWMutex
	users  map[string]User
	nextID int
}

// NewStore returns a store seeded with the given users.
func NewStore(seed ...User) *Store {
	s := &Store{users: map[string]User{}, nextID: 1}
	for _, u := range seed {
		s.users[u.ID] = u
		if n, err := strconv.Atoi(u.ID); err == nil && n >= s.nextID {
			s.nextID = n + 1
		}
	}
	return s
}

// DefaultSeed is the data the service starts with.
func DefaultSeed() []User {
	return []User{
		{ID: "1", Name: "Ada Lovelace", Email: "ada@example.com"},
		{ID: "2", Name: "Alan Turing", Email: "alan@example.com"},
		{ID: "3", Name: "Grace Hopper", Email: "grace@example.com"},
	}
}

// List returns all users sorted by numeric ID.
func (s *Store) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	sort.Slice(out, func(i, j int) bool {
		a, _ := strconv.Atoi(out[i].ID)
		b, _ := strconv.Atoi(out[j].ID)
		return a < b
	})
	return out
}

// Get returns one user or ErrNotFound.
func (s *Store) Get(id string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

// Create assigns an ID and stores the user.
func (s *Store) Create(u User) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = strconv.Itoa(s.nextID)
	s.nextID++
	s.users[u.ID] = u
	return u
}
