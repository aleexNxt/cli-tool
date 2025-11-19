package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/aleexNxt/cli-tool/internal/usecase"
)

func newDeployCommand() *cobra.Command {
	var (
		environment string
		dryRun      bool
		force       bool
		tag         string
		configFile  string
	)

	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deployment-Operationen",
		Long:  "Führt Deployment-Operationen für verschiedene Umgebungen aus",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			opts := &usecase.DeployOptions{
				Environment: environment,
				DryRun:      dryRun,
				Force:       force,
				Tag:         tag,
				ConfigFile:  configFile,
			}

			return deps.DeployUseCase.Execute(context.Background(), opts)
		},
	}

	deployCmd.Flags().StringVarP(&environment, "env", "e", "production", "Ziel-Umgebung")
	deployCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulation ohne tatsächliches Deployment")
	deployCmd.Flags().BoolVar(&force, "force", false, "Erzwinge Deployment")
	deployCmd.Flags().StringVarP(&tag, "tag", "t", "", "Version/Tag für Deployment")
	deployCmd.Flags().StringVarP(&configFile, "config-file", "f", "", "Deployment-Konfigurationsdatei")

	// Subcommands
	deployCmd.AddCommand(newDeployRollbackCommand())
	deployCmd.AddCommand(newDeployStatusCommand())

	return deployCmd
}

func newDeployRollbackCommand() *cobra.Command {
	var (
		environment string
		version     string
	)

	rollbackCmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback zu vorheriger Version",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.DeployUseCase.Rollback(context.Background(), environment, version)
		},
	}

	rollbackCmd.Flags().StringVarP(&environment, "env", "e", "production", "Ziel-Umgebung")
	rollbackCmd.Flags().StringVarP(&version, "version", "v", "", "Version für Rollback")
	rollbackCmd.MarkFlagRequired("version")

	return rollbackCmd
}

func newDeployStatusCommand() *cobra.Command {
	var environment string

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Zeigt Deployment-Status",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.DeployUseCase.Status(context.Background(), environment)
		},
	}

	statusCmd.Flags().StringVarP(&environment, "env", "e", "production", "Ziel-Umgebung")

	return statusCmd
}
