package logger

import (
	"fmt"
	"time"
)

// Color codes für Terminal-Ausgabe
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

// Logger implementiert domain.Logger
type Logger struct {
	verbose bool
}

// NewLogger erstellt einen neuen Logger
func NewLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
	}
}

// Info loggt eine Info-Nachricht
func (l *Logger) Info(msg string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s[%s] %sℹ%s  %s\n", colorGray, timestamp, colorBlue, colorReset, formatted)
}

// Error loggt eine Fehler-Nachricht
func (l *Logger) Error(msg string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s[%s] %s✗%s  %s\n", colorGray, timestamp, colorRed, colorReset, formatted)
}

// Debug loggt eine Debug-Nachricht
func (l *Logger) Debug(msg string, args ...interface{}) {
	if !l.verbose {
		return
	}
	timestamp := time.Now().Format("15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s[%s] 🔍  %s%s\n", colorGray, timestamp, formatted, colorReset)
}

// Warn loggt eine Warn-Nachricht
func (l *Logger) Warn(msg string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s[%s] %s⚠%s  %s\n", colorGray, timestamp, colorYellow, colorReset, formatted)
}

// Success loggt eine Erfolgs-Nachricht
func (l *Logger) Success(msg string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s[%s] %s✓%s  %s\n", colorGray, timestamp, colorGreen, colorReset, formatted)
}
