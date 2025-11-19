package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileSystem implementiert domain.FileSystem
type FileSystem struct{}

// NewFileSystem erstellt ein neues FileSystem
func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

// Exists prüft ob ein Pfad existiert
func (fs *FileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory prüft ob ein Pfad ein Verzeichnis ist
func (fs *FileSystem) IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ReadFile liest eine Datei
func (fs *FileSystem) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// WriteFile schreibt eine Datei
func (fs *FileSystem) WriteFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}

// CreateDirectory erstellt ein Verzeichnis
func (fs *FileSystem) CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// ListFiles listet Dateien in einem Verzeichnis
func (fs *FileSystem) ListFiles(path string) ([]string, error) {
	var files []string

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		files = append(files, filepath.Join(path, entry.Name()))
	}

	return files, nil
}

// GetWorkingDirectory gibt das aktuelle Arbeitsverzeichnis zurück
func (fs *FileSystem) GetWorkingDirectory() (string, error) {
	return os.Getwd()
}
