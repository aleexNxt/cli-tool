package domain

import "context"

// CommandExecutor definiert das Interface für Command-Ausführung
type CommandExecutor interface {
	Execute(ctx context.Context, task *Task) (*TaskResult, error)
	ExecuteWithOutput(ctx context.Context, task *Task) (*TaskResult, error)
}

// FileSystem definiert das Interface für Dateisystem-Operationen
type FileSystem interface {
	Exists(path string) bool
	IsDirectory(path string) bool
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	CreateDirectory(path string) error
	ListFiles(path string) ([]string, error)
	GetWorkingDirectory() (string, error)
}

// ConfigRepository definiert das Interface für Konfigurationsverwaltung
type ConfigRepository interface {
	Load() (*Config, error)
	Save(config *Config) error
	Get(key string) (string, error)
	Set(key, value string) error
}

// Logger definiert das Interface für Logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Success(msg string, args ...interface{})
}

// Config repräsentiert die Anwendungskonfiguration
type Config struct {
	BuildCommand   string            `json:"build_command,omitempty"`
	TestCommand    string            `json:"test_command,omitempty"`
	DeployCommand  string            `json:"deploy_command,omitempty"`
	LintCommand    string            `json:"lint_command,omitempty"`
	DockerRegistry string            `json:"docker_registry,omitempty"`
	Environment    map[string]string `json:"environment,omitempty"`
	Timeout        int               `json:"timeout,omitempty"` // in seconds
}

// NewConfig erstellt eine neue Standardkonfiguration
func NewConfig() *Config {
	return &Config{
		BuildCommand:  "go build -v ./...",
		TestCommand:   "go test -v ./...",
		LintCommand:   "golangci-lint run",
		Environment:   make(map[string]string),
		Timeout:       300, // 5 minutes default
	}
}
