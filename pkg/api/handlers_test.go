package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockDispatcher implements Dispatcher for testing.
type mockDispatcher struct {
	result string
	err    error
	tools  []map[string]any
	info   map[string]any
}

func (m *mockDispatcher) ProcessDirectWithChannel(_ context.Context, _, _, _, _ string) (string, error) {
	return m.result, m.err
}

func (m *mockDispatcher) GetStartupInfo() map[string]any       { return m.info }
func (m *mockDispatcher) GetToolDefinitions() []map[string]any { return m.tools }

func newTestMux(d Dispatcher) *http.ServeMux {
	mux := http.NewServeMux()
	h := NewHandlers(d)
	h.Register(mux)
	return mux
}

func TestDispatch_ValidRequest(t *testing.T) {
	mux := newTestMux(&mockDispatcher{result: "hello back"})

	body := `{"content":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/dispatch", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp dispatchResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Content != "hello back" {
		t.Errorf("expected 'hello back', got %q", resp.Content)
	}
	if resp.FinishReason != "stop" {
		t.Errorf("expected finish_reason 'stop', got %q", resp.FinishReason)
	}
}

func TestDispatch_EmptyContent(t *testing.T) {
	mux := newTestMux(&mockDispatcher{})

	body := `{"content":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/dispatch", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestDispatch_MissingContent(t *testing.T) {
	mux := newTestMux(&mockDispatcher{})

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/dispatch", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestDispatch_WrongMethod(t *testing.T) {
	mux := newTestMux(&mockDispatcher{})

	req := httptest.NewRequest(http.MethodGet, "/api/dispatch", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestDispatch_ProcessError(t *testing.T) {
	mux := newTestMux(&mockDispatcher{err: context.DeadlineExceeded})

	body := `{"content":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/dispatch", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}

	var resp dispatchResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.FinishReason != "error" {
		t.Errorf("expected finish_reason 'error', got %q", resp.FinishReason)
	}
}

func TestTools(t *testing.T) {
	defs := []map[string]any{
		{"name": "web_search", "description": "Search the web"},
		{"name": "message", "description": "Send a message"},
	}
	mux := newTestMux(&mockDispatcher{tools: defs})

	req := httptest.NewRequest(http.MethodGet, "/api/tools", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	count, ok := resp["count"].(float64)
	if !ok {
		t.Fatalf("expected count to be a number, got %T", resp["count"])
	}
	if int(count) != 2 {
		t.Errorf("expected count 2, got %v", count)
	}
}

func TestTools_NilDefinitions(t *testing.T) {
	mux := newTestMux(&mockDispatcher{tools: nil})

	req := httptest.NewRequest(http.MethodGet, "/api/tools", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	count, ok := resp["count"].(float64)
	if !ok {
		t.Fatalf("expected count to be a number, got %T", resp["count"])
	}
	if int(count) != 0 {
		t.Errorf("expected count 0, got %v", count)
	}
}

func TestStatus(t *testing.T) {
	info := map[string]any{
		"tools": map[string]any{"count": 5, "names": []string{"a", "b", "c", "d", "e"}},
	}
	mux := newTestMux(&mockDispatcher{info: info})

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	tools, ok := resp["tools"].(map[string]any)
	if !ok {
		t.Fatal("expected tools key in response")
	}
	count, ok := tools["count"].(float64)
	if !ok {
		t.Fatalf("expected count to be a number, got %T", tools["count"])
	}
	if int(count) != 5 {
		t.Errorf("expected tool count 5, got %v", count)
	}
}
