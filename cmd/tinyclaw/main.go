// TinyClaw - Ultra-lightweight personal AI agent
// Inspired by Sipeed PicoClaw, itself inspired by OpenClaw
// License: MIT
//
// Copyright (c) 2026 TinyClaw contributors

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal"
	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal/agent"
	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal/auth"
	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal/cron"
	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal/gateway"
	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal/migrate"
	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal/onboard"
	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal/skills"
	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal/status"
	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal/version"
)

func NewTinyclawCommand() *cobra.Command {
	short := fmt.Sprintf("%s tinyclaw - Personal AI Assistant v%s\n\n", internal.Logo, internal.GetVersion())

	cmd := &cobra.Command{
		Use:     "tinyclaw",
		Short:   short,
		Example: "tinyclaw list",
	}

	cmd.AddCommand(
		onboard.NewOnboardCommand(),
		agent.NewAgentCommand(),
		auth.NewAuthCommand(),
		gateway.NewGatewayCommand(),
		status.NewStatusCommand(),
		cron.NewCronCommand(),
		migrate.NewMigrateCommand(),
		skills.NewSkillsCommand(),
		version.NewVersionCommand(),
	)

	return cmd
}

func main() {
	cmd := NewTinyclawCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
