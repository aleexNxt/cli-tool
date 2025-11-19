package usecase

import (
	"context"
	"fmt"

	"github.com/aleexNxt/cli-tool/internal/domain"
)

// TestUseCase behandelt Test-Operationen
type TestUseCase struct {
	executor   domain.CommandExecutor
	fileSystem domain.FileSystem
	logger     domain.Logger
	config     *domain.Config
}

// NewTestUseCase erstellt einen neuen TestUseCase
func NewTestUseCase(
	executor domain.CommandExecutor,
	fileSystem domain.FileSystem,
	logger domain.Logger,
	config *domain.Config,
) *TestUseCase {
	return &TestUseCase{
		executor:   executor,
		fileSystem: fileSystem,
		logger:     logger,
		config:     config,
	}
}

// TestOptions definiert Optionen für Tests
type TestOptions struct {
	Verbose     bool
	Coverage    bool
	Race        bool
	Short       bool
	Package     string
	Run         string
	Count       int
	Parallel    int
	Timeout     string
	BenchMem    bool
}

// Execute führt Tests aus
func (uc *TestUseCase) Execute(ctx context.Context, opts *TestOptions) error {
	uc.logger.Info("Running tests...")

	if !uc.fileSystem.Exists("go.mod") {
		return fmt.Errorf("go.mod not found in current directory")
	}

	task := uc.createTestTask(opts)

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		uc.logger.Error("Test execution failed: %v", err)
		return err
	}

	if !result.Success {
		uc.logger.Error("Tests failed")
		return fmt.Errorf("tests failed: %s", result.Output)
	}

	uc.logger.Success("All tests passed in %v", result.Duration)
	return nil
}

// RunBenchmarks führt Benchmarks aus
func (uc *TestUseCase) RunBenchmarks(ctx context.Context, opts *TestOptions) error {
	uc.logger.Info("Running benchmarks...")

	args := []string{"test", "-bench=.", "-benchmem"}

	if opts != nil {
		if opts.Package != "" {
			args = append(args, opts.Package)
		} else {
			args = append(args, "./...")
		}

		if opts.Run != "" {
			args = append(args, "-run", opts.Run)
		}

		if opts.Count > 0 {
			args = append(args, "-count", fmt.Sprintf("%d", opts.Count))
		}
	} else {
		args = append(args, "./...")
	}

	task := domain.NewTask("benchmark", domain.TaskTypeTest, "go", args)

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("benchmarks failed")
	}

	uc.logger.Success("Benchmarks completed successfully")
	return nil
}

// GenerateCoverage generiert Coverage-Report
func (uc *TestUseCase) GenerateCoverage(ctx context.Context, outputFile string) error {
	uc.logger.Info("Generating coverage report...")

	// Run tests with coverage
	args := []string{"test", "-coverprofile=" + outputFile, "./..."}
	task := domain.NewTask("coverage", domain.TaskTypeTest, "go", args)

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("coverage generation failed")
	}

	uc.logger.Success("Coverage report generated: %s", outputFile)

	// Generate HTML report
	htmlFile := outputFile + ".html"
	htmlTask := domain.NewTask("coverage-html", domain.TaskTypeTest, "go",
		[]string{"tool", "cover", "-html=" + outputFile, "-o", htmlFile})

	htmlResult, err := uc.executor.Execute(ctx, htmlTask)
	if err == nil && htmlResult.Success {
		uc.logger.Success("HTML coverage report: %s", htmlFile)
	}

	return nil
}

// Vet führt go vet aus
func (uc *TestUseCase) Vet(ctx context.Context) error {
	uc.logger.Info("Running go vet...")

	task := domain.NewTask("vet", domain.TaskTypeTest, "go", []string{"vet", "./..."})

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return err
	}

	if !result.Success {
		uc.logger.Error("Vet found issues: %s", result.Output)
		return fmt.Errorf("vet failed")
	}

	uc.logger.Success("No issues found by vet")
	return nil
}

func (uc *TestUseCase) createTestTask(opts *TestOptions) *domain.Task {
	args := []string{"test"}

	if opts == nil {
		opts = &TestOptions{Verbose: true}
	}

	if opts.Verbose {
		args = append(args, "-v")
	}

	if opts.Coverage {
		args = append(args, "-cover")
	}

	if opts.Race {
		args = append(args, "-race")
	}

	if opts.Short {
		args = append(args, "-short")
	}

	if opts.Count > 0 {
		args = append(args, "-count", fmt.Sprintf("%d", opts.Count))
	}

	if opts.Parallel > 0 {
		args = append(args, "-parallel", fmt.Sprintf("%d", opts.Parallel))
	}

	if opts.Timeout != "" {
		args = append(args, "-timeout", opts.Timeout)
	}

	if opts.Run != "" {
		args = append(args, "-run", opts.Run)
	}

	if opts.Package != "" {
		args = append(args, opts.Package)
	} else {
		args = append(args, "./...")
	}

	return domain.NewTask("test", domain.TaskTypeTest, "go", args)
}
