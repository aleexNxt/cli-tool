package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/aleexNxt/cli-tool/internal/usecase"
)

func newLintCommand() *cobra.Command {
	var (
		fast       bool
		fix        bool
		configFile string
		linters    []string
	)

	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Führt Code-Linting aus",
		Long:  "Führt Code-Linting mit golangci-lint oder go vet aus",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			opts := &usecase.LintOptions{
				Fast:       fast,
				Fix:        fix,
				Verbose:    verbose,
				ConfigFile: configFile,
				Linters:    linters,
			}

			return deps.LintUseCase.Execute(context.Background(), opts)
		},
	}

	lintCmd.Flags().BoolVar(&fast, "fast", false, "Schnelles Linting (weniger Linters)")
	lintCmd.Flags().BoolVar(&fix, "fix", false, "Auto-Fix für Probleme")
	lintCmd.Flags().StringVarP(&configFile, "config-file", "f", "", "Pfad zur Linter-Konfiguration")
	lintCmd.Flags().StringSliceVarP(&linters, "enable", "e", []string{}, "Spezifische Linters aktivieren")

	// Subcommands
	lintCmd.AddCommand(newLintFixCommand())
	lintCmd.AddCommand(newLintFormatCommand())
	lintCmd.AddCommand(newLintImportsCommand())

	return lintCmd
}

func newLintFixCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fix",
		Short: "Führt Linting mit Auto-Fix aus",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.LintUseCase.Fix(context.Background())
		},
	}
}

func newLintFormatCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fmt",
		Short: "Formatiert Code mit gofmt",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.LintUseCase.Format(context.Background())
		},
	}
}

func newLintImportsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "imports",
		Short: "Organisiert Imports mit goimports",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.LintUseCase.Imports(context.Background())
		},
	}
}
