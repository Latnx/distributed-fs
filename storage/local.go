package storage

import (
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BaseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{BaseDir: baseDir}
}

func (s *LocalStorage) WriteFile(filename string, data []byte) error {
	filePath := filepath.Join(s.BaseDir, filename)
	return os.WriteFile(filePath, data, 0644) // Replaced ioutil.WriteFile with os.WriteFile
}

func (s *LocalStorage) ReadFile(filename string) ([]byte, error) {
	filePath := filepath.Join(s.BaseDir, filename)
	return os.ReadFile(filePath) // Replaced ioutil.ReadFile with os.ReadFile
}

func (s *LocalStorage) DeleteFile(filename string) error {
	filePath := filepath.Join(s.BaseDir, filename)
	return os.Remove(filePath)
}

func (s *LocalStorage) ListFiles() ([]string, error) {
	files := []string{}
	err := filepath.Walk(s.BaseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})
	return files, err
}
