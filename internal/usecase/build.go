package usecase

import (
	"context"
	"fmt"

	"github.com/aleexNxt/cli-tool/internal/domain"
)

// BuildUseCase behandelt Build-Operationen
type BuildUseCase struct {
	executor   domain.CommandExecutor
	fileSystem domain.FileSystem
	logger     domain.Logger
	config     *domain.Config
}

// NewBuildUseCase erstellt einen neuen BuildUseCase
func NewBuildUseCase(
	executor domain.CommandExecutor,
	fileSystem domain.FileSystem,
	logger domain.Logger,
	config *domain.Config,
) *BuildUseCase {
	return &BuildUseCase{
		executor:   executor,
		fileSystem: fileSystem,
		logger:     logger,
		config:     config,
	}
}

// BuildOptions definiert Optionen für den Build-Prozess
type BuildOptions struct {
	OutputDir   string
	CGOEnabled  bool
	Race        bool
	Verbose     bool
	Tags        []string
	GOOS        string
	GOARCH      string
	LDFlags     string
	CustomCmd   string
}

// Execute führt den Build-Prozess aus
func (uc *BuildUseCase) Execute(ctx context.Context, opts *BuildOptions) error {
	uc.logger.Info("Starting build process...")

	// Prüfe ob go.mod existiert
	if !uc.fileSystem.Exists("go.mod") {
		return fmt.Errorf("go.mod not found in current directory")
	}

	// Erstelle Build-Task
	task := uc.createBuildTask(opts)

	// Führe Build aus
	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		uc.logger.Error("Build failed: %v", err)
		return err
	}

	if !result.Success {
		uc.logger.Error("Build failed")
		return fmt.Errorf("build failed: %s", result.Output)
	}

	uc.logger.Success("Build completed successfully in %v", result.Duration)
	return nil
}

// Clean entfernt Build-Artefakte
func (uc *BuildUseCase) Clean(ctx context.Context) error {
	uc.logger.Info("Cleaning build artifacts...")

	task := domain.NewTask("clean", domain.TaskTypeBuild, "go", []string{"clean", "-cache", "-testcache"})

	result, err := uc.executor.Execute(ctx, task)
	if err != nil {
		return fmt.Errorf("clean failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("clean failed: %s", result.Output)
	}

	uc.logger.Success("Clean completed successfully")
	return nil
}

// BuildAll erstellt alle Packages
func (uc *BuildUseCase) BuildAll(ctx context.Context, opts *BuildOptions) error {
	uc.logger.Info("Building all packages...")

	task := domain.NewTask("build-all", domain.TaskTypeBuild, "go", []string{"build", "-v", "./..."})

	if opts != nil && opts.Verbose {
		task.Args = append(task.Args, "-x")
	}

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return err
	}

	if !result.Success {
		uc.logger.Error("Build failed: %s", result.Output)
		return fmt.Errorf("build failed")
	}

	uc.logger.Success("All packages built successfully")
	return nil
}

func (uc *BuildUseCase) createBuildTask(opts *BuildOptions) *domain.Task {
	args := []string{"build"}

	if opts == nil {
		opts = &BuildOptions{}
	}

	// Verwende custom command falls angegeben
	if opts.CustomCmd != "" {
		task := domain.NewTask("custom-build", domain.TaskTypeBuild, "sh", []string{"-c", opts.CustomCmd})
		return task
	}

	// Baue Build-Argumente
	if opts.Verbose {
		args = append(args, "-v")
	}

	if opts.Race {
		args = append(args, "-race")
	}

	if len(opts.Tags) > 0 {
		tagStr := ""
		for i, tag := range opts.Tags {
			if i > 0 {
				tagStr += ","
			}
			tagStr += tag
		}
		args = append(args, "-tags", tagStr)
	}

	if opts.LDFlags != "" {
		args = append(args, "-ldflags", opts.LDFlags)
	}

	if opts.OutputDir != "" {
		args = append(args, "-o", opts.OutputDir)
	}

	args = append(args, "./...")

	task := domain.NewTask("build", domain.TaskTypeBuild, "go", args)

	// Setze Umgebungsvariablen
	if opts.GOOS != "" {
		task.Environment["GOOS"] = opts.GOOS
	}
	if opts.GOARCH != "" {
		task.Environment["GOARCH"] = opts.GOARCH
	}
	if !opts.CGOEnabled {
		task.Environment["CGO_ENABLED"] = "0"
	}

	return task
}
