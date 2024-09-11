package util

import (
	"os"
	"path/filepath"
	"strings"
)

func IsRunningInsideDocker() bool {
	_, err := os.Stat(filepath.Join("/", ".dockerenv"))
	if err == nil {
		return true
	}

	return isRunningInsideDockerCGroup()
}

func isRunningInsideDockerCGroup() bool {
	data, err := os.ReadFile(filepath.Join("/", "proc", "self", "cgroup"))
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
