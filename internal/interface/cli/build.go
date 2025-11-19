package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aleexNxt/cli-tool/internal/usecase"
)

func newBuildCommand() *cobra.Command {
	var (
		outputDir  string
		race       bool
		tags       string
		goos       string
		goarch     string
		ldflags    string
		cgoEnabled bool
		clean      bool
	)

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build das Projekt",
		Long:  "Erstellt das Go-Projekt mit den angegebenen Optionen",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			ctx := context.Background()

			// Clean if requested
			if clean {
				if err := deps.BuildUseCase.Clean(ctx); err != nil {
					return err
				}
			}

			// Parse tags
			var tagList []string
			if tags != "" {
				tagList = strings.Split(tags, ",")
			}

			opts := &usecase.BuildOptions{
				OutputDir:  outputDir,
				Race:       race,
				Verbose:    verbose,
				Tags:       tagList,
				GOOS:       goos,
				GOARCH:     goarch,
				LDFlags:    ldflags,
				CGOEnabled: cgoEnabled,
			}

			return deps.BuildUseCase.Execute(ctx, opts)
		},
	}

	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output-Verzeichnis oder -Datei")
	buildCmd.Flags().BoolVar(&race, "race", false, "Race-Detector aktivieren")
	buildCmd.Flags().StringVarP(&tags, "tags", "t", "", "Build-Tags (kommagetrennt)")
	buildCmd.Flags().StringVar(&goos, "goos", "", "Ziel-Betriebssystem (GOOS)")
	buildCmd.Flags().StringVar(&goarch, "goarch", "", "Ziel-Architektur (GOARCH)")
	buildCmd.Flags().StringVarP(&ldflags, "ldflags", "l", "", "Linker-Flags")
	buildCmd.Flags().BoolVar(&cgoEnabled, "cgo", false, "CGO aktivieren")
	buildCmd.Flags().BoolVar(&clean, "clean", false, "Vor dem Build aufräumen")

	// Subcommands
	buildCmd.AddCommand(newBuildAllCommand())
	buildCmd.AddCommand(newBuildCleanCommand())

	return buildCmd
}

func newBuildAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Build alle Packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			ctx := context.Background()
			opts := &usecase.BuildOptions{
				Verbose: verbose,
			}

			return deps.BuildUseCase.BuildAll(ctx, opts)
		},
	}
}

func newBuildCleanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Entfernt Build-Artefakte",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.BuildUseCase.Clean(context.Background())
		},
	}
}
