# DevTool (dt)

CLI-Tool zur Automatisierung von Entwicklungs- und Deployment-Aufgaben in Go nach Clean Architecture Prinzipien.

## Features

- **Build**: Cross-Compilation, Race-Detector, Custom Flags
- **Test**: Coverage, Benchmarks, go vet
- **Lint**: golangci-lint, gofmt, goimports
- **Deploy**: Multi-Environment, Rollback, Dry-Run
- **Docker**: Build, Push, Container Management
- **Config**: Persistente Konfiguration

## Installation

```bash
# Clone & Build
git clone https://github.com/aleexNxt/cli-tool.git
cd cli-tool
make build

# Optional: Installieren
make install
```

## Verwendung

```bash
# Beide Commands verfügbar
devtool <command>
dt <command>

# Build
dt build --race --verbose
dt build --goos linux --goarch amd64

# Test
dt test --coverage
dt test bench
dt test vet

# Lint
dt lint
dt lint --fix
dt lint fmt

# Docker
dt docker build --image myapp --tag v1.0.0
dt docker push --image myapp --tag v1.0.0
dt docker run --image myapp:latest --name app --port 8080:80

# Deploy
dt deploy --env production --dry-run
dt deploy --env staging --tag v1.0.0
dt deploy rollback --env production --version v0.9.0

# Config
dt config init
dt config show
dt config set build_command "go build -v ./..."
```

## Architektur

```
cli-tool/
├── cmd/
│   ├── devtool/          # Main binary
│   └── dt/               # Alias binary
├── internal/
│   ├── domain/           # Entities & Interfaces
│   ├── usecase/          # Business Logic
│   ├── infrastructure/   # Executor, FileSystem, Config
│   └── interface/cli/    # CLI Commands
└── pkg/logger/           # Utilities
```

**Clean Architecture Layers:**
- **Domain**: Core entities ohne Dependencies
- **UseCase**: Business Logic (Build, Test, Lint, Deploy, Docker)
- **Infrastructure**: Command Executor, FileSystem, Config
- **Interface**: CLI mit Cobra Framework


## Makefile Targets

```bash
make build         # Build beide Binaries (devtool & dt)
make build-all     # Multi-Platform Build
make test          # Tests ausführen
make test-coverage # Coverage Report
make clean         # Aufräumen
make install       # Binaries installieren
```

## Konfiguration

`~/.devtool.json`:
```json
{
  "build_command": "go build -v ./...",
  "test_command": "go test -v ./...",
  "lint_command": "golangci-lint run",
  "deploy_command": "./scripts/deploy.sh",
  "docker_registry": "registry.example.com",
  "environment": {
    "GO111MODULE": "on",
    "CGO_ENABLED": "0"
  },
  "timeout": 300
}
```

## Beispiele

### CI/CD Pipeline
```bash
dt lint
dt test --race --coverage
dt build --goos linux --goarch amd64
dt docker build --image myapp --tag $(git rev-parse --short HEAD)
dt deploy --env production --tag $(git rev-parse --short HEAD)
```

### Docker Workflow
```bash
dt docker build --image myapp --tag v1.0.0 --no-cache
dt docker push --image myapp --tag v1.0.0
dt docker run --image myapp:v1.0.0 --name myapp \
  --port 8080:80 --env DATABASE_URL=postgres://db/app
```

## Lizenz

MIT License - siehe [LICENSE](LICENSE)
