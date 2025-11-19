package domain

import (
	"time"
)

// TaskType definiert den Typ einer Aufgabe
type TaskType string

const (
	TaskTypeBuild  TaskType = "build"
	TaskTypeTest   TaskType = "test"
	TaskTypeDeploy TaskType = "deploy"
	TaskTypeLint   TaskType = "lint"
	TaskTypeDocker TaskType = "docker"
)

// TaskStatus repräsentiert den Status einer Aufgabe
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

// Task repräsentiert eine ausführbare Aufgabe
type Task struct {
	ID          string
	Name        string
	Type        TaskType
	Status      TaskStatus
	Command     string
	Args        []string
	WorkDir     string
	Environment map[string]string
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	Error       error
	Output      string
}

// TaskResult repräsentiert das Ergebnis einer Aufgabe
type TaskResult struct {
	Success  bool
	Output   string
	Error    error
	Duration time.Duration
}

// ExecutionContext enthält Kontext-Informationen für die Ausführung
type ExecutionContext struct {
	WorkDir     string
	Environment map[string]string
	Timeout     time.Duration
	Verbose     bool
}

// NewTask erstellt eine neue Task
func NewTask(name string, taskType TaskType, command string, args []string) *Task {
	return &Task{
		ID:          generateTaskID(),
		Name:        name,
		Type:        taskType,
		Status:      TaskStatusPending,
		Command:     command,
		Args:        args,
		Environment: make(map[string]string),
		CreatedAt:   time.Now(),
	}
}

// Start markiert den Beginn der Aufgabe
func (t *Task) Start() {
	now := time.Now()
	t.Status = TaskStatusRunning
	t.StartedAt = &now
}

// Complete markiert die erfolgreiche Fertigstellung
func (t *Task) Complete(output string) {
	now := time.Now()
	t.Status = TaskStatusCompleted
	t.CompletedAt = &now
	t.Output = output
}

// Fail markiert das Fehlschlagen der Aufgabe
func (t *Task) Fail(err error, output string) {
	now := time.Now()
	t.Status = TaskStatusFailed
	t.CompletedAt = &now
	t.Error = err
	t.Output = output
}

// Duration berechnet die Ausführungsdauer
func (t *Task) Duration() time.Duration {
	if t.StartedAt == nil {
		return 0
	}
	if t.CompletedAt == nil {
		return time.Since(*t.StartedAt)
	}
	return t.CompletedAt.Sub(*t.StartedAt)
}

func generateTaskID() string {
	return time.Now().Format("20060102150405")
}
