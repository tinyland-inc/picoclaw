package migrate

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tinyland-inc/tinyclaw/pkg/migrate"
)

func NewMigrateCommand() *cobra.Command {
	var opts migrate.Options

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate configuration between formats",
		Example: `  tinyclaw migrate
  tinyclaw migrate --dry-run
  tinyclaw migrate --to-dhall
  tinyclaw migrate --to-dhall --config /path/to/config.json`,
	}

	// Default: OpenClaw -> TinyClaw migration
	openclawCmd := &cobra.Command{
		Use:   "openclaw",
		Short: "Migrate from OpenClaw to TinyClaw",
		Args:  cobra.NoArgs,
		Example: `  tinyclaw migrate openclaw
  tinyclaw migrate openclaw --dry-run
  tinyclaw migrate openclaw --refresh
  tinyclaw migrate openclaw --force`,
		RunE: func(_ *cobra.Command, _ []string) error {
			result, err := migrate.Run(opts)
			if err != nil {
				return err
			}
			if !opts.DryRun {
				migrate.PrintSummary(result)
			}
			return nil
		},
	}

	openclawCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false,
		"Show what would be migrated without making changes")
	openclawCmd.Flags().BoolVar(&opts.Refresh, "refresh", false,
		"Re-sync workspace files from OpenClaw (repeatable)")
	openclawCmd.Flags().BoolVar(&opts.ConfigOnly, "config-only", false,
		"Only migrate config, skip workspace files")
	openclawCmd.Flags().BoolVar(&opts.WorkspaceOnly, "workspace-only", false,
		"Only migrate workspace files, skip config")
	openclawCmd.Flags().BoolVar(&opts.Force, "force", false,
		"Skip confirmation prompts")
	openclawCmd.Flags().StringVar(&opts.OpenClawHome, "openclaw-home", "",
		"Override OpenClaw home directory (default: ~/.openclaw)")
	openclawCmd.Flags().StringVar(&opts.TinyClawHome, "tinyclaw-home", "",
		"Override TinyClaw home directory (default: ~/.tinyclaw)")

	// to-dhall: JSON -> Dhall migration
	var dhallOpts migrate.ToDhallOptions

	toDhallCmd := &cobra.Command{
		Use:   "to-dhall",
		Short: "Convert JSON config to Dhall format",
		Args:  cobra.NoArgs,
		Example: `  tinyclaw migrate to-dhall
  tinyclaw migrate to-dhall --dry-run
  tinyclaw migrate to-dhall --config ~/.tinyclaw/config.json
  tinyclaw migrate to-dhall --output ~/.tinyclaw/config.dhall --force`,
		RunE: func(_ *cobra.Command, _ []string) error {
			result, err := migrate.RunToDhall(dhallOpts)
			if err != nil {
				return err
			}
			if !dhallOpts.DryRun {
				fmt.Printf("Dhall config written to %s\n", result.OutputPath)
			}
			if len(result.Warnings) > 0 {
				fmt.Println("\nWarnings:")
				for _, w := range result.Warnings {
					fmt.Printf("  - %s\n", w)
				}
			}
			return nil
		},
	}

	toDhallCmd.Flags().StringVar(&dhallOpts.ConfigPath, "config", "",
		"JSON config file path (default: ~/.tinyclaw/config.json)")
	toDhallCmd.Flags().StringVar(&dhallOpts.OutputPath, "output", "",
		"Dhall output file path (default: same dir as input, .dhall extension)")
	toDhallCmd.Flags().BoolVar(&dhallOpts.DryRun, "dry-run", false,
		"Print generated Dhall without writing")
	toDhallCmd.Flags().BoolVar(&dhallOpts.Force, "force", false,
		"Overwrite existing output file")

	// Legacy: bare `tinyclaw migrate` still runs OpenClaw migration for compatibility
	cmd.RunE = func(_ *cobra.Command, _ []string) error {
		result, err := migrate.Run(opts)
		if err != nil {
			return err
		}
		if !opts.DryRun {
			migrate.PrintSummary(result)
		}
		return nil
	}
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false,
		"Show what would be migrated without making changes")
	cmd.Flags().BoolVar(&opts.Force, "force", false,
		"Skip confirmation prompts")

	cmd.AddCommand(openclawCmd)
	cmd.AddCommand(toDhallCmd)

	return cmd
}
