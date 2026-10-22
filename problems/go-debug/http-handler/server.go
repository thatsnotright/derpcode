// Package api implements a small JSON API for managing todos.
package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var t Todo
	if err := json.Unmarshal(body, &t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
	}

	s.mu.Lock()
	s.nextID++
	t.ID = s.nextID
	s.todos[t.ID] = t
	s.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	t, ok := s.todos[id]
	s.mu.Unlock()
	if !ok {
		http.Error(w, "todo not found", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(t)
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	t, err := decodeTodo(r)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}
	t.ID = id
	s.todos[id] = t

	json.NewEncoder(w).Encode(t)
}

func decodeTodo(r *http.Request) (t Todo, err error) {
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			log.Printf("decode todo: %v", err)
		}
	}
	return t, err
}
