// Package api implements a small JSON API for managing todos.
package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	maxBodyBytes   = 1 << 20
	maxTitleLength = 200
)

// Todo is a single todo item.
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Server holds the todo store and serves the HTTP API.
type Server struct {
	mu     sync.Mutex
	todos  map[int]Todo
	nextID int
}

// NewServer returns an empty Server.
func NewServer() *Server {
	return &Server{todos: map[int]Todo{}}
}

// Routes returns the HTTP handler for the API.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /todos", s.create)
	mux.HandleFunc("GET /todos/{id}", s.get)
	mux.HandleFunc("PUT /todos/{id}", s.update)
	return mux
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	t, err := decodeTodo(w, r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	s.nextID++
	t.ID = s.nextID
	s.todos[t.ID] = t
	s.mu.Unlock()

	writeJSON(w, http.StatusCreated, t)
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	s.mu.Lock()
	t, ok := s.todos[id]
	s.mu.Unlock()
	if !ok {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	t, err := decodeTodo(w, r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	_, ok := s.todos[id]
	if ok {
		t.ID = id
		s.todos[id] = t
	}
	s.mu.Unlock()
	if !ok {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}

	writeJSON(w, http.StatusOK, t)
}

// decodeTodo reads and validates a Todo from the request body.
func decodeTodo(w http.ResponseWriter, r *http.Request) (Todo, error) {
	var t Todo
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&t); err != nil {
		return Todo{}, errors.New("invalid JSON body")
	}
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		return Todo{}, errors.New("title is required")
	}
	if len(t.Title) > maxTitleLength {
		return Todo{}, errors.New("title is too long")
	}
	return t, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
