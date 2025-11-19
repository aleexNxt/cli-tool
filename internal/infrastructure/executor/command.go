package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/aleexNxt/cli-tool/internal/domain"
)

// CommandExecutor implementiert domain.CommandExecutor
type CommandExecutor struct {
	workDir string
	timeout time.Duration
}

// NewCommandExecutor erstellt einen neuen CommandExecutor
func NewCommandExecutor(workDir string, timeout time.Duration) *CommandExecutor {
	if timeout == 0 {
		timeout = 5 * time.Minute
	}
	return &CommandExecutor{
		workDir: workDir,
		timeout: timeout,
	}
}

// Execute führt einen Command aus ohne Output zu streamen
func (e *CommandExecutor) Execute(ctx context.Context, task *domain.Task) (*domain.TaskResult, error) {
	task.Start()

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(timeoutCtx, task.Command, task.Args...)

	// Set working directory
	if task.WorkDir != "" {
		cmd.Dir = task.WorkDir
	} else if e.workDir != "" {
		cmd.Dir = e.workDir
	}

	// Set environment variables
	if len(task.Environment) > 0 {
		cmd.Env = append(cmd.Env, e.buildEnvVars(task.Environment)...)
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	// Combine output
	output := stdout.String()
	if stderr.Len() > 0 {
		if len(output) > 0 {
			output += "\n"
		}
		output += stderr.String()
	}

	// Create result
	result := &domain.TaskResult{
		Success:  err == nil,
		Output:   output,
		Error:    err,
		Duration: duration,
	}

	// Update task
	if err != nil {
		task.Fail(err, output)
	} else {
		task.Complete(output)
	}

	return result, err
}

// ExecuteWithOutput führt einen Command aus und streamt Output
func (e *CommandExecutor) ExecuteWithOutput(ctx context.Context, task *domain.Task) (*domain.TaskResult, error) {
	task.Start()

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(timeoutCtx, task.Command, task.Args...)

	// Set working directory
	if task.WorkDir != "" {
		cmd.Dir = task.WorkDir
	} else if e.workDir != "" {
		cmd.Dir = e.workDir
	}

	// Set environment variables
	if len(task.Environment) > 0 {
		cmd.Env = append(cmd.Env, e.buildEnvVars(task.Environment)...)
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	// Combine output
	output := stdout.String()
	if stderr.Len() > 0 {
		if len(output) > 0 {
			output += "\n"
		}
		output += stderr.String()
	}

	// Create result
	result := &domain.TaskResult{
		Success:  err == nil,
		Output:   output,
		Error:    err,
		Duration: duration,
	}

	// Update task
	if err != nil {
		task.Fail(err, output)
	} else {
		task.Complete(output)
	}

	return result, err
}

func (e *CommandExecutor) buildEnvVars(envMap map[string]string) []string {
	var envVars []string
	for key, val := range envMap {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, val))
	}
	return envVars
}
