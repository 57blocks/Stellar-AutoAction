package util

import (
	"os"
	"path/filepath"
	"strings"
)

type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
}

type RealFileSystem struct{}

func (RealFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (RealFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func IsRunningInsideDocker(fs FileSystem) bool {
	_, err := fs.Stat(filepath.Join("/", ".dockerenv"))
	if err == nil {
		return true
	}

	return isRunningInsideDockerCGroup(fs)
}

func isRunningInsideDockerCGroup(fs FileSystem) bool {
	data, err := fs.ReadFile(filepath.Join("/", "proc", "self", "cgroup"))
	if err != nil {
		return false
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "/")
		for _, part := range parts {
			if part == "docker" {
				return true
			}
		}
	}

	return false
}
