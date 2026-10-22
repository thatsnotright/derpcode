package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCreate(t *testing.T) {
	h := NewServer().Routes()
	rec := do(t, h, http.MethodPost, "/todos", `{"title":"write tests"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusCreated)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q; want application/json", ct)
	}
	var got Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if got.ID != 1 || got.Title != "write tests" {
		t.Fatalf("got %+v", got)
	}
}

func TestCreateBadJSON(t *testing.T) {
	h := NewServer().Routes()
	rec := do(t, h, http.MethodPost, "/todos", `{"title":`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusBadRequest)
	}
	// nothing should have been stored
	if rec := do(t, h, http.MethodGet, "/todos/1", ""); rec.Code != http.StatusNotFound {
		t.Fatalf("GET /todos/1 status = %d; want %d", rec.Code, http.StatusNotFound)
	}
}

func TestGet(t *testing.T) {
	h := NewServer().Routes()
	do(t, h, http.MethodPost, "/todos", `{"title":"a"}`)

	rec := do(t, h, http.MethodGet, "/todos/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"title":"a"`) {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestGetNotFound(t *testing.T) {
	h := NewServer().Routes()
	rec := do(t, h, http.MethodGet, "/todos/42", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
	}
}

func TestUpdate(t *testing.T) {
	h := NewServer().Routes()
	do(t, h, http.MethodPost, "/todos", `{"title":"a"}`)

	rec := do(t, h, http.MethodPut, "/todos/1", `{"title":"b","done":true}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rec.Code)
	}
	var got Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got != (Todo{ID: 1, Title: "b", Done: true}) {
		t.Fatalf("got %+v", got)
	}
}

func TestUpdateBadJSON(t *testing.T) {
	h := NewServer().Routes()
	do(t, h, http.MethodPost, "/todos", `{"title":"keep me"}`)

	rec := do(t, h, http.MethodPut, "/todos/1", `not json`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusBadRequest)
	}
	rec = do(t, h, http.MethodGet, "/todos/1", "")
	if !strings.Contains(rec.Body.String(), "keep me") {
		t.Fatalf("todo was overwritten: %s", rec.Body.String())
	}
}

func TestRealServer(t *testing.T) {
	srv := httptest.NewServer(NewServer().Routes())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/todos", "application/json", strings.NewReader(`{"title":"x"}`))
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(b), `"id":1`) {
		t.Fatalf("body = %s", b)
	}

	resp, _ = http.Get(srv.URL + "/todos/1")
	b, _ = io.ReadAll(resp.Body)
	if !strings.Contains(string(b), `"title":"x"`) {
		t.Fatalf("body = %s", b)
	}
}
