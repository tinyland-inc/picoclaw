// Package adapters provides backend adapter interfaces for campaign dispatch.
//
// Each adapter implements the BackendAdapter interface, allowing campaigns
// to target different agent backends (TinyClaw, IronClaw, HexStrike-AI).
package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/tinyland-inc/tinyclaw/pkg/campaign"
)

// TinyClawAdapter dispatches campaign steps to the local TinyClaw agent loop.
type TinyClawAdapter struct {
	// ProcessFn is the function to call for processing a message.
	// This is wired to agent.AgentLoop.ProcessDirect at gateway startup.
	ProcessFn func(ctx context.Context, prompt, sessionKey string) (string, error)
}

// NewTinyClawAdapter creates a new adapter for the local TinyClaw agent.
func NewTinyClawAdapter() *TinyClawAdapter {
	return &TinyClawAdapter{}
}

func (a *TinyClawAdapter) Execute(ctx context.Context, agentID, prompt string, tools []string) (string, error) {
	if a.ProcessFn == nil {
		return "", errors.New("tinyclaw adapter not initialized: ProcessFn is nil")
	}
	sessionKey := "campaign:" + agentID
	return a.ProcessFn(ctx, prompt, sessionKey)
}

func (a *TinyClawAdapter) Name() string { return "tinyclaw" }

// Ensure TinyClawAdapter implements BackendAdapter
var _ campaign.BackendAdapter = (*TinyClawAdapter)(nil)

// StubAdapter is a test/development adapter that returns canned responses.
type StubAdapter struct {
	BackendName string
	Response    string
}

func (a *StubAdapter) Execute(_ context.Context, agentID, prompt string, _ []string) (string, error) {
	if a.Response != "" {
		return a.Response, nil
	}
	return fmt.Sprintf("[%s/%s] processed: %s", a.BackendName, agentID, prompt[:min(50, len(prompt))]), nil
}

func (a *StubAdapter) Name() string { return a.BackendName }

var _ campaign.BackendAdapter = (*StubAdapter)(nil)
