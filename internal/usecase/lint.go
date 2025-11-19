package usecase

import (
	"context"
	"fmt"

	"github.com/aleexNxt/cli-tool/internal/domain"
)

// LintUseCase behandelt Code-Linting-Operationen
type LintUseCase struct {
	executor   domain.CommandExecutor
	fileSystem domain.FileSystem
	logger     domain.Logger
	config     *domain.Config
}

// NewLintUseCase erstellt einen neuen LintUseCase
func NewLintUseCase(
	executor domain.CommandExecutor,
	fileSystem domain.FileSystem,
	logger domain.Logger,
	config *domain.Config,
) *LintUseCase {
	return &LintUseCase{
		executor:   executor,
		fileSystem: fileSystem,
		logger:     logger,
		config:     config,
	}
}

// LintOptions definiert Optionen für Linting
type LintOptions struct {
	Fast       bool
	Fix        bool
	Verbose    bool
	ConfigFile string
	Linters    []string
}

// Execute führt den Linter aus
func (uc *LintUseCase) Execute(ctx context.Context, opts *LintOptions) error {
	uc.logger.Info("Running linter...")

	if !uc.fileSystem.Exists("go.mod") {
		return fmt.Errorf("go.mod not found in current directory")
	}

	task := uc.createLintTask(opts)

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		uc.logger.Error("Linter execution failed: %v", err)
		return err
	}

	if !result.Success {
		uc.logger.Error("Linting issues found")
		return fmt.Errorf("linting failed: %s", result.Output)
	}

	uc.logger.Success("No linting issues found in %v", result.Duration)
	return nil
}

// Fix führt den Linter mit Auto-Fix aus
func (uc *LintUseCase) Fix(ctx context.Context) error {
	uc.logger.Info("Running linter with auto-fix...")

	opts := &LintOptions{
		Fix:     true,
		Verbose: true,
	}

	return uc.Execute(ctx, opts)
}

// Format formatiert den Code mit gofmt
func (uc *LintUseCase) Format(ctx context.Context) error {
	uc.logger.Info("Formatting code with gofmt...")

	task := domain.NewTask("format", domain.TaskTypeLint, "gofmt", []string{"-w", "."})

	result, err := uc.executor.Execute(ctx, task)
	if err != nil {
		return fmt.Errorf("formatting failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("formatting failed")
	}

	uc.logger.Success("Code formatted successfully")
	return nil
}

// Imports organisiert Imports mit goimports
func (uc *LintUseCase) Imports(ctx context.Context) error {
	uc.logger.Info("Organizing imports with goimports...")

	task := domain.NewTask("imports", domain.TaskTypeLint, "goimports", []string{"-w", "."})

	result, err := uc.executor.Execute(ctx, task)
	if err != nil {
		uc.logger.Warn("goimports not found, trying with go fmt...")
		return uc.Format(ctx)
	}

	if !result.Success {
		return fmt.Errorf("import organization failed")
	}

	uc.logger.Success("Imports organized successfully")
	return nil
}

func (uc *LintUseCase) createLintTask(opts *LintOptions) *domain.Task {
	// Verwende golangci-lint falls konfiguriert, sonst go vet
	cmd := uc.config.LintCommand
	if cmd == "" {
		cmd = "go vet ./..."
	}

	// Wenn golangci-lint verwendet wird
	if cmd == "golangci-lint run" || opts != nil {
		args := []string{"run"}

		if opts != nil {
			if opts.Fast {
				args = append(args, "--fast")
			}

			if opts.Fix {
				args = append(args, "--fix")
			}

			if opts.Verbose {
				args = append(args, "-v")
			}

			if opts.ConfigFile != "" {
				args = append(args, "--config", opts.ConfigFile)
			}

			if len(opts.Linters) > 0 {
				for _, linter := range opts.Linters {
					args = append(args, "--enable", linter)
				}
			}
		}

		return domain.NewTask("lint", domain.TaskTypeLint, "golangci-lint", args)
	}

	// Fallback zu go vet
	return domain.NewTask("lint", domain.TaskTypeLint, "go", []string{"vet", "./..."})
}
