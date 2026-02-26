package skills

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal"
	"github.com/tinyland-inc/tinyclaw/pkg/skills"
)

func newInstallCommand(installerFn func() (*skills.SkillInstaller, error)) *cobra.Command {
	var registry string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install skill from GitHub",
		Example: `
tinyclaw skills install tinyland-inc/tinyclaw-skills/weather
tinyclaw skills install --registry clawhub github
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if registry != "" {
				if len(args) != 2 {
					return errors.New("when --registry is set, exactly 2 arguments are required: <name> <slug>")
				}
				return nil
			}

			if len(args) != 1 {
				return errors.New("exactly 1 argument is required: <github>")
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			installer, err := installerFn()
			if err != nil {
				return err
			}

			if registry != "" {
				cfg, err := internal.LoadConfig()
				if err != nil {
					return err
				}

				return skillsInstallFromRegistry(cfg, args[0], args[1])
			}

			return skillsInstallCmd(installer, args[0])
		},
	}

	cmd.Flags().StringVar(&registry, "registry", "", "Install from registry: --registry <name> <slug>")

	return cmd
}
