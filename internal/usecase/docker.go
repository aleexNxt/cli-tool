package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/aleexNxt/cli-tool/internal/domain"
)

// DockerUseCase behandelt Docker-Operationen
type DockerUseCase struct {
	executor   domain.CommandExecutor
	fileSystem domain.FileSystem
	logger     domain.Logger
	config     *domain.Config
}

// NewDockerUseCase erstellt einen neuen DockerUseCase
func NewDockerUseCase(
	executor domain.CommandExecutor,
	fileSystem domain.FileSystem,
	logger domain.Logger,
	config *domain.Config,
) *DockerUseCase {
	return &DockerUseCase{
		executor:   executor,
		fileSystem: fileSystem,
		logger:     logger,
		config:     config,
	}
}

// BuildOptions definiert Optionen für Docker Build
type DockerBuildOptions struct {
	ImageName  string
	Tag        string
	Dockerfile string
	Context    string
	NoBuild    bool
	NoCache    bool
	Platform   string
	BuildArgs  map[string]string
}

// Build erstellt ein Docker Image
func (uc *DockerUseCase) Build(ctx context.Context, opts *DockerBuildOptions) error {
	if opts == nil {
		return fmt.Errorf("build options are required")
	}

	uc.logger.Info("Building Docker image: %s:%s", opts.ImageName, opts.Tag)

	// Check if Dockerfile exists
	dockerfile := opts.Dockerfile
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}

	if !uc.fileSystem.Exists(dockerfile) {
		return fmt.Errorf("Dockerfile not found: %s", dockerfile)
	}

	task := uc.createBuildTask(opts, dockerfile)

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return fmt.Errorf("docker build failed: %w", err)
	}

	if !result.Success {
		uc.logger.Error("Docker build failed: %s", result.Output)
		return fmt.Errorf("docker build failed")
	}

	uc.logger.Success("Docker image built successfully: %s:%s", opts.ImageName, opts.Tag)
	return nil
}

// Push pusht ein Image zur Registry
func (uc *DockerUseCase) Push(ctx context.Context, imageName, tag string) error {
	fullImage := fmt.Sprintf("%s:%s", imageName, tag)
	uc.logger.Info("Pushing Docker image: %s", fullImage)

	// Add registry prefix if configured
	if uc.config.DockerRegistry != "" {
		fullImage = fmt.Sprintf("%s/%s", uc.config.DockerRegistry, fullImage)
	}

	task := domain.NewTask("docker-push", domain.TaskTypeDocker, "docker",
		[]string{"push", fullImage})

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return fmt.Errorf("docker push failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("docker push failed: %s", result.Output)
	}

	uc.logger.Success("Image pushed successfully: %s", fullImage)
	return nil
}

// Run startet einen Container
func (uc *DockerUseCase) Run(ctx context.Context, imageName, containerName string, ports []string, env map[string]string) error {
	uc.logger.Info("Starting container: %s", containerName)

	args := []string{"run", "-d"}

	if containerName != "" {
		args = append(args, "--name", containerName)
	}

	// Add port mappings
	for _, port := range ports {
		args = append(args, "-p", port)
	}

	// Add environment variables
	for key, val := range env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, val))
	}

	args = append(args, imageName)

	task := domain.NewTask("docker-run", domain.TaskTypeDocker, "docker", args)

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return fmt.Errorf("docker run failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("docker run failed: %s", result.Output)
	}

	containerID := strings.TrimSpace(result.Output)
	uc.logger.Success("Container started: %s (ID: %s)", containerName, containerID)
	return nil
}

// Stop stoppt einen Container
func (uc *DockerUseCase) Stop(ctx context.Context, containerName string) error {
	uc.logger.Info("Stopping container: %s", containerName)

	task := domain.NewTask("docker-stop", domain.TaskTypeDocker, "docker",
		[]string{"stop", containerName})

	result, err := uc.executor.Execute(ctx, task)
	if err != nil {
		return fmt.Errorf("docker stop failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("docker stop failed")
	}

	uc.logger.Success("Container stopped: %s", containerName)
	return nil
}

// Remove entfernt einen Container
func (uc *DockerUseCase) Remove(ctx context.Context, containerName string, force bool) error {
	uc.logger.Info("Removing container: %s", containerName)

	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, containerName)

	task := domain.NewTask("docker-rm", domain.TaskTypeDocker, "docker", args)

	result, err := uc.executor.Execute(ctx, task)
	if err != nil {
		return fmt.Errorf("docker rm failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("docker rm failed")
	}

	uc.logger.Success("Container removed: %s", containerName)
	return nil
}

// ListImages listet Docker Images
func (uc *DockerUseCase) ListImages(ctx context.Context) error {
	uc.logger.Info("Listing Docker images...")

	task := domain.NewTask("docker-images", domain.TaskTypeDocker, "docker",
		[]string{"images"})

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return fmt.Errorf("docker images failed: %w", err)
	}

	if result.Success {
		fmt.Println(result.Output)
	}

	return nil
}

// ListContainers listet Docker Container
func (uc *DockerUseCase) ListContainers(ctx context.Context, all bool) error {
	uc.logger.Info("Listing Docker containers...")

	args := []string{"ps"}
	if all {
		args = append(args, "-a")
	}

	task := domain.NewTask("docker-ps", domain.TaskTypeDocker, "docker", args)

	result, err := uc.executor.ExecuteWithOutput(ctx, task)
	if err != nil {
		return fmt.Errorf("docker ps failed: %w", err)
	}

	if result.Success {
		fmt.Println(result.Output)
	}

	return nil
}

func (uc *DockerUseCase) createBuildTask(opts *DockerBuildOptions, dockerfile string) *domain.Task {
	args := []string{"build"}

	if opts.NoCache {
		args = append(args, "--no-cache")
	}

	if opts.Platform != "" {
		args = append(args, "--platform", opts.Platform)
	}

	// Add build args
	for key, val := range opts.BuildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", key, val))
	}

	args = append(args, "-f", dockerfile)

	tag := opts.Tag
	if tag == "" {
		tag = "latest"
	}
	args = append(args, "-t", fmt.Sprintf("%s:%s", opts.ImageName, tag))

	context := opts.Context
	if context == "" {
		context = "."
	}
	args = append(args, context)

	return domain.NewTask("docker-build", domain.TaskTypeDocker, "docker", args)
}
