# DevTool - CLI Automation Tool

Ein professionelles Kommandozeilen-Tool in Go zur Automatisierung wiederkehrender Entwicklungs- und Deployment-Aufgaben, strukturiert nach Clean Code und Clean Architecture Prinzipien.

## 🚀 Features

- **Build-Automation**: Flexible Build-Konfigurationen mit Cross-Compilation Support
- **Test-Automation**: Umfassende Test-Suite mit Coverage-Reports und Benchmarks
- **Deployment**: Automatisierte Deployments mit Rollback-Funktionalität
- **Docker-Integration**: Vollständige Docker-Lifecycle-Verwaltung
- **Konfigurationsmanagement**: Zentrale, persistente Konfiguration
- **Clean Architecture**: Modularer, testbarer und wartbarer Code

## 📋 Inhaltsverzeichnis

- [Installation](#installation)
- [Architektur](#architektur)
- [Verwendung](#verwendung)
  - [Build](#build)
  - [Test](#test)
  - [Deploy](#deploy)
  - [Docker](#docker)
  - [Konfiguration](#konfiguration)
- [Entwicklung](#entwicklung)
- [Beispiele](#beispiele)

## 🔧 Installation

### Aus dem Quellcode

```bash
# Repository klonen
git clone https://github.com/aleexNxt/cli-tool.git
cd cli-tool

# Dependencies installieren
go mod download

# Binary bauen
make build

# Optional: Binary installieren
make install
```

### Mit Go Install

```bash
go install github.com/aleexNxt/cli-tool/cmd/devtool@latest
```

## 🏗️ Architektur

Das Projekt folgt den Prinzipien der **Clean Architecture** mit klarer Trennung der Schichten:

```
cli-tool/
├── cmd/
│   └── devtool/              # Entry Point
│       └── main.go
├── internal/
│   ├── domain/               # Business Entities & Interfaces
│   │   ├── task.go          # Core Domain Models
│   │   └── repository.go    # Repository Interfaces
│   ├── usecase/              # Business Logic
│   │   ├── build.go         # Build-Operationen
│   │   ├── test.go          # Test-Operationen
│   │   ├── deploy.go        # Deployment-Operationen
│   │   └── docker.go        # Docker-Operationen
│   ├── infrastructure/       # External Dependencies
│   │   ├── executor/        # Command Execution
│   │   ├── filesystem/      # File System Operations
│   │   └── config/          # Configuration Management
│   └── interface/            # Interface Adapters
│       └── cli/             # CLI Commands (Cobra)
├── pkg/                      # Public Libraries
│   └── logger/              # Logging Utilities
└── Makefile                 # Build Automation
```

### Clean Architecture Prinzipien

1. **Domain Layer**: Enthält die Kerngeschäftslogik ohne externe Dependencies
2. **Use Case Layer**: Implementiert Anwendungsfälle mit Dependency Injection
3. **Infrastructure Layer**: Implementiert Interfaces für externe Systeme
4. **Interface Layer**: CLI-Adapter mit Cobra Framework

## 📖 Verwendung

### Build

Build-Operationen für Go-Projekte:

```bash
# Einfacher Build
devtool build

# Build mit Optionen
devtool build --output bin/myapp --race --verbose

# Cross-Compilation
devtool build --goos linux --goarch amd64

# Build mit Custom Flags
devtool build --ldflags "-X main.version=1.0.0"

# Build mit Tags
devtool build --tags "integration,e2e"

# Alle Packages bauen
devtool build all

# Build-Artefakte entfernen
devtool build clean
```

**Optionen:**
- `--output, -o`: Output-Verzeichnis oder -Datei
- `--race`: Race-Detector aktivieren
- `--tags, -t`: Build-Tags (kommagetrennt)
- `--goos`: Ziel-Betriebssystem
- `--goarch`: Ziel-Architektur
- `--ldflags, -l`: Linker-Flags
- `--cgo`: CGO aktivieren
- `--clean`: Vor dem Build aufräumen
- `--verbose, -v`: Ausführlicher Output

### Test

Test-Automation mit Coverage und Benchmarks:

```bash
# Tests ausführen
devtool test

# Tests mit Coverage
devtool test --coverage

# Tests mit Race-Detector
devtool test --race

# Spezifisches Package testen
devtool test --package ./internal/usecase

# Spezifische Tests ausführen
devtool test --run TestBuild

# Tests mit Optionen
devtool test --short --timeout 5m --parallel 4

# Coverage-Report generieren
devtool test coverage --output coverage.out

# Benchmarks ausführen
devtool test bench

# Go vet ausführen
devtool test vet
```

**Optionen:**
- `--coverage, -c`: Coverage-Report erstellen
- `--race, -r`: Race-Detector aktivieren
- `--short, -s`: Kurze Tests ausführen
- `--package, -p`: Spezifisches Package
- `--run`: Test-Muster
- `--count`: Anzahl der Durchläufe
- `--parallel`: Parallele Ausführung
- `--timeout, -t`: Timeout für Tests

### Deploy

Deployment-Automation:

```bash
# Deployment ausführen
devtool deploy --env production

# Dry-Run (Simulation)
devtool deploy --env staging --dry-run

# Deployment mit Tag
devtool deploy --env production --tag v1.0.0

# Rollback
devtool deploy rollback --env production --version v0.9.0

# Deployment-Status prüfen
devtool deploy status --env production
```

**Optionen:**
- `--env, -e`: Ziel-Umgebung (default: production)
- `--dry-run`: Simulation ohne Deployment
- `--force`: Erzwinge Deployment
- `--tag, -t`: Version/Tag für Deployment
- `--config-file, -f`: Deployment-Konfigurationsdatei

### Docker

Docker-Lifecycle-Management:

```bash
# Docker Image bauen
devtool docker build --image myapp --tag v1.0.0

# Image mit Build-Args
devtool docker build --image myapp --build-arg VERSION=1.0.0

# Multi-Platform Build
devtool docker build --image myapp --platform linux/amd64,linux/arm64

# Image zur Registry pushen
devtool docker push --image myapp --tag v1.0.0

# Container starten
devtool docker run --image myapp:v1.0.0 --name myapp-container \
  --port 8080:80 --env ENV=production

# Container stoppen
devtool docker stop --name myapp-container

# Container entfernen
devtool docker rm --name myapp-container --force

# Images auflisten
devtool docker images

# Container auflisten
devtool docker ps --all
```

**Build-Optionen:**
- `--image, -i`: Image-Name (erforderlich)
- `--tag, -t`: Image-Tag (default: latest)
- `--file, -f`: Pfad zum Dockerfile
- `--context, -C`: Build-Context
- `--no-cache`: Cache nicht verwenden
- `--platform`: Ziel-Platform
- `--build-arg`: Build-Argumente

### Konfiguration

Zentrale Konfigurationsverwaltung:

```bash
# Konfiguration initialisieren
devtool config init

# Aktuelle Konfiguration anzeigen
devtool config show

# Wert setzen
devtool config set build_command "go build -v ./..."
devtool config set docker_registry "registry.example.com"

# Wert abrufen
devtool config get build_command
```

**Konfigurationsdatei**: `~/.devtool.json`

**Verfügbare Keys:**
- `build_command`: Standard-Build-Command
- `test_command`: Standard-Test-Command
- `deploy_command`: Standard-Deploy-Command
- `lint_command`: Standard-Lint-Command
- `docker_registry`: Docker Registry URL
- `timeout`: Standard-Timeout in Sekunden

**Beispiel-Konfiguration:**
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

## 🛠️ Entwicklung

### Voraussetzungen

- Go 1.21 oder höher
- Make (optional, für Makefile-Targets)

### Projekt bauen

```bash
# Dependencies installieren
make deps

# Build
make build

# Tests ausführen
make test

# Code formatieren
make fmt

# Linting
make vet

# Alles bauen (Clean + Build)
make all
```

### Makefile-Targets

- `make build` - Binary bauen
- `make build-all` - Für mehrere Plattformen bauen
- `make test` - Tests ausführen
- `make test-coverage` - Tests mit Coverage
- `make bench` - Benchmarks ausführen
- `make vet` - Go vet ausführen
- `make fmt` - Code formatieren
- `make clean` - Build-Artefakte entfernen
- `make install` - Binary installieren
- `make run` - Bauen und ausführen
- `make deps` - Dependencies aktualisieren

### Projektstruktur erweitern

#### Neuen Use Case hinzufügen

1. **Domain Interface** in `internal/domain/repository.go` definieren
2. **Use Case** in `internal/usecase/` implementieren
3. **CLI Command** in `internal/interface/cli/` hinzufügen
4. Command in `internal/interface/cli/root.go` registrieren

#### Beispiel: Neuen "Lint" Use Case hinzufügen

```go
// internal/usecase/lint.go
package usecase

type LintUseCase struct {
    executor domain.CommandExecutor
    logger   domain.Logger
}

func (uc *LintUseCase) Execute(ctx context.Context) error {
    // Implementation
}

// internal/interface/cli/lint.go
func newLintCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "lint",
        Short: "Führt Code-Linting aus",
        RunE: func(cmd *cobra.Command, args []string) error {
            deps, _ := InitDependencies()
            return deps.LintUseCase.Execute(context.Background())
        },
    }
}
```

## 📝 Beispiele

### Beispiel 1: Build-Pipeline

```bash
# Projekt bauen und testen
devtool build --clean --verbose
devtool test --coverage --race
devtool test vet
```

### Beispiel 2: Docker-Workflow

```bash
# Image bauen und pushen
devtool docker build --image myapp --tag v1.0.0 --no-cache
devtool docker push --image myapp --tag v1.0.0

# Container starten
devtool docker run --image myapp:v1.0.0 --name myapp \
  --port 8080:80 \
  --env DATABASE_URL=postgres://localhost/mydb
```

### Beispiel 3: Deployment-Pipeline

```bash
# Tests ausführen
devtool test --race

# Build für Production
devtool build --goos linux --goarch amd64 --ldflags "-s -w"

# Docker Image erstellen
devtool docker build --image myapp --tag $(git rev-parse --short HEAD)

# Deploy mit Dry-Run
devtool deploy --env production --dry-run --tag $(git rev-parse --short HEAD)

# Tatsächliches Deployment
devtool deploy --env production --tag $(git rev-parse --short HEAD)
```

### Beispiel 4: CI/CD Integration

```yaml
# .github/workflows/ci.yml
name: CI
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Install DevTool
        run: go install github.com/aleexNxt/cli-tool/cmd/devtool@latest
      - name: Run Tests
        run: devtool test --coverage
      - name: Build
        run: devtool build
```

## 🧪 Testing

```bash
# Unit Tests
go test ./...

# Mit Coverage
go test -cover ./...

# Benchmarks
go test -bench=. ./...

# Race-Detector
go test -race ./...
```

## 📄 Lizenz

Dieses Projekt ist unter der MIT-Lizenz lizenziert - siehe [LICENSE](LICENSE) Datei für Details.

## 🤝 Beitragen

Contributions sind willkommen! Bitte öffne ein Issue oder Pull Request.

1. Fork das Projekt
2. Erstelle einen Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit deine Änderungen (`git commit -m 'Add some AmazingFeature'`)
4. Push zum Branch (`git push origin feature/AmazingFeature`)
5. Öffne einen Pull Request

## 📞 Support

Bei Fragen oder Problemen öffne bitte ein Issue auf GitHub.

## 🙏 Danksagungen

- [Cobra](https://github.com/spf13/cobra) - Leistungsstarkes CLI-Framework
- Clean Architecture Prinzipien von Robert C. Martin
- Go Community für ausgezeichnete Tools und Libraries
