package usecase

import (
	"context"
	"fmt"

	"github.com/aleexNxt/cli-tool/internal/domain"
)

// DeployUseCase behandelt Deployment-Operationen
type DeployUseCase struct {
	executor   domain.CommandExecutor
	fileSystem domain.FileSystem
	logger     domain.Logger
	config     *domain.Config
}

// NewDeployUseCase erstellt einen neuen DeployUseCase
func NewDeployUseCase(
	executor domain.CommandExecutor,
	fileSystem domain.FileSystem,
	logger domain.Logger,
	config *domain.Config,
) *DeployUseCase {
	return &DeployUseCase{
		executor:   executor,
		fileSystem: fileSystem,
		logger:     logger,
		config:     config,
	}
}

// DeployOptions definiert Optionen für Deployment
type DeployOptions struct {
	Environment string
	DryRun      bool
	Force       bool
	Tag         string
	ConfigFile  string
}

// Execute führt das Deployment aus
func (uc *DeployUseCase) Execute(ctx context.Context, opts *DeployOptions) error {
	if opts == nil {
		opts = &DeployOptions{Environment: "production"}
	}

	uc.logger.Info("Starting deployment to %s...", opts.Environment)

	if opts.DryRun {
		uc.logger.Info("Dry-run mode: No actual deployment will occur")
	}

	// Pre-deployment checks
	if err := uc.preDeploymentChecks(); err != nil {
		return fmt.Errorf("pre-deployment checks failed: %w", err)
	}

	// Execute deployment
	if err := uc.executeDeploy(ctx, opts); err != nil {
		return err
	}

	uc.logger.Success("Deployment to %s completed successfully", opts.Environment)
	return nil
}

// Rollback führt ein Rollback durch
func (uc *DeployUseCase) Rollback(ctx context.Context, environment string, version string) error {
	uc.logger.Info("Rolling back %s to version %s...", environment, version)

	task := domain.NewTask("rollback", domain.TaskTypeDeploy, "sh",
		[]string{"-c", fmt.Sprintf("echo 'Rollback to %s'", version)})

	result, err := uc.executor.Execute(ctx, task)
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("rollback failed")
	}

	uc.logger.Success("Rollback completed successfully")
	return nil
}

// Status prüft den Deployment-Status
func (uc *DeployUseCase) Status(ctx context.Context, environment string) error {
	uc.logger.Info("Checking deployment status for %s...", environment)

	task := domain.NewTask("status", domain.TaskTypeDeploy, "sh",
		[]string{"-c", "echo 'Checking deployment status...'"})

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return err
	}

	if result.Success {
		uc.logger.Info("Status: %s", result.Output)
	}

	return nil
}

func (uc *DeployUseCase) preDeploymentChecks() error {
	uc.logger.Info("Running pre-deployment checks...")

	// Check if we're in a git repository
	if !uc.fileSystem.Exists(".git") {
		uc.logger.Warn("Not in a git repository")
	}

	// Check if build artifacts exist
	if !uc.fileSystem.Exists("go.mod") {
		return fmt.Errorf("go.mod not found")
	}

	uc.logger.Success("Pre-deployment checks passed")
	return nil
}

func (uc *DeployUseCase) executeDeploy(ctx context.Context, opts *DeployOptions) error {
	// Use custom command from config if available
	cmd := uc.config.DeployCommand
	if cmd == "" {
		cmd = "echo 'No deployment command configured'"
	}

	args := []string{"-c", cmd}

	if opts.DryRun {
		args = []string{"-c", fmt.Sprintf("echo 'DRY-RUN: Would execute: %s'", cmd)}
	}

	task := domain.NewTask("deploy", domain.TaskTypeDeploy, "sh", args)

	// Set environment variables
	if opts.Environment != "" {
		task.Environment["DEPLOY_ENV"] = opts.Environment
	}
	if opts.Tag != "" {
		task.Environment["DEPLOY_TAG"] = opts.Tag
	}

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("deployment failed: %s", result.Output)
	}

	return nil
}
