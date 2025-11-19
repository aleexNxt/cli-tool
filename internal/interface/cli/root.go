package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aleexNxt/cli-tool/internal/domain"
	"github.com/aleexNxt/cli-tool/internal/infrastructure/config"
	"github.com/aleexNxt/cli-tool/internal/infrastructure/executor"
	"github.com/aleexNxt/cli-tool/internal/infrastructure/filesystem"
	"github.com/aleexNxt/cli-tool/internal/usecase"
	"github.com/aleexNxt/cli-tool/pkg/logger"
)

var (
	verbose    bool
	configPath string
)

// Dependencies hält alle Use Cases und Services
type Dependencies struct {
	BuildUseCase  *usecase.BuildUseCase
	TestUseCase   *usecase.TestUseCase
	DeployUseCase *usecase.DeployUseCase
	DockerUseCase *usecase.DockerUseCase
	LintUseCase   *usecase.LintUseCase
	Logger        domain.Logger
	Config        *domain.Config
}

// NewRootCommand erstellt das Root-Command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "devtool",
		Short: "DevTool - Automatisierung für Entwicklungs- und Deployment-Aufgaben",
		Long: `DevTool ist ein Kommandozeilen-Tool zur Automatisierung von
wiederkehrenden Entwicklungs- und Deployment-Aufgaben.

Es bietet Befehle für:
  - Build-Automation (build)
  - Test-Automation (test)
  - Code-Linting (lint)
  - Deployment (deploy)
  - Docker-Operationen (docker)
  - Konfigurationsverwaltung (config)`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// This runs before any subcommand
		},
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Ausführlicher Output")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Pfad zur Konfigurationsdatei")

	// Add subcommands
	rootCmd.AddCommand(newBuildCommand())
	rootCmd.AddCommand(newTestCommand())
	rootCmd.AddCommand(newLintCommand())
	rootCmd.AddCommand(newDeployCommand())
	rootCmd.AddCommand(newDockerCommand())
	rootCmd.AddCommand(newConfigCommand())
	rootCmd.AddCommand(newVersionCommand())

	return rootCmd
}

// Execute führt das Root-Command aus
func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// InitDependencies initialisiert alle Dependencies
func InitDependencies() (*Dependencies, error) {
	// Logger
	log := logger.NewLogger(verbose)

	// FileSystem
	fs := filesystem.NewFileSystem()

	// Config Repository
	var configRepo domain.ConfigRepository
	if configPath != "" {
		configRepo = config.NewConfigRepositoryWithPath(configPath)
	} else {
		configRepo = config.NewConfigRepository()
	}

	// Load Config
	cfg, err := configRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Working Directory
	workDir, err := fs.GetWorkingDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Command Executor
	exec := executor.NewCommandExecutor(workDir, 0)

	// Use Cases
	buildUC := usecase.NewBuildUseCase(exec, fs, log, cfg)
	testUC := usecase.NewTestUseCase(exec, fs, log, cfg)
	lintUC := usecase.NewLintUseCase(exec, fs, log, cfg)
	deployUC := usecase.NewDeployUseCase(exec, fs, log, cfg)
	dockerUC := usecase.NewDockerUseCase(exec, fs, log, cfg)

	return &Dependencies{
		BuildUseCase:  buildUC,
		TestUseCase:   testUC,
		LintUseCase:   lintUC,
		DeployUseCase: deployUC,
		DockerUseCase: dockerUC,
		Logger:        log,
		Config:        cfg,
	}, nil
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Zeigt die Version an",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("DevTool v1.0.0")
		},
	}
}
