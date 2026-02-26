// Package api provides HTTP handlers for the TinyClaw API.
package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Dispatcher abstracts the agent loop for testability.
type Dispatcher interface {
	ProcessDirectWithChannel(ctx context.Context, content, sessionKey, channel, chatID string) (string, error)
	GetStartupInfo() map[string]any
	GetToolDefinitions() []map[string]any
}

// RouteRegistrar accepts new HTTP handler routes.
type RouteRegistrar interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// Handlers holds references needed by the API endpoints.
type Handlers struct {
	dispatcher Dispatcher
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(d Dispatcher) *Handlers {
	return &Handlers{dispatcher: d}
}

// Register adds all API routes to the given registrar.
func (h *Handlers) Register(r RouteRegistrar) {
	r.HandleFunc("POST /api/dispatch", h.handleDispatch)
	r.HandleFunc("GET /api/tools", h.handleTools)
	r.HandleFunc("GET /api/status", h.handleStatus)
}

type dispatchRequest struct {
	Content    string `json:"content"`
	SessionKey string `json:"session_key"`
	Channel    string `json:"channel"`
	ChatID     string `json:"chat_id"`
}

type dispatchResponse struct {
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason"`
	Error        string `json:"error,omitempty"`
}

func (h *Handlers) handleDispatch(w http.ResponseWriter, r *http.Request) {
	var req dispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dispatchResponse{Error: "invalid request body"})
		return
	}

	if req.Content == "" {
		writeJSON(w, http.StatusBadRequest, dispatchResponse{Error: "content is required"})
		return
	}

	// Apply defaults.
	if req.Channel == "" {
		req.Channel = "api"
	}
	if req.ChatID == "" {
		req.ChatID = "dispatch"
	}
	if req.SessionKey == "" {
		req.SessionKey = "api:" + req.ChatID
	}

	// Best-effort: extend the write deadline — dispatch can take minutes.
	// Errors are ignored because not all ResponseWriter implementations support this.
	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Now().Add(10 * time.Minute))

	result, err := h.dispatcher.ProcessDirectWithChannel(
		r.Context(), req.Content, req.SessionKey, req.Channel, req.ChatID,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dispatchResponse{Error: err.Error(), FinishReason: "error"})
		return
	}

	writeJSON(w, http.StatusOK, dispatchResponse{Content: result, FinishReason: "stop"})
}

func (h *Handlers) handleTools(w http.ResponseWriter, _ *http.Request) {
	defs := h.dispatcher.GetToolDefinitions()
	if defs == nil {
		defs = []map[string]any{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"tools": defs,
		"count": len(defs),
	})
}

func (h *Handlers) handleStatus(w http.ResponseWriter, _ *http.Request) {
	info := h.dispatcher.GetStartupInfo()
	if info == nil {
		info = map[string]any{}
	}
	writeJSON(w, http.StatusOK, info)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("api: failed to encode response: %v", err)
	}
}
