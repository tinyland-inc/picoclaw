package tools

import (
	"context"
	"fmt"
)

// ApprovalNotifyTool sends approval request notifications to a configured
// channel (typically XMPP). When gnucash-bridge returns require_approval from
// security_check, the agent invokes this tool to notify the user.
type ApprovalNotifyTool struct {
	sendCallback SendCallback
	channel      string // default notification channel (e.g. "xmpp")
	chatID       string // default notification recipient
}

func NewApprovalNotifyTool() *ApprovalNotifyTool {
	return &ApprovalNotifyTool{}
}

func (t *ApprovalNotifyTool) Name() string {
	return "approval_notify"
}

func (t *ApprovalNotifyTool) Description() string {
	return "Send an approval request notification to the user. Use when a financial operation requires human approval before proceeding."
}

func (t *ApprovalNotifyTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"approval_id": map[string]any{
				"type":        "string",
				"description": "The approval request ID from security_check",
			},
			"agent": map[string]any{
				"type":        "string",
				"description": "The agent requesting approval",
			},
			"tool": map[string]any{
				"type":        "string",
				"description": "The tool/operation that needs approval",
			},
			"args_summary": map[string]any{
				"type":        "string",
				"description": "Human-readable summary of the operation arguments",
			},
			"channel": map[string]any{
				"type":        "string",
				"description": "Optional: override notification channel",
			},
			"chat_id": map[string]any{
				"type":        "string",
				"description": "Optional: override notification recipient",
			},
		},
		"required": []string{"approval_id", "agent", "tool"},
	}
}

func (t *ApprovalNotifyTool) SetContext(channel, chatID string) {
	t.channel = channel
	t.chatID = chatID
}

func (t *ApprovalNotifyTool) SetSendCallback(callback SendCallback) {
	t.sendCallback = callback
}

func (t *ApprovalNotifyTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	approvalID, _ := args["approval_id"].(string)
	agent, _ := args["agent"].(string)
	tool, _ := args["tool"].(string)
	argsSummary, _ := args["args_summary"].(string)

	if approvalID == "" {
		return ErrorResult("approval_id is required")
	}
	if agent == "" {
		return ErrorResult("agent is required")
	}
	if tool == "" {
		return ErrorResult("tool is required")
	}

	channel, _ := args["channel"].(string)
	chatID, _ := args["chat_id"].(string)
	if channel == "" {
		channel = t.channel
	}
	if chatID == "" {
		chatID = t.chatID
	}

	if t.sendCallback == nil {
		return ErrorResult("notification delivery not configured")
	}

	if channel == "" || chatID == "" {
		return ErrorResult("notification channel and chat_id are required")
	}

	// Format the approval notification message
	message := fmt.Sprintf(
		"Approval needed: %s wants to %s",
		agent, tool,
	)
	if argsSummary != "" {
		message += fmt.Sprintf(" with: %s", argsSummary)
	}
	message += fmt.Sprintf("\n\nReply /approve %s or /reject %s <reason>", approvalID, approvalID)

	if err := t.sendCallback(channel, chatID, message); err != nil {
		return ErrorResult(fmt.Sprintf("failed to send notification: %s", err.Error()))
	}

	return SilentResult(fmt.Sprintf("Approval notification sent for %s (ID: %s)", tool, approvalID))
}
